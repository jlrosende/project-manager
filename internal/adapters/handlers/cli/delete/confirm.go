package delete

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
	"strings"

	"golang.org/x/term"

	"github.com/jlrosende/project-manager/internal/adapters/handlers/tui"
	"github.com/jlrosende/project-manager/internal/core/domain"
)

// ConfirmationRequest captures the information presented to the user when the
// delete command requests confirmation.
type ConfirmationRequest struct {
	Target            domain.ProjectIdentifier
	Scope             domain.DeleteScope
	Force             bool
	DryRun            bool
	BackupDestination string
}

// ConfirmationFunc represents the function signature used to obtain user
// confirmation prior to executing a delete.
type ConfirmationFunc func(context.Context, ConfirmationRequest) (bool, error)

func defaultConfirmation(in io.Reader, out io.Writer, palette map[string]string) ConfirmationFunc {
	colors := palette
	if colors == nil {
		colors = tui.DefaultPalette()
	} else {
		colors = clonePalette(colors)
	}

	return func(ctx context.Context, req ConfirmationRequest) (bool, error) {
		if ctx != nil {
			if err := ctx.Err(); err != nil {
				return false, err
			}
		}

		if confirmed, handled, err := tryConfirmWithTUI(ctx, req, in, out, colors); handled {
			return confirmed, err
		}

		return promptYesNo(in, out, req)
	}
}

func tryConfirmWithTUI(
	ctx context.Context,
	req ConfirmationRequest,
	in io.Reader,
	out io.Writer,
	palette map[string]string,
) (bool, bool, error) {
	stdin, okIn := in.(*os.File)

	stdout, okOut := out.(*os.File)
	if !okIn || !okOut {
		return false, false, nil
	}

	if !term.IsTerminal(int(stdin.Fd())) || !term.IsTerminal(int(stdout.Fd())) {
		return false, false, nil
	}

	confirmed, err := confirmWithModal(ctx, req, stdin, stdout, palette)

	return confirmed, true, err
}

func promptYesNo(in io.Reader, out io.Writer, req ConfirmationRequest) (bool, error) {
	label := targetLabel(req.Target)
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

func clonePalette(src map[string]string) map[string]string {
	dup := make(map[string]string, len(src))
	for k, v := range src {
		dup[k] = v
	}

	return dup
}

func targetLabel(identifier domain.ProjectIdentifier) string {
	label := strings.TrimSpace(identifier.Name)
	if label != "" {
		return label
	}

	label = strings.TrimSpace(identifier.Path)
	if label != "" {
		return label
	}

	return "project"
}
