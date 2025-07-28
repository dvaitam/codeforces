package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

type testCaseA []int

func genTestsA() []testCaseA {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	tests := make([]testCaseA, 100)
	for i := range tests {
		n := rng.Intn(8) + 2 // 2..9
		p := rng.Perm(n)
		for j := range p {
			p[j]++
		}
		tests[i] = append(testCaseA(nil), p...)
	}
	return tests
}

func isLocalMin(p []int, idx int) bool {
	n := len(p)
	val := p[idx]
	left := int(1e9)
	right := int(1e9)
	if idx > 0 {
		left = p[idx-1]
	}
	if idx+1 < n {
		right = p[idx+1]
	}
	return val < left && val < right
}

func runCase(bin string, perm []int) error {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return err
	}
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}
	cmd.Stderr = os.Stderr
	if err := cmd.Start(); err != nil {
		return err
	}
	inW := bufio.NewWriter(stdin)
	outR := bufio.NewReader(stdout)
	fmt.Fprintln(inW, len(perm))
	inW.Flush()

	queries := 0
	for {
		line, err := outR.ReadString('\n')
		if err != nil {
			cmd.Process.Kill()
			return fmt.Errorf("read error: %v", err)
		}
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "?") {
			queries++
			if queries > 100 {
				cmd.Process.Kill()
				return fmt.Errorf("too many queries")
			}
			parts := strings.Fields(line)
			if len(parts) != 2 {
				cmd.Process.Kill()
				return fmt.Errorf("invalid query: %s", line)
			}
			idx, err := strconv.Atoi(parts[1])
			if err != nil || idx < 1 || idx > len(perm) {
				cmd.Process.Kill()
				return fmt.Errorf("query out of range")
			}
			fmt.Fprintln(inW, perm[idx-1])
			inW.Flush()
		} else if strings.HasPrefix(line, "!") {
			parts := strings.Fields(line)
			if len(parts) != 2 {
				cmd.Process.Kill()
				return fmt.Errorf("invalid answer: %s", line)
			}
			idx, err := strconv.Atoi(parts[1])
			if err != nil || idx < 1 || idx > len(perm) {
				cmd.Process.Kill()
				return fmt.Errorf("answer out of range")
			}
			if !isLocalMin(perm, idx-1) {
				cmd.Process.Kill()
				return fmt.Errorf("wrong answer: %d not local minimum", idx)
			}
			stdin.Close()
			cmd.Wait()
			return nil
		} else {
			cmd.Process.Kill()
			return fmt.Errorf("invalid output: %s", line)
		}
	}
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := genTestsA()
	for i, p := range tests {
		if err := runCase(bin, p); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}
