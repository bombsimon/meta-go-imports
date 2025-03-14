# meta-go-imports

To make go import statements work the URI specified need to resolve to a web
page with a `go-import` meta tag. This can be found e.g. at
[GitHub](https://github.com) and looks something like this:

```html
<meta
  name="go-import"
  content="github.com/user/package git https://github.com/user/package.git"
/>
```

If you want to host your internal packages and ensure that you can run `go get`
and thus use go modules you need a server responding with this tag.

This is a simple HTTP server that will generate a proper meta tag with whatever
import path you specify. You start the server by setting your host/port, package
path and where to clone.

```sh
go build -o meta-go-imports main.go
./meta-go-imports \
    --http-listen ":4080" \
    --package-path "dev.internal.se" \
    --clone-path "git+ssh://git@another.internal.se:7999" \
    --cert-file "/path/to/cert.pem" \
    --key-file "/path/to/key.pem
```

You can ommit `--cert-file` and `--key-file` if you don't want to use TLS, but
then you must add `-insecure` when running `go get`.

When you try to fetch `some-package` from `some-project` with something like
`dev.internal.se/some-project/some-package`, the server will respond with the
following:

```html
<html>
  <head>
    <meta
      name="go-import"
      content="dev.internal.se/some-project/some-package git git+ssh://git@another.internal.se:7999/some-project/some-package.git"
    />
  </head>
</html>
```

## Docker

The project comes with a `Dockerfile` which you can build an run. The file sets
the same default values as the compiled program (just for visualization) but you
can build and run it with your own values.

```sh
docker build --tag meta-go-imports .
docker run \
    -it --rm -p "4080:4080" \
    -e HTTP_LISTEN=":4080" \
    -e PACKAGE_PATH="dev.internal.se" \
    -e CLONE_PATH="git+ssh://git@another.internal.se:7999" \
    -e CERT_FILE="/certificates/cert.pem" \
    -e KEY_FILE="/certificates/key.pem" \
    -v $(pwd)/certificates:/certificates \
    bombsimon/meta-go-imports
```

### Downloading

The image is also hosted both on
[Dockerhub](https://hub.docker.com/repository/docker/bombsimon/meta-go-imports/general)
and
[GitHub](https://github.com/users/bombsimon/packages/container/package/meta-go-imports)
so you can download it without building the image yourslef.

```sh
# From Docker Hub
docker run bombsimon/meta-go-imports

# From GitHub
docker run ghcr.io/bombsimon/meta-go-imports
```
