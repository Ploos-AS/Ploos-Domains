# Database migrations

PostgreSQL 15+ schema migrations are stored in `db/migrations` using numbered `*.up.sql` and `*.down.sql` files. No migration runner or production database connection is included in this PR.

## Local smoke test (disposable database only)

```sh
createdb ploos_domains_dev
psql -v ON_ERROR_STOP=1 -d ploos_domains_dev -f db/migrations/000001_initial.up.sql
psql -v ON_ERROR_STOP=1 -d ploos_domains_dev -f db/tests/000001_constraints.sql
```

The SQL test runs in a transaction that rolls back. Use a dedicated development database. Do not run down migrations on customer data.

## Guarantees and known limitations

- Membership roles and scope are constrained by the database.
- A unit's parent and a membership's unit must belong to the same organization.
- Duplicate memberships in the same scope are rejected.
- Direct self-parenting is rejected.
- **Longer parent cycles are not yet prevented**. The service must validate hierarchy updates transactionally; a database trigger or closure table should be added before production.
- The schema does not yet implement authentication, authorization, row-level security, DNS allocations or domain ownership verification.
- Deleting an organization cascades into its units and memberships; destructive operations must be protected by service-layer authorization and audit before exposure.
