package shell

const (
	testAndroidHome          = "ANDROID_HOME"
	testCaseHomeInTemplate   = "Home in template"
	testCaseNoTemplate       = "No template"
	testGitBangEchoWorld     = "!echo world"
	testHelloEnvName         = "HELLO"
	testHomeLinux            = "/home/trajano"
	testHomeUnix             = "/Users/jan"
	testHomeWindows          = `C:\Users\trajano`
	testNameBarUpper         = "BAR"
	testNameFooUpper         = "FOO"
	testPathPwshWindows      = `$env:PATH = "C:\Users\jan\.tools\bin" + ';' + $env:PATH`
	testShellUnknown         = "FOO"
	testUsrLocalBin          = "/usr/local/bin"
	testUsrLocalBinAndUsrBin = "/usr/local/bin\n/usr/bin"
	testUsrLocalShare        = "/usr/local/share"
	testUsrBin               = "/usr/bin"
	testUsrShare             = "/usr/share"
	testWindowsBins          = "C:\\bin\nD:\\bin"
)
