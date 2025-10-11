package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"gopkg.in/yaml.v3"

	"github.com/jlrosende/project-manager/internal/core/domain"
)

var (
	projectEditFieldOrder = []string{
		domain.ProjectFieldDescription,
		domain.ProjectFieldShell,
		domain.ProjectFieldEnvVarsFile,
		domain.ProjectFieldDefaultEnv,
	}

	projectEditFieldLabels = map[string]string{
		domain.ProjectFieldDescription: "Description",
		domain.ProjectFieldShell:       "Shell",
		domain.ProjectFieldEnvVarsFile: "Env Vars File",
		domain.ProjectFieldDefaultEnv:  "Default Environment",
	}

	environmentEditFieldOrder = []string{
		domain.EnvironmentFieldColor,
		domain.EnvironmentFieldEnvVarsMode,
		domain.EnvironmentFieldEnvVarsFile,
	}

	environmentEditFieldLabels = map[string]string{
		domain.EnvironmentFieldColor:       "Color",
		domain.EnvironmentFieldEnvVarsMode: "Env Vars Mode",
		domain.EnvironmentFieldEnvVarsFile: "Env Vars File",
	}

	projectEditAllowedKeys = map[string]struct{}{
		"description":   {},
		"shell":         {},
		"env_vars_file": {},
		"default_env":   {},
	}

	environmentEditAllowedKeys = map[string]struct{}{
		"color":         {},
		"env_vars_mode": {},
		"env_vars_file": {},
	}

	projectEditFieldMap = map[string]string{
		"description":   domain.ProjectFieldDescription,
		"shell":         domain.ProjectFieldShell,
		"env_vars_file": domain.ProjectFieldEnvVarsFile,
		"default_env":   domain.ProjectFieldDefaultEnv,
	}

	environmentEditFieldMap = map[string]string{
		"color":         domain.EnvironmentFieldColor,
		"env_vars_mode": domain.EnvironmentFieldEnvVarsMode,
		"env_vars_file": domain.EnvironmentFieldEnvVarsFile,
	}
)

// ErrProjectEditNoChanges indicates the caller attempted to save an edit without providing any mutations.
var ErrProjectEditNoChanges = errors.New("no edits supplied")

// ErrEnvironmentNotFound indicates the requested environment could not be located.
var ErrEnvironmentNotFound = errors.New("environment not found")

// ProjectEditFlags captures CLI flag overrides for project-level edits.
type ProjectEditFlags struct {
	Description    string
	DescriptionSet bool
	Shell          string
	ShellSet       bool
	EnvVarsFile    string
	EnvVarsFileSet bool
	DefaultEnv     string
	DefaultEnvSet  bool
}

// EnvironmentEditFlags captures CLI flag overrides for environment-level edits.
type EnvironmentEditFlags struct {
	Color       string
	ColorSet    bool
	EnvVarsMode string
	ModeSet     bool
	EnvVarsFile string
	EnvFileSet  bool
}

// ProjectEditOptions aggregates all inputs required to perform a project or environment edit.
type ProjectEditOptions struct {
	Name            string
	EnvironmentName string
	AllowUnknown    bool
	AllowUnknownEnv bool
	CLIInputPath    string
	DryRun          bool
	Flags           ProjectEditFlags
	EnvFlags        EnvironmentEditFlags
}

// ProjectEditChange describes a single field mutation applied (or previewed) during an edit.
type ProjectEditChange struct {
	Field string `json:"field"`
	Label string `json:"label"`
	Old   string `json:"old"`
	New   string `json:"new"`
}

// ProjectEditResult conveys the outcome of a project or environment edit request.
type ProjectEditResult struct {
	ProjectName     string              `json:"project"`
	EnvironmentName string              `json:"environment,omitempty"`
	Scope           string              `json:"scope"`
	DryRun          bool                `json:"dry_run"`
	Changes         []ProjectEditChange `json:"changes"`
	Errors          []ProjectEditChange `json:"errors,omitempty"`
}

// ProjectEditValidationError wraps validation failures encountered during a project edit.
type ProjectEditValidationError struct {
	result *ProjectEditResult
	errs   domain.ProjectValidationErrors
}

// Error implements the error interface.
func (e *ProjectEditValidationError) Error() string {
	if e == nil {
		return "project edit validation failed"
	}

	if len(e.errs) == 0 {
		return "project edit validation failed"
	}

	return e.errs.Error()
}

// Unwrap exposes the aggregated validation errors for use with errors.As.
func (e *ProjectEditValidationError) Unwrap() error {
	if e == nil {
		return nil
	}

	if len(e.errs) == 0 {
		return nil
	}

	return e.errs
}

