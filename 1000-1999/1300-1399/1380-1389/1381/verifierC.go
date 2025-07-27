package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

func buildIfGo(path string) (string, func(), error) {
	if strings.HasSuffix(path, ".go") {
		tmp, err := os.CreateTemp("", "solbin*")
		if err != nil {
			return "", nil, err
		}
		tmp.Close()
		if out, err := exec.Command("go", "build", "-o", tmp.Name(), path).CombinedOutput(); err != nil {
			os.Remove(tmp.Name())
			return "", nil, fmt.Errorf("build failed: %v\n%s", err, out)
		}
		return tmp.Name(), func() { os.Remove(tmp.Name()) }, nil
	}
	return path, func() {}, nil
}

func buildOracle() (string, error) {
	tmp, err := os.CreateTemp("", "oracle*")
	if err != nil {
		return "", err
	}
	tmp.Close()
	if out, err := exec.Command("go", "build", "-o", tmp.Name(), "1381C.go").CombinedOutput(); err != nil {
		os.Remove(tmp.Name())
		return "", fmt.Errorf("build oracle failed: %v\n%s", err, out)
	}
	return tmp.Name(), nil
}

func runProg(exe, input string) (string, error) {
	cmd := exec.Command(exe)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func parseCandidate(out string, n int) (bool, []int, error) {
	r := strings.NewReader(out)
	var token string
	if _, err := fmt.Fscan(r, &token); err != nil {
		return false, nil, err
	}
	token = strings.ToUpper(token)
	if token == "NO" {
		return false, nil, nil
	}
	arr := make([]int, n)
	for i := 0; i < n; i++ {
		if _, err := fmt.Fscan(r, &arr[i]); err != nil {
			return false, nil, err
		}
	}
	return true, arr, nil
}

func checkValidity(n, x, y int, b, ans []int) bool {
	if len(ans) != n {
		return false
	}
	cntB := make(map[int]int)
	cntA := make(map[int]int)
	matches := 0
	for i := 0; i < n; i++ {
		if ans[i] < 1 || ans[i] > n+1 {
			return false
		}
		if ans[i] == b[i] {
			matches++
		}
		cntA[ans[i]]++
		cntB[b[i]]++
	}
	if matches != x {
		return false
	}
	inter := 0
	for c, ca := range cntA {
		if cb, ok := cntB[c]; ok {
			if ca < cb {
				inter += ca
			} else {
				inter += cb
			}
		}
	}
	return inter == y
}

type testCase struct {
	n int
	x int
	y int
	b []int
}

func genTests() []testCase {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	tests := make([]testCase, 100)
	for i := 0; i < 100; i++ {
		n := rng.Intn(6) + 1
		x := rng.Intn(n + 1)
		y := x + rng.Intn(n-x+1)
		b := make([]int, n)
		for j := 0; j < n; j++ {
			b[j] = rng.Intn(n+1) + 1
		}
		tests[i] = testCase{n, x, y, b}
	}
	return tests
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	exe, cleanup, err := buildIfGo(os.Args[1])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer cleanup()
	oracle, err := buildOracle()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer os.Remove(oracle)

	tests := genTests()
	for i, tc := range tests {
		var sb strings.Builder
		sb.WriteString("1\n")
		sb.WriteString(fmt.Sprintf("%d %d %d\n", tc.n, tc.x, tc.y))
		for j, v := range tc.b {
			if j > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprintf("%d", v))
		}
		sb.WriteByte('\n')
		input := sb.String()
		want, err := runProg(oracle, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "oracle fail case %d: %v\ninput:%s", i+1, err, input)
			os.Exit(1)
		}
		wantHas := strings.HasPrefix(strings.TrimSpace(want), "YES")
		got, err := runProg(exe, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: %v\ninput:%s", i+1, err, input)
			os.Exit(1)
		}
		has, arr, err := parseCandidate(got, tc.n)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: parse error: %v\noutput:%s", i+1, err, got)
			os.Exit(1)
		}
		if wantHas != has {
			fmt.Fprintf(os.Stderr, "case %d: expected %v got %v\ninput:%s", i+1, wantHas, has, input)
			os.Exit(1)
		}
		if has {
			if !checkValidity(tc.n, tc.x, tc.y, tc.b, arr) {
				fmt.Fprintf(os.Stderr, "case %d: invalid solution\ninput:%soutput:%s", i+1, input, got)
				os.Exit(1)
			}
		}
	}
	fmt.Println("All tests passed")
}
