package xgoutils

var WINDOWS_ALLOWED_ARCHITECTURES = []XGOArchitecture{ARCH_386, ARCH_AMD64, ARCH_ARM64}

type XGOWindowsCompileSettings struct {
	BaseXGOPlatformCompileSettings
}

func NewWindowsCompileSettings(architectures []XGOArchitecture) *XGOWindowsCompileSettings {

	lcs := &XGOWindowsCompileSettings{
		BaseXGOPlatformCompileSettings: newBaseXGOPlatformCompileSettings(WINDOWS),
	}

	lcs.addArchitectures(architectures, WINDOWS_ALLOWED_ARCHITECTURES)

	return lcs

}