// Result returns the structured edit result associated with the validation failure.
func (e *ProjectEditValidationError) Result() *ProjectEditResult {
	if e == nil {
		return nil
	}

	return e.result
}

// ValidationErrors exposes a copy of the aggregated validation failures.
func (e *ProjectEditValidationError) ValidationErrors() domain.ProjectValidationErrors {
	if e == nil {
		return nil
	}

	dup := make(domain.ProjectValidationErrors, len(e.errs))
	copy(dup, e.errs)

	return dup
}

// RenderProjectEditSkeleton serializes the editable fields for the requested scope in the given format.
func (svc *ProjectService) RenderProjectEditSkeleton(ctx context.Context, projectName, environmentName string, format SkeletonFormat) ([]byte, error) {
	_, project, err := svc.loadProjectForEdit(ctx, projectName)
	if err != nil {
		return nil, err
	}

	envName := strings.TrimSpace(environmentName)
	var payload map[string]any

	if envName == "" {
		payload = map[string]any{
			"project": map[string]any{
				"description":   project.Description,
				"shell":         project.Shell,
				"env_vars_file": project.EnvVarsFile,
				"default_env":   project.DefaultEnv,
			},
		}
	} else {
		env := environmentByName(project, envName)
		if env == nil {
			return nil, fmt.Errorf("environment %q not found", envName)
		}

		payload = map[string]any{
			"environment": map[string]any{
				"color":         env.Color,
				"env_vars_mode": env.EnvVarsMode,
				"env_vars_file": env.EnvVarsFile,
			},
		}
	}

	switch format {
	case SkeletonFormatJSON:
		data, marshalErr := json.MarshalIndent(payload, "", "  ")
		if marshalErr != nil {
			return nil, fmt.Errorf("render edit skeleton: %w", marshalErr)
		}

		if len(data) == 0 || data[len(data)-1] != '\n' {
			data = append(data, '\n')
		}

		return data, nil
	case SkeletonFormatYAML:
		data, marshalErr := yaml.Marshal(payload)
		if marshalErr != nil {
			return nil, fmt.Errorf("render edit skeleton: %w", marshalErr)
		}

		if len(data) == 0 || data[len(data)-1] != '\n' {
			data = append(data, '\n')
		}

		return data, nil
	default:
		return nil, fmt.Errorf("unsupported skeleton format: %s", format)
	}
}

// GenerateProjectEditSkeleton writes an edit skeleton to the provided destination path.
func (svc *ProjectService) GenerateProjectEditSkeleton(ctx context.Context, projectName, environmentName, path string, format SkeletonFormat) error {
	path = strings.TrimSpace(path)
	if path == "" {
		return errors.New("skeleton output path is empty")
	}

	if info, err := os.Stat(path); err == nil {
		if info.IsDir() {
			return fmt.Errorf("%s is a directory", path)
		}

		return fmt.Errorf("file %s already exists", path)
	} else if !errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("check skeleton destination: %w", err)
	}

	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return fmt.Errorf("prepare skeleton destination: %w", err)
	}

	data, err := svc.RenderProjectEditSkeleton(ctx, projectName, environmentName, format)
	if err != nil {
		return err
	}

	if err := os.WriteFile(path, data, 0o600); err != nil {
		return fmt.Errorf("write skeleton file: %w", err)
	}

	return nil
}

// ProjectEdit performs a project or environment edit by applying the supplied change set, optionally as a dry-run.
func (svc *ProjectService) ProjectEdit(ctx context.Context, opts ProjectEditOptions) (*ProjectEditResult, error) {
	name := strings.TrimSpace(opts.Name)
	if name == "" {
		return nil, errors.New("project name is required")
	}

	identifier, current, err := svc.loadProjectForEdit(ctx, name)
	if err != nil {
		return nil, err
	}

	envName := strings.TrimSpace(opts.EnvironmentName)
	if envName == "" {
		return svc.projectScopeEdit(ctx, identifier, current, opts)
	}

	return svc.environmentScopeEdit(ctx, identifier, current, envName, opts)
}

