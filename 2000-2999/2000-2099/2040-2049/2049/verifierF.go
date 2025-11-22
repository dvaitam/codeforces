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

type update struct {
	i int
	x int
}

type testCase struct {
	n   int
	q   int
	a   []int
	ops []update
}

func buildRef() (string, error) {
	_, self, _, _ := runtime.Caller(0)
	dir := filepath.Dir(self)
	ref := filepath.Join(dir, "refF.bin")
	cmd := exec.Command("go", "build", "-o", ref, filepath.Join(dir, "2049F.go"))
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
		fmt.Fprintf(&sb, "%d %d\n", tc.n, tc.q)
		for i, v := range tc.a {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.Itoa(v))
		}
		sb.WriteByte('\n')
		for _, op := range tc.ops {
			fmt.Fprintf(&sb, "%d %d\n", op.i, op.x)
		}
	}
	return sb.String()
}

func parseOutputs(out string, total int) ([]int, error) {
	fields := strings.Fields(out)
	if len(fields) != total {
		return nil, fmt.Errorf("expected %d answers, got %d", total, len(fields))
	}
	ans := make([]int, total)
	for i, f := range fields {
		v, err := strconv.Atoi(f)
		if err != nil {
			return nil, fmt.Errorf("failed to parse answer %q: %v", f, err)
		}
		ans[i] = v
	}
	return ans, nil
}

func randomCase(r *rand.Rand, nLim, qLim int) testCase {
	n := r.Intn(nLim) + 1
	q := r.Intn(qLim) + 1
	a := make([]int, n)
	for i := range a {
		a[i] = r.Intn(n + 1)
	}
	ops := make([]update, q)
	for i := 0; i < q; i++ {
		pos := r.Intn(n) + 1 // 1-based for input
		x := r.Intn(n) + 1
		ops[i] = update{i: pos, x: x}
	}
	return testCase{n: n, q: q, a: a, ops: ops}
}

func genCases() []testCase {
	cases := []testCase{
		// sample 1 case with 3 updates
		{
			n: 6, q: 3,
			a:   []int{0, 0, 1, 0, 1, 0},
			ops: []update{{i: 6, x: 1}, {i: 3, x: 2}, {i: 6, x: 3}},
		},
	}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	sumN, sumQ := 6, 3
	for len(cases) < 25 && sumN < 90000 && sumQ < 90000 {
		tc := randomCase(r, 400, 400)
		if sumN+tc.n > 100000 || sumQ+tc.q > 100000 {
			break
		}
		sumN += tc.n
		sumQ += tc.q
		cases = append(cases, tc)
	}

	// One larger case while staying inside global limits to stress performance.
	remN := 100000 - sumN
	remQ := 100000 - sumQ
	if remN > 0 && remQ > 0 {
		nBig := remN
		if nBig > 20000 {
			nBig = 20000
		}
		qBig := remQ
		if qBig > 20000 {
			qBig = 20000
		}
		a := make([]int, nBig)
		for i := range a {
			a[i] = r.Intn(nBig + 1)
		}
		ops := make([]update, qBig)
		for i := 0; i < qBig; i++ {
			ops[i] = update{i: r.Intn(nBig) + 1, x: r.Intn(nBig) + 1}
		}
		cases = append(cases, testCase{n: nBig, q: qBig, a: a, ops: ops})
	}
	return cases
}

func check(bin, ref string, cases []testCase) error {
	input := buildInput(cases)
	totalAns := 0
	for _, tc := range cases {
		totalAns += tc.q
	}

	refOut, err := runBinary(ref, input)
	if err != nil {
		return fmt.Errorf("reference failed: %v", err)
	}
	expect, err := parseOutputs(refOut, totalAns)
	if err != nil {
		return fmt.Errorf("reference output invalid: %v", err)
	}

	out, err := runBinary(bin, input)
	if err != nil {
		return err
	}
	got, err := parseOutputs(out, totalAns)
	if err != nil {
		return err
	}
	for i := range expect {
		if expect[i] != got[i] {
			return fmt.Errorf("answer mismatch at position %d: expected %d got %d", i+1, expect[i], got[i])
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
		fmt.Println("usage: go run verifierF.go /path/to/binary")
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
	if err := check(bin, ref, cases); err != nil {
		fmt.Fprintf(os.Stderr, "%v\ninput:\n%s", err, buildInput(cases))
		exitCode = 1
		return
	}
	fmt.Printf("All %d test cases passed\n", len(cases))
}
