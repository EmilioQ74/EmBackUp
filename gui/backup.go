package gui

import (
	"context"
	"log/slog"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/EmilioQ74/EmBackUp/internal/util"
)

func newBackupTab(log *slog.Logger) fyne.CanvasObject {
	logOutput := widget.NewMultiLineEntry()
	logOutput.Disable()
	logOutput.SetPlaceHolder("Backup output will appear here...")
	logOutput.SetMinRowsVisible(10)

	addLog := func(msg string) {
		ts := time.Now().Format("15:04:05")
		logOutput.SetText(logOutput.Text + "[" + ts + "]" + msg + "\n")
	}

	progress := widget.NewProgressBarInfinite()
	progress.Hide()

	backupBtn := widget.NewButton("Run Backup Now", func() {})
	backupBtn.Importance = widget.HighImportance

	backupBtn.OnTapped = func() {
		backupBtn.Disable()
		progress.Show()
		addLog("Starting backup...")

		go func() {
			eng, cfg, err := util.BuildEngine(log)
			if err != nil {
				fyne.Do(func() {
					addLog("ERROR: " + err.Error())
					backupBtn.Enable()
					progress.Hide()
				})
				return
			}

			if err := eng.Backup(context.Background(), cfg); err != nil {
				fyne.Do(func() {
					addLog("ERROR: " + err.Error())
					backupBtn.Enable()
					progress.Hide()
				})
				return
			}

			fyne.Do(func() {
				addLog("Backup completed successfully")
				backupBtn.Enable()
				progress.Hide()
			})
		}()
	}

	clearBtn := widget.NewButton("Clear Log", func() {
		logOutput.SetText("")
	})

	buttons := container.NewHBox(backupBtn, clearBtn)
	content := container.NewBorder(
		container.NewVBox(
			widget.NewLabel("One-shot backup using current settings."),
			buttons,
			progress,
		),
		nil, nil, nil, container.NewScroll(logOutput),
	)
	return content
}
