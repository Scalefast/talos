package dast

import (
	"archive/tar"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/signal"
	"path/filepath"
	"regexp"
	"strings"
	"syscall"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/archive"
	"github.com/docker/docker/pkg/jsonmessage"
	"github.com/docker/go-connections/nat"
	"github.com/moby/term"
	"github.com/spf13/viper"
	"github.com/scalefast/talos/logger"
)

// Manages user <Ctrl-C>
// Stops running Docker container
func SetupapiCloseHandler(cli *client.Client, cid string) {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		fmt.Println("\r- Ctrl+C pressed in Terminal")
		if err := stopAndRemoveContainer(cli, cid); err != nil {
			fmt.Println("Error stopping container")
			fmt.Println("Error: " + err.Error())
		}
		os.Exit(0)
	}()
}

// Manages user <Ctrl-C>
// Stops running Docker container
func SetupCloseHandler() {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		fmt.Println("\r- Ctrl+C pressed in Terminal")
		os.Exit(0)
	}()
}

// Stop and remove a container
func stopAndRemoveContainer(client *client.Client, containername string) error {
	ctx := context.Background()

	if err := client.ContainerStop(ctx, containername, nil); err != nil {
		fmt.Printf("Unable to stop container %s: %s", containername, err)
	}

	removeOptions := types.ContainerRemoveOptions{
		RemoveVolumes: true,
		Force:         true,
	}

	if err := client.ContainerRemove(ctx, containername, removeOptions); err != nil {
		fmt.Printf("Unable to remove container: %s", err)
		return err
	}

	return nil
}

// Returns a container, so we can access its information.
func getContainer(ctx context.Context, cli client.Client, ContainerName string) *types.Container {
	containers, err := cli.ContainerList(ctx, types.ContainerListOptions{})
	if err != nil {
		panic(err)
	}
	for _, cnt := range containers {
		if ContainerName == cnt.Names[0][1:] {
			// The container specified exists.
			// We have to get the container ID, to see it's settings, get
			// The network it's connected to, and then connect
			// the new contanier to it
			return &cnt
		}
	}
	return nil
}

// Returns the network ID the container is attached to, or a empty network and a error
func getNetworkID(ctx context.Context, cli client.Client, networkName string) (networkID string, err error) {

	// Get available networks
	n, err := cli.NetworkList(ctx, types.NetworkListOptions{})
	if err != nil {
		return "", errors.New("dast.network has invalid value")
	}

	// Check if the network specified exists.
	for _, nw := range n {
		if networkName == nw.Name {
			// The network specified exists, let's see if the Docker image is
			// Connected to this network.
			networkID = nw.ID
		}
	}
	return networkID, err
}

// createZAPContainer creates a Docker container based on the settings
// passed through the config file
func createZAPContainer(ctx context.Context, cli *client.Client, config *viper.Viper) (c container.ContainerCreateCreatedBody) {
	containerConfig := container.Config{
		Hostname:     "",
		Domainname:   "",
		User:         "",
		AttachStdin:  false,
		AttachStdout: false,
		AttachStderr: false,
		ExposedPorts: map[nat.Port]struct{}{},
		Tty:          false,
		OpenStdin:    false,
		StdinOnce:    false,
		Env:          []string{},
		Cmd:          []string{"bash", "/zap/helpers/entrypoint.sh"},
		//Cmd:             []string{"ls", "-la", "/zap"},
		Healthcheck:     &container.HealthConfig{},
		ArgsEscaped:     false,
		Image:           "owasp/zap2docker-weekly",
		Volumes:         map[string]struct{}{},
		WorkingDir:      "",
		Entrypoint:      []string{},
		NetworkDisabled: false,
		MacAddress:      "",
		OnBuild:         []string{},
		Labels:          map[string]string{},
		StopSignal:      "",
		StopTimeout:     new(int),
		Shell:           []string{},
	}
	//docker run --network=bridge owasp/zap2docker-weekly ./zap.sh -autorun config.yaml
	c, err := cli.ContainerCreate(ctx, &containerConfig, nil, nil, nil, "")
	if err != nil {
		panic(err)
	}
	return c
}

func InitSettings(c *viper.Viper) (s *Settings, err error) {

	s = new(Settings)

	//Initialize context and new client
	s.ctx = context.Background()
	s.cli, err = client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, err
	}
	s.ScanType = strings.ToLower(c.GetString("dast.ScanType"))

	// Check if the ZAPConfigFileName is specified, if not, exit.
	// We need this parameter
	if ok := c.GetString("dast.ZAPConfigFileName"); ok == "" {
		return s, errors.New("please specify field ZAPConfigFileName under 'dast' key in talos config file")
	}

	switch s.ScanType {
	case "website":
		err = SetwebsiteSettings(s, c)
		if err != nil {
			return nil, err
		}
	case "api":
		s.targetContainer = getContainer(s.ctx, *s.cli, c.GetString("dast.ImageName"))
		if s.targetContainer == nil {
			return nil, errors.New("No container named: " + c.GetString("dast.ImageName") + " running")
		}
		s.OpenApiConfigFile = filepath.Join(c.GetString("dast.OpenApiConfigFileDir"), c.GetString("dast.OpenApiConfigFileName"))
		if ok := exists(s.OpenApiConfigFile); !ok {
			return nil, errors.New("OpenApi config file location" + s.OpenApiConfigFile)
		}
		err = SetapiSettings(s, c)
		if err != nil {
			return nil, err
		}
	default:
		return nil, errors.New("dast.scantype value must be one of website or api")
	}

	// Check if Docker network is set in parameters
	// If the network has not been set, then we assume it's the default network
	n := c.GetString("dast.network")
	if n == "" {
		n = "bridge"
	}

	// Check if the network exists, and get it's ID
	// We cannot do a scan to the Image, if we are not attached to the same network
	// So we will exit the program if it is not defined properly.
	if netID, err := getNetworkID(s.ctx, *s.cli, n); err == nil {
		s.TargetContainerNetworkID = netID
		s.ZAPContainer = createZAPContainer(s.ctx, s.cli, c)
	} else {
		return nil, errors.New("network dastnetwork")
	}

	s.ZAPConfigFile = filepath.Join(c.GetString("dast.ZAPConfigFileDir"), c.GetString("dast.ZAPConfigFileName"))
	if ok := exists(s.ZAPConfigFile); !ok {
		return nil, errors.New("File not found: " + s.ZAPConfigFile)
	}

	return s, err
}

