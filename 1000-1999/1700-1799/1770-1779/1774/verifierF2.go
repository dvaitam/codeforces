package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

// Embedded reference solver for 1774F2.

const refMOD int64 = 998244353
const refCAP int64 = 1000000000

func refSolve(input string) int64 {
	data := []byte(input)
	pos := 0
	nextInt := func() int64 {
		for pos < len(data) && (data[pos] < '0' || data[pos] > '9') {
			pos++
		}
		var v int64
		for pos < len(data) && data[pos] >= '0' && data[pos] <= '9' {
			v = v*10 + int64(data[pos]-'0')
			pos++
		}
		return v
	}

	n := int(nextInt())
	typ := make([]int, n+1)
	val := make([]int64, n+1)

	for i := 1; i <= n; i++ {
		t := int(nextInt())
		typ[i] = t
		if t != 3 {
			val[i] = nextInt()
		}
	}

	pow2 := make([]int64, n+1)
	pow2[0] = 1
	for i := 1; i <= n; i++ {
		pow2[i] = (pow2[i-1] * 2) % refMOD
	}

	pref := make([]int64, n+1)
	for i := 1; i <= n; i++ {
		switch typ[i] {
		case 1:
			pref[i] = pref[i-1]
		case 2:
			s := pref[i-1] + val[i]
			if s > refCAP {
				s = refCAP
			}
			pref[i] = s
		case 3:
			s := pref[i-1] * 2
			if s > refCAP {
				s = refCAP
			}
			pref[i] = s
		}
	}

	countSubsets := func(weights []int64, limit int64) int64 {
		if limit < 0 {
			return 0
		}
		var cnt int64
		m := len(weights)
		for i, w := range weights {
			if limit >= w {
				cnt += int64(1) << uint(m-i-1)
				limit -= w
			}
		}
		return cnt + 1
	}

	var ans int64
	var suffixDamage int64
	zeroRepeats := 0
	weights := make([]int64, 0, 32)

	for i := n; i >= 1; i-- {
		if typ[i] == 1 {
			limit := val[i] - 1 - suffixDamage
			if limit >= 0 {
				c := countSubsets(weights, limit) % refMOD
				ans = (ans + c*pow2[zeroRepeats]) % refMOD
			}
		}

		switch typ[i] {
		case 2:
			suffixDamage += val[i]
			if suffixDamage > refCAP {
				suffixDamage = refCAP
			}
		case 3:
			t := pref[i-1]
			if t == 0 {
				zeroRepeats++
			} else if t < refCAP {
				weights = append(weights, t)
			}
		}
	}

	return ans % refMOD
}

type testCase struct {
	input string
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierF2.go /path/to/binary")
		os.Exit(1)
	}
	candidate := os.Args[1]

	tests := generateTests()
	for idx, tc := range tests {
		want := refSolve(tc.input)

		gotOut, err := runProgram(candidate, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate failed on test %d: %v\nInput:\n%s\n", idx+1, err, tc.input)
			os.Exit(1)
		}
		got, err := parseOutput(gotOut)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate produced invalid output on test %d: %v\nOutput:\n%s\n", idx+1, err, gotOut)
			os.Exit(1)
		}

		if want != got {
			fmt.Fprintf(os.Stderr, "test %d mismatch: expected %d, got %d\nInput:\n%s\nCandidate output:\n%s\n", idx+1, want, got, tc.input, gotOut)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed.\n", len(tests))
}

func runProgram(target, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(target, ".go") {
		cmd = exec.Command("go", "run", target)
	} else {
		cmd = exec.Command(target)
	}
	cmd.Stdin = strings.NewReader(input)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\nstderr:\n%s", err, stderr.String())
	}
	return stdout.String(), nil
}

func parseOutput(out string) (int64, error) {
	fields := strings.Fields(out)
	if len(fields) == 0 {
		return 0, fmt.Errorf("empty output")
	}
	val, err := strconv.ParseInt(fields[0], 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid integer %q", fields[0])
	}
	return val, nil
}

func generateTests() []testCase {
	var tests []testCase
	tests = append(tests,
		buildTest([]op{{typ: 1, val: 1}}),
		buildTest([]op{{typ: 2, val: 3}, {typ: 1, val: 2}}),
		buildTest([]op{{typ: 3}, {typ: 1, val: 10}, {typ: 2, val: 5}}),
	)
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for len(tests) < 50 {
		n := rng.Intn(50) + 1
		tests = append(tests, randomTest(rng, n))
	}
	tests = append(tests, randomTest(rng, 2000))
	return tests
}

type op struct {
	typ int
	val int64
}

func buildTest(ops []op) testCase {
	var b strings.Builder
	b.WriteString(fmt.Sprintf("%d\n", len(ops)))
	for _, op := range ops {
		if op.typ == 3 {
			b.WriteString(fmt.Sprintf("%d\n", op.typ))
		} else {
			b.WriteString(fmt.Sprintf("%d %d\n", op.typ, op.val))
		}
	}
	return testCase{input: b.String()}
}

func randomTest(rng *rand.Rand, n int) testCase {
	ops := make([]op, n)
	for i := 0; i < n; i++ {
		t := rng.Intn(3) + 1
		ops[i].typ = t
		if t != 3 {
			ops[i].val = rng.Int63n(1_000_000_000)
		}
	}
	return buildTest(ops)
}