func (svc *ProjectService) projectScopeEdit(
	ctx context.Context,
	identifier domain.ProjectIdentifier,
	current *domain.Project,
	opts ProjectEditOptions,
) (*ProjectEditResult, error) {
	if opts.EnvFlags.ColorSet || opts.EnvFlags.ModeSet || opts.EnvFlags.EnvFileSet {
		return nil, errors.New("environment flags require an environment name")
	}

	mutations := map[string]*domain.FieldMutation{}
	var validationErrors domain.ProjectValidationErrors

	if path := strings.TrimSpace(opts.CLIInputPath); path != "" {
		fileMutations, unknown, loadErr := loadProjectEditFile(path)
		if loadErr != nil {
			return nil, loadErr
		}

		if len(unknown) > 0 && !opts.AllowUnknown {
			sort.Strings(unknown)
			return nil, fmt.Errorf(
				"unknown project fields: %s (use --allow-unknown to ignore)",
				strings.Join(unknown, ", "),
			)
		}

		for field, mutation := range fileMutations {
			mutations[field] = mutation
		}
	}

	applyFlag := func(set bool, value string, field string) {
		if !set {
			return
		}

		mutations[field] = &domain.FieldMutation{
			Value:  value,
			Source: domain.ChangeSourceFlag,
		}
	}

	applyFlag(opts.Flags.DescriptionSet, opts.Flags.Description, domain.ProjectFieldDescription)
	applyFlag(opts.Flags.ShellSet, opts.Flags.Shell, domain.ProjectFieldShell)
	applyFlag(opts.Flags.EnvVarsFileSet, opts.Flags.EnvVarsFile, domain.ProjectFieldEnvVarsFile)
	applyFlag(opts.Flags.DefaultEnvSet, opts.Flags.DefaultEnv, domain.ProjectFieldDefaultEnv)

	if len(mutations) == 0 {
		return nil, ErrProjectEditNoChanges
	}

	if err := pruneUnchangedProjectMutations(current, mutations); err != nil {
		return nil, err
	}

	if len(mutations) == 0 {
		return nil, ErrProjectEditNoChanges
	}

	preview := cloneProject(current)
	if err := applyProjectMutationsInMemory(preview, mutations); err != nil {
		return nil, err
	}

	if errs, err := validateProjectEdit(current, preview, mutations); err != nil {
		return nil, err
	} else if len(errs) > 0 {
		validationErrors = append(validationErrors, errs...)
	}

	if errs := validateProjectClearances(current, preview, mutations); len(errs) > 0 {
		validationErrors = append(validationErrors, errs...)
	}

	result := &ProjectEditResult{
		ProjectName: identifier.Name,
		Scope:       "project",
	}

	if len(validationErrors) > 0 {
		result.Errors = formatValidationErrors(validationErrors)
		return result, newProjectEditValidationError(result, validationErrors)
	}

	changeSet := &domain.EditChangeSet{
		Project: &domain.ProjectChangeSet{Fields: mutations},
	}

	lock, err := svc.project.AcquireEditLock(ctx, identifier)
	if err != nil {
		return nil, err
	}
	defer svc.releaseEditLock(lock, identifier.Name)

	snapshot := cloneProject(current)
	if err := applyProjectMutationsInMemory(snapshot, mutations); err != nil {
		return nil, err
	}

	summary := summarizeProjectChanges(current, snapshot, mutations)

	if opts.DryRun {
		svc.logInfo(
			"project edit preview",
			field("name", identifier.Name),
			field("changes", len(summary)),
		)

		result.DryRun = true
		result.Changes = summary
		return result, nil
	}

	if _, err := svc.project.ApplyEditChangeSet(ctx, identifier, changeSet); err != nil {
		return nil, err
	}

	svc.logInfo(
		"project edit applied",
		field("name", identifier.Name),
		field("changes", len(summary)),
	)

	result.DryRun = false
	result.Changes = summary
	return result, nil
}

