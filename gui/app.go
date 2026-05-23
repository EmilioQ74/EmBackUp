package gui

import (
	"log/slog"
	"os"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
)

func Start() {
	log := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	a := app.New()
	a.Settings().SetTheme(theme.DarkTheme())

	w := a.NewWindow("EmBackUp - Database Backup Manager")
	w.Resize(fyne.NewSize(800, 600))
	w.SetFixedSize(false)

	tabs := container.NewAppTabs(
		container.NewTabItemWithIcon("Backup", theme.StorageIcon(), newBackupTab(log)),
		container.NewTabItemWithIcon("Restore", theme.HistoryIcon(), newRestoreTab(log, w)),
		container.NewTabItemWithIcon("Schedule", theme.MediaPlayIcon(), newScheduleTab(log)),
		container.NewTabItemWithIcon("Settings", theme.SettingsIcon(), newSettingsTab()),
	)
	tabs.SetTabLocation(container.TabLocationLeading) // tabs on the left side

	w.SetContent(tabs)
	w.ShowAndRun()
}
