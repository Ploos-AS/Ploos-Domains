package server

import (
 "context"
 "encoding/json"
 "errors"
 "fmt"
 "io"
 "net/http"
 "net/url"
 "strings"
 "time"
)

var ErrOIDCCodeExchange = errors.New("OIDC code exchange failed")

// OIDCCodeExchangeConfig must be configured from trusted OIDC discovery and
// application settings, never from a callback request.
type OIDCCodeExchangeConfig struct {
 TokenEndpoint string
 ClientID string
 ClientSecret string
 RedirectURI string
 HTTPClient *http.Client
}

type OIDCCodeExchange struct { Config OIDCCodeExchangeConfig }

func (x OIDCCodeExchange) ExchangeCode(ctx context.Context,code,pkceVerifier string)(string,error){
 c:=x.Config
 endpoint,err:=url.Parse(c.TokenEndpoint)
 if err!=nil||endpoint.Scheme!="https"||endpoint.Host==""||endpoint.User!=nil||endpoint.Fragment!=""||endpoint.RawQuery!=""||
  c.ClientID==""||c.ClientSecret==""||code==""||len(pkceVerifier)<43||len(pkceVerifier)>128||ctx==nil{
  return "",ErrOIDCCodeExchange
 }
 callback,err:=url.Parse(c.RedirectURI)
 if err!=nil||callback.Scheme!="https"||callback.Host==""||callback.User!=nil||callback.Fragment!=""{
  return "",ErrOIDCCodeExchange
 }
 for _,ch:=range pkceVerifier{
  if !((ch>='a'&&ch<='z')||(ch>='A'&&ch<='Z')||(ch>='0'&&ch<='9')||ch=='-'||ch=='.'||ch=='_'||ch=='~'){
   return "",ErrOIDCCodeExchange
  }
 }
 form:=url.Values{
  "grant_type":{"authorization_code"},
  "code":{code},
  "redirect_uri":{c.RedirectURI},
  "code_verifier":{pkceVerifier},
 }
 req,err:=http.NewRequestWithContext(ctx,http.MethodPost,c.TokenEndpoint,strings.NewReader(form.Encode()))
 if err!=nil{return "",ErrOIDCCodeExchange}
 req.Header.Set("Content-Type","application/x-www-form-urlencoded")
 req.Header.Set("Accept","application/json")
 req.SetBasicAuth(c.ClientID,c.ClientSecret)
 client:=c.HTTPClient
 if client==nil{client=&http.Client{Timeout:10*time.Second}}
 // Do not forward credentials or authorization codes through redirects.
 safeClient:=*client
 safeClient.CheckRedirect=func(*http.Request,[]*http.Request)error{return http.ErrUseLastResponse}
 resp,err:=safeClient.Do(req)
 if err!=nil{return "",fmt.Errorf("%w: request failed",ErrOIDCCodeExchange)}
 defer resp.Body.Close()
 if resp.StatusCode!=http.StatusOK{return "",ErrOIDCCodeExchange}
 if !strings.HasPrefix(strings.ToLower(resp.Header.Get("Content-Type")),"application/json"){return "",ErrOIDCCodeExchange}
 body:=io.LimitReader(resp.Body,65537)
 var payload struct{ IDToken string `json:"id_token"` }
 raw,err:=io.ReadAll(body)
 if err!=nil||len(raw)>65536{return "",ErrOIDCCodeExchange}
 if err:=json.Unmarshal(raw,&payload);err!=nil||payload.IDToken==""||len(payload.IDToken)>32768{return "",ErrOIDCCodeExchange}
 return payload.IDToken,nil
}
