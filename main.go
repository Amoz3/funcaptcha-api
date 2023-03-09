package main

import (
	"fmt"
	"funcaptcha_api/gamevariants"
	"funcaptcha_api/signals"
	"io"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/gin-gonic/gin"
)

var solverMap = make(map[string]gamevariants.AudioSolver)

func main() {
	// signals.PatternDetect(signals.OptionOnePattern, "./test/3_footsteps.wav")
	signals.SavePattern("./test/option_two.wav", "./test/pattern_two.txt")
	signals.SavePattern("./test/option_three.wav", "./test/pattern_three.txt")
	return
	gin := gin.Default()

	gin.POST("/solve", solveHandler)

	solverMap["crowdsound"] = gamevariants.CrowdSoundSolver

	go jannyFunc()
	gin.Run(":9911")
}

func solveHandler(c *gin.Context) {

	gamevariant := c.PostForm("gamevariant")
	timestamp := time.Now().UnixNano()

	solver, contains := solverMap[gamevariant]
	if !contains {
		c.String(210, "game not supported")
	}

	file, err := c.FormFile("file")
	if err != nil {
		if strings.Contains(err.Error(), "no such file") {
			c.String(http.StatusBadRequest, "no file included")
			return
		}
		panic(err)
	}

	destinationDir := fmt.Sprintf("./%s/%d/", gamevariant, timestamp)
	destinationPath := fmt.Sprintf("%s%s", destinationDir, file.Filename)

	src, err := file.Open()
	if err != nil {
		panic(err)
	}
	defer src.Close()

	dst, err := createInFolder(destinationPath)
	if err != nil {
		panic(err)
	}

	io.Copy(dst, src)
	fmt.Printf("File path = %s\n", dst.Name())
	c.String(http.StatusOK, fmt.Sprintf("%d", solver.Solve(destinationPath, destinationDir))) // todo return int 1 - 3 for which one was solved
}

func createInFolder(p string) (*os.File, error) {
	err := os.MkdirAll(filepath.Dir(p), 0770)
	if err != nil {
		return nil, err
	}
	return os.Create(p)
}

// files are save in a dir with the nanosecond timestamp of the request
// delete all dirs like that older than x time (2 minute, could probably be like <30 seconds)
func jannyFunc() {
	const maxTime int64 = 120000000000 // 2 mins in nanoseconds
	fmt.Println(time.Now().UnixNano())
	for {
		filepath.Walk("./", func(path string, info fs.FileInfo, err error) error {
			if err != nil {
				return nil
			}

			if !info.IsDir() {
				return nil
			}

			if strings.Contains(path, ".git") {
				return nil
			}

			dirtimestamp, err := strconv.Atoi(info.Name())
			if err != nil {
				return nil // should just mean this is not a number
			}

			if int64(dirtimestamp)+maxTime <= time.Now().UnixNano() {
				color.Red("Jannied " + path)
				os.RemoveAll(path)
			}
			return nil
		})
	}
}
