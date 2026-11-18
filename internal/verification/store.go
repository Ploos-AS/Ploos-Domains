package verification

import (
 "context"
 "database/sql"
 "errors"
 "fmt"
 "time"
)

var (
 ErrChallengeUnavailable = errors.New("verification challenge missing, expired or already consumed")
 ErrDomainAlreadyVerified = errors.New("domain already verified")
)

// CompleteVerification commits an already validated DNS proof to PostgreSQL.
// The caller MUST first authenticate and authorize the request, verify the
// challenge TXT record using Verify, and obtain domainID/challengeID from
// trusted server-side state. This method does not perform those steps.
//
// The challenge is locked and consumed exactly once, and the domain is
// marked verified in the same transaction. No external network calls occur
// while holding the transaction.
func CompleteVerification(ctx context.Context, db *sql.DB, domainID, challengeID string, now time.Time) (err error) {
 if db == nil || domainID == "" || challengeID == "" { return ErrInvalid }
 tx, err := db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelReadCommitted})
 if err != nil { return fmt.Errorf("begin verification transaction: %w", err) }
 defer func(){ _ = tx.Rollback() }()

 var existing sql.NullTime
 if err = tx.QueryRowContext(ctx,
  `SELECT verified_at FROM domains WHERE id=$1 FOR UPDATE`, domainID,
 ).Scan(&existing); err != nil {
  if errors.Is(err,sql.ErrNoRows) { return ErrChallengeUnavailable }
  return fmt.Errorf("lock domain: %w",err)
 }
 if existing.Valid { return ErrDomainAlreadyVerified }

 // Lock the challenge and check its exact domain and expiry. The
 // caller must have validated the associated TXT digest beforehand.
 var expiresAt time.Time
 var consumed sql.NullTime
 if err = tx.QueryRowContext(ctx,
  `SELECT expires_at,consumed_at FROM domain_verification_challenges
    WHERE id=$1 AND domain_id=$2 FOR UPDATE`, challengeID, domainID,
 ).Scan(&expiresAt,&consumed); err != nil {
  if errors.Is(err,sql.ErrNoRows) { return ErrChallengeUnavailable }
  return fmt.Errorf("lock challenge: %w",err)
 }
 if consumed.Valid || !now.Before(expiresAt) { return ErrChallengeUnavailable }

 result,err:=tx.ExecContext(ctx,
  `UPDATE domain_verification_challenges SET consumed_at=$3
    WHERE id=$1 AND domain_id=$2 AND consumed_at IS NULL AND expires_at>$3`,
  challengeID,domainID,now)
 if err!=nil { return fmt.Errorf("consume challenge: %w",err) }
 affected,err:=result.RowsAffected()
 if err!=nil { return fmt.Errorf("challenge rows affected: %w",err) }
 if affected!=1 { return ErrChallengeUnavailable }

 result,err=tx.ExecContext(ctx,
  `UPDATE domains SET verified_at=$2 WHERE id=$1 AND verified_at IS NULL`,
  domainID,now)
 if err!=nil { return fmt.Errorf("verify domain: %w",err) }
 affected,err=result.RowsAffected()
 if err!=nil { return fmt.Errorf("domain rows affected: %w",err) }
 if affected!=1 { return ErrDomainAlreadyVerified }

 if err=tx.Commit();err!=nil { return fmt.Errorf("commit verification: %w",err) }
 return nil
}
