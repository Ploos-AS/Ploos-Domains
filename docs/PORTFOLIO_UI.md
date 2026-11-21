# Portfolio UI foundation

Static, responsive, bilingual (Norwegian/English) and keyboard-accessible HTML/CSS/JS shell under `web/portfolio/`.

This is **not yet connected to customer data**: there is no OIDC login, organization selector, API fetch, token handling or production mounting. The page intentionally shows an empty state and a disabled next-page button. It must not be described as a working customer portal.

Next milestones:
1. Select a secure OIDC browser session strategy (prefer backend-managed HttpOnly/Secure/SameSite cookie; do not persist bearer tokens in localStorage).
2. Resolve the current user's authorized organizations server-side.
3. Wire paginated portfolio API and handle loading, errors, empty states and navigation.
4. Add automated accessibility tests and deploy only after authentication and integration tests.

Preview by serving `web/portfolio` as static files in a local development environment. No registrar or DNS changes.
