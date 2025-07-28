package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

// solveA implements the official solution for 1668A.
func solveA(input string) string {
	in := bufio.NewReader(strings.NewReader(input))
	var t int
	fmt.Fscan(in, &t)
	var sb strings.Builder
	for ; t > 0; t-- {
		var n, m int64
		fmt.Fscan(in, &n, &m)
		if n == 1 && m == 1 {
			sb.WriteString("0\n")
			continue
		}
		if n == 1 || m == 1 {
			if (n == 1 && m == 2) || (m == 1 && n == 2) {
				sb.WriteString("1\n")
			} else {
				sb.WriteString("-1\n")
			}
			continue
		}
		if n < m {
			n, m = m, n
		}
		diff := n - m
		if diff%2 == 0 {
			sb.WriteString(fmt.Sprintf("%d\n", 2*n-2))
		} else {
			sb.WriteString(fmt.Sprintf("%d\n", 2*n-3))
		}
	}
	return sb.String()
}

func buildCaseA(n, m int64) string {
	return fmt.Sprintf("1\n%d %d\n", n, m)
}

func generateRandomCaseA(rng *rand.Rand) string {
	n := rng.Int63n(1e9) + 1
	m := rng.Int63n(1e9) + 1
	return buildCaseA(n, m)
}

func runBinary(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	err := cmd.Run()
	if err != nil {
		return out.String(), fmt.Errorf("%v\n%s", err, errBuf.String())
	}
	return out.String(), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	var cases []string
	// predetermined edge cases
	cases = append(cases, buildCaseA(1, 1))
	cases = append(cases, buildCaseA(1, 2))
	cases = append(cases, buildCaseA(2, 1))
	cases = append(cases, buildCaseA(1, 3))
	cases = append(cases, buildCaseA(3, 1))
	cases = append(cases, buildCaseA(5, 5))
	cases = append(cases, buildCaseA(10, 7))
	cases = append(cases, buildCaseA(9, 6))
	cases = append(cases, buildCaseA(2, 2))
	cases = append(cases, buildCaseA(1000000000, 1000000000))
	for len(cases) < 100 {
		cases = append(cases, generateRandomCaseA(rng))
	}
	for i, tc := range cases {
		expect := strings.TrimSpace(solveA(tc))
		got, err := runBinary(bin, tc)
		if err != nil {
			fmt.Printf("case %d: runtime error: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != expect {
			fmt.Printf("case %d failed\ninput:\n%sexpected:%sq\ngot:%sq\n", i+1, tc, expect, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
