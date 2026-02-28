package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

const refSource = "2064C.go"

func buildBinary(path string) (string, func(), error) {
	if !strings.HasSuffix(path, ".go") {
		return path, func() {}, nil
	}
	tmp, err := os.CreateTemp("", "cf-2064C-*")
	if err != nil {
		return "", func() {}, err
	}
	tmp.Close()
	os.Remove(tmp.Name())

	cmd := exec.Command("go", "build", "-o", tmp.Name(), filepath.Base(path))
	cmd.Dir = filepath.Dir(path)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		os.Remove(tmp.Name())
		return "", func() {}, fmt.Errorf("go build error: %v\n%s", err, stderr.String())
	}
	return tmp.Name(), func() { os.Remove(tmp.Name()) }, nil
}

func runProgram(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(stdout.String()), nil
}

func fixedTests() []string {
	return []string{
		"3\n6\n3 1 4 -1 -5 -9\n6\n-10 -3 -17 1 19 20\n1\n1\n",
		"4\n1\n-7\n2\n5 -5\n3\n1 -2 3\n5\n-1 -2 -3 -4 -5\n",
	}
}

func generateRandomInput(rng *rand.Rand) string {
	t := rng.Intn(20) + 1
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", t))
	totalN := 0
	for i := 0; i < t; i++ {
		remaining := 200000 - totalN
		if remaining <= 0 {
			sb.WriteString("1\n1\n")
			totalN++
			continue
		}
		var n int
		switch rng.Intn(6) {
		case 0:
			n = 1
		case 1:
			n = rng.Intn(5) + 1
		case 2:
			n = rng.Intn(20) + 1
		case 3:
			n = rng.Intn(200) + 1
		case 4:
			n = rng.Intn(2000) + 1
		default:
			n = rng.Intn(50000) + 1
		}
		if n > remaining {
			n = remaining
		}
		totalN += n
		sb.WriteString(fmt.Sprintf("%d\n", n))
		for j := 0; j < n; j++ {
			if j > 0 {
				sb.WriteByte(' ')
			}
			val := rng.Int63n(2_000_000_000) - 1_000_000_000
			if val >= 0 {
				val++
			}
			sb.WriteString(fmt.Sprintf("%d", val))
		}
		sb.WriteByte('\n')
	}
	return sb.String()
}

func bigEdgeCase(rng *rand.Rand) string {
	n := 200000
	var sb strings.Builder
	sb.WriteString("1\n")
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for i := 0; i < n; i++ {
		if i > 0 {
			sb.WriteByte(' ')
		}
		val := int64(1_000_000_000)
		if rng.Intn(2) == 0 {
			val = -val
		}
		sb.WriteString(fmt.Sprintf("%d", val))
	}
	sb.WriteByte('\n')
	return sb.String()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	candPath, err := filepath.Abs(os.Args[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to resolve candidate path: %v\n", err)
		os.Exit(1)
	}
	refPath, err := filepath.Abs(refSource)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to resolve reference path: %v\n", err)
		os.Exit(1)
	}

	refBin, refCleanup, err := buildBinary(refPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to build reference: %v\n", err)
		os.Exit(1)
	}
	defer refCleanup()

	candBin, candCleanup, err := buildBinary(candPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to build candidate: %v\n", err)
		os.Exit(1)
	}
	defer candCleanup()

	tests := fixedTests()
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	tests = append(tests, bigEdgeCase(rng))
	for i := 0; i < 30; i++ {
		tests = append(tests, generateRandomInput(rng))
	}

	for idx, input := range tests {
		exp, err := runProgram(refBin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference failed on case %d: %v\ninput:\n%s", idx+1, err, input)
			os.Exit(1)
		}
		got, err := runProgram(candBin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate failed on case %d: %v\ninput:\n%s", idx+1, err, input)
			os.Exit(1)
		}
		if exp != got {
			fmt.Fprintf(os.Stderr, "mismatch on case %d\nexpected:\n%s\n\ngot:\n%s\ninput:\n%s", idx+1, exp, got, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
