# snap

Simple POC to upload a Maven example to Minio.

## Setup

Checkout this git repo.

Connect to a Openshift namespace.

Build the application with:
```
make
```

## Publish

To deploy the example project:

```
./snap
```

Should be successful.

## Verify

Expose the minio service as route:

```
oc expose service/minio-service
```

Get the Minio endpoint from the Route:

```
export MINIO_ENDPOINT=$(oc get route minio-service -o jsonpath='{.spec.host}')
```

Download the `mc` tool from Minio, then:

```
mc config host add minio http://$MINIO_ENDPOINT minio minio123
```

List the data:
```
mc ls --recursive minio
```

Example result:
```
[2020-01-18 01:01:03 CET]  2.8KiB snap/me/nicolaferraro/snap/example/1.0-SNAPSHOT/example-1.0-20200118.000050-1.jar
[2020-01-18 01:01:03 CET]     32B snap/me/nicolaferraro/snap/example/1.0-SNAPSHOT/example-1.0-20200118.000050-1.jar.md5
[2020-01-18 01:01:03 CET]     40B snap/me/nicolaferraro/snap/example/1.0-SNAPSHOT/example-1.0-20200118.000050-1.jar.sha1
[2020-01-18 01:01:03 CET]  2.6KiB snap/me/nicolaferraro/snap/example/1.0-SNAPSHOT/example-1.0-20200118.000050-1.pom
[2020-01-18 01:01:03 CET]     32B snap/me/nicolaferraro/snap/example/1.0-SNAPSHOT/example-1.0-20200118.000050-1.pom.md5
[2020-01-18 01:01:03 CET]     40B snap/me/nicolaferraro/snap/example/1.0-SNAPSHOT/example-1.0-20200118.000050-1.pom.sha1
[2020-01-18 01:01:03 CET]    772B snap/me/nicolaferraro/snap/example/1.0-SNAPSHOT/maven-metadata.xml
[2020-01-18 01:01:03 CET]     32B snap/me/nicolaferraro/snap/example/1.0-SNAPSHOT/maven-metadata.xml.md5
[2020-01-18 01:01:03 CET]     40B snap/me/nicolaferraro/snap/example/1.0-SNAPSHOT/maven-metadata.xml.sha1
[2020-01-18 01:01:03 CET]    286B snap/me/nicolaferraro/snap/example/maven-metadata.xml
[2020-01-18 01:01:03 CET]     32B snap/me/nicolaferraro/snap/example/maven-metadata.xml.md5
[2020-01-18 01:01:03 CET]     40B snap/me/nicolaferraro/snap/example/maven-metadata.xml.sha1
```
