package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/Kunde21/markdownfmt/v2/markdownfmt"
	"github.com/gunk/docgen/generate"
	"github.com/gunk/docgen/parser"
	"github.com/gunk/plugin"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/pluginpb"
)

const modeAppend = "append"

func main() {
	plugin.RunMain(new(docPlugin))
}

type docPlugin struct{}

func (p *docPlugin) Generate(req *pluginpb.CodeGeneratorRequest) (*pluginpb.CodeGeneratorResponse, error) {
	var (
		lang            []string
		mode, dir       string
		customHeaderIds generate.CustomHeaderIdsOpt
		onlyExternal    parser.OnlyExternalOpt
	)
	if param := req.GetParameter(); param != "" {
		ps := strings.Split(param, ",")
		for _, p := range ps {
			s := strings.Split(p, "=")
			if len(s) != 2 {
				return nil, fmt.Errorf("could not parse parameter: %s", p)
			}
			k, v := s[0], s[1]
			switch k {
			case "lang":
				lang = append(lang, v)
			case "mode":
				if v != modeAppend {
					return nil, fmt.Errorf("unknown mode: %s", v)
				}
				mode = v
			case "out":
				dir = v
			case "custom-header-ids":
				boolVal, err := strconv.ParseBool(v)
				if err != nil {
					return nil, err
				}
				if boolVal {
					customHeaderIds = generate.CustomHeaderIdsOptEnabled
				}
			case "only_external":
				boolVal, err := strconv.ParseBool(v)
				if err != nil {
					return nil, err
				}
				if boolVal {
					onlyExternal = parser.OnlyExternalEnabled
				}
			default:
				return nil, fmt.Errorf("unknown parameter: %s", k)
			}
		}
	}
	// find the source by looping through the protofiles and finding the one
	// that matches the FileToGenerate
	var source *parser.FileDescWrapper
	for _, f := range req.GetProtoFile() {
		if strings.Contains(f.GetName(), "all.proto") {
			for _, fileToGenerate := range req.FileToGenerate {
				if fileToGenerate == f.GetName() {
					source = &parser.FileDescWrapper{FileDescriptorProto: f}
					break
				}
			}
		}
	}
	if source == nil {
		return nil, fmt.Errorf("no file to generate")
	}
	base := filepath.Join(filepath.Dir(source.GetName()))
	source.DependencyMap = parser.GenerateDependencyMap(source.FileDescriptorProto, req.GetProtoFile())
	var buf bytes.Buffer
	err := generate.Run(&buf, source, lang, customHeaderIds, onlyExternal)
	if err != nil {
		// file does not include a service -> just don't do anything
		if err == generate.ErrNoServices {
			return &pluginpb.CodeGeneratorResponse{
				File: []*pluginpb.CodeGeneratorResponse_File{},
			}, nil
		}
		return nil, fmt.Errorf("failed markdown generation: %v", err)
	}

	if mode == modeAppend {
		// Load content from existing all.md
		e, err := ioutil.ReadFile(filepath.Join(dir, "all.md"))
		if err != nil && !os.IsNotExist(err) {
			return nil, err
		}
		buf = *bytes.NewBuffer(append(e, buf.Bytes()...))
	}
	formatted, err := markdownfmt.Process("", buf.Bytes())
	if err != nil {
		return nil, err
	}
	return &pluginpb.CodeGeneratorResponse{
		File: []*pluginpb.CodeGeneratorResponse_File{
			{
				Name:    proto.String(filepath.Join(base, "all.md")),
				Content: proto.String(string(formatted)),
			},
		},
	}, nil
}
