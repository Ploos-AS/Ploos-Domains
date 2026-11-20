package domeneshop

import (
 "context"
 "errors"
 "io"
 "net/http"
 "net/url"
 "strings"
 "time"
)

var ErrTransportUnavailable = errors.New("Domeneshop HTTP transport is not configured")
var ErrInvalidRequest = errors.New("invalid Domeneshop HTTP request")

// HTTPDoer allows tests to inject an in-memory transport. No API endpoint,
// authentication or provider request is implemented by this package.
type HTTPDoer interface {
 Do(*http.Request) (*http.Response,error)
}

// Transport is an explicitly injected, read-only HTTP boundary. It cannot
// issue a request without a caller-supplied HTTPDoer.
type Transport struct {
 base *url.URL
 client HTTPDoer
 timeout time.Duration
}

func NewTransport(cfg APIConfig,client HTTPDoer)(*Transport,error){
 if err:=cfg.Validate();err!=nil {return nil,err}
 if client==nil {return nil,ErrTransportUnavailable}
 u,err:=url.Parse(cfg.BaseURL)
 if err!=nil {return nil,ErrInvalidConfiguration}
 return &Transport{base:u,client:client,timeout:cfg.Timeout},nil
}

// Get performs a read-only request to a single relative path segment. It
// rejects absolute URLs, path traversal, queries and fragments. This is a
// low-level testing primitive, NOT an authorized Domeneshop API operation.
func (t *Transport) Get(ctx context.Context,segment string)(*http.Response,error){
 if t==nil || t.client==nil {return nil,ErrTransportUnavailable}
 if ctx==nil || segment=="" || segment=="." || segment==".." ||
  strings.ContainsAny(segment,"/?#\\%") || strings.Contains(segment,":") {
  return nil,ErrInvalidRequest
 }
 u:=*t.base
 u.Path="/"+segment
 bounded,cancel:=context.WithTimeout(ctx,t.timeout)
 req,err:=http.NewRequestWithContext(bounded,http.MethodGet,u.String(),nil)
 if err!=nil {cancel();return nil,ErrInvalidRequest}
 resp,err:=t.client.Do(req)
 if err!=nil {cancel();return nil,err}
 if resp==nil || resp.Body==nil {cancel();return nil,ErrTransportUnavailable}
 resp.Body=&cancelOnClose{ReadCloser:resp.Body,cancel:cancel}
 return resp,nil
}

type cancelOnClose struct{
 io.ReadCloser
 cancel context.CancelFunc
}
func (b *cancelOnClose) Close()error{
 err:=b.ReadCloser.Close()
 b.cancel()
 return err
}
