package util

import (
	"fmt"
	"log/slog"

	"github.com/EmilioQ74/EmBackUp/internal/adapters"
	"github.com/EmilioQ74/EmBackUp/internal/compress"
	"github.com/EmilioQ74/EmBackUp/internal/engine"
	"github.com/EmilioQ74/EmBackUp/internal/storage"
	"github.com/spf13/viper"
)

func BuildEngine(log *slog.Logger) (*engine.Engine, adapters.DBConfig, error) {
	dbType := viper.GetString("db.type")
	var adapter adapters.Adapter
	switch dbType {
	case "postgres":
		adapter = &adapters.PostgresAdapter{}
	case "mysql":
		adapter = &adapters.MySQLAdapter{}
	case "sqlite":
		adapter = &adapters.SQLiteAdapter{}
	case "mongodb":
		adapter = &adapters.MongoAdapter{}
	default:
		return nil, adapters.DBConfig{}, fmt.Errorf("unknown db.type %q", dbType)
	}

	var store storage.Storage
	var err error
	storageType := viper.GetString("storage.type")
	switch storageType {
	case "local":
		store, err = storage.NewLocal(viper.GetString("storage.path"))
	case "s3":
		accessKey := viper.GetString("storage.access_key")
		secretKey := viper.GetString("storage.secret_key")
		if accessKey != "" && secretKey != "" {
			store, err = storage.NewS3WithCredentials(
				viper.GetString("storage.bucket"),
				viper.GetString("storage.region"),
				accessKey,
				secretKey,
			)
		} else {
			store, err = storage.NewS3(
				viper.GetString("storage.bucket"),
				viper.GetString("storage.region"),
			)
		}
	default:
		return nil, adapters.DBConfig{}, fmt.Errorf("unknown storage.type %q", storageType)
	}
	if err != nil {
		return nil, adapters.DBConfig{}, fmt.Errorf("storage init: %w", err)
	}

	retention := viper.GetInt("retention")
	algo := compress.Algorithm(viper.GetString("compress"))

	cfg := adapters.DBConfig{
		Host:     viper.GetString("db.host"),
		Port:     viper.GetInt("db.port"),
		User:     viper.GetString("db.user"),
		Password: viper.GetString("db.password"),
		Database: viper.GetString("db.name"),
	}

	eng := engine.New(adapter, store, algo, retention, log)
	return eng, cfg, nil
}

func FormatSize(b int64) string {
	switch {
	case b >= 1<<30:
		return fmt.Sprintf("%.1f GB", float64(b)/(1<<30))
	case b >= 1<<20:
		return fmt.Sprintf("%.1f MB", float64(b)/(1<<20))
	case b >= 1<<10:
		return fmt.Sprintf("%.1f KB", float64(b)/(1<<10))
	default:
		return fmt.Sprintf("%d B", b)
	}
}
