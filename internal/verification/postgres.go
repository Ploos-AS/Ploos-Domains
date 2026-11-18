package verification

import (
 "context"
 "database/sql"
 "errors"
 "fmt"
 "time"
)

// PostgresRepository implements ChallengeRepository using database/sql.
// The application must supply a PostgreSQL-compatible *sql.DB; a driver and
// connection configuration are intentionally not selected in this package.
type PostgresRepository struct { DB *sql.DB }

func (r PostgresRepository) LoadPending(ctx context.Context, domainID, challengeID string) (PendingChallenge,error) {
 if r.DB==nil || domainID=="" || challengeID=="" { return PendingChallenge{},ErrInvalid }
 var c PendingChallenge
 var hash []byte
 var consumed sql.NullTime
 err:=r.DB.QueryRowContext(ctx,`
 SELECT c.id::text,c.domain_id::text,d.organization_id::text,c.record_name,
        c.token_sha256,c.expires_at,c.consumed_at
 FROM domain_verification_challenges c
 JOIN domains d ON d.id=c.domain_id
 WHERE c.id=$1::uuid AND c.domain_id=$2::uuid
 `,challengeID,domainID).Scan(
  &c.ID,&c.DomainID,&c.OrganizationID,&c.RecordName,&hash,&c.ExpiresAt,&consumed,
 )
 if errors.Is(err,sql.ErrNoRows) { return PendingChallenge{},ErrChallengeUnavailable }
 if err!=nil { return PendingChallenge{},fmt.Errorf("load challenge: %w",err) }
 if len(hash)!=32 { return PendingChallenge{},ErrInvalid }
 copy(c.TokenHash[:],hash)
 c.Consumed=consumed.Valid
 return c,nil
}

func (r PostgresRepository) Complete(ctx context.Context,domainID,challengeID string,now time.Time) error {
 return CompleteVerification(ctx,r.DB,domainID,challengeID,now)
}
