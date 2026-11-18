package auth

import (
 "errors"
 "testing"
 "time"
)

func TestAuditEventValidation(t *testing.T){
 valid:=AuditEvent{OrganizationID:"org",Type:AuditInvitationIssued,OccurredAt:time.Now()}
 if err:=valid.Validate();err!=nil{t.Fatal(err)}
 valid.Type="arbitrary"
 if err:=valid.Validate();!errors.Is(err,ErrInvalidAuditEvent){t.Fatal(err)}
 valid.Type=AuditInvitationRedeemed
 valid.OrganizationID=""
 if err:=valid.Validate();!errors.Is(err,ErrInvalidAuditEvent){t.Fatal(err)}
}
