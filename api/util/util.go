package util

import (
	"github.com/bigmamallc/env"
	"io/ioutil"
)

func MustWriteToTmpFile(data string) string {
	file, err := ioutil.TempFile("/tmp", "*")
	if err != nil {
		panic(err)
	}
	if _, err := file.WriteString(data); err != nil {
		panic(err)
	}
	return file.Name()
}

func MustSetDefaultCfg(cfg interface{}) {
	if err := env.SetDefaultOnly(cfg, ""); err != nil {
		panic(err)
	}
}
