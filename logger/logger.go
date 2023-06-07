package logger

import (
	"fmt"
	"go.uber.org/zap"
)

func Initialize(level string) (sugar *zap.SugaredLogger, error error) {
	lvl, err := zap.ParseAtomicLevel(level)
	if err != nil {
		return nil, fmt.Errorf("failed to parse the logging level %s: %w", level, err)
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
