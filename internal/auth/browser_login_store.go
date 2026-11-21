package auth

import (
 "context"
 "crypto/sha256"
 "database/sql"
 "errors"
 "fmt"
 "time"
)

var ErrLoginTransactionUnavailable=errors.New("login transaction store unavailable")
var ErrLoginTransactionExpired=errors.New("login transaction invalid or expired")

type BrowserLoginStore struct{DB *sql.DB}

func (s BrowserLoginStore) Create(ctx context.Context,tx BrowserLoginTransaction,binding string,ttl time.Duration)error{
 if s.DB==nil{return ErrLoginTransactionUnavailable}
 if ctx==nil||tx.State==""||tx.Nonce==""||len(tx.PKCEVerifier)<43||binding==""||ttl<time.Minute||ttl>10*time.Minute{return ErrInvalidLoginTransaction}
 state:=sha256.Sum256([]byte(tx.State))
 bound:=sha256.Sum256([]byte(binding))
 _,err:=s.DB.ExecContext(ctx,`INSERT INTO oidc_login_transactions(state_sha256,browser_binding_sha256,nonce,pkce_verifier,expires_at)
 VALUES($1,$2,$3,$4,now()+($5 * interval '1 second'))`,state[:],bound[:],tx.Nonce,tx.PKCEVerifier,int64(ttl.Seconds()))
 if err!=nil{return fmt.Errorf("create login transaction: %w",err)}
 return nil
}

// Consume atomically marks a transaction as used and returns secrets only
// when both state and browser binding match and the transaction is unexpired.
func (s BrowserLoginStore) Consume(ctx context.Context,state,binding string)(BrowserLoginTransaction,error){
 if s.DB==nil{return BrowserLoginTransaction{},ErrLoginTransactionUnavailable}
 if ctx==nil||state==""||binding==""{return BrowserLoginTransaction{},ErrInvalidLoginTransaction}
 stateHash:=sha256.Sum256([]byte(state))
 bindingHash:=sha256.Sum256([]byte(binding))
 var tx BrowserLoginTransaction
 err:=s.DB.QueryRowContext(ctx,`UPDATE oidc_login_transactions SET consumed_at=now()
 WHERE state_sha256=$1 AND browser_binding_sha256=$2 AND consumed_at IS NULL AND expires_at>now()
 RETURNING nonce,pkce_verifier`,stateHash[:],bindingHash[:]).Scan(&tx.Nonce,&tx.PKCEVerifier)
 if errors.Is(err,sql.ErrNoRows){return BrowserLoginTransaction{},ErrLoginTransactionExpired}
 if err!=nil{return BrowserLoginTransaction{},fmt.Errorf("consume login transaction: %w",err)}
 tx.State=state
 return tx,nil
}
