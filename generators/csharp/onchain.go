package csharp

import (
	"cpm/generators"
	"fmt"
	"os"
	"text/template"

	"github.com/epicchainlabs/epicchain-go/pkg/smartcontract"
	"github.com/iancoleman/strcase"
	log "github.com/sirupsen/logrus"
)

const csharpSrcTmpl = `
{{- define "METHOD" }}
        public static {{.ReturnType}} {{.Name}}({{range $index, $arg := .Arguments -}}
       {{- if ne $index 0}}, {{end}}
          {{- .Type}} {{.Name}}
       {{- end}}) {
            return ({{.ReturnType}}) Contract.Call(ScriptHash, "{{.NameABI}}", CallFlags.All
{{- range $arg := .Arguments -}}, {{ .Name -}} {{ else }}, new object[0]{{end}});
		} 
{{- end -}}
using EpicChain;
using EpicChain.Cryptography.ECC;
using EpicChain.SmartContract;
using EpicChain.SmartContract.Framework;
using EpicChain.SmartContract.Framework.Services;
using EpicChain.SmartContract.Framework.Attributes;

namespace cpm {
    public class {{ .ContractName }}  {

        [InitialValue("{{.Hash}}", ContractParameterType.Hash160)]
        static readonly UInt160 ScriptHash;

        {{- range $m := .Methods}}
        {{ template "METHOD" $m -}}
        {{end}}
    }
}
`

func GenerateCsharpSDK(cfg *generators.GenerateCfg) error {
	err := createCsharpPackage(cfg)
	defer cfg.ContractOutput.Close()
	if err != nil {
		return err
	}

	cfg.MethodNameConverter = strcase.ToCamel
	cfg.ParamTypeConverter = scTypeToCsharp
	ctr, err := generators.TemplateFromManifest(cfg)
	if err != nil {
		return fmt.Errorf("failed to parse manifest into contract template: %v", err)
	}

	tmp, err := template.New("generate").Parse(csharpSrcTmpl)
	if err != nil {
		return fmt.Errorf("failed to parse C# source template: %v", err)
	}

	err = tmp.Execute(cfg.ContractOutput, ctr)
	if err != nil {
		return fmt.Errorf("failed to generate C# code using template: %v", err)
	}

	wd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get working directory: %v", err)
	}

	sdkLocation := wd + "/" + cfg.SdkDestination + generators.UpperFirst(cfg.Manifest.Name) + ".cs"
	log.Infof("Created SDK for contract '%s' at %s with contract hash 0x%s", cfg.Manifest.Name, sdkLocation, cfg.ContractHash.StringLE())

	return nil
}

func createCsharpPackage(cfg *generators.GenerateCfg) error {
	dir := cfg.SdkDestination
	err := os.MkdirAll(dir, 0755)
	if err != nil {
		return fmt.Errorf("can't create directory %s: %w", dir, err)
	}

	filename := generators.UpperFirst(cfg.Manifest.Name)
	f, err := os.Create(fmt.Sprintf(dir+"%s.cs", filename))
	if err != nil {
		f.Close()
		return fmt.Errorf("can't create %s.cs file: %w", filename, err)
	} else {
		cfg.ContractOutput = f
	}

	return nil
}

func scTypeToCsharp(typ smartcontract.ParamType) string {
	switch typ {
	case smartcontract.AnyType:
		return "object"
	case smartcontract.BoolType:
		return "bool"
	case smartcontract.InteropInterfaceType:
		return "object"
	case smartcontract.IntegerType:
		return "BigInteger"
	case smartcontract.ByteArrayType:
		return "byte[]"
	case smartcontract.StringType:
		return "string"
	case smartcontract.Hash160Type:
		return "UInt160"
	case smartcontract.Hash256Type:
		return "UInt256"
	case smartcontract.PublicKeyType:
		return "ECPoint"
	case smartcontract.ArrayType:
		return "object[]"
	case smartcontract.MapType:
		return "Map<object, object>"
	case smartcontract.VoidType:
		return "void"
	default:
		panic(fmt.Sprintf("unknown type: %T %s", typ, typ))
	}
}
