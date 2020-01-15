# snap

Simple POC to upload a Maven example to Minio.

## Setup

Checkout this git repo.

Connect to a Openshift namespace, then apply:

```
oc apply -f deploy/minio-standalone-pvc.yaml
oc apply -f deploy/minio-standalone-deployment.yaml
oc apply -f deploy/minio-standalone-service.yaml
oc expose service/minio-service
```

Get the Minio endpoint from the Route:

```
export MINIO_ENDPOINT=$(oc get route minio-service -o jsonpath='{.spec.host}')
```

Build the application with:
```
go build -o snap ./cmd/
```

## Publish

To deploy the example project:

```
./snap
```

Should be successful.

## Verify

Download the `mc` tool from Minio.

```
mc config host add minio http://$MINIO_ENDPOINT minio minio123
```

Example result:
```
[nferraro@localhost snap]$ mc ls --recursive minio
[2020-01-15 12:56:51 CET]  2.8KiB extensions/ext1/me/nicolaferraro/snap/example/1.0-SNAPSHOT/example-1.0-20200115.115652-1.jar
[2020-01-15 12:56:51 CET]     32B extensions/ext1/me/nicolaferraro/snap/example/1.0-SNAPSHOT/example-1.0-20200115.115652-1.jar.md5
[2020-01-15 12:56:51 CET]     40B extensions/ext1/me/nicolaferraro/snap/example/1.0-SNAPSHOT/example-1.0-20200115.115652-1.jar.sha1
[2020-01-15 12:56:51 CET]  2.6KiB extensions/ext1/me/nicolaferraro/snap/example/1.0-SNAPSHOT/example-1.0-20200115.115652-1.pom
[2020-01-15 12:56:51 CET]     32B extensions/ext1/me/nicolaferraro/snap/example/1.0-SNAPSHOT/example-1.0-20200115.115652-1.pom.md5
[2020-01-15 12:56:51 CET]     40B extensions/ext1/me/nicolaferraro/snap/example/1.0-SNAPSHOT/example-1.0-20200115.115652-1.pom.sha1
[2020-01-15 12:56:51 CET]    772B extensions/ext1/me/nicolaferraro/snap/example/1.0-SNAPSHOT/maven-metadata.xml
[2020-01-15 12:56:51 CET]     32B extensions/ext1/me/nicolaferraro/snap/example/1.0-SNAPSHOT/maven-metadata.xml.md5
[2020-01-15 12:56:51 CET]     40B extensions/ext1/me/nicolaferraro/snap/example/1.0-SNAPSHOT/maven-metadata.xml.sha1
[2020-01-15 12:56:51 CET]    286B extensions/ext1/me/nicolaferraro/snap/example/maven-metadata.xml
[2020-01-15 12:56:51 CET]     32B extensions/ext1/me/nicolaferraro/snap/example/maven-metadata.xml.md5
[2020-01-15 12:56:51 CET]     40B extensions/ext1/me/nicolaferraro/snap/example/maven-metadata.xml.sha1

```

