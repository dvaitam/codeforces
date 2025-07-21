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

type testCase struct{ input string }

func compileRef() (string, error) {
	exe := filepath.Join(os.TempDir(), "ref44G")
	cmd := exec.Command("go", "build", "-o", exe, "44G.go")
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return "", err
	}
	return exe, nil
}

func generateRandomCase(rng *rand.Rand) testCase {
	n := rng.Intn(3) + 1
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for i := 0; i < n; i++ {
		xl := rng.Intn(10)
		xr := xl + rng.Intn(3) + 1
		yl := rng.Intn(10)
		yr := yl + rng.Intn(3) + 1
		z := rng.Intn(10)
		sb.WriteString(fmt.Sprintf("%d %d %d %d %d\n", xl, xr, yl, yr, z))
	}
	m := rng.Intn(3) + 1
	sb.WriteString(fmt.Sprintf("%d\n", m))
	for i := 0; i < m; i++ {
		x := rng.Intn(12)
		y := rng.Intn(12)
		sb.WriteString(fmt.Sprintf("%d %d\n", x, y))
	}
	return testCase{input: sb.String()}
}

func runCase(bin, ref string, tc testCase) error {
	run := func(path string) (string, error) {
		cmd := exec.Command(path)
		cmd.Stdin = strings.NewReader(tc.input)
		var out bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = &out
		err := cmd.Run()
		return strings.TrimSpace(out.String()), err
	}
	exp, err := run(ref)
	if err != nil {
		return fmt.Errorf("ref error: %v", err)
	}
	got, err := run(bin)
	if err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, got)
	}
	if exp != got {
		return fmt.Errorf("expected %q got %q", exp, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierG.go /path/to/binary")
		os.Exit(1)
	}
	ref, err := compileRef()
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to build reference:", err)
		os.Exit(1)
	}
	defer os.Remove(ref)
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	bin := os.Args[1]
	cases := []testCase{{input: "1\n0 1 0 1 0\n1\n0 0\n"}}
	for i := 0; i < 100; i++ {
		cases = append(cases, generateRandomCase(rng))
	}
	for i, tc := range cases {
		if err := runCase(bin, ref, tc); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, tc.input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
