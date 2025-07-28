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

func expected(n int, k, h []int64) int64 {
	curL := k[n-1] - h[n-1] + 1
	curR := k[n-1]
	var ans int64
	for i := n - 2; i >= 0; i-- {
		start := k[i] - h[i] + 1
		end := k[i]
		if end >= curL {
			if start < curL {
				curL = start
			}
		} else {
			length := curR - curL + 1
			ans += length * (length + 1) / 2
			curL = start
			curR = end
		}
	}
	length := curR - curL + 1
	ans += length * (length + 1) / 2
	return ans
}

func runCase(bin string, n int, k, h []int64) error {
	var sb strings.Builder
	sb.WriteString("1\n")
	sb.WriteString(strconv.Itoa(n))
	sb.WriteByte('\n')
	for i := 0; i < n; i++ {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprint(k[i]))
	}
	sb.WriteByte('\n')
	for i := 0; i < n; i++ {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprint(h[i]))
	}
	sb.WriteByte('\n')
	input := sb.String()

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
		return fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	gotStr := strings.TrimSpace(out.String())
	exp := fmt.Sprint(expected(n, k, h))
	if gotStr != exp {
		return fmt.Errorf("expected %s got %s", exp, gotStr)
	}
	return nil
}

func randCase(rng *rand.Rand) (int, []int64, []int64) {
	n := rng.Intn(10) + 1
	k := make([]int64, n)
	h := make([]int64, n)
	cur := int64(0)
	for i := 0; i < n; i++ {
		cur += int64(rng.Intn(5) + 1)
		k[i] = cur
		h[i] = int64(rng.Intn(5) + 1)
	}
	return n, k, h
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	var cases []struct {
		n    int
		k, h []int64
	}
	cases = append(cases, struct {
		n    int
		k, h []int64
	}{1, []int64{1}, []int64{1}})
	cases = append(cases, struct {
		n    int
		k, h []int64
	}{2, []int64{2, 4}, []int64{1, 2}})
	for i := 0; i < 100; i++ {
		n, k, h := randCase(rng)
		cases = append(cases, struct {
			n    int
			k, h []int64
		}{n, k, h})
	}
	for idx, tc := range cases {
		if err := runCase(bin, tc.n, tc.k, tc.h); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
