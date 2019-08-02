// Copyright 2018 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package plugin

import (
	"context"

	"github.com/drone/drone-go/drone"
	"github.com/drone/drone-go/plugin/config"
	"golang.org/x/oauth2"

	"github.com/google/go-github/github"
)

// New returns a new configuration plugin.
func New(namespace, name, path, branch, token string) config.Plugin {
	return &plugin{
		namespace: namespace,
		name:      name,
		path:      path,
		branch:    branch,
		token:     token,
	}
}

type plugin struct {
	namespace string
	name      string
	path      string
	branch    string
	token     string
}

func (p *plugin) Find(ctx context.Context, req *config.Request) (*drone.Config, error) {
	// creates a github client used to fetch the yaml.
	trans := oauth2.NewClient(ctx, oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: p.token},
	))
	client := github.NewClient(trans)

	// HACK: the drone-go library does not currently work
	// with 0.9 which means the configuration file path is
	// always empty. default to .drone.yml. This can be
	// removed as soon as drone-go is fully updated for 0.9.
	path := req.Repo.Config
	if path == "" {
		path = ".drone.yml"
	}

	// get the configuration file from the github
	// repository for the build ref.
	data, _, _, err := client.Repositories.GetContents(ctx, req.Repo.Namespace, req.Repo.Name, path, &github.RepositoryContentGetOptions{Ref: req.Build.After})
	if err == nil && data != nil {
		// get the file contents.
		content, err := data.GetContent()
		if err != nil {
			return nil, err
		}
		return &drone.Config{
			Data: content,
		}, nil
	}

	// if the configuration file does not exist,
	// we should fallback to a global configuration
	// file stored in a central repository.
	data, _, _, err = client.Repositories.GetContents(ctx, p.namespace, p.name, p.path, &github.RepositoryContentGetOptions{Ref: p.branch})
	if err != nil {
		return nil, err
	}
	// get the file contents.
	content, err := data.GetContent()
	if err != nil {
		return nil, err
	}
	return &drone.Config{
		Data: content,
	}, nil
}
