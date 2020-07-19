// Copyright 2016 Google Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// TODO: High-level file comment.
package main

import (
	"flag"
	"fmt"
	"inu1255/sshdog/daemon"
	"io/ioutil"
	"os"
)

type Debugger bool

func (d *Debugger) Debug(format string, args ...interface{}) {
	if *d {
		msg := fmt.Sprintf(format, args...)
		fmt.Fprintf(os.Stderr, "[DEBUG] %s\n", msg)
	}
}

var dbg *Debugger = (*Debugger)(flag.Bool("debug", false, "enable debug"))
var daemonMode *bool = flag.Bool("daemon", false, "enable daemon")
var port *int = flag.Int("port", 2222, "port")
var hostkey *string = flag.String("hostkey", "", "hostkey default random")
var authkey *string = flag.String("authkey", "", "authkey")

func main() {
	flag.Parse()
	if *authkey == "" {
		flag.Usage()
		return
	}

	if *daemonMode {
		if err := daemon.Daemonize(daemonStart); err != nil {
			dbg.Debug("Error daemonizing: %v", err)
		}
	} else {
		waitFunc, _ := daemonStart()
		if waitFunc != nil {
			waitFunc()
		}
	}
}

// Actually run the implementation of the daemon
func daemonStart() (waitFunc func(), stopFunc func()) {
	server := NewServer()

	if *hostkey == "" {
		if err := server.RandomHostkey(); err != nil {
			dbg.Debug("Error adding random hostkey: %v", err)
			return
		}
	} else {
		keyData, err := ioutil.ReadFile(*hostkey)
		if err != nil {
			dbg.Debug("Error read hostkey: %v", err)
			return
		}
		dbg.Debug("Adding hostkey file: %s", *hostkey)
		if err = server.AddHostkey(keyData); err != nil {
			dbg.Debug("Error adding public key: %v", err)
			return
		}
	}

	if *authkey != "" {
		authData, err := ioutil.ReadFile(*authkey)
		if err != nil {
			dbg.Debug("Error read hostkey: %v", err)
			return
		}
		dbg.Debug("Adding authorized_keys")
		server.AddAuthorizedKeys(authData)
	}

	server.ListenAndServe(int16(*port))
	return server.Wait, server.Stop
}
