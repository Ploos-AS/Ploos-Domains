package server

import (
 "context"
 "errors"
 "io"
 "net/http"
 "net/url"
 "strings"
 "testing"
)

type exchangeTransport func(*http.Request)(*http.Response,error)
func (f exchangeTransport) RoundTrip(r *http.Request)(*http.Response,error){return f(r)}
func TestOIDCCodeExchange(t *testing.T){
 called:=0
 client:=&http.Client{Transport:exchangeTransport(func(r *http.Request)(*http.Response,error){
  called++
  if r.Method!="POST"||r.URL.String()!="https://issuer.example/token"{t.Fatal("unexpected request")}
  if r.Header.Get("Content-Type")!="application/x-www-form-urlencoded"{t.Fatal("missing content type")}
  u,p,ok:=r.BasicAuth();if !ok||u!="client"||p!="secret"{t.Fatal("missing basic auth")}
  raw,_:=io.ReadAll(r.Body);v,err:=url.ParseQuery(string(raw));if err!=nil{t.Fatal(err)}
  if v.Get("code")!="authorization-code"||v.Get("grant_type")!="authorization_code"||v.Get("redirect_uri")!="https://domains.example/callback"||v.Get("code_verifier")!=strings.Repeat("a",43){t.Fatal("wrong exchange fields")}
  return &http.Response{StatusCode:200,Header:http.Header{"Content-Type":[]string{"application/json"}},Body:io.NopCloser(strings.NewReader(`{"id_token":"signed-token"}`))},nil
 })}
 x:=OIDCCodeExchange{Config:OIDCCodeExchangeConfig{TokenEndpoint:"https://issuer.example/token",ClientID:"client",ClientSecret:"secret",RedirectURI:"https://domains.example/callback",HTTPClient:client}}
 got,err:=x.ExchangeCode(context.Background(),"authorization-code",strings.Repeat("a",43))
 if err!=nil||got!="signed-token"||called!=1{t.Fatalf("token=%q err=%v calls=%d",got,err,called)}
}
func TestOIDCCodeExchangeRejectsBadConfig(t *testing.T){
 for _,endpoint:=range []string{"http://issuer.example/token","https://issuer.example/token?x=1","https://user:pass@issuer.example/token","/token"}{
  x:=OIDCCodeExchange{Config:OIDCCodeExchangeConfig{TokenEndpoint:endpoint,ClientID:"client",ClientSecret:"secret",RedirectURI:"https://domains.example/callback"}}
  if _,err:=x.ExchangeCode(context.Background(),"code",strings.Repeat("a",43));!errors.Is(err,ErrOIDCCodeExchange){t.Fatalf("accepted %q: %v",endpoint,err)}
 }
}
func TestOIDCCodeExchangeRejectsRedirect(t *testing.T){
 called:=0
 client:=&http.Client{Transport:exchangeTransport(func(r *http.Request)(*http.Response,error){
  called++
  return &http.Response{StatusCode:302,Header:http.Header{"Location":[]string{"https://other.example/steal"}},Body:io.NopCloser(strings.NewReader(""))},nil
 })}
 x:=OIDCCodeExchange{Config:OIDCCodeExchangeConfig{TokenEndpoint:"https://issuer.example/token",ClientID:"client",ClientSecret:"secret",RedirectURI:"https://domains.example/callback",HTTPClient:client}}
 if _,err:=x.ExchangeCode(context.Background(),"code",strings.Repeat("a",43));!errors.Is(err,ErrOIDCCodeExchange){t.Fatalf("unexpected error: %v",err)}
 if called!=1{t.Fatalf("unexpected redirect followed: %d",called)}
}
