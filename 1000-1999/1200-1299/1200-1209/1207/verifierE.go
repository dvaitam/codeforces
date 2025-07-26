package main

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

func buildIfGo(path string) (string, func(), error) {
	if strings.HasSuffix(path, ".go") {
		tmp, err := os.CreateTemp("", "solbin*")
		if err != nil {
			return "", nil, err
		}
		tmp.Close()
		out, err := exec.Command("go", "build", "-o", tmp.Name(), path).CombinedOutput()
		if err != nil {
			os.Remove(tmp.Name())
			return "", nil, fmt.Errorf("build failed: %v\n%s", err, out)
		}
		return tmp.Name(), func() { os.Remove(tmp.Name()) }, nil
	}
	return path, func() {}, nil
}

func runCase(bin string, x int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, bin)
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return err
	}
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Start(); err != nil {
		return err
	}
	inw := bufio.NewWriter(stdin)
	outr := bufio.NewReader(stdout)
	queries := 0
	for {
		line, err := outr.ReadString('\n')
		if err != nil {
			return fmt.Errorf("read error: %v stderr:%s", err, stderr.String())
		}
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "?") {
			queries++
			fields := strings.Fields(line)
			if len(fields) != 101 {
				return fmt.Errorf("bad query")
			}
			nums := make([]int, 100)
			for i := 0; i < 100; i++ {
				v, _ := strconv.Atoi(fields[i+1])
				nums[i] = v
			}
			resp := nums[0] ^ x
			fmt.Fprintln(inw, resp)
			inw.Flush()
			if queries == 2 {
				continue
			}
		} else if strings.HasPrefix(line, "!") {
			valStr := strings.TrimSpace(line[1:])
			val, _ := strconv.Atoi(valStr)
			stdin.Close()
			err := cmd.Wait()
			if err != nil {
				return fmt.Errorf("runtime error: %v stderr:%s", err, stderr.String())
			}
			if val != x {
				return fmt.Errorf("wrong answer: expected %d got %d", x, val)
			}
			return nil
		}
	}
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		return
	}
	bin, cleanup, err := buildIfGo(os.Args[1])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer cleanup()
	for i := 0; i < 100; i++ {
		if err := runCase(bin, i); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
