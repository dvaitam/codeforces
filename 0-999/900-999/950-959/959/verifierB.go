package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

func runCandidate(bin string, input string) (string, error) {
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

func solveB(data string) string {
	reader := bufio.NewReader(strings.NewReader(data))
	var n, k, m int
	fmt.Fscan(reader, &n, &k, &m)
	words := make([]string, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &words[i])
	}
	costs := make([]int64, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &costs[i])
	}
	group := make([]int, n)
	minCost := make([]int64, k)
	for g := 0; g < k; g++ {
		var x int
		fmt.Fscan(reader, &x)
		minVal := int64(^uint64(0) >> 1)
		for j := 0; j < x; j++ {
			var idx int
			fmt.Fscan(reader, &idx)
			idx--
			group[idx] = g
			if costs[idx] < minVal {
				minVal = costs[idx]
			}
		}
		minCost[g] = minVal
	}
	wordIndex := make(map[string]int, n)
	for i, w := range words {
		wordIndex[w] = i
	}
	var total int64
	for i := 0; i < m; i++ {
		var w string
		fmt.Fscan(reader, &w)
		idx := wordIndex[w]
		g := group[idx]
		total += minCost[g]
	}
	return fmt.Sprintf("%d", total)
}

func genCaseB(rng *rand.Rand) string {
	n := rng.Intn(10) + 1
	k := rng.Intn(n) + 1
	m := rng.Intn(10) + 1
	words := make([]string, n)
	for i := 0; i < n; i++ {
		words[i] = fmt.Sprintf("w%d", i+1)
	}
	costs := make([]int, n)
	for i := range costs {
		costs[i] = rng.Intn(1000) + 1
	}
	groups := make([][]int, k)
	perm := rng.Perm(n)
	for i := 0; i < k; i++ {
		groups[i] = append(groups[i], perm[i])
	}
	for i := k; i < n; i++ {
		g := rng.Intn(k)
		groups[g] = append(groups[g], perm[i])
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d %d\n", n, k, m))
	for i, w := range words {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(w)
	}
	sb.WriteByte('\n')
	for i, c := range costs {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(c))
	}
	sb.WriteByte('\n')
	for g := 0; g < k; g++ {
		sb.WriteString(strconv.Itoa(len(groups[g])))
		for _, idx := range groups[g] {
			sb.WriteByte(' ')
			sb.WriteString(strconv.Itoa(idx + 1))
		}
		sb.WriteByte('\n')
	}
	for i := 0; i < m; i++ {
		if i > 0 {
			sb.WriteByte(' ')
		}
		idx := rng.Intn(n)
		sb.WriteString(words[idx])
	}
	sb.WriteByte('\n')
	return sb.String()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in := genCaseB(rng)
		expect := solveB(in)
		got, err := runCandidate(bin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
		if got != expect {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %q got %q\ninput:\n%s", i+1, expect, got, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
