package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"
	"time"
)

type testCase struct {
	n int
	p []int
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
		return "", fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func verifyCase(in, out []int) error {
	n := len(in)
	if len(out) != n {
		return fmt.Errorf("expected %d numbers, got %d", n, len(out))
	}
	used := make([]bool, n+1)
	for _, v := range out {
		if v < 1 || v > n {
			return fmt.Errorf("value %d out of range 1..%d", v, n)
		}
		if used[v] {
			return fmt.Errorf("value %d repeated", v)
		}
		used[v] = true
	}
	same := true
	for i := 0; i < n; i++ {
		if in[i] != out[i] {
			same = false
			break
		}
	}
	if same {
		return fmt.Errorf("output permutation equals input")
	}
	if n > 1 {
		inFp := make([]int, n-1)
		outFp := make([]int, n-1)
		for i := 0; i < n-1; i++ {
			inFp[i] = in[i] + in[i+1]
			outFp[i] = out[i] + out[i+1]
		}
		sort.Ints(inFp)
		sort.Ints(outFp)
		for i := range inFp {
			if inFp[i] != outFp[i] {
				return fmt.Errorf("fingerprint mismatch")
			}
		}
	}
	return nil
}

func genCase(rng *rand.Rand) testCase {
	n := rng.Intn(50) + 2
	p := rng.Perm(n)
	for i := range p {
		p[i]++
	}
	return testCase{n: n, p: p}
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	cases := []testCase{{n: 2, p: []int{1, 2}}, {n: 3, p: []int{3, 1, 2}}}
	for len(cases) < 100 {
		cases = append(cases, genCase(rng))
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", len(cases)))
	for _, tc := range cases {
		sb.WriteString(fmt.Sprintf("%d\n", tc.n))
		for i := 0; i < tc.n; i++ {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprint(tc.p[i]))
		}
		sb.WriteByte('\n')
	}
	out, err := runCandidate(bin, sb.String())
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	fields := strings.Fields(out)
	pos := 0
	for idx, tc := range cases {
		if pos+tc.n > len(fields) {
			fmt.Fprintf(os.Stderr, "case %d: expected %d numbers, got %d\n", idx+1, tc.n, len(fields)-pos)
			os.Exit(1)
		}
		arr := make([]int, tc.n)
		for i := 0; i < tc.n; i++ {
			v, err := strconv.Atoi(fields[pos])
			if err != nil {
				fmt.Fprintf(os.Stderr, "case %d: invalid integer %q\n", idx+1, fields[pos])
				os.Exit(1)
			}
			arr[i] = v
			pos++
		}
		if err := verifyCase(tc.p, arr); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput: %v\noutput: %v\n", idx+1, err, tc.p, arr)
			os.Exit(1)
		}
	}
	if pos != len(fields) {
		fmt.Fprintf(os.Stderr, "extra output tokens: %d\n", len(fields)-pos)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}
