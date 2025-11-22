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

const refSource = "2123G.go"

type query struct {
	t int
	i int
	x int
	k int
}

type testCase struct {
	name string
	n    int
	m    int
	a    []int
	qs   []query
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierG.go /path/to/candidate_binary")
		os.Exit(1)
	}
	candidate := os.Args[1]

	tests := buildTests()
	input := buildInput(tests)
	expected := evaluateAll(tests)

	refBin, cleanup, err := buildReference()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to build reference: %v\n", err)
		os.Exit(1)
	}
	defer cleanup()

	refOut, err := runProgram(refBin, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "reference runtime error: %v\ninput:\n%soutput:\n%s", err, input, refOut)
		os.Exit(1)
	}
	refAns, err := parseAnswers(refOut, len(expected))
	if err != nil {
		fmt.Fprintf(os.Stderr, "reference output invalid: %v\ninput:\n%soutput:\n%s", err, input, refOut)
		os.Exit(1)
	}
	for idx, exp := range expected {
		if refAns[idx] != exp {
			fmt.Fprintf(os.Stderr, "reference answer mismatch at #%d: expected %s got %s\ninput:\n%soutput:\n%s",
				idx+1, exp, refAns[idx], input, refOut)
			os.Exit(1)
		}
	}

	candOut, err := runProgram(candidate, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "candidate runtime error: %v\ninput:\n%soutput:\n%s", err, input, candOut)
		os.Exit(1)
	}
	candAns, err := parseAnswers(candOut, len(expected))
	if err != nil {
		fmt.Fprintf(os.Stderr, "candidate output invalid: %v\ninput:\n%soutput:\n%s", err, input, candOut)
		os.Exit(1)
	}
	for idx, exp := range expected {
		if candAns[idx] != exp {
			fmt.Fprintf(os.Stderr, "candidate failed at answer #%d: expected %s got %s\ninput:\n%soutput:\n%s",
				idx+1, exp, candAns[idx], input, candOut)
			os.Exit(1)
		}
	}

	fmt.Printf("All %d tests passed\n", len(expected))
}

func buildReference() (string, func(), error) {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		return "", nil, fmt.Errorf("cannot locate verifier path")
	}
	dir := filepath.Dir(file)
	tmpDir, err := os.MkdirTemp("", "oracle-2123G-")
	if err != nil {
		return "", nil, err
	}
	binPath := filepath.Join(tmpDir, "oracleG")
	cmd := exec.Command("go", "build", "-o", binPath, refSource)
	cmd.Dir = dir
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		os.RemoveAll(tmpDir)
		return "", nil, fmt.Errorf("build failed: %v\n%s", err, stderr.String())
	}
	cleanup := func() {
		_ = os.RemoveAll(tmpDir)
	}
	return binPath, cleanup, nil
}

func runProgram(bin, input string) (string, error) {
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
	err := cmd.Run()
	return out.String(), err
}

func buildInput(tests []testCase) string {
	var sb strings.Builder
	sb.Grow(256)
	sb.WriteString(strconv.Itoa(len(tests)))
	sb.WriteByte('\n')
	for _, tc := range tests {
		sb.WriteString(fmt.Sprintf("%d %d %d\n", tc.n, tc.m, len(tc.qs)))
		for i, v := range tc.a {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.Itoa(v))
		}
		sb.WriteByte('\n')
		for _, qu := range tc.qs {
			if qu.t == 1 {
				sb.WriteString(fmt.Sprintf("1 %d %d\n", qu.i+1, qu.x))
			} else {
				sb.WriteString(fmt.Sprintf("2 %d\n", qu.k))
			}
		}
	}
	return sb.String()
}

func parseAnswers(out string, need int) ([]string, error) {
	fields := strings.Fields(out)
	if len(fields) != need {
		return nil, fmt.Errorf("expected %d answers, got %d", need, len(fields))
	}
	res := make([]string, need)
	for i, f := range fields {
		up := strings.ToUpper(f)
		if strings.HasPrefix(up, "Y") {
			res[i] = "YES"
		} else if strings.HasPrefix(up, "N") {
			res[i] = "NO"
		} else {
			return nil, fmt.Errorf("answer %d: invalid token %q", i+1, f)
		}
	}
	return res, nil
}