func (svc *ProjectService) environmentScopeEdit(
	ctx context.Context,
	identifier domain.ProjectIdentifier,
	current *domain.Project,
	envName string,
	opts ProjectEditOptions,
) (*ProjectEditResult, error) {
	if opts.Flags.DescriptionSet || opts.Flags.ShellSet || opts.Flags.EnvVarsFileSet || opts.Flags.DefaultEnvSet {
		return nil, errors.New("project flags cannot be used when editing an environment")
	}

	env := environmentByName(current, envName)
	if env == nil {
		return nil, fmt.Errorf("%w: %s", ErrEnvironmentNotFound, envName)
	}

	mutations := map[string]*domain.FieldMutation{}
	var validationErrors domain.ProjectValidationErrors

	if path := strings.TrimSpace(opts.CLIInputPath); path != "" {
		fileMutations, unknown, loadErr := loadEnvironmentEditFile(path)
		if loadErr != nil {
			return nil, loadErr
		}

		if len(unknown) > 0 && !opts.AllowUnknownEnv {
			sort.Strings(unknown)
			return nil, fmt.Errorf(
				"unknown environment fields: %s (use --allow-unknown to ignore)",
				strings.Join(unknown, ", "),
			)
		}

		for field, mutation := range fileMutations {
			mutations[field] = mutation
		}
	}

	applyEnvFlag := func(set bool, value string, field string) {
		if !set {
			return
		}

		mutations[field] = &domain.FieldMutation{
			Value:  value,
			Source: domain.ChangeSourceFlag,
		}
	}

	applyEnvFlag(opts.EnvFlags.ColorSet, opts.EnvFlags.Color, domain.EnvironmentFieldColor)
	applyEnvFlag(opts.EnvFlags.ModeSet, opts.EnvFlags.EnvVarsMode, domain.EnvironmentFieldEnvVarsMode)
	applyEnvFlag(opts.EnvFlags.EnvFileSet, opts.EnvFlags.EnvVarsFile, domain.EnvironmentFieldEnvVarsFile)

	if len(mutations) == 0 {
		return nil, ErrProjectEditNoChanges
	}

	if err := pruneUnchangedEnvironmentMutations(env, mutations); err != nil {
		return nil, err
	}

	if len(mutations) == 0 {
		return nil, ErrProjectEditNoChanges
	}

	preview := cloneProject(current)
	if err := applyEnvironmentMutationsInMemory(preview, envName, mutations); err != nil {
		return nil, err
	}

	if errs, err := validateEnvironmentEdit(current, preview, envName, mutations); err != nil {
		return nil, err
	} else if len(errs) > 0 {
		validationErrors = append(validationErrors, errs...)
	}

	if errs := validateEnvironmentClearances(env, environmentByName(preview, envName), mutations); len(errs) > 0 {
		validationErrors = append(validationErrors, errs...)
	}

	envResult := &ProjectEditResult{
		ProjectName:     identifier.Name,
		EnvironmentName: envName,
		Scope:           "environment",
	}

	if len(validationErrors) > 0 {
		envResult.Errors = formatValidationErrors(validationErrors)
		return envResult, newProjectEditValidationError(envResult, validationErrors)
	}

	changeSet := &domain.EditChangeSet{
		Environment: &domain.EnvironmentChangeSet{Name: envName, Fields: mutations},
	}

	lock, err := svc.project.AcquireEditLock(ctx, identifier)
	if err != nil {
		return nil, err
	}
	defer svc.releaseEditLock(lock, identifier.Name)

	snapshot := cloneProject(current)
	if err := applyEnvironmentMutationsInMemory(snapshot, envName, mutations); err != nil {
		return nil, err
	}

	summary := summarizeEnvironmentChanges(env, environmentByName(snapshot, envName), mutations)

	if opts.DryRun {
		svc.logInfo(
			"environment edit preview",
			field("name", identifier.Name),
			field("environment", envName),
			field("changes", len(summary)),
		)

		envResult.DryRun = true
		envResult.Changes = summary
		return envResult, nil
	}

	if _, err := svc.project.ApplyEditChangeSet(ctx, identifier, changeSet); err != nil {
		return nil, err
	}

	svc.logInfo(
		"environment edit applied",
		field("name", identifier.Name),
		field("environment", envName),
		field("changes", len(summary)),
	)

	envResult.DryRun = false
	envResult.Changes = summary
	return envResult, nil
}

func (svc *ProjectService) releaseEditLock(lock *domain.ProjectLock, projectName string) {
	if lock == nil {
		return
	}

	if err := lock.Release(); err != nil {
		svc.logWarn(
			"release project edit lock failed",
			field("name", projectName),
			field("err", err),
		)
	}
}

func (svc *ProjectService) loadProjectForEdit(
	ctx context.Context,
	name string,
) (domain.ProjectIdentifier, *domain.Project, error) {
	identifier, err := svc.project.ResolveIdentifier(ctx, domain.ProjectIdentifier{Name: name})
	if err != nil {
		svc.logError("resolve project for edit failed", field("name", name), field("err", err))
		return domain.ProjectIdentifier{}, nil, err
	}

	project, err := svc.project.LoadProjectDefinition(ctx, identifier)
	if err != nil {
		svc.logError("load project definition failed", field("name", identifier.Name), field("err", err))
		return domain.ProjectIdentifier{}, nil, err
	}

	return identifier, project, nil
}

