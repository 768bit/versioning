package common

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/768bit/vutils"
	"io/ioutil"
	"path/filepath"
	"strings"
)

const VERSION_DIR_PATH = ".verman"

func getVersionDataFilePath(cwd string, check bool) (string, error) {

	versionFilePath := filepath.Join(cwd, "version.json")

	if check && !vutils.Files.CheckPathExists(versionFilePath) {
		return "", errors.New(fmt.Sprintf("Unable to find version.json file in %s", cwd))
	}

	return versionFilePath, nil

}

func getVersionDataFilePathForNamespace(cwd string, namespace string, check bool) (string, error) {

	versionFilePath := filepath.Join(cwd, VERSION_DIR_PATH, namespace+".json")

	if check && !vutils.Files.CheckPathExists(versionFilePath) {
		return "", errors.New(fmt.Sprintf("Unable to find version.json file in %s", cwd))
	}

	return versionFilePath, nil

}

func getVersionDataFileForNamespace(cwd string, namespace string) (*VersionData, error) {

	versionFilePath, err := getVersionDataFilePathForNamespace(cwd, namespace, true)
	if err != nil {
		return nil, err
	}

	bContent, err := ioutil.ReadFile(versionFilePath)

	if err != nil {
		return nil, err
	}

	//marshall the versionData file into VersionData

	var vdata VersionData

	err = json.Unmarshal(bContent, &vdata)

	if err != nil {
		return nil, err
	} else {
		vdata.GetGitCommit()
	}

	return &vdata, nil

}

func getVersionDataFile(cwd string) (*VersionData, error) {

	versionFilePath, err := getVersionDataFilePath(cwd, true)
	if err != nil {
		return nil, err
	}

	bContent, err := ioutil.ReadFile(versionFilePath)

	if err != nil {
		return nil, err
	}

	//marshall the versionData file into VersionData

	var vdata VersionData

	err = json.Unmarshal(bContent, &vdata)

	if err != nil {
		return nil, err
	} else {
		vdata.GetGitCommit()
	}

	return &vdata, nil

}

func saveVersionDataFile(cwd string, vdata *VersionData) error {

	versionFilePath, err := getVersionDataFilePath(cwd, false)
	if err != nil {
		return err
	}

	bContent, err := json.Marshal(vdata)

	if err != nil {
		return err
	}

	err = ioutil.WriteFile(versionFilePath, bContent, 0640)

	if err != nil {
		return err
	}

	return nil

}

func saveVersionDataFileForNamespace(cwd string, namespace string, vdata *VersionData) error {

	versionFilePath, err := getVersionDataFilePathForNamespace(cwd, namespace, false)
	if err != nil {
		return err
	}

	bContent, err := json.Marshal(vdata)

	if err != nil {
		return err
	}

	nsVermanPath := filepath.Join(cwd, VERSION_DIR_PATH)

	if !vutils.Files.PathExists(nsVermanPath) {
		vutils.Files.CreateDirIfNotExist(nsVermanPath)
	}

	err = ioutil.WriteFile(versionFilePath, bContent, 0640)

	if err != nil {
		return err
	}

	return nil

}

func CheckAllowedOs(os string) (bool, string) {

	os = strings.ToLower(os)

	switch os {

	case "ubuntu":
		fallthrough
	case "debian":
		fallthrough
	case "redhat":
		fallthrough
	case "centos":
		fallthrough
	case "fedora":
		fallthrough
	case "suse":
		fallthrough
	case "arch":
		fallthrough
	case "alpine":
		fallthrough
	case "gentoo":
		fallthrough
	case "mandrake":
		fallthrough
	case "slackware":
		fallthrough
	case "macos":
		fallthrough
	case "freebsd":
		fallthrough
	case "openbsd":
		fallthrough
	case "netbsd":
		fallthrough
	case "windows":
		return true, os
	default:
		return false, ""

	}

}

func PkgOsVersionString(os string, vdata *VersionData) (bool, string) {

	switch os {

	case "ubuntu":
		fallthrough
	case "debian":
		return true, fmt.Sprintf("%d.%d.%d+%s-%d", vdata.Major, vdata.Minor, vdata.Revision, vdata.ShortID, vdata.PkgRevisions[os])
	case "redhat":
		fallthrough
	case "centos":
		fallthrough
	case "fedora":
		return true, fmt.Sprintf("%d.%d.%d-%d", vdata.Major, vdata.Minor, vdata.Revision, vdata.PkgRevisions[os])
	case "suse":
		fallthrough
	case "arch":
		fallthrough
	case "alpine":
		fallthrough
	case "gentoo":
		fallthrough
	case "mandrake":
		fallthrough
	case "slackware":
		fallthrough
	case "macos":
		fallthrough
	case "freebsd":
		fallthrough
	case "openbsd":
		fallthrough
	case "netbsd":
		fallthrough
	case "windows":
		return true, fmt.Sprintf("%d.%d.%d", vdata.Major, vdata.Minor, vdata.Revision)
	default:
		return false, ""

	}

}

func resetRevisions(vdata *VersionData) {

	if vdata == nil || vdata.PkgRevisions == nil {
		return
	}

	for key, _ := range vdata.PkgRevisions {
		vdata.PkgRevisions[key] = 0
	}

}
