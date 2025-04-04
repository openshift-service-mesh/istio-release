package main

import (
	"fmt"
	"os"
	"path"

	"istio.io/istio/pkg/log"
	"istio.io/release-builder/pkg"
	"istio.io/release-builder/pkg/build"
	"istio.io/release-builder/pkg/model"
	"istio.io/release-builder/pkg/util"
)

// This program imports some code from the upstream release-builder
// to build release-ready helm chart tarballs. The only argument it
// takes is a version name, which must match a directory in the
// repository containing a manifest.yaml.
//
// Note that this DOES NOT PRODUCE A FULL RELEASE. It only packages
// the helm charts.

func main() {
	if len(os.Args) < 2 {
		fmt.Println("You need to pass a directory name")
		os.Exit(1)
	}
	dirName := os.Args[1]
	manifestFilename := path.Join(dirName, "manifest.yaml")
	if !util.FileExists(manifestFilename) {
		fmt.Println(manifestFilename + " not found")
		os.Exit(1)
	}

	inManifest, err := pkg.ReadInManifest(manifestFilename)
	if err != nil {
		panic(fmt.Errorf("failed to unmarshal manifest: %v", err))
	}

	manifest, err := pkg.InputManifestToManifest(inManifest)
	if err != nil {
		panic(fmt.Errorf("failed to setup manifest: %v", err))
	}

	// Save these values as they are needed for git commits and PRs
	savedIstioGit := inManifest.Dependencies.Get()["istio"].Git
	savedIstioBranch := inManifest.Dependencies.Get()["istio"].Branch
	log.Infof("Saved Istio git:\n%+v", savedIstioGit)
	log.Infof("Saved Istio branch:\n%+v", savedIstioBranch)

	if err := pkg.SetupWorkDir(manifest.Directory); err != nil {
		panic(fmt.Errorf("failed to setup work dir: %v", err))
	}
	if err := pkg.Sources(manifest); err != nil {
		panic(fmt.Errorf("failed to fetch sources: %v", err))
	}
	log.Infof("Fetched all sources and setup working directory at %v", manifest.WorkDir())

	if err := build.SanitizeAllCharts(manifest); err != nil {
		panic(fmt.Errorf("failed to sanitize charts: %v", err))
	}
	if !util.IsValidSemver(manifest.Version) {
		panic("Invalid Semantic Version. Skipping Charts build")
	}
	if _, f := manifest.BuildOutputs[model.Helm]; f {
		if err := build.HelmCharts(manifest); err != nil {
			panic(fmt.Errorf("failed to build HelmCharts: %v", err))
		}
	}
	util.CopyDir(path.Join(manifest.OutDir(), "helm"), dirName)
}
