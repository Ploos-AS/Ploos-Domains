// Package domains validates canonical ASCII DNS names for subdomain allocation.
package domains

import (
 "errors"
 "strings"
)

var ErrInvalidName = errors.New("invalid DNS name")

// CanonicalName accepts ASCII DNS hostnames only (not IDNs in Unicode form).
// Internationalized names must be converted to A-labels by a separate IDNA
// validation layer before calling this function.
func CanonicalName(name string) (string, error) {
 name = strings.TrimSuffix(name, ".")
 if len(name) == 0 || len(name) > 253 { return "", ErrInvalidName }
 name = strings.ToLower(name)
 for _, c := range name {
  if c > 127 { return "", ErrInvalidName }
 }
 for _, label := range strings.Split(name, ".") {
  if len(label) == 0 || len(label) > 63 || label[0] == '-' || label[len(label)-1] == '-' { return "", ErrInvalidName }
  for _, c := range label {
   if !((c >= 'a' && c <= 'z') || (c >= '0' && c <= '9') || c == '-') { return "", ErrInvalidName }
  }
 }
 return name, nil
}

// AllocationName ensures a candidate is strictly below a verified zone.
// The caller must separately verify domain control and organization scope.
func AllocationName(candidate, zone string) (string, error) {
 c, err := CanonicalName(candidate)
 if err != nil { return "", err }
 z, err := CanonicalName(zone)
 if err != nil { return "", err }
 if len(c) <= len(z)+1 || !strings.HasSuffix(c, "."+z) { return "", ErrInvalidName }
 return c, nil
}
