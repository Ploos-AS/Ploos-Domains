BEGIN;
DROP TABLE IF EXISTS subdomain_allocations;
DROP TABLE IF EXISTS domains;
DROP FUNCTION IF EXISTS check_subdomain_allocation();
DROP FUNCTION IF EXISTS prevent_domain_identity_change();
COMMIT;
