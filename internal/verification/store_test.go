package verification

import (
 "context"
 "errors"
 "testing"
 "time"
)

func TestCompleteVerificationRejectsInvalidInputs(t *testing.T) {
 for _,tc:=range []struct{name,domain,challenge string}{
  {"missing all","",""},
  {"missing challenge","domain",""},
  {"missing domain","","challenge"},
 } {
  t.Run(tc.name,func(t *testing.T){
   err:=CompleteVerification(context.Background(),nil,tc.domain,tc.challenge,time.Now())
   if !errors.Is(err,ErrInvalid){t.Fatalf("got %v, want ErrInvalid",err)}
  })
 }
}
