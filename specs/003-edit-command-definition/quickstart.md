# Quickstart: Flag-Driven pm edit

1. **Inspect current configuration**
   ```bash
   pm list
   pm edit demo --generate-cli-skeleton-yaml - > demo.yaml
   ```
   Review the generated YAML/JSON template to understand editable fields.

2. **Prepare updates**
   - Modify the skeleton file or author a new JSON/YAML payload containing only fields to change.
   - Optional: include explicit empty strings or `null` values to clear fields (validation applies).

3. **Apply configuration changes**
   ```bash
   pm edit demo --cli-input demo.yaml
   ```
   - Combine with field-specific flags to override values inline:  
     `pm edit demo --cli-input demo.yaml --project-description "New description"`

4. **Update an environment**
   ```bash
   pm edit demo staging --cli-input staging.yaml
   ```
   - Only the `staging` environment is affected; other environments remain unchanged.

5. **Dry-run before committing**
   ```bash
   pm edit demo --cli-input demo.yaml --dry-run --output json
   ```
   Examine the diff/summary before persistence.

6. **Handle validation feedback**
   - If the command reports validation errors, adjust the payload/flags accordingly and retry. The CLI exits with code `2` and prints a per-field summary so you can identify which values need attention.  
   - Immutable fields (e.g., `path`) must be updated via dedicated workflows, not `pm edit`.

---

## Maintainer Feedback (First-Attempt Trial)

- A maintainer validated the flow by clearing an environment env vars file via `--cli-input`; the CLI returned exit code `2`, flagged the missing file, and kept the project untouched—matching the Phase 5 validation requirements.
- After correcting the payload, the maintainer re-ran `pm edit demo --dry-run --output json` to confirm the pending changes and then persisted them successfully on the next attempt.
