package domains

import (
 "errors"
 "strings"
)

var ErrInvalidTLDPolicy = errors.New("invalid TLD capability policy")

// TLDOperation identifies an explicitly confirmed provider capability.
type TLDOperation string
const (
 OperationAvailability TLDOperation = "availability"
 OperationRegistration TLDOperation = "registration"
 OperationRenewal TLDOperation = "renewal"
 OperationTransferIn TLDOperation = "transfer_in"
 OperationTransferOut TLDOperation = "transfer_out"
 OperationNameservers TLDOperation = "nameservers"
)

var validOperations=map[TLDOperation]bool{
 OperationAvailability:true,OperationRegistration:true,OperationRenewal:true,
 OperationTransferIn:true,OperationTransferOut:true,OperationNameservers:true,
}

// TLDConfiguration is an operator-approved policy record. A true capability
// means verified provider permission, not that an operation has been built.
type TLDConfiguration struct {
 TLD string
 Enabled bool
 IDNSupported bool
 Operations []TLDOperation
}

// TLDCatalog fails closed when an extension or capability is unconfigured.
type TLDCatalog struct { entries map[string]TLDConfiguration }

func NewTLDCatalog(configs []TLDConfiguration)(TLDCatalog,error){
 c:=TLDCatalog{entries:make(map[string]TLDConfiguration,len(configs))}
 for _,cfg:=range configs {
  if cfg.TLD=="" || cfg.TLD!=strings.ToLower(cfg.TLD) || strings.Contains(cfg.TLD,".") {
   return TLDCatalog{},ErrInvalidTLDPolicy
  }
  name,err:=CanonicalName("example."+cfg.TLD)
  if err!=nil || name!="example."+cfg.TLD { return TLDCatalog{},ErrInvalidTLDPolicy }
  if _,exists:=c.entries[cfg.TLD];exists { return TLDCatalog{},ErrInvalidTLDPolicy }
  seen:=map[TLDOperation]bool{}
  for _,op:=range cfg.Operations {
   if !validOperations[op] || seen[op] { return TLDCatalog{},ErrInvalidTLDPolicy }
   seen[op]=true
  }
  cfg.Operations=append([]TLDOperation(nil),cfg.Operations...)
  c.entries[cfg.TLD]=cfg
 }
 return c,nil
}

func (c TLDCatalog) Supports(tld string,op TLDOperation) bool {
 cfg,ok:=c.entries[tld]
 if !ok || !cfg.Enabled || !validOperations[op] { return false }
 for _,allowed:=range cfg.Operations { if allowed==op { return true } }
 return false
}

func (c TLDCatalog) Policy(tld string)(TLDConfiguration,bool){
 cfg,ok:=c.entries[tld]
 if !ok { return TLDConfiguration{},false }
 cfg.Operations=append([]TLDOperation(nil),cfg.Operations...)
 return cfg,true
}

// SearchRegistry enables only TLDs explicitly approved for availability
// queries; a configured TLD is not automatically purchasable.
func (c TLDCatalog) SearchRegistry()(TLDRegistry,error){
 policies:=make([]TLDPolicy,0,len(c.entries))
 for tld,cfg:=range c.entries {
  policies=append(policies,TLDPolicy{TLD:tld,Enabled:cfg.Enabled && c.Supports(tld,OperationAvailability)})
 }
 return NewTLDRegistry(policies)
}
