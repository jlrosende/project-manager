## Why
Existing cli-new specification omits critical behaviors discovered during review: repeated runs should no-op once a project exists, environment creation must fail when the project is missing, `--allow-unknown` should be documented for silent acceptance, and precedence rules for config-only environment names need clarification.

## What Changes
- Document idempotent project creation behavior and required messaging for repeat runs.
- Clarify environment shortcut failure when the project is absent, and describe precedence when environment names originate from configuration files.
- Add coverage for the `--allow-unknown` flag explicitly enabling silent acceptance of unknown configuration keys.

## Impact
- Affected specs: `cli-new`
- Affected code: specification only (no code changes)
