# Portfolio HTTP pagination

`GET /api/v1/organizations/{organizationID}/portfolio?limit=50&cursor=example.no`

- Default limit: 50; allowed range: 1–100.
- `cursor` is the last FQDN returned on the preceding page; omit for first page.
- Response: `{ "domains": [...], "next_cursor": "..." }`; `next_cursor` omitted on final page.
- OIDC and organization-level membership checks remain mandatory.
- The endpoint now calls `ListPageForUser`; the unbounded reader is not used by this handler.
- Cursor is not an authorization token. Membership is rechecked on each page.
- Cursor order is lexical; concurrent additions/deletions may affect traversal. A snapshot is not guaranteed.
- This is not yet a production rollout. PostgreSQL cursor traversal and HTTP end-to-end integration tests should be completed before deployment.
