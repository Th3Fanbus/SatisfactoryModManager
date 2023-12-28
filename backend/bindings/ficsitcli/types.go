package ficsitcli

import (
	"github.com/satisfactorymodding/ficsit-cli/cli"

	"github.com/satisfactorymodding/SatisfactoryModManager/backend/installfinders/common"
)

type InstallationInfo struct {
	Installation *cli.Installation    `json:"installation"`
	Info         *common.Installation `json:"info"`
}

type Progress struct {
	Item     string  `json:"item"`
	Message  string  `json:"message"`
	Progress float64 `json:"progress"`
}
