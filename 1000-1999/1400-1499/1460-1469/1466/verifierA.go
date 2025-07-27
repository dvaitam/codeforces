package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

// solveA computes expected output for one test case.
func solveA(n int, xs []int) string {
	m := make(map[int]struct{})
	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			d := xs[j] - xs[i]
			m[d] = struct{}{}
		}
	}
	return fmt.Sprint(len(m))
}

// genCases generates at least 100 random test cases.
func genCases() []string {
	rand.Seed(1)
	cases := make([]string, 100)
	for i := 0; i < 100; i++ {
		n := rand.Intn(10) + 3 // 3..12
		xs := make([]int, n)
		used := make(map[int]bool)
		for j := 0; j < n; j++ {
			x := rand.Intn(50) + 1
			for used[x] {
				x = rand.Intn(50) + 1
			}
			used[x] = true
			xs[j] = x
		}
		// sort xs
		for j := 0; j < n; j++ {
			for k := j + 1; k < n; k++ {
				if xs[k] < xs[j] {
					xs[k], xs[j] = xs[j], xs[k]
				}
			}
		}
		sb := strings.Builder{}
		sb.WriteString("1\n")
		sb.WriteString(fmt.Sprint(n))
		sb.WriteByte('\n')
		for j := 0; j < n; j++ {
			if j > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprint(xs[j]))
		}
		sb.WriteByte('\n')
		cases[i] = sb.String()
	}
	return cases
}

func runCase(bin string, input string) (string, error) {
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
		fmt.Println("Usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	cases := genCases()
	for i, tc := range cases {
		// parse to compute expected answer
		lines := strings.Split(strings.TrimSpace(tc), "\n")
		var n int
		fmt.Sscan(lines[1], &n)
		xs := make([]int, n)
		parts := strings.Fields(lines[2])
		for j := 0; j < n; j++ {
			fmt.Sscan(parts[j], &xs[j])
		}
		want := solveA(n, xs)
		got, err := runCase(bin, tc)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Runtime error on case %d: %v\n", i+1, err)
			os.Exit(1)
		}
		if got != want {
			fmt.Fprintf(os.Stderr, "Wrong answer on case %d\nInput:\n%sExpected: %s Got: %s\n", i+1, tc, want, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}
