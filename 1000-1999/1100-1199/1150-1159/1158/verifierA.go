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

func solveA(n, m int, b, g []int64) (int64, bool) {
	sort.Slice(b, func(i, j int) bool { return b[i] < b[j] })
	sort.Slice(g, func(i, j int) bool { return g[i] < g[j] })
	bMax := b[n-1]
	if g[0] < bMax {
		return -1, false
	}
	var sumB int64
	for i := 0; i < n; i++ {
		sumB += b[i]
	}
	ans := sumB * int64(m)
	for j := 0; j < m; j++ {
		ans += g[j] - bMax
	}
	if g[0] > bMax {
		bSec := b[n-2]
		ans += bMax - bSec
	}
	return ans, true
}

func run(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for t := 0; t < 100; t++ {
		n := rng.Intn(5) + 2
		m := rng.Intn(5) + 2
		b := make([]int64, n)
		g := make([]int64, m)
		for i := 0; i < n; i++ {
			b[i] = rand.Int63n(20)
		}
		for j := 0; j < m; j++ {
			g[j] = rand.Int63n(20)
		}
		// ensure g contains at least one >= max(b)
		maxB := b[0]
		for _, v := range b {
			if v > maxB {
				maxB = v
			}
		}
		g[0] = maxB + int64(rng.Intn(3)) // ensure feasible sometimes
		input := fmt.Sprintf("%d %d\n", n, m)
		for i, v := range b {
			if i > 0 {
				input += " "
			}
			input += fmt.Sprint(v)
		}
		input += "\n"
		for j, v := range g {
			if j > 0 {
				input += " "
			}
			input += fmt.Sprint(v)
		}
		input += "\n"
		expVal, ok := solveA(n, m, append([]int64(nil), b...), append([]int64(nil), g...))
		var expected string
		if ok {
			expected = fmt.Sprint(expVal)
		} else {
			expected = "-1"
		}
		got, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:%s", t+1, err, input)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != strings.TrimSpace(expected) {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:%s", t+1, expected, got, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
