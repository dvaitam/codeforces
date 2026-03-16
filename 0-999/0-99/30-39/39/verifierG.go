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

var refBinPath string

func buildRef() error {
	refPath := os.Getenv("REFERENCE_SOURCE_PATH")
	if refPath == "" {
		return fmt.Errorf("REFERENCE_SOURCE_PATH not set")
	}
	outBin := filepath.Join(os.TempDir(), "ref39G")
	content, err := os.ReadFile(refPath)
	if err != nil {
		return fmt.Errorf("read reference: %v", err)
	}
	if strings.Contains(string(content), "#include") {
		cppPath := filepath.Join(os.TempDir(), "ref39G.cpp")
		if err := os.WriteFile(cppPath, content, 0644); err != nil {
			return err
		}
		cmd := exec.Command("g++", "-O2", "-o", outBin, cppPath)
		if o, err := cmd.CombinedOutput(); err != nil {
			return fmt.Errorf("build ref (c++) failed: %v\n%s", err, o)
		}
	} else {
		cmd := exec.Command("go", "build", "-o", outBin, refPath)
		if o, err := cmd.CombinedOutput(); err != nil {
			return fmt.Errorf("build ref failed: %v\n%s", err, o)
		}
	}
	refBinPath = outBin
	return nil
}

func solveG(input string) string {
	cmd := exec.Command(refBinPath)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	if err != nil {
		return ""
	}
	return out.String()
}

func generateCaseG(rng *rand.Rand) string {
	typ := rng.Intn(3)
	switch typ {
	case 0:
		target := rng.Intn(32768)
		return fmt.Sprintf("%d\nreturn n;\n", target)
	case 1:
		c := rng.Intn(10) + 1
		n := rng.Intn(32768)
		var target int
		if n > c {
			target = (n - c) % 32768
		} else {
			target = n
		}
		code := fmt.Sprintf("if (n>%d) return n-%d; return n;", c, c)
		return fmt.Sprintf("%d\n%s\n", target, code)
	default:
		c := rng.Intn(10) + 1
		n := rng.Intn(32768)
		var target int
		if n < c {
			target = n
		} else {
			target = n - c
		}
		code := fmt.Sprintf("if (n<%d) return n; return n-%d;", c, c)
		return fmt.Sprintf("%d\n%s\n", target, code)
	}
}

func runBinary(path, input string) (string, error) {
	cmd := exec.Command(path)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierG.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	if err := buildRef(); err != nil {
		fmt.Fprintf(os.Stderr, "failed to build reference: %v\n", err)
		os.Exit(1)
	}
	defer os.Remove(refBinPath)
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	cases := make([]string, 100)
	for i := 0; i < 100; i++ {
		cases[i] = generateCaseG(rng)
	}
	for i, tc := range cases {
		expect := solveG(tc)
		got, err := runBinary(bin, tc)
		if err != nil {
			fmt.Printf("case %d: runtime error: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != strings.TrimSpace(expect) {
			fmt.Printf("case %d failed\ninput:\n%sexpected:%sq\ngot:%sq\n", i+1, tc, expect, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
