package python

import (
	"cpm/generators"
	"fmt"
	"os"
	"strings"
	"text/template"

	"github.com/iancoleman/strcase"
	"github.com/epicchainlabs/epicchain-go/pkg/smartcontract"
	log "github.com/sirupsen/logrus"
)

/*
	Creates a Python SDK that can be easily used when writing smart contracts with neo3-boa.
	The output is a Python package. For example for a contract named `samplecontract` it results in the folder structure
		.
		├── samplecontract
		│   ├── __init__.py
		│   └── contract.py

	which can be used in your neo3-boa contract with

		from samplecontract import Samplecontract
		Samplecontract.func1()
*/

const pythonSrcTmpl = `
{{- define "METHOD" }}
    @staticmethod
    {{- if ne .NameABI .Name}}
    @display_name("{{.NameABI }}"){{end}}
    def {{.Name}}({{range $index, $arg := .Arguments -}}
       {{- if ne $index 0}}, {{end}}
          {{- .Name}}: {{.Type}}
       {{- end}}) -> {{if .ReturnType }}{{ .ReturnType }}: {{ else }} None: {{ end }}
        pass
{{- end -}}
from boa3.builtin.type import UInt160, UInt256, ECPoint
from boa3.builtin.compile_time import contract, display_name
from typing import cast, Any


@contract('{{ .Hash }}')
class {{ .ContractName }}:
    hash: UInt160
{{- range $m := .Methods}}
{{ template "METHOD" $m -}}
{{end}}`

func generateOnchainSDK(cfg *generators.GenerateCfg) error {
	err := createPythonPackage(cfg)
	defer cfg.ContractOutput.Close()
	if err != nil {
		return err
	}

	cfg.MethodNameConverter = strcase.ToSnake
	cfg.ParamTypeConverter = scTypeToPython
	ctr, err := generators.TemplateFromManifest(cfg)
	if err != nil {
		return fmt.Errorf("failed to parse manifest into contract template: %v", err)
	}

	tmp, err := template.New("generate").Parse(pythonSrcTmpl)
	if err != nil {
		return fmt.Errorf("failed to parse Python source template: %v", err)
	}

	err = tmp.Execute(cfg.ContractOutput, ctr)
	if err != nil {
		return fmt.Errorf("failed to generate Python code using template: %v", err)
	}

	wd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get working directory: %v", err)
	}
	sdkLocation := wd + "/" + cfg.SdkDestination + generators.UpperFirst(cfg.Manifest.Name)
	log.Infof("Created SDK for contract '%s' at %s with contract hash 0x%s", cfg.Manifest.Name, sdkLocation, cfg.ContractHash.StringLE())

	return nil
}

// create the Python package structure and set the ContractOutput to the open file handle
func createPythonPackage(cfg *generators.GenerateCfg) error {
	sdkDir := cfg.SdkDestination + strings.ToLower(cfg.Manifest.Name)
	err := os.MkdirAll(sdkDir, 0755)
	if err != nil {
		return fmt.Errorf("can't create on-chain directory %s: %w", sdkDir, err)
	}

	f, err := os.Create(sdkDir + "/__init__.py")
	if err != nil {
		f.Close()
		return fmt.Errorf("can't create on-chain __init__.py file: %w", err)
	}
	f.Close()

	f, err = os.Create(sdkDir + "/contract.py")
	if err != nil {
		f.Close()
		return fmt.Errorf("can't create on-chain contract.py file: %w", err)
	} else {
		cfg.ContractOutput = f
	}
	return nil
}

func scTypeToPython(typ smartcontract.ParamType) string {
	switch typ {
	case smartcontract.AnyType, smartcontract.InteropInterfaceType:
		return "Any"
	case smartcontract.BoolType:
		return "bool"
	case smartcontract.IntegerType:
		return "int"
	case smartcontract.ByteArrayType:
		return "bytes"
	case smartcontract.StringType:
		return "str"
	case smartcontract.Hash160Type:
		return "UInt160"
	case smartcontract.Hash256Type:
		return "UInt256"
	case smartcontract.PublicKeyType:
		return "ECPoint"
	case smartcontract.ArrayType:
		return "list"
	case smartcontract.MapType:
		return "dict"
	case smartcontract.VoidType:
		return "None"
	default:
		panic(fmt.Sprintf("unknown type: %T %s", typ, typ))
	}
}
