package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func buildRef() (string, error) {
	ref := "refD.bin"
	cmd := exec.Command("go", "build", "-o", ref, "999D.go")
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("failed to build reference: %v\n%s", err, out)
	}
	return ref, nil
}

func runBinary(bin, input string) (string, string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err := cmd.Run()
	return strings.TrimSpace(out.String()), stderr.String(), err
}

type testCase struct{ input string }

func genCase(rng *rand.Rand) testCase {
	m := rng.Intn(10) + 1
	n := (rng.Intn(10) + 1) * m
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", n, m))
	for i := 0; i < n; i++ {
		sb.WriteString(strconv.Itoa(rng.Intn(20) - 10))
		if i+1 < n {
			sb.WriteByte(' ')
		}
	}
	sb.WriteByte('\n')
	return testCase{sb.String()}
}

func runCase(bin, ref string, tc testCase, idx int) error {
	expOut, expErr, err := runBinary(ref, tc.input)
	if err != nil {
		return fmt.Errorf("reference failed: %v\nstderr: %s", err, expErr)
	}
	gotOut, gotErr, err := runBinary(bin, tc.input)
	if err != nil {
		return fmt.Errorf("test %d: runtime error: %v\nstderr: %s", idx, err, gotErr)
	}
	if strings.TrimSpace(gotOut) != strings.TrimSpace(expOut) {
		return fmt.Errorf("test %d failed\nexpected: %s\n got: %s", idx, expOut, gotOut)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	ref, err := buildRef()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer os.Remove(ref)

	rng := rand.New(rand.NewSource(99))
	for i := 1; i <= 100; i++ {
		tc := genCase(rng)
		if err := runCase(bin, ref, tc, i); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
