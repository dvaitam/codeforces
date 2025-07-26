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

func expectedA(a []int) int {
	n := len(a)
	pre1 := make([]int, n+1)
	for i := 1; i <= n; i++ {
		pre1[i] = pre1[i-1]
		if a[i-1] == 1 {
			pre1[i]++
		}
	}
	suf2 := make([]int, n+2)
	for i := n; i >= 1; i-- {
		suf2[i] = suf2[i+1]
		if a[i-1] == 2 {
			suf2[i]++
		}
	}
	ans := 0
	for r := 1; r <= n; r++ {
		dp1, dp2 := 0, 0
		for l := r; l >= 1; l-- {
			x := a[l-1]
			if x == 1 {
				dp1++
				if dp2 < dp1 {
					dp2 = dp1
				}
			} else {
				if dp2+1 > dp1+1 {
					dp2 = dp2 + 1
				} else {
					dp2 = dp1 + 1
				}
			}
			cand := pre1[l-1] + dp2 + suf2[r+1]
			if cand > ans {
				ans = cand
			}
		}
	}
	return ans
}

func generateCaseA(rng *rand.Rand) []int {
	n := rng.Intn(2000) + 1
	arr := make([]int, n)
	for i := range arr {
		if rng.Intn(2) == 0 {
			arr[i] = 1
		} else {
			arr[i] = 2
		}
	}
	return arr
}

func runCaseA(bin string, arr []int) error {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", len(arr)))
	for i, v := range arr {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprint(v))
	}
	sb.WriteByte('\n')
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(sb.String())
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
	exp := expectedA(arr)
	if got != exp {
		return fmt.Errorf("expected %d got %d", exp, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	edge := [][]int{
		{1}, {2}, {1, 1}, {2, 2}, {1, 2}, {2, 1},
	}
	for idx, arr := range edge {
		if err := runCaseA(bin, arr); err != nil {
			fmt.Fprintf(os.Stderr, "edge case %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
	}

	for i := 0; i < 100; i++ {
		arr := generateCaseA(rng)
		if err := runCaseA(bin, arr); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
