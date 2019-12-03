# meta-go-imports

To make go import statements work the URI specified need to resolve to a web
page with a `go-import` meta tag. This can be found e.g. at
[GitHub](https://github.com) and looks something like this:

```html
<meta name="go-import" content="github.com/user/package git https://github.com/user/package.git">
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
    --package-path "dev.internal.com" \
    --clone-path  "git@another.internal.se:7999"
```

When you try to fetch `some-package` from `some-project` with something like
`dev.internal.se/some-project/some-package`, the server will respond with the
following:

```html
<html>
  <head>
    <meta name="go-import" content="dev.internal.com/some-project/some-package git git@another.internal.se:7999:some-project/some-package.git">
  </head>
</html>
```
