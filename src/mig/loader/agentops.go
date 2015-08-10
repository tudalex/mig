// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// Contributor: Aaron Meihm ameihm@mozilla.com [:alm]

package main

import (
	"fmt"
	"mig"
	"os/exec"
)

func runTriggers() (err error) {
	defer func() {
		if e := recover(); e != nil {
			err = fmt.Errorf("runTriggers() -> %v", e)
		}
	}()
	ctx.Channels.Log <- mig.Log{Desc: "running triggers due to modification"}
	err = terminateAgent()
	if err != nil {
		ctx.Channels.Log <- mig.Log{Desc: fmt.Sprintf("%v (ignored)", err)}
	}
	// XXX Run agent service initialization here.
	return
}

// Terminate any running agent on the system.
func terminateAgent() (err error) {
	defer func() {
		if e := recover(); e != nil {
			err = fmt.Errorf("terminateAgent() -> %v", e)
		}
	}()
	hb, err := mig.GetHostBundle()
	if err != nil {
		panic(err)
	}
	var abe mig.BundleDictionaryEntry
	found := false
	for _, x := range hb {
		if x.Name == "agent" {
			abe = x
			found = true
			break
		}
	}
	if !found {
		panic("no agent entry in host bundle")
	}

	migcomm := exec.Command(abe.Path, "-q", "shutdown")
	err = migcomm.Run()
	if err != nil {
		panic(err)
	}

	return
}
