package main

import (
	"bytes"
	"fmt"
	"math/bits"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strings"
	"time"
)

func nextPow2(x int) int {
	if x <= 1 {
		return 1
	}
	return 1 << (bits.Len(uint(x - 1)))
}

func expected(n int, arr []int) int {
	freq := make([]int, n+2)
	for _, x := range arr {
		if x >= 1 && x <= n {
			freq[x]++
		}
	}
	pre := make([]int, n+1)
	for i := 1; i <= n; i++ {
		pre[i] = pre[i-1] + freq[i]
	}
	powers := []int{}
	for p := 1; p <= 2*n; p <<= 1 {
		powers = append(powers, p)
	}
	ans := int(1e9)
	for i := 0; i <= n; i++ {
		c1 := pre[i]
		cost1 := nextPow2(c1) - c1
		for _, limit := range powers {
			target := pre[i] + limit
			j := sort.Search(len(pre), func(k int) bool { return pre[k] > target }) - 1
			if j < i+1 {
				continue
			}
			c2 := pre[j] - pre[i]
			c3 := n - pre[j]
			cost2 := nextPow2(c2) - c2
			cost3 := nextPow2(c3) - c3
			total := cost1 + cost2 + cost3
			if total < ans {
				ans = total
			}
		}
	}
	return ans
}

func runCase(bin string, n int, arr []int) error {
	var sb strings.Builder
	sb.WriteString("1\n")
	sb.WriteString(fmt.Sprint(n))
	sb.WriteByte('\n')
	for i, v := range arr {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprint(v))
	}
	sb.WriteByte('\n')
	input := sb.String()
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	got := strings.TrimSpace(out.String())
	exp := fmt.Sprint(expected(n, arr))
	if got != exp {
		return fmt.Errorf("expected %s got %s", exp, got)
	}
	return nil
}

func randCase(rng *rand.Rand) (int, []int) {
	n := rng.Intn(10) + 1
	arr := make([]int, n)
	for i := 0; i < n; i++ {
		arr[i] = rng.Intn(n*2 + 1)
	}
	return n, arr
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	var cases []struct {
		n   int
		arr []int
	}
	cases = append(cases, struct {
		n   int
		arr []int
	}{3, []int{1, 2, 3}})
	cases = append(cases, struct {
		n   int
		arr []int
	}{5, []int{1, 1, 1, 1, 1}})
	for i := 0; i < 100; i++ {
		n, arr := randCase(rng)
		cases = append(cases, struct {
			n   int
			arr []int
		}{n, arr})
	}
	for idx, tc := range cases {
		if err := runCase(bin, tc.n, tc.arr); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