func evaluateAll(tests []testCase) []string {
	ans := make([]string, 0, 1024)
	for _, tc := range tests {
		ans = append(ans, solveCase(tc)...)
	}
	return ans
}

func solveCase(tc testCase) []string {
	gset := make(map[int]struct{})
	for _, q := range tc.qs {
		if q.t == 2 {
			gset[gcd(q.k, tc.m)] = struct{}{}
		}
	}
	glist := make([]int, 0, len(gset))
	for g := range gset {
		glist = append(glist, g)
	}
	desc := make(map[int]int)
	for _, g := range glist {
		cnt := 0
		for i := 0; i+1 < tc.n; i++ {
			if tc.a[i]%g > tc.a[i+1]%g {
				cnt++
			}
		}
		desc[g] = cnt
	}

	ans := make([]string, 0, len(tc.qs))
	a := make([]int, len(tc.a))
	copy(a, tc.a)
	for _, qu := range tc.qs {
		if qu.t == 1 {
			pos := qu.i
			oldVal := a[pos]
			newVal := qu.x
			for _, g := range glist {
				cnt := desc[g]
				left := pos > 0
				right := pos+1 < tc.n
				if left {
					lv := a[pos-1] % g
					if lv > oldVal%g {
						cnt--
					}
					if lv > newVal%g {
						cnt++
					}
				}
				if right {
					rv := a[pos+1] % g
					if oldVal%g > rv {
						cnt--
					}
					if newVal%g > rv {
						cnt++
					}
				}
				desc[g] = cnt
			}
			a[pos] = newVal
		} else {
			g := gcd(qu.k, tc.m)
			if g == 1 {
				ans = append(ans, "YES")
				continue
			}
			if desc[g] < tc.m/g {
				ans = append(ans, "YES")
			} else {
				ans = append(ans, "NO")
			}
		}
	}
	return ans
}

func gcd(a, b int) int {
	for b != 0 {
		a, b = b, a%b
	}
	if a < 0 {
		return -a
	}
	return a
}

func buildTests() []testCase {
	tests := []testCase{
		{
			name: "simple_yes_no",
			n:    3,
			m:    6,
			a:    []int{5, 1, 2},
			qs: []query{
				{t: 2, k: 2},
				{t: 2, k: 3},
			},
		},
		{
			name: "with_updates",
			n:    4,
			m:    10,
			a:    []int{9, 0, 5, 5},
			qs: []query{
				{t: 2, k: 5},
				{t: 1, i: 2, x: 3},
				{t: 2, k: 4},
				{t: 1, i: 0, x: 1},
				{t: 2, k: 6},
			},
		},
	}

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	totalN, totalQ := 0, 0
	for _, tc := range tests {
		totalN += tc.n
		totalQ += len(tc.qs)
	}

	for len(tests) < 60 && totalN < 100000 && totalQ < 100000 {
		n := rng.Intn(5000) + 2
		if totalN+n > 100000 {
			n = 100000 - totalN
		}
		m := rng.Intn(500000-2) + 2
		q := rng.Intn(4000) + 1
		if totalQ+q > 100000 {
			q = 100000 - totalQ
		}
		a := make([]int, n)
		for i := range a {
			a[i] = rng.Intn(m)
		}
		qs := make([]query, 0, q)
		for i := 0; i < q; i++ {
			if rng.Intn(3) == 0 {
				pos := rng.Intn(n)
				val := rng.Intn(m)
				qs = append(qs, query{t: 1, i: pos, x: val})
			} else {
				k := rng.Intn(m-1) + 1
				qs = append(qs, query{t: 2, k: k})
			}
		}
		tests = append(tests, testCase{
			name: fmt.Sprintf("rnd_%d", len(tests)+1),
			n:    n,
			m:    m,
			a:    a,
			qs:   qs,
		})
		totalN += n
		totalQ += q
	}

	return tests
}
