package main

import (
	"path"
	"runtime"
)

func getRootDir() string {
	_, filename, _, _ := runtime.Caller(1)
	return path.Dir(filename)
}

func initParser() error {
	dir := getRootDir()

	err := InitMorph(path.Join(dir, "morph.bin"))
	if err != nil {
		return err
	}

	err = InitRules(path.Join(dir, "rules.yaml"))
	if err != nil {
		return err
	}

	InitCache(256)

	return nil
}

func finalizeParser() {
	FinalizeMorph()
	FinalizeRules()
	FinalizeCache()
}
