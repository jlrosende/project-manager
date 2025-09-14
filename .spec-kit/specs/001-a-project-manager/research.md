# Research: TUI environments column

## Decisions
- Use bubbles list for environments pane
- Split layout horizontally: left projects list, right environments list
- Update envs list on ProjectSelectedMsg in Update()
- Collapse environments pane if width < 60 cols
- Precompute lipgloss styles in model init

## Rationale
- Bubbles list integrates with Bubble Tea and matches existing patterns
- Horizontal split mirrors mental model: select project → see envs
- Update on selection ensures single source of truth in Elm model
- Collapse rule maintains usability on small terminals
- Precomputed styles reduce per-frame cost

## Alternatives Considered
- Table component: heavier API, unnecessary for simple names
- Vertical split with tabs: wastes vertical space; less discoverable
- Recomputing styles per View(): simpler code but slower rendering