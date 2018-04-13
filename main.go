package main

import (
	"github.com/fsouza/fake-gcs-server/fakestorage"
	"time"
)

func main() {
	_, err := fakestorage.NewServerWithHostPort(nil, "0.0.0.0", 8060)
	if err != nil {
		panic(err)
	}

	for ;; {
		time.Sleep(time.Duration(time.Second * 5))
	}
}

