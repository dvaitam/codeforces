package main

import (
	"bytes"
	"fmt"
	"math"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

func isqrt(x int64) int64 {
	r := int64(math.Sqrt(float64(x)))
	for (r+1)*(r+1) <= x {
		r++
	}
	for r*r > x {
		r--
	}
	return r
}

func solveD(n int64) string {
	x := (1 + isqrt(1+8*n)) / 2
	for x*(x-1)/2 > n {
		x--
	}
	tri := x * (x - 1) / 2
	res := x + (n - tri)
	return fmt.Sprint(res)
}

func genCases() []string {
	rand.Seed(4)
	cases := make([]string, 100)
	for i := 0; i < 100; i++ {
		n := rand.Int63n(1_000_000_000_000) + 1
		sb := strings.Builder{}
		sb.WriteString("1\n")
		sb.WriteString(fmt.Sprint(n))
		sb.WriteByte('\n')
		cases[i] = sb.String()
	}
	return cases
}

func runCase(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	cases := genCases()
	for i, tc := range cases {
		var n int64
		fmt.Sscan(strings.TrimSpace(tc), &n)
		want := solveD(n)
		got, err := runCase(bin, tc)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Runtime error on case %d: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != want {
			fmt.Fprintf(os.Stderr, "Wrong answer on case %d\nInput:\n%sExpected: %s Got: %s\n", i+1, tc, want, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}
