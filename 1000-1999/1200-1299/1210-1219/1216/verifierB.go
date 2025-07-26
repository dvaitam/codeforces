package main

import (
	"bytes"
	"fmt"
	"io"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"
	"time"
)

func runCandidate(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func expectedCost(a []int) int64 {
	vals := append([]int(nil), a...)
	sort.Slice(vals, func(i, j int) bool { return vals[i] > vals[j] })
	var res int64
	for i, v := range vals {
		res += int64(i*v + 1)
	}
	return res
}

func genCase(rng *rand.Rand) (string, []int) {
	n := rng.Intn(10) + 1
	vals := make([]int, n)
	for i := 0; i < n; i++ {
		vals[i] = rng.Intn(20) + 1
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for i := 0; i < n; i++ {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(vals[i]))
	}
	sb.WriteByte('\n')
	return sb.String(), vals
}

func parseOutput(out string, n int) (int64, []int, error) {
	reader := strings.NewReader(out)
	var cost int64
	if _, err := fmt.Fscan(reader, &cost); err != nil {
		return 0, nil, fmt.Errorf("parse cost: %v", err)
	}
	perm := make([]int, n)
	for i := 0; i < n; i++ {
		if _, err := fmt.Fscan(reader, &perm[i]); err != nil {
			return 0, nil, fmt.Errorf("parse perm: %v", err)
		}
	}
	var extra string
	if _, err := fmt.Fscan(reader, &extra); err != io.EOF {
		return 0, nil, fmt.Errorf("extra output")
	}
	return cost, perm, nil
}

func checkCase(vals []int, out string) error {
	n := len(vals)
	cost, perm, err := parseOutput(out, n)
	if err != nil {
		return err
	}
	used := make([]bool, n)
	var candCost int64
	for i, v := range perm {
		if v < 1 || v > n || used[v-1] {
			return fmt.Errorf("invalid permutation")
		}
		used[v-1] = true
		candCost += int64(i*vals[v-1] + 1)
	}
	if cost != candCost {
		return fmt.Errorf("reported cost %d but computed %d", cost, candCost)
	}
	exp := expectedCost(vals)
	if candCost != exp {
		return fmt.Errorf("expected cost %d got %d", exp, candCost)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		return
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		input, vals := genCase(rng)
		out, err := runCandidate(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
		if err := checkCase(vals, out); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s\noutput:\n%s\n", i+1, err, input, out)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
