// Copyright 2016 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package help

var version = "v1.0.0"
//go:generate govendor license -o licenses.go -template gen-license.template

var HelpFull = `goassem (` + version + `): 
	-licenses    Show goassem's licenses.
	-version     Show goassem version

Sub-Commands

	init     Create the "assembly.json" file.
	package  List and filter existing dependencies and packages.

Ignoring files with build tags, or excluding packages from being vendored:
	The "vendor.json" file contains a string field named "ignore".
	It may contain a space separated list of build tags to ignore when
	listing and copying files.
	This list may also contain package prefixes (containing a "/", possibly
	as last character) to exclude when copying files in the vendor folder.
	If "foo/" appears in this field, then package "foo" and all its sub-packages
	("foo/bar", …) will be excluded (but package "bar/foo" will not).
	By default the init command adds the "test" tag to the ignore list.
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
