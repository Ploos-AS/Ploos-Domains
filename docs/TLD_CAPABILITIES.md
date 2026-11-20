# TLD capabilities catalogue

`TLDCatalog` is an in-memory, provider-neutral operator policy layer. It does **not** claim any TLD is available from Domeneshop or that the application has implemented purchase/transfer operations.

- Each TLD has an explicit enabled flag and confirmed operation list.
- Availability, registration, renewal, transfer-in, transfer-out and nameserver changes are separate permissions.
- `SearchRegistry` derives only explicitly enabled availability capabilities for the existing multi-TLD search contract.
- IDN support is recorded as a policy flag only. Unicode input conversion and per-registry IDN rules are **not implemented**.
- A policy's `registration` flag alone does not authorize orders; checkout, prices, identity, billing and an actual provider adapter remain mandatory.
- Disabled or unknown TLDs fail closed.
- This first version deliberately supports single-label TLDs only; multi-label public suffixes need a separate registry policy and matching engine.

## Before real Domeneshop integration

Confirm partner agreement, documented API capabilities per TLD, prices and registry restrictions. Store only confirmed capabilities in the catalogue, and add validation, versioning, cache invalidation and observability before exposing it through a public API.
