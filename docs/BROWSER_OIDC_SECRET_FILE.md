# Browser OIDC confidential client secret file

The standalone `loadBrowserOIDCClientSecret` helper accepts an absolute, normalized path to a regular file with no group or other access bits (mode 0600 or stricter). It rejects symlinks, empty files, oversized files, embedded newlines and leading/trailing whitespace. Secret values are never included in returned errors.

This is an isolated building block and **is not yet called during startup**. The runtime service should mount a dedicated secret file with restrictive ownership and permissions. Container secret volume modes and UID/GID must be verified before wiring this helper into the startup sequence. Browser OIDC routes remain unmounted and login disabled.
