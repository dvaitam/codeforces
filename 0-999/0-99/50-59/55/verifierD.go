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

var (
	lcmVals   []int
	lcmIndex  map[int]int
	nextLcm   [][]int
	totalMods = 2520
	divCnt    int
	dp        [][][]int64
	used      [][][]bool
	digits    []int
)

func gcd(a, b int) int {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

func initPrecomp() {
	for i := 1; i <= totalMods; i++ {
		if totalMods%i == 0 {
			lcmVals = append(lcmVals, i)
		}
	}
	divCnt = len(lcmVals)
	lcmIndex = make(map[int]int, divCnt)
	for i, v := range lcmVals {
		lcmIndex[v] = i
	}
	nextLcm = make([][]int, divCnt)
	for i := 0; i < divCnt; i++ {
		nextLcm[i] = make([]int, 10)
		for d := 0; d < 10; d++ {
			if d == 0 {
				nextLcm[i][d] = i
			} else {
				cur := lcmVals[i]
				g := gcd(cur, d)
				nl := cur / g * d
				nextLcm[i][d] = lcmIndex[nl]
			}
		}
	}
	dp = make([][][]int64, 20)
	used = make([][][]bool, 20)
	for i := 0; i < 20; i++ {
		dp[i] = make([][]int64, divCnt)
		used[i] = make([][]bool, divCnt)
		for j := 0; j < divCnt; j++ {
			dp[i][j] = make([]int64, totalMods)
			used[i][j] = make([]bool, totalMods)
		}
	}
}

func dfs(pos, lidx, rem int, tight bool) int64 {
	if pos == len(digits) {
		if rem%lcmVals[lidx] == 0 {
			return 1
		}
		return 0
	}
	if !tight && used[pos][lidx][rem] {
		return dp[pos][lidx][rem]
	}
	var res int64
	limit := 9
	if tight {
		limit = digits[pos]
	}
	for d := 0; d <= limit; d++ {
		nt := tight && d == limit
		nlidx := nextLcm[lidx][d]
		nrem := (rem*10 + d) % totalMods
		res += dfs(pos+1, nlidx, nrem, nt)
	}
	if !tight {
		used[pos][lidx][rem] = true
		dp[pos][lidx][rem] = res
	}
	return res
}

func solveNumber(n uint64) int64 {
	s := strconv.FormatUint(n, 10)
	digits = make([]int, len(s))
	for i, ch := range s {
		digits[i] = int(ch - '0')
	}
	for i := 0; i <= len(digits); i++ {
		for j := 0; j < divCnt; j++ {
			for k := 0; k < totalMods; k++ {
				used[i][j][k] = false
			}
		}
	}
	total := dfs(0, lcmIndex[1], 0, true)
	return total - 1
}

func expected(ranges [][2]uint64) string {
	var sb strings.Builder
	for idx, r := range ranges {
		l := r[0]
		rr := r[1]
		left := uint64(0)
		if l > 0 {
			left = l - 1
		}
		ans := solveNumber(rr) - solveNumber(left)
		if idx+1 == len(ranges) {
			fmt.Fprintf(&sb, "%d", ans)
		} else {
			fmt.Fprintf(&sb, "%d\n", ans)
		}
	}
	return sb.String()
}

func generateCase() (string, string) {
	t := rand.Intn(3) + 1
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", t)
	ranges := make([][2]uint64, t)
	for i := 0; i < t; i++ {
		l := rand.Uint64()%1000000000000 + 1
		r := l + uint64(rand.Intn(1000000))
		ranges[i] = [2]uint64{l, r}
		fmt.Fprintf(&sb, "%d %d\n", l, r)
	}
	return sb.String(), expected(ranges)
}

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

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rand.Seed(time.Now().UnixNano())
	initPrecomp()
	for i := 0; i < 100; i++ {
		input, exp := generateCase()
		got, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
		if got != exp {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, exp, got, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