func loadProjectEditFile(path string) (map[string]*domain.FieldMutation, []string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, nil, fmt.Errorf("read CLI input: %w", err)
	}

	if len(data) == 0 {
		return nil, nil, fmt.Errorf("cli input %s is empty", path)
	}

	var raw map[string]any
	if err := yaml.Unmarshal(data, &raw); err != nil {
		return nil, nil, fmt.Errorf("parse CLI input %s: %w", path, err)
	}

	fields, unknown := extractProjectEditFields(raw)
	mutations := make(map[string]*domain.FieldMutation, len(fields))

	for key, value := range fields {
		fieldName, ok := projectEditFieldMap[key]
		if !ok {
			continue
		}

		mutation, err := mutationFromInterface(value, domain.ChangeSourceInput)
		if err != nil {
			return nil, nil, fmt.Errorf("project.%s: %w", key, err)
		}

		mutations[fieldName] = mutation
	}

	return mutations, unknown, nil
}

func loadEnvironmentEditFile(path string) (map[string]*domain.FieldMutation, []string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, nil, fmt.Errorf("read CLI input: %w", err)
	}

	if len(data) == 0 {
		return nil, nil, fmt.Errorf("cli input %s is empty", path)
	}

	var raw map[string]any
	if err := yaml.Unmarshal(data, &raw); err != nil {
		return nil, nil, fmt.Errorf("parse CLI input %s: %w", path, err)
	}

	fields, unknown := extractEnvironmentEditFields(raw)
	mutations := make(map[string]*domain.FieldMutation, len(fields))

	for key, value := range fields {
		fieldName, ok := environmentEditFieldMap[key]
		if !ok {
			continue
		}

		mutation, err := mutationFromInterface(value, domain.ChangeSourceInput)
		if err != nil {
			return nil, nil, fmt.Errorf("environment.%s: %w", key, err)
		}

		mutations[fieldName] = mutation
	}

	return mutations, unknown, nil
}

func newProjectEditValidationError(result *ProjectEditResult, errs domain.ProjectValidationErrors) error {
	if len(errs) == 0 {
		return nil
	}

	dup := make(domain.ProjectValidationErrors, len(errs))
	copy(dup, errs)

	return &ProjectEditValidationError{
		result: result,
		errs:   dup,
	}
}

func pruneUnchangedProjectMutations(current *domain.Project, mutations map[string]*domain.FieldMutation) error {
	for field, mutation := range mutations {
		if mutation == nil {
			delete(mutations, field)
			continue
		}

		value, err := mutation.StringValue()
		if err != nil {
			return fmt.Errorf("project mutation %s: %w", field, err)
		}

		normalized := normalizeProjectFieldValue(field, value)
		if normalizeProjectFieldValue(field, projectFieldValue(current, field)) == normalized {
			delete(mutations, field)
		}
	}

	return nil
}

func pruneUnchangedEnvironmentMutations(env *domain.Environment, mutations map[string]*domain.FieldMutation) error {
	for field, mutation := range mutations {
		if mutation == nil {
			delete(mutations, field)
			continue
		}

		value, err := mutation.StringValue()
		if err != nil {
			return fmt.Errorf("environment mutation %s: %w", field, err)
		}

		normalized := normalizeEnvironmentFieldValue(field, value)
		if normalizeEnvironmentFieldValue(field, environmentFieldValue(env, field)) == normalized {
			delete(mutations, field)
		}
	}

	return nil
}

func normalizeProjectFieldValue(field, value string) string {
	switch field {
	case domain.ProjectFieldShell, domain.ProjectFieldEnvVarsFile, domain.ProjectFieldDefaultEnv:
		return strings.TrimSpace(value)
	default:
		return value
	}
}

func normalizeEnvironmentFieldValue(field, value string) string {
	switch field {
	case domain.EnvironmentFieldColor:
		return strings.TrimSpace(value)
	case domain.EnvironmentFieldEnvVarsMode:
		return strings.TrimSpace(strings.ToLower(value))
	case domain.EnvironmentFieldEnvVarsFile:
		return strings.TrimSpace(value)
	default:
		return value
	}
}

func cloneProject(in *domain.Project) *domain.Project {
	if in == nil {
		return nil
	}

	out := *in
	if len(in.Environments) > 0 {
		out.Environments = make([]*domain.Environment, len(in.Environments))
		for i, env := range in.Environments {
			if env == nil {
				continue
			}

			copyEnv := *env
			out.Environments[i] = &copyEnv
		}
	}

	return &out
}

func environmentByName(project *domain.Project, name string) *domain.Environment {
	if project == nil {
		return nil
	}

	for _, env := range project.Environments {
		if env != nil && env.Name == name {
			return env
		}
	}

	return nil
}

