package dast

import (
	"encoding/json"
	"errors"
	"io/fs"
	"os"
	"path/filepath"
)

type report struct {
	filepath string
	filename string
	data     *DASTReport
}

type DASTReport struct {
	Version   string `json:"@version"`
	Generated string `json:"@generated"`
	Hostname  string `json:"@hostname"`
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

// Wrapper function that manages all processes
// required for wazuh implementation.
func StoreReportForWazhuToDir(f string) (err error) {

	r, err := getReport()
	if err != nil {
		return errors.New("unable to get report: " + err.Error())
	}

	r.aggregateHostname()
	if err != nil {
		return errors.New("unable to aggregate hostname to report: " + err.Error())
	}

	r.writeTo(f)
	if err != nil {
		return errors.New("unable to write to " + f + ": " + err.Error())
	}

	return nil
}

// Returns a report that matches the 'talos-report.json'
// file description
func getReport() (r report, err error) {
	// Check if report file exists.
	pwd, _ := os.Getwd()
	filepath.WalkDir(pwd, func(s string, d fs.DirEntry, err error) error {
		//TODO: Prevent hardcoding the file name in the regexp
		// This forces us to keep config TALOS DAST file static.
		// Unflexible.
		// Check if file meets standards. (Ends in json, and is named a certain way)
		if !d.IsDir() {
			if d.Name() == "talos-report.json" {
				r = report{s, d.Name(), nil}
			}
		}
		return nil
	})

	if r.filepath == "" {
		return r, errors.New("could not find talos-report.json file")
	}

	file, err := os.ReadFile(r.filepath)
	if err != nil {
		return r, errors.New("could not read " + r.filepath + " file")
	}
	if json.Valid(file) {
		json.Unmarshal(file, &r.data)
		if err != nil {
			return r, errors.New("could not unmarshall report")
		}
	} else {
		return r, errors.New("invalid json")
	}
	return r, nil
}

// Adds field '@hostname' to json report
// So we can use the report with wazuh
func (r report) aggregateHostname() (err error) {

	hn, err := os.Hostname()
	if err != nil {
		hn = "Undefined"
	}
	r.data.Hostname = hn
	return nil
}

// Writes the report to path location
func (r report) writeTo(path string) (err error) {

	md, err := json.MarshalIndent(r.data, "", "  ")
	if err != nil {
		return errors.New("could not marshall new json file")
	}

	err = os.WriteFile(path, md, 0644)
	if err != nil {
		return errors.New("could not write " + path + " file")
	}
	return nil
}
