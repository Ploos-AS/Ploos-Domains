package organizations

import (
 "errors"
 "testing"
)

func TestAuthorizationBoundaries(t *testing.T) {
 cases := []struct {
  name string
  member Membership
  org, unit string
  action Action
  want error
 }{
  {"owner can manage organization", Membership{"alice","org1","",RoleOwner},"org1","",ManageOrganization,nil},
  {"admin can manage units", Membership{"alice","org1","",RoleAdmin},"org1","sirdal",ManageUnits,nil},
  {"unit admin manages own allocation", Membership{"bob","org1","sirdal",RoleUnitAdmin},"org1","sirdal",ManageAllocatedDNS,nil},
  {"unit admin cannot manage sibling", Membership{"bob","org1","sirdal",RoleUnitAdmin},"org1","bergen",ManageAllocatedDNS,ErrForbidden},
  {"unit admin cannot manage parent zone", Membership{"bob","org1","sirdal",RoleUnitAdmin},"org1","",ManageAllocatedDNS,ErrForbidden},
  {"unit admin cannot manage organization", Membership{"bob","org1","sirdal",RoleUnitAdmin},"org1","sirdal",ManageOrganization,ErrForbidden},
  {"tenant isolation", Membership{"alice","org1","",RoleOwner},"org2","sirdal",ManageAllocatedDNS,ErrForbidden},
  {"viewer read only", Membership{"bob","org1","sirdal",RoleViewer},"org1","sirdal",View,nil},
  {"viewer cannot write", Membership{"bob","org1","sirdal",RoleViewer},"org1","sirdal",ManageAllocatedDNS,ErrForbidden},
  {"org viewer cannot write", Membership{"bob","org1","",RoleViewer},"org1","sirdal",ManageUnits,ErrForbidden},
  {"unit-scoped owner invalid", Membership{"bob","org1","sirdal",RoleOwner},"org1","sirdal",View,ErrInvalidMembership},
  {"unknown role denied", Membership{"bob","org1","",Role("superuser")},"org1","",View,ErrForbidden},
  {"unknown action denied", Membership{"alice","org1","",RoleOwner},"org1","",Action("delete_all"),ErrForbidden},
  {"missing identity denied", Membership{"","org1","",RoleOwner},"org1","",View,ErrForbidden},
 }
 for _, tc := range cases {
  t.Run(tc.name, func(t *testing.T) {
   err := Authorize(tc.member,tc.org,tc.unit,tc.action)
   if tc.want == nil && err != nil { t.Fatalf("unexpected error: %v",err) }
   if tc.want != nil && !errors.Is(err,tc.want) { t.Fatalf("got %v, want %v",err,tc.want) }
  })
 }
}

func TestValidateUnit(t *testing.T) {
 valid := Unit{ID:"sirdal",OrganizationID:"seniornett",Name:"Sirdal"}
 if err:=ValidateUnit(valid);err!=nil { t.Fatal(err) }
 valid.ParentID="sirdal"
 if err:=ValidateUnit(valid);err==nil { t.Fatal("expected self-parent rejection") }
 valid=Unit{ID:"x",OrganizationID:"org",Name:"  "}
 if err:=ValidateUnit(valid);err==nil { t.Fatal("expected blank name rejection") }
}
