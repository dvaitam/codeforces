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
	refSrc := os.Getenv("REFERENCE_SOURCE_PATH")
	if refSrc == "" {
		return "", fmt.Errorf("REFERENCE_SOURCE_PATH not set")
	}
	content, err := os.ReadFile(refSrc)
	if err != nil {
		return "", fmt.Errorf("read reference source: %v", err)
	}
	bin := os.TempDir() + "/1945F_ref.bin"
	if strings.Contains(string(content), "#include") {
		cppSrc := os.TempDir() + "/1945F_ref.cpp"
		if err := os.WriteFile(cppSrc, content, 0644); err != nil {
			return "", fmt.Errorf("write cpp source: %v", err)
		}
		cmd := exec.Command("g++", "-O2", "-o", bin, cppSrc)
		if out, err := cmd.CombinedOutput(); err != nil {
			return "", fmt.Errorf("build c++ reference failed: %v\n%s", err, string(out))
		}
	} else {
		cmd := exec.Command("go", "build", "-o", bin, refSrc)
		if out, err := cmd.CombinedOutput(); err != nil {
			return "", fmt.Errorf("build reference failed: %v\n%s", err, string(out))
		}
	}
	return bin, nil
}

func generatePerm(rng *rand.Rand, n int) []int {
	perm := rng.Perm(n)
	for i := range perm {
		perm[i]++
	}
	return perm
}

func generateCase(rng *rand.Rand) string {
	n := rng.Intn(20) + 1
	v := make([]int, n)
	for i := range v {
		v[i] = rng.Intn(100) + 1
	}
	p := generatePerm(rng, n)
	var sb strings.Builder
	sb.WriteString("1\n")
	fmt.Fprintf(&sb, "%d\n", n)
	for i, val := range v {
		if i > 0 {
			sb.WriteByte(' ')
		}
		fmt.Fprintf(&sb, "%d", val)
	}
	sb.WriteByte('\n')
	for i, val := range p {
		if i > 0 {
			sb.WriteByte(' ')
		}
		fmt.Fprintf(&sb, "%d", val)
	}
	sb.WriteByte('\n')
	return sb.String()
}

func runCase(exe, ref, input string) error {
	cmdRef := exec.Command(ref)
	cmdRef.Stdin = strings.NewReader(input)
	var refOut bytes.Buffer
	cmdRef.Stdout = &refOut
	cmdRef.Stderr = &refOut
	if err := cmdRef.Run(); err != nil {
		return fmt.Errorf("reference runtime error: %v\n%s", err, refOut.String())
	}
	expected := strings.TrimSpace(refOut.String())

	cmd := exec.Command(exe)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(out.String())
	if got != expected {
		return fmt.Errorf("expected %s got %s", expected, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: verifierF /path/to/binary")
		os.Exit(1)
	}
	exe := os.Args[1]
	ref, err := buildRef()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer os.Remove(ref)
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in := generateCase(rng)
		if err := runCase(exe, ref, in); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
