package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/bits"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

type runResult struct {
	out string
	err error
}

func runBinary(path, input string) runResult {
	cmd := exec.Command(path)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		err = fmt.Errorf("%v: %s", err, stderr.String())
	}
	return runResult{strings.TrimSpace(out.String()), err}
}

func expectedOps(arr []int64) int {
	var minA, maxA int64 = arr[0], arr[0]
	for _, v := range arr[1:] {
		if v < minA {
			minA = v
		}
		if v > maxA {
			maxA = v
		}
	}
	d := maxA - minA
	if d == 0 {
		return 0
	}
	return bits.Len64(uint64(d))
}

func parseInts(out string) ([]int64, error) {
	s := bufio.NewScanner(strings.NewReader(out))
	s.Split(bufio.ScanWords)
	res := make([]int64, 0)
	for s.Scan() {
		tok := s.Text()
		v, err := strconv.ParseInt(tok, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid integer token %q", tok)
		}
		res = append(res, v)
	}
	if err := s.Err(); err != nil {
		return nil, err
	}
	return res, nil
}

func validateOutput(arr []int64, out string) error {
	toks, err := parseInts(out)
	if err != nil {
		return err
	}
	if len(toks) == 0 {
		return fmt.Errorf("empty output")
	}
	kExp := expectedOps(arr)
	k := toks[0]
	if k < 0 {
		return fmt.Errorf("k must be non-negative, got %d", k)
	}
	if int(k) != kExp {
		return fmt.Errorf("non-minimal/incorrect operation count: got %d, expected %d", k, kExp)
	}
	if k == 0 || int(k) > len(arr) {
		return nil
	}
	if len(toks) < 1+int(k) {
		return fmt.Errorf("expected %d operation values, got %d", k, len(toks)-1)
	}
	cur := append([]int64(nil), arr...)
	for i := 0; i < int(k); i++ {
		x := toks[1+i]
		if x < 0 || x > 1_000_000_000_000_000_000 {
			return fmt.Errorf("x[%d]=%d out of range [0, 1e18]", i, x)
		}
		for j := range cur {
			cur[j] = (cur[j] + x) / 2
		}
	}
	for i := 1; i < len(cur); i++ {
		if cur[i] != cur[0] {
			return fmt.Errorf("after %d operations array is not equalized: %v", k, cur)
		}
	}
	return nil
}

func genTest() string {
	n := rand.Intn(30) + 1
	arr := make([]int64, n)
	for i := 0; i < n; i++ {
		switch rand.Intn(5) {
		case 0:
			arr[i] = 0
		case 1:
			arr[i] = 1_000_000_000
		default:
			arr[i] = rand.Int63n(1_000_000_001)
		}
	}
	var sb strings.Builder
	fmt.Fprintf(&sb, "1\n%d\n", n)
	for i, v := range arr {
		if i > 0 {
			sb.WriteByte(' ')
		}
		fmt.Fprintf(&sb, "%d", v)
	}
	sb.WriteByte('\n')
	return sb.String()
}

func parseInputArray(input string) ([]int64, error) {
	toks, err := parseInts(input)
	if err != nil {
		return nil, err
	}
	if len(toks) < 3 || toks[0] != 1 {
		return nil, fmt.Errorf("unexpected generated input format")
	}
	n := int(toks[1])
	if len(toks) != 2+n {
		return nil, fmt.Errorf("bad generated input, n=%d tokens=%d", n, len(toks))
	}
	return append([]int64(nil), toks[2:]...), nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	binary := os.Args[1]
	rand.Seed(3)
	for i := 0; i < 200; i++ {
		tc := genTest()
		arr, err := parseInputArray(tc)
		if err != nil {
			fmt.Fprintf(os.Stderr, "internal generator error on test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		got := runBinary(binary, tc)
		if got.err != nil {
			fmt.Fprintf(os.Stderr, "binary failed on test %d: %v\n", i+1, got.err)
			os.Exit(1)
		}
		if err := validateOutput(arr, got.out); err != nil {
			fmt.Printf("wrong answer on test %d\ninput:\n%s\noutput:\n%s\nreason: %v\n", i+1, tc, got.out, err)
			os.Exit(1)
		}
	}
	fmt.Println("all tests passed")
}
