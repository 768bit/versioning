package xgoutils

import (
	"fmt"
	"github.com/768bit/vpkg/common"
	"github.com/768bit/vpkg/pkgutils"
	"github.com/768bit/vutils"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strconv"
)

var LINUX_ALLOWED_ARCHITECTURES = []XGOArchitecture{ARCH_386, ARCH_AMD64, ARCH_ARM_5, ARCH_ARM_6, ARCH_ARM_7, ARCH_ARM64, ARCH_MIPS, ARCH_MIPS64, ARCH_MIPS_LE, ARCH_MIPS64_LE}

type XGOLinuxCompileSettings struct {
	BaseXGOPlatformCompileSettings
	PackagingOptions *LinuxPackagingOptions
	_packageQueue    PackagingFuncMap
}

func (linux *XGOLinuxCompileSettings) AddPackagingOptions(opts *LinuxPackagingOptions) *XGOLinuxCompileSettings {
	opts.platformCompileSettings = linux
	linux.PackagingOptions = opts
	return linux
}

func NewLinuxPackagingOptions() *LinuxPackagingOptions {
	return &LinuxPackagingOptions{}
}

type LinuxPackagingOptions struct {
	platformCompileSettings *XGOLinuxCompileSettings
	Debian                  *DebianLinuxPackagingOptions
}

func (linpkg *LinuxPackagingOptions) AddDebian(debianOptions *DebianLinuxPackagingOptions) *LinuxPackagingOptions {
	debianOptions.platformCompileSettings = linpkg.platformCompileSettings
	linpkg.Debian = debianOptions
	return linpkg
}

type DebianLinuxPackagingOptions struct {
	platformCompileSettings *XGOLinuxCompileSettings
	BinPath                 string
	contentMap              map[string]string
	archContentsMap         map[XGOArchitecture]*vutils.ContentsMap
	pkgBuildRoot            string
}

var DestFileRX = regexp.MustCompile("^([0-7]{4}:)?(?=/)(/(?=[^/\\0])[^/\\0]+)+/?$")

func (debpkg *DebianLinuxPackagingOptions) runBuild(arch XGOArchitecture) (string, error) {

	//get the contentmap so it can be written...
	err := pkgutils.Debian.WriteContentsMap(debpkg.pkgBuildRoot, arch, debpkg.archContentsMap[arch])
	if err != nil {
		return "", err
	}

	cwd, _ := os.Getwd()

	debpkg.platformCompileSettings.compileSettings.vdata.NewPkgRevision("debian")
	debpkg.platformCompileSettings.compileSettings.vdata.Save(cwd)

	buildOutputRoot := filepath.Join(cwd, "build", "pkg", LINUX, DebianPackage)

	_, pkgName := common.PkgOsVersionString("debian", debpkg.platformCompileSettings.compileSettings.vdata)
	pkgName = fmt.Sprintf("%s_%s_%s.deb", debpkg.platformCompileSettings.compileSettings.packageMetadata.Name, pkgName, arch)

	buildOutputFullPath := filepath.Join(buildOutputRoot, pkgName)

	pkgutils.Debian.CleanPreviousPackages(debpkg.platformCompileSettings.compileSettings.packageMetadata.Name, buildOutputRoot, arch)

	err := pkgutils.Debian.BuildDebianPackage(debpkg.pkgBuildRoot, arch, buildOutputFullPath)
	return buildOutputFullPath, err

}

