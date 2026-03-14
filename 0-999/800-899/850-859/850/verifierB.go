package main

import (
	"bytes"
	"context"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

type Test struct {
	input string
}

const timeout = 2 * time.Minute

func runExe(path string, input string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	var cmd *exec.Cmd
	if strings.HasSuffix(path, ".go") {
		cmd = exec.CommandContext(ctx, "go", "run", path)
	} else {
		cmd = exec.CommandContext(ctx, path)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

func buildRef() (string, error) {
	srcPath := os.Getenv("REFERENCE_SOURCE_PATH")
	if srcPath == "" {
		srcPath = "850B.go"
	}
	ref, err := filepath.Abs("./refB.bin")
	if err != nil {
		return "", err
	}
	ext := filepath.Ext(srcPath)
	var cmd *exec.Cmd
	switch ext {
	case ".go":
		cmd = exec.Command("go", "build", "-o", ref, srcPath)
	case ".cpp", ".cc", ".cxx":
		cmd = exec.Command("g++", "-O2", "-o", ref, srcPath)
	case ".c":
		cmd = exec.Command("gcc", "-O2", "-o", ref, srcPath)
	default:
		return "", fmt.Errorf("unsupported reference source extension: %s", ext)
	}
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("build reference failed: %v: %s", err, string(out))
	}
	return ref, nil
}

func genTests() []Test {
	rand.Seed(0)
	tests := make([]Test, 0, 101)
	for i := 0; i < 100; i++ {
		n := rand.Intn(8) + 1
		x := rand.Intn(10) + 1
		y := rand.Intn(10) + 1
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d %d %d\n", n, x, y))
		for j := 0; j < n; j++ {
			if j > 0 {
				sb.WriteByte(' ')
			}
			val := rand.Intn(20) + 1
			sb.WriteString(strconv.Itoa(val))
		}
		sb.WriteByte('\n')
		tests = append(tests, Test{sb.String()})
	}
	tests = append(tests, Test{"1 1 1\n5\n"})
	return tests
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: verifierB /path/to/candidate")
		return
	}
	bin := os.Args[1]
	ref, err := buildRef()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer os.Remove(ref)

	tests := genTests()
	for i, tc := range tests {
		exp, err := runExe(ref, tc.input)
		if err != nil {
			fmt.Printf("reference runtime error on test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		got, err := runExe(bin, tc.input)
		if err != nil {
			fmt.Printf("candidate runtime error on test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(exp) != strings.TrimSpace(got) {
			fmt.Printf("Test %d failed\nInput:\n%sExpected:\n%sGot:\n%s\n", i+1, tc.input, exp, got)
			os.Exit(1)
		}
	}
	fmt.Println("all tests passed")
}
