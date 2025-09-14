# Data Model: TUI environments column

## Entities
- Project
  - name: string (unique)
  - environments: []string

- UIModel
  - projects: []Project
  - selectedProjectIdx: int
  - envsForSelected: []string
  - width: int
  - height: int
  - styles: struct{ left, right, title }

## Messages (contracts)
- ProjectsLoadedMsg{ projects []Project }
- ProjectSelectedMsg{ index int }
- WindowSizeMsg{ width, height int }

## Validation Rules
- selectedProjectIdx in [0, len(projects)) else envsForSelected = nil
- envsForSelected = projects[selectedProjectIdx].environments
- collapseRight = width < 60

## State Transitions
- ProjectsLoadedMsg → projects set, selectedProjectIdx=0, derive envsForSelected
- ProjectSelectedMsg → selectedProjectIdx updated, derive envsForSelected
- WindowSizeMsg → width/height set, affects collapseRight