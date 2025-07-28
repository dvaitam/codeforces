package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strings"
	"time"
)

type pair struct {
	diff int64
	idx  int
}

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

func abs64(x int64) int64 {
	if x < 0 {
		return -x
	}
	return x
}

func winCount(x int64, arr []int64, m int64, k int) int64 {
	n := len(arr)
	arr2 := append([]int64(nil), arr...)
	arr2 = append(arr2, x)
	count := int64(0)
	for t := int64(0); t <= m; t++ {
		diffs := make([]pair, n+1)
		for i := 0; i < n+1; i++ {
			diffs[i] = pair{diff: abs64(arr2[i] - t), idx: i}
		}
		sort.Slice(diffs, func(i, j int) bool {
			if diffs[i].diff == diffs[j].diff {
				return diffs[i].idx < diffs[j].idx
			}
			return diffs[i].diff < diffs[j].diff
		})
		winners := diffs
		if len(winners) > k {
			winners = winners[:k]
		}
		for _, p := range winners {
			if p.idx == n { // Bytek index
				count++
				break
			}
		}
	}
	return count
}

func bruteB(n int, m int64, k int, arr []int64) (int64, int64) {
	var bestX int64
	bestVal := int64(-1)
	for x := int64(0); x <= m; x++ {
		val := winCount(x, arr, m, k)
		if val > bestVal || (val == bestVal && x < bestX) {
			bestVal = val
			bestX = x
		}
	}
	return bestVal, bestX
}

func generateCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(4) + 2
	m := int64(rng.Intn(10) + 5)
	k := rng.Intn(n) + 1
	arr := make([]int64, n)
	for i := 0; i < n; i++ {
		arr[i] = int64(rng.Intn(int(m + 1)))
	}
	val, x := bruteB(n, m, k, arr)
	// build input
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d %d\n", n, m, k)
	for i := 0; i < n; i++ {
		if i > 0 {
			sb.WriteByte(' ')
		}
		fmt.Fprintf(&sb, "%d", arr[i])
	}
	sb.WriteByte('\n')
	input := sb.String()
	expected := fmt.Sprintf("%d %d", val, x)
	return input, expected
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := generateCase(rng)
		out, err := runCandidate(bin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:%s", i+1, err, in)
			os.Exit(1)
		}
		if strings.TrimSpace(out) != strings.TrimSpace(exp) {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:%s", i+1, exp, out, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
