package auth

import (
 "context"
 "errors"
 "fmt"
 "net/url"
 "strings"

 "github.com/coreos/go-oidc/v3/oidc"
)

// OIDCVerifier validates signed OIDC ID tokens using provider discovery and
// remotely managed JWKS. It is not an OAuth2 access-token introspector.
type OIDCVerifier struct {
 issuer string
 audience string
 verifier *oidc.IDTokenVerifier
}

// NewOIDCVerifier discovers a configured HTTPS issuer and its JWKS endpoint.
// Never derive the issuer from a token or from an untrusted request.
func NewOIDCVerifier(ctx context.Context, issuer, audience string) (*OIDCVerifier,error) {
 if err:=ValidateOIDCConfig(issuer,audience);err!=nil {return nil,err}
 provider,err:=oidc.NewProvider(ctx,issuer)
 if err!=nil{return nil,fmt.Errorf("discover OIDC provider: %w",err)}
 if provider.Endpoint().AuthURL=="" {return nil,errors.New("OIDC provider missing authorization endpoint")}
 return &OIDCVerifier{
  issuer:issuer,audience:audience,
  verifier:provider.Verifier(&oidc.Config{ClientID:audience,SupportedSigningAlgs:[]string{"RS256","ES256"}}),
 },nil
}

// ValidateOIDCConfig rejects non-HTTPS, ambiguous and userinfo-bearing issuers.
func ValidateOIDCConfig(issuer,audience string)error{
 if issuer==""||strings.TrimSpace(audience)!=audience||audience==""{return ErrInvalidIdentity}
 u,err:=url.Parse(issuer)
 if err!=nil||u.Scheme!="https"||u.Host==""||u.User!=nil||u.RawQuery!=""||u.Fragment!=""||strings.HasSuffix(issuer,"/") {
  return ErrInvalidIdentity
 }
 return nil
}

func (v *OIDCVerifier) VerifyToken(ctx context.Context,rawToken string)(Identity,error){
 if v==nil||v.verifier==nil||rawToken==""{return Identity{},ErrUnauthenticated}
 token,err:=v.verifier.Verify(ctx,rawToken)
 if err!=nil{return Identity{},fmt.Errorf("%w: %v",ErrUnauthenticated,err)}
 if token.Subject==""||token.Issuer!=v.issuer {return Identity{},ErrInvalidIdentity}
 return Identity{Issuer:token.Issuer,Subject:token.Subject,Audience:v.audience},nil
}

var _ Verifier = (*OIDCVerifier)(nil)
