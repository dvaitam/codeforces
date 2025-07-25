package main

import (
	"bytes"
	"fmt"
	"math/bits"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

type testCaseC struct {
	n, m int
	rows []string
}

func fwt(a []int64, invert bool) {
	n := len(a)
	for step := 1; step < n; step <<= 1 {
		for i := 0; i < n; i += step << 1 {
			for j := 0; j < step; j++ {
				x := a[i+j]
				y := a[i+j+step]
				a[i+j] = x + y
				a[i+j+step] = x - y
			}
		}
	}
	if invert {
		for i := 0; i < n; i++ {
			a[i] /= int64(n)
		}
	}
}

func solveC(tc testCaseC) int64 {
	size := 1 << tc.n
	cnt := make([]int64, size)
	for col := 0; col < tc.m; col++ {
		mask := 0
		for i := 0; i < tc.n; i++ {
			if tc.rows[i][col] == '1' {
				mask |= 1 << i
			}
		}
		cnt[mask]++
	}
	weight := make([]int64, size)
	for mask := 0; mask < size; mask++ {
		k := bits.OnesCount(uint(mask))
		if k > tc.n-k {
			k = tc.n - k
		}
		weight[mask] = int64(k)
	}
	fwt(cnt, false)
	fwt(weight, false)
	for i := 0; i < size; i++ {
		cnt[i] *= weight[i]
	}
	fwt(cnt, true)
	minVal := cnt[0]
	for _, v := range cnt {
		if v < minVal {
			minVal = v
		}
	}
	return minVal
}

func runCaseC(bin string, tc testCaseC) error {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", tc.n, tc.m))
	for _, row := range tc.rows {
		sb.WriteString(row)
		sb.WriteByte('\n')
	}
	input := sb.String()
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	var got int64
	if _, err := fmt.Sscan(strings.TrimSpace(out.String()), &got); err != nil {
		return fmt.Errorf("bad output: %v", err)
	}
	exp := solveC(tc)
	if got != exp {
		return fmt.Errorf("expected %d got %d", exp, got)
	}
	return nil
}

func genCaseC(rng *rand.Rand) testCaseC {
	n := rng.Intn(5) + 1
	m := rng.Intn(8) + 1
	rows := make([]string, n)
	for i := 0; i < n; i++ {
		var sb strings.Builder
		for j := 0; j < m; j++ {
			if rng.Intn(2) == 0 {
				sb.WriteByte('0')
			} else {
				sb.WriteByte('1')
			}
		}
		rows[i] = sb.String()
	}
	return testCaseC{n, m, rows}
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		tc := genCaseC(rng)
		if err := runCaseC(bin, tc); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
