package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/candidate")
		os.Exit(1)
	}
	candidate := os.Args[1]

	refPath := os.Getenv("REFERENCE_SOURCE_PATH")
	if refPath == "" {
		fmt.Fprintln(os.Stderr, "REFERENCE_SOURCE_PATH not set")
		os.Exit(1)
	}

	refBin, err := buildBinary(refPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to build reference: %v\n", err)
		os.Exit(1)
	}
	defer os.Remove(refBin)

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 50; i++ {
		input := genTest(rng)
		refOut, err := runProgram(refBin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference failed on test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		candOut, err := runCandidate(candidate, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate failed on test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(refOut) != strings.TrimSpace(candOut) {
			fmt.Fprintf(os.Stderr, "test %d mismatch\ninput:\n%s\nreference:\n%s\ncandidate:\n%s\n",
				i+1, input, refOut, candOut)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed.")
}

func buildBinary(src string) (string, error) {
	tmp, err := os.CreateTemp("", "1903E-ref-*")
	if err != nil {
		return "", err
	}
	tmp.Close()
	cmd := exec.Command("go", "build", "-o", tmp.Name(), src)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		os.Remove(tmp.Name())
		return "", fmt.Errorf("%v\n%s", err, out.String())
	}
	return tmp.Name(), nil
}

func runProgram(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	return out.String(), nil
}

func runCandidate(path, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(path, ".go") {
		cmd = exec.Command("go", "run", path)
	} else {
		cmd = exec.Command(path)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	return out.String(), nil
}

func genTest(rng *rand.Rand) string {
	t := rng.Intn(5) + 1
	var sb strings.Builder
	fmt.Fprintln(&sb, t)
	for i := 0; i < t; i++ {
		n := rng.Intn(8) + 1
		fmt.Fprintln(&sb, n)
		for j := 0; j < n; j++ {
			if j > 0 {
				sb.WriteByte(' ')
			}
			fmt.Fprintf(&sb, "%d", rng.Intn(100)+1)
		}
		sb.WriteByte('\n')
	}
	return sb.String()
}
