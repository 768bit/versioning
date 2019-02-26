package xgoutils

import (
	"errors"
	"fmt"
	"github.com/768bit/vutils"
	"github.com/bmatcuk/doublestar"
	"gitlab.768bit.com/pub/vpkg/common"
	"gitlab.768bit.com/pub/vpkg/pkgutils"
	"log"
	"os"
	"path/filepath"
	"regexp"
)

var QueueModeDisablesDirectBuildsError = errors.New("Cannot perform a direct build if aqueue mode has been enabled.")
var QueueDataAlreadySetError = errors.New("Cannot set queue data as it has already been set.")
var QueueNotReadyError = errors.New("Cannot add a build to the queue as it hasn't been initialised. Run SetQueueData.")
var QueueCurrentlyRunningError = errors.New("Cannot proceed with request as the queue is currently being processed.")
var BuildCurrentlyRunningError = errors.New("Cannot proceed with request as a build is currently running.")

func (cs *XGOCompileSettings) SetQueueData(pkgMeta *pkgutils.PackageMetadata, vdata *common.VersionData, buildDir string, platform ...XGOPlatform) error {

	if cs.isQueue {
		return QueueDataAlreadySetError
	}
	cs.isQueue = true

	if platform == nil || len(platform) == 0 {
		cs.buildPlatformsForQueue = []XGOPlatform{ANDROID, DARWIN, IOS, LINUX, WINDOWS}
	} else {
		cs.buildPlatformsForQueue = platform
	}

	cs.vdata = vdata
	cs.packageMetadata = pkgMeta

	cs.buildPathForQueue = buildDir

	cs._packagesQueue = map[XGOPlatform]packagingFunc{}

	return nil

}

func (cs *XGOCompileSettings) AddBuildToQueue(name string, targetFile string, ldflags string) error {

	if !cs.isQueue {
		return QueueNotReadyError
	}

	if cs.isRunning {
		return QueueCurrentlyRunningError
	}

	if cs.compileQueue == nil {
		cs.compileQueue = []compileTarget{{
			name:       name,
			targetFile: targetFile,
			ldflags:    ldflags,
		}}
	} else {
		cs.compileQueue = append(cs.compileQueue, compileTarget{
			name:       name,
			targetFile: targetFile,
			ldflags:    ldflags,
		})
	}
	return nil

}

func (cs *XGOCompileSettings) BuildAllPlatforms(pkgMeta *pkgutils.PackageMetadata, vdata *common.VersionData, buildDir string, targetFile string, ldflags string) ([]XGOCompileResult, error) {

	return cs.BuildPlatforms(pkgMeta, vdata, buildDir, ldflags, ANDROID, DARWIN, IOS, LINUX, WINDOWS)

}

func makePlatformMissingError(platform XGOPlatform) error {
	return errors.New(fmt.Sprintf("The selected platform %s isnt configured for XGO compilation.", platform))
}

