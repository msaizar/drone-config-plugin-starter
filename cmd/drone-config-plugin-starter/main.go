// Copyright 2018 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"net/http"

	"github.com/drone/drone-config-plugin-starter/plugin"
	"github.com/drone/drone-go/plugin/config"

	_ "github.com/joho/godotenv/autoload"
	"github.com/kelseyhightower/envconfig"
	"github.com/sirupsen/logrus"
)

type spec struct {
	Debug     bool   `envconfig:"PLUGIN_DEBUG"`
	Address   string `envconfig:"PLUGIN_ADDRESS" default:":3000"`
	Secret    string `envconfig:"PLUGIN_SECRET"`
	Token     string `envconfig:"GITHUB_TOKEN"`
	Namespace string `envconfig:"GITHUB_REPO_OWNER"`
	Name      string `envconfig:"GITHUB_REPO_NAME"`
	Branch    string `envconfig:"GITHUB_REPO_BRANCH" default:"master"`
	Path      string `envconfig:"GITHUB_YAML_PATH" default:".drone.yml"`
}

func main() {
	spec := new(spec)
	err := envconfig.Process("", spec)
	if err != nil {
		logrus.Fatal(err)
	}

	if spec.Debug {
		logrus.SetLevel(logrus.DebugLevel)
	}
	if spec.Secret == "" {
		logrus.Fatalln("missing secret key")
	}
	if spec.Token == "" {
		logrus.Warnln("missing github token")
	}
	if spec.Namespace == "" {
		logrus.Warnln("missing github repository owner")
	}
	if spec.Name == "" {
		logrus.Warnln("missing github repository name")
	}
	if spec.Address == "" {
		spec.Address = ":3000"
	}

	handler := config.Handler(
		plugin.New(
			spec.Namespace,
			spec.Name,
			spec.Path,
			spec.Branch,
			spec.Token,
		),
		spec.Secret,
		logrus.StandardLogger(),
	)

	logrus.Infof("server listening on address %s", spec.Address)

	http.Handle("/", handler)
	logrus.Fatal(http.ListenAndServe(spec.Address, nil))
}
