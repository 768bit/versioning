package pkgutils

import (
	"fmt"
	"github.com/768bit/vutils"
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

	return contentsMap.DoCopy(pkgBuildPath)

}

func (deb *debianUtils) BuildDebianPackage(pkgBuildDir string, packageDetails *PackageDetails, buildOutputPath string) error {

	//clean and create directory structure for package...

	pkgBuildPath := filepath.Join(pkgBuildDir, "deb", packageDetails.Architecture)

	//dpkg-deb --build pkg/debian build/vcli-$(call GetOSPackageVersion,debian).deb
	debPkgBuildCmd := vutils.Exec.CreateAsyncCommand("dpkg-deb", false, "--build", pkgBuildPath, buildOutputPath)
	err := debPkgBuildCmd.BindToStdoutAndStdErr().StartAndWait()
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil

}

var Debian = &debianUtils{}
