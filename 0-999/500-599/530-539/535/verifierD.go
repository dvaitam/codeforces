package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
)

const (
	mod                  = 1000000007
	referenceSolutionRel = "0-999/500-599/530-539/535/535D.go"
)

var referenceSolutionPath string

func init() {
	referenceSolutionPath = referenceSolutionRel
	if _, file, _, ok := runtime.Caller(0); ok {
		dir := filepath.Dir(file)
		candidate := filepath.Join(dir, "535D.go")
		if _, err := os.Stat(candidate); err == nil {
			referenceSolutionPath = candidate
			return
		}
	}
	if abs, err := filepath.Abs(referenceSolutionRel); err == nil {
		if _, err := os.Stat(abs); err == nil {
			referenceSolutionPath = abs
		}
	}
}

type testCase struct {
	n  int
	p  string
	ys []int
}

func makeZ(s string) []int {
	n := len(s)
	if n == 0 {
		return nil
	}
	z := make([]int, n)
	z[0] = n
	l, r := 0, 0
	for i := 1; i < n; i++ {
		if i <= r {
			k := i - l
			if z[k] < r-i+1 {
				z[i] = z[k]
			} else {
				j := r + 1
				for j < n && s[j] == s[j-i] {
					j++
				}
				z[i] = j - i
				l, r = i, j-1
			}
		} else {
			j := 0
			for i+j < n && s[j] == s[i+j] {
				j++
			}
			z[i] = j
			if j > 0 {
				l, r = i, i+j-1
			}
		}
	}
	return z
}

func modPow(a, e int64) int64 {
	res := int64(1)
	a %= mod
	for e > 0 {
		if e&1 == 1 {
			res = res * a % mod
		}
		a = a * a % mod
		e >>= 1
	}
	return res
}

func calcAnswer(tc testCase) (int64, bool) {
	n := tc.n
	p := tc.p
	lenP := len(p)
	ys := tc.ys
	if lenP > n {
		if len(ys) > 0 {
			return 0, false
		}
		return modPow(26, int64(n)), true
	}
	maxStart := n - lenP + 1
	for i := 0; i < len(ys); i++ {
		if ys[i] < 1 || ys[i] > maxStart {
			return 0, false
		}
		if i > 0 && ys[i] <= ys[i-1] {
			return 0, false
		}
	}
	z := makeZ(p)
	lastEnd := 0
	var forced int64
	for i, y := range ys {
		start := y
		end := y + lenP - 1
		if i > 0 {
			d := y - ys[i-1]
			if d < lenP {
				if z[d] < lenP-d {
					return 0, true
				}
			}
		}
		if start > lastEnd {
			forced += int64(lenP)
		} else {
			overlap := lastEnd - start + 1
			if overlap < lenP {
				forced += int64(lenP - overlap)
			}
		}
		if end > lastEnd {
			lastEnd = end
		}
	}
	free := int64(n) - forced
	if free < 0 {
		free = 0
	}
	return modPow(26, free), true
}

func inputString(tc testCase) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", tc.n, len(tc.ys)))
	sb.WriteString(tc.p)
	sb.WriteByte('\n')
	for i, v := range tc.ys {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", v))
	}
	sb.WriteByte('\n')
	return sb.String()
}

