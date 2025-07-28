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
	exe := "./refB.bin"
	cmd := exec.Command("go", "build", "-o", exe, "1631B.go")
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("failed to build reference: %v\n%s", err, out)
	}
	return exe, nil
}

func runCmd(path, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(path, ".go") {
		tmp, err := os.CreateTemp("", "candidate-*")
		if err != nil {
			return "", err
		}
		tmp.Close()
		exe := tmp.Name()
		build := exec.Command("go", "build", "-o", exe, path)
		if out, err := build.CombinedOutput(); err != nil {
			os.Remove(exe)
			return "", fmt.Errorf("build error: %v\n%s", err, out)
		}
		defer os.Remove(exe)
		cmd = exec.Command(exe)
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
	return strings.TrimSpace(out.String()), nil
}

func generateCase(rng *rand.Rand) string {
	n := rng.Intn(200) + 1
	var sb strings.Builder
	sb.WriteString("1\n")
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for i := 0; i < n; i++ {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", rng.Intn(n)+1))
	}
	sb.WriteByte('\n')
	return sb.String()
}

func runCase(candidate, reference, input string) error {
	want, err := runCmd(reference, input)
	if err != nil {
		return fmt.Errorf("reference failed: %v", err)
	}
	got, err := runCmd(candidate, input)
	if err != nil {
		return err
	}
	if want != got {
		return fmt.Errorf("expected %q got %q", want, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	candidate := os.Args[1]

	ref, err := buildRef()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer os.Remove(ref)

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	tests := []string{
		"1\n1\n1\n",
		"1\n2\n1 1\n",
	}
	for len(tests) < 100 {
		tests = append(tests, generateCase(rng))
	}

	for i, tc := range tests {
		if err := runCase(candidate, ref, tc); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, tc)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
