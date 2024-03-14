package main

import (
	"flag"
	"github.com/saisesai/doghole"
)

func init() {
	var err error
	flag.Parse()
	err = doghole.PrepareLog()
	if err != nil {
		panic(err)
	}
}
