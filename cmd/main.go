package main

import (
	"io/ioutil"
	"log"
	"os"

	"github.com/nicolaferraro/snap/v2/pkg/deployer/java"
	"github.com/nicolaferraro/snap/v2/pkg/publisher"
)

func main() {
	dir, err := ioutil.TempDir("", "snap-")
	if err != nil {
		log.Fatal("can't create temporary dir: ", err)
	}
	defer os.RemoveAll(dir)

	deployer := java.JavaDeployer{}
	if err := deployer.Deploy("./example", dir); err != nil {
		log.Fatal("error while creating deployment for java source: ", err)
	}

	pub := publisher.Publisher{}
	if err := pub.Publish(dir, "extensions/ext1"); err != nil {
		log.Fatal("cannot publish to server: ", err)
	}

	println("Uploaded!!")
}
