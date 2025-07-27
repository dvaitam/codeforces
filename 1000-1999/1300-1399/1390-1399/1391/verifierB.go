package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

func buildRef() (string, error) {
	ref := "refB.bin"
	cmd := exec.Command("go", "build", "-o", ref, "1391B.go")
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
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	err := cmd.Run()
	return strings.TrimSpace(out.String()), errBuf.String(), err
}

type testCase struct{ input string }

func genCase(rng *rand.Rand) testCase {
	t := rng.Intn(3) + 1
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", t))
	for caseIdx := 0; caseIdx < t; caseIdx++ {
		n := rng.Intn(10) + 1
		m := rng.Intn(10) + 1
		sb.WriteString(fmt.Sprintf("%d %d\n", n, m))
		for i := 0; i < n; i++ {
			for j := 0; j < m; j++ {
				if rng.Intn(2) == 0 {
					sb.WriteByte('R')
				} else {
					sb.WriteByte('D')
				}
			}
			sb.WriteByte('\n')
		}
	}
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
		fmt.Println("usage: go run verifierB.go /path/to/binary")
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
