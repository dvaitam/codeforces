package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

type testCaseB struct {
	n   int
	arr []int
}

func isHillOrValley(arr []int, i int) bool {
	if i <= 0 || i >= len(arr)-1 {
		return false
	}
	return (arr[i] > arr[i-1] && arr[i] > arr[i+1]) || (arr[i] < arr[i-1] && arr[i] < arr[i+1])
}

func solveB(n int, a []int) int {
	if n <= 2 {
		return 0
	}
	hv := make([]bool, n)
	total := 0
	for i := 1; i < n-1; i++ {
		hv[i] = isHillOrValley(a, i)
		if hv[i] {
			total++
		}
	}
	ans := total
	for i := 0; i < n; i++ {
		orig := a[i]
		before := 0
		for j := i - 1; j <= i+1; j++ {
			if j > 0 && j < n-1 && hv[j] {
				before++
			}
		}
		if i > 0 {
			a[i] = a[i-1]
			after := 0
			for j := i - 1; j <= i+1; j++ {
				if j > 0 && j < n-1 && isHillOrValley(a, j) {
					after++
				}
			}
			if total-before+after < ans {
				ans = total - before + after
			}
		}
		if i < n-1 {
			a[i] = a[i+1]
			after := 0
			for j := i - 1; j <= i+1; j++ {
				if j > 0 && j < n-1 && isHillOrValley(a, j) {
					after++
				}
			}
			if total-before+after < ans {
				ans = total - before + after
			}
		}
		a[i] = orig
	}
	return ans
}

func buildInputB(tc testCaseB) string {
	var sb strings.Builder
	sb.WriteString("1\n")
	sb.WriteString(fmt.Sprint(tc.n))
	sb.WriteByte('\n')
	for i, v := range tc.arr {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprint(v))
	}
	sb.WriteByte('\n')
	return sb.String()
}

func runCaseB(bin string, tc testCaseB) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(buildInputB(tc))
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	var got int
	if _, err := fmt.Fscan(strings.NewReader(out.String()), &got); err != nil {
		return fmt.Errorf("bad output: %v", err)
	}
	expect := solveB(tc.n, append([]int(nil), tc.arr...))
	if got != expect {
		return fmt.Errorf("expected %d got %d", expect, got)
	}
	return nil
}

func generateCasesB() []testCaseB {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	cases := make([]testCaseB, 0, 100)
	cases = append(cases, testCaseB{n: 1, arr: []int{5}}, testCaseB{n: 2, arr: []int{1, 2}})
	for len(cases) < 100 {
		n := rng.Intn(20) + 1
		arr := make([]int, n)
		for i := 0; i < n; i++ {
			arr[i] = rng.Intn(30)
		}
		cases = append(cases, testCaseB{n: n, arr: arr})
	}
	return cases
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	cases := generateCasesB()
	for i, tc := range cases {
		if err := runCaseB(bin, tc); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v (n=%d)\n", i+1, err, tc.n)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
