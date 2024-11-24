package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/bombsimon/http-helpers/middleware"
	"github.com/bombsimon/http-helpers/server"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

const htmlContent = `<html>
  <head>
    <meta name="go-import" content="%s/%s/%s git %s/%s/%s.git">
  </head>
</html>`

type config struct {
	httpListen  string
	packagePath string
	clonePath   string
	certFile    string
	keyFile     string
}

func envOrDefault(env string, defaultValue string) string {
	if v := os.Getenv(env); v != "" {
		return v
	}

	return defaultValue
}

func main() {
	cfg := config{}

	flag.StringVar(&cfg.httpListen, "http-listen", envOrDefault("HTTP_LISTEN", ":4080"), "The host and/or port to listen on")
	flag.StringVar(&cfg.packagePath, "package-path", envOrDefault("PACKAGE_PATH", "github.com"), "The default path for the package")
	flag.StringVar(&cfg.clonePath, "clone-path", envOrDefault("CLONE_PATH", "https://github.com"), "The default path to clone the package")
	flag.StringVar(&cfg.certFile, "cert-file", os.Getenv("CERT_FILE"), "Path to the certificate for TLS")
	flag.StringVar(&cfg.keyFile, "key-file", os.Getenv("KEY_FILE"), "Path to the key file for TLS")
	flag.Parse()

	runForever(&cfg)
}

func runForever(cfg *config) {
	logger := logrus.New().WithFields(logrus.Fields{
		"listen": cfg.httpListen,
		"pkg":    cfg.packagePath,
		"clone":  cfg.clonePath,
	})
	logger.Logger.SetFormatter(&logrus.JSONFormatter{})

	httpServer := &http.Server{
		Addr:    cfg.httpListen,
		Handler: createHandler(cfg, logger),
	}

	idleConnsClosed := server.GracefulShutdown(
		httpServer,
		10*time.Second,
		logger,
	)

	logger.Infof(
		"will create meta tag for all packages under %s and point to %s",
		cfg.packagePath,
		cfg.clonePath,
	)

	withTLS := cfg.certFile != "" && cfg.keyFile != ""

	logger.WithField("tls", withTLS).Info("server listening")

	if withTLS {
		if err := httpServer.ListenAndServeTLS(cfg.certFile, cfg.keyFile); err != nil {
			logger.Fatal(err)
		}
	} else {
		if err := httpServer.ListenAndServe(); err != nil {
			logger.Fatal(err)
		}
	}

	<-idleConnsClosed
}

func createHandler(cfg *config, logger logrus.FieldLogger) http.Handler {
	r := mux.NewRouter()
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
			w,
			htmlContent,
			cfg.packagePath,
			project,
			pkg,
			cfg.clonePath,
			project,
			pkg,
		)
	})

	return middleware.AddMiddlewares(
		r,
		middleware.PanicRecovery(logger),
		middleware.Logger(logger),
	)
}
