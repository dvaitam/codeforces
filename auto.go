package main

import (
	"flag"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
)

type task struct {
	cppPath string
	goPath  string
	logPath string
	dirName string
	suffix  string
}

func main() {
	var concurrency int
	flag.IntVar(&concurrency, "c", 5, "concurrency level")
	flag.Parse()

	if concurrency < 1 {
		concurrency = 1
	}

	logDir := "/root/log"
	err := os.MkdirAll(logDir, 0755)
	if err != nil {
		fmt.Printf("Error creating log directory: %v\n", err)
		return
	}

	var tasks []task
	err = filepath.WalkDir(".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		name := d.Name()
		if !strings.HasPrefix(name, "sol") || !strings.HasSuffix(name, ".cpp") {
			return nil
		}
		dir := filepath.Dir(path)
		basedir := filepath.Base(dir)
		suffix := strings.TrimPrefix(strings.TrimSuffix(name, ".cpp"), "sol")
		goFile := filepath.Join(dir, basedir+suffix+".go")
		_, err = os.Stat(goFile)
		if err == nil {
			// File exists, skip
			return nil
		}
		if !os.IsNotExist(err) {
			// Some other error, log but continue
			fmt.Printf("Error checking %s: %v\n", goFile, err)
			return nil
		}
		logFile := filepath.Join(logDir, basedir+suffix+".log")
		tasks = append(tasks, task{
			cppPath: path,
			goPath:  goFile,
			logPath: logFile,
			dirName: basedir,
			suffix:  suffix,
		})
		return nil
	})
	if err != nil {
		fmt.Printf("Error walking directory: %v\n", err)
		return
	}

	if len(tasks) == 0 {
		fmt.Println("No matching files found to process.")
		return
	}

	sem := make(chan struct{}, concurrency)
	var wg sync.WaitGroup

	for _, t := range tasks {
		sem <- struct{}{} // Acquire slot
		wg.Add(1)
		go func(t task) {
			defer wg.Done()
			defer func() { <-sem }()
			fmt.Printf("Starting codex for %s %s\n", t.dirName, t.suffix)
			cmdStr := fmt.Sprintf("convert %s to go as %s make sure it builds", t.cppPath, t.goPath)
			cmd := exec.Command("codex", "exec", "--full-auto", cmdStr)
			output, err := cmd.CombinedOutput()
			// Save output to log file
			writeErr := os.WriteFile(t.logPath, output, 0644)
			if writeErr != nil {
				fmt.Printf("Error writing log for %s: %v\n", t.cppPath, writeErr)
			}
			if err != nil {
				fmt.Printf("Error running command for %s: %v\nLog saved to: %s\n", t.cppPath, err, t.logPath)
			} else {
				fmt.Printf("Successfully processed %s to %s\nLog saved to: %s\n", t.cppPath, t.goPath, t.logPath)
			}
		}(t)
	}

	wg.Wait()
	fmt.Println("All tasks completed.")
}
