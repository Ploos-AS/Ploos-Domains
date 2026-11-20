package domeneshop

import (
 "context"
 "errors"
 "fmt"
 "io"
 "net/http"
)

var ErrUnexpectedStatus = errors.New("unexpected registrar HTTP status")
var ErrResponseTooLarge = errors.New("registrar response exceeds limit")
var ErrInvalidResponseLimit = errors.New("invalid registrar response size limit")

// ReadResponse is a bounded, read-only response helper. It never interprets
// provider-specific payloads or includes response bodies in error messages.
// A production client must separately disable redirects in http.Client.
func (t *Transport) ReadResponse(ctx context.Context, segment string, maxBytes int64) ([]byte,error) {
 if maxBytes<=0 || maxBytes>4<<20 { return nil,ErrInvalidResponseLimit }
 resp,err:=t.Get(ctx,segment)
 if err!=nil { return nil,err }
 defer resp.Body.Close()
 if resp.StatusCode!=http.StatusOK {
  return nil,fmt.Errorf("%w: %d",ErrUnexpectedStatus,resp.StatusCode)
 }
 body,err:=io.ReadAll(io.LimitReader(resp.Body,maxBytes+1))
 if err!=nil { return nil,err }
 if int64(len(body))>maxBytes { return nil,ErrResponseTooLarge }
 return body,nil
}

// NewRestrictedHTTPClient creates a real HTTP client with redirect following
// disabled. It is intentionally not used by the inert adapter automatically.
func NewRestrictedHTTPClient(timeoutConfig APIConfig)(*http.Client,error) {
 if err:=timeoutConfig.Validate();err!=nil { return nil,err }
 return &http.Client{
  Timeout:timeoutConfig.Timeout,
  CheckRedirect:func(_ *http.Request,_ []*http.Request)error {
   return http.ErrUseLastResponse
  },
 },nil
}
