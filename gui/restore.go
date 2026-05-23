package gui

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/EmilioQ74/EmBackUp/internal/storage"
	"github.com/EmilioQ74/EmBackUp/internal/util"
)

func newRestoreTab(log *slog.Logger, w fyne.Window) fyne.CanvasObject {
	var backups []storage.BackupMeta
	selectedKey := ""

	statusLabel := widget.NewLabel("Select a backup to restore.")

	backupList := widget.NewList(
		func() int { return len(backups) },
		func() fyne.CanvasObject { return widget.NewLabel("") },
		func(i widget.ListItemID, o fyne.CanvasObject) {
			o.(*widget.Label).SetText(
				backups[i].Key + " — " + util.FormatSize(backups[i].Size),
			)
		},
	)
	var refreshBtn *widget.Button
	refreshBtn = widget.NewButton("Refresh List", func() {
		statusLabel.SetText("Loading...")
		refreshBtn.Disable()

		go func() {
			defer fyne.Do(func() { refreshBtn.Enable() })

			eng, cfg, err := util.BuildEngine(log)
			if err != nil {
				fyne.Do(func() {
					statusLabel.SetText("Error: " + err.Error())
				})
				return
			}

			items, err := eng.List(context.Background(), cfg.Database)
			if err != nil {
				fyne.Do(func() {
					statusLabel.SetText("Error: " + err.Error())
				})
				return
			}

			fyne.Do(func() {
				backups = items
				selectedKey = ""
				backupList.UnselectAll()
				backupList.Refresh()
				statusLabel.SetText(fmt.Sprintf("Found %d backup(s).", len(items)))
			})
		}()
	})

	var restoreBtn *widget.Button
	restoreBtn = widget.NewButton("Restore Selected", func() {
		if selectedKey == "" {
			return
		}
		dialog.ShowConfirm(
			"Confirm Restore",
			fmt.Sprintf("Overwrite the database with %q?\nThis cannot be undone.", selectedKey),
			func(confirmed bool) {
				if !confirmed {
					return
				}
				restoreBtn.Disable()
				statusLabel.SetText("Restoring…")

				go func() {
					defer fyne.Do(func() { restoreBtn.Enable() })

					ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
					defer cancel()

					eng, cfg, err := util.BuildEngine(log)
					if err != nil {
						fyne.Do(func() {
							statusLabel.SetText("Error: " + err.Error())
						})
						return
					}
					if err := eng.Restore(ctx, cfg, selectedKey); err != nil {
						fyne.Do(func() {
							statusLabel.SetText("Restore failed: " + err.Error())
						})
						return
					}
					fyne.Do(func() {
						statusLabel.SetText("Restore completed successfully.")
						dialog.ShowInformation(
							"Restore Complete",
							"The database was restored successfully.",
							w,
						)
					})
				}()
			},
			w,
		)
	})
	restoreBtn.Importance = widget.DangerImportance

	backupList.OnSelected = func(i widget.ListItemID) {
		if i < len(backups) {
			selectedKey = backups[i].Key
			statusLabel.SetText("Selected: " + selectedKey)
		}
	}

	buttons := container.NewHBox(refreshBtn, restoreBtn)

	return container.NewBorder(
		container.NewVBox(
			widget.NewLabel("Select a backup and click Restore."),
			buttons,
			statusLabel,
		),
		nil, nil, nil,
		backupList,
	)
}
