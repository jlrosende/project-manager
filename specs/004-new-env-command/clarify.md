# `pm new` Clarification Questions

## Q1. Project Existence Detection
- **Description:** Determine how the CLI decides whether the named project already exists before interpreting the second positional argument.
- **Options:**
  - **A.** Query `ProjectService.Load` by name only; success means existing project.
  - **B.** Check for `.project.hcl` at resolved path if provided; existence implies project.
  - **C.** Require both name match via repository list and `.project.hcl` on disk.
  - **D.** Treat any second argument as an environment unless an explicit `--path` override is present.

- Answer: C

## Q2. Environment Flag Set
- **Description:** Define the set of CLI flags that configure environment creation in shortcut mode.
- **Options:**
  - **A.** No new flags; positional env name plus defaults only.
  - **B.** Introduce flags for `--env-name`, `--env-color`, `--env-file`, `--env-mode`.
  - **C.** Mirror existing project-level flags (e.g., metadata) plus an env vars file path.
  - **D.** Postpone flags; require `--cli-input` for advanced environment configuration.

- Answer: B but without the `--env-name` flag

## Q3. Config File Schema Updates
- **Description:** Decide how CLI input files should represent environment creation while staying backward compatible.
- **Options:**
  - **A.** Add top-level `environment` object and keep current `environments` map untouched.
  - **B.** Deprecate `environments` map in favor of a structured `environment` object.
  - **C.** Support both, merging entries; conflicts resolved by favoring `environment` object.
  - **D.** Require environment definitions to move under a new `environments` array of objects.

- Answer: B

## Q4. Dry Run Output Shape
- **Description:** Specify fields included in dry-run previews when an environment would be added.
- **Options:**
  - **A.** Include env name and env file only.
  - **B.** Include full environment struct (name, mode, file, color) plus env vars count.
  - **C.** Include environment data only in JSON format; text stays minimal.
  - **D.** Extend both text and JSON outputs with the same rich environment details.

- Answer: B

## Q5. Skeleton Template Backward Compatibility
- **Description:** Determine how skeleton generation should expose the new environment fields.
- **Options:**
  - **A.** Append an `environment` object while retaining existing keys.
  - **B.** Replace existing content with an environment-focused schema.
  - **C.** Generate separate project and environment skeleton sections with comments.
  - **D.** Provide two skeleton types selected via new flag (project vs environment).

- Answer: B

## Q6. Force Flag Semantics
- **Description:** Clarify whether `--force` affects environment addition when invoked via shortcut mode.
- **Options:**
  - **A.** `--force` applies to both project scaffolding and environment overwrites.
  - **B.** `--force` only impacts project creation; env shortcut ignores it.
  - **C.** Introduce a separate `--env-force` flag for environment overwrites.
  - **D.** Reject `--force` when adding environments; require manual edit workflow.

- Answer: A

## Q7. Error Messaging Strategy
- **Description:** Define how errors are surfaced for missing projects, duplicate environments, or ambiguous inputs.
- **Options:**
  - **A.** Return raw service errors verbatim.
  - **B.** Wrap service errors with CLI-friendly context strings.
  - **C.** Provide custom error codes/messages for each failure mode.
  - **D.** Introduce structured output when `--output json` is used, plain text otherwise.

- Answer: C

## Q8. Uniqueness Enforcement Layer
- **Description:** Select where name and path uniqueness rules should be enforced for the new behavior.
- **Options:**
  - **A.** CLI validates before calling services.
  - **B.** Services enforce uniqueness; CLI passes through errors.
  - **C.** Repository layer is authoritative; CLI does minimal checks.
  - **D.** Shared helper validates and is reused by CLI and services.

- Answer: D

## Q9. Existing Project Without Second Argument
- **Description:** Decide outcome when `pm new <project>` is run for an existing project with no path/env argument.
- **Options:**
  - **A.** Fail with guidance to use `pm edit` or env shortcut.
  - **B.** Treat as no-op success with informational message.
  - **C.** Re-run project initialization (potentially destructive).
  - **D.** Prompt interactively if TUI available; fail in headless mode.

- Answer: B return a message indicating the project already exist

## Q10. Ambiguous Second Argument Values
- **Description:** Handle cases where users pass a filesystem-looking second argument intending a path for an existing project.
- **Options:**
  - **A.** Treat any existing project + second arg as env name regardless of format.
  - **B.** Heuristically treat inputs with `/` or `.` as paths, otherwise env names.
  - **C.** Require explicit `--path` flag to disambiguate; positional becomes env-only.
  - **D.** Prompt for confirmation when ambiguity detected (fallback to env mode on non-interactive).

- Answer: A

