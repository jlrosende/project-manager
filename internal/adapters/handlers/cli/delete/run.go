package delete

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	"github.com/jlrosende/project-manager/configs"
	"github.com/jlrosende/project-manager/internal/adapters/handlers/tui"
	"github.com/jlrosende/project-manager/internal/bootstrap"
	"github.com/jlrosende/project-manager/internal/core/domain"
	"github.com/jlrosende/project-manager/internal/core/services"
)

func run(cmd *cobra.Command, args []string) error {
	return runWithOptions(cmd, args, Options{})
}

func runWithOptions(cmd *cobra.Command, args []string, opts Options) error {
	ctx := cmd.Context()
	if ctx == nil {
		ctx = context.Background()
	}

	service := opts.Service
	if service == nil {
		container, err := bootstrap.New(bootstrap.Options{Logger: bootstrap.WrapSlog(slog.Default())})
		if err != nil {
			return fmt.Errorf("bootstrap delete command: %w", err)
		}

		service = container.ProjectService
	}

	palette := opts.Palette
	cfg := opts.Config

	if palette != nil {
		palette = clonePalette(palette)
	}

	if palette == nil || cfg == nil {
		resolvedPalette, resolvedCfg, err := resolvePalette(cmd)
		if err != nil {
			return err
		}

		if resolvedPalette == nil {
			resolvedPalette = tui.DefaultPalette()
		}

		if palette == nil {
			palette = clonePalette(resolvedPalette)
		}

		if cfg == nil {
			cfg = resolvedCfg
		}
	}

	if palette == nil {
		palette = tui.DefaultPalette()
	}

	confirm := opts.Confirm
	if confirm == nil {
		confirm = defaultConfirmation(cmd.InOrStdin(), cmd.OutOrStdout(), palette)
	}

	flags, err := collectFlags(cmd)
	if err != nil {
		return err
	}

	flags.DefaultBackupDirectory = determineBackupDirectory(cfg)

	target, err := resolveTarget(args[0])
	if err != nil {
		return err
	}

	options, err := services.BuildDeleteOptions(target, flags)
	if err != nil {
		return err
	}

	if !options.DryRun && !options.Force {
		req := ConfirmationRequest{
			Target: options.Target,
			Scope:  options.Scope,
			Force:  options.Force,
			DryRun: options.DryRun,
		}

		if options.Backup != nil {
			req.BackupDestination = options.Backup.Destination
		}

		proceed, err := confirm(ctx, req)
		if err != nil {
			return err
		}

		if !proceed {
			fmt.Fprintln(cmd.OutOrStdout(), "operation cancelled")
			return nil
		}
	}

	result, execErr := service.DeleteProject(ctx, options)
	if result != nil {
		renderResult(cmd, result, options.DryRun)
	}

	return execErr
}

func collectFlags(cmd *cobra.Command) (services.DeleteCLIFlags, error) {
	var flags services.DeleteCLIFlags

	all, err := cmd.Flags().GetBool(flagAll)
	if err != nil {
		return flags, err
	}

	keepFiles, err := cmd.Flags().GetBool(flagKeepFiles)
	if err != nil {
		return flags, err
	}

	onlyEnv, err := cmd.Flags().GetBool(flagOnlyEnv)
	if err != nil {
		return flags, err
	}

	dryRun, err := cmd.Flags().GetBool(flagDryRun)
	if err != nil {
		return flags, err
	}

	force, err := cmd.Flags().GetBool(flagForce)
	if err != nil {
		return flags, err
	}

	backup, err := cmd.Flags().GetBool(flagBackup)
	if err != nil {
		return flags, err
	}

	destination, err := cmd.Flags().GetString(flagBackupDestination)
	if err != nil {
		return flags, err
	}

	if strings.TrimSpace(destination) != "" && !backup {
		return flags, fmt.Errorf("--%s requires --%s", flagBackupDestination, flagBackup)
	}

	flags = services.DeleteCLIFlags{
		All:               all,
		KeepFiles:         keepFiles,
		OnlyEnv:           onlyEnv,
		DryRun:            dryRun,
		Force:             force,
		Backup:            backup,
		BackupDestination: destination,
	}

	return flags, nil
}

func resolveTarget(input string) (domain.ProjectIdentifier, error) {
	trimmed := strings.TrimSpace(input)
	if trimmed == "" {
		return domain.ProjectIdentifier{}, fmt.Errorf("target is required")
	}

	if looksLikePath(trimmed) {
		return domain.ProjectIdentifier{Path: trimmed}, nil
	}

	return domain.ProjectIdentifier{Name: trimmed}, nil
}

func looksLikePath(value string) bool {
	if filepath.IsAbs(value) {
		return true
	}

	if strings.HasPrefix(value, "./") || strings.HasPrefix(value, "../") || strings.HasPrefix(value, "~/") {
		return true
	}

	return strings.Contains(value, "/") || strings.Contains(value, "\\")
}

