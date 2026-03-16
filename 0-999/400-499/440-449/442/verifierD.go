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

func buildReference() (string, error) {
	refPath := os.Getenv("REFERENCE_SOURCE_PATH")
	if refPath == "" {
		return "", fmt.Errorf("REFERENCE_SOURCE_PATH not set")
	}
	outBin := filepath.Join(os.TempDir(), "ref442D")
	content, err := os.ReadFile(refPath)
	if err != nil {
		return "", fmt.Errorf("read reference: %v", err)
	}
	if strings.Contains(string(content), "#include") {
		cppPath := filepath.Join(os.TempDir(), "ref442D.cpp")
		if err := os.WriteFile(cppPath, content, 0644); err != nil {
			return "", err
		}
		cmd := exec.Command("g++", "-O2", "-o", outBin, cppPath)
		if o, err := cmd.CombinedOutput(); err != nil {
			return "", fmt.Errorf("build ref (c++) failed: %v\n%s", err, o)
		}
	} else {
		cmd := exec.Command("go", "build", "-o", outBin, refPath)
		if o, err := cmd.CombinedOutput(); err != nil {
			return "", fmt.Errorf("build ref failed: %v\n%s", err, o)
		}
	}
	return outBin, nil
}

func runBinary(path, input string) (string, error) {
	cmd := exec.Command(path)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return out.String(), nil
}

func generateCase(rng *rand.Rand) string {
	n := rng.Intn(20) + 1
	var b strings.Builder
	fmt.Fprintf(&b, "%d\n", n)
	for i := 1; i <= n; i++ {
		if i > 1 {
			b.WriteByte(' ')
		}
		fmt.Fprintf(&b, "%d", rng.Intn(i)+1)
	}
	b.WriteByte('\n')
	return b.String()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	refBin, err := buildReference()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to build reference: %v\n", err)
		os.Exit(1)
	}
	defer os.Remove(refBin)

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		tc := generateCase(rng)
		expect, err := runBinary(refBin, tc)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: ref error: %v\ninput:\n%s", i+1, err, tc)
			os.Exit(1)
		}
		out, err := runBinary(bin, tc)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, tc)
			os.Exit(1)
		}
		if strings.TrimSpace(out) != strings.TrimSpace(expect) {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %sinput:\n%s", i+1, expect, out, tc)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
