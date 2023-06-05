package logger

import (
	"go.uber.org/zap"
)

func Initialize(level string) (sugar *zap.SugaredLogger, error error) {
	lvl, err := zap.ParseAtomicLevel(level)
	if err != nil {
		return nil, err
	}
	cfg := zap.NewDevelopmentConfig()
	cfg.Level = lvl
	zl, err := cfg.Build()
	if err != nil {
		return nil, err
	}
	sugar = zl.Sugar()
	return sugar, nil
}
