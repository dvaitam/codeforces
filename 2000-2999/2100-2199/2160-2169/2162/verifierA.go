package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

func buildReference() (string, error) {
	path := "./2162A_ref.bin"
	cmd := exec.Command("go", "build", "-o", path, "2162A.go")
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("failed to build reference: %v\n%s", err, out)
	}
	return path, nil
}

func runProgram(target, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(target, ".go") {
		cmd = exec.Command("go", "run", target)
	} else {
		cmd = exec.Command(target)
	}
	cmd.Stdin = strings.NewReader(input)
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("%v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(stdout.String()), nil
}

func buildFromSlices(cases [][]int) string {
	var sb strings.Builder
	sb.WriteString(strconv.Itoa(len(cases)))
	sb.WriteByte('\n')
	for _, arr := range cases {
		sb.WriteString(strconv.Itoa(len(arr)))
		sb.WriteByte('\n')
		for i, v := range arr {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.Itoa(v))
		}
		sb.WriteByte('\n')
	}
	return sb.String()
}

func randomInput(rng *rand.Rand, t int) string {
	var sb strings.Builder
	sb.WriteString(strconv.Itoa(t))
	sb.WriteByte('\n')
	for i := 0; i < t; i++ {
		n := rng.Intn(10) + 1
		sb.WriteString(strconv.Itoa(n))
		sb.WriteByte('\n')
		for j := 0; j < n; j++ {
			if j > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.Itoa(rng.Intn(10) + 1))
		}
		sb.WriteByte('\n')
	}
	return sb.String()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	var refPath string
	fail := func(format string, args ...interface{}) {
		if refPath != "" {
			_ = os.Remove(refPath)
		}
		fmt.Fprintf(os.Stderr, format+"\n", args...)
		os.Exit(1)
	}

	var err error
	refPath, err = buildReference()
	if err != nil {
		fail("%v", err)
	}
	defer os.Remove(refPath)

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	tests := []string{}
	tests = append(tests, buildFromSlices([][]int{{1}}))
	tests = append(tests, buildFromSlices([][]int{{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}}))
	tests = append(tests, buildFromSlices([][]int{{10}, {1, 1, 1, 1, 1, 1, 1, 1, 1, 1}}))

	for i := 0; i < 200; i++ {
		t := rng.Intn(20) + 1
		tests = append(tests, randomInput(rng, t))
	}
	tests = append(tests, randomInput(rng, 10000))

	for idx, input := range tests {
		expect, err := runProgram(refPath, input)
		if err != nil {
			fail("reference failed on test %d: %v", idx+1, err)
		}
		got, err := runProgram(bin, input)
		if err != nil {
			fail("test %d: runtime error: %v", idx+1, err)
		}
		if expect != got {
			fail("test %d failed:\ninput:\n%s\nexpected:\n%s\ngot:\n%s", idx+1, input, expect, got)
		}
	}
	fmt.Println("All tests passed")
}
