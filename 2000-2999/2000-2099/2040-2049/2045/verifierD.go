package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"
)

type testCase struct {
	N      int
	D      int
	Ts     int
	Tf     int
	Tw     int
	powers []int
}

func buildRef() (string, error) {
	_, self, _, _ := runtime.Caller(0)
	dir := filepath.Dir(self)
	ref := filepath.Join(dir, "refD.bin")
	cmd := exec.Command("go", "build", "-o", ref, filepath.Join(dir, "2045D.go"))
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
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func buildInput(tc testCase) string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d %d %d %d\n", tc.N, tc.D, tc.Ts, tc.Tf, tc.Tw)
	for i, v := range tc.powers {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(v))
	}
	sb.WriteByte('\n')
	return sb.String()
}

func parseOutput(out string) (int64, error) {
	fields := strings.Fields(out)
	if len(fields) == 0 {
		return 0, fmt.Errorf("empty output")
	}
	val, err := strconv.ParseInt(fields[0], 10, 64)
	if err != nil {
		return 0, fmt.Errorf("cannot parse answer %q: %v", fields[0], err)
	}
	return val, nil
}

func randomCase(r *rand.Rand, nLimit int) testCase {
	N := r.Intn(nLimit-1) + 2
	D := r.Intn(200000) + 1
	Ts := r.Intn(200000) + 1
	Tf := r.Intn(200000) + 1
	Tw := r.Intn(200000) + 1
	powers := make([]int, N)
	for i := range powers {
		powers[i] = r.Intn(200000) + 1
	}
	return testCase{N: N, D: D, Ts: Ts, Tf: Tf, Tw: Tw, powers: powers}
}

func genCases() []testCase {
	// Deterministic small cases (including samples).
	cases := []testCase{
		{N: 5, D: 4, Ts: 2, Tf: 9, Tw: 1, powers: []int{1, 2, 4, 2, 1}},
		{N: 5, D: 4, Ts: 2, Tf: 1, Tw: 1, powers: []int{1, 2, 4, 2, 1}},
		{N: 3, D: 4, Ts: 2, Tf: 10, Tw: 1, powers: []int{3, 1, 2}},
		{N: 2, D: 1, Ts: 1, Tf: 5, Tw: 1, powers: []int{1, 1}},
	}

	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	// Moderate random cases.
	for i := 0; i < 12; i++ {
		cases = append(cases, randomCase(r, 400))
	}

	// Larger stress but still comfortable for verifier.
	cases = append(cases, randomCase(r, 5000))
	cases = append(cases, randomCase(r, 200000))

	return cases
}

func checkCase(bin, ref string, tc testCase, idx int) error {
	input := buildInput(tc)
	refOut, err := runBinary(ref, input)
	if err != nil {
		return fmt.Errorf("reference failed: %v", err)
	}
	expect, err := parseOutput(refOut)
	if err != nil {
		return fmt.Errorf("reference output invalid: %v", err)
	}

	out, err := runBinary(bin, input)
	if err != nil {
		return fmt.Errorf("runtime error: %v", err)
	}
	got, err := parseOutput(out)
	if err != nil {
		return err
	}
	if expect != got {
		return fmt.Errorf("case %d mismatch: expected %d got %d", idx, expect, got)
	}
	return nil
}

func main() {
	exitCode := 0
	cleanup := func() {}
	defer func() {
		cleanup()
		if exitCode != 0 {
			os.Exit(exitCode)
		}
	}()

	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		exitCode = 1
		return
	}
	bin := os.Args[1]
	ref, err := buildRef()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		exitCode = 1
		return
	}
	cleanup = func() { _ = os.Remove(ref) }

	cases := genCases()
	for i, tc := range cases {
		if err := checkCase(bin, ref, tc, i+1); err != nil {
			fmt.Fprintf(os.Stderr, "%v\ninput:\n%s", err, buildInput(tc))
			exitCode = 1
			return
		}
	}
	fmt.Printf("All %d test cases passed\n", len(cases))
}
