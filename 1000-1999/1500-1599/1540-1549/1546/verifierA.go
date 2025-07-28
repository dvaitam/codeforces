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

type testCase struct {
	n        int
	a        []int
	b        []int
	possible bool
}

func runCandidate(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return out.String() + errBuf.String(), fmt.Errorf("runtime error: %v", err)
	}
	return out.String(), nil
}

func genCase(rng *rand.Rand) testCase {
	n := rng.Intn(6) + 1 // 1..6
	a := make([]int, n)
	b := make([]int, n)
	for i := 0; i < n; i++ {
		a[i] = rng.Intn(6)
		b[i] = a[i]
	}
	possible := rng.Intn(2) == 0
	if possible {
		ops := rng.Intn(5*n + 1)
		for k := 0; k < ops; k++ {
			i := rng.Intn(n)
			j := rng.Intn(n)
			if b[i] > 0 {
				b[i]--
				b[j]++
			}
		}
	} else {
		diff := rng.Intn(5) + 1
		if rng.Intn(2) == 0 {
			idx := rng.Intn(n)
			a[idx] += diff
		} else {
			idx := rng.Intn(n)
			b[idx] += diff
		}
	}
	return testCase{n: n, a: a, b: b, possible: possible}
}

func verify(tc testCase, out string) error {
	tokens := strings.Fields(strings.TrimSpace(out))
	sumA, sumB := 0, 0
	for _, v := range tc.a {
		sumA += v
	}
	for _, v := range tc.b {
		sumB += v
	}
	if sumA != sumB {
		if len(tokens) == 0 || tokens[0] != "-1" {
			return fmt.Errorf("expected -1 for impossible case, got %q", out)
		}
		return nil
	}
	if len(tokens) == 0 {
		return fmt.Errorf("empty output")
	}
	m, err := strconv.Atoi(tokens[0])
	if err != nil {
		return fmt.Errorf("cannot parse operation count: %v", err)
	}
	if m < 0 || m > 100 {
		return fmt.Errorf("invalid operation count %d", m)
	}
	if len(tokens) != 1+2*m {
		return fmt.Errorf("expected %d numbers, got %d", 1+2*m, len(tokens))
	}
	arr := make([]int, len(tc.a))
	copy(arr, tc.a)
	idx := 1
	for k := 0; k < m; k++ {
		from, err1 := strconv.Atoi(tokens[idx])
		to, err2 := strconv.Atoi(tokens[idx+1])
		idx += 2
		if err1 != nil || err2 != nil {
			return fmt.Errorf("invalid integers in operation %d", k+1)
		}
		if from < 1 || from > tc.n || to < 1 || to > tc.n {
			return fmt.Errorf("indices out of range in operation %d", k+1)
		}
		arr[from-1]--
		arr[to-1]++
		if arr[from-1] < 0 {
			return fmt.Errorf("negative value after operation %d", k+1)
		}
	}
	for i := 0; i < tc.n; i++ {
		if arr[i] != tc.b[i] {
			return fmt.Errorf("arrays not equal after operations; expected %v got %v", tc.b, arr)
		}
	}
	return nil
}

func main() {
	var bin string
	if len(os.Args) == 2 {
		bin = os.Args[1]
	} else if len(os.Args) == 3 && os.Args[1] == "--" {
		bin = os.Args[2]
	} else {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	var cases []testCase
	// deterministic edge cases
	cases = append(cases, testCase{n: 1, a: []int{0}, b: []int{1}, possible: false})
	cases = append(cases, testCase{n: 1, a: []int{2}, b: []int{2}, possible: true})
	for len(cases) < 100 {
		cases = append(cases, genCase(rng))
	}
	for i, tc := range cases {
		var input strings.Builder
		fmt.Fprintln(&input, 1)
		fmt.Fprintln(&input, tc.n)
		for j, v := range tc.a {
			if j > 0 {
				input.WriteByte(' ')
			}
			input.WriteString(strconv.Itoa(v))
		}
		input.WriteByte('\n')
		for j, v := range tc.b {
			if j > 0 {
				input.WriteByte(' ')
			}
			input.WriteString(strconv.Itoa(v))
		}
		input.WriteByte('\n')
		out, err := runCandidate(bin, input.String())
		if err != nil {
			fmt.Fprintf(os.Stderr, "test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		if err := verify(tc, out); err != nil {
			fmt.Fprintf(os.Stderr, "test %d failed: %v\ninput:\n%soutput:\n%s\n", i+1, err, input.String(), out)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}
