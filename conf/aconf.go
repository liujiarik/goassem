/*
 conf package define the struct object of assembly.json file,and to parse it.
 */
package conf

import (
	"io/ioutil"
	"encoding/json"
	"os"
	"strings"
	"path/filepath"
	"io"
	"fmt"
)

const (
	CONF_FILE = "assembly.json"
)

type AssemblyConf struct {
	Name      string      `json:"name"`
	Version   string      `json:"version"`
	Format    string      `json:"format"`
	Main      string      `json:"main"`
	BinDir    string      `json:"binDir"`
	BuildArgs []string    `json:"buildArgs"`
	FileSets  []*FileSet  `json:"fileSets"`
	Platforms []*Platform `json:"platforms"`
}

type FileSet struct {
	Directory       string   `json:"directory"`
	OutputDirectory string   `json:"outputDirectory"`
	Includes        []string `json:"includes"`
}

type Platform struct {
	Arch string `json:"arch"`
	Os   string `json:"os"`
}

func (a AssemblyConf) Print(w io.Writer) {

	fmt.Fprintln(w, "Name : ", a.Name)
	fmt.Fprintln(w, "Version : ", a.Version)
	fmt.Fprintln(w, "Format : ", a.Format)
	fmt.Fprintln(w, "Main Go File : ", a.Main)
	fmt.Fprintln(w, "BuildArgs : ", a.BuildArgs)

}

func NewDefaultConf() []*AssemblyConf {
	a := &AssemblyConf{}
	a.Format = "zip"
	a.Main = "main"
	a.BinDir = "bin"
	a.Name = GetProjectName()
	a.Version = "0.0.0"
	a.BuildArgs = make([]string, 0)
	a.FileSets = make([]*FileSet, 0)
	a.Platforms = make([]*Platform, 0)

	l := make([]*AssemblyConf, 0)
	l = append(l, a)
	return l
}

//Load conf file
func LoadFile(fileName string) ([]AssemblyConf, error) {

	a := make([]AssemblyConf, 0, 100)

	data, err := ioutil.ReadFile(fileName)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(data, &a); err != nil {
		return nil, err
	}
	return a, nil
}

func Load() ([]AssemblyConf, error) {
	return LoadFile(CONF_FILE)
}

func GetProjectName() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		panic(err)
	}
	s := strings.Split(dir, "/")
	return s[len(s)-1]
}
