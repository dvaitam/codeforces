package main

import (
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strings"
)

func runBinary(bin string, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	out, err := cmd.CombinedOutput()
	return string(out), err
}

func solve(n int, m int64, a []int64) int {
	sort.Slice(a, func(i, j int) bool { return a[i] < a[j] })
	check := func(x int) bool {
		var tmp int64
		cnt := 0
		cur := int64(0)
		for i := n - 1; i >= 0; i-- {
			if a[i] > cur {
				tmp += a[i] - cur
			}
			cnt++
			if cnt == x {
				cur++
				cnt = 0
			}
		}
		return tmp >= m
	}
	l, r := 1, n
	ans := -1
	for l <= r {
		mid := (l + r) / 2
		if check(mid) {
			ans = mid
			r = mid - 1
		} else {
			l = mid + 1
		}
	}
	return ans
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierD2.go /path/to/binary")
		return
	}
	bin := os.Args[1]
	rand.Seed(5)
	for t := 1; t <= 100; t++ {
		n := rand.Intn(8) + 1
		m := rand.Int63n(100) + 1
		a := make([]int64, n)
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d %d\n", n, m))
		for i := 0; i < n; i++ {
			a[i] = rand.Int63n(50) + 1
			sb.WriteString(fmt.Sprintf("%d ", a[i]))
		}
		sb.WriteString("\n")
		expected := solve(n, m, a)
		out, err := runBinary(bin, sb.String())
		if err != nil {
			fmt.Printf("test %d runtime error: %v\n", t, err)
			fmt.Println(out)
			return
		}
		var got int
		fmt.Fscan(strings.NewReader(out), &got)
		if got != expected {
			fmt.Printf("test %d failed\ninput:\n%sexpected %d got %d\n", t, sb.String(), expected, got)
			return
		}
	}
	fmt.Println("all tests passed")
}
