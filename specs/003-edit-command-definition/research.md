# Research Findings: pm edit flag-driven updates

## Decision: Performance Targets for pm edit
- **Rationale**: Align command latency with repository success criteria while accounting for bootstrap overhead and configuration size. Targets ensure responsive scripting and parity with other `pm` commands. Small configs (≤200 lines) should complete within p95 750 ms, medium (≤800 lines) within 1.4 s, large within 2.2 s, keeping absolute maxima under 3 s. Throughput expectations (≥20 edits/min sequentially, bursts of 5 edits ≤10 s) enable automation without overloading typical developer hardware.
- **Alternatives Considered**: (1) Rely solely on existing 5 s SLA—too lax for flag-driven automation. (2) Set uniform sub-500 ms target—unrealistic for cold-start and large configs. Adopted tiered targets as balanced compromise.

## Decision: Operational Constraints & Safeguards
- **Rationale**: Preserve hexagonal architecture guarantees, enforce immutability, and avoid partial writes. Lock per project across read/validate/write, treat `project.Path` and environment identity as immutable, and require explicit clears for non-empty metadata. Reuse validation to ensure schema integrity, enforce atomic filesystem operations, and preflight path resolution. This protects downstream tooling and avoids accidental data corruption.
- **Alternatives Considered**: (1) Allow direct file writes from CLI—would bypass repositories and risk stale projections. (2) Skip locks relying on fast operations—risks race conditions with delete/new commands. Adopted repository services + locking for safety.
