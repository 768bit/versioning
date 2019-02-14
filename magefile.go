// +build mage

package main

import (
  "fmt"
  "github.com/768bit/verman"
  "github.com/768bit/verman/pkgutils"
  "github.com/768bit/vutils"
  "github.com/768bit/isokit"
  "github.com/bmatcuk/doublestar"
  "github.com/magefile/mage/mg" // mg contains helpful utility functions, like Deps
  "os"
  "errors"
  "path/filepath"
)

var VDATA *verman.VersionData

func InitialiseVersionData() error {

	vd, err := verman.LoadVersionData()
	if err != nil {
		fmt.Println(err)
		vd = verman.NewVersionData()
		err = vd.Save()
		if err != nil {
			return err
		}
		fmt.Println("Created a new version.json file.")
		return InitialiseVersionData()
	}
	VDATA = vd
	return nil

}

// Default target to run when none is specified
// If not set, running mage will list available targets
var Default = Build

var LD_FLAGS_FMT_STR = "-w -s -X main.Version=%s -X main.Build=%s -X \"main.BuildDate=%s\" -X main.BuildUUID=%s -X main.GitCommit=%s"

// A build step that requires additional params, or platform specific steps for example
func Build() error {
	cwd, _ := os.Getwd()

	mg.Deps(InitialiseVersionData, PackAssets)
	VDATA.NewBuild()
	VDATA.Save()
  InitialiseVersionData()
  vermanCliFolder := filepath.Join(cwd, "cli")
	vermanBinary := filepath.Join(cwd, "build", "verman")
	vutils.Files.CreateDirIfNotExist(filepath.Join(cwd, "build"))

	fmt.Println("Building verman...")
	ldflags := fmt.Sprintf(LD_FLAGS_FMT_STR, VDATA.FullVersionString(), VDATA.ShortID, VDATA.DateString, VDATA.UUID, VDATA.GitCommit)
	fmt.Printf("Building with LDFLAGS: %s\n", ldflags)
	cmd := vutils.Exec.CreateAsyncCommand("go", false, "build", "-v", "-ldflags", ldflags, "-o", vermanBinary)
	err := cmd.BindToStdoutAndStdErr().CopyEnv().SetWorkingDir(vermanCliFolder).StartAndWait()
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

func PackAssets() error {
  cwd, _ := os.Getwd()

  mg.Deps(BuildTemplates)
  fmt.Println("Packing Assets...")

  pkrBuildCmd := vutils.Exec.CreateAsyncCommand("packr", false, "-z")
  err := pkrBuildCmd.BindToStdoutAndStdErr().SetWorkingDir(filepath.Join(cwd, "pkgutils")).StartAndWait()
  if err != nil {
    fmt.Println(err)
    return err
  }
  return nil
}

func BuildDeb() error {
  mg.Deps(Build)
  cwd, _ := os.Getwd()

  VDATA.NewPkgRevision("ubuntu")
  VDATA.Save()
  InitialiseVersionData()

  debPkgRoot := filepath.Join(cwd, "build", "pkg")
  allowed, vermanPackageName := verman.PkgOsVersionString("ubuntu", VDATA)
  if !allowed {
    return errors.New("unable to make version")
  }

  pkgDetails := &pkgutils.PackageDetails{
    Name: "verman",
    Architecture: "amd64",
    Maintainer: "Craig Smith <craig.smith@768bit.com>",
    Version: VDATA.FullVersionString(),
    Description: "verman command line version.json management utility.",
  }

  cm := vutils.NewContentsMap(true)
  binSrc := filepath.Join(cwd, "build", "verman")
  err := cm.AddFile(binSrc, "/usr/bin/verman", 0755)
  if err != nil {
    return err
  }

  err = pkgutils.Debian.NewDebianPackage(debPkgRoot, pkgDetails, cm)
  if err != nil {
    return err
  }

  targets, err := doublestar.Glob(filepath.Join(cwd, "build", fmt.Sprintf("verman_*_%s.deb", pkgDetails.Architecture)))
  if err == nil && len(targets) > 0 {
    for _, target := range targets {
      fmt.Println("Removing Package File", target)
      os.Remove(target)
    }
  } else if err != nil {
    fmt.Println(err)
  }

  vermanPackageName = fmt.Sprintf("verman_%s_%s.deb", vermanPackageName, pkgDetails.Architecture)
  vermanPackageBuildPath := filepath.Join(cwd, "build", vermanPackageName)

  return pkgutils.Debian.BuildDebianPackage(debPkgRoot, pkgDetails, vermanPackageBuildPath)

}


func BuildTemplates() error {
  cwd, _ := os.Getwd()
  fmt.Println("Building templates for verman assets box...")

  templatesPath := filepath.Join(cwd, "templates")

  templateBundleOutPath := filepath.Join(cwd, "assets", "templates", "template.bundle")

  os.RemoveAll(templateBundleOutPath)

  ts := isokit.NewTemplateSet()
  ts.GatherTemplatesFromPath("verman", templatesPath)

  vutils.Files.CreateDirIfNotExist(filepath.Join(cwd, "assets", "templates"))
  fmt.Println(ts)

  return ts.PersistTemplateBundleToDisk(templateBundleOutPath)

}