func applyProjectMutationsInMemory(project *domain.Project, fields map[string]*domain.FieldMutation) error {
	if project == nil {
		return errors.New("project is nil")
	}

	for field, mutation := range fields {
		if mutation == nil {
			continue
		}

		value, err := mutation.StringValue()
		if err != nil {
			return fmt.Errorf("project mutation %s: %w", field, err)
		}

		normalized := normalizeProjectFieldValue(field, value)

		switch field {
		case domain.ProjectFieldDescription:
			project.Description = normalized
		case domain.ProjectFieldShell:
			project.Shell = normalized
		case domain.ProjectFieldEnvVarsFile:
			project.EnvVarsFile = normalized
		case domain.ProjectFieldDefaultEnv:
			project.DefaultEnv = normalized
		default:
			return fmt.Errorf("unsupported project field mutation %q", field)
		}
	}

	return nil
}

func applyEnvironmentMutationsInMemory(project *domain.Project, envName string, fields map[string]*domain.FieldMutation) error {
	if project == nil {
		return errors.New("project is nil")
	}

	env := environmentByName(project, envName)
	if env == nil {
		return fmt.Errorf("environment %s not found", envName)
	}

	for field, mutation := range fields {
		if mutation == nil {
			continue
		}

		value, err := mutation.StringValue()
		if err != nil {
			return fmt.Errorf("environment mutation %s: %w", field, err)
		}

		normalized := normalizeEnvironmentFieldValue(field, value)

		switch field {
		case domain.EnvironmentFieldColor:
			env.Color = normalized
		case domain.EnvironmentFieldEnvVarsMode:
			if normalized == "" {
				env.EnvVarsMode = domain.EnvVarsModeMerge
			} else {
				env.EnvVarsMode = normalized
			}
		case domain.EnvironmentFieldEnvVarsFile:
			env.EnvVarsFile = normalized
		default:
			return fmt.Errorf("unsupported environment field mutation %q", field)
		}
	}

	return nil
}

func projectFieldValue(project *domain.Project, field string) string {
	if project == nil {
		return ""
	}

	switch field {
	case domain.ProjectFieldDescription:
		return project.Description
	case domain.ProjectFieldShell:
		return project.Shell
	case domain.ProjectFieldEnvVarsFile:
		return project.EnvVarsFile
	case domain.ProjectFieldDefaultEnv:
		return project.DefaultEnv
	default:
		return ""
	}
}

func environmentFieldValue(env *domain.Environment, field string) string {
	if env == nil {
		return ""
	}

	switch field {
	case domain.EnvironmentFieldColor:
		return env.Color
	case domain.EnvironmentFieldEnvVarsMode:
		return env.EnvVarsMode
	case domain.EnvironmentFieldEnvVarsFile:
		return env.EnvVarsFile
	default:
		return ""
	}
}

func summarizeProjectChanges(before, after *domain.Project, fields map[string]*domain.FieldMutation) []ProjectEditChange {
	if before == nil || after == nil || len(fields) == 0 {
		return nil
	}

	changes := make([]ProjectEditChange, 0, len(fields))
	for _, field := range projectEditFieldOrder {
		if _, ok := fields[field]; !ok {
			continue
		}

		oldVal := projectFieldValue(before, field)
		newVal := projectFieldValue(after, field)
		if oldVal == newVal {
			continue
		}

		changes = append(changes, ProjectEditChange{
			Field: field,
			Label: projectEditFieldLabels[field],
			Old:   oldVal,
			New:   newVal,
		})
	}

	return changes
}

func summarizeEnvironmentChanges(before *domain.Environment, after *domain.Environment, fields map[string]*domain.FieldMutation) []ProjectEditChange {
	if before == nil || after == nil || len(fields) == 0 {
		return nil
	}

	changes := make([]ProjectEditChange, 0, len(fields))
	for _, field := range environmentEditFieldOrder {
		if _, ok := fields[field]; !ok {
			continue
		}

		oldVal := environmentFieldValue(before, field)
		newVal := environmentFieldValue(after, field)
		if oldVal == newVal {
			continue
		}

		changes = append(changes, ProjectEditChange{
			Field: field,
			Label: environmentEditFieldLabels[field],
			Old:   oldVal,
			New:   newVal,
		})
	}

	return changes
}

