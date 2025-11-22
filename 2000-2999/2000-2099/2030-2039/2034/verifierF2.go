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
	n  int
	m  int
	k  int
	rm [][2]int
}

func buildRef() (string, error) {
	_, self, _, _ := runtime.Caller(0)
	dir := filepath.Dir(self)
	ref := filepath.Join(dir, "refF2.bin")
	cmd := exec.Command("go", "build", "-o", ref, filepath.Join(dir, "2034F2.go"))
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

func buildInput(cases []testCase) string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", len(cases))
	for _, tc := range cases {
		fmt.Fprintf(&sb, "%d %d %d\n", tc.n, tc.m, tc.k)
		for _, p := range tc.rm {
			fmt.Fprintf(&sb, "%d %d\n", p[0], p[1])
		}
	}
	return sb.String()
}

func parseOutputs(out string, t int) ([]int, error) {
	fields := strings.Fields(out)
	if len(fields) != t {
		return nil, fmt.Errorf("expected %d answers, got %d", t, len(fields))
	}
	ans := make([]int, t)
	for i, f := range fields {
		v, err := strconv.Atoi(f)
		if err != nil {
			return nil, fmt.Errorf("failed to parse answer %q: %v", f, err)
		}
		if v < 0 || v >= 998244353 {
			return nil, fmt.Errorf("answer out of range: %d", v)
		}
		ans[i] = v
	}
	return ans, nil
}

func randomCase(r *rand.Rand, limit int) testCase {
	n := r.Intn(limit-1) + 1
	m := r.Intn(limit-1) + 1
	if n < 1 {
		n = 1
	}
	if m < 1 {
		m = 1
	}
	maxK := r.Intn(10) + 1
	k := r.Intn(maxK + 1)
	if k > 5000 {
		k = 5000
	}
	rm := make([][2]int, 0, k)
	seen := make(map[[2]int]struct{})
	for len(rm) < k {
		rv := r.Intn(n + 1)
		bv := r.Intn(m + 1)
		if rv+bv == 0 || rv+bv >= n+m {
			continue
		}
		key := [2]int{rv, bv}
		if _, ok := seen[key]; ok {
			continue
		}
		seen[key] = struct{}{}
		rm = append(rm, key)
	}
	return testCase{n: n, m: m, k: len(rm), rm: rm}
}

func genCases() []testCase {
	cases := []testCase{
		{n: 1, m: 1, k: 0},
		{n: 3, m: 4, k: 1, rm: [][2]int{{1, 0}}},
		{n: 5, m: 5, k: 3, rm: [][2]int{{2, 2}, {1, 4}, {4, 1}}},
	}
	sumBudget := 180000
	used := 0
	for _, tc := range cases {
		used += tc.n + tc.m
	}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for len(cases) < 40 && used < sumBudget {
		tc := randomCase(r, 400)
		if used+tc.n+tc.m > sumBudget {
			break
		}
		used += tc.n + tc.m
		cases = append(cases, tc)
	}
	// Add one larger stress within budget.
	if used+150000 <= sumBudget {
		tc := testCase{n: 100000, m: 50000, k: 2, rm: [][2]int{{50000, 20000}, {40000, 10000}}}
		cases = append(cases, tc)
	}
	return cases
}

func checkCase(bin, ref string, tc []testCase) error {
	input := buildInput(tc)
	refOut, err := runBinary(ref, input)
	if err != nil {
		return fmt.Errorf("reference failed: %v", err)
	}
	expect, err := parseOutputs(refOut, len(tc))
	if err != nil {
		return fmt.Errorf("reference output invalid: %v", err)
	}

	out, err := runBinary(bin, input)
	if err != nil {
		return err
	}
	got, err := parseOutputs(out, len(tc))
	if err != nil {
		return err
	}
	for i := range expect {
		if expect[i] != got[i] {
			return fmt.Errorf("test %d mismatch: expected %d got %d", i+1, expect[i], got[i])
		}
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
		fmt.Println("usage: go run verifierF2.go /path/to/binary")
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
	if err := checkCase(bin, ref, cases); err != nil {
		fmt.Fprintf(os.Stderr, "%v\ninput:\n%s", err, buildInput(cases))
		exitCode = 1
		return
	}
	fmt.Printf("All %d test cases passed\n", len(cases))
}
