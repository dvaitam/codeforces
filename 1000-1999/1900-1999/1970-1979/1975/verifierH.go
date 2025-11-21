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

const refSource = "1000-1999/1900-1999/1970-1979/1975/1975H.go"

func buildBinary(path string) (string, func(), error) {
	if !strings.HasSuffix(path, ".go") {
		return path, func() {}, nil
	}
	tmp, err := os.CreateTemp("", "cf-1975H-*")
	if err != nil {
		return "", func() {}, err
	}
	tmp.Close()
	os.Remove(tmp.Name())

	cmd := exec.Command("go", "build", "-o", tmp.Name(), path)
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

func generateRandomInput(rng *rand.Rand) string {
	t := rng.Intn(15) + 1
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", t))
	totalLen := 0
	for i := 0; i < t; i++ {
		var n int
		switch rng.Intn(5) {
		case 0:
			n = 1
		case 1:
			n = rng.Intn(5) + 1
		case 2:
			n = rng.Intn(20) + 1
		case 3:
			n = rng.Intn(100) + 1
		default:
			n = rng.Intn(1000) + 1
		}
		if totalLen+n > 30000 {
			n = 1
		}
		totalLen += n
		sb.WriteString(fmt.Sprintf("%d\n", n))
		sb.WriteString(randomString(rng, n))
		sb.WriteByte('\n')
	}
	return sb.String()
}

func randomString(rng *rand.Rand, n int) string {
	pattern := rng.Intn(4)
	buf := make([]byte, n)
	switch pattern {
	case 0: // all same char
		ch := byte('a' + rng.Intn(26))
		for i := range buf {
			buf[i] = ch
		}
	case 1: // skewed towards a max char
		maxCh := byte('a' + rng.Intn(26))
		for i := range buf {
			if rng.Intn(5) == 0 {
				buf[i] = maxCh
			} else {
				buf[i] = byte('a' + rng.Intn(int(maxCh-'a'+1)))
			}
		}
	case 2: // ascending run ending with max char
		maxCh := byte('a' + rng.Intn(26))
		for i := range buf {
			if i == n-1 {
				buf[i] = maxCh
				continue
			}
			buf[i] = byte('a' + rng.Intn(int(maxCh-'a'+1)))
		}
	default: // fully random
		for i := range buf {
			buf[i] = byte('a' + rng.Intn(26))
		}
	}
	return string(buf)
}

func fixedTests() []string {
	samples := []string{
		"5\n1\na\n1\nz\n2\nab\n3\nqaq\n5\nabczz\n",
		"3\n4\nbbbb\n6\naaaaaa\n6\nzzzzzz\n",
		"4\n5\nazzzz\n3\nzzz\n6\nzaaaaa\n7\nbzbzbzb\n",
	}
	return samples
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierH.go /path/to/binary")
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
	for i := 0; i < 50; i++ {
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
