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
	"runtime"
)

const (
	OUT_PATH = "_out"
)

/*
@author:liujiarik
@since:2018/1/31
*/
func Run(w io.Writer, appArgs []string) (string, error) {

	if len(appArgs) == 1 {
		return help.HelpFull, nil
	}
	flags := flag.NewFlagSet("goassem", flag.ContinueOnError)
	version := flags.Bool("version", false, "show goassem version")
	licenses := flags.Bool("licenses", false, "show goassem's licenses")
	helps := flags.Bool("help", false, "show goassem's helps")

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
	if *helps {
		return help.HelpFull, nil
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

	return "", nil

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
		fmt.Fprintln(w, "==start to assem package NO.", i)
		af := c[i] // assembly conf
		af.Print(w)

		if af.Platforms == nil || len(af.Platforms) == 0 {
			// no mention of the platform.It use local platform
			platforms := make([]*conf.Platform, 0)
			platforms = append(platforms, GetLocalPlatform())
			af.Platforms = platforms
		}

		k := len(af.Platforms)
		for j := 0; j < k; j++ {

			platform := af.Platforms[j]                                  // platform info
			packageName := GetPackageName(af.Name, af.Version, platform) // package name
			packageDir := filepath.Join(OUT_PATH, packageName)           //package dir
			CheckDirExistsAndCreate(packageDir)                          // create package dir
			binDir := filepath.Join(packageDir, af.BinDir);              // bin dir
			CheckDirExistsAndCreate(binDir)                              // create bin dir
			if r := goBuild(w, af.BuildArgs, af.Main+".go", platform); r == nil {
				binName := getBinName(af.Main)
				CopyFile(filepath.Join(binDir, binName), binName)
				os.Remove(binName)
			} else {
				return help.Fail, r
			}
			fileSet(w, af.FileSets, packageDir) // deal fileSets
			switch af.Format {
			case "zip":
				zipFile := packageDir + ".zip"
				zip.ArchiveFile(packageDir, zipFile, nil) // ArchiveFile
				break
			case "tar.gz":
				return af.Format + "is'surport", nil
			default:
				return af.Format + "is'surport", nil
			}
		}
	}
	return help.AllDone, nil
}

func getBinName(goFile string) string {

	f := strings.Split(goFile, "/")
	return f[len(f)-1]
}

func Clear(w io.Writer, appArgs []string) (string, error) {

	os.RemoveAll(OUT_PATH);
	return "", nil

}

func GetPackageName(name string, version string, p *conf.Platform) string {

	fileItem := make([]string, 0)
	fileItem = append(fileItem, name)
	fileItem = append(fileItem, version)
	fileItem = append(fileItem, p.Os)
	fileItem = append(fileItem, p.Arch)
	return strings.Join(fileItem, "-")
}
func fileSet(w io.Writer, fileSets []*conf.FileSet, packageDir string) {
	l := len(fileSets)
	for i := 0; i < l; i++ {
		fs := fileSets[i]

		walkFileList(fs.Directory, func(file string) error {
			if CheckMatchInFileRegexps(file, fs.Includes) {
				sourceFile := filepath.Join(fs.Directory, file);
				targetDir := filepath.Join(packageDir, fs.OutputDirectory);
				CheckDirExistsAndCreate(targetDir)
				targetFile := filepath.Join(targetDir, file);
				CheckFileExistsAndCreate(targetFile)
				_, e := CopyFile(targetFile, sourceFile)
				if e != nil {
					fmt.Fprintln(w, e.Error())
				}
			}
			return nil
		})
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

func execCommand(w io.Writer, commandName string, env []string, arg ...string) error {
	cmd := exec.Command(commandName, arg...)
	fmt.Fprintln(w, cmd.Args)
	cmd.Stdin = strings.NewReader("some input")

	if env != nil && len(env) > 0 {
		cmd.Env = append(os.Environ(), env...)
	}

	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return err
	}
	fmt.Fprint(w, "exec :\n", out.String())

	return nil
}

func goBuild(w io.Writer, buildParams []string, goFile string, p *conf.Platform) error {

	checkIsLocalEnv(p)

	fmt.Fprintln(w, "start to build...")
	arg := make([]string, 0)
	arg = append(arg, "build")
	arg = append(arg, buildParams...)
	arg = append(arg, goFile)

	if checkIsLocalEnv(p) {
		return execCommand(w, "go", nil, arg...)
	} else {
		env := append(os.Environ(),
			"CGO_ENABLED=0",
			"GOOS="+strings.ToLower(p.Os),
			"GOARCH="+strings.ToLower(p.Arch),
		)
		return execCommand(w, "go", env, arg...)
	}

}

func GetLocalPlatform() *conf.Platform {
	return &conf.Platform{runtime.GOARCH, runtime.GOOS}
}

func checkIsLocalEnv(p *conf.Platform) bool {

	return strings.ToLower(p.Arch) == runtime.GOARCH && strings.ToLower(p.Os) == runtime.GOOS
}