## Q11. CLI Input vs Positional Environment Precedence
- **Description:** Resolve conflicts when CLI input file defines an environment name different from the positional argument.
- **Options:**
  - **A.** Positional argument always overrides config.
  - **B.** Config file takes precedence; positional is ignored.
  - **C.** Error out on mismatch and request explicit clarification.
  - **D.** Merge: positional provides name, config supplies remaining fields.

- Answer: D

## Q12. Environment Defaults Without Positional Name
- **Description:** Define behavior when `--cli-input` provides environment data but command-line omits a second positional argument.
- **Options:**
  - **A.** Use `environment.name` from config and proceed.
  - **B.** Require positional name; otherwise reject execution.
  - **C.** Derive name from project + mode (e.g., `default`).
  - **D.** Skip environment creation entirely unless positional argument supplied.

- Answer: A

## Q13. Combined Project and Environment Creation Flow
- **Description:** Decide sequencing when a single invocation should both create a project and add an environment via config or flags.
- **Options:**
  - **A.** Always create project first, then environment if project creation succeeded.
  - **B.** Allow environment creation only on subsequent runs; reject combined attempts.
  - **C.** Support atomic create+env with rollback if either step fails.
  - **D.** Gate combined workflow behind a new `--with-environment` flag.

- Answer: B

## Q14. Dry Run With CLI Input and Environment Data
- **Description:** Establish expected output when `--dry-run` is used alongside `--cli-input` containing environment configuration.
- **Options:**
  - **A.** Show both project and environment data sourced from config.
  - **B.** Show project only; note that environment would be added but omit details.
  - **C.** Require positional env argument in dry run; otherwise skip environment preview.
  - **D.** Error if config tries to add environment during dry run.

- Answer: A

## Q15. Case Sensitivity Policy
- **Description:** Confirm whether project and environment names remain case-sensitive in detection and duplication checks.
- **Options:**
  - **A.** Keep strict case sensitivity for both project and environment names.
  - **B.** Case-insensitive comparisons for environments only; projects stay sensitive.
  - **C.** Case-insensitive for both names and paths.
  - **D.** Normalize to lowercase before storage and comparison.

- Answer: A

## Q16. Final Environment Flag Set
- **Description:** Clarify the exact flag surface for environment customization given the decision to add flags from option B without `--env-name`.
- **Options:**
  - **A.** Provide `--env-color`, `--env-file`, and `--env-mode` flags only.
  - **B.** Expose all flags from option B but alias `--env-name` to the positional argument for backward compatibility.
  - **C.** Introduce a single `--environment` flag accepting JSON/YAML to cover advanced fields.
  - **D.** Publish `--env-color`, `--env-file`, `--env-mode`, and reserve `--env-name` for future use with clear documentation.

- Answer: A

## Q17. Deprecation Path for Legacy `environments` Map
- **Description:** Decide how existing CLI input files using the map-based `environments` field should behave after adopting the structured `environment` object.
- **Options:**
  - **A.** Continue parsing the legacy map with a warning, favoring the new object when both are present.
  - **B.** Fail fast when the legacy map is detected, requiring migration before continuing.
  - **C.** Support both schemas indefinitely but emit telemetry/metrics for legacy usage.
  - **D.** Provide an automatic converter that rewrites legacy configs on load.

- Answer: B

## Q18. Config-Driven Environment on Project Creation
- **Description:** Determine behavior when `--cli-input` supplies environment data but the target project does not yet exist (given the choice to disallow combined create+env flows).
- **Options:**
  - **A.** Reject the run with guidance to re-run once the project is created.
  - **B.** Ignore environment data silently during project creation.
  - **C.** Create the project, then stage the environment addition but require manual confirmation before applying.
  - **D.** Automatically invoke a second internal pass to add the environment after creation completes.

- Answer: A

## Q19. Dry Run Preview for Config-Only Environments
- **Description:** Align dry-run output with the requirement to include full environment details while also requiring positional arguments for previews.
- **Options:**
  - **A.** Permit config-sourced environments to appear in dry-run previews even without a positional argument.
  - **B.** Enforce the positional requirement by emitting a warning that environment details are hidden.
  - **C.** Require an explicit `--with-environment-preview` flag to include config-sourced environments in dry runs.
  - **D.** Automatically downgrade dry-run output to project-only when no positional argument is present, documenting the behavior.

- Answer: A

## Q20. `--force` Behavior for Existing Environments
- **Description:** Specify what "force" means when an environment with the same name already exists.
- **Options:**
  - **A.** Overwrite environment metadata and env vars file contents in place.
  - **B.** Allow re-generation of env vars file but keep existing metadata untouched.
  - **C.** Only allow force when the environment is marked inactive or archived elsewhere.
  - **D.** Disallow overwriting existing environments; `--force` only regenerates project scaffolding.

- Answer: A