func renderResult(cmd *cobra.Command, result *domain.ProjectDeleteResult, dryRun bool) {
	out := cmd.OutOrStdout()
	errOut := cmd.ErrOrStderr()

	scope := result.Scope.String()
	if scope == "" && result.Plan != nil {
		scope = result.Plan.Scope.String()
	}

	if scope != "" {
		fmt.Fprintf(out, "Scope: %s\n", scope)
	}

	if result.Summary != "" {
		fmt.Fprintln(out, result.Summary)
	}

	if dryRun && result.Plan != nil && len(result.Plan.Artifacts) > 0 {
		fmt.Fprintln(out, "Planned artifacts:")

		for _, artifact := range result.Plan.Artifacts {
			fmt.Fprintf(out, "  - %s\n", describeArtifact(artifact))
		}
	}

	if len(result.ArtifactsRemoved) > 0 {
		fmt.Fprintln(out, "Removed:")

		for _, artifact := range result.ArtifactsRemoved {
			fmt.Fprintf(out, "  - %s\n", describeArtifact(artifact))
		}
	}

	if len(result.ArtifactsSkipped) > 0 {
		fmt.Fprintln(out, "Skipped:")

		for _, artifact := range result.ArtifactsSkipped {
			fmt.Fprintf(out, "  - %s\n", describeArtifact(artifact))
		}
	}

	if result.BackupPath != "" {
		fmt.Fprintf(out, "Backup: %s\n", result.BackupPath)
	}

	if dryRun {
		fmt.Fprintln(out, "No changes were applied.")
	}

	if len(result.Errors) > 0 {
		fmt.Fprintln(errOut, "Errors:")

		for _, e := range result.Errors {
			if e == nil {
				continue
			}

			fmt.Fprintf(errOut, "  - %v\n", e)
		}
	}
}

func determineBackupDirectory(cfg *configs.Config) string {
	if cfg != nil {
		if dir := strings.TrimSpace(cfg.BackupDirectory); dir != "" {
			return dir
		}
	}

	return filepath.Join("~", ".pm", "backups")
}

func describeArtifact(artifact domain.DeletionArtifact) string {
	if desc := strings.TrimSpace(artifact.Description); desc != "" {
		return desc
	}

	if path := strings.TrimSpace(artifact.Path); path != "" {
		return path
	}

	if artifact.Type != "" {
		return string(artifact.Type)
	}

	return "artifact"
}

func resolvePalette(cmd *cobra.Command) (map[string]string, *configs.Config, error) {
	root := cmd.Root()

	themeFlag, _ := root.PersistentFlags().GetString("theme")
	configFlag, _ := root.PersistentFlags().GetString("config")

	cfgPath := strings.TrimSpace(configFlag)
	if cfgPath == "" {
		if v, ok := os.LookupEnv("PM_CONFIG"); ok && strings.TrimSpace(v) != "" {
			cfgPath = v
		}
	}

	cfg, cfgErr := configs.GetConfig(cfgPath)
	if cfgPath != "" && cfgErr != nil {
		return nil, nil, cfgErr
	}

	theme := strings.TrimSpace(themeFlag)
	if theme == "" {
		if v, ok := os.LookupEnv("PM_THEME"); ok && strings.TrimSpace(v) != "" {
			theme = v
		} else if cfg != nil && strings.TrimSpace(cfg.Theme) != "" {
			theme = cfg.Theme
		}
	}

	overrides := map[string]string{}

	if cfg != nil {
		for _, ct := range cfg.CustomThemes {
			if strings.EqualFold(ct.Name, theme) {
				if ct.Title != "" {
					overrides["title"] = ct.Title
				}

				if ct.Section != "" {
					overrides["section"] = ct.Section
				}

				if ct.Subtext != "" {
					overrides["subtext"] = ct.Subtext
				}

				if ct.Text != "" {
					overrides["text"] = ct.Text
				}

				if ct.Placeholder != "" {
					overrides["placeholder"] = ct.Placeholder
				}

				if ct.Border != "" {
					overrides["border"] = ct.Border
				}

				if ct.Error != "" {
					overrides["error"] = ct.Error
				}

				if ct.ButtonDefFg != "" {
					overrides["buttonDefFg"] = ct.ButtonDefFg
				}

				if ct.ButtonDefBg != "" {
					overrides["buttonDefBg"] = ct.ButtonDefBg
				}

				if ct.ButtonSelFg != "" {
					overrides["buttonSelFg"] = ct.ButtonSelFg
				}

				if ct.ButtonSelBg != "" {
					overrides["buttonSelBg"] = ct.ButtonSelBg
				}

				if ct.SelectedFg != "" {
					overrides["selectedFg"] = ct.SelectedFg
				}

				if ct.SelectedBg != "" {
					overrides["selectedBg"] = ct.SelectedBg
				}

				if ct.Help != "" {
					overrides["help"] = ct.Help
				}
			}
		}
	}

	return tui.PaletteForTheme(theme, overrides), cfg, nil
}
