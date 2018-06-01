package core

import (
	"os"
	"os/exec"
	"io/ioutil"
	"strings"
	"regexp"
	"github.com/zekroTJA/serverManager2/util"
	"path/filepath"
)


type Screen struct {
	Uid int
	Id, Name, Started string
}

func GetRunningScreens() []Screen {
	out, _ := exec.Command("screen", "-ls").Output()
	outsplit := strings.Split(string(out), "\n")
	regex := regexp.MustCompile(`[()]`)
	
	screens := []Screen {}
	for i, e := range outsplit[1:len(outsplit)-3]  {
		fields := strings.Fields(e)
		nameandid := strings.Split(fields[0], ".")
		screens = append(screens, Screen {
			i, 
			nameandid[0], 
			nameandid[1], 
			regex.ReplaceAllString(fields[1] + " " +  fields[2], ""),
		})
	}

	return screens
}

func GetServers(location string) []Screen {
	screens := []Screen {}
	filepath.Walk(location, func(path string, info os.FileInfo, err error) error {
		folder := strings.Replace(path, location, "", -1)
		pathsplit := strings.Split(folder, "/")
		if len(pathsplit) == 2 {
			screens = append(screens, Screen {
				Uid: len(screens),
				Name: folder[1:] })
		}
		return err
	})
	return screens
}

func SliceContainsServer(slc []Screen, server Screen) bool {
	for _, e := range slc {
		if e.Name == server.Name {
			return true
		}
	}
	return false
}

// SCREEN ACTION FUNCTIONS

func StartScreen(screen Screen, screens []Screen, config util.Conf, runInLoop bool) {
	if SliceContainsServer(screens, screen) {
		util.LogError("Screen '" + screen.Name + "' is still running!")
		pause()
		return
	}

	startfile := config.ServerLocation + "/" + screen.Name + "/run.sh"

	_, err := ioutil.ReadFile(startfile)
	if os.IsNotExist(err) {
		util.LogError("This server has no 'run.sh' file specified.\n" + 
					  "Please create this file in the root directory of the server with the command to start.")
		pause()
		return
	} else if err != nil {
		util.LogError("An unexpected error occured opening 'run.sh' of this server:\n" + err.Error())
		pause()
		return
	}

	if runInLoop {
		ioutil.WriteFile(".runner", []byte(
			"# This is a autogenerated file.\n" +
			"# Please, do not delete this file!\n\n" +
			"while true; do bash $1; done"), 0644)
		exec.Command("screen", "-dmLS", screen.Name, "bash", ".runner", startfile).Run()
	} else {
		exec.Command("screen", "-dmS", screen.Name, "bash", startfile).Run()
	}
}

func StopScreen(screen Screen, screens []Screen, config util.Conf) {
	if !SliceContainsServer(screens, screen) {
		util.LogError("Screen '" + screen.Name + "' is not running!")
		pause()
		return
	}

	exec.Command("screen", "-XS", screen.Name, "quit").Run()
}

func ResumeScreen(screen Screen, screens []Screen, config util.Conf) {
	if !SliceContainsServer(screens, screen) {
		util.LogError("Screen '" + screen.Name + "' is not running!")
		pause()
		return
	}

	cmd := exec.Command("screen", "-r", screen.Name)
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Run()
}

func RestartScreen(screen Screen, screens []Screen, config util.Conf, runInLoop bool) {

}