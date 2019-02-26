package pkgutils

import (
	"fmt"
	"github.com/768bit/vutils"
	"github.com/bmatcuk/doublestar"
	"io/ioutil"
	"os"
	"path/filepath"
)

type debianUtils struct{}

func (deb *debianUtils) NewDebianPackage(pkgBuildDir string, packageDetails *PackageDetails, contentsMap *vutils.ContentsMap) error {

	//clean and create directory structure for package...

	pkgBuildPath := filepath.Join(pkgBuildDir, "deb", packageDetails.Architecture)

	err := os.RemoveAll(pkgBuildPath)
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	err = vutils.Files.CreateDirIfNotExist(pkgBuildPath)
	if err != nil {
		return err
	}

	debControlPath := filepath.Join(pkgBuildPath, "DEBIAN")

	err = vutils.Files.CreateDirIfNotExist(debControlPath)
	if err != nil {
		return err
	}

	debControlPath = filepath.Join(debControlPath, "control")

	err = os.RemoveAll(debControlPath)
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	ts, err := getTemplateSet()
	if err != nil {
		return err
	}

	if contents, err := ts.RenderSimple("verman/debian-control", packageDetails); err != nil {
		return err
	} else if err := ioutil.WriteFile(debControlPath, contents, 0644); err != nil {
		return err
	}

	//now we can copy the assets...

	if contentsMap == nil {
		return nil
	}

	return contentsMap.DoCopy(pkgBuildPath)

}

func (deb *debianUtils) WriteContentsMap(pkgBuildDir string, arch string, contentsMap *vutils.ContentsMap) error {

	//clean and create directory structure for package...

	pkgBuildPath := filepath.Join(pkgBuildDir, "deb", arch)

	return contentsMap.DoCopy(pkgBuildPath)

}

func (deb *debianUtils) BuildDebianPackage(pkgBuildDir string, arch string, buildOutputPath string) error {

	//clean and create directory structure for package...

	pkgBuildPath := filepath.Join(pkgBuildDir, "deb", arch)

	//dpkg-deb --build pkg/debian build/vcli-$(call GetOSPackageVersion,debian).deb
	debPkgBuildCmd := vutils.Exec.CreateAsyncCommand("dpkg-deb", false, "--build", pkgBuildPath, buildOutputPath)
	err := debPkgBuildCmd.CaptureStdoutAndStdErr(true, true).StartAndWait()
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil

}

func (deb *debianUtils) CleanPreviousPackages(name string, buildOutputRoot string, arch string) {

	//clean and create directory structure for package...

	targets, err := doublestar.Glob(filepath.Join(buildOutputRoot, fmt.Sprintf("%s_*_%s.deb", name, arch)))
	if err == nil && len(targets) > 0 {
		for _, target := range targets {
			fmt.Println("Removing Package File", target)
			os.Remove(target)
		}
	} else if err != nil {
		fmt.Println(err)
	}

}

var Debian = &debianUtils{}
