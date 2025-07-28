package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"
	"time"
)

const MaxA = 100000

var preDivs = make([][]int, MaxA+1)

func init() {
	for d := 1; d <= MaxA; d++ {
		for m := d; m <= MaxA; m += d {
			preDivs[m] = append(preDivs[m], d)
		}
	}
}

func solveCase(n, m int, a []int) int {
	sort.Ints(a)
	allCnt := make([]int, m+1)
	allCov := 0
	for i := 0; i < n; i++ {
		for _, d := range preDivs[a[i]] {
			if d > m {
				break
			}
			if allCnt[d] == 0 {
				allCov++
			}
			allCnt[d]++
		}
	}
	if allCov < m {
		return -1
	}
	cnt := make([]int, m+1)
	covered := 0
	ans := int(1<<31 - 1)
	l := 0
	for r := 0; r < n; r++ {
		for _, d := range preDivs[a[r]] {
			if d > m {
				break
			}
			cnt[d]++
			if cnt[d] == 1 {
				covered++
			}
		}
		for covered == m {
			if diff := a[r] - a[l]; diff < ans {
				ans = diff
			}
			for _, d := range preDivs[a[l]] {
				if d > m {
					break
				}
				cnt[d]--
				if cnt[d] == 0 {
					covered--
				}
			}
			l++
		}
	}
	if ans == int(1<<31-1) {
		return -1
	}
	return ans
}

func generateCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(10) + 1
	m := rng.Intn(10) + 1
	a := make([]int, n)
	for i := range a {
		a[i] = rng.Intn(100) + 1
	}
	ans := solveCase(n, m, a)
	var sb strings.Builder
	sb.WriteString("1\n")
	sb.WriteString(fmt.Sprintf("%d %d\n", n, m))
	for i, v := range a {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(v))
	}
	sb.WriteByte('\n')
	return sb.String(), fmt.Sprintf("%d\n", ans)
}

func runCase(exe, input, expected string) error {
	cmd := exec.Command(exe)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	outStr := strings.TrimSpace(out.String())
	exp := strings.TrimSpace(expected)
	if outStr != exp {
		return fmt.Errorf("expected %s got %s", exp, outStr)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	exe := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := generateCase(rng)
		if err := runCase(exe, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
