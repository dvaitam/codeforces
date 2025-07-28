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
	ref := "./refE.bin"
	cmd := exec.Command("go", "build", "-o", ref, "1603E.go")
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("failed to build reference: %v\n%s", err, out)
	}
	return ref, nil
}

var primes = []int{1000000007, 1000000009, 1000000033, 1000000087, 1000000093, 1000000097}

func generateCase(rng *rand.Rand) string {
	n := rng.Intn(20) + 1
	M := primes[rng.Intn(len(primes))]
	return fmt.Sprintf("%d %d\n", n, M)
}

func runCmd(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		tmp, err := os.CreateTemp("", "candidate-*")
		if err != nil {
			return "", err
		}
		tmp.Close()
		exe := tmp.Name()
		build := exec.Command("go", "build", "-o", exe, bin)
		if out, err := build.CombinedOutput(); err != nil {
			os.Remove(exe)
			return "", fmt.Errorf("build error: %v\n%s", err, out)
		}
		defer os.Remove(exe)
		cmd = exec.Command(exe)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func runCase(exe, ref, input string) error {
	expected, err := runCmd(ref, input)
	if err != nil {
		return err
	}
	got, err := runCmd(exe, input)
	if err != nil {
		return err
	}
	if got != expected {
		return fmt.Errorf("expected %s got %s", expected, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
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
