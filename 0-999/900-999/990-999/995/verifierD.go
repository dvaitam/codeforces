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

func solveD(n int, arr []int64, queries [][2]int64) []float64 {
	size := 1 << uint(n)
	sum := int64(0)
	for _, v := range arr {
		sum += v
	}
	res := make([]float64, len(queries)+1)
	denom := float64(size)
	res[0] = float64(sum) / denom
	for i, q := range queries {
		x := q[0]
		y := q[1]
		sum += y - arr[x]
		arr[x] = y
		res[i+1] = float64(sum) / denom
	}
	return res
}

func genCaseD(rng *rand.Rand) (int, []int64, [][2]int64) {
	n := rng.Intn(6) + 1
	size := 1 << uint(n)
	arr := make([]int64, size)
	for i := 0; i < size; i++ {
		arr[i] = rng.Int63n(1000)
	}
	q := rng.Intn(20) + 1
	queries := make([][2]int64, q)
	for i := 0; i < q; i++ {
		queries[i][0] = int64(rng.Intn(size))
		queries[i][1] = rng.Int63n(1000)
	}
	return n, arr, queries
}

func runCaseD(bin string, n int, arr []int64, queries [][2]int64) error {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d\n", n, len(queries))
	for i, v := range arr {
		if i > 0 {
			sb.WriteByte(' ')
		}
		fmt.Fprintf(&sb, "%d", v)
	}
	sb.WriteByte('\n')
	for _, q := range queries {
		fmt.Fprintf(&sb, "%d %d\n", q[0], q[1])
	}

	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(sb.String())
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	fields := strings.Fields(out.String())
	expected := solveD(n, append([]int64(nil), arr...), queries)
	if len(fields) != len(expected) {
		return fmt.Errorf("expected %d numbers got %d", len(expected), len(fields))
	}
	for i, f := range fields {
		val, err := strconv.ParseFloat(f, 64)
		if err != nil {
			return fmt.Errorf("bad float %q", f)
		}
		if absFloat(val-expected[i]) > 1e-6*maxFloat(1, absFloat(expected[i])) {
			return fmt.Errorf("expected %.7f got %.7f", expected[i], val)
		}
	}
	return nil
}

func absFloat(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}
func maxFloat(a, b float64) float64 {
	if a > b {
		return a
	}
	return b
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		n, arr, qs := genCaseD(rng)
		if err := runCaseD(bin, n, arr, qs); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