## Q21. Custom Error Code Format
- **Description:** Define the structure for the custom error codes/messages promised in Q7.
- **Options:**
  - **A.** Human-readable slug strings (e.g., `ERR_PROJECT_NOT_FOUND`).
  - **B.** Numeric codes mapped to documentation.
  - **C.** Hierarchical codes combining category and detail (e.g., `E/PROJECT/NOT_FOUND`).
  - **D.** JSON objects with `code` and `message` fields when `--output json` is active, plain text otherwise.

- Answer: A

## Q22. Missing `.project.hcl` During Existence Check
- **Description:** Handle cases where the repository lists a project name but the `.project.hcl` file is absent at the stored path (e.g., file manually deleted).
- **Options:**
  - **A.** Treat the project as existing and continue with environment mode, logging a warning about the missing file.
  - **B.** Treat the project as non-existent and fall back to project creation mode.
  - **C.** Fail with an actionable error instructing the user to repair or delete the project entry.
  - **D.** Attempt to regenerate `.project.hcl` automatically before proceeding.

- Answer: B

## Q23. Dry Run with Config Environments for Nonexistent Projects
- **Description:** When `--cli-input` supplies environment data but the project does not yet exist (Answer: Q18=A), determine how `--dry-run` should behave.
- **Options:**
  - **A.** Fail the dry run with the same error as the real execution.
  - **B.** Produce a preview including environment details but emit a warning that execution will be blocked.
  - **C.** Skip environment output in dry run and explain that creation is disallowed until the project exists.
  - **D.** Allow preview and automatically downgrade the subsequent real run to project-only creation.

- Answer: A

## Q24. Environment Field Precedence
- **Description:** Clarify how `--env-color`, `--env-file`, and `--env-mode` flags interact with values provided in the `environment` object from `--cli-input` (Answer: Q11=D, Q16=A).
- **Options:**
  - **A.** CLI flags override config values for the corresponding fields; config provides defaults.
  - **B.** Config values override CLI flags to preserve declarative input.
  - **C.** Merge per field: positional/flags for name and mode, config for file/color when flags omitted.
  - **D.** Treat providing both as an error requiring the user to choose one source.

- Answer: A

## Q25. Force Overwrite Semantics
- **Description:** Define the exact side effects when `--force` overwrites an existing environment (Answers: Q6=A, Q20=A).
- **Options:**
  - **A.** Replace metadata and rewrite env vars file without backups.
  - **B.** Replace metadata but create a timestamped backup of the previous env vars file.
  - **C.** Prompt for confirmation unless `--yes` is supplied, then overwrite.
  - **D.** Only overwrite metadata; env vars file rotation must be handled separately by the user.

- Answer: A

## Q26. Stale Repository Entries
- **Description:** After treating a missing `.project.hcl` as non-existent (Answer: Q22=B), decide whether `pm new` should clean up stale registry metadata before proceeding with creation.
- **Options:**
  - **A.** Automatically remove the stale entry prior to creating the new project.
  - **B.** Leave the stale entry untouched and rely on `pm edit/delete` for cleanup.
  - **C.** Abort with instructions to resolve the stale entry manually before retrying.
  - **D.** Create the project but mark the stale entry for later reconciliation.

- Answer: B

## Q27. Uniqueness Checks with Stale Entries
- **Description:** Given Q26=B leaves stale registry entries in place, clarify how the shared uniqueness helper (Q8=D) should treat projects whose `.project.hcl` is missing.
- **Options:**
  - **A.** Ignore registry entries lacking `.project.hcl` when enforcing uniqueness.
  - **B.** Treat them as conflicts and block creation until the entry is cleaned up.
  - **C.** Attempt automatic cleanup during the uniqueness check before proceeding.
  - **D.** Warn the user but allow creation, marking the stale entry for review.

- Answer: A

## Q28. Environment Flags During Initial Creation
- **Description:** With Q13=B disallowing create+env in one run, define the behavior when users pass `--env-color`, `--env-file`, or `--env-mode` while the target project does not yet exist.
- **Options:**
  - **A.** Fail immediately with guidance to run project creation first, then re-run with the flags.
  - **B.** Ignore the environment flags and proceed with project creation only, emitting a warning.
  - **C.** Stage the environment configuration and prompt the user to rerun after creation completes.
  - **D.** Allow the run but skip applying env changes, requiring an explicit follow-up command.

- Answer: A

## Q29. Stale Entry User Feedback
- **Description:** When the uniqueness helper ignores stale registry entries (Q27=A) and creation proceeds, decide whether the CLI should still inform the user that a stale entry was detected.
- **Options:**
  - **A.** Emit a warning in both normal and dry-run modes describing the stale entry and suggesting cleanup.
  - **B.** Only log the stale entry at debug level; no user-facing output.
  - **C.** Surface the warning only when `--verbose` or `--output json` is used.
  - **D.** Record the issue in logs and return a non-zero exit code after successful creation.

