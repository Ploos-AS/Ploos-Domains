package server

import (
 "context"
 "io"
 "net/http"
 "strings"
 "testing"
)

func TestOIDCDiscoveryMatchesConfiguredProvider(t *testing.T){
 cfg:=validOIDCDeploymentConfig()
 calls:=0
 d:=OIDCDiscoveryClient{HTTPClient:&http.Client{Transport:exchangeTransport(func(r *http.Request)(*http.Response,error){
  calls++
  if r.URL.String()!="https://issuer.example/.well-known/openid-configuration"{t.Fatalf("wrong discovery URL: %s",r.URL)}
  return &http.Response{StatusCode:200,Header:http.Header{"Content-Type":[]string{"application/json"}},Body:io.NopCloser(strings.NewReader(`{"issuer":"https://issuer.example","authorization_endpoint":"https://issuer.example/authorize","token_endpoint":"https://issuer.example/token","jwks_uri":"https://issuer.example/keys"}`))},nil
 })}}
 if err:=d.Validate(context.Background(),cfg);err!=nil{t.Fatal(err)}
 if calls!=1{t.Fatalf("calls=%d",calls)}
}
func TestOIDCDiscoveryRejectsMismatch(t *testing.T){
 cfg:=validOIDCDeploymentConfig()
 cases:=[]string{
  `{"issuer":"https://other.example","authorization_endpoint":"https://issuer.example/authorize","token_endpoint":"https://issuer.example/token","jwks_uri":"https://issuer.example/keys"}`,
  `{"issuer":"https://issuer.example","authorization_endpoint":"https://evil.example/authorize","token_endpoint":"https://issuer.example/token","jwks_uri":"https://issuer.example/keys"}`,
  `{"issuer":"https://issuer.example","authorization_endpoint":"https://issuer.example/authorize","token_endpoint":"https://evil.example/token","jwks_uri":"https://issuer.example/keys"}`,
  `{"issuer":"https://issuer.example","authorization_endpoint":"https://issuer.example/authorize","token_endpoint":"https://issuer.example/token","jwks_uri":"http://issuer.example/keys"}`,
 }
 for _,body:=range cases{
  d:=OIDCDiscoveryClient{HTTPClient:&http.Client{Transport:exchangeTransport(func(*http.Request)(*http.Response,error){
   return &http.Response{StatusCode:200,Header:http.Header{"Content-Type":[]string{"application/json"}},Body:io.NopCloser(strings.NewReader(body))},nil
  })}}
  if err:=d.Validate(context.Background(),cfg);err==nil{t.Fatalf("accepted mismatched metadata: %s",body)}
 }
}
func TestOIDCDiscoveryRejectsRedirect(t *testing.T){
 cfg:=validOIDCDeploymentConfig()
 calls:=0
 d:=OIDCDiscoveryClient{HTTPClient:&http.Client{Transport:exchangeTransport(func(*http.Request)(*http.Response,error){
  calls++
  return &http.Response{StatusCode:302,Header:http.Header{"Location":[]string{"https://evil.example/metadata"}},Body:io.NopCloser(strings.NewReader(""))},nil
 })}}
 if err:=d.Validate(context.Background(),cfg);err==nil{t.Fatal("accepted discovery redirect")}
 if calls!=1{t.Fatalf("followed redirect: %d",calls)}
}
