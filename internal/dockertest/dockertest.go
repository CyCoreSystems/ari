package dockertest

// from https://divan.github.io/posts/integration_testing/

// build: +test

import (
	"bytes"
	"fmt"
	"log"
	"net"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/fsouza/go-dockerclient"
)

// waitReachable waits for hostport to became reachable for the maxWait time.
func waitReachable(hostport string, maxWait time.Duration) error {

	time.Sleep(maxWait)

	done := time.Now().Add(maxWait)
	for time.Now().Before(done) {
		c, err := net.Dial("tcp", hostport)
		if err == nil {
			c.Close()
			return nil
		}
		time.Sleep(100 * time.Millisecond)
	}
	return fmt.Errorf("cannot connect %v for %v", hostport, maxWait)
}

// waitStarted waits for a container to start for the maxWait time.
func waitStarted(client *docker.Client, id string, maxWait time.Duration) error {
	done := time.Now().Add(maxWait)
	for time.Now().Before(done) {
		c, err := client.InspectContainer(id)
		if err != nil {
			break
		}
		if c.State.Running {
			return nil
		}
		time.Sleep(100 * time.Millisecond)
	}
	return fmt.Errorf("cannot start container %s for %v", id, maxWait)
}

// CreateOptions creates the options for launching the test asterisk container
func CreateOptions() docker.CreateContainerOptions {
	ports := make(map[docker.Port]struct{})
	ports["8088"] = struct{}{}
	opts := docker.CreateContainerOptions{
		Config: &docker.Config{
			Image:        "test-asterisk:13.8",
			ExposedPorts: ports,
		},
	}

	return opts
}

// StartOptions creates the options for starting the docker container
func StartOptions(bindPorts bool) *docker.HostConfig {
	portBinds := make(map[docker.Port][]docker.PortBinding)
	if bindPorts {
		portBinds["8088"] = []docker.PortBinding{
			docker.PortBinding{HostPort: "8088"},
		}
	}
	conf := docker.HostConfig{
		PortBindings: portBinds,
	}

	return &conf
}

// StartAsterisk starts asterisk in docker, returning the URL to asterisk and the cancel function
func StartAsterisk() (string, func(), error) {
	client, err := docker.NewClientFromEnv()
	if err != nil {
		return "", nil, err
	}
	c, err := client.CreateContainer(CreateOptions())
	if err != nil {
		return "", nil, err
	}
	deferFn := func() {
		if err := client.RemoveContainer(docker.RemoveContainerOptions{
			ID:    c.ID,
			Force: true,
		}); err != nil {
			log.Printf("cannot remove container: %v\n", err)
		}
	}

	// VM IP is the IP of DockerMachine VM, if running (used in non-Linux OSes)
	vmIP := strings.TrimSpace(DockerMachineIP())
	nonLinux := (vmIP != "")

	err = client.StartContainer(c.ID, StartOptions(nonLinux))
	if err != nil {
		deferFn()
		return "", nil, err
	}

	// wait for container to wake up
	if err := waitStarted(client, c.ID, 10*time.Second); err != nil {
		deferFn()
		return "", nil, err
	}
	if c, err = client.InspectContainer(c.ID); err != nil {
		deferFn()
		return "", nil, err
	}

	// determine IP address for MySQL
	ip := ""
	if vmIP != "" {
		ip = vmIP
	} else if c.NetworkSettings != nil {
		ip = strings.TrimSpace(c.NetworkSettings.IPAddress)
	}

	// wait asterisk ARI to wake up
	if err := waitReachable(ip+":8088", 10*time.Second); err != nil {
		deferFn()
		return "", nil, err
	}

	url := ariURL(ip)

	return url, deferFn, nil
}

// ariURL returns valid url to be used with the ARI client
func ariURL(ip string) string {
	return fmt.Sprintf("http://%s:8088/ari", ip)
}

// DockerMachineIP returns IP of docker-machine or boot2docker VM instance.
//
// If docker-machine or boot2docker is running and has IP, it will be used to
// connect to dockerized services (MySQL, etc).
//
// Basically, it adds support for MacOS X and Windows.
func DockerMachineIP() string {
	// Docker-machine is a modern solution for docker in MacOS X.
	// Try to detect it, with fallback to boot2docker
	var dockerMachine bool
	machine := os.Getenv("DOCKER_MACHINE_NAME")
	if machine != "" {
		dockerMachine = true
	}

	var buf bytes.Buffer

	var cmd *exec.Cmd
	if dockerMachine {
		cmd = exec.Command("docker-machine", "ip", machine)
	} else {
		cmd = exec.Command("boot2docker", "ip")
	}
	cmd.Stdout = &buf

	if err := cmd.Run(); err != nil {
		// ignore error, as it's perfectly OK on Linux
		return ""
	}

	return buf.String()
}
