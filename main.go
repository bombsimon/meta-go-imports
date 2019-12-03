package main

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/alecthomas/kingpin"
	"github.com/bombsimon/http-helpers/middleware"
	"github.com/bombsimon/http-helpers/server"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

const (
	projectSeparatorSlash = "/"
	projectSeparatorColon = ":"
)

const (
	htmlContent = `<html>
  <head>
    <meta name="go-import" content="%s/%s/%s git %s%s%s/%s.git">
  </head>
</html>`
)

func main() {
	var (
		httpListen  = kingpin.Flag("http-listen", "The host and/or port to listen on").Default(":4080").Envar("HTTP_LISTEN").String()
		packagePath = kingpin.Flag("package-path", "The default path for the package").Default("github.com").Envar("PACKAGE_PATH").String()
		clonePath   = kingpin.Flag("clone-path", "The default path to clone the package").Default("https://github.com").Envar("CLONE_PATH").String()
	)

	kingpin.Parse()

	var (
		r      = mux.NewRouter()
		logger = logrus.New().WithFields(logrus.Fields{
			"listen": *httpListen,
			"pkg":    *packagePath,
			"clone":  *clonePath,
		})
		separator = projectSeparatorSlash
	)

	if !strings.HasPrefix(*clonePath, "http") {
		separator = projectSeparatorColon
	}

	r.HandleFunc("/{project}/{pkg}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		project, ok := vars["project"]
		if !ok {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "could not find project in path")

			return
		}

		pkg, ok := vars["pkg"]
		if !ok {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "could not find package in path")

			return
		}

		fmt.Fprintf(
			w, htmlContent,
			*packagePath,
			project,
			pkg,
			*clonePath,
			separator,
			project,
			pkg,
		)
	})

	handlers := middleware.AddMiddlewares(
		r,
		middleware.PanicRecovery(logger),
		middleware.Logger(logger),
	)

	s := &http.Server{
		Addr:    *httpListen,
		Handler: handlers,
	}

	idleConnsClosed := server.GracefulShutdown(
		s,              // The HTTP server
		10*time.Second, // Wait time
		logrus.New(),   // Optional logger
	)

	logger.Infof("will create meta tag for all packages under %s and point to %s", *packagePath, *clonePath)
	logger.Infof("listening on '%s'\n", *httpListen)

	if err := s.ListenAndServe(); err != nil {
		logger.Fatal(err)
	}

	<-idleConnsClosed
}
