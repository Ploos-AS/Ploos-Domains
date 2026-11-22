#!/usr/bin/env bash
set -euo pipefail

# Manual, opt-in public-metadata smoke tests against a dedicated test provider.
# This script does not configure or deploy Authelia and does not enable login.
if [[ "${PLOOS_TEST_LIVE_OIDC_ENABLED:-}" != "true" ]]; then
  echo "Refusing live OIDC test: set PLOOS_TEST_LIVE_OIDC_ENABLED=true explicitly." >&2
  exit 2
fi

required=(
  PLOOS_TEST_LIVE_OIDC_ISSUER
  PLOOS_TEST_LIVE_OIDC_AUTHORIZATION_ENDPOINT
  PLOOS_TEST_LIVE_OIDC_TOKEN_ENDPOINT
  PLOOS_TEST_LIVE_OIDC_CLIENT_ID
  PLOOS_TEST_LIVE_OIDC_REDIRECT_URI
  PLOOS_TEST_LIVE_OIDC_PUBLIC_HOST
)
for key in "${required[@]}"; do
  if [[ -z "${!key:-}" ]]; then
    echo "Missing required test configuration: $key" >&2
    exit 2
  fi
done

case "$PLOOS_TEST_LIVE_OIDC_ISSUER" in
  https://*) ;;
  *) echo "Test issuer must use HTTPS." >&2; exit 2 ;;
esac
case "$PLOOS_TEST_LIVE_OIDC_REDIRECT_URI" in
  https://*) ;;
  *) echo "Test callback must use HTTPS." >&2; exit 2 ;;
esac

echo "Running three opt-in OIDC public-provider smoke tests (no login or secrets)."
go test -tags=integration -count=1 -run '^(TestBrowserOIDCLiveProviderDiscovery|TestBrowserOIDCLiveVerifierInitialization|TestBrowserOIDCLiveJWKSReadiness)$' ./internal/server
