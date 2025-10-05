package domain

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var projectNamePattern = regexp.MustCompile(`^[a-zA-Z0-9-_]+$`)

// ProjectValidationError represents a single validation failure on a specific field.
type ProjectValidationError struct {
	Field   string
	Message string
}

// ProjectValidationErrors aggregates multiple validation failures.
type ProjectValidationErrors []ProjectValidationError

// Error implements the error interface.
func (v ProjectValidationErrors) Error() string {
	if len(v) == 0 {
		return ""
	}

	parts := make([]string, 0, len(v))
	for _, err := range v {
		parts = append(parts, fmt.Sprintf("%s: %s", err.Field, err.Message))
	}

	return strings.Join(parts, "; ")
}

func (v *ProjectValidationErrors) add(field, message string) {
	*v = append(*v, ProjectValidationError{Field: field, Message: message})
}

// ValidateProjectDefinition ensures the supplied project definition satisfies the documented rules.
// It normalizes the definition in-place (e.g., trimming fields and resolving paths).
func ValidateProjectDefinition(def *ProjectDefinition) error {
	if def == nil {
		return errors.New("project definition is nil")
	}

	var errs ProjectValidationErrors

	def.Name = strings.TrimSpace(def.Name)
	if def.Name == "" {
		errs.add("name", "must be provided")
	} else if !projectNamePattern.MatchString(def.Name) {
		errs.add("name", "must match pattern [a-zA-Z0-9-_]+")
	}

	path := strings.TrimSpace(def.Path)
	def.Path = path

	if def.Here && path != "" {
		errs.add("destination", "cannot set both --here and path")
	}

	if !def.Here && path == "" {
		errs.add("destination", "must provide either --here or an absolute path")
	}

	if path != "" {
		abs, err := absPath(path)
		if err != nil {
			errs.add("path", fmt.Sprintf("invalid path: %v", err))
		} else if strings.TrimSpace(abs) == "" {
			errs.add("path", "resolved absolute path is empty")
		} else {
			def.Path = abs
		}
	}

	if len(errs) == 0 {
		return nil
	}

	return errs
}

func absPath(path string) (string, error) {
	if strings.HasPrefix(path, "~/") {
		home, err := os.UserHomeDir()
		if err == nil {
			path = filepath.Join(home, path[2:])
		}
	}

	return filepath.Abs(path)
}
