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

type player struct {
	num  int
	role int // 0 goalie, 1 defender, 2 forward
}

func comb(n, k int) int64 {
	if n < k || k < 0 {
		return 0
	}
	if k > n-k {
		k = n - k
	}
	res := int64(1)
	for i := 1; i <= k; i++ {
		res = res * int64(n-k+i) / int64(i)
	}
	return res
}

func solveC(g, d, f int, gNums, dNums, fNums []int) int64 {
	players := make([]player, g+d+f)
	idx := 0
	for i := 0; i < g; i++ {
		players[idx] = player{gNums[i], 0}
		idx++
	}
	for i := 0; i < d; i++ {
		players[idx] = player{dNums[i], 1}
		idx++
	}
	for i := 0; i < f; i++ {
		players[idx] = player{fNums[i], 2}
		idx++
	}
	sort.Slice(players, func(i, j int) bool { return players[i].num < players[j].num })
	n := len(players)
	prefixG := make([]int, n+1)
	prefixD := make([]int, n+1)
	prefixF := make([]int, n+1)
	for i := 0; i < n; i++ {
		prefixG[i+1] = prefixG[i]
		prefixD[i+1] = prefixD[i]
		prefixF[i+1] = prefixF[i]
		switch players[i].role {
		case 0:
			prefixG[i+1]++
		case 1:
			prefixD[i+1]++
		case 2:
			prefixF[i+1]++
		}
	}
	var ans int64
	j := 0
	for i := 0; i < n; i++ {
		if j < i+1 {
			j = i + 1
		}
		for j < n && players[j].num <= 2*players[i].num {
			j++
		}
		numG := prefixG[j] - prefixG[i]
		numD := prefixD[j] - prefixD[i]
		numF := prefixF[j] - prefixF[i]
		switch players[i].role {
		case 0:
			if numD >= 2 && numF >= 3 {
				ans += comb(numD, 2) * comb(numF, 3)
			}
		case 1:
			if numG >= 1 && numD >= 2 && numF >= 3 {
				ans += int64(numG) * comb(numD-1, 1) * comb(numF, 3)
			}
		case 2:
			if numG >= 1 && numD >= 2 && numF >= 3 {
				ans += int64(numG) * comb(numD, 2) * comb(numF-1, 2)
			}
		}
	}
	return ans
}

func generateCase(rng *rand.Rand) (string, string) {
	g := rng.Intn(5) + 1
	d := rng.Intn(5) + 1
	f := rng.Intn(5) + 1
	gNums := make([]int, g)
	dNums := make([]int, d)
	fNums := make([]int, f)
	for i := range gNums {
		gNums[i] = rng.Intn(100) + 1
	}
	for i := range dNums {
		dNums[i] = rng.Intn(100) + 1
	}
	for i := range fNums {
		fNums[i] = rng.Intn(100) + 1
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d %d\n", g, d, f))
	for i, v := range gNums {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(v))
	}
	sb.WriteByte('\n')
	for i, v := range dNums {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(v))
	}
	sb.WriteByte('\n')
	for i, v := range fNums {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(v))
	}
	sb.WriteByte('\n')
	ans := solveC(g, d, f, gNums, dNums, fNums)
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
	if strings.TrimSpace(out.String()) != strings.TrimSpace(expected) {
		return fmt.Errorf("expected %q got %q", strings.TrimSpace(expected), strings.TrimSpace(out.String()))
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierC.go /path/to/binary")
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
