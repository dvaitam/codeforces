package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

const MOD int64 = 1_000_000_007
const INV2 int64 = (MOD + 1) / 2

func run(bin, input string) (string, error) {
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
		return "", fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func solve(c []int) (int64, int64) {
	m := len(c)
	maxC := 0
	for _, v := range c {
		if v > maxC {
			maxC = v
		}
	}
	arrOdd := make([]int64, maxC)
	arrEven := make([]int64, maxC)
	for _, v := range c {
		if v%2 == 1 {
			arrOdd[v-1]++
		} else {
			arrEven[v-1]++
		}
	}
	prefOdd := make([]int64, maxC+1)
	prefEven := make([]int64, maxC+1)
	for t := maxC - 1; t >= 0; t-- {
		prefOdd[t] = prefOdd[t+1] + arrOdd[t]
		prefEven[t] = prefEven[t+1] + arrEven[t]
	}
	fact := make([]int64, m+1)
	fact[0] = 1
	for i := 1; i <= m; i++ {
		fact[i] = fact[i-1] * int64(i) % MOD
	}
	pos := int64(1)
	ans := int64(0)
	ways := int64(1)
	for x := -maxC + 1; x <= maxC-1; x++ {
		t := x
		if t < 0 {
			t = -t
		}
		if t >= maxC {
			continue
		}
		var cnt int64
		if t%2 == 0 {
			cnt = prefOdd[t]
		} else {
			cnt = prefEven[t]
		}
		if cnt == 0 {
			continue
		}
		cntMod := cnt % MOD
		posMod := pos % MOD
		sumPosMod := (cntMod * posMod) % MOD
		sumPosMod = (sumPosMod + cntMod*((cnt-1)%MOD)%MOD*INV2%MOD) % MOD
		xMod := int64(x % int(MOD))
		if xMod < 0 {
			xMod += MOD
		}
		ans = (ans + xMod*sumPosMod%MOD) % MOD
		ways = ways * fact[int(cnt)] % MOD
		pos += cnt
	}
	return ans, ways
}

func parseOutput(output string) (int64, int64, error) {
	parts := strings.Fields(output)
	if len(parts) != 2 {
		return 0, 0, fmt.Errorf("expected two integers")
	}
	a, err1 := strconv.ParseInt(parts[0], 10, 64)
	b, err2 := strconv.ParseInt(parts[1], 10, 64)
	if err1 != nil || err2 != nil {
		return 0, 0, fmt.Errorf("invalid integers")
	}
	a %= MOD
	b %= MOD
	if a < 0 {
		a += MOD
	}
	if b < 0 {
		b += MOD
	}
	return a, b, nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierG.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rand.Seed(time.Now().UnixNano())
	for t := 0; t < 100; t++ {
		m := rand.Intn(5) + 1
		c := make([]int, m)
		var sb strings.Builder
		fmt.Fprintf(&sb, "%d\n", m)
		for i := 0; i < m; i++ {
			c[i] = rand.Intn(5) + 1
			fmt.Fprintf(&sb, "%d ", c[i])
		}
		sb.WriteByte('\n')
		input := sb.String()
		out, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", t+1, err)
			os.Exit(1)
		}
		gotAns, gotWays, err := parseOutput(out)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d invalid output: %v\n", t+1, err)
			os.Exit(1)
		}
		ans, ways := solve(c)
		if gotAns != ans%MOD || gotWays != ways%MOD {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %d %d got %d %d\n", t+1, ans%MOD, ways%MOD, gotAns, gotWays)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
