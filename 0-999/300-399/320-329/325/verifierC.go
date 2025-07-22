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

func buildRef() (string, error) {
	refBin := filepath.Join(os.TempDir(), "ref325C")
	cmd := exec.Command("go", "build", "-o", refBin, "325C.go")
	var out bytes.Buffer
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("failed to build reference: %v\n%s", err, out.String())
	}
	return refBin, nil
}

func runRef(bin string, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("ref runtime error: %v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func generateCase(rng *rand.Rand) string {
	m := rng.Intn(10) + 1
	n := rng.Intn(8) + 1
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", m, n))
	for i := 0; i < m; i++ {
		mi := rng.Intn(n) + 1
		li := rng.Intn(3) + 1
		sb.WriteString(fmt.Sprintf("%d %d", mi, li))
		hasDiam := false
		for j := 0; j < li; j++ {
			if rng.Intn(3) == 0 {
				sb.WriteString(" -1")
				hasDiam = true
			} else {
				sb.WriteString(fmt.Sprintf(" %d", rng.Intn(n)+1))
			}
		}
		if !hasDiam {
			sb.WriteString(" -1")
		}
		sb.WriteByte('\n')
	}
	return sb.String()
}

func runCase(bin, refBin, input string) error {
	exp, err := runRef(refBin, input)
	if err != nil {
		return err
	}
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(out.String())
	if got != exp {
		return fmt.Errorf("expected:\n%s\n got:\n%s", exp, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	refBin, err := buildRef()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer os.Remove(refBin)
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		input := generateCase(rng)
		if err := runCase(bin, refBin, input); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
