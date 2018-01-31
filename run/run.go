package run

import (
	"io"
	"github.com/liujiarik/goassem/conf"
	"github.com/liujiarik/goassem/help"
	"flag"
	"fmt"
	"os"
	"encoding/json"
	"path/filepath"
	"strings"
	"os/exec"
	"bytes"
	"github.com/pierrre/archivefile/zip"
)

const (
	OUT_PATH = "_out"
)

/*
@author:liujia43
@since:2018/1/31
*/
func Run(w io.Writer, appArgs []string) (string, error) {

	if len(appArgs) == 1 {
		return help.HelpFull, nil
	}
	flags := flag.NewFlagSet("goassem", flag.ContinueOnError)
	version := flags.Bool("version", false, "show goassem version")
	licenses := flags.Bool("licenses", false, "show goassem's licenses")

	err := flags.Parse(appArgs[1:])
	if err != nil {
		return help.HelpFull, err
	}

	if *version {
		return help.VersionMsg, nil
	}
	if *licenses {
		return help.VersionMsg, nil
	}

	args := flags.Args()

	cmd := args[0]
	switch cmd {
	case "init":
		return Init(w)
	case "package":
		return Package(w, appArgs)
	case "clear":
		return Clear(w, appArgs)
	default:
		return help.HelpFull, fmt.Errorf("Unknown command %q", cmd)
	}

	return "", nil
}

func Init(w io.Writer) (string, error) {

	// check file exist
	if _, err := os.Stat(conf.CONF_FILE); err == nil {
		return "assembly.json exist!", nil
	}
	f, err := os.Create(conf.CONF_FILE)
	if err != nil {
		return "", fmt.Errorf("create %s happen error | %q", conf.CONF_FILE, err)
	}
	defer f.Close()

	c := conf.NewDefaultConf()
	cb, _ := json.Marshal(c)
	f.WriteString(string(cb))
	f.Sync()

	return help.InitMsg, nil

}
func Package(w io.Writer, appArgs []string) (string, error) {
	c, er := conf.Load()
	if er != nil {
		fmt.Fprintf(os.Stderr, "Load assembly.json happen error | %+v\n", er)
		os.Exit(2)
	}
	l := len(c)
	if l == 0 {
		return "assembly.json is empty", nil
	}
	if ie, _ := FileOrDirIsExists(OUT_PATH); !ie {
		os.Mkdir(OUT_PATH, os.ModePerm) //create out dir
	}

	for i := 0; i < l; i++ {
		fmt.Fprintln(w, "##start to assem package NO.", i)
		af := c[i] // assembly conf
		af.Print(w)

		pn := GetPackageName(&c[i])
		k := len(pn)
		for j := 0; j < k; j++ {
			packageDir := filepath.Join(OUT_PATH, pn[j])    //package dir
			CheckDirExistsAndCreate(packageDir)             // create package dir
			binDir := filepath.Join(packageDir, af.BinDir); // bin dir
			CheckDirExistsAndCreate(binDir)                 // create bin dir
			if r := goBuild(w, af.BuildArgs, af.Main+".go"); r == nil {
				CopyFile(filepath.Join(binDir, af.Main), af.Main)
			} else {
				return help.Fail, r
			}
			fileSet(w, af.FileSets, packageDir)
			switch af.Format {
			case "zip":
				zipFile := packageDir + ".zip"
				zip.ArchiveFile(packageDir, zipFile, nil)
				return help.AllDone, nil
			case "tar.gz":
				return help.AllDone, nil
			default:
				return af.Format + "is'surport", nil

			}
		}

	}

	return help.AllDone, nil

}

func Clear(w io.Writer, appArgs []string) (string, error) {

	os.RemoveAll(OUT_PATH);
	return "", nil

}

func GetPackageName(ac *conf.AssemblyConf) []string {

	packageNames := make([]string, 0)

	if ac.Platforms == nil || len(ac.Platforms) == 0 {
		fileItem := make([]string, 0)
		fileItem = append(fileItem, ac.Name)
		fileItem = append(fileItem, ac.Version)
		packageNames = append(packageNames, strings.Join(fileItem, "-"))
	} else {
		for i := 0; i < len(ac.Platforms); i++ {
			fileItem := make([]string, 0)
			fileItem = append(fileItem, ac.Name)
			fileItem = append(fileItem, ac.Version)
			fileItem = append(fileItem, ac.Platforms[i].Os)
			fileItem = append(fileItem, ac.Platforms[i].Arch)
			packageNames = append(packageNames, strings.Join(fileItem, "-"))
		}

	}

	return packageNames

}
func fileSet(w io.Writer, fileSets []*conf.FileSet, packageDir string) {
	l := len(fileSets)
	for i := 0; i < l; i++ {
		fs := fileSets[i]
		k := len(fs.Includes)
		for j := 0; j < k; j++ {
			include := fs.Includes[j]
			sourceFile := filepath.Join(fs.Directory, include);
			targetDir := filepath.Join(packageDir, fs.OutputDirectory);
			CheckDirExistsAndCreate(targetDir)
			targetFile := filepath.Join(targetDir, include);
			CheckFileExistsAndCreate(targetFile)
			_, e := CopyFile(sourceFile, targetFile)
			if e != nil {

				fmt.Fprintln(w, e.Error())
			}

		}
	}

}

func CheckDirExistsAndCreate(path string) {
	if ie, _ := FileOrDirIsExists(path); !ie {
		os.MkdirAll(path, os.ModePerm) // create  dir
	}
}

func CheckFileExistsAndCreate(path string) {
	if ie, _ := FileOrDirIsExists(path); !ie {
		os.Create(path) // create  dir
	}
}

func execCommand(w io.Writer, commandName string, arg ...string) error {
	cmd := exec.Command(commandName, arg...)
	fmt.Fprintln(w, cmd.Args)
	cmd.Stdin = strings.NewReader("some input")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return err
	}
	fmt.Fprint(w, "exec :%q\n", out.String())

	return nil
}

func goBuild(w io.Writer, buildParams []string, goFile string) error {

	fmt.Fprintln(w, "start to build...")
	arg := make([]string, 0)
	arg = append(arg, "build")
	arg = append(arg, buildParams...)
	arg = append(arg, goFile)
	return execCommand(w, "go", arg...)
}
