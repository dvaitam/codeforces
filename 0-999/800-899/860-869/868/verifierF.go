package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

type Test struct {
	in  string
	out string
}

func run(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func oracle(input string) string {
	in := bufio.NewReader(strings.NewReader(input))
	var n, k int
	fmt.Fscan(in, &n, &k)
	arr := make([]int, n+1)
	for i := 1; i <= n; i++ {
		fmt.Fscan(in, &arr[i])
	}
	const INF int64 = 1 << 60
	dpPrev := make([]int64, n+1)
	for i := 1; i <= n; i++ {
		dpPrev[i] = INF
	}
	cnt := make([]int, n+1)
	curL, curR := 1, 0
	var curCost int64
	add := func(idx int) {
		v := arr[idx]
		curCost += int64(cnt[v])
		cnt[v]++
	}
	remove := func(idx int) {
		v := arr[idx]
		cnt[v]--
		curCost -= int64(cnt[v])
	}
	setSeg := func(l, r int) {
		for curL > l {
			curL--
			add(curL)
		}
		for curR < r {
			curR++
			add(curR)
		}
		for curL < l {
			remove(curL)
			curL++
		}
		for curR > r {
			remove(curR)
			curR--
		}
	}
	dpCur := make([]int64, n+1)
	var solve func(L, R, optL, optR int)
	solve = func(L, R, optL, optR int) {
		if L > R {
			return
		}
		mid := (L + R) / 2
		bestPos := optL
		bestVal := INF
		start := optL
		end := optR
		if start < 0 {
			start = 0
		}
		if end > mid-1 {
			end = mid - 1
		}
		for i := start; i <= end; i++ {
			setSeg(i+1, mid)
			val := dpPrev[i] + curCost
			if val < bestVal {
				bestVal = val
				bestPos = i
			}
		}
		dpCur[mid] = bestVal
		solve(L, mid-1, optL, bestPos)
		solve(mid+1, R, bestPos, optR)
	}
	dpPrev[0] = 0
	for seg := 1; seg <= k; seg++ {
		for i := 0; i <= n; i++ {
			dpCur[i] = INF
		}
		curL, curR, curCost = 1, 0, 0
		for i := range cnt {
			cnt[i] = 0
		}
		solve(1, n, 0, n-1)
		dpPrev, dpCur = dpPrev, dpCur
		dpPrev, dpCur = dpCur, dpPrev
	}
	return fmt.Sprintf("%d", dpPrev[n])
}

func genCase(rng *rand.Rand) Test {
	n := rng.Intn(15) + 1
	k := rng.Intn(min(5, n)) + 1
	arr := make([]int, n)
	for i := range arr {
		arr[i] = rng.Intn(n) + 1
	}
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d\n", n, k)
	for i, v := range arr {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", v))
	}
	sb.WriteByte('\n')
	input := sb.String()
	out := oracle(input)
	return Test{input, out}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(6))
	for i := 0; i < 100; i++ {
		tc := genCase(rng)
		got, err := run(bin, tc.in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, tc.in)
			os.Exit(1)
		}
		if got != tc.out {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, tc.out, got, tc.in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
