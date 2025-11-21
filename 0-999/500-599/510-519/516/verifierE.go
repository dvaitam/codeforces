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
	id    string
	input string
}

func main() {
	if len(os.Args) < 2 || len(os.Args) > 3 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/candidate")
		os.Exit(1)
	}
	target := os.Args[len(os.Args)-1]
	if target == "--" {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/candidate")
		os.Exit(1)
	}

	base := currentDir()
	refBin, err := buildReference(base)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to build reference: %v\n", err)
		os.Exit(1)
	}
	defer os.Remove(refBin)

	tests := generateTests()
	for i, tc := range tests {
		exp, err := runProgram(refBin, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference error on %s: %v\ninput:\n%s", tc.id, err, tc.input)
			os.Exit(1)
		}
		got, err := runProgram(target, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate error on %s: %v\ninput:\n%s", tc.id, err, tc.input)
			os.Exit(1)
		}
		if strings.TrimSpace(exp) == "" {
			fmt.Fprintf(os.Stderr, "reference produced empty output on %s\n", tc.id)
			os.Exit(1)
		}
		expVal, err := parseInt(exp)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to parse reference output on %s: %v\noutput:\n%s", tc.id, err, exp)
			os.Exit(1)
		}
		gotVal, err := parseInt(got)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to parse candidate output on %s: %v\noutput:\n%s", tc.id, err, got)
			os.Exit(1)
		}
		if expVal != gotVal {
			fmt.Fprintf(os.Stderr, "wrong answer on %s\nInput:\n%sExpected: %d\nGot: %d\n", tc.id, tc.input, expVal, gotVal)
			os.Exit(1)
		}
		if (i+1)%10 == 0 {
			fmt.Fprintf(os.Stderr, "validated %d/%d tests...\n", i+1, len(tests))
		}
	}
	fmt.Printf("All %d tests passed.\n", len(tests))
}

func currentDir() string {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		panic("cannot determine current file")
	}
	return filepath.Dir(file)
}

func buildReference(dir string) (string, error) {
	out := filepath.Join(dir, "ref516E.bin")
	src := filepath.Join(dir, "516E.go")
	code, err := os.ReadFile(src)
	if err != nil {
		return "", fmt.Errorf("read reference: %v", err)
	}
	patch := bytes.Replace(code, []byte("const INF64 = 1<<63 - 1"), []byte("const INF64 int64 = 1<<63 - 1"), 1)
	tmp, err := os.CreateTemp(dir, "516E_ref_*.go")
	if err != nil {
		return "", fmt.Errorf("create temp file: %v", err)
	}
	tmpPath := tmp.Name()
	if _, err := tmp.Write(patch); err != nil {
		tmp.Close()
		os.Remove(tmpPath)
		return "", fmt.Errorf("write temp file: %v", err)
	}
	tmp.Close()
	cmd := exec.Command("go", "build", "-o", out, tmpPath)
	cmd.Dir = dir
	if data, err := cmd.CombinedOutput(); err != nil {
		os.Remove(tmpPath)
		return "", fmt.Errorf("go build failed: %v\n%s", err, data)
	}
	os.Remove(tmpPath)
	return out, nil
}

func runProgram(target, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(target, ".go") {
		cmd = exec.Command("go", "run", target)
	} else {
		cmd = exec.Command(target)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return out.String(), fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func parseInt(out string) (int64, error) {
	reader := strings.NewReader(out)
	var v int64
	if _, err := fmt.Fscan(reader, &v); err != nil {
		return 0, err
	}
	return v, nil
}

func generateTests() []testCase {
	var tests []testCase
	// deterministic edge cases
	tests = append(tests, formatCase("simple-1",
		4, 6,
		[]int64{0, 2},
		[]int64{1, 3}))
	tests = append(tests, formatCase("missing-residue",
		4, 6,
		[]int64{0},
		[]int64{2}))
	tests = append(tests, formatCase("larger-gcd",
		6, 9,
		[]int64{0, 3},
		[]int64{0, 6}))
	tests = append(tests, formatCase("tricky",
		18, 24,
		[]int64{0, 5, 10, 15},
		[]int64{3, 9, 12}))

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 80; i++ {
		tests = append(tests, randomCase(rng, fmt.Sprintf("rand-%02d", i+1)))
	}
	return tests
}

func formatCase(id string, n, m int64, boys, girls []int64) testCase {
	input := formatInput(n, m, boys, girls)
	return testCase{id: id, input: input}
}

func randomCase(rng *rand.Rand, id string) testCase {
	n := int64(rng.Intn(1000) + 1)
	m := int64(rng.Intn(1000) + 1)
	bCnt := rng.Intn(int(minInt64(n, 60))) + 1
	gCnt := rng.Intn(int(minInt64(m, 60))) + 1
	boys := pickDistinct(rng, n, bCnt)
	girls := pickDistinct(rng, m, gCnt)

	// ensure not all friends already happy
	if int64(len(boys)) == n && int64(len(girls)) == m {
		girls = girls[:len(girls)-1]
	}

	// occasionally force failure by removing coverage for a residue
	if rng.Float64() < 0.3 {
		d := gcdInt64(n, m)
		if d > 1 {
			miss := int64(rng.Intn(int(d)))
			boys = filterResidue(boys, d, miss)
			girls = filterResidue(girls, d, miss)
			// ensure we still have at least one of each gender to avoid reference bug
			if len(boys) == 0 {
				boys = append(boys, int64((miss+1)%d))
			}
			if len(girls) == 0 {
				girls = append(girls, int64((miss+2)%d))
			}
		}
	}

	// ensure both sides have at least one element (reference requirement)
	if len(boys) == 0 {
		boys = []int64{0}
	}
	if len(girls) == 0 {
		girls = []int64{0}
	}

	return formatCase(id, n, m, boys, girls)
}

func formatInput(n, m int64, boys, girls []int64) string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d\n", n, m)
	fmt.Fprintf(&sb, "%d", len(boys))
	for _, b := range boys {
		fmt.Fprintf(&sb, " %d", b%n)
	}
	sb.WriteByte('\n')
	fmt.Fprintf(&sb, "%d", len(girls))
	for _, g := range girls {
		fmt.Fprintf(&sb, " %d", g%m)
	}
	sb.WriteByte('\n')
	return sb.String()
}

func pickDistinct(rng *rand.Rand, mod int64, count int) []int64 {
	if count <= 0 {
		return nil
	}
	set := make(map[int64]struct{})
	for len(set) < count {
		val := int64(rng.Intn(int(mod)))
		set[val] = struct{}{}
	}
	res := make([]int64, 0, len(set))
	for v := range set {
		res = append(res, v)
	}
	return res
}

func filterResidue(arr []int64, d, miss int64) []int64 {
	res := arr[:0]
	for _, v := range arr {
		if v%d != miss {
			res = append(res, v)
		}
	}
	return append([]int64(nil), res...)
}

func gcdInt64(a, b int64) int64 {
	if a < 0 {
		a = -a
	}
	if b < 0 {
		b = -b
	}
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

func minInt64(a, b int64) int64 {
	if a < b {
		return a
	}
	return b
}
