package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strings"
)

type testCaseA struct {
	n     int
	coins []int
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rand.Seed(42)
	var tests []testCaseA
	// generate 100 random tests
	for i := 0; i < 100; i++ {
		n := rand.Intn(100) + 1
		coins := make([]int, n)
		for j := range coins {
			coins[j] = rand.Intn(100) + 1
		}
		tests = append(tests, testCaseA{n, coins})
	}
	// add some edge cases
	tests = append(tests, testCaseA{1, []int{1}})
	tests = append(tests, testCaseA{2, []int{1, 1}})

	for i, t := range tests {
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d\n", t.n))
		for j, v := range t.coins {
			if j > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprint(v))
		}
		sb.WriteByte('\n')
		input := sb.String()
		expect := solveA(strings.NewReader(input))
		out, err := runBinary(bin, input)
		if err != nil {
			fmt.Printf("test %d: execution failed: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(out) != strings.TrimSpace(expect) {
			fmt.Printf("test %d failed: expected %q got %q\n", i+1, strings.TrimSpace(expect), strings.TrimSpace(out))
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

func runBinary(path, input string) (string, error) {
	cmd := exec.Command(path)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	return out.String(), err
}

func solveA(r io.Reader) string {
	in := bufio.NewReader(r)
	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return ""
	}
	coins := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &coins[i])
	}
	sort.Sort(sort.Reverse(sort.IntSlice(coins)))
	total := 0
	for _, v := range coins {
		total += v
	}
	sum := 0
	for i, v := range coins {
		sum += v
		if sum > total-sum {
			return fmt.Sprintf("%d\n", i+1)
		}
	}
	return ""
}