func (cs *XGOCompileSettings) BuildPlatforms(pkgMeta *pkgutils.PackageMetadata, vdata *common.VersionData, buildDir string, targetFile string, ldflags string, platform ...XGOPlatform) ([]XGOCompileResult, error) {

	if cs.isRunning {
		if cs.isQueue {
			return nil, QueueCurrentlyRunningError
		}
		return nil, BuildCurrentlyRunningError
	}

	if cs.isQueue {
		return nil, QueueModeDisablesDirectBuildsError
	}

	cs.isRunning = true
	defer cs.resetRunning()

	cs.packageMetadata = pkgMeta
	cs.vdata = vdata

	compileResults := []XGOCompileResult{}

	for _, platformItem := range platform {

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

		cr := XGOCompileResult{
			Platform:      platformItem,
			Architectures: archList,
		}

		platformBuildPath := filepath.Join(buildDir, string(platformItem))
		os.RemoveAll(platformBuildPath)
		vutils.Files.CreateDirIfNotExist(platformBuildPath)

		cr.Output, cr.Error = cs.runPlatformBuild(platformBuildPath, targetFile, ldflags, platformItem, platformTargets)
		if cr.Error == nil {

			//get binary list..

			binTargets, err := doublestar.Glob(filepath.Join(platformBuildPath, "xgo-build-*"))
			if err == nil && len(binTargets) > 0 {

				cleanBinTargets := map[XGOArchitecture]string{}

				for _, binTarget := range binTargets {

					//get the last item from path so we can check what the situation is...

					getPlatformAndArchForBinary(binTarget)

				}

				pkgBuildPath := filepath.Join(platformBuildPath, "pkg")
				if ipkg, err := cs.getPackagingInterface(platformItem); err != nil {
					return nil, err
				} else if pkgRes, err := ipkg.BuildPackages(pkgMeta, vdata, pkgBuildPath); err != nil {
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

	return compileResults, nil

}

func (cs *XGOCompileSettings) runPlatformBuild(platformBuildPath string, targetFile string, ldflags string, platform XGOPlatform, platformTargets string) (string, error) {

	cmd := vutils.Exec.CreateAsyncCommand("xgo", false, platformTargets, "-out", "xgo-build", "-v", "-ldflags", ldflags, "-dest", platformBuildPath+"/", targetFile)
	err := cmd.CaptureStdoutAndStdErr(true, true).CopyEnv().StartAndWait()
	if err != nil {
		return "", err
	}

	return string(cmd.GetStdoutBuffer()), nil

}

func (cs *XGOCompileSettings) getPackagingInterface(platformItem XGOPlatform) (IPackaging, error) {

	switch platformItem {
	case ANDROID:
		if cs.Android == nil {
			return nil, makePlatformMissingError(platformItem)
		}
		return cs.Android, nil
	case DARWIN:
		if cs.Darwin == nil {
			return nil, makePlatformMissingError(platformItem)
		}
		return cs.Darwin, nil
	case IOS:
		if cs.IOS == nil {
			return nil, makePlatformMissingError(platformItem)
		}
		return cs.IOS, nil
	case LINUX:
		if cs.Linux == nil {
			return nil, makePlatformMissingError(platformItem)
		}
		return cs.Linux, nil
	case WINDOWS:
		if cs.Windows == nil {
			return nil, makePlatformMissingError(platformItem)
		}
		return cs.Windows, nil
	default:
		return nil, errors.New(fmt.Sprintf("The selected platform %s is invalid.", platformItem))
	}

}

var BinaryRX = regexp.MustCompile("^xgo-build-([a-z]+)(?:-\\d+(?:\\.\\d+)?)?-([a-zA-Z0-9\\-]+)(?:\\.exe)?$")
var BinaryNameInvalidError = errors.New("The binary name supplied isnt in the correct format")
var BinaryArchitectureNameInvalidError = errors.New("Thearchitecture used for the binary isnt valid for this platform.")

func getPlatformAndArchForBinary(binaryPath string) (XGOPlatform, XGOArchitecture, error) {

	dir, file := filepath.Split(binaryPath)

	if !BinaryRX.MatchString(file) {
		return "", "", BinaryNameInvalidError
	}

	matches := BinaryRX.FindStringSubmatch(file)

	if len(matches) != 3 {
		return "", "", BinaryNameInvalidError
	}

	return checkPlatformAndArchitecture(matches[1], matches[2])

}

func checkPlatformAndArchitecture(plat, arch string) (oplatform XGOPlatform, oarch XGOArchitecture, err error) {

	err = nil
	oplatform = XGOPlatform(plat)
	oarch = XGOArchitecture(arch)
	var archList []XGOArchitecture
	switch oplatform {
	case ANDROID:
		archList = ANDROID_ALLOWED_ARCHITECTURES
		break
	case DARWIN:
		archList = DARWIN_ALLOWED_ARCHITECTURES
		break
	case IOS:
		archList = IOS_ALLOWED_ARCHITECTURES
		break
	case LINUX:
		archList = LINUX_ALLOWED_ARCHITECTURES
		break
	case WINDOWS:
		archList = WINDOWS_ALLOWED_ARCHITECTURES
		break
	default:
		err = makePlatformMissingError(oplatform)
		return
	}

	if archList == nil || len(archList) == 0 {
		err = BinaryArchitectureNameInvalidError
		return
	}

	if archInSlice(oarch, archList) {
		return
	}
	err = BinaryArchitectureNameInvalidError
	return

}
