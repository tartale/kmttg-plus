package logz

func init() {

	err := InitLoggers()
	if err != nil {
		panic(err)
	}

	err = InitThirdPartyLoggers()
	if err != nil {
		panic(err)
	}

	err = InitDebugDir()
	if err != nil {
		panic(err)
	}

}
