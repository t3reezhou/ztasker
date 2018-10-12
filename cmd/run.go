package cmd

import (
	"bufio"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/t3reezhou/ztasker/bar"
	l "github.com/xtimeline/gox/log"
)

var logFileBuffer *bufio.Writer

type Tasker interface {
	TaskNum() int
	TaskName() string
	TaskConf() interface{}
	Run(cc chan error, args []string)
}

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "run tasks",
}

func run(task Tasker, args []string) error {
	if task.TaskNum() <= 0 {
		panic("nothing to tasker")
	}
	cc := make(chan error, task.TaskNum())
	logFile, err := os.OpenFile(fmt.Sprintf("%s.log", task.TaskName()), os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return l.Log(err)
	}
	defer logFile.Close()
	logFileBuffer = bufio.NewWriter(logFile)
	// task.SetLog(logFileBuffer)
	go task.Run(cc, args)

	finish := 0
	for err := range cc {
		if err != nil {
			logFileBuffer.WriteString(fmt.Sprintf("%s\n", err.Error()))
		}
		finish++
		ter, _ := bar.TerminalWidth()
		fmt.Print("\r %s\r", bar.Icon(finish, task.TaskNum(), ter))
		if finish >= task.TaskNum() {
			close(cc)
		}
		logFileBuffer.Flush()
	}
	return nil
}

func init() {
	RootCmd.AddCommand(runCmd)
}
