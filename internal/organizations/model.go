// Package organizations defines the domain model and authorization boundaries
// for organizations and their local units. Persistence is deliberately separate.
package organizations

import (
 "errors"
 "strings"
)

type Role string

const (
 RoleOwner Role = "owner"
 RoleAdmin Role = "admin"
 RoleUnitAdmin Role = "unit_admin"
 RoleViewer Role = "viewer"
)

type Organization struct {
 ID string
 Name string
}

type Unit struct {
 ID string
 OrganizationID string
 ParentID string
 Name string
}

type Membership struct {
 UserID string
 OrganizationID string
 UnitID string // Empty means organization-level membership.
 Role Role
}

type Action string

const (
 View Action = "view"
 ManageOrganization Action = "manage_organization"
 ManageUnits Action = "manage_units"
 ManageAllocatedDNS Action = "manage_allocated_dns"
)

var (
 ErrInvalidMembership = errors.New("invalid membership scope")
 ErrForbidden = errors.New("forbidden")
)

// Authorize enforces explicit tenant and unit scopes. The caller must pass
// the target unit's organization ID from trusted storage, never user input.
// Descendant-unit inheritance is intentionally NOT supported in v1.
func Authorize(m Membership, targetOrgID, targetUnitID string, action Action) error {
 if m.UserID == "" || m.OrganizationID == "" || targetOrgID == "" || m.OrganizationID != targetOrgID {
  return ErrForbidden
 }
 if m.UnitID != "" {
  if m.Role != RoleUnitAdmin && m.Role != RoleViewer { return ErrInvalidMembership }
  if targetUnitID == "" || m.UnitID != targetUnitID { return ErrForbidden }
  switch action {
  case View: return nil
  case ManageAllocatedDNS:
   if m.Role == RoleUnitAdmin { return nil }
  }
  return ErrForbidden
 }
 switch m.Role {
 case RoleOwner, RoleAdmin:
  switch action {
  case View, ManageOrganization, ManageUnits, ManageAllocatedDNS: return nil
  }
 case RoleViewer:
  if action == View { return nil }
 }
 return ErrForbidden
}

// ValidateUnit performs basic validation. Parent ownership and cycle checks
// must be enforced transactionally by the persistence layer.
func ValidateUnit(u Unit) error {
 if strings.TrimSpace(u.ID) == "" || strings.TrimSpace(u.OrganizationID) == "" || strings.TrimSpace(u.Name) == "" {
  return errors.New("unit id, organization id and name are required")
 }
 if u.ID == u.ParentID { return errors.New("unit cannot be its own parent") }
 return nil
}
