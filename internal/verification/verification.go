// Package verification implements DNS TXT proof-of-control primitives.
// It does not grant permissions or update database verification state.
package verification

import (
 "context"
 "crypto/rand"
 "crypto/sha256"
 "crypto/subtle"
 "encoding/hex"
 "errors"
 "fmt"
 "net"
 "strings"
 "time"

 "github.com/Ploos-AS/Ploos-Domains/internal/domains"
)

const Prefix = "_ploos-verify"
var (
 ErrExpired = errors.New("verification challenge expired")
 ErrMismatch = errors.New("verification TXT record not found")
 ErrInvalid = errors.New("invalid verification challenge")
)

type TXTResolver interface {
 LookupTXT(context.Context, string) ([]string,error)
}

type Challenge struct {
 Domain string
 RecordName string
 Token string
 TokenHash [32]byte
 ExpiresAt time.Time
}

// NewChallenge generates a random 256-bit token; only TokenHash should be
// persisted. Token must be returned to the requester once over HTTPS.
func NewChallenge(domain string, now time.Time, ttl time.Duration) (Challenge,error) {
 fqdn,err:=domains.CanonicalName(domain)
 if err!=nil || ttl<=0 || ttl>24*time.Hour { return Challenge{},ErrInvalid }
 raw:=make([]byte,32)
 if _,err=rand.Read(raw);err!=nil { return Challenge{},fmt.Errorf("generate token: %w",err) }
 token:="ploos-domain-verification="+hex.EncodeToString(raw)
 return Challenge{Domain:fqdn,RecordName:Prefix+"."+fqdn,Token:token,TokenHash:sha256.Sum256([]byte(token)),ExpiresAt:now.Add(ttl)},nil
}

// Verify checks a fresh DNS TXT lookup against the stored token digest.
// Never accept an expired challenge. A successful result must be committed
// atomically with challenge consumption by the persistence layer.
func Verify(ctx context.Context, resolver TXTResolver, name string, expected [32]byte, expiresAt,timeNow time.Time) error {
 if resolver==nil || !timeNow.Before(expiresAt) { return ErrExpired }
 canonical,err:=domains.CanonicalName(name)
 if err!=nil || !strings.HasPrefix(canonical,Prefix+".") { return ErrInvalid }
 records,err:=resolver.LookupTXT(ctx,canonical)
 if err!=nil { return fmt.Errorf("lookup TXT: %w",err) }
 for _,record:=range records {
  digest:=sha256.Sum256([]byte(record))
  if subtle.ConstantTimeCompare(digest[:],expected[:])==1 { return nil }
 }
 return ErrMismatch
}

var _ TXTResolver = (*net.Resolver)(nil)
