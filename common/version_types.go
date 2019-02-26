package common

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/768bit/vutils"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type VersionData struct {
	Name          string              `json:"name,omitempty"`
	Major         int                 `json:"major"`
	Minor         int                 `json:"minor"`
	Revision      int                 `json:"revision"`
	BuildRevision int                 `json:"buildRevision"`
	UUID          string              `json:"uuid"`
	DateString    string              `json:"dateString"`
	Date          time.Time           `json:"date"`
	ShortID       string              `json:"shortID"`
	GitCommit     string              `json:"gitCommit"`
	PkgRevisions  PkgVersionRevisions `json:"revisions"`
	TargetFolder  string              `json:"targetFolder,omitempty"`
}

type PkgVersionRevisions map[string]int

func NewVersionData() *VersionData {

	cwd, _ := os.Getwd()

	vd := VersionData{
		Major:         0,
		Minor:         0,
		Revision:      0,
		BuildRevision: 0,
		PkgRevisions:  PkgVersionRevisions{},
		TargetFolder:  cwd,
	}

	vd.NewBuild()

	return &vd
}

func NewVersionDataNamespace(cwd string, namespace string, targetFolder string) *VersionData {

	//nneed to ensure target folder is a child folder (relative) of CWD

	if targetFolder[0] == '/' {
		targetFolder = strings.Replace(targetFolder, cwd, "", 1)
		if targetFolder[0] == '/' {
			targetFolder = targetFolder[1:]
		}
	} else if targetFolder[0:1] == "./" {
		targetFolder = targetFolder[2:]
	}

	fullpath := filepath.Join(cwd, targetFolder)
	if !vutils.Files.PathExists(fullpath) {
		log.Println("Cannot find path", fullpath)
		panic(errors.New("Cannot find target folder for namespaced version.json operations"))
	}

	vd := VersionData{
		TargetFolder:  targetFolder,
		Name:          namespace,
		Major:         0,
		Minor:         0,
		Revision:      0,
		BuildRevision: 0,
		PkgRevisions:  PkgVersionRevisions{},
	}

	vd.NewBuild()

	return &vd
}

func LoadVersionData(cwd string) (*VersionData, error) {

	return getVersionDataFile(cwd)
}

func LoadVersionDataNamespace(cwd string, namespace string) (*VersionData, error) {

	return getVersionDataFileForNamespace(cwd, namespace)
}

func (vd *VersionData) NewBuild() {

	vd.NewBuildUUID()
	vd.NewBuildDate()
	vd.NewBuildRevision()

}

func (vd *VersionData) NewBuildUUID() {

	//generate a new build ID that will be used an put it in the version data file...

	verUUID, _ := vutils.UUID.MakeUUIDString()
	vd.UUID = verUUID

}

func (vd *VersionData) NewBuildDate() {

	//generate a new build ID that will be used an put it in the version data file...

	vd.Date = time.Now()

	vd.DateString = vd.Date.Format(time.UnixDate)

	vd.NewBuildShortID()

}

func (vd *VersionData) NewBuildShortID() {

	//generate a new build ID that will be used an put it in the version data file...

	vd.ShortID = vd.Date.Format("20060102")

}

func (vd *VersionData) NewBuildRevision() {

	//generate a new build ID that will be used an put it in the version data file...

	vd.BuildRevision += 1

}

func (vd *VersionData) NewMajorVersion() {

	vd.Major = vd.Major + 1
	vd.Minor = 0
	vd.Revision = 0
	vd.BuildRevision = 0
	resetRevisions(vd)

}

func (vd *VersionData) NewMinorVersion() {

	vd.Minor = vd.Minor + 1
	vd.Revision = 0
	vd.BuildRevision = 0
	resetRevisions(vd)

}

func (vd *VersionData) NewRevision() {

	vd.Revision = vd.Revision + 1
	vd.BuildRevision = 0
	resetRevisions(vd)

}

func (vd *VersionData) NewPkgRevision(os string) {

	if _, ok := vd.PkgRevisions[os]; !ok {
		vd.PkgRevisions[os] = 0
	}

	vd.PkgRevisions[os] = vd.PkgRevisions[os] + 1

}

func (vd *VersionData) GetPkgRevision(os string) (int, error) {

	if _, ok := vd.PkgRevisions[os]; !ok {
		return 0, errors.New(fmt.Sprintf("No revisions have been made for os type %s", os))
	}

	return vd.PkgRevisions[os], nil

}

func (vd *VersionData) Save(cwd string) error {

	if _, err := vd.GetGitCommit(); err != nil {
		fmt.Println("Error getting git commit hash for versioning")
	}

	if vd.TargetFolder != "" && vd.Name != "" {
		return saveVersionDataFileForNamespace(cwd, vd.Name, vd)
	}

	return saveVersionDataFile(cwd, vd)

}

func (vd *VersionData) FullVersionString() string {

	return fmt.Sprintf("%d.%d.%d", vd.Major, vd.Minor, vd.Revision)

}

func (vd *VersionData) GetGitCommit() (string, error) {

	cwd, _ := os.Getwd()

	gitCommitCmd := vutils.Exec.CreateAsyncCommand("git", false, "rev-parse", "--short", "HEAD").CopyEnv().CaptureStdoutAndStdErr()

	if vd.TargetFolder != "" && vd.Name != "" {
		gitCommitCmd = gitCommitCmd.SetWorkingDir(filepath.Join(cwd, vd.TargetFolder))
	}

	err := gitCommitCmd.StartAndWait()
	if err != nil {
		return "", err
	}

	out := gitCommitCmd.GetStdoutBuffer()

	vd.GitCommit = strings.TrimSpace(string(out))

	return vd.GitCommit, nil

}

func (vd *VersionData) TagGitCommit(message string) error {

	cwd, _ := os.Getwd()

	gitCommitCmd := vutils.Exec.CreateAsyncCommand("git", false, "tag", "-a", "-m", message, "-f", fmt.Sprintf("V%s", vd.FullVersionString()), "HEAD").CopyEnv()

	if vd.TargetFolder != "" && vd.Name != "" {
		gitCommitCmd = gitCommitCmd.SetWorkingDir(filepath.Join(cwd, vd.TargetFolder))
	}

	err := gitCommitCmd.StartAndWait()
	if err != nil {
		return err
	}

	return nil

}

func (vd *VersionData) PerformGitCommit(message string) error {

	cwd, _ := os.Getwd()

	gitCommitCmd := vutils.Exec.CreateAsyncCommand("git", false, "commit", "-a", "-m", message).CopyEnv()

	if vd.TargetFolder != "" && vd.Name != "" {
		gitCommitCmd = gitCommitCmd.SetWorkingDir(filepath.Join(cwd, vd.TargetFolder))
	}

	err := gitCommitCmd.StartAndWait()
	if err != nil {
		return err
	}

	return nil

}

func (vd *VersionData) ToJSON() (string, error) {

	bContent, err := json.Marshal(vd)

	if err != nil {
		return "", err
	}

	return string(bContent), nil

}
