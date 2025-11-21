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

type restriction struct {
	u int
	v int
	y int64
}

type testCase struct {
	n int
	m int
	e []restriction
}

const limitTime int64 = 1e18

func buildOracle() (string, func(), error) {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		return "", nil, fmt.Errorf("cannot determine verifier path")
	}
	dir := filepath.Dir(file)
	tmpDir, err := os.MkdirTemp("", "oracle-1842D-")
	if err != nil {
		return "", nil, err
	}
	path := filepath.Join(tmpDir, "oracleD")
	cmd := exec.Command("go", "build", "-o", path, "1842D.go")
	cmd.Dir = dir
	if out, err := cmd.CombinedOutput(); err != nil {
		os.RemoveAll(tmpDir)
		return "", nil, fmt.Errorf("failed to build oracle: %v\n%s", err, out)
	}
	cleanup := func() { os.RemoveAll(tmpDir) }
	return path, cleanup, nil
}

func runBinary(bin string, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(stdout.String()), nil
}

func deterministicTests() []testCase {
	return []testCase{
		{
			n: 4, m: 1,
			e: []restriction{{u: 2, v: 3, y: 1}},
		},
		{
			n: 3, m: 0,
			e: []restriction{},
		},
		{
			n: 5, m: 2,
			e: []restriction{
				{u: 1, v: 3, y: 2},
				{u: 2, v: 4, y: 1},
			},
		},
	}
}

func randomTest(rng *rand.Rand) testCase {
	n := rng.Intn(7) + 2
	pairs := make([][2]int, 0)
	for u := 1; u < n; u++ {
		for v := u + 1; v <= n; v++ {
			pairs = append(pairs, [2]int{u, v})
		}
	}
	rng.Shuffle(len(pairs), func(i, j int) { pairs[i], pairs[j] = pairs[j], pairs[i] })
	m := rng.Intn(len(pairs) + 1)
	edges := make([]restriction, 0, m)
	for i := 0; i < m; i++ {
		edges = append(edges, restriction{
			u: pairs[i][0],
			v: pairs[i][1],
			y: int64(rng.Intn(10) + 1),
		})
	}
	return testCase{n: n, m: len(edges), e: edges}
}

func buildInput(tc testCase) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", tc.n, tc.m))
	for _, r := range tc.e {
		sb.WriteString(fmt.Sprintf("%d %d %d\n", r.u, r.v, r.y))
	}
	return sb.String()
}

func parseOutput(out string, n int) (string, int64, []struct {
	mask string
	t    int64
}, error) {
	out = strings.TrimSpace(out)
	if out == "inf" {
		return "inf", 0, nil, nil
	}
	lines := strings.Split(out, "\n")
	if len(lines) == 0 {
		return "", 0, nil, fmt.Errorf("empty output")
	}
	header := strings.Fields(lines[0])
	if len(header) != 2 {
		return "", 0, nil, fmt.Errorf("invalid header line: %s", lines[0])
	}
	total, err := strconv.ParseInt(header[0], 10, 64)
	if err != nil {
		return "", 0, nil, fmt.Errorf("invalid total time: %v", err)
	}
	k, err := strconv.Atoi(header[1])
	if err != nil {
		return "", 0, nil, fmt.Errorf("invalid k: %v", err)
	}
	if k < 0 || k > n*n {
		return "", 0, nil, fmt.Errorf("k out of bounds")
	}
	if len(lines)-1 != k {
		return "", 0, nil, fmt.Errorf("expected %d games, got %d", k, len(lines)-1)
	}
	if total < 0 || total > limitTime {
		return "", 0, nil, fmt.Errorf("total time out of bounds")
	}
	games := make([]struct {
		mask string
		t    int64
	}, k)
	for i := 0; i < k; i++ {
		fields := strings.Fields(lines[i+1])
		if len(fields) != 2 {
			return "", 0, nil, fmt.Errorf("invalid game line: %s", lines[i+1])
		}
		mask := fields[0]
		if len(mask) != n {
			return "", 0, nil, fmt.Errorf("mask length mismatch")
		}
		for _, ch := range mask {
			if ch != '0' && ch != '1' {
				return "", 0, nil, fmt.Errorf("mask contains invalid character")
			}
		}
		dur, err := strconv.ParseInt(fields[1], 10, 64)
		if err != nil {
			return "", 0, nil, fmt.Errorf("invalid duration: %v", err)
		}
		if dur < 0 || dur > limitTime {
			return "", 0, nil, fmt.Errorf("duration out of bounds")
		}
		games[i] = struct {
			mask string
			t    int64
		}{mask: mask, t: dur}
	}
	return "finite", total, games, nil
}

