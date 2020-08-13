package controller

import (
	"github.com/ChubaoMonitor/pkg/controller/chubaomonitor"
)

func init() {
	// AddToManagerFuncs is a list of functions to create controllers and add them to a manager.
	AddToManagerFuncs = append(AddToManagerFuncs, chubaomonitor.Add)
}
