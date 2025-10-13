# Quickstart: `pm new` Environment Enhancements

1. **Create a project (unchanged)**
   ```bash
   pm new demo-project ~/work/demo-project
   ```
   - Creates `.project.hcl` and registers the project. Environment options are ignored until the project exists.

2. **Add an environment with defaults**
   ```bash
   pm new demo-project staging
   ```
   - Adds a `staging` environment with file `./.staging.env`, merge mode, and no additional metadata.

3. **Customize environment via flags**
   ```bash
   pm new demo-project production \
     --env-file .env.production \
     --env-mode replace \
     --env-color 160 \
     --env-var API_URL=https://api.example.com \
     --env-var FEATURE_X=true
   ```
   - Overrides file/mode/color, merges env vars with command-line precedence.

4. **Use config input**
   ```yaml
   # env-config.yaml
   name: demo-project
   environment:
     name: staging
     env_vars_file: .env.staging
     env_vars_mode: merge
     color: 110
     env_vars:
       API_URL: https://staging.example.com
   ```
   ```bash
   pm new demo-project staging --cli-input env-config.yaml --env-var FEATURE_X=true
   ```
   - Config supplies base values; CLI flags win on conflicts (`FEATURE_X` added via CLI).

5. **Preview changes**
   ```bash
   pm new demo-project staging --dry-run --output json
   ```
   - Displays merged project + environment payload; no files are written.

6. **Handle stale registry entries**
   - If `.project.hcl` was removed manually, rerunning `pm new demo-project` recreates the project. The CLI ignores stale registry entries that lack the project file and avoids writing duplicate include rules.

7. **Legacy configs**
   - Config files using the old `environments` map fail fast unless `--allow-unknown` is set. When allowed, legacy data is ignored silently—update configs to the new `environment` object for full functionality.
