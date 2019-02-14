package pkgutils

import (
	"errors"
	"github.com/768bit/isokit"
	"path/filepath"
)

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
