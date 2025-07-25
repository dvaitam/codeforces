package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

type testCaseB struct {
	n   int
	k   int64
	arr []int
	exp string
}

func solveB(n int, k int64, arr []int) string {
	current := arr[0]
	var wins int64
	for i := 1; i < n && wins < k; i++ {
		if arr[i] > current {
			current = arr[i]
			wins = 1
		} else {
			wins++
		}
	}
	return fmt.Sprint(current)
}

func generateTests() []testCaseB {
	rng := rand.New(rand.NewSource(2))
	cases := make([]testCaseB, 100)
	for i := range cases {
		n := rng.Intn(20) + 2
		k := int64(rng.Intn(2*n) + 1)
		perm := rng.Perm(n)
		for j := range perm {
			perm[j]++
		}
		cases[i] = testCaseB{n: n, k: k, arr: append([]int(nil), perm...), exp: solveB(n, k, perm)}
	}
	return cases
}

func run(bin string, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
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
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	cases := generateTests()
	for i, tc := range cases {
		var sb strings.Builder
		fmt.Fprintf(&sb, "%d %d\n", tc.n, tc.k)
		for j, v := range tc.arr {
			if j > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprint(v))
		}
		sb.WriteByte('\n')
		got, err := run(bin, sb.String())
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != tc.exp {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\n", i+1, tc.exp, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
