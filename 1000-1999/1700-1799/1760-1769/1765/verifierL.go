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
	"strconv"
	"strings"
	"time"
)

var dayNames = []string{"Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday", "Sunday"}

type testCase struct {
	n        int
	m        int
	k        int
	workDays [][]int
	holidays []int
	projects [][]int
}

func buildOracle() (string, func(), error) {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		return "", nil, fmt.Errorf("cannot determine verifier path")
	}
	dir := filepath.Dir(file)
	tmpDir, err := os.MkdirTemp("", "oracle-1765L-")
	if err != nil {
		return "", nil, err
	}
	path := filepath.Join(tmpDir, "oracleL")
	cmd := exec.Command("go", "build", "-o", path, "1765L.go")
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

func parseOutput(out string, k int) ([]int, error) {
	fields := strings.Fields(out)
	if len(fields) != k {
		return nil, fmt.Errorf("expected %d numbers, got %d", k, len(fields))
	}
	res := make([]int, k)
	for i, f := range fields {
		val, err := strconv.Atoi(f)
		if err != nil {
			return nil, fmt.Errorf("invalid integer %q", f)
		}
		if val < 0 {
			return nil, fmt.Errorf("negative completion day")
		}
		res[i] = val
	}
	return res, nil
}

func buildInput(tc testCase) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d %d\n", tc.n, tc.m, tc.k))
	for i := 0; i < tc.n; i++ {
		sb.WriteString(strconv.Itoa(len(tc.workDays[i])))
		for _, d := range tc.workDays[i] {
			sb.WriteString(" ")
			sb.WriteString(dayNames[d])
		}
		sb.WriteByte('\n')
	}
	if tc.m > 0 {
		for i := 0; i < tc.m; i++ {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.Itoa(tc.holidays[i]))
		}
		sb.WriteByte('\n')
	} else {
		sb.WriteByte('\n')
	}
	for i := 0; i < tc.k; i++ {
		sb.WriteString(strconv.Itoa(len(tc.projects[i])))
		for _, v := range tc.projects[i] {
			sb.WriteString(" ")
			sb.WriteString(strconv.Itoa(v))
		}
		sb.WriteByte('\n')
	}
	return sb.String()
}

func deterministicTests() []testCase {
	return []testCase{
		{
			n: 1, m: 0, k: 1,
			workDays: [][]int{{0}},
			holidays: []int{},
			projects: [][]int{{1}},
		},
		{
			n: 2, m: 1, k: 2,
			workDays: [][]int{{0, 2, 4}, {1, 3, 5}},
			holidays: []int{3},
			projects: [][]int{{1, 2}, {2}},
		},
		{
			n: 3, m: 2, k: 2,
			workDays: [][]int{{0, 1, 2}, {3, 4}, {5, 6}},
			holidays: []int{4, 8},
			projects: [][]int{{1, 2, 3}, {3, 2}},
		},
	}
}

func randomTest(rng *rand.Rand) testCase {
	n := rng.Intn(4) + 1
	m := rng.Intn(4)
	k := rng.Intn(4) + 1
	workDays := make([][]int, n)
	for i := 0; i < n; i++ {
		t := rng.Intn(7) + 1
		set := make(map[int]struct{})
		for len(set) < t {
			set[rng.Intn(7)] = struct{}{}
		}
		arr := make([]int, 0, t)
		for d := range set {
			arr = append(arr, d)
		}
		sort.Ints(arr)
		workDays[i] = arr
	}
	holidays := make([]int, 0, m)
	if m > 0 {
		used := make(map[int]struct{})
		for len(holidays) < m {
			val := rng.Intn(30) + 1
			if _, ok := used[val]; ok {
				continue
			}
			used[val] = struct{}{}
			holidays = append(holidays, val)
		}
		sort.Ints(holidays)
	}
	projects := make([][]int, k)
	totalParts := 0
	for i := 0; i < k; i++ {
		maxParts := 3
		if totalParts >= 10 {
			maxParts = 1
		}
		p := rng.Intn(maxParts) + 1
		totalParts += p
		arr := make([]int, p)
		for j := 0; j < p; j++ {
			arr[j] = rng.Intn(n) + 1
		}
		projects[i] = arr
	}
	return testCase{n: n, m: len(holidays), k: k, workDays: workDays, holidays: holidays, projects: projects}
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierL.go /path/to/binary")
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
			fmt.Fprintf(os.Stderr, "oracle failed on test %d: %v\ninput:\n%s\n", idx+1, err, input)
			os.Exit(1)
		}
		expVals, err := parseOutput(expOut, tc.k)
		if err != nil {
			fmt.Fprintf(os.Stderr, "invalid oracle output on test %d: %v\noutput:\n%s\n", idx+1, err, expOut)
			os.Exit(1)
		}

		gotOut, err := runBinary(target, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "target runtime error on test %d: %v\ninput:\n%s\n", idx+1, err, input)
			os.Exit(1)
		}
		gotVals, err := parseOutput(gotOut, tc.k)
		if err != nil {
			fmt.Fprintf(os.Stderr, "target output invalid on test %d: %v\noutput:\n%s\ninput:\n%s\n", idx+1, err, gotOut, input)
			os.Exit(1)
		}

		for i := 0; i < tc.k; i++ {
			if expVals[i] != gotVals[i] {
				fmt.Fprintf(os.Stderr, "mismatch on test %d project %d: expected %d got %d\ninput:\n%s\n", idx+1, i+1, expVals[i], gotVals[i], input)
				os.Exit(1)
			}
		}
	}

	fmt.Printf("All %d tests passed\n", len(tests))
}
