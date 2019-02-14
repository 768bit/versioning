package verman

import (
	"encoding/json"
	"fmt"
	"errors"
	"github.com/768bit/vutils"
	"strings"
	"time"
)

type VersionData struct {
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
}

type PkgVersionRevisions map[string]int

func NewVersionData() *VersionData {

	vd := VersionData{
		Major:         0,
		Minor:         0,
		Revision:      0,
		BuildRevision: 0,
		PkgRevisions:  PkgVersionRevisions{},
	}

	vd.NewBuild()

	return &vd
}

func LoadVersionData() (*VersionData, error) {

	return getVersionDataFile()
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

func (vd *VersionData) Save() error {

	if _, err := vd.GetGitCommit(); err != nil {
		fmt.Println("Error getting git commit hash for versioning")
	}

	return saveVersionDataFile(vd)

}

func (vd *VersionData) FullVersionString() string {

	return fmt.Sprintf("%d.%d.%d", vd.Major, vd.Minor, vd.Revision)

}

func (vd *VersionData) GetGitCommit() (string, error) {

	comm, err := vutils.Exec.ExecCommandShowStdErrReturnOutput("git", "rev-parse", "--short", "HEAD")

	if err != nil {
		return "", err
	}

	vd.GitCommit = strings.TrimSpace(comm)

	return vd.GitCommit, nil

}

func (vd *VersionData) TagGitCommit(message string) (error) {

  _, err := vutils.Exec.ExecCommandShowStdErrReturnOutput("git", "tag", "-a", "-m", message, "-f", fmt.Sprintf("V%s", vd.FullVersionString()), "HEAD")
  
  return err

}

func (vd *VersionData) ToJSON() (string, error) {

	bContent, err := json.Marshal(vd)

	if err != nil {
		return "", err
	}

	return string(bContent), nil

}
