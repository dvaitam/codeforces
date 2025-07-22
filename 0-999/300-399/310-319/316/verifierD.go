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

const mod = 1000000007

func nextPerm(a []int) bool {
	n := len(a)
	i := n - 2
	for i >= 0 && a[i] >= a[i+1] {
		i--
	}
	if i < 0 {
		return false
	}
	j := n - 1
	for a[j] <= a[i] {
		j--
	}
	a[i], a[j] = a[j], a[i]
	for l, r := i+1, n-1; l < r; l, r = l+1, r-1 {
		a[l], a[r] = a[r], a[l]
	}
	return true
}

func validPerm(p, cap []int) bool {
	n := len(p)
	vis := make([]bool, n)
	for i := 0; i < n; i++ {
		if vis[i] {
			continue
		}
		cur := i
		cycle := []int{}
		for !vis[cur] {
			vis[cur] = true
			cycle = append(cycle, cur)
			cur = p[cur]
		}
		k := len(cycle)
		if k <= 1 {
			continue
		}
		b0, b1, b2 := 0, 0, 0
		for _, v := range cycle {
			if cap[v] <= 0 {
				b0++
			} else if cap[v] == 1 {
				b1++
			} else {
				b2++
			}
		}
		if b0 > 0 {
			return false
		}
		if b1 > 2 {
			return false
		}
		if k >= 3 && b2 < k-2 {
			return false
		}
	}
	return true
}

func solveCase(caps []int) int64 {
	n := len(caps)
	p := make([]int, n)
	for i := 0; i < n; i++ {
		p[i] = i
	}
	var cnt int64 = 0
	for {
		if validPerm(p, caps) {
			cnt++
		}
		if !nextPerm(p) {
			break
		}
	}
	return cnt % mod
}

func genCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(5) + 1
	caps := make([]int, n)
	for i := 0; i < n; i++ {
		caps[i] = rng.Intn(3)
	}
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", n)
	for i := 0; i < n; i++ {
		if i > 0 {
			sb.WriteByte(' ')
		}
		fmt.Fprintf(&sb, "%d", caps[i])
	}
	sb.WriteByte('\n')
	ans := solveCase(caps)
	return sb.String(), fmt.Sprintf("%d\n", ans)
}

func runCase(bin, input, expected string) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	if strings.TrimSpace(out.String()) != strings.TrimSpace(expected) {
		return fmt.Errorf("expected %s got %s", expected, out.String())
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := genCase(rng)
		if err := runCase(bin, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
