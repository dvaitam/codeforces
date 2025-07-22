package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

type testCase struct {
	pl, pr, vl, vr, k int
}

func compileRef() (string, error) {
	_, file, _, _ := runtime.Caller(0)
	dir := filepath.Dir(file)
	exe := filepath.Join(os.TempDir(), "ref109B")
	cmd := exec.Command("go", "build", "-o", exe, filepath.Join(dir, "109B.go"))
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return "", err
	}
	return exe, nil
}

func generateCase(rng *rand.Rand) testCase {
	pl := rng.Intn(1000000000) + 1
	pr := pl + rng.Intn(1000)
	vl := rng.Intn(1000000000) + 1
	vr := vl + rng.Intn(1000)
	k := rng.Intn(10) + 1
	return testCase{pl, pr, vl, vr, k}
}

func runCase(bin, ref string, tc testCase) error {
	input := fmt.Sprintf("%d %d %d %d %d\n", tc.pl, tc.pr, tc.vl, tc.vr, tc.k)
	run := func(path string) (string, error) {
		cmd := exec.Command(path)
		cmd.Stdin = strings.NewReader(input)
		var out bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = &out
		err := cmd.Run()
		return strings.TrimSpace(out.String()), err
	}
	expected, err := run(ref)
	if err != nil {
		return fmt.Errorf("reference error: %v", err)
	}
	got, err := run(bin)
	if err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, got)
	}
	if expected == "" {
		expected = "0"
	}
	if got == "" {
		got = "0"
	}
	var expVal, gotVal float64
	if _, err := fmt.Sscan(expected, &expVal); err != nil {
		return fmt.Errorf("bad ref output: %v", err)
	}
	if _, err := fmt.Sscan(got, &gotVal); err != nil {
		return fmt.Errorf("bad output: %v", err)
	}
	diff := gotVal - expVal
	if diff < -1e-6 || diff > 1e-6 {
		return fmt.Errorf("expected %.6f got %.6f", expVal, gotVal)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	ref, err := compileRef()
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to build reference solution:", err)
		os.Exit(1)
	}
	defer os.Remove(ref)
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	cases := []testCase{
		{1, 10, 1, 10, 2},
		{1, 100, 1, 100, 3},
	}
	for i := 0; i < 100; i++ {
		cases = append(cases, generateCase(rng))
	}
	for i, tc := range cases {
		if err := runCase(bin, ref, tc); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput: %v\n", i+1, err, tc)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
