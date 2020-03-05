# Snap

Snap is a simple tool that allows publishing SNAPSHOT libraries into a Kubernetes
cluster, so that they can be used by on-cluster build tools.

The backend that stores the artifacts is based on [Minio](https://min.io). 

## Usage

Download a version of Snap from the [release page](https://github.com/container-tools/snap/releases)
and put it into your OS path.

To deploy the example project:

```
./snap example/library
```

Should be successful.

## To Verify

On OpenShift, expose the Minio service as route:

```
oc expose service/snap-minio-service
```

Then get the Minio endpoint from the Route:

```
export MINIO_ENDPOINT=$(oc get route snap-minio-service -o jsonpath='{.spec.host}')
```

On vanilla Kubernetes, it depends on your installation, but you need to expose the service as well.

Get the username and password from the secret:

```
kubectl get secret snap-minio-credentials -o yaml
# to convert each entry
# echo --type-the-base64-data-here-- | base64 -d
```

Download the `mc` tool from Minio, then:

```
mc config host add minio http://$MINIO_ENDPOINT minio-access-key-got-before minio-access-secret-got-before
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
