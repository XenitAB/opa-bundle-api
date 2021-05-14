package main

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"os"

	"github.com/xenitab/opa-bundle-api/pkg/bundle"
	"github.com/xenitab/opa-bundle-api/pkg/config"
	"github.com/xenitab/opa-bundle-api/pkg/rule"
)

var (
	// Version is set at build time to print the released version using --version
	Version = "v0.0.0-dev"
	// Revision is set at build time to print the release git commit sha using --version
	Revision = ""
	// Created is set at build time to print the timestamp for when it was built using --version
	Created = ""
)

func main() {
	cfg, err := newConfigClient()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to generate config: %q\n", err)
		os.Exit(1)
	}

	err = start(cfg)
	if err != nil {
		os.Exit(1)
	}

	os.Exit(0)
}

func start(cfg config.Client) error {
	rules := rule.NewRepository()
	for i := 1; i <= 5; i++ {
		rand, err := generateRandomString(10)
		if err != nil {
			return err
		}

		_, err = rules.Add(rule.Options{
			Country:  fmt.Sprintf("Sweden-%d", i),
			City:     fmt.Sprintf("Gothenburg-%d", i),
			Building: fmt.Sprintf("HQ-%d", i),
			Role:     fmt.Sprintf("admin-%s", rand),
			Action:   rule.ActionAllow,
		})

		if err != nil {
			return err
		}
	}

	pol := bundle.Policies{
		{
			Name:    "test-name",
			Content: "test-content",
		},
	}

	bundle, err := bundle.GenerateBundle(&rules, pol)
	if err != nil {
		return err
	}

	fmt.Println(bundle)

	return nil
}

func newConfigClient() (config.Client, error) {
	opts := config.Options{
		Version:  Version,
		Revision: Revision,
		Created:  Created,
	}

	return config.NewClient(opts)
}

func generateRandomString(n int) (string, error) {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	ret := make([]byte, n)
	for i := 0; i < n; i++ {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(letters))))
		if err != nil {
			return "", err
		}
		ret[i] = letters[num.Int64()]
	}

	return string(ret), nil
}