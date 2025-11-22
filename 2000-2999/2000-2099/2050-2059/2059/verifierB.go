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
	n int
	k int
	a []int
}

func buildRef() (string, error) {
	_, self, _, _ := runtime.Caller(0)
	dir := filepath.Dir(self)
	ref := filepath.Join(dir, "refB.bin")
	cmd := exec.Command("go", "build", "-o", ref, filepath.Join(dir, "2059B.go"))
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
		fmt.Fprintf(&sb, "%d %d\n", tc.n, tc.k)
		for i, v := range tc.a {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.Itoa(v))
		}
		sb.WriteByte('\n')
	}
	return sb.String()
}

func parseOutputs(out string, t int) ([]int64, error) {
	fields := strings.Fields(out)
	if len(fields) != t {
		return nil, fmt.Errorf("expected %d answers, got %d", t, len(fields))
	}
	ans := make([]int64, t)
	for i, f := range fields {
		v, err := strconv.ParseInt(f, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("cannot parse answer %q: %v", f, err)
		}
		ans[i] = v
	}
	return ans, nil
}

func randomCase(r *rand.Rand, nLimit int) testCase {
	n := r.Intn(nLimit-1) + 2
	k := r.Intn(n/2)*2 + 2 // even between 2 and n
	if k > n {
		k = n
		if k%2 == 1 {
			k--
		}
		if k < 2 {
			k = 2
		}
	}
	a := make([]int, n)
	for i := range a {
		// keep values reasonable but within 1..1e9
		a[i] = r.Intn(1_000_000_000) + 1
	}
	return testCase{n: n, k: k, a: a}
}

func genCases() []testCase {
	cases := []testCase{
		{n: 3, k: 2, a: []int{1, 1, 1}},
		{n: 8, k: 8, a: []int{1, 1, 2, 2, 3, 3, 4, 4}},
		{n: 5, k: 4, a: []int{1, 1, 1_000_000_000, 2, 2}},
	}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	sumN := 0
	for _, tc := range cases {
		sumN += tc.n
	}
	for sumN < 190000 && len(cases) < 60 {
		tc := randomCase(r, 2000)
		if sumN+tc.n > 200000 {
			break
		}
		sumN += tc.n
		cases = append(cases, tc)
	}
	return cases
}

func check(bin, ref string, cases []testCase) error {
	input := buildInput(cases)
	refOut, err := runBinary(ref, input)
	if err != nil {
		return fmt.Errorf("reference failed: %v", err)
	}
	expect, err := parseOutputs(refOut, len(cases))
	if err != nil {
		return fmt.Errorf("reference output invalid: %v", err)
	}

	out, err := runBinary(bin, input)
	if err != nil {
		return err
	}
	got, err := parseOutputs(out, len(cases))
	if err != nil {
		return err
	}
	for i := range expect {
		if expect[i] != got[i] {
			return fmt.Errorf("answer mismatch on test %d: expected %d got %d", i+1, expect[i], got[i])
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
		fmt.Println("usage: go run verifierB.go /path/to/binary")
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
