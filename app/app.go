package app

import (
	"fmt"
	"time"

	"github.com/TrevorEdris/banner"
	"github.com/TrevorEdris/go-csv/app/config"
	"github.com/TrevorEdris/go-csv/app/csv"
	"github.com/TrevorEdris/go-csv/app/log"
	"go.uber.org/zap"
)

type App struct {
	cfg    *config.Config
	logger *zap.Logger
}

func New(runtimeCfg *config.Runtime) (*App, error) {
	a := App{}
	a.logger = log.New(runtimeCfg.LogLevel)

	var err error
	a.cfg, err = config.New(runtimeCfg)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize config: %w", err)
	}

	return &a, nil
}

func (a *App) Generate() error {
	a.logger.Info("Generating CSV")
	start := time.Now()

	f := csv.New(a.cfg.Columns, a.cfg.Metadata.RowCount, a.cfg.Metadata.DelimiterRune)
	for i := 0; i < a.cfg.Metadata.RowCount; i++ {
		f.AddRow()
		a.logger.Debug(banner.New(fmt.Sprintf("%d", i), banner.WithChar('-'), banner.WithLength(40), banner.Yellow()))
		if (i+1)%100000 == 0 {
			a.logger.Info(banner.New(fmt.Sprintf("%d", i), banner.WithChar('-'), banner.WithLength(40), banner.Yellow()))
		}
		a.logger.Debug(fmt.Sprintf("%+v", f.Rows[f.RecordCount-1]))
	}

	err := f.Write(a.cfg.Runtime.Output)
	if err != nil {
		return fmt.Errorf("failed to write to file %s: %w", a.cfg.Runtime.Output, err)
	}

	elapsed := time.Since(start)
	a.logger.Info(fmt.Sprintf("Completed in %v", elapsed))
	return nil
}
