package pkgutils

import (
	"errors"
	"github.com/768bit/isokit"
	"path/filepath"
)

type PackageMetadata struct {
	Name        string
	Maintainer  string
	Description string
	parent      *PackageMetadata
}

func (pmeta *PackageMetadata) MakePackageDetails(version string, architecture string) *PackageDetails {
	return &PackageDetails{
		Name:         pmeta.Name,
		Version:      version,
		Maintainer:   pmeta.Maintainer,
		Architecture: architecture,
		Description:  pmeta.Description,
	}
}

func (pmeta *PackageMetadata) MakeChild(name string) *PackageMetadata {
	return &PackageMetadata{
		Name:        name,
		Maintainer:  pmeta.Maintainer,
		Description: pmeta.Description,
		parent:      pmeta,
	}
}

func (pmeta *PackageMetadata) IsChild() bool {
	return pmeta.parent != nil
}

func (pmeta *PackageMetadata) GetParent() *PackageMetadata {
	return pmeta.parent
}

type PackageDetails struct {
	Name         string
	Version      string
	Maintainer   string
	Architecture string
	Description  string
}

var templateSet *isokit.TemplateSet

func getTemplateSet() (*isokit.TemplateSet, error) {

	if templateSet == nil {
		templateSet = isokit.NewTemplateSet()
		baseBundlePath := filepath.Join("templates", "template.bundle")
		if contents, err := AssetsBox.Find(baseBundlePath); err != nil || len(contents) == 0 {
			return nil, errors.New("Unable to load templates binary from Assets Box :: " + baseBundlePath)
		} else {
			if err := templateSet.RestoreTemplateBundleFromBinary(contents); err != nil {
				return nil, err
			}
		}
	}

	return templateSet, nil

}