func (debpkg *DebianLinuxPackagingOptions) createContentMapForArch(arch XGOArchitecture, binaryMap map[string]string, pkgBuildRoot string) (*vutils.ContentsMap, error) {

	debpkg.pkgBuildRoot = pkgBuildRoot

	if currContentsMap, ok := debpkg.archContentsMap[arch]; !ok || currContentsMap == nil {

		vContentsMap := vutils.NewContentsMap(true)

		if debpkg.contentMap != nil {

			for source, dest := range debpkg.contentMap {

				if !DestFileRX.MatchString(dest) {
					log.Println("Unable to add content map entry", source, "->", dest, ":: It hasn't been formatted correctly")
				} else if !vutils.Files.PathExists(source) {
					log.Println("Unable to add content map entry", source, "->", dest, ":: The source doesn't exist.")
				} else {

					sourceStat, err := os.Lstat(source)
					if err != nil {
						vContentsMap = nil
						return nil, err
					}

					matches := DestFileRX.FindStringSubmatch(dest)

					matchesLen := len(matches)

					destPath := ""
					destPathPerms := sourceStat.Mode()

					if matchesLen <= 1 {
						log.Println("Unable to add content map entry", source, "->", dest, ":: It hasn't been formatted correctly")
					} else if matchesLen == 2 {
						//only a path item!
						destPath = matches[1]
					} else if matchesLen == 3 {
						destPath = matches[2]
						vval, err := strconv.ParseUint(matches[1][0:3], 10, 32)
						if err != nil {
							log.Println("Unable to get file mode from path element.", matches[1][0:3], "for path", dest, "Error:", err, "Using the source items permissions instead:", destPathPerms)
						} else {
							destPathPerms = os.FileMode(vval)
						}
					}

					if sourceStat.IsDir() {
						err := vContentsMap.AddDirectory(source, destPath, destPathPerms, true)
						if err != nil {
							vContentsMap = nil
							return nil, err
						}
					} else if sourceStat.Mode()&os.ModeSymlink != 0 {
						//resolve the link so we can figure out what to do..
						resolvedSourcePath, resolvedSourceInfo, err := vutils.Files.ResolveLink(source)
						if err != nil {
							vContentsMap = nil
							return nil, err
						}
						if resolvedSourceInfo.IsDir() {
							err := vContentsMap.AddDirectory(resolvedSourcePath, destPath, destPathPerms, true)
							if err != nil {
								vContentsMap = nil
								return nil, err
							}
						} else {
							err := vContentsMap.AddFile(resolvedSourcePath, destPath, destPathPerms)
							if err != nil {
								vContentsMap = nil
								return nil, err
							}
						}
					} else {
						err := vContentsMap.AddFile(source, destPath, destPathPerms)
						if err != nil {
							vContentsMap = nil
							return nil, err
						}
					}

				}

			}

		}

		//now we need to process the binaries...

		if binaryMap != nil {

			for binSource, binaryName := range binaryMap {

				destPath := filepath.Join(debpkg.BinPath, binaryName)

				if !vutils.Files.PathExists(binSource) {
					log.Println("Unable to add binary entry", binSource, "->", destPath, ":: The source binary doesn't exist.")
				} else {
					err := vContentsMap.AddFile(binSource, destPath, 0755)
					if err != nil {
						vContentsMap = nil
						return nil, err
					}
				}

			}

		}

		debpkg.archContentsMap[arch] = vContentsMap

		return debpkg.archContentsMap[arch], nil

	} else {

		//map already exists...

		if binaryMap != nil {

			for binSource, binaryName := range binaryMap {

				destPath := filepath.Join(debpkg.BinPath, binaryName)

				if !vutils.Files.PathExists(binSource) {
					log.Println("Unable to add binary entry", binSource, "->", destPath, ":: The source binary doesn't exist.")
				} else {
					err := debpkg.archContentsMap[arch].AddFile(binSource, destPath, 0755)
					if err != nil {
						debpkg.archContentsMap[arch] = nil
						return nil, err
					}
				}

			}

		}

		return debpkg.archContentsMap[arch], nil

	}

}

func NewDebianLinuxPackagingOptions(binPath string, contentMap map[string]string) *DebianLinuxPackagingOptions {
	return &DebianLinuxPackagingOptions{
		BinPath:         binPath,
		contentMap:      contentMap,
		archContentsMap: map[XGOArchitecture]*vutils.ContentsMap{},
	}
}

