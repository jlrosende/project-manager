## Why
The cli-edit specification does not document expected behavior for dry-run/preview flags or the confirmation messaging that distinguishes project-level versus environment-specific updates. These details were highlighted during the recent review and should be captured before implementation work begins.

## What Changes
- Add requirements detailing how `pm edit` communicates the scope of successful edits in its confirmation output.
- Describe dry-run behavior for the edit command, ensuring it validates without persisting changes and reports planned updates.

## Impact
- Affected specs: `cli-edit`
- Affected code: specification only
