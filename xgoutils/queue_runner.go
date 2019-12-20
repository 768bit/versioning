package xgoutils

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/768bit/vutils"
	"github.com/bmatcuk/doublestar"
)

func (cs *XGOCompileSettings) resetRunning() {
	cs.isRunning = false
}

func (cs *XGOCompileSettings) BuildQueue() ([]XGOCompileResult, error) {

	if cs.isRunning {
		if cs.isQueue {
			return nil, QueueCurrentlyRunningError
		}
		return nil, BuildCurrentlyRunningError
	}

	if !cs.isQueue {
		return nil, QueueNotReadyError
	}

	cs.isRunning = true
	defer cs.resetRunning()

	compileResults := []XGOCompileResult{}

	for _, platformItem := range cs.buildPlatformsForQueue {

		var archList []XGOArchitecture
		var platformTargets string

		switch platformItem {
		case ANDROID:
			if cs.Android == nil {
				return nil, makePlatformMissingError(platformItem)
			}
			cs.Android.compileSettings = cs
			archList = cs.Android.Architectures
			platformTargets = cs.Android.GetXGOCompileTarget()
			break
		case DARWIN:
			if cs.Darwin == nil {
				return nil, makePlatformMissingError(platformItem)
			}
			cs.Darwin.compileSettings = cs
			archList = cs.Darwin.Architectures
			platformTargets = cs.Darwin.GetXGOCompileTarget()
			break
		case IOS:
			if cs.IOS == nil {
				return nil, makePlatformMissingError(platformItem)
			}
			cs.IOS.compileSettings = cs
			archList = cs.IOS.Architectures
			platformTargets = cs.IOS.GetXGOCompileTarget()
			break
		case LINUX:
			if cs.Linux == nil {
				return nil, makePlatformMissingError(platformItem)
			}
			cs.Linux.compileSettings = cs
			archList = cs.Linux.Architectures
			platformTargets = cs.Linux.GetXGOCompileTarget()
			break
		case WINDOWS:
			if cs.Windows == nil {
				return nil, makePlatformMissingError(platformItem)
			}
			cs.Windows.compileSettings = cs
			archList = cs.Windows.Architectures
			platformTargets = cs.Windows.GetXGOCompileTarget()
			break
		default:
			return nil, errors.New(fmt.Sprintf("The selected platform %s is invalid.", platformItem))
		}

		//so now we run the compile...

		if archList == nil || len(archList) == 0 {
			return nil, NoArchitecturesToBuild
		}

		platformBuildPath := filepath.Join(cs.buildPathForQueue, string(platformItem))
		os.RemoveAll(platformBuildPath)
		vutils.Files.CreateDirIfNotExist(platformBuildPath)

		for _, ctarget := range cs.compileQueue {

			cr := XGOCompileResult{
				Platform:      platformItem,
				Architectures: archList,
			}

			targetBuildPath := filepath.Join(platformBuildPath, string(ctarget.name))
			os.RemoveAll(targetBuildPath)
			vutils.Files.CreateDirIfNotExist(targetBuildPath)

			cr.Output, cr.Error = cs.runPlatformBuild(targetBuildPath, ctarget.targetFile, ctarget.ldflags, platformItem, platformTargets)
			if cr.Error == nil {

				//get binary list..

				targetPkgMeta := cs.packageMetadata.MakeChild(ctarget.name)

				binTargets, err := doublestar.Glob(filepath.Join(targetBuildPath, "xgo-build-*"))
				if err == nil && len(binTargets) > 0 {

					cleanBinTargets := map[XGOArchitecture]string{}

					for _, binTarget := range binTargets {

						//get the last item from path so we can check what the situation is...

						p, a, err := getPlatformAndArchForBinary(binTarget)
						if err != nil {
							log.Println("Error getting platform and architecture from found binary target", binTarget, "ERR:", err)
						} else if p != platformItem {
							log.Println("The platform found from binary target doesnt match the current platform", binTarget, "Extracted:", p, "Needed:", platformItem)
						} else if !archInSlice(a, archList) {
							log.Println("The architecture found from binary target isnt in the allowed list of architectures for this platform", binTarget, "Extracted:", a, "Allowed:", strings.Join(archList, ", "))
						} else {
							cleanBinTargets[a] = binTarget
						}

					}

					pkgBuildPath := filepath.Join(platformBuildPath, "pkg")
					if ipkg, err := cs.getPackagingInterface(platformItem); err != nil {
						return nil, err
					} else if pkgRes, err := ipkg.BuildPackages(targetPkgMeta, cs.vdata, pkgBuildPath, cleanBinTargets); err != nil {
						if err == NoPackagesToBuildError {
							log.Println("No packages to build for platform", platformItem)
						} else {
							return nil, err
						}
					} else {
						cr.Packages = pkgRes
					}

				}

			}

			compileResults = append(compileResults, cr)

		}

	}

	return compileResults, nil

}
