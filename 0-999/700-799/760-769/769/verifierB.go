package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/candidate")
		os.Exit(1)
	}
	candidate := os.Args[1]

	refBin, err := buildReference()
	if err != nil {
		fail("failed to build reference: %v", err)
	}
	defer os.Remove(refBin)

	rand.Seed(time.Now().UnixNano())

	for i := 0; i < 50; i++ {
		n := rand.Intn(99) + 2 // 2 to 100
		a := make([]int, n)
		for j := 0; j < n; j++ {
			a[j] = rand.Intn(101) // 0 to 100
		}

		input := fmt.Sprintf("%d\n", n)
		for j, val := range a {
			if j > 0 {
				input += " "
			}
			input += fmt.Sprintf("%d", val)
		}
		input += "\n"

		refOut, err := runProgram(refBin, input)
		if err != nil {
			fail("reference failed on test %d: %v", i+1, err)
		}

		candOut, err := runProgram(candidate, input)
		if err != nil {
			fail("candidate failed on test %d: %v", i+1, err)
		}

		if err := verify(n, a, refOut, candOut); err != nil {
			fmt.Printf("Input:\n%s\n", input)
			fmt.Printf("Reference Output:\n%s\n", refOut)
			fmt.Printf("Candidate Output:\n%s\n", candOut)
			fail("test %d failed: %v", i+1, err)
		}
	}
	fmt.Println("All 50 tests passed")
}

func verify(n int, a []int, refOut, candOut string) error {
	refPossible := strings.TrimSpace(refOut) != "-1"
	candPossible := strings.TrimSpace(candOut) != "-1"

	if refPossible != candPossible {
		return fmt.Errorf("possibility mismatch: reference says %v, candidate says %v", refPossible, candPossible)
	}

	if !candPossible {
		return nil
	}

	scanner := bufio.NewScanner(strings.NewReader(candOut))
	if !scanner.Scan() {
		return fmt.Errorf("empty candidate output")
	}
	kStr := strings.TrimSpace(scanner.Text())
	k, err := strconv.Atoi(kStr)
	if err != nil {
		return fmt.Errorf("invalid first line (expected integer k): %v", err)
	}

	informed := make(map[int]bool)
	informed[1] = true
	
	caps := make([]int, n+1)
	for i := 0; i < n; i++ {
		caps[i+1] = a[i]
	}

	count := 0
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" { continue }
		parts := strings.Fields(line)
		if len(parts) != 2 {
			return fmt.Errorf("invalid line format (expected 'u v'): %s", line)
		}
		u, err1 := strconv.Atoi(parts[0])
		v, err2 := strconv.Atoi(parts[1])
		if err1 != nil || err2 != nil {
			return fmt.Errorf("invalid student numbers: %s", line)
		}

		if u < 1 || u > n || v < 1 || v > n {
			return fmt.Errorf("student number out of range 1..%d: %d->%d", n, u, v)
		}

		if !informed[u] {
			return fmt.Errorf("student %d sends message but doesn't know news", u)
		}
		if caps[u] <= 0 {
			return fmt.Errorf("student %d has no capacity left", u)
		}
		caps[u]--
		informed[v] = true
		count++
	}

	if count != k {
		return fmt.Errorf("declared %d messages, but found %d lines", k, count)
	}

	if len(informed) != n {
		return fmt.Errorf("only %d students informed, expected %d", len(informed), n)
	}

	return nil
}

func buildReference() (string, error) {
	_, currentFile, _, ok := runtime.Caller(0)
	if !ok {
		return "", fmt.Errorf("failed to determine current file path")
	}
	dir := filepath.Dir(currentFile)
	refSource := filepath.Join(dir, "769B.go")

	tmp, err := os.CreateTemp("", "769B-ref-*")
	if err != nil {
		return "", err
	}
	tmp.Close()

	cmd := exec.Command("go", "build", "-o", tmp.Name(), refSource)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		os.Remove(tmp.Name())
		return "", fmt.Errorf("%v\n%s", err, out.String())
	}
	return tmp.Name(), nil
}

func runProgram(path string, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(path, ".go") {
		cmd = exec.Command("go", "run", path)
	} else {
		cmd = exec.Command(path)
	}
	cmd.Stdin = strings.NewReader(input)
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("%v\n%s", err, stderr.String())
	}
	return stdout.String(), nil
}

func fail(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format+"\n", args...)
	os.Exit(1)
}