func validateProjectEdit(current, updated *domain.Project, mutations map[string]*domain.FieldMutation) (domain.ProjectValidationErrors, error) {
	var errs domain.ProjectValidationErrors

	if current == nil || updated == nil {
		return nil, errors.New("project is nil")
	}

	if _, ok := mutations[domain.ProjectFieldEnvVarsFile]; ok {
		if ok, _, err := ensurePathWithinProject(current.Path, updated.EnvVarsFile); err != nil {
			return nil, fmt.Errorf("validate project env vars file: %w", err)
		} else if !ok {
			errs = append(errs, domain.ProjectValidationError{
				Field:   domain.ProjectFieldEnvVarsFile,
				Message: fmt.Sprintf("env vars file must be within project root (%s)", current.Path),
			})
		}
	}

	if _, ok := mutations[domain.ProjectFieldDefaultEnv]; ok {
		value := strings.TrimSpace(updated.DefaultEnv)
		if value != "" && environmentByName(updated, value) == nil {
			errs = append(errs, domain.ProjectValidationError{
				Field:   domain.ProjectFieldDefaultEnv,
				Message: fmt.Sprintf("environment %q does not exist", value),
			})
		}
	}

	return errs, nil
}

func validateProjectClearances(_ *domain.Project, updated *domain.Project, mutations map[string]*domain.FieldMutation) domain.ProjectValidationErrors {
	var errs domain.ProjectValidationErrors

	if updated == nil {
		return errs
	}

	if _, ok := mutations[domain.ProjectFieldEnvVarsFile]; ok {
		if strings.TrimSpace(updated.EnvVarsFile) == "" {
			errs = append(errs, domain.ProjectValidationError{
				Field:   domain.ProjectFieldEnvVarsFile,
				Message: "env vars file cannot be empty",
			})
		}
	}

	if _, ok := mutations[domain.ProjectFieldShell]; ok {
		if strings.TrimSpace(updated.Shell) == "" {
			errs = append(errs, domain.ProjectValidationError{
				Field:   domain.ProjectFieldShell,
				Message: "shell cannot be empty",
			})
		}
	}

	return errs
}

func validateEnvironmentEdit(current, updated *domain.Project, envName string, mutations map[string]*domain.FieldMutation) (domain.ProjectValidationErrors, error) {
	var errs domain.ProjectValidationErrors

	if current == nil || updated == nil {
		return nil, errors.New("project is nil")
	}

	env := environmentByName(updated, envName)
	if env == nil {
		return nil, fmt.Errorf("environment %s not found", envName)
	}

	if _, ok := mutations[domain.EnvironmentFieldEnvVarsMode]; ok {
		mode := strings.TrimSpace(env.EnvVarsMode)
		if mode == "" {
			errs = append(errs, domain.ProjectValidationError{
				Field:   domain.EnvironmentFieldEnvVarsMode,
				Message: "env vars mode cannot be empty",
			})
		} else if mode != domain.EnvVarsModeMerge && mode != domain.EnvVarsModeReplace {
			errs = append(errs, domain.ProjectValidationError{
				Field:   domain.EnvironmentFieldEnvVarsMode,
				Message: fmt.Sprintf("env vars mode must be %q or %q", domain.EnvVarsModeMerge, domain.EnvVarsModeReplace),
			})
		}
	}

	if _, ok := mutations[domain.EnvironmentFieldEnvVarsFile]; ok {
		if ok, _, err := ensurePathWithinProject(current.Path, env.EnvVarsFile); err != nil {
			return nil, fmt.Errorf("validate environment env vars file: %w", err)
		} else if !ok {
			errs = append(errs, domain.ProjectValidationError{
				Field:   domain.EnvironmentFieldEnvVarsFile,
				Message: fmt.Sprintf("env vars file must be within project root (%s)", current.Path),
			})
		}
	}

	return errs, nil
}

func validateEnvironmentClearances(_ *domain.Environment, updated *domain.Environment, mutations map[string]*domain.FieldMutation) domain.ProjectValidationErrors {
	var errs domain.ProjectValidationErrors

	if updated == nil {
		return errs
	}

	if _, ok := mutations[domain.EnvironmentFieldEnvVarsFile]; ok {
		if strings.TrimSpace(updated.EnvVarsFile) == "" {
			errs = append(errs, domain.ProjectValidationError{
				Field:   domain.EnvironmentFieldEnvVarsFile,
				Message: "env vars file cannot be empty",
			})
		}
	}

	if _, ok := mutations[domain.EnvironmentFieldEnvVarsMode]; ok {
		if strings.TrimSpace(updated.EnvVarsMode) == "" {
			errs = append(errs, domain.ProjectValidationError{
				Field:   domain.EnvironmentFieldEnvVarsMode,
				Message: "env vars mode cannot be empty",
			})
		}
	}

	return errs
}

