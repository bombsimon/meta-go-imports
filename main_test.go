package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/net/html"
)

func TestHandler(t *testing.T) {
	logger := logrus.New()
	logger.SetLevel(logrus.FatalLevel)

	for _, tc := range []struct {
		description          string
		path                 string
		cfg                  *config
		expectedResponsecode int
		expectedContent      string
	}{
		{
			description:          "Invalid path",
			path:                 "http://localhost/invalid",
			cfg:                  &config{},
			expectedResponsecode: http.StatusNotFound,
		},
		{
			description: "HTTP clone path",
			path:        "http://localhost/my/repo",
			cfg: &config{
				packagePath: "github.com",
				clonePath:   "https://github.com",
			},
			expectedResponsecode: http.StatusOK,
			expectedContent:      "github.com/my/repo git https://github.com/my/repo.git",
		},
		{
			description: "SSH clone path",
			path:        "http://localhost/my-org/ssh-repo",
			cfg: &config{
				packagePath: "dev.myorg.net",
				clonePath:   "git+ssh://git@internal.server.myorg.net:7999",
			},
			expectedResponsecode: http.StatusOK,
			expectedContent:      "dev.myorg.net/my-org/ssh-repo git git+ssh://git@internal.server.myorg.net:7999/my-org/ssh-repo.git",
		},
	} {
		t.Run(tc.description, func(t *testing.T) {
			w := httptest.NewRecorder()
			handler := createHandler(tc.cfg, logger)
			req := httptest.NewRequest("GET", tc.path, nil)
			handler.ServeHTTP(w, req)

			assert.Equal(t, tc.expectedResponsecode, w.Code)
			if w.Code != http.StatusOK {
				return
			}

			doc, err := html.Parse(w.Body)
			require.NoError(t, err)

			//              html       head       meta
			metaNode := doc.FirstChild.FirstChild.FirstChild.NextSibling
			require.Len(t, metaNode.Attr, 2)

			nameNode := metaNode.Attr[0]
			contentNode := metaNode.Attr[1]

			assert.Equal(t, "go-import", nameNode.Val)
			assert.Equal(t, tc.expectedContent, contentNode.Val)
		})
	}
}
