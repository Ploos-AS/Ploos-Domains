# Browser OIDC startup opt-in guard

`PLOOS_BROWSER_OIDC_LOGIN_ENABLED` defaults to disabled. An explicit value of `false` also disables browser login. Any other non-empty value currently **fails application startup**, because browser OIDC handlers are not yet registered and the full provider, database, proxy and secret wiring is incomplete.

This is an intentional interim safety guard, not functional browser-login activation. It does not change the separate audit API opt-in. A future PR must replace this guard with real `PrepareOIDCLoginHandlers` startup wiring, proper discovery verification, protected client secret loading, route mounting, and integration tests before accepting `true`. Do not set the flag to `true` in production yet.
