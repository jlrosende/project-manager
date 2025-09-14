# Contracts: UI Messages

## ProjectsLoadedMsg
- Input: projects []Project
- Effect: populate projects, select first, derive envs

## ProjectSelectedMsg
- Input: index int
- Effect: update selection, derive envs

## WindowSizeMsg
- Input: width int, height int
- Effect: update layout/collapse behavior