// +build mage

package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/768bit/isokit"
	"github.com/768bit/vpkg/pkgutils"
	"github.com/768bit/vpkg/xgoutils"
	"github.com/768bit/vutils"
	"github.com/bmatcuk/doublestar"
	"github.com/magefile/mage/mg" // mg contains helpful utility functions, like Deps
)

var VDATA *vpkg.VersionData

func InitialiseVersionData() error {

	cwd, _ := os.Getwd()

	vd, err := vpkg.LoadVersionData(cwd)
	if err != nil {
		fmt.Println(err)
		vd = vpkg.NewVersionData()
		err = vd.Save(cwd)
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
	VDATA.Save(cwd)
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

	pkrBuildCmd := vutils.Exec.CreateAsyncCommand("packr2", false)
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
	VDATA.Save(cwd)
	InitialiseVersionData()

	debPkgRoot := filepath.Join(cwd, "build", "pkg")
	allowed, vermanPackageName := vpkg.PkgOsVersionString("ubuntu", VDATA)
	if !allowed {
		return errors.New("unable to make version")
	}

	pkgDetails := &pkgutils.PackageDetails{
		Name:         "verman",
		Architecture: "amd64",
		Maintainer:   "Craig Smith <craig.smith@768bit.com>",
		Version:      VDATA.FullVersionString(),
		Description:  "verman command line version.json management utility.",
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

var XGOBuildSettings = xgoutils.XGOCompileSettings{
	Android: xgoutils.NewAndroidCompileSettings(xgoutils.ARCH_ALL),
	Darwin:  xgoutils.NewDarwinCompileSettings(xgoutils.ARCH_ALL),
	IOS:     xgoutils.NewIosCompileSettings(xgoutils.ARCH_ALL),
	Linux: xgoutils.NewLinuxCompileSettings(xgoutils.ARCH_ALL).
		AddPackagingOptions(xgoutils.NewLinuxPackagingOptions().AddDebian(xgoutils.NewDebianLinuxPackagingOptions("/usr/bin", nil))),
	Windows: xgoutils.NewWindowsCompileSettings(xgoutils.ARCH_ALL),
}

func BuildXGO() error {

	return doXGOBuild(true)

}

func doXGOBuild(isRelease bool) error {

	mg.Deps(InitialiseVersionData, PackAssets)

	VDATA.NewBuild()
	VDATA.Save()
	cwd, _ := os.Getwd()

	releaseVer := "verman"
	vcliOut := "build"

	if !isRelease {
		releaseVer += "-dev"
	} else {
		releaseVer += "-" + VDATA.FullVersionString()
	}

	fmt.Println("Creating Cross Platform Build: " + releaseVer)

	fmt.Println("Building Verman using XGO for cross compilation...")
	ldflags := fmt.Sprintf(LD_FLAGS_FMT_STR, VDATA.FullVersionString(), VDATA.ShortID, VDATA.DateString, VDATA.UUID, VDATA.GitCommit)
	fmt.Printf("Building with LDFLAGS: %s\n", ldflags)
	cmd := vutils.Exec.CreateAsyncCommand("xgo", false, "-out", releaseVer, "--targets=darwin/amd64,linux/*,windows/*", "-v", "-ldflags", ldflags, "-dest", vcliOut+"/", ".")
	err = cmd.BindToStdoutAndStdErr().SetWorkingDir(cwd).CopyEnv().StartAndWait()
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil

}
