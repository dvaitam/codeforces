package main

import (
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"
)

type seg struct{ l, r int }

func expected(n, k int, temps []int) int {
	cold := 0
	var segs []seg
	inSeg := false
	var l int
	for i, t := range temps {
		if t < 0 {
			cold++
			if !inSeg {
				inSeg = true
				l = i
			}
		} else {
			if inSeg {
				segs = append(segs, seg{l, i - 1})
				inSeg = false
			}
		}
	}
	if inSeg {
		segs = append(segs, seg{l, n - 1})
	}
	if cold > k {
		return -1
	}
	if len(segs) == 0 {
		return 0
	}
	m := len(segs)
	gaps := []int{}
	for i := 0; i < m-1; i++ {
		gap := segs[i+1].l - segs[i].r - 1
		if gap > 0 {
			gaps = append(gaps, gap)
		}
	}
	tail := 0
	if segs[m-1].r < n-1 {
		tail = n - 1 - segs[m-1].r
	}
	defaultCost := 2 * m
	if tail == 0 {
		defaultCost--
	}
	budget := k - cold
	sort.Ints(gaps)
	save := 0
	for _, gap := range gaps {
		if gap <= budget {
			budget -= gap
			save += 2
		} else {
			break
		}
	}
	if tail > 0 && tail <= budget {
		save++
	}
	return defaultCost - save
}

func runCase(bin string, n, k int, temps []int) error {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", n, k))
	for i, t := range temps {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(t))
	}
	sb.WriteByte('\n')
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(sb.String())
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("runtime error: %v", err)
	}
	got, err := strconv.Atoi(strings.TrimSpace(string(out)))
	if err != nil {
		return fmt.Errorf("bad output: %v", err)
	}
	exp := expected(n, k, temps)
	if got != exp {
		return fmt.Errorf("expected %d got %d", exp, got)
	}
	return nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run verifierD.go /path/to/binary")
		return
	}
	bin := os.Args[1]
	rand.Seed(1)
	for tcase := 1; tcase <= 120; tcase++ {
		n := rand.Intn(20) + 1
		k := rand.Intn(n + 1)
		temps := make([]int, n)
		for i := range temps {
			temps[i] = rand.Intn(41) - 20
		}
		if err := runCase(bin, n, k, temps); err != nil {
			fmt.Printf("Test %d failed: %v\n", tcase, err)
			return
		}
	}
	fmt.Println("OK")
}
