#  Go automatic assembly package tools: gossam
[中文wiki](https://github.com/liujiarik/goassem/wiki/goassem-chinese-wiki)
Goassem reference  the design idea of maven assembly build plug-ins .By executing  `goassem package ` command, the  Go source codes  will be compiled into different platforms executable files. At the same time,  documentation, configuration, external script files will be archived into one or more compressed package orderly .It is very convenient to release and deploy Go applications.

## Feature
1. build a variety of platforms  Go executable file in batches.
2. archive executable file and external resource file (scripts, configuration, documentation, etc.) into package files. Furthermore, developer can define the structure and version of package as their want.

## Install
 installation is very simple
```
go get -u github.com/liujiarik/goassem
```

 A brief description
```
goassem -help
```

## Initial
Goassem use **assembly. json** file to save packaging process.Therefore, we need to create this file in project root path.
```
goassem init
```
After executing,  the **assembly. json** template file will be created. The developer can edit this file to describe packaging process.

## The instructions of  assembly.json file
 It is crucial to understand assembly.json for user. The file define the process of build and packaging。
The following is a simple assembly.json file :

``` json
[
  {
    "name": "demo",
    "version": "1.0.1",
    "format": "zip",
    "main": "main",
    "binDir": "bin",
    "buildArgs": [],
    "fileSets": [
      {
        "directory": "conf",
        "outputDirectory": "conf",
        "includes": [
          "*"
        ]
      },
      {
        "directory": "sh",
        "outputDirectory": "/",
        "includes": [
          "*.sh"
        ]
      },
      {
        "directory": "./",
        "outputDirectory": "/",
        "includes": [
          "README"
        ]
      }
    ],
    "platforms": [
      {
        "arch": "amd64",
        "os": "linux"
      },
      {
        "arch": "amd64",
        "os": "darwin"
      }
    ]
  }
]
```

### Parameter description：
**name** ：package name

**version** ：package version

**format** ：package Archive format （Currently only support zip ）

**main** ：Main Go file name without the suffix

**binDir** ：bin dir.

**buildArgs** ：go buildArgs.It is a string array

**fileSets** : A set of file copy command

* ~directory~ ：source directory .the path is relative to project root path. If  directory is the project root directory, please use ". /"
* ~outputDirectory~:  target directory. The path is relative to package directory
* ~includes~: A set of file which are need to copy ,and It is a string array. includes parameter support ‘*’  wildcard

**platforms** :  Cross-compilation support
* ~arch~ ：target  arch
* ~os~ :  target  os

## Package
```
goassem package
```

Goassem read assembly. json  file and package Go project. All outputs  will be  located  at  ` _out ` folder

the example  above will package two archive files :

1. demo-1.0.1-linux-amd64.zip
2. demo-1.0.1-darwin-amd64.zip


## Clear Output
```
goassem clear
```

remove all file and dir in `_out` folder