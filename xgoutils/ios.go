package xgoutils

var IOS_ALLOWED_ARCHITECTURES = []XGOArchitecture{ARCH_ARM64}

type XGOIosCompileSettings struct {
	BaseXGOPlatformCompileSettings
}

func NewIosCompileSettings(architectures []XGOArchitecture) *XGOIosCompileSettings {

	lcs := &XGOIosCompileSettings{
		BaseXGOPlatformCompileSettings: newBaseXGOPlatformCompileSettings(IOS),
	}

	lcs.addArchitectures(architectures, IOS_ALLOWED_ARCHITECTURES)

	return lcs

}
