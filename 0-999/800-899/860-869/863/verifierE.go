package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func genTestsE() []string {
	rand.Seed(5)
	tests := make([]string, 0, 100)
	for len(tests) < 100 {
		n := rand.Intn(8) + 2
		l := make([]int, n)
		r := make([]int, n)
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d\n", n))
		for i := 0; i < n; i++ {
			l[i] = rand.Intn(100) + 1
			r[i] = l[i] + rand.Intn(20)
			sb.WriteString(fmt.Sprintf("%d %d\n", l[i], r[i]))
		}
		tests = append(tests, sb.String())
	}
	return tests
}

func runBinary(bin string, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func isRedundant(idx int, l, r []int) bool {
	// Coverage map for integer coordinates
	cnt := make(map[int]int)
	n := len(l)
	for i := 0; i < n; i++ {
		for k := l[i]; k <= r[i]; k++ {
			cnt[k]++
		}
	}

	// Check if removing interval idx leaves every integer point in it covered
	for k := l[idx]; k <= r[idx]; k++ {
		if cnt[k] <= 1 {
			return false
		}
	}
	return true
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "Usage: go run verifierE.go <binary>\n")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := genTestsE()
	for idx, t := range tests {
		lines := strings.Split(strings.TrimSpace(t), "\n")
		var n int
		fmt.Sscanf(lines[0], "%d", &n)
		l := make([]int, n)
		r := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Sscanf(lines[1+i], "%d %d", &l[i], &r[i])
		}

		got, err := runBinary(bin, t)
		if err != nil {
			fmt.Printf("Test %d: runtime error: %v\n", idx+1, err)
			os.Exit(1)
		}

		gotInt, err := strconv.Atoi(got)
		if err != nil {
			fmt.Printf("Test %d failed. Invalid output format: '%s'\n", idx+1, got)
			os.Exit(1)
		}

		if gotInt == -1 {
			// Check if we missed a valid solution
			for i := 0; i < n; i++ {
				if isRedundant(i, l, r) {
					fmt.Printf("Test %d failed.\nInput:\n%s\nGot: -1\nExpected: %d (interval %d is redundant)\n", idx+1, t, i+1, i+1)
					os.Exit(1)
				}
			}
		} else {
			// Check if returned index is valid
			if gotInt < 1 || gotInt > n {
				fmt.Printf("Test %d failed. Output index out of bounds: %d (n=%d)\n", idx+1, gotInt, n)
				os.Exit(1)
			}
			if !isRedundant(gotInt-1, l, r) {
				fmt.Printf("Test %d failed.\nInput:\n%s\nGot: %d\nError: Interval %d is NOT redundant.\n", idx+1, t, gotInt, gotInt)
				os.Exit(1)
			}
		}
	}
	fmt.Printf("All %d tests passed.\n", len(tests))
}

