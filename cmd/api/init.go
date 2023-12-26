package main

import (
	"gomodel/cmd/api/rest"
)

func Teardown() {

}

func Init(
	rest *rest.RestServer,
) {
	rest.Run()
}
