package launcher

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"runtime"

	"github.com/Masterminds/semver/v3"
)

type launcher1155 struct {
	json LauncherjsonX115
	flag []string
	*Gameinfo
	//run launcher1155.cp to set this
	fixlog4j bool
}

func (l *launcher1155) GetLauncherjsonX115() LauncherjsonX115 {
	return l.json
}

func newlauncher1155(json LauncherjsonX115) *launcher1155 {
	flag := make([]string, 0)
	return &launcher1155{json: json, flag: flag}
}

func (l launcher1155) Launcher115() error {
	fmt.Println(l.flag)
	var cmd *exec.Cmd
	if l.JavePath == "" {
		l.JavePath = "java"
	}
	if l.Log {
		cmd = exec.Command(l.JavePath, l.flag...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Dir = l.Gamedir
		err := cmd.Run()
		if err != nil {
			if err != nil {
				return fmt.Errorf("launcher1155.Launcher115: %w", err)
			}
		}
	} else {
		if runtime.GOOS == "windows" && l.JavePath == "java" {
			cmd = exec.Command("javaw", l.flag...)
		} else {
			l.flag = append(l.flag, "&")
			cmd = exec.Command(l.JavePath, l.flag...)
		}
		cmd.Dir = l.Gamedir
		err := cmd.Start()
		if err != nil {
			return fmt.Errorf("launcher1155.Launcher115: %w", err)
		}
	}
	return nil
}

func (l *launcher1155) cp() string {
	path := l.Minecraftpath + `/libraries/`
	b := bytes.NewBuffer(nil)
	const log4jpackname = "org.apache.logging.log4jlog4j-core"
	for _, p := range l.json.Libraries {
		if !Ifallow(p) {
			continue
		}
		pack := FullLibraryX115(&p, "")
		if p.Downloads.Artifact.Path == "" {
			continue
		}
		key := pack[0] + pack[1]
		v, ok := l.Gameinfo.flag[key]
		add := func() {
			b.WriteString(path)
			if p.Downloads.Classifiers != nil && p.Rules != nil {
				b.WriteString(p.Downloads.Classifiers["natives"+p.Rules[0].Os.Name].Path)
			} else {
				b.WriteString(p.Downloads.Artifact.Path)
			}
			b.WriteByte(os.PathListSeparator)
			if key == log4jpackname && needFixlog4j(pack[2]) {
				l.fixlog4j = true
			}
		}
		if ok {
			if v == pack[2] {
				add()
			}
		} else {
			add()
		}
	}
	b.WriteString(l.Minecraftpath + `/versions/` + l.json.ID + `/` + l.json.ID + `.jar`)
	return b.String()
}

func needFixlog4j(ver string) bool {
	v, err := semver.NewVersion(ver)
	if err != nil {
		return true
	}
	// CVE-2021-45046 CVE-2021-44228 has been addressed in 2.16.0
	if v.Major() >= 2 && v.LessThan(semver.MustParse("2.16.0")) {
		return true
	}
	return false
}

// Deprecated: 之前想清理安装多余的库，就添加了这个函数用来导出某个版本所引入的库。可惜 forge 并不会在 json 中写上所有它导入的库，因此这个函数也就没有意义了。
func (l *launcher1155) CP() []string {
	path := l.Minecraftpath + `/libraries/`
	list := make([]string, 0, len(l.json.Libraries))
	for _, p := range l.json.Libraries {
		pack := Name2path(p.Name)
		v, ok := l.Gameinfo.flag[pack[0]+pack[1]]
		add := func() {
			list = append(list, path+p.Downloads.Artifact.Path)
		}
		if ok {
			if v == pack[2] {
				add()
			}
		} else {
			add()
		}

	}
	return list
}

func Ifallow(l LibraryX115) bool {
	if l.Rules != nil {
		var allow bool
		for _, r := range l.Rules {
			if r.Action == "disallow" && osbool(r.Os.Name) {
				return false
			}
			if r.Action == "allow" && (r.Os.Name == "" || osbool(r.Os.Name)) {
				allow = true
			}
		}
		return allow
	}
	return true
}

func osbool(os string) bool {
	GOOS := runtime.GOOS
	if GOOS == "darwin" {
		GOOS = "osx"
	}
	return os == GOOS
}
