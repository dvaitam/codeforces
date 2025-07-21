package main

import (
	"bufio"
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

func bitsTrailing(x uint) int {
	return bitsTrailingDeBruijn(x)
}

var deBruijn = uint(0x077CB531)
var idx32 = [32]int{
	0, 1, 28, 2, 29, 14, 24, 3,
	30, 22, 20, 15, 25, 17, 4, 8,
	31, 27, 13, 23, 21, 19, 16, 7,
	26, 12, 18, 6, 11, 5, 10, 9,
}

func bitsTrailingDeBruijn(v uint) int {
	return idx32[((v&-v)*deBruijn)>>27]
}

func solve(data string) string {
	in := bufio.NewReader(strings.NewReader(data))
	var N, K int
	var Tlimit int64
	fmt.Fscan(in, &N, &K, &Tlimit)
	var s string
	fmt.Fscan(in, &s)
	tAll := make([]int64, K)
	for i := 0; i < K; i++ {
		fmt.Fscan(in, &tAll[i])
	}
	aAll := make([][]int64, K)
	for i := 0; i < K; i++ {
		aAll[i] = make([]int64, K)
		for j := 0; j < K; j++ {
			fmt.Fscan(in, &aAll[i][j])
		}
	}
	present := make([]bool, K)
	for i := 0; i < N; i++ {
		present[s[i]-'A'] = true
	}
	typeMap := make([]int, K)
	revMap := make([]int, 0, K)
	for i := 0; i < K; i++ {
		if present[i] {
			typeMap[i] = len(revMap)
			revMap = append(revMap, i)
		}
	}
	m := len(revMap)
	if m == 0 {
		return "0\n"
	}
	A := make([]int, N)
	for i := 0; i < N; i++ {
		A[i] = typeMap[s[i]-'A']
	}
	nextpos := make([][]int, m)
	for x := 0; x < m; x++ {
		nextpos[x] = make([]int, N+1)
		next := N
		for i := N - 1; i >= 0; i-- {
			if A[i] == x {
				next = i
			}
			nextpos[x][i] = next
		}
		nextpos[x][N] = N
	}
	size := 1 << m
	totalOut := make([]int64, size)
	updates := make([][]struct {
		mask uint32
		val  int64
	}, m)
	type pair struct{ pos, idx int }
	ord := make([]pair, 0, m)
	for p := 0; p < N; p++ {
		i := A[p]
		ord = ord[:0]
		for x := 0; x < m; x++ {
			np := nextpos[x][p]
			if np < N {
				ord = append(ord, pair{np, x})
			}
		}
		sort.Slice(ord, func(i, j int) bool { return ord[i].pos < ord[j].pos })
		var mask uint32
		for _, pr := range ord {
			j := pr.idx
			cost := aAll[revMap[i]][revMap[j]]
			totalOut[mask] += cost
			updates[j] = append(updates[j], struct {
				mask uint32
				val  int64
			}{mask, cost})
			mask |= 1 << j
		}
	}
	for b := 0; b < m; b++ {
		for mask := 0; mask < size; mask++ {
			if mask&(1<<b) != 0 {
				totalOut[mask] += totalOut[mask^(1<<b)]
			}
		}
	}
	h := make([]int64, size)
	for j := 0; j < m; j++ {
		W := make([]int64, size)
		for _, u := range updates[j] {
			W[u.mask] += u.val
		}
		updates[j] = nil
		for b := 0; b < m; b++ {
			for mask := 0; mask < size; mask++ {
				if mask&(1<<b) != 0 {
					W[mask] += W[mask^(1<<b)]
				}
			}
		}
		bit := 1 << j
		for mask := bit; mask < size; mask++ {
			if mask&bit != 0 {
				h[mask] += W[mask]
			}
		}
	}
	tMap := make([]int64, m)
	for idx, orig := range revMap {
		tMap[idx] = tAll[orig]
	}
	tMask := make([]int64, size)
	for mask := 1; mask < size; mask++ {
		lsb := mask & -mask
		b := bitsTrailing(uint(lsb))
		tMask[mask] = tMask[mask^lsb] + tMap[b]
	}
	full := size - 1
	var ans int64
	for mask := 0; mask < size; mask++ {
		if mask == full {
			continue
		}
		riskAdj := totalOut[mask] - h[mask]
		if riskAdj+tMask[mask] <= Tlimit {
			ans++
		}
	}
	return fmt.Sprintf("%d\n", ans)
}

func generateCase(rng *rand.Rand) (string, string) {
	N := rng.Intn(8) + 1
	K := rng.Intn(3) + 1
	Tlimit := int64(rng.Intn(50) + 10)
	letters := make([]byte, N)
	for i := 0; i < N; i++ {
		letters[i] = byte('A' + rng.Intn(K))
	}
	tVals := make([]int64, K)
	for i := range tVals {
		tVals[i] = int64(rng.Intn(10) + 1)
	}
	a := make([][]int64, K)
	for i := 0; i < K; i++ {
		a[i] = make([]int64, K)
		for j := 0; j < K; j++ {
			a[i][j] = int64(rng.Intn(5))
		}
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d %d\n", N, K, Tlimit))
	sb.WriteString(string(letters) + "\n")
	for i, v := range tVals {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.FormatInt(v, 10))
	}
	sb.WriteByte('\n')
	for i := 0; i < K; i++ {
		for j := 0; j < K; j++ {
			if j > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.FormatInt(a[i][j], 10))
		}
		sb.WriteByte('\n')
	}
	input := sb.String()
	expected := solve(input)
	return input, expected
}

func runCase(exe string, in, expected string) error {
	cmd := exec.Command(exe)
	cmd.Stdin = strings.NewReader(in)
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