func runProgram(path, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(path, ".go") {
		cmd = exec.Command("go", "run", path)
	} else {
		cmd = exec.Command(path)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

func parseAnswer(out string) (int64, error) {
	reader := strings.NewReader(out)
	var val int64
	if _, err := fmt.Fscan(reader, &val); err != nil {
		return 0, fmt.Errorf("failed to parse integer from output: %v\nfull output:\n%s", err, out)
	}
	val %= mod
	if val < 0 {
		val += mod
	}
	return val, nil
}

func buildReferenceBinary() (string, func(), error) {
	tmpDir, err := os.MkdirTemp("", "535D-ref")
	if err != nil {
		return "", nil, err
	}
	binPath := filepath.Join(tmpDir, "ref_535D")
	cmd := exec.Command("go", "build", "-o", binPath, referenceSolutionPath)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		os.RemoveAll(tmpDir)
		return "", nil, fmt.Errorf("failed to build reference: %v\n%s", err, out.String())
	}
	cleanup := func() {
		os.RemoveAll(tmpDir)
	}
	return binPath, cleanup, nil
}

func randomLowerString(rng *rand.Rand, length int) string {
	if length <= 0 {
		return ""
	}
	b := make([]byte, length)
	for i := range b {
		b[i] = byte('a' + rng.Intn(26))
	}
	return string(b)
}

func randomValidCase(rng *rand.Rand) testCase {
	n := rng.Intn(50) + 1
	lenP := rng.Intn(n) + 1
	base := randomLowerString(rng, n)
	start := rng.Intn(n - lenP + 1)
	p := base[start : start+lenP]
	var occ []int
	for i := 0; i+lenP <= n; i++ {
		if base[i:i+lenP] == p {
			occ = append(occ, i+1)
		}
	}
	var ys []int
	for _, pos := range occ {
		if rng.Intn(2) == 0 {
			ys = append(ys, pos)
		}
	}
	if len(ys) == 0 && len(occ) > 0 && rng.Intn(3) == 0 {
		ys = append(ys, occ[rng.Intn(len(occ))])
		sort.Ints(ys)
	}
	return testCase{n: n, p: p, ys: ys}
}

func randomGeneralCase(rng *rand.Rand) testCase {
	n := rng.Intn(60) + 1
	lenP := rng.Intn(n) + 1
	p := randomLowerString(rng, lenP)
	maxStart := n - lenP + 1
	count := 0
	if maxStart > 0 {
		count = rng.Intn(maxStart + 1)
	}
	var ys []int
	if count > 0 {
		perm := rng.Perm(maxStart)
		for i := 0; i < count; i++ {
			ys = append(ys, perm[i]+1)
		}
		sort.Ints(ys)
	}
	return testCase{n: n, p: p, ys: ys}
}

func genTests() []testCase {
	rng := rand.New(rand.NewSource(20250212))
	tests := []testCase{
		{n: 1, p: "a"},
		{n: 5, p: "ab", ys: []int{1}},
		{n: 5, p: "aa", ys: []int{1, 2, 3}},
		{n: 6, p: "abc", ys: []int{1, 2}},
		{n: 8, p: "aba", ys: []int{1, 3}},
	}
	for i := 0; i < 150; i++ {
		tests = append(tests, randomValidCase(rng))
	}
	for i := 0; i < 150; i++ {
		tests = append(tests, randomGeneralCase(rng))
	}
	return tests
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[len(os.Args)-1]
	if bin == "--" {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	refBin, cleanup, err := buildReferenceBinary()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer cleanup()
	tests := genTests()
	for i, tc := range tests {
		in := inputString(tc)
		expected, ok := calcAnswer(tc)
		if !ok {
			expected = 0
		}
		refOut, err := runProgram(refBin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference runtime error on test %d: %v\ninput:\n%soutput:\n%s\n", i+1, err, in, refOut)
			os.Exit(1)
		}
		refAns, err := parseAnswer(refOut)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference output parse error on test %d: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
		if refAns != expected {
			fmt.Fprintf(os.Stderr, "reference mismatch on test %d: expected %d got %d\ninput:\n%soutput:\n%s\n", i+1, expected, refAns, in, refOut)
			os.Exit(1)
		}
		out, runErr := runProgram(bin, in)
		if runErr != nil {
			fmt.Fprintf(os.Stderr, "test %d runtime error: %v\ninput:\n%soutput:\n%s\n", i+1, runErr, in, out)
			os.Exit(1)
		}
		ans, err := parseAnswer(out)
		if err != nil {
			fmt.Fprintf(os.Stderr, "test %d parse error: %v\ninput:\n%soutput:\n%s\n", i+1, err, in, out)
			os.Exit(1)
		}
		if ans != expected {
			fmt.Fprintf(os.Stderr, "test %d failed: expected %d got %d\ninput:\n%soutput:\n%s\n", i+1, expected, ans, in, out)
			os.Exit(1)
		}
	}
	fmt.Println("all tests passed")
}