func copyFromContainer(ctx context.Context, cli *client.Client, l *logger.StandardLogger, src string, dst string, cont string) (err error) {
	// Create a TAR file with the config, so we can copy it to the container
	// Copy the file to the container

	// _ FileStats, can be discarded, as we know what we want
	// Maybe, parse dast config file and retrieve the location configured.
	// Used for last parameter.
	tarReader, _, err := cli.CopyFromContainer(ctx, cont, src)

	if err != nil {
		panic(err)
	}
	tr := tar.NewReader(tarReader)

	re, err := regexp.Compile("^zap/report/*")

	for {
		header, err := tr.Next()

		switch {

		case err == io.EOF: // if no more files are found return
			return nil
		case err != nil: // return any other error
			l.ECustom("Error: " + err.Error())
			panic(err)
		case header == nil: // if the header is nil, just skip it (not sure how this happens)
			continue
		}

		if match := re.Match([]byte(header.Name)); match {
			// the target location where the dir/file should be created
			target := filepath.Join(dst, header.Name)
			l.ICustom("Report being saved at: " + target)

			// check the file type
			switch header.Typeflag {

			// if its a dir and it doesn't exist create it
			case tar.TypeDir:
				if _, err := os.Stat(target); err != nil {
					if err := os.MkdirAll(target, 0755); err != nil {
						panic(err)
					}
				}

			// if it's a file create it
			case tar.TypeReg:
				f, err := os.OpenFile(target, os.O_CREATE|os.O_RDWR, os.FileMode(header.Mode))
				if err != nil {
					panic(err)
				}

				// copy over contents
				if _, err := io.Copy(f, tr); err != nil {
					l.ECustom("Error: " + err.Error())
				}

				// manually close here after each file operation; defering would cause each file close
				// to wait until all operations have completed.
				f.Close()
			}

		}
	}
}

// TODO: Needs refactor.
// Create a single tar archive and upload only that, instead
// of three archives.
func transferFilesToContainer(cli *client.Client, s *Settings, f string) (err error) {
	ctx := context.Background()
	fmt.Println("ZAP Config file: " + s.ZAPConfigFile)
	zcf, status := archive.Tar(s.ZAPConfigFile, 0)
	if status != nil {
		return errors.New("Error: " + status.Error())
	}
	fmt.Println("Reached ZAPConfigFile copy: " + s.ZAPConfigFile)
	defer zcf.Close()
	if err = cli.CopyToContainer(ctx, s.ZAPContainer.ID, "/zap", zcf, types.CopyToContainerOptions{}); err != nil {
		return errors.New("unable to copy DAST config file to ZAP container")
	}
	ocf, status := archive.Tar(s.OpenApiConfigFile, 0)
	if status != nil {
		return errors.New("Error: " + status.Error())
	}
	fmt.Println("Reached OpenAPIConfigFile copy: " + s.OpenApiConfigFile)
	defer ocf.Close()
	if err = cli.CopyToContainer(ctx, s.ZAPContainer.ID, "/zap", ocf, types.CopyToContainerOptions{}); err != nil {
		return errors.New("unable to copy DAST config file to ZAP container")
	}
	// Copy helper scripts to tar file
	tf, status := archive.Tar(f, 0)
	if status != nil {
		return errors.New("Error: " + status.Error())
	}
	fmt.Println("Reached Helpers copy: " + f)
	fi, err := os.Lstat(f)
	fmt.Println("Check if file exists: " + fi.Name())
	defer tf.Close()
	if copyErr := cli.CopyToContainer(ctx, s.ZAPContainer.ID, "/zap", tf, types.CopyToContainerOptions{AllowOverwriteDirWithFile: true}); copyErr != nil {
		return errors.New(copyErr.Error())
	}

	return nil
}

func exists(path string) bool {
	_, err := os.Stat(path)
	return !errors.Is(err, os.ErrNotExist)
}

func (s *Settings) UpdateImage() (err error) {
	// Pull the Docker image from the container registry
	events, err := s.cli.ImagePull(s.ctx, "owasp/zap2docker-weekly", types.ImagePullOptions{})
	if err != nil {
		return errors.New("error pulling the zap docker image")
	}
	defer events.Close()
	termFd, isTerm := term.GetFdInfo(os.Stderr)
	jsonmessage.DisplayJSONMessagesStream(events, os.Stderr, termFd, isTerm, nil)
	return nil
}
