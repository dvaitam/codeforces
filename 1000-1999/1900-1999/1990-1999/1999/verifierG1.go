package main

import (
	"bufio"
	"fmt"
	"io"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

func commandFor(path string) *exec.Cmd {
	switch strings.ToLower(filepath.Ext(path)) {
	case ".go":
		return exec.Command("go", "run", path)
	case ".py":
		return exec.Command("python3", path)
	default:
		return exec.Command(path)
	}
}

func measured(y, x int) int {
	if y < x {
		return y
	}
	return y + 1
}

func interact(target string, values []int) error {
	cmd := commandFor(target)
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return fmt.Errorf("failed to open stdin pipe: %v", err)
	}
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("failed to open stdout pipe: %v", err)
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return fmt.Errorf("failed to open stderr pipe: %v", err)
	}

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start target: %v", err)
	}

	reader := bufio.NewReader(stdout)
	writer := bufio.NewWriter(stdin)
	stderrCh := make(chan string, 1)
	go func() {
		data, _ := io.ReadAll(stderr)
		stderrCh <- string(data)
	}()

	if _, err := fmt.Fprintln(writer, len(values)); err != nil {
		_ = cmd.Process.Kill()
		return fmt.Errorf("failed to send t: %v", err)
	}
	if err := writer.Flush(); err != nil {
		_ = cmd.Process.Kill()
		return fmt.Errorf("failed to flush t: %v", err)
	}

	for tc, x := range values {
		queries := 0
		for {
			var op string
			if _, err := fmt.Fscan(reader, &op); err != nil {
				_ = cmd.Process.Kill()
				if err == io.EOF {
					return fmt.Errorf("target terminated before answering test %d", tc+1)
				}
				return fmt.Errorf("failed to read operation on test %d: %v", tc+1, err)
			}

			switch op {
			case "?":
				var a, b int
				if _, err := fmt.Fscan(reader, &a, &b); err != nil {
					_ = cmd.Process.Kill()
					return fmt.Errorf("invalid query format on test %d: %v", tc+1, err)
				}
				if a < 1 || a > 1000 || b < 1 || b > 1000 {
					_ = cmd.Process.Kill()
					return fmt.Errorf("query out of bounds on test %d: ? %d %d", tc+1, a, b)
				}
				queries++
				if queries > 10 {
					_ = cmd.Process.Kill()
					return fmt.Errorf("too many queries on test %d: %d", tc+1, queries)
				}

				resp := measured(a, x) * measured(b, x)
				if _, err := fmt.Fprintln(writer, resp); err != nil {
					_ = cmd.Process.Kill()
					return fmt.Errorf("failed to send response on test %d: %v", tc+1, err)
				}
				if err := writer.Flush(); err != nil {
					_ = cmd.Process.Kill()
					return fmt.Errorf("failed to flush response on test %d: %v", tc+1, err)
				}
			case "!":
				var ans int
				if _, err := fmt.Fscan(reader, &ans); err != nil {
					_ = cmd.Process.Kill()
					return fmt.Errorf("invalid answer format on test %d: %v", tc+1, err)
				}
				if ans != x {
					_ = cmd.Process.Kill()
					return fmt.Errorf("wrong answer on test %d: expected %d, got %d", tc+1, x, ans)
				}
				goto nextTest
			default:
				_ = cmd.Process.Kill()
				return fmt.Errorf("unexpected token %q on test %d, expected '?' or '!'", op, tc+1)
			}
		}
	nextTest:
	}

	_ = stdin.Close()
	targetStderr := <-stderrCh
	if err := cmd.Wait(); err != nil {
		return fmt.Errorf("target exited with error: %v\nstderr:\n%s", err, strings.TrimSpace(string(targetStderr)))
	}
	return nil
}

func deterministicTests() []int {
	return []int{2, 4, 100, 500, 999, 678, 345, 876, 3, 998}
}

func randomTests() []int {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	vals := make([]int, 200)
	for i := range vals {
		vals[i] = rng.Intn(998) + 2
	}
	return vals
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierG1.go /path/to/candidate")
		os.Exit(1)
	}

	allTests := append(deterministicTests(), randomTests()...)
	if err := interact(os.Args[1], allTests); err != nil {
		fmt.Fprintf(os.Stderr, "verification failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("All tests passed:", strconv.Itoa(len(allTests)))
}
