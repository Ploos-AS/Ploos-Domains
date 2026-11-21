# Domain portfolio keyset pagination

`PortfolioReader.ListPageForUser(ctx, trustedUserID, organizationID, limit, cursor)` provides a bounded organization-authorized PostgreSQL read. The cursor is the last returned canonical domain name (exclusive), sorted lexicographically. Page sizes must be 1–100. A query fetches at most `limit+1` rows to detect whether a next page exists.

Membership authorization is checked within the same SQL query; unit-only membership does not grant access. The cursor is scoped to the organization query and cannot authorize access. The underlying `domain_portfolio` primary key supports organization-scoped FQDN ordering.

**Not yet wired into HTTP**: the existing API still uses unpaged `ListForUser`. A follow-up PR must replace it with the paginated method and enforce a page size on all customer-facing requests. Do not advertise HTTP pagination before that change. Integration tests with PostgreSQL, including cursor traversal and membership revocation, are also required.

No registrar operations, migrations or production deployment.
