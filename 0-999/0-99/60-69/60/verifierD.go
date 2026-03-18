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

func buildRef() (string, error) {
	src := os.Getenv("REFERENCE_SOURCE_PATH")
	if src == "" {
		return "", fmt.Errorf("REFERENCE_SOURCE_PATH not set")
	}
	data, err := os.ReadFile(src)
	if err != nil {
		return "", fmt.Errorf("read reference: %v", err)
	}
	ref := "./refD.bin"
	if strings.Contains(string(data), "#include") {
		cppPath := "refD.cpp"
		if err := os.WriteFile(cppPath, data, 0644); err != nil {
			return "", fmt.Errorf("write cpp: %v", err)
		}
		cmd := exec.Command("g++", "-O2", "-o", ref, cppPath)
		if out, err := cmd.CombinedOutput(); err != nil {
			return "", fmt.Errorf("build reference cpp: %v: %s", err, string(out))
		}
	} else {
		cmd := exec.Command("go", "build", "-o", ref, src)
		if out, err := cmd.CombinedOutput(); err != nil {
			return "", fmt.Errorf("build reference: %v: %s", err, string(out))
		}
	}
	return ref, nil
}

func runBinary(path, input string) (string, error) {
	cmd := exec.Command(path)
	cmd.Stdin = strings.NewReader(input)
	var out, stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func generateCaseD(rng *rand.Rand) string {
	n := rng.Intn(10) + 1
	vals := make([]int, n)
	for i := 0; i < n; i++ {
		vals[i] = rng.Intn(40) + 1
	}
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", n)
	for i, v := range vals {
		if i > 0 {
			sb.WriteByte(' ')
		}
		fmt.Fprintf(&sb, "%d", v)
	}
	sb.WriteByte('\n')
	return sb.String()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: verifierD /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	ref, err := buildRef()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer os.Remove(ref)

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		tc := generateCaseD(rng)
		expect, err := runBinary(ref, tc)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: reference error: %v\ninput:\n%s", i+1, err, tc)
			os.Exit(1)
		}
		got, err := runBinary(bin, tc)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, tc)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != strings.TrimSpace(expect) {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, expect, got, tc)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
