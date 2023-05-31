package main

import (
	"log"

	"github.com/nikhilsbhat/opensource-insight-exporter/cmd"
	"github.com/nikhilsbhat/opensource-insight-exporter/pkg/common"
	docgen "github.com/nikhilsbhat/urfavecli-docgen"
)

//go:generate go run github.com/nikhilsbhat/opensource-insight-exporter/docs
func main() {
	if err := docgen.GenerateDocs(cmd.App(), common.OpenSourceInsightExporterName); err != nil {
		log.Fatalln(err)
	}
}
