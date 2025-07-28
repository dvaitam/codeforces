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
	ref := "./refF.bin"
	cmd := exec.Command("go", "build", "-o", ref, "1603F.go")
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("failed to build reference: %v\n%s", err, out)
	}
	return ref, nil
}

func generateCase(rng *rand.Rand) string {
	n := rng.Int63n(1000) + 1
	k := rng.Int63n(30)
	var xbound int64
	if k < 20 {
		xbound = int64(1) << k
	} else {
		xbound = int64(1) << 20
	}
	var x int64
	if xbound > 0 {
		x = rng.Int63n(xbound)
	}
	return fmt.Sprintf("%d %d %d\n", n, k, x)
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
	input = "1\n" + input
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
		fmt.Fprintln(os.Stderr, "usage: go run verifierF.go /path/to/binary")
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
