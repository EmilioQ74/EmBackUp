package gui

import (
	"log/slog"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/EmilioQ74/EmBackUp/internal/scheduler"
	"github.com/EmilioQ74/EmBackUp/internal/util"
)

func newScheduleTab(log *slog.Logger) fyne.CanvasObject {
	cronEntry := widget.NewEntry()
	cronEntry.SetText("0 2 * * *")
	cronEntry.SetPlaceHolder("cron expression e.g. 0 2 * * *")

	statusLabel := widget.NewLabel("Scheduler is Stopped.")
	var sc *scheduler.Scheduler

	startBtn := widget.NewButton("Start Scheduler", func() {})
	stopBtn := widget.NewButton("Stop scheduler", func() {})
	stopBtn.Disable()

	startBtn.OnTapped = func() {
		eng, cfg, err := util.BuildEngine(log)
		if err != nil {
			statusLabel.SetText("Error: " + err.Error())
			return
		}

		sc, err = scheduler.New(log)
		if err != nil {
			statusLabel.SetText("Error: " + err.Error())
			return
		}

		if err := sc.Add(cronEntry.Text, eng, cfg); err != nil {
			statusLabel.SetText("Invalid cron: " + err.Error())
			return
		}

		sc.Start()
		startBtn.Disable()
		stopBtn.Enable()
		statusLabel.SetText("Scheduler running — cron: " + cronEntry.Text)
	}

	stopBtn.OnTapped = func() {
		if sc != nil {
			sc.Stop()
		}
		startBtn.Enable()
		stopBtn.Disable()
		statusLabel.SetText("Scheduler stopped.")
	}

	// cron presets
	presets := widget.NewSelect([]string{
		"Every minute   — * * * * *",
		"Every hour     — 0 * * * *",
		"Daily 2am      — 0 2 * * *",
		"Weekly Sunday  — 0 2 * * 0",
		"Monthly 1st    — 0 2 1 * *",
	}, func(s string) {
		// extract the cron part after —
		parts := strings.Split(s, "— ")
		if len(parts) == 2 {
			cronEntry.SetText(strings.TrimSpace(parts[1]))
		}
	})
	presets.PlaceHolder = "Quick presets..."

	return container.NewVBox(
		widget.NewLabel("Schedule automatic backups using a cron expression."),
		presets,
		cronEntry,
		container.NewHBox(startBtn, stopBtn),
		statusLabel,
	)
}
