package commands

import (
	"fmt"
	"os"

	"encoding/json"

	"github.com/scalefast/talos/tools/csa"
	"github.com/scalefast/talos/tools/dast"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

type Settings struct {
	Csa  csa.Settings  `json:"csa" yaml:"CSA"`
	Dast dast.Settings `json:"dast" yaml:"DAST"`
}

var cmdAll = &cobra.Command{
	Use:   "all (--format [yaml|json]) (--o Path)",
	Short: "Generate all config settings",
	Long:  `Generate all config settings and output to a file for later use`,
	Run: func(cmd *cobra.Command, args []string) {

		s := new(Settings)

		s.Csa.Image = "#Required; Can be a local image, or a image from a container registry"

		s.Dast.ScanType = "#Required; one of website or api"
		s.Dast.TargetContainerNetworkID = "#Optional; Change if network is other than bridge"
		s.Dast.ImageName = "#Required for api scan type; Running Docker image"
		s.Dast.ZAPConfigFile = "#Required; ZAP config file location"
		s.Dast.OpenApiConfigFile = "#Required; Name of OpenAPI Config file"
		s.Dast.Auth.Website.AccessToken = "#Required for website; AccessToken stands for authentication header, sent in every request"
		s.Dast.Auth.Api.ClientID = "#Required for api; Will be used to get authorization from API"
		s.Dast.Auth.Api.ClientSecret = "#Required for api; Will be used to get authorization from API"
		s.Dast.Auth.Api.GrantType = "#Required for api; Will be used to get authorization from API"
		s.Dast.Auth.Api.Username = "#Required for api; Will be used to get authorization from API"
		s.Dast.Auth.Api.Password = "#Required for api; Will be used to get authorization from API"

		var data []byte
		if filetype == "yaml" {
			data, _ = yaml.Marshal(s)
		} else if filetype == "json" {
			data, _ = json.Marshal(s)
		} else {
			fmt.Print("Filetype not supported. Must be yaml or json")
		}
		err := os.WriteFile(filename, data, 0644)
		if err != nil {
			fmt.Printf("Error generating the file %v", err)
		}
	},
}