func NewLinuxCompileSettings(architectures []XGOArchitecture) *XGOLinuxCompileSettings {

	lcs := &XGOLinuxCompileSettings{
		BaseXGOPlatformCompileSettings: newBaseXGOPlatformCompileSettings(LINUX),
	}

	lcs._packageQueue = PackagingFuncMap{}
	lcs.addArchitectures(architectures, LINUX_ALLOWED_ARCHITECTURES)

	return lcs

}

const (
	DebianPackage PackageType = "deb"
)

func (linux *XGOLinuxCompileSettings) BuildPackages(packageMeta *pkgutils.PackageMetadata, vdata *common.VersionData, pkgBuildRoot string, binaryTargets BinaryTargetsArchitectureMap) (PackageBuildMap, error) {

	//if the xgo compile went ok on linux we can no package (if there are settings available and the host we are compiling on supports it)

	if linux.PackagingOptions == nil {
		return nil, NoPackagesToBuildError
	} else {
		if runtime.GOOS != "linux" {
			return nil, HostIncompatibleBuildPackageError
		}

		pkgBuildMap := PackageBuildMap{}

		if !packageMeta.IsChild() {
			os.RemoveAll(pkgBuildRoot)
		}
		vutils.Files.CreateDirIfNotExist(pkgBuildRoot)

		for _, arch := range linux.Architectures {

			//now we can start to see if anythign needs to be done....

			if binary, ok := binaryTargets[arch]; !ok || binary == "" {
				log.Println("No binary available for architecture", arch)
				continue
			} else if pkgTypeMap, err := linux.buildPackage(arch, packageMeta, vdata, pkgBuildRoot, binary); err != nil {
				return nil, err
			} else if !packageMeta.IsChild() && pkgTypeMap != nil {
				pkgBuildMap[arch] = pkgTypeMap
			}

		}

		return pkgBuildMap, nil

	}

}

func (linux *XGOLinuxCompileSettings) buildPackage(arch XGOArchitecture, packageMeta *pkgutils.PackageMetadata, vdata *common.VersionData, pkgBuildRoot string, binaryTarget string) (map[PackageType]string, error) {

	//if the xgo compile went ok on linux we can no package (if there are settings available and the host we are compiling on supports it)

	if linux.PackagingOptions != nil {

		pkgTypeMap := map[PackageType]string{}

		if linux.PackagingOptions.Debian != nil {

			var pkgDetails *pkgutils.PackageDetails

			if packageMeta.IsChild() {

				pkgDetails = packageMeta.GetParent().MakePackageDetails(vdata.FullVersionString(), arch)

			} else {

				pkgDetails = packageMeta.MakePackageDetails(vdata.FullVersionString(), arch)

			}

			_, err := linux.PackagingOptions.Debian.createContentMapForArch(arch, map[string]string{binaryTarget: packageMeta.Name}, pkgBuildRoot)
			if err != nil {
				return nil, err
			}

			debPkgPath := filepath.Join(pkgBuildRoot, "deb", arch)
			if !vutils.Files.PathExists(debPkgPath) {
				if packageMeta.IsChild() {
					err := pkgutils.Debian.NewDebianPackage(pkgBuildRoot, pkgDetails, nil)
					if err != nil {
						return nil, err
					}
				} else {
					err := pkgutils.Debian.NewDebianPackage(pkgBuildRoot, pkgDetails, nil)
					if err != nil {
						return nil, err
					}
				}
			}

			if packageMeta.IsChild() {
				linux._packageQueue.AddPackageBuildFuncToArch(arch, DebianPackage, linux.PackagingOptions.Debian.runBuild)
				linux.compileSettings.addPackageProcessToQueue(LINUX, linux.ProcessPackageQueue, false)
			} else {
				path, err := linux.PackagingOptions.Debian.runBuild(arch)
				if err != nil {
					return nil, err
				}
				pkgTypeMap[DebianPackage] = path
			}

		}

		return pkgTypeMap, nil

	}

	return nil, nil

}

func (linux *XGOLinuxCompileSettings) ProcessPackageQueue() (PackageBuildMap, error) {

	return nil, NoPackagesToBuildError

}
