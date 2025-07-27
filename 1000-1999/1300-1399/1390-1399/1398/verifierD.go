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

func brute(r, g, b []int) int {
	sort.Slice(r, func(i, j int) bool { return r[i] > r[j] })
	sort.Slice(g, func(i, j int) bool { return g[i] > g[j] })
	sort.Slice(b, func(i, j int) bool { return b[i] > b[j] })
	memo := make(map[[3]int]int)
	var f func(int, int, int) int
	f = func(i, j, k int) int {
		key := [3]int{i, j, k}
		if v, ok := memo[key]; ok {
			return v
		}
		best := 0
		if i < len(r) && j < len(g) {
			v := r[i]*g[j] + f(i+1, j+1, k)
			if v > best {
				best = v
			}
		}
		if i < len(r) && k < len(b) {
			v := r[i]*b[k] + f(i+1, j, k+1)
			if v > best {
				best = v
			}
		}
		if j < len(g) && k < len(b) {
			v := g[j]*b[k] + f(i, j+1, k+1)
			if v > best {
				best = v
			}
		}
		memo[key] = best
		return best
	}
	return f(0, 0, 0)
}

func solveCase(input string) string {
	fields := strings.Fields(input)
	idx := 0
	R, _ := strconv.Atoi(fields[idx])
	idx++
	G, _ := strconv.Atoi(fields[idx])
	idx++
	B, _ := strconv.Atoi(fields[idx])
	idx++
	r := make([]int, R)
	gArr := make([]int, G)
	bArr := make([]int, B)
	for i := 0; i < R; i++ {
		val, _ := strconv.Atoi(fields[idx])
		idx++
		r[i] = val
	}
	for i := 0; i < G; i++ {
		val, _ := strconv.Atoi(fields[idx])
		idx++
		gArr[i] = val
	}
	for i := 0; i < B; i++ {
		val, _ := strconv.Atoi(fields[idx])
		idx++
		bArr[i] = val
	}
	ans := brute(r, gArr, bArr)
	return fmt.Sprint(ans)
}

func generateCase(rng *rand.Rand) string {
	R := rng.Intn(3) + 1
	G := rng.Intn(3) + 1
	B := rng.Intn(3) + 1
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d %d\n", R, G, B))
	for i := 0; i < R; i++ {
		sb.WriteString(fmt.Sprintf("%d ", rng.Intn(9)+1))
	}
	sb.WriteByte('\n')
	for i := 0; i < G; i++ {
		sb.WriteString(fmt.Sprintf("%d ", rng.Intn(9)+1))
	}
	sb.WriteByte('\n')
	for i := 0; i < B; i++ {
		sb.WriteString(fmt.Sprintf("%d ", rng.Intn(9)+1))
	}
	sb.WriteByte('\n')
	return sb.String()
}

func runCase(bin, input, expected string) error {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(out.String())
	expected = strings.TrimSpace(expected)
	if got != expected {
		return fmt.Errorf("expected %q got %q", expected, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in := generateCase(rng)
		exp := solveCase(in)
		if err := runCase(bin, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
