// Copyright 2016 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package help

var version = "v1.0.0"
//go:generate govendor license -o licenses.go -template gen-license.template

var HelpFull = `goassem (` + version + `): 
	-licenses    Show goassem's licenses.
	-version     Show goassem version
	-help     	 Show goassem help

Sub-Commands

	init     Create the "assembly.json" file.
	package  package all file
	clear	 clear all file in _out

	@author liujiarik
`
var helpConf = `goassem can't load assembly.json.: `

var VersionMsg = version + `
`
var LicensesMsg = version + `
`
var InitMsg = `assembly.json has already been created!`

var AllDone = `success!`
var Fail = ` fail!`