- Answer: B

## Q30. Environment Variable Payload Source
- **Description:** Specify how environment-specific key/value pairs should be supplied now that the legacy `environments` map is deprecated (Q3=B) and the environment object is preferred.
- **Options:**
  - **A.** Add an `env_vars` map within the `environment` object; CLI flags cannot set individual keys.
  - **B.** Require users to provide a separate `--env-var KEY=VALUE` repeatable flag alongside the environment object.
  - **C.** Continue reading top-level `.env` style maps for backward compatibility but translate them into the environment object internally.
  - **D.** Do not support specifying env vars during shortcut creation; require manual editing after creation.

- Answer: B

## Q31. Config-Provided Environment Variables
- **Description:** With CLI env vars supplied via `--env-var` (Q30=B), define how configuration files should express environment key/value pairs.
- **Options:**
  - **A.** Allow an `env_vars` map inside the `environment` object when using `--cli-input`.
  - **B.** Require config files to list env vars under a new top-level `env_vars` array of `KEY=VALUE` strings.
  - **C.** Disallow env var definitions in config; users must rely on CLI flags or manual editing.
  - **D.** Support both `env_vars` map and repeatable `env_var` list, merging them with CLI flags according to Q24 precedence.

- Answer: A

## Q32. Env Var Precedence
- **Description:** Clarify how env vars from `--env-var` flags (Q30=B) interact with the `env_vars` map provided in config files (Q31=A).
- **Options:**
  - **A.** CLI `--env-var` flags override config-provided keys.
  - **B.** Config-provided `env_vars` override CLI flags to ensure reproducible declarative configs.
  - **C.** Treat duplicate keys as an error requiring user resolution.
  - **D.** Merge keys; if collisions occur, prefer the value with the latest timestamp (requires tracking order).

- Answer: A

## Q33. Registry Deduplication After Recreation
- **Description:** When a project is recreated after a stale entry (Q26=B, Q27=A, Q29=B), decide whether `pm new` should remove duplicate includeIf entries once the new `.project.hcl` is written.
- **Options:**
  - **A.** Automatically remove any stale includeIf entries pointing to the same path after successful creation.
  - **B.** Leave duplicate entries untouched and rely on separate maintenance commands.
  - **C.** Skip creating a new includeIf entry if one already exists for the path, even if previously stale.
  - **D.** Warn (via logs) about duplicates but leave resolution to the user.

- Answer: C

## Q34. Legacy Map Handling with `--allow-unknown`
- **Description:** With Q17=B failing fast on legacy `environments` maps, clarify whether the `--allow-unknown` flag can bypass that failure.
- **Options:**
  - **A.** `--allow-unknown` has no effect; legacy maps always cause an error.
  - **B.** Respect `--allow-unknown` by downgrading the failure to a warning.
  - **C.** Require an explicit `--allow-legacy-environments` flag to bypass the failure.
  - **D.** Ignore the legacy map entirely (no environments created) when `--allow-unknown` is set.

- Answer: D

## Q35. Simultaneous Config and Flag Env Vars
- **Description:** Given config `env_vars` maps (Q31=A) and CLI `--env-var` flags (Q30=B, Q32=A), decide how dry-run output should display combined env var sources.
- **Options:**
  - **A.** Show a merged map with CLI overrides applied, marking overridden keys.
  - **B.** List config-derived vars and CLI-derived vars separately in the output.
  - **C.** Display only the final merged values without source attribution.
  - **D.** Include source attribution only when `--verbose` or `--output json` is used.

- Answer: C

## Q36. Legacy Map Ignored Feedback
- **Description:** When `--allow-unknown` is present and legacy `environments` maps are ignored (Q34=D), clarify what feedback the CLI should give to the user.
- **Options:**
  - **A.** Emit a warning stating the legacy map was ignored and no environments were created.
  - **B.** Only note the ignored map in verbose/debug logs; keep default output silent.
  - **C.** Treat ignored legacy maps as informational events surfaced only in JSON output.
  - **D.** Fail the command unless `--quiet` is specified to explicitly acknowledge the loss of environment data.

- Answer: no feedback the legacy environments are consider unknown values

## Q37. Silent Handling Confirmation
- **Description:** The answer to Q36 specifies "no feedback" for ignored legacy environments, which is stricter than option B. Confirm the intended behavior.
- **Options:**
  - **A.** Absolutely no output or logging; legacy environments are discarded silently.
  - **B.** Silent in user-facing output but still logged at debug level for diagnostics.
  - **C.** Provide a single aggregated warning after processing all ignored legacy entries.
  - **D.** Require an explicit `--silent-legacy` flag to suppress warnings; otherwise warn.

- Answer: A