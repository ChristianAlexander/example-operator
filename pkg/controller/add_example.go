package controller

import (
	"github.com/christianalexander/example-operator/pkg/controller/example"
)

func init() {
	// AddToManagerFuncs is a list of functions to create controllers and add them to a manager.
	AddToManagerFuncs = append(AddToManagerFuncs, example.Add)
}
