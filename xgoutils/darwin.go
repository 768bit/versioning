package xgoutils

var DARWIN_ALLOWED_ARCHITECTURES = []XGOArchitecture{ARCH_AMD64}

type XGODarwinCompileSettings struct {
	BaseXGOPlatformCompileSettings
}

func NewDarwinCompileSettings(architectures []XGOArchitecture) *XGODarwinCompileSettings {

	lcs := &XGODarwinCompileSettings{
		BaseXGOPlatformCompileSettings: newBaseXGOPlatformCompileSettings(DARWIN),
	}

	lcs.addArchitectures(architectures, DARWIN_ALLOWED_ARCHITECTURES)

	return lcs

}