func formatValidationErrors(errs domain.ProjectValidationErrors) []ProjectEditChange {
	if len(errs) == 0 {
		return nil
	}

	formatted := make([]ProjectEditChange, 0, len(errs))
	for _, err := range errs {
		formatted = append(formatted, ProjectEditChange{
			Field: err.Field,
			Label: validationLabel(err.Field),
			Old:   "",
			New:   err.Message,
		})
	}

	return formatted
}

func validationLabel(field string) string {
	if label, ok := projectEditFieldLabels[field]; ok {
		return label
	}

	if label, ok := environmentEditFieldLabels[field]; ok {
		return label
	}

	return field
}

func extractProjectEditFields(raw map[string]any) (map[string]any, []string) {
	fields := make(map[string]any)
	unknown := []string{}

	if raw == nil {
		return fields, unknown
	}

	for key, value := range raw {
		trimmed := strings.TrimSpace(key)
		if trimmed == "" {
			continue
		}

		switch trimmed {
		case "project":
			obj, ok := value.(map[string]any)
			if !ok {
				unknown = append(unknown, trimmed)
				continue
			}

			for fieldKey, fieldValue := range obj {
				k := strings.TrimSpace(fieldKey)
				if k == "" {
					continue
				}

				if _, allowed := projectEditAllowedKeys[k]; !allowed {
					unknown = append(unknown, "project."+k)
					continue
				}

				fields[k] = fieldValue
			}
		default:
			unknown = append(unknown, trimmed)
		}
	}

	return fields, unknown
}

func extractEnvironmentEditFields(raw map[string]any) (map[string]any, []string) {
	fields := make(map[string]any)
	unknown := []string{}

	if raw == nil {
		return fields, unknown
	}

	for key, value := range raw {
		trimmed := strings.TrimSpace(key)
		if trimmed == "" {
			continue
		}

		switch trimmed {
		case "environment":
			obj, ok := value.(map[string]any)
			if !ok {
				unknown = append(unknown, trimmed)
				continue
			}

			for fieldKey, fieldValue := range obj {
				k := strings.TrimSpace(fieldKey)
				if k == "" {
					continue
				}

				if _, allowed := environmentEditAllowedKeys[k]; !allowed {
					unknown = append(unknown, "environment."+k)
					continue
				}

				fields[k] = fieldValue
			}
		default:
			unknown = append(unknown, trimmed)
		}
	}

	return fields, unknown
}

func mutationFromInterface(value any, source domain.ChangeSource) (*domain.FieldMutation, error) {
	switch v := value.(type) {
	case nil:
		return &domain.FieldMutation{Value: nil, Source: source}, nil
	case string:
		return &domain.FieldMutation{Value: v, Source: source}, nil
	case fmt.Stringer:
		return &domain.FieldMutation{Value: v.String(), Source: source}, nil
	case []byte:
		return &domain.FieldMutation{Value: string(v), Source: source}, nil
	case int, int8, int16, int32, int64,
		uint, uint8, uint16, uint32, uint64,
		float32, float64,
		bool:
		return &domain.FieldMutation{Value: fmt.Sprint(v), Source: source}, nil
	default:
		return nil, fmt.Errorf("unsupported value type %T", value)
	}
}

func ensurePathWithinProject(root, candidate string) (bool, string, error) {
	root = strings.TrimSpace(root)
	if root == "" {
		return false, "", errors.New("project root is empty")
	}

	resolvedRoot := root
	if !filepath.IsAbs(resolvedRoot) {
		abs, err := filepath.Abs(resolvedRoot)
		if err != nil {
			return false, "", err
		}
		resolvedRoot = abs
	}
	resolvedRoot = filepath.Clean(resolvedRoot)

	candidate = strings.TrimSpace(candidate)
	if candidate == "" {
		return false, "", nil
	}

	resolvedCandidate := candidate
	if !filepath.IsAbs(resolvedCandidate) {
		resolvedCandidate = filepath.Join(resolvedRoot, resolvedCandidate)
	}
	resolvedCandidate = filepath.Clean(resolvedCandidate)

	rel, err := filepath.Rel(resolvedRoot, resolvedCandidate)
	if err != nil {
		return false, resolvedCandidate, err
	}

	if rel == "." {
		return true, resolvedCandidate, nil
	}

	if strings.HasPrefix(rel, ".."+string(filepath.Separator)) || rel == ".." {
		return false, resolvedCandidate, nil
	}

	return true, resolvedCandidate, nil
}
