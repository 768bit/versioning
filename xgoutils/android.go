package xgoutils

var ANDROID_ALLOWED_ARCHITECTURES = []XGOArchitecture{ARCH_ARM_7, ARCH_ARM_6, ARCH_ARM64}

type XGOAndroidCompileSettings struct {
	BaseXGOPlatformCompileSettings
}

func NewAndroidCompileSettings(architectures []XGOArchitecture) *XGOAndroidCompileSettings {

	lcs := &XGOAndroidCompileSettings{
		BaseXGOPlatformCompileSettings: newBaseXGOPlatformCompileSettings(ANDROID),
	}

	lcs.addArchitectures(architectures, ANDROID_ALLOWED_ARCHITECTURES)

	return lcs

}
