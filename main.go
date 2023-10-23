package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"gopkg.in/yaml.v2"
)

type Channel struct {
	Releases []Release `json:"releases"`
}

type Release struct {
	Version string `json:"version"`
}

// This program is used to generate Google Kubernetes Engine JSON versions for Renovate custom datasource.
// Custom datasource docs: https://docs.renovatebot.com/modules/datasource/custom/
func main() {
	ctx := context.Background()
	requestedChannel := flag.String("channel", "stable", "Channel to scrape")
	requestedLocation := flag.String("location", "us-central-1c", "GCP location to check versions for (they might differ per location)")
	flag.Parse()

	channel, err := scrapeChannel(ctx, *requestedChannel, *requestedLocation)
	if err != nil {
		fmt.Println("Error scraping channel:", err)
		os.Exit(1)
	}

	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", " ")
	if err := encoder.Encode(channel); err != nil {
		fmt.Println("Error encoding JSON:", err)
		os.Exit(1)
	}
}

func scrapeChannel(ctx context.Context, channel string, location string) (Channel, error) {
	cmd := exec.CommandContext(ctx, "gcloud", "container", "get-server-config",
		"--location", location,
		"--flatten", "channels",
		"--filter", "channels.channel="+channel,
		"--format", "yaml(channels.channel,channels.validVersions)",
	)
	cmdOut, err := cmd.CombinedOutput()
	if err != nil {
		return Channel{}, fmt.Errorf("running gcloud command: %w, OUTPUT:\n %s", err, string(cmdOut))
	}

	releases, err := extractReleases(string(cmdOut))
	if err != nil {
		return Channel{}, fmt.Errorf("extracting releases: %w", err)
	}
	return Channel{Releases: releases}, nil
}

func extractReleases(cmdOut string) ([]Release, error) {
	parts := strings.SplitN(cmdOut, "---", 2)
	if len(parts) != 2 {
		return nil, fmt.Errorf("unexpected command output: %s", cmdOut)
	}
	channelsYAML := parts[1]

	// Sample command output:
	// Fetching server config for us-central1-c
	// ---
	// channels:
	//  channel: RAPID
	//  validVersions:
	//  - 1.28.2-gke.1157000
	//  - 1.27.6-gke.1506000
	//  - 1.27.6-gke.1445000
	//  - 1.27.6-gke.1248000
	//  - 1.27.5-gke.200
	//  - 1.27.4-gke.900
	//  - 1.26.9-gke.1507000
	//  - 1.26.9-gke.1437000
	//  - 1.26.8-gke.200
	//  - 1.26.7-gke.500
	//  - 1.25.14-gke.1474000
	//  - 1.25.14-gke.1421000
	//  - 1.25.13-gke.200
	//  - 1.25.12-gke.500
	//  - 1.24.17-gke.2155000
	//  - 1.24.17-gke.2113000
	//  - 1.24.17-gke.200`
	var decoded struct {
		Channels struct {
			Channel       string   `yaml:"channel"`
			ValidVersions []string `yaml:"validVersions"`
		} `yaml:"channels"`
	}

	if err := yaml.Unmarshal([]byte(channelsYAML), &decoded); err != nil {
		return nil, fmt.Errorf("unmarshaling YAML: %w", err)
	}

	var releases []Release
	for _, v := range decoded.Channels.ValidVersions {
		// We want to remove the "-gke.*" suffix from the version.
		// Example: 1.24.17-gke.200 -> 1.24.17
		v = strings.Split(v, "-")[0]
		releases = append(releases, Release{Version: v})
	}
	return releases, nil
}
