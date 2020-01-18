package java

import (
	"fmt"
	"io"
	"os/exec"

	"github.com/nicolaferraro/snap/pkg/deployer"
	"github.com/nicolaferraro/snap/pkg/util/log"
)

type JavaDeployer struct {
	stdOut io.Writer
	stdErr io.Writer
}

var (
	logger = log.WithName("java-deployer")
)

func NewJavaDeployer(stdOut, stdErr io.Writer) deployer.Deployer {
	return &JavaDeployer{
		stdOut: stdOut,
		stdErr: stdErr,
	}
}

func (d *JavaDeployer) Deploy(source, destination string) error {
	logger.Infof("Executing maven release phase on project %s", source)
	cmd := exec.Command("./mvnw", "deploy", "-DskipTests", fmt.Sprintf("-DaltDeploymentRepository=snapshot-repo::default::file:%s", destination))
	cmd.Dir = source
	cmd.Stdout = d.stdOut
	cmd.Stderr = d.stdErr
	err := cmd.Run()
	if err != nil {
		return err
	}
	logger.Infof("Maven release phase completed for project %s", source)
	return nil
}
