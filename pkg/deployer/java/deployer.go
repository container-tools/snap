package java

import (
	"fmt"
	"os"
	"os/exec"
)

type JavaDeployer struct{}

func (*JavaDeployer) Deploy(source, destination string) error {
	cmd := exec.Command("./mvnw", "deploy", "-DskipTests", fmt.Sprintf("-DaltDeploymentRepository=snapshot-repo::default::file:%s", destination))
	cmd.Dir = source
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
