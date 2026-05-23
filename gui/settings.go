package gui

import (
	"fmt"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/spf13/viper"
)

func newSettingsTab() fyne.CanvasObject {
	// Database
	dbType := widget.NewSelect(
		[]string{"postgres", "mysql", "sqlite", "mongodb"}, nil,
	)
	dbType.SetSelected(viper.GetString("db.type"))

	host := widget.NewEntry()
	host.SetText(viper.GetString("db.host"))

	port := widget.NewEntry()
	port.SetText(fmt.Sprint(viper.GetInt("db.port")))

	user := widget.NewEntry()
	user.SetText(viper.GetString("db.user"))

	password := widget.NewPasswordEntry()
	password.SetText(viper.GetString("db.password"))

	dbName := widget.NewEntry()
	dbName.SetText(viper.GetString("db.name"))

	// Storage
	storageType := widget.NewSelect([]string{"local", "s3"}, nil)
	storageType.SetSelected(viper.GetString("storage.type"))

	storagePath := widget.NewEntry()
	storagePath.SetText(viper.GetString("storage.path"))

	// S3
	bucket := widget.NewEntry()
	bucket.SetText(viper.GetString("storage.bucket"))

	region := widget.NewEntry()
	region.SetText(viper.GetString("storage.region"))

	accessKey := widget.NewEntry()
	accessKey.SetText(viper.GetString("storage.access_key"))
	accessKey.SetPlaceHolder("Leave blank to use env / IAM role")

	secretKey := widget.NewPasswordEntry()
	secretKey.SetText(viper.GetString("storage.secret_key"))
	secretKey.SetPlaceHolder("Leave blank to use env / IAM role")

	// Compression
	compressType := widget.NewSelect([]string{"gzip", "zstd", "none"}, nil)
	compressType.SetSelected(viper.GetString("compress"))

	// Retention
	retention := widget.NewEntry()
	retention.SetText(fmt.Sprint(viper.GetInt("retention")))
	retention.SetPlaceHolder("0 = disabled")

	saveBtn := widget.NewButton("Save Settings", func() {
		p, _ := strconv.Atoi(port.Text)
		r, _ := strconv.Atoi(retention.Text)

		viper.Set("db.type", dbType.Selected)
		viper.Set("db.host", host.Text)
		viper.Set("db.port", p)
		viper.Set("db.user", user.Text)
		viper.Set("db.password", password.Text)
		viper.Set("db.name", dbName.Text)

		viper.Set("storage.type", storageType.Selected)
		viper.Set("storage.path", storagePath.Text)
		viper.Set("storage.bucket", bucket.Text)
		viper.Set("storage.region", region.Text)
		viper.Set("storage.access_key", accessKey.Text)
		viper.Set("storage.secret_key", secretKey.Text)

		viper.Set("compress", compressType.Selected)
		viper.Set("retention", r)

		if err := viper.WriteConfig(); err != nil {
			_ = viper.SafeWriteConfig()
		}
	})
	saveBtn.Importance = widget.HighImportance

	form := container.NewVBox(
		widget.NewLabel("Database"),
		widget.NewForm(
			widget.NewFormItem("Type", dbType),
			widget.NewFormItem("Host", host),
			widget.NewFormItem("Port", port),
			widget.NewFormItem("User", user),
			widget.NewFormItem("Password", password),
			widget.NewFormItem("Database", dbName),
		),

		widget.NewSeparator(),

		widget.NewLabel("Storage"),
		widget.NewForm(
			widget.NewFormItem("Type", storageType),
			widget.NewFormItem("Path / prefix", storagePath),
			widget.NewFormItem("S3 bucket", bucket),
			widget.NewFormItem("S3 region", region),
			widget.NewFormItem("Access key", accessKey),
			widget.NewFormItem("Secret key", secretKey),
		),

		widget.NewSeparator(),

		widget.NewLabel("Compression"),
		widget.NewForm(
			widget.NewFormItem("Algorithm", compressType),
		),

		widget.NewSeparator(),

		widget.NewLabel("Retention"),
		widget.NewForm(
			widget.NewFormItem("Keep last N backups", retention),
		),

		widget.NewSeparator(),

		saveBtn,
	)

	return container.NewScroll(form)
}
