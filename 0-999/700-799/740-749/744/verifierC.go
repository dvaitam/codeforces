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

type pair struct{ r, b int }

func bitsTrailing(v int) int {
	for i := 0; i < 32; i++ {
		if v>>i&1 == 1 {
			return i
		}
	}
	return 0
}

func prune(arr []pair) []pair {
	for i := 1; i < len(arr); i++ {
		j := i
		for j > 0 && (arr[j].r > arr[j-1].r || (arr[j].r == arr[j-1].r && arr[j].b > arr[j-1].b)) {
			arr[j], arr[j-1] = arr[j-1], arr[j]
			j--
		}
	}
	res := make([]pair, 0, len(arr))
	maxB := -1
	for _, p := range arr {
		if p.b > maxB {
			res = append(res, p)
			maxB = p.b
		}
	}
	return res
}

func can(t, n int, rreq, breq []int, isRed []bool, rcnt, bcnt []int, goal int) bool {
	N := 1 << n
	dp := make([][]pair, N)
	dp[0] = []pair{{t, t}}
	for mask := 0; mask < N; mask++ {
		states := dp[mask]
		if len(states) == 0 {
			continue
		}
		dp[mask] = prune(states)
		if mask == goal && len(dp[mask]) > 0 {
			return true
		}
		A := rcnt[mask]
		B := bcnt[mask]
		for _, st := range dp[mask] {
			rrem, brem := st.r, st.b
			for i := 0; i < n; i++ {
				bit := 1 << i
				if mask&bit != 0 {
					continue
				}
				needR := rreq[i] - A
				if needR < 0 {
					needR = 0
				}
				needB := breq[i] - B
				if needB < 0 {
					needB = 0
				}
				if rrem >= needR && brem >= needB {
					nm := mask | bit
					dp[nm] = append(dp[nm], pair{rrem - needR, brem - needB})
				}
			}
		}
	}
	return len(dp[goal]) > 0
}

func solveC(n int, isRed []bool, rreq, breq []int) int {
	N := 1 << n
	rcnt := make([]int, N)
	bcnt := make([]int, N)
	for mask := 1; mask < N; mask++ {
		lsb := mask & -mask
		i := bitsTrailing(lsb)
		prev := mask ^ lsb
		if isRed[i] {
			rcnt[mask] = rcnt[prev] + 1
			bcnt[mask] = bcnt[prev]
		} else {
			bcnt[mask] = bcnt[prev] + 1
			rcnt[mask] = rcnt[prev]
		}
	}
	sumR, sumB := 0, 0
	for i := 0; i < n; i++ {
		sumR += rreq[i]
		sumB += breq[i]
	}
	lo, hi := 0, sumR+sumB
	goal := N - 1
	for lo < hi {
		mid := (lo + hi) / 2
		if can(mid, n, rreq, breq, isRed, rcnt, bcnt, goal) {
			hi = mid
		} else {
			lo = mid + 1
		}
	}
	return lo + n
}

func genCase(rng *rand.Rand) (string, int) {
	n := rng.Intn(4) + 1
	isRed := make([]bool, n)
	rreq := make([]int, n)
	breq := make([]int, n)
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", n)
	for i := 0; i < n; i++ {
		if rng.Intn(2) == 0 {
			isRed[i] = true
			sb.WriteByte('R')
		} else {
			sb.WriteByte('B')
		}
		rreq[i] = rng.Intn(5)
		breq[i] = rng.Intn(5)
		fmt.Fprintf(&sb, " %d %d\n", rreq[i], breq[i])
	}
	ans := solveC(n, isRed, rreq, breq)
	return sb.String(), ans
}

func run(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for t := 1; t <= 100; t++ {
		in, expect := genCase(rng)
		outStr, err := run(bin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", t, err, in)
			os.Exit(1)
		}
		var got int
		if _, err := fmt.Sscan(outStr, &got); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: invalid output %q\n", t, outStr)
			os.Exit(1)
		}
		if got != expect {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %d got %d\ninput:\n%s", t, expect, got, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
