package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strings"
	"time"
)

const mod = 998244353

func solve(arr []int) int64 {
	n := len(arr)
	type node struct{ id, val int }
	type node1 struct{ l, r, val int }
	a := make([]node, n)
	for i, v := range arr {
		a[i] = node{i, v}
	}
	sort.Slice(a, func(i, j int) bool {
		if a[i].val != a[j].val {
			return a[i].val < a[j].val
		}
		return a[i].id < a[j].id
	})
	b := make([]node1, 0)
	for i := 1; i < n; i++ {
		if a[i].val == a[i-1].val {
			b = append(b, node1{a[i-1].id, a[i].id, a[i].val})
		}
	}
	cnt := len(b)
	sort.Slice(b, func(i, j int) bool { return b[i].l < b[j].l })
	ans := make([]int, n)
	if cnt > 0 {
		lp, ip, col := b[0].l, 0, 1
		for lp < n && ip < cnt {
			for lp > b[ip].r {
				ip++
				if ip >= cnt {
					break
				}
				if b[ip].val != b[ip-1].val && lp <= b[ip].l {
					col++
				}
				if lp < b[ip].l {
					lp = b[ip].l
				}
			}
			if ip >= cnt {
				break
			}
			ans[lp] = col
			lp++
		}
	}
	var res int64 = 1
	for i := 1; i < n; i++ {
		if ans[i] == 0 || ans[i] != ans[i-1] {
			res = res * 2 % mod
		}
	}
	return res
}

func generateCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(20) + 2
	arr := make([]int, n)
	for i := range arr {
		arr[i] = rng.Intn(20)
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for i, v := range arr {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", v))
	}
	sb.WriteByte('\n')
	ans := solve(arr)
	return sb.String(), fmt.Sprintf("%d", ans)
}

func runCase(bin, input, expected string) error {
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
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	if strings.TrimSpace(out.String()) != strings.TrimSpace(expected) {
		return fmt.Errorf("expected %q got %q", expected, strings.TrimSpace(out.String()))
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := generateCase(rng)
		if err := runCase(bin, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
