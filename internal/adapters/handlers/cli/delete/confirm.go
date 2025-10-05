package delete

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/jlrosende/project-manager/internal/core/domain"
)

// ConfirmationRequest captures the information presented to the user when the
// delete command requests confirmation.
type ConfirmationRequest struct {
	Target domain.ProjectIdentifier
	Scope  domain.DeleteScope
	Force  bool
}

// ConfirmationFunc represents the function signature used to obtain user
// confirmation prior to executing a delete.
type ConfirmationFunc func(context.Context, ConfirmationRequest) (bool, error)

func defaultConfirmation(in io.Reader, out io.Writer) ConfirmationFunc {
	return func(ctx context.Context, req ConfirmationRequest) (bool, error) {
		if ctx != nil {
			if err := ctx.Err(); err != nil {
				return false, err
			}
		}

		label := strings.TrimSpace(req.Target.Name)
		if label == "" {
			label = strings.TrimSpace(req.Target.Path)
		}
		if label == "" {
			label = "project"
		}

		fmt.Fprintf(out, "Delete %s (%s)? (y/N): ", label, req.Scope.String())

		reader := bufio.NewReader(in)
		response, err := reader.ReadString('\n')
		if err != nil && err != io.EOF {
			return false, err
		}

		answer := strings.TrimSpace(strings.ToLower(response))
		if answer == "y" || answer == "yes" {
			return true, nil
		}

		return false, nil
	}
}