func simulate(tc testCase, plan []struct {
	mask string
	t    int64
}) (bool, string, int64) {
	type edge struct {
		u, v int
		y    int64
	}
	limits := make(map[[2]int]edge, len(tc.e))
	for _, r := range tc.e {
		key := [2]int{r.u, r.v}
		limits[key] = edge{u: r.u, v: r.v, y: r.y}
	}
	n := tc.n
	total := int64(0)
	for _, g := range plan {
		if g.t < 0 {
			return false, "negative duration", total
		}
		if len(g.mask) != n {
			return false, "mask length mismatch", total
		}
		if g.mask[0] != '1' {
			return false, "friend 1 not included", total
		}
		if g.mask[n-1] != '0' {
			return false, "friend n included", total
		}
		for key, e := range limits {
			u := e.u
			v := e.v
			if g.mask[u-1] != g.mask[v-1] {
				if g.t > e.y {
					return false, fmt.Sprintf("restriction exceeded for pair (%d,%d)", u, v), total
				}
				e.y -= g.t
				limits[key] = e
			}
		}
		total += g.t
	}
	return true, "", total
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	target := os.Args[1]

	oracle, cleanup, err := buildOracle()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	defer cleanup()

	tests := deterministicTests()
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 200; i++ {
		tests = append(tests, randomTest(rng))
	}

	for idx, tc := range tests {
		input := buildInput(tc)

		expOut, err := runBinary(oracle, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "oracle failed on test %d: %v\n", idx+1, err)
			os.Exit(1)
		}
		expMode, expTotal, _, err := parseOutput(expOut, tc.n)
		if err != nil {
			fmt.Fprintf(os.Stderr, "invalid oracle output on test %d: %v\n", idx+1, err)
			os.Exit(1)
		}

		gotOut, err := runBinary(target, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "test %d: %v\n", idx+1, err)
			os.Exit(1)
		}
		gotMode, gotTotal, gotPlan, err := parseOutput(gotOut, tc.n)
		if err != nil {
			fmt.Fprintf(os.Stderr, "test %d: %v\n", idx+1, err)
			os.Exit(1)
		}

		if expMode == "inf" {
			if gotMode != "inf" {
				fmt.Fprintf(os.Stderr, "test %d: expected 'inf' but got finite output\n", idx+1)
				os.Exit(1)
			}
			continue
		}
		if gotMode == "inf" {
			fmt.Fprintf(os.Stderr, "test %d: unexpected 'inf' output\n", idx+1)
			os.Exit(1)
		}
		if gotTotal != expTotal {
			fmt.Fprintf(os.Stderr, "test %d: expected total %d got %d\n", idx+1, expTotal, gotTotal)
			os.Exit(1)
		}
		ok, msg, simTotal := simulate(tc, gotPlan)
		if !ok {
			fmt.Fprintf(os.Stderr, "test %d: %s\n", idx+1, msg)
			os.Exit(1)
		}
		if simTotal != gotTotal {
			fmt.Fprintf(os.Stderr, "test %d: reported total %d but simulated %d\n", idx+1, gotTotal, simTotal)
			os.Exit(1)
		}
	}

	fmt.Printf("All %d tests passed\n", len(tests))
}
