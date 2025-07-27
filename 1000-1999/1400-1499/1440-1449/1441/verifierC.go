package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

type arrCase struct {
	len  int
	vals []int64
}

type testCase struct {
	n   int
	k   int
	arr []arrCase
}

func (tc testCase) Input() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", tc.n, tc.k))
	for _, a := range tc.arr {
		sb.WriteString(strconv.Itoa(a.len))
		for _, v := range a.vals {
			sb.WriteByte(' ')
			sb.WriteString(strconv.FormatInt(v, 10))
		}
		sb.WriteByte('\n')
	}
	return sb.String()
}

func buildOracle() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	oracle := filepath.Join(dir, "oracleC")
	cmd := exec.Command("go", "build", "-o", oracle, "1441C.go")
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("build oracle failed: %v\n%s", err, out)
	}
	return oracle, nil
}

func genCase(rng *rand.Rand) testCase {
	n := rng.Intn(4) + 1
	k := rng.Intn(5) + 1
	arr := make([]arrCase, n)
	for i := 0; i < n; i++ {
		l := rng.Intn(5) + 1
		vals := make([]int64, l)
		prev := int64(rng.Intn(3))
		vals[0] = prev
		for j := 1; j < l; j++ {
			prev += int64(rng.Intn(3))
			vals[j] = prev
		}
		arr[i] = arrCase{l, vals}
	}
	return testCase{n, k, arr}
}

func runExe(path, input string) (string, error) {
	cmd := exec.Command(path)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	oracle, err := buildOracle()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer os.Remove(oracle)
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	cases := []testCase{
		{n: 1, k: 1, arr: []arrCase{{len: 1, vals: []int64{5}}}},
	}
	for i := 0; i < 100; i++ {
		cases = append(cases, genCase(rng))
	}
	for i, tc := range cases {
		exp, err := runExe(oracle, tc.Input())
		if err != nil {
			fmt.Fprintf(os.Stderr, "oracle error on case %d: %v\n", i+1, err)
			os.Exit(1)
		}
		got, err := runExe(bin, tc.Input())
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, tc.Input())
			os.Exit(1)
		}
		if got != exp {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, exp, got, tc.Input())
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
