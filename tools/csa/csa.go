package csa

import (
	"context"
	"os"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/jsonmessage"
	"github.com/docker/docker/pkg/stdcopy"
	"github.com/moby/term"
	"github.com/scalefast/talos/logger"
	"github.com/spf13/viper"
)

// Analyze serves as the csa tool entrypoint
// its parameters are *viper.Viper and logger.StandardLogger
// viper provides direct access to the user settings, defined in
// the config.yaml file
// logger provides a interface to log messages into a default error
// This method does not return anything because the program finishes
// when this method does, so it has to be autocontained
func Analyze(c *viper.Viper, l *logger.StandardLogger, user string, network string, output string) {

	ctx2 := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}

	l.ICustom("Checking if we have the latest ZAP Docker image")
	events, err := cli.ImagePull(ctx2, "aquasec/trivy", types.ImagePullOptions{})
	if err != nil {
		l.ECustom("Error pulling the trivy docker image")
		os.Exit(1)
	}
	defer events.Close()
	termFd, isTerm := term.GetFdInfo(os.Stderr)
	jsonmessage.DisplayJSONMessagesStream(events, os.Stderr, termFd, isTerm, nil)

	var image string = c.GetString("csa.image")

	if image == "" {
		l.EMissingArg("csa.image")
		os.Exit(1) //If the int is "0" it generates a log with panic, if not does not generate that log. (This is why is "1")
	}

	//docker run -u root -v /var/run/docker.sock:/var/run/docker.sock aquasec/trivy:latest testimage:latest (Command that really is executed)
	resp, err := cli.ContainerCreate(ctx2, &container.Config{
		Image: "aquasec/trivy:latest",
		User:  user,
		Cmd:   []string{"image", image},
		Tty:   false,
	}, &container.HostConfig{
		Binds:       []string{"/var/run/docker.sock:/var/run/docker.sock"},
		NetworkMode: container.NetworkMode(network),
	}, nil, nil, "")
	if err != nil {
		panic(err)
	}

	if err := cli.ContainerStart(ctx2, resp.ID, types.ContainerStartOptions{}); err != nil {
		panic(err)
	}

	statusCh, errCh := cli.ContainerWait(ctx2, resp.ID, container.WaitConditionNotRunning)
	select {
	case err := <-errCh:
		if err != nil {
			panic(err)
		}
	case <-statusCh:
	}

	out, err := cli.ContainerLogs(ctx2, resp.ID, types.ContainerLogsOptions{ShowStdout: true})
	if err != nil {
		panic(err)
	}

	f, err := os.Create(output)
	if err != nil {
		l.ECustom("Error, could not create the file")
		panic(err)
	}
	l.ICustom("Report saved to file report.txt")
	stdcopy.StdCopy(f, os.Stderr, out)
	defer f.Close()
}
