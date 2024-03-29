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

	gin := gin.Default()

	gin.POST("/solve", solveHandler)
	
	// solverMap["smeep"] = *gamevariants.RepeatingSongRegressionSolver.AudioSolver

	solverMap["crowdsound"] = gamevariants.CrowdSoundSolver
	// solverMap["count_3_footsteps"] = gamevariants.ThreeFootstepSolver
	solverMap["fake_cats"] = gamevariants.FakeCatsSolver
	solverMap["fake_bees"] = gamevariants.FakeBeesSolver
	solverMap["recognise_drums"] = gamevariants.RecogniseDrumsSolver

	go jannyFunc()
	gin.Run(":" + os.Getenv("PORT"))
}

func savePatterns(path string, fileInfo fs.FileInfo, err error) error {
	if err != nil || fileInfo.IsDir() || !strings.Contains(fileInfo.Name(), "converted") {
		return err
	}

	signals.SavePattern(path, fmt.Sprintf("./crowdsound_whole_%s.txt", fileInfo.Name()))
	return nil
}

func splitAll(path string, fileInfo fs.FileInfo, err error) error {
	if err != nil {
		return err
	}

	if fileInfo.IsDir() {
		return nil
	}

	gamevariants.SplitRepeatingSongs(path, "./test/splits/")
	return nil
}

func solveHandler(c *gin.Context) {

	gamevariant := c.PostForm("gamevariant")
	timestamp := time.Now().UnixNano()

	solver, contains := solverMap[gamevariant]
	if !contains {
		color.Cyan("Game upsupported: " + gamevariant)
		c.String(210, "game not supported")
		return
	}

	file, err := c.FormFile("file")
	if err != nil {
		if strings.Contains(err.Error(), "no such file") {
			c.String(http.StatusBadRequest, "no file included")
			return
		}
		panic(err)
	}

	destinationDir := fmt.Sprintf("./working/%s/%d/", gamevariant, timestamp)
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
