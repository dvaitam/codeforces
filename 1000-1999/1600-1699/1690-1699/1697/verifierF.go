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

type testCase struct {
	input string
}

func buildRef() (string, error) {
	ref := "./refF.bin"
	cmd := exec.Command("go", "build", "-o", ref, "1697F.go")
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("failed to build reference: %v\n%s", err, out)
	}
	return ref, nil
}

func runBinary(bin, input string) (string, error) {
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
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func genCase(rng *rand.Rand) testCase {
	n := rng.Intn(3) + 1
	m := rng.Intn(3) + 1
	K := rng.Intn(3) + 2
	var sb strings.Builder
	sb.WriteString("1\n")
	sb.WriteString(fmt.Sprintf("%d %d %d\n", n, m, K))
	for i := 0; i < m; i++ {
		typ := rng.Intn(3) + 1
		if typ == 1 {
			idx := rng.Intn(n) + 1
			v := rng.Intn(K) + 1
			sb.WriteString(fmt.Sprintf("1 %d %d\n", idx, v))
		} else if typ == 2 {
			a := rng.Intn(n) + 1
			b := rng.Intn(n) + 1
			v := rng.Intn(2*K) + 1
			sb.WriteString(fmt.Sprintf("2 %d %d %d\n", a, b, v))
		} else {
			a := rng.Intn(n) + 1
			b := rng.Intn(n) + 1
			v := rng.Intn(2*K) + 1
			sb.WriteString(fmt.Sprintf("3 %d %d %d\n", a, b, v))
		}
	}
	return testCase{input: sb.String()}
}

func runCase(bin, ref string, tc testCase) error {
	expected, err := runBinary(ref, tc.input)
	if err != nil {
		return fmt.Errorf("reference failed: %v", err)
	}
	got, err := runBinary(bin, tc.input)
	if err != nil {
		return err
	}
	if strings.TrimSpace(got) != strings.TrimSpace(expected) {
		return fmt.Errorf("expected %s got %s", expected, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierF.go /path/to/binary")
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

	cases := make([]testCase, 100)
	for i := 0; i < 100; i++ {
		cases[i] = genCase(rng)
	}

	for i, tc := range cases {
		if err := runCase(bin, ref, tc); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, tc.input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
