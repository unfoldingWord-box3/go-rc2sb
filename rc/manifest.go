// Package rc provides types and parsing for Resource Container (RC) manifest files.
package rc

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// Manifest represents the top-level structure of an RC manifest.yaml file.
type Manifest struct {
	DublinCore DublinCore `yaml:"dublin_core"`
	Checking   Checking   `yaml:"checking"`
	Projects   []Project  `yaml:"projects"`
}

// DublinCore holds the dublin_core metadata from the RC manifest.
type DublinCore struct {
	ConformsTo  string   `yaml:"conformsto"`
	Contributor []string `yaml:"contributor"`
	Creator     string   `yaml:"creator"`
	Description string   `yaml:"description"`
	Format      string   `yaml:"format"`
	Identifier  string   `yaml:"identifier"`
	Issued      string   `yaml:"issued"`
	Language    Language  `yaml:"language"`
	Modified    string   `yaml:"modified"`
	Publisher   string   `yaml:"publisher"`
	Relation    []string `yaml:"relation"`
	Rights      string   `yaml:"rights"`
	Source      []Source `yaml:"source"`
	Subject     string   `yaml:"subject"`
	Title       string   `yaml:"title"`
	Type        string   `yaml:"type"`
	Version     string   `yaml:"version"`
}

// Language describes the language in the RC manifest.
type Language struct {
	Direction  string `yaml:"direction"`
	Identifier string `yaml:"identifier"`
	Title      string `yaml:"title"`
}

// Source describes a source reference in the RC manifest.
type Source struct {
	Identifier string `yaml:"identifier"`
	Language   string `yaml:"language"`
	Version    string `yaml:"version"`
}

// Checking holds the checking metadata from the RC manifest.
type Checking struct {
	CheckingEntity []string `yaml:"checking_entity"`
	CheckingLevel  string   `yaml:"checking_level"`
}

// Project describes a single project entry in the RC manifest.
type Project struct {
	Categories    []string `yaml:"categories"`
	Identifier    string   `yaml:"identifier"`
	Path          string   `yaml:"path"`
	Sort          int      `yaml:"sort"`
	Title         string   `yaml:"title"`
	Versification string   `yaml:"versification"`
}

// LoadManifest reads and parses a manifest.yaml file from the given directory.
func LoadManifest(dir string) (*Manifest, error) {
	path := filepath.Join(dir, "manifest.yaml")
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("not a valid Resource Container: manifest.yaml not found in %s", dir)
		}
		return nil, fmt.Errorf("reading manifest.yaml: %w", err)
	}

	var m Manifest
	if err := yaml.Unmarshal(data, &m); err != nil {
		return nil, fmt.Errorf("parsing manifest.yaml: %w", err)
	}

	return &m, nil
}
