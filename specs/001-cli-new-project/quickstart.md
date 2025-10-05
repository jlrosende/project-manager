# Quickstart: `pm new`

## Create at a path
```
pm new myproj /absolute/path/to/myproj
```

## Initialize current directory
```
pm new myproj --here
```

## Use a CLI input file (JSON or YAML)
```
pm new myproj --cli-input ./project.yaml
```

## Arguments override CLI input
```
pm new override /abs/path --cli-input ./project.yaml
```

## Dry run
```
pm new demo /abs/path --dry-run
```

## Force overwrite existing files
```
pm new demo /abs/path --force
```

## Allow unknown CLI input fields
```
pm new demo --cli-input ./project.yaml --allow-unknown
```

## Generate CLI input skeletons
```
pm new --generate-cli-skeleton-json ./project.json
pm new --generate-cli-skeleton-yaml ./project.yaml
```

Omit the path to print the skeleton to standard output:
```
pm new --generate-cli-skeleton-json
```
