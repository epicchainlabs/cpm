package main

import (
	"cpm/generators"
	_ "embed"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/epicchainlabs/epicchain-go/pkg/util"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

//go:embed cpm.yaml.default
var defaultConfig []byte
var cfg = &CPMConfig{
	Defaults: Defaults{
		ContractGenerateSdk: false,
		ContractDownload:    true,
	},
}

type ContractConfig struct {
	Label         string          `yaml:"label"`
	ScriptHash    util.Uint160    `yaml:"script-hash"`
	SourceNetwork *string         `yaml:"source-network,omitempty"`
	GenerateSdk   *bool           `yaml:"generate-sdk,omitempty"`
	Download      *bool           `yaml:"download,omitempty"`
	OnChain       *GenerateConfig `yaml:"on-chain,omitempty"`
	OffChain      *GenerateConfig `yaml:"off-chain,omitempty"`
}

type GenerateConfig struct {
	Languages       []string       `yaml:"languages"`
	SdkDestinations SdkDestination `yaml:"destinations"`
}

type SdkDestination struct {
	Csharp *string `yaml:"csharp,omitempty"`
	Golang *string `yaml:"go,omitempty"`
	Java   *string `yaml:"java,omitempty"`
	Python *string `yaml:"python,omitempty"`
	TS     *string `yaml:"ts,omitempty"`
}

type Defaults struct {
	ContractSourceNetwork string          `yaml:"contract-source-network"`
	ContractDestination   string          `yaml:"contract-destination"`
	ContractGenerateSdk   bool            `yaml:"contract-generate-sdk"`
	ContractDownload      bool            `yaml:"contract-download,omitempty"`
	OnChain               *GenerateConfig `yaml:"on-chain,omitempty"`
	OffChain              *GenerateConfig `yaml:"off-chain,omitempty"`
}

type CPMConfig struct {
	Defaults  Defaults         `yaml:"defaults"`
	Contracts []ContractConfig `yaml:"contracts"`
	Tools     struct {
		EpicChainExpress struct {
			CanGenerateSDK      bool    `yaml:"canGenerateSDK"`
			CanDownloadContract bool    `yaml:"canDownloadContract"`
			ExecutablePath      *string `yaml:"executable-path,omitempty"`
			ConfigPath          string  `yaml:"config-path"`
		} `yaml:"epicchain-express"`
	} `yaml:"tools"`
	Networks []struct {
		Label string   `yaml:"label"`
		Hosts []string `yaml:"hosts"`
	} `yaml:"networks"`
}

func LoadConfig() {
	f, err := os.Open(DEFAULT_CONFIG_FILE)
	if err != nil {
		if os.IsNotExist(err) {
			log.Fatalf("Config file %s not found. Run `cpm init` to create a default config", DEFAULT_CONFIG_FILE)
		} else {
			log.Fatal(err)
		}
	}
	defer f.Close()

	yamlData, _ := io.ReadAll(f)
	if err := yaml.Unmarshal(yamlData, &cfg); err != nil {
		log.Fatal(fmt.Errorf("failed to parse config file: %w", err))
	}

	// ensure all contract configs can be worked with directly
	for i, c := range cfg.Contracts {
		if c.SourceNetwork == nil {
			cfg.Contracts[i].SourceNetwork = &cfg.Defaults.ContractSourceNetwork
		}
		if c.GenerateSdk == nil {
			cfg.Contracts[i].GenerateSdk = &cfg.Defaults.ContractGenerateSdk
		}
		if c.Download == nil {
			cfg.Contracts[i].Download = &cfg.Defaults.ContractDownload
		}
	}
}

func CreateDefaultConfig() {
	if _, err := os.Stat(DEFAULT_CONFIG_FILE); os.IsNotExist(err) {
		err = os.WriteFile(DEFAULT_CONFIG_FILE, defaultConfig, 0644)
		if err != nil {
			log.Fatal(err)
		}
		log.Infof("Written %s\n", DEFAULT_CONFIG_FILE)
	} else {
		log.Fatalf("%s already exists", DEFAULT_CONFIG_FILE)
	}
}

func (c *CPMConfig) addContract(label string, scriptHash util.Uint160) {
	for _, c := range cfg.Contracts {
		if c.ScriptHash.Equals(scriptHash) {
			return
		}
	}
	cfg.Contracts = append(cfg.Contracts, ContractConfig{Label: label, ScriptHash: scriptHash})
}

func (c *CPMConfig) getHosts(networkLabel string) []string {
	for _, network := range c.Networks {
		if network.Label == networkLabel {
			return network.Hosts
		}
	}
	log.Fatalf("Could not find hosts for label: %s", networkLabel)
	return nil
}

func (c *CPMConfig) getSdkDestination(forLanguage string, sdkType string) string {
	defaultLocation := generators.OutputRoot + sdkType + "/" + forLanguage + "/"

	if c == nil {
		return defaultLocation
	}

	var sdkTypePath SdkDestination
	if sdkType == generators.SDKOnChain {
		if c.Defaults.OnChain == nil {
			return defaultLocation
		}
		sdkTypePath = c.Defaults.OnChain.SdkDestinations
	} else {
		if c.Defaults.OffChain == nil {
			return defaultLocation
		}
		sdkTypePath = c.Defaults.OffChain.SdkDestinations
	}

	switch forLanguage {
	case LANG_PYTHON:
		if path := sdkTypePath.Python; path != nil {
			return EnsureSuffix(*path)
		}
		return defaultLocation
	case LANG_GO:
		if path := sdkTypePath.Golang; path != nil {
			return EnsureSuffix(*path)
		}
		return defaultLocation
	case LANG_JAVA:
		if path := sdkTypePath.Java; path != nil {
			return EnsureSuffix(*path)
		}
		return defaultLocation
	case LANG_CSHARP:
		if path := sdkTypePath.Csharp; path != nil {
			return EnsureSuffix(*path)
		}
		return defaultLocation
	case LANG_TYPESCRIPT:
		if path := sdkTypePath.TS; path != nil {
			return EnsureSuffix(*path)
		}
		return defaultLocation
	default:
		return defaultLocation
	}
}

func (c *CPMConfig) saveToDisk() {
	f, err := os.Create(DEFAULT_CONFIG_FILE)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	data, err := yaml.Marshal(cfg)
	if err != nil {
		return
	}

	_, err = f.Write(data)
	if err != nil {
		log.Fatal(err)
	}
}

type EnumValue struct {
	Enum     []string
	Default  string
	selected string
}

func (e *EnumValue) Set(value string) error {
	for _, enum := range e.Enum {
		if enum == value {
			e.selected = value
			return nil
		}
	}

	return fmt.Errorf("allowed values are %s", strings.Join(e.Enum, ", "))
}

func (e EnumValue) String() string {
	if e.selected == "" {
		return e.Default
	}
	return e.selected
}
