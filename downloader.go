package main

import (
	"bytes"
	"fmt"
	"os/exec"
	"runtime"
	"strings"

	"github.com/epicchainlabs/epicchain-go/pkg/util"
	log "github.com/sirupsen/logrus"
)

type Downloader interface {
	downloadContract(scriptHash util.Uint160, host string) (string, error)
}

type EpicChainExpressDownloader struct {
	expressConfigPath *string
}

func NewEpicChainExpressDownloader(configPath string) Downloader {
	executablePath := cfg.Tools.EpicChainExpress.ExecutablePath
	if executablePath == nil {
		var cmd *exec.Cmd
		if runtime.GOOS == "darwin" {
			cmd = exec.Command("bash", "-c", "epicchainxp -h")
		} else {
			cmd = exec.Command("epicchainxp", "-h")
		}
		err := cmd.Run()
		if err != nil {
			log.Fatal("Could not find 'epicchainxp' executable in $PATH. Please install epicchainxp globally using " +
				"'dotnet tool install EpicChain.Express -g'" +
				" or specify the 'executable-path' in cpm.yaml in the epicchain-express tools section")
		}
	} else {
		// Verify path works by calling help (which has a 0 exit code)
		cmd := exec.Command(*executablePath, "-h")
		err := cmd.Run()
		if err != nil {
			log.Fatal(fmt.Errorf("could not find 'epicchainxp' executable in the configured executable-path: %w", err))
		}
	}
	return &EpicChainExpressDownloader{
		expressConfigPath: &configPath,
	}
}

func (ned *EpicChainExpressDownloader) downloadContract(scriptHash util.Uint160, host string) (string, error) {
	// the name and arguments supplied to exec.Command differ slightly depending on the OS and whether epicchainxp is
	// installed globally. the following are the base arguments that hold for all scenarios
	args := []string{"contract", "download", "-i", cfg.Tools.EpicChainExpress.ConfigPath, "--force", "0x" + scriptHash.StringLE(), host}

	// global default
	executable := "epicchainxp"

	if cfg.Tools.EpicChainExpress.ExecutablePath != nil {
		executable = *cfg.Tools.EpicChainExpress.ExecutablePath
	} else if runtime.GOOS == "darwin" {
		executable = "bash"
		tmp := append([]string{"epicchainxp"}, args...)
		args = []string{"-c", strings.Join(tmp, " ")}
	}

	cmd := exec.Command(executable, args...)
	var errOut bytes.Buffer
	cmd.Stderr = &errOut
	out, err := cmd.Output()
	if err != nil {
		return "[epicchainxp]" + errOut.String(), err
	} else {
		return "[epicchainxp]" + string(out), nil
	}
}
