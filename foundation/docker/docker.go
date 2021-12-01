// Package docker provides support for starting and stopping docker containers
// for running tests.
package docker

import (
	"bytes"
	"encoding/json"
	"net"
	"os/exec"
	"testing"
)

// Container tracks information about the docker container started for tests.
type Container struct {
	ID   string
	Host string // IP:Port
}

// StartContainer starts the specified container for running tests.
// port parameter is the internal container port
func StartContainer(t *testing.T, image string, port string, args ...string) *Container {
	arg := []string{"run", "-P", "-d"}
	arg = append(arg, args...)
	arg = append(arg, image)

	// build the command
	// set the commands stdOut to a buffer so we can get output of the run
	cmd := exec.Command("docker", arg...)
	var out bytes.Buffer
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		// fail the test if the docker command fails
		// much better than returning an error
		t.Fatalf("could not start container %s: %v", image, err)
	}

	// get the container id
	id := out.String()[:12]

	// take the id and run it through docker inspect
	// docker inspect returns a json document
	cmd = exec.Command("docker", "inspect", id)
	out.Reset()
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		t.Fatalf("could not inspect container %s: %v", id, err)
	}

	// unmarshal the docker inspect response into a map
	var doc []map[string]interface{}
	if err := json.Unmarshal(out.Bytes(), &doc); err != nil {
		t.Fatalf("could not decode json: %v", err)
	}

	// extract the external port docker assigned to map to the parameterized internal port
	ip, randPort := extractIPPort(t, doc, port)

	c := Container{
		ID:   id,
		Host: net.JoinHostPort(ip, randPort),
	}

	t.Logf("Image:       %s", image)
	t.Logf("ContainerID: %s", c.ID)
	t.Logf("Host:        %s", c.Host)

	return &c
}

// StopContainer stops and removes the specified container.
func StopContainer(t *testing.T, id string) {
	if err := exec.Command("docker", "stop", id).Run(); err != nil {
		t.Fatalf("could not stop container: %v", err)
	}
	t.Log("Stopped:", id)

	// -v flag removes the temporary volume that is created
	if err := exec.Command("docker", "rm", id, "-v").Run(); err != nil {
		t.Fatalf("could not remove container: %v", err)
	}
	t.Log("Removed:", id)
}

// DumpContainerLogs outputs logs from the running docker container.
func DumpContainerLogs(t *testing.T, id string) {
	out, err := exec.Command("docker", "logs", id).CombinedOutput()
	if err != nil {
		t.Fatalf("could not log container: %v", err)
	}
	t.Logf("Logs for %s\n%s:", id, out)
}

func extractIPPort(t *testing.T, doc []map[string]interface{}, port string) (string, string) {
	nw, exists := doc[0]["NetworkSettings"]
	if !exists {
		t.Fatal("could not get network settings")
	}
	// type assert nw to another map[string]interface{}
	ports, exists := nw.(map[string]interface{})["Ports"]
	if !exists {
		t.Fatal("could not get network port settings")
	}
	tcp, exists := ports.(map[string]interface{})[port+"/tcp"]
	if !exists {
		t.Fatal("could not get network ports/tcp settings")
	}
	// type assert to a collection using the empty interface
	list, exists := tcp.([]interface{})
	if !exists {
		t.Fatal("could not get network ports/tcp list settings")
	}

	var hostIP string
	var hostPort string

	// range over the collection
	for _, l := range list {
		data, exists := l.(map[string]interface{})
		if !exists {
			t.Fatal("could not get network ports/tcp list data")
		}
		hostIP = data["HostIp"].(string)
		if hostIP != "::" {
			hostPort = data["HostPort"].(string)
		}
	}

	return hostIP, hostPort
}
