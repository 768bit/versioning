package support

import (
	"fmt"
	"gitlab.768bit.com/pub/vpkg"
)

var VDATA *vpkg.VersionData

func GetVersionData() *vpkg.VersionData {

	return VDATA

}

func InitialiseVersionData(namespace string, target string, root string) error {

	if target != "" && namespace != "" {
		vd, err := vpkg.LoadVersionDataNamespace(root, namespace)
		if err != nil {
			fmt.Println(err)
			vd = vpkg.NewVersionDataNamespace(root, namespace, target)
			err = vd.Save(root)
			if err != nil {
				return err
			}
			fmt.Println("Created a new version.json file for namespace and target.")
			return InitialiseVersionData(root, namespace, target)
		}
		VDATA = vd
		return nil
	}

	vd, err := vpkg.LoadVersionData(root)
	if err != nil {
		fmt.Println(err)
		vd = vpkg.NewVersionData()
		err = vd.Save(root)
		if err != nil {
			return err
		}
		fmt.Println("Created a new version.json file.")
		return InitialiseVersionData(root, "", "")
	}
	VDATA = vd
	return nil

}
