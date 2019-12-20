package xgoutils

import (
	"errors"
	"fmt"
	"github.com/768bit/vpkg/common"
	"github.com/768bit/vpkg/pkgutils"
	"strings"
)

type XGOPlatform = string

const (
	ANDROID XGOPlatform = "android"
	DARWIN  XGOPlatform = "darwin"
	IOS     XGOPlatform = "ios"
	LINUX   XGOPlatform = "linux"
	WINDOWS XGOPlatform = "windows"
)

type XGOArchitecture = string

const (
	ARCH_386       XGOArchitecture = "386"
	ARCH_AMD64     XGOArchitecture = "amd64"
	ARCH_ARM_5     XGOArchitecture = "arm-5"
	ARCH_ARM_6     XGOArchitecture = "arm-6"
	ARCH_ARM_7     XGOArchitecture = "arm-7"
	ARCH_ARM64     XGOArchitecture = "arm64"
	ARCH_MIPS      XGOArchitecture = "mips"
	ARCH_MIPS_LE   XGOArchitecture = "mipsle"
	ARCH_MIPS64    XGOArchitecture = "mips64"
	ARCH_MIPS64_LE XGOArchitecture = "mips64le"
)

var ARCH_ALL []XGOArchitecture = []XGOArchitecture{"*"}

type compileTarget struct {
	name       string
	targetFile string
	ldflags    string
}

type XGOCompileSettings struct {
	Android                *XGOAndroidCompileSettings
	Darwin                 *XGODarwinCompileSettings
	IOS                    *XGOIosCompileSettings
	Linux                  *XGOLinuxCompileSettings
	Windows                *XGOWindowsCompileSettings
	compileQueue           []compileTarget
	isQueue                bool
	isRunning              bool
	buildPathForQueue      string
	buildPlatformsForQueue []XGOPlatform
	packageMetadata        *pkgutils.PackageMetadata
	vdata                  *common.VersionData
	_packagesQueue         map[XGOPlatform]packagingFunc
}

type PackagingFuncMap map[XGOArchitecture]map[PackageType]func(arch XGOArchitecture) (string, error)

func (pfm PackagingFuncMap) AddPackageBuildFuncToArch(arch XGOArchitecture, pkgType PackageType, buildFn func(arch XGOArchitecture) (string, error)) {

	if existing, ok := pfm[arch]; !ok || existing == nil {
		pfm[arch] = map[PackageType]func(arch XGOArchitecture) (string, error){
			pkgType: buildFn,
		}
	} else if existingPkg, ok := pfm[arch][pkgType]; !ok || existingPkg == nil {
		pfm[arch][pkgType] = buildFn
	} else {
		pfm[arch][pkgType] = buildFn
	}

}

type packagingFunc = func() (PackageBuildMap, error)

type BaseXGOPlatformCompileSettings struct {
	Architectures   []XGOArchitecture
	Platform        XGOPlatform
	compileSettings *XGOCompileSettings
}

func (base *BaseXGOPlatformCompileSettings) addArchitectures(archs []XGOArchitecture, allowed []XGOArchitecture) {

	for _, arch := range archs {
		if arch == "*" {
			base.Architectures = allowed
			return
		} else if !archInSlice(arch, base.Architectures) && archInSlice(arch, allowed) {
			base.Architectures = append(base.Architectures, arch)
		}
	}

}

func newBaseXGOPlatformCompileSettings(platform XGOPlatform) BaseXGOPlatformCompileSettings {
	return BaseXGOPlatformCompileSettings{
		Architectures: []XGOArchitecture{},
		Platform:      platform,
	}
}

func archInSlice(a XGOArchitecture, list []XGOArchitecture) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func (base *BaseXGOPlatformCompileSettings) GetXGOCompileTarget() string {

	targetsList := make([]string, len(base.Architectures))

	for i, arch := range base.Architectures {
		targetsList[i] = base.Platform + "/" + arch
	}

	return fmt.Sprintf("--targets=%s", strings.Join(targetsList, ","))

}

func (cs *XGOCompileSettings) addPackageProcessToQueue(platform XGOPlatform, fn packagingFunc, replace bool) {

	if existing, ok := cs._packagesQueue[platform]; !ok || existing == nil {
		cs._packagesQueue[platform] = fn
	} else if replace {
		cs._packagesQueue[platform] = fn
	}

}

func (cs *XGOCompileSettings) processQueue() (PlatformPackageBuildMap, error) {

	pbm := PlatformPackageBuildMap{}

	for p, fn := range cs._packagesQueue {
		if bm, err := fn(); err != nil {
			return nil, err
		} else {
			pbm[p] = bm
		}
	}
	return pbm, nil

}

var NoArchitecturesToBuild = errors.New("There are no architectures to build for this platform.")
var NoPackagesToBuildError = errors.New("There are no package packaging settings for this platform.")
var HostIncompatibleBuildPackageError = errors.New("There are packaging settings for this platform but the compilation host is unable to build packages.")

type IPackaging interface {
	BuildPackages(packageMeta *pkgutils.PackageMetadata, vdata *common.VersionData, pkgBuildRoot string, binaryTargets BinaryTargetsArchitectureMap) (PackageBuildMap, error)
}

type PackageType = string

const (
	TarPackage        PackageType = "tar.gz"
	ZipPackage        PackageType = "zip"
	SourceCodePackage PackageType = "src.tar.gz"
	BinaryOnlyPackage PackageType = "bin"
)

type XGOCompileResult struct {
	Platform      XGOPlatform
	Architectures []XGOArchitecture
	Output        string
	Error         error
	Packages      PackageBuildMap
}

type PackageBuildMap map[XGOArchitecture]map[PackageType]string
type PlatformPackageBuildMap map[XGOPlatform]PackageBuildMap

func (base *BaseXGOPlatformCompileSettings) BuildPackages(packageMeta *pkgutils.PackageMetadata, vdata *common.VersionData, pkgBuildRoot string, binaryTargets BinaryTargetsArchitectureMap) (PackageBuildMap, error) {

	return nil, NoPackagesToBuildError

}

func (base *BaseXGOPlatformCompileSettings) ProcessPackageQueue() (PackageBuildMap, error) {

	return nil, NoPackagesToBuildError

}

type BinaryTargetsArchitectureMap = map[XGOArchitecture]string

func (pbm PackageBuildMap) AddBuiltPackage(arch XGOArchitecture, packageType PackageType, packagePath string) {

	if existing, ok := pbm[arch]; !ok || existing == nil {
		pbm[arch] = map[PackageType]string{
			packageType: packagePath,
		}
	} else {
		pbm[arch][packageType] = packagePath
	}

}
