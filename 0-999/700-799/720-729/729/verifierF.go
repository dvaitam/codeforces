package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

type state struct {
	l, r, k int
	turn    int
}

func run(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func getSum(l, r int, pref []int64) int64 {
	if l > r {
		return 0
	}
	return pref[r+1] - pref[l]
}

func solve(l, r, k, turn int, pref []int64, memo map[state]int64) int64 {
	if l > r {
		return 0
	}
	st := state{l, r, k, turn}
	if v, ok := memo[st]; ok {
		return v
	}
	var res int64
	rem := r - l + 1
	moves := []int{}
	if k == 0 {
		moves = []int{1, 2}
	} else {
		moves = []int{k, k + 1}
	}
	if turn == 0 {
		res = -1 << 60
		for _, x := range moves {
			if x <= rem {
				sum := getSum(l, l+x-1, pref)
				val := sum + solve(l+x, r, x, 1, pref, memo)
				if val > res {
					res = val
				}
			}
		}
	} else {
		res = 1 << 60
		for _, x := range moves {
			if x <= rem {
				sum := getSum(r-x+1, r, pref)
				val := solve(l, r-x, x, 0, pref, memo) - sum
				if val < res {
					res = val
				}
			}
		}
	}
	if res > 1<<59 || res < -(1<<59) {
		res = 0
	}
	memo[st] = res
	return res
}

func expected(a []int64) string {
	n := len(a)
	pref := make([]int64, n+1)
	for i := 0; i < n; i++ {
		pref[i+1] = pref[i] + a[i]
	}
	memo := make(map[state]int64)
	ans := solve(0, n-1, 0, 0, pref, memo)
	return fmt.Sprintf("%d", ans)
}

func main() {
	if len(os.Args) != 2 && !(len(os.Args) == 3 && os.Args[1] == "--") {
		fmt.Println("usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[len(os.Args)-1]
	r := rand.New(rand.NewSource(6))
	for tc := 1; tc <= 100; tc++ {
		n := r.Intn(8) + 1
		arr := make([]int64, n)
		for i := 0; i < n; i++ {
			arr[i] = int64(r.Intn(21) - 10)
		}
		var sb strings.Builder
		fmt.Fprintf(&sb, "%d\n", n)
		for i := 0; i < n; i++ {
			if i > 0 {
				sb.WriteByte(' ')
			}
			fmt.Fprintf(&sb, "%d", arr[i])
		}
		sb.WriteByte('\n')
		input := sb.String()
		expect := expected(arr)
		out, err := run(bin, input)
		if err != nil {
			fmt.Printf("test %d: %v\n", tc, err)
			os.Exit(1)
		}
		if out != expect {
			fmt.Printf("test %d failed\ninput:\n%sexpected: %s\ngot: %s\n", tc, input, expect, out)
			os.Exit(1)
		}
	}
	fmt.Println("All 100 tests passed")
}
