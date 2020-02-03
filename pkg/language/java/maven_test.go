package java

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMavenParse(t *testing.T) {
	pom := `
<?xml version="1.0" encoding="UTF-8"?>

<project xmlns="http://maven.apache.org/POM/4.0.0" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"
  xsi:schemaLocation="http://maven.apache.org/POM/4.0.0 http://maven.apache.org/xsd/maven-4.0.0.xsd">
  <modelVersion>4.0.0</modelVersion>

  <groupId>me.nicolaferraro.snap</groupId>
  <artifactId>example</artifactId>
  <version>1.0-SNAPSHOT</version>
  <packaging>jar</packaging>

  <properties>
    <project.build.sourceEncoding>UTF-8</project.build.sourceEncoding>
    <maven.compiler.source>1.7</maven.compiler.source>
    <maven.compiler.target>1.7</maven.compiler.target>
  </properties>

  <dependencies>
    <dependency>
      <groupId>junit</groupId>
      <artifactId>junit</artifactId>
      <version>4.11</version>
      <scope>test</scope>
    </dependency>
  </dependencies>

  <build>
  </build>
</project>
`
	model, err := parsePomData(pom)
	assert.Nil(t, err)
	assert.Equal(t, "me.nicolaferraro.snap", model.GetGroupID())
	assert.Equal(t, "example", model.GetArtifactID())
	assert.Equal(t, "1.0-SNAPSHOT", model.GetVersion())
	assert.Equal(t, "jar", model.GetPackaging())
	assert.Equal(t, "me.nicolaferraro.snap:example:1.0-SNAPSHOT", model.GetID())
}

func TestMavenParseWithParent(t *testing.T) {
	pom := `
<?xml version="1.0" encoding="UTF-8"?>

<project xmlns="http://maven.apache.org/POM/4.0.0" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"
  xsi:schemaLocation="http://maven.apache.org/POM/4.0.0 http://maven.apache.org/xsd/maven-4.0.0.xsd">
  <modelVersion>4.0.0</modelVersion>
  <parent>
    <groupId>me.nicolaferraro.snap</groupId>
    <artifactId>example-Parent</artifactId>
    <version>1.0-SNAPSHOT</version>
  </parent>
  <modelVersion>4.0.0</modelVersion>

  <artifactId>example</artifactId>

  <properties>
    <project.build.sourceEncoding>UTF-8</project.build.sourceEncoding>
    <maven.compiler.source>1.7</maven.compiler.source>
    <maven.compiler.target>1.7</maven.compiler.target>
  </properties>

  <dependencies>
    <dependency>
      <groupId>junit</groupId>
      <artifactId>junit</artifactId>
      <version>4.11</version>
      <scope>test</scope>
    </dependency>
  </dependencies>

  <build>
  </build>
</project>
`
	model, err := parsePomData(pom)
	assert.Nil(t, err)
	assert.Equal(t, "me.nicolaferraro.snap", model.GetGroupID())
	assert.Equal(t, "example", model.GetArtifactID())
	assert.Equal(t, "1.0-SNAPSHOT", model.GetVersion())
	assert.Equal(t, "", model.GetPackaging())
	assert.Equal(t, "me.nicolaferraro.snap:example:1.0-SNAPSHOT", model.GetID())
}
