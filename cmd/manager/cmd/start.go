package cmd

import (
	"strings"

	"github.com/mirzakhany/pm/cmd/manager/engine"
	"github.com/mirzakhany/pm/cmd/manager/templates"
	"github.com/spf13/cobra"
)

func start(serviceName string) {

	n := strings.ToLower(serviceName + "s")
	data := engine.Data{
		Pkg: engine.Package{
			Name:        serviceName,
			NamePlural:  serviceName + "s",
			EntityAlias: string(n[0]),
		},
	}

	tmpls := []struct {
		Tmpl     string
		Dest     string
		Filename string
	}{
		{
			Tmpl:     templates.EntityTmpl,
			Dest:     "./internal/entity/",
			Filename: n + ".go",
		},
		{
			Tmpl:     templates.ProtoTmpl,
			Dest:     "./protobuf/" + n,
			Filename: n + ".proto",
		},
		{
			Tmpl:     templates.ApiTmpl,
			Dest:     "./internal/" + n,
			Filename: "api.go",
		},
		{
			Tmpl:     templates.RepositoryTmpl,
			Dest:     "./internal/" + n,
			Filename: "repository.go",
		},
		{
			Tmpl:     templates.RepositoryTestTmpl,
			Dest:     "./internal/" + n,
			Filename: "repository_test.go",
		},
		{
			Tmpl:     templates.ServiceTmpl,
			Dest:     "./internal/" + n,
			Filename: "service.go",
		},
		{
			Tmpl:     templates.ServiceTestTmpl,
			Dest:     "./internal/" + n,
			Filename: "/service_test.go",
		},
	}

	var err error
	for _, tmpl := range tmpls {
		err = engine.Render(
			tmpl.Tmpl,
			tmpl.Dest,
			tmpl.Filename,
			data,
		)
		if err != nil {
			panic(err)
		}
	}
}

var startCmd = &cobra.Command{
	Use:   "start [Service name]",
	Short: "create a service",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		start(args[0])
	},
}

func init() {
	rootCmd.AddCommand(startCmd)
}
