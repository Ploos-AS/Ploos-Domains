# Browser OIDC runtime lifecycle

The unmounted browser OIDC runtime owns its PostgreSQL pool. `Close` is safe to call repeatedly and concurrently; the pool is closed at most once, and the first close error is returned consistently. Initialization rejects an already canceled context before any Discovery or database network work.

These changes do not register login routes or enable browser login. The existing startup opt-in path still closes the temporary runtime and returns a fail-closed error. No production deployment.
