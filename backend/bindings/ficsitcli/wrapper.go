package ficsitcli

import (
	"context"
	"log/slog"
	"time"

	"github.com/mitchellh/go-ps"
	"github.com/pkg/errors"
	"github.com/satisfactorymodding/ficsit-cli/cli"
	"github.com/satisfactorymodding/ficsit-cli/cli/provider"
	wailsRuntime "github.com/wailsapp/wails/v2/pkg/runtime"

	"github.com/satisfactorymodding/SatisfactoryModManager/backend/settings"
)

type FicsitCLI struct {
	ctx                  context.Context
	ficsitCli            *cli.GlobalContext
	installations        []*InstallationInfo
	installFindErrors    []error
	selectedInstallation *InstallationInfo
	progress             *Progress
	isGameRunning        bool
}

func MakeFicsitCLI() (*FicsitCLI, error) {
	f := &FicsitCLI{}
	err := f.Init()
	if err != nil {
		return nil, err
	}
	return f, nil
}

func (f *FicsitCLI) Startup(ctx context.Context) {
	f.ctx = ctx
}

func (f *FicsitCLI) Init() error {
	if f.ficsitCli != nil {
		return errors.New("FicsitCLIWrapper already initialized")
	}
	var err error
	f.ficsitCli, err = cli.InitCLI(false)
	if err != nil {
		return errors.Wrap(err, "failed to initialize ficsit-cli")
	}
	f.ficsitCli.Provider.(*provider.MixedProvider).Offline = settings.Settings.Offline
	err = f.initInstallations()
	if err != nil {
		return errors.Wrap(err, "failed to initialize installations")
	}
	gameRunningTicker := time.NewTicker(5 * time.Second)
	go func() {
		for range gameRunningTicker.C {
			processes, err := ps.Processes()
			if err != nil {
				slog.Error("failed to get processes", slog.Any("error", err))
				continue
			}
			f.isGameRunning = false
			for _, process := range processes {
				if process.Executable() == "FactoryGame-Win64-Shipping.exe" || process.Executable() == "FactoryGame-Win64-Shipping" {
					f.isGameRunning = true
					break
				}
			}
			wailsRuntime.EventsEmit(f.ctx, "isGameRunning", f.isGameRunning)
		}
	}()
	return nil
}

func (f *FicsitCLI) setProgress(p *Progress) {
	f.progress = p
	wailsRuntime.EventsEmit(f.ctx, "progress", p)
}
