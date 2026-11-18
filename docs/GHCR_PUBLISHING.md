# GHCR image publication

The `Publish container image` workflow builds and pushes `linux/amd64` images to `ghcr.io/ploos-as/ploos-domains` on merges to `main`, version tags (`v*`) and manual dispatch. It uses the built-in `GITHUB_TOKEN` with `packages: write`; **no paid GitHub plan or personal access token is required for this workflow**.

Images receive immutable-commit-derived `sha-...` tags and release tags. Deployments should use the **digest** printed in the workflow's final step, for example `ghcr.io/ploos-as/ploos-domains@sha256:<digest>`. Do not use mutable tags as the deployment source of truth.

The package may initially be private even though the repository is public. Confirm GHCR package visibility and grant your deployment host read access, or explicitly make the package public if appropriate. Do not place credentials in Compose files or the repository.

This workflow publishes images but **does not deploy** them. Deploy with Ploos Deploy only after CI succeeds and the image is available in GHCR. Keep production configuration and secrets outside the repository.

Security follow-up: pin third-party GitHub Actions by full commit SHA, add image vulnerability scanning and attestations, and consider publishing only after required CI checks have passed.
