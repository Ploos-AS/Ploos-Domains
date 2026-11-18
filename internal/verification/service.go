package verification

import (
 "context"
 "errors"
 "fmt"
 "time"
)

var ErrUnauthorized = errors.New("not authorized to verify domain")

// PendingChallenge is loaded from trusted persistence, never from HTTP input.
type PendingChallenge struct {
 ID string
 DomainID string
 OrganizationID string
 RecordName string
 TokenHash [32]byte
 ExpiresAt time.Time
 Consumed bool
}

type ChallengeRepository interface {
 LoadPending(ctx context.Context, domainID, challengeID string) (PendingChallenge,error)
 Complete(ctx context.Context, domainID, challengeID string, now time.Time) error
}

type DomainAuthorizer interface {
 CanVerifyDomain(ctx context.Context, actorID, organizationID, domainID string) (bool,error)
}

// Service coordinates authorization, fresh DNS proof and atomic completion.
// Concrete adapters and authenticated HTTP handlers will be added separately.
type Service struct {
 Repository ChallengeRepository
 Authorizer DomainAuthorizer
 Resolver TXTResolver
 Clock func() time.Time
}

func (s Service) VerifyDomain(ctx context.Context, actorID, organizationID, domainID, challengeID string) error {
 if actorID=="" || organizationID=="" || domainID=="" || challengeID=="" ||
    s.Repository==nil || s.Authorizer==nil || s.Resolver==nil || s.Clock==nil {
  return ErrInvalid
 }
 allowed,err:=s.Authorizer.CanVerifyDomain(ctx,actorID,organizationID,domainID)
 if err!=nil { return fmt.Errorf("authorize domain verification: %w",err) }
 if !allowed { return ErrUnauthorized }

 challenge,err:=s.Repository.LoadPending(ctx,domainID,challengeID)
 if err!=nil { return fmt.Errorf("load verification challenge: %w",err) }
 if challenge.ID!=challengeID || challenge.DomainID!=domainID ||
    challenge.OrganizationID!=organizationID || challenge.Consumed {
  return ErrChallengeUnavailable
 }
 now:=s.Clock()
 if err:=Verify(ctx,s.Resolver,challenge.RecordName,challenge.TokenHash,challenge.ExpiresAt,now);err!=nil {
  return err
 }
 // Repository.Complete must recheck domain ownership, challenge state and
 // expiry transactionally. The authorization decision is not a substitute
 // for those checks.
 if err:=s.Repository.Complete(ctx,domainID,challengeID,s.Clock());err!=nil {
  return fmt.Errorf("complete verification: %w",err)
 }
 return nil
}
