package dast

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	tester "github.com/scalefast/talos/tester"
)

// Test the report is not available
func TestStoreReportToWazuhDirDoesntFindFile(t *testing.T) {

	nt := tester.NewTester(t)
	wd := nt.TempDir()
	expected := "unable to get report: could not find talos-report.json file"

	nt.Given("There is no file to be read")
	{
		nt.When("We want to retrieve the report")
		{
			err := StoreReportForWazhuToDir(wd)
			nt.Then("A error should appear, stating we where unable to retrieve the report")
			{
				if strings.Compare(err.Error(), expected) != 0 {
					nt.Error("error", expected, err.Error())
				}
			}
		}
	}
}

// Report that OWASP ZAP generates
type DASTOficialReport struct {
	Version   string `json:"@version"`
	Generated string `json:"@generated"`
	Site      []struct {
		Name   string `json:"@name"`
		Host   string `json:"@host"`
		Port   string `json:"@port"`
		Ssl    string `json:"@ssl"`
		Alerts []struct {
			Pluginid   string `json:"pluginid"`
			AlertRef   string `json:"alertRef"`
			Alert      string `json:"alert"`
			Name       string `json:"name"`
			Riskcode   string `json:"riskcode"`
			Confidence string `json:"confidence"`
			Riskdesc   string `json:"riskdesc"`
			Desc       string `json:"desc"`
			Instances  []struct {
				URI      string `json:"uri"`
				Method   string `json:"method"`
				Param    string `json:"param"`
				Attack   string `json:"attack"`
				Evidence string `json:"evidence"`
			} `json:"instances"`
			Count     string `json:"count"`
			Solution  string `json:"solution"`
			Otherinfo string `json:"otherinfo"`
			Reference string `json:"reference"`
			Cweid     string `json:"cweid"`
			Wascid    string `json:"wascid"`
			Sourceid  string `json:"sourceid"`
		} `json:"alerts"`
	} `json:"site"`
}

// TODO: Finish; wd and tr are not working correctly.
// Test the json file is available and valid
func TestStoreReportToWazuhDirValidJson(t *testing.T) {

	nt := tester.NewTester(t)
	td, _ := os.Getwd()
	tr := filepath.Join(td, "talos-report.json")

	nt.Given("We have a valid json file")
	{
		// Create a file called talos-report.json
		// that does not have a Hostname field
		// because OWASP ZAP report will be generated
		// without the report key.
		data := &DASTOficialReport{Version: "V1.0", Generated: "today"}
		f, _ := os.OpenFile(tr, os.O_RDWR|os.O_CREATE, 0755)
		json.NewEncoder(f).Encode(data)
		defer f.Close()
		// Arrange for project cleanup
		nt.Cleanup(func() {
			os.Remove(tr)
		})
		nt.When("We aggregate the hostname")
		{
			wd := filepath.Join(nt.TempDir(), "talos-report.json")
			err := StoreReportForWazhuToDir(wd)
			if err != nil {
				nt.Error("error", "No error", err.Error())
			}
			nt.Then("The report should contain a hostname key")
			{
				// Check the resulting file has a field named: hostname
				// And has a correct the Version (Ensures json has been
				// parsed correctly)
				generated, _ := os.ReadFile(wd)
				// Unmarshall the report, to see what fields it has.
				genStruct := &DASTReport{}
				json.Unmarshal(generated, genStruct)
				hn, _ := os.Hostname()
				if strings.Compare(genStruct.Hostname, hn) != 0 {
					nt.Error("hostname not present or invalid value", hn, genStruct.Hostname)
				}
				nt.And("The report shound have the correct value for Generated key")
				{
					if strings.Compare(genStruct.Generated, "today") != 0 {
						nt.Error("value for Generated invalid or not initialized", "today", genStruct.Generated)
					}
				}
				nt.And("The report shound have the correct value for Version key")
				{
					if strings.Compare(genStruct.Version, "V1.0") != 0 {
						nt.Error("value for Version invalid or not initialized", "V1.0", genStruct.Version)
					}
				}
				nt.Cleanup(func() {
					os.Remove(filepath.Join(os.TempDir(), "talos-report.json"))
				})
			}
		}
	}

}
