package dast

import (
	"embed"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/pkg/stdcopy"
	"github.com/spf13/viper"
	"github.com/scalefast/talos/logger"
)

// content holds our helper files
//go:embed helpers/*
// These files will be copied to the Docker conteiner
// So they must be embedded into the Talos binary
var helpers embed.FS

// Analyze serves as the dast tool entrypoint
// its parameters are *viper.Viper and logger.StandardLogger
// viper provides direct access to the user settings, defined in
// the config.yaml file
// logger provides a interface to log messages into a default error
// This method does not return anything because the program finishes
// when this method does, so it has to be autocontained
func Analyze(c *viper.Viper, l *logger.StandardLogger) {

	// Initialize the settings used in this script
	s, err := InitSettings(c)
	if err != nil {
		l.ECustom("Could not initialize settings: " + err.Error())
		os.Exit(1)
	}

	err = s.UpdateImage()
	if err != nil {
		l.ECustom("Could not update Docker image: " + err.Error())
		os.Exit(1)
	}

	// Manage started container.
	// It stops and removes the running container
	// if the user presses <Ctrl+C>
	if strings.Compare(s.ScanType, "api") == 0 {
		SetupapiCloseHandler(s.cli, s.ZAPContainer.ID)
	} else {
		SetupCloseHandler()
	}

	// If scan type is api, Generate the .env file and auth_post_request.json
	// to configure the DAST technique we are using for the Docker container.
	// If we are using normal DAST against web stores, then .env file
	// should have a Basic token, instead of the Bearer token configured.
	// Use Talos config.yaml file to set the content of the ACCES_TOKEN
	// env variable.
	// If there is no key in config.yaml file, then we assume we are
	// going to run a api-API-DAST analysis
	f, err := PrepareContainerEnv(s)
	l.Debug("Reached end of prepareContainereEnv!")
	if err != nil {
		l.ECustom("Error: DAST env could not be correctly set-up")
		l.ECustom(err.Error())
		os.Exit(1)
	}
	defer os.RemoveAll(f)

	// Should transfer files f to container.
	// All files should have been stored in
	// f directory
	l.Debug("Path to create tmp dir: " + f)
	if err = transferFilesToContainer(s.cli, s, f); err != nil {
		l.ECustom("error: " + err.Error())
		l.ECustom("error: Files could not be transfered into Docker container")
		os.Exit(1)
	}

	// Connect the zap docker container to the network the docker container we want to test
	// is in
	if err = s.cli.NetworkConnect(s.ctx, s.TargetContainerNetworkID, s.ZAPContainer.ID, &network.EndpointSettings{}); err != nil {
		l.ECustom(err.Error())
		l.EInvalidValue("Network ID", s.TargetContainerNetworkID)
		os.Exit(1)
	}

	// Start the container and perform the tests
	if err = s.cli.ContainerStart(s.ctx, s.ZAPContainer.ID, types.ContainerStartOptions{}); err != nil {
		l.ECustom("Error: ZAP container was unable to start" + err.Error())
		os.Exit(1)
	}
	l.ICustom("DAST Tests have started, they may take a while, be patient")

	attach, err := s.cli.ContainerAttach(s.ctx, s.ZAPContainer.ID, types.ContainerAttachOptions{Stream: true, Stdin: false, Stdout: true, Stderr: true})
	if err != nil {
		l.ECustom("error attaching to container!" + err.Error())
	}
	defer attach.Close()
	go io.Copy(os.Stdout, attach.Reader)

	// TODO: Feature: Add a progress bar, or something.
	statusCh, errCh := s.cli.ContainerWait(s.ctx, s.ZAPContainer.ID, container.WaitConditionNotRunning)
	select {
	case err = <-errCh:
		if err != nil {
			l.ICustom("error:" + err.Error())
		}
	case <-statusCh:
		l.ICustom("Testing...")
	}

	out, err := s.cli.ContainerLogs(s.ctx, s.ZAPContainer.ID, types.ContainerLogsOptions{ShowStdout: true})
	if err != nil {
		l.ECustom("Container stopped unexpectedly")
		os.Exit(1)
	}
	// Change to store logs to a file
	stdcopy.StdCopy(os.Stdout, os.Stderr, out)

	// Print user report
	pwd, _ := os.Getwd()
	if err = copyFromContainer(s.ctx, s.cli, l, "/zap/", pwd, s.ZAPContainer.ID); err != nil {
		l.ECustom("Could not transfer report from container")
		os.Exit(1)
	}

	wazuhDir := filepath.Join(os.TempDir(), "talos-report.json")
	StoreReportForWazhuToDir(wazuhDir)
	if err != nil {
		l.ECustom("error: " + err.Error())
		os.Exit(1)
	}
	l.ICustom("Report stored for monitoring")
}
