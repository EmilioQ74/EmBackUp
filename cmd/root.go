package cmd

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/EmilioQ74/EmBackUp/internal/scheduler"
	"github.com/EmilioQ74/EmBackUp/internal/util"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var rootCmd = &cobra.Command{
	Use:   "EmbackUps",
	Short: "Multi-DBMS Backup utility",
}

func ExecuteContext(ctx context.Context) {
	log := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	// Load config before any command runs
	viper.SetConfigName(".emBackup")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("$HOME")
	viper.AddConfigPath(".")
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	_ = viper.ReadInConfig()

	rootCmd.AddCommand(
		backupCmd(log),
		scheduleCmd(log),
		listCmd(log),
		restoreCmd(log),
	)

	if err := rootCmd.ExecuteContext(ctx); err != nil {
		os.Exit(1)
	}
}

func backupCmd(log *slog.Logger) *cobra.Command {
	return &cobra.Command{
		Use:   "backup",
		Short: "Run a one-shot backup now",
		RunE: func(cmd *cobra.Command, _ []string) error {
			eng, cfg, err := util.BuildEngine(log)
			if err != nil {
				return err
			}
			return eng.Backup(cmd.Context(), cfg)
		},
	}
}

func restoreCmd(log *slog.Logger) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "restore [backup-key]",
		Short: "Restore a backup by its storage key",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			eng, cfg, err := util.BuildEngine(log)
			if err != nil {
				return err
			}
			key := args[0]
			return eng.Restore(cmd.Context(), cfg, key)
		},
	}
	return cmd
}

func listCmd(log *slog.Logger) *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List all available backups",
		RunE: func(cmd *cobra.Command, _ []string) error {
			eng, cfg, err := util.BuildEngine(log)
			if err != nil {
				return err
			}
			items, err := eng.List(cmd.Context(), cfg.Database)
			if err != nil {
				return err
			}
			if len(items) == 0 {
				fmt.Println("no backups found")
				return nil
			}
			fmt.Printf("%-50s  %-10s  %s\n", "KEY", "SIZE", "CREATED")
			for _, m := range items {
				fmt.Printf("%-50s  %-10s  %s\n",
					m.Key,
					util.FormatSize(m.Size),
					m.CreatedAt.Format("2006-01-02 15:04:05 UTC"),
				)
			}
			return nil
		},
	}
}

func scheduleCmd(log *slog.Logger) *cobra.Command {
	var cron string
	cmd := &cobra.Command{
		Use:   "schedule",
		Short: "Run backups on a cron schedule (blocks until stopped)",
		RunE: func(cmd *cobra.Command, _ []string) error {
			eng, cfg, err := util.BuildEngine(log)
			if err != nil {
				return err
			}

			sc, err := scheduler.New(log)
			if err != nil {
				return fmt.Errorf("scheduler init: %w", err)
			}

			if err := sc.Add(cron, eng, cfg); err != nil {
				return fmt.Errorf("invalid cron expression: %w", err)
			}

			log.Info("scheduler started", "cron", cron, "db", cfg.Database)
			sc.Start()

			<-cmd.Context().Done()
			sc.Stop()
			log.Info("scheduler stopped")
			return nil
		},
	}

	cmd.Flags().StringVar(&cron, "cron", "0 2 * * *", `cron expression, e.g. "0 2 * * *" = daily at 02:00`)
	return cmd
}
