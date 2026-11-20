package domeneshop

import (
 "errors"
 "testing"
 "time"
)

func TestAPIConfigValidation(t *testing.T){
 valid:=APIConfig{BaseURL:"https://api.example.test",Timeout:5*time.Second}
 if err:=valid.Validate();err!=nil {t.Fatal(err)}
 for _,url:=range []string{"http://api.example.test","https://user:pass@api.example.test","https://api.example.test/path","https://api.example.test?x=1","https://api.example.test/#fragment","https://api.example.test:8080","https://"," https://api.example.test"} {
  cfg:=valid;cfg.BaseURL=url
  if err:=cfg.Validate();!errors.Is(err,ErrInvalidConfiguration){t.Errorf("accepted %q: %v",url,err)}
 }
 for _,timeout:=range []time.Duration{0,-time.Second,31*time.Second} {
  cfg:=valid;cfg.Timeout=timeout
  if err:=cfg.Validate();!errors.Is(err,ErrInvalidConfiguration){t.Errorf("accepted timeout %v: %v",timeout,err)}
 }
}
