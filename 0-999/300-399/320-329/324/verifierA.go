package main

import (
	"bytes"
	"fmt"
	"math"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

type testCaseA struct {
	input string
	a     []int64
	sum   int64
}

func solveA(a []int64) int64 {
	n := len(a) - 1 // 1-indexed for easier logic
	posPref := make([]int64, n+1)
	for i := 1; i <= n; i++ {
		posPref[i] = posPref[i-1]
		if a[i] > 0 {
			posPref[i] += a[i]
		}
	}
	posMap := make(map[int64][]int)
	for i := 1; i <= n; i++ {
		v := a[i]
		posMap[v] = append(posMap[v], i)
	}
	bestSum := int64(math.MinInt64)
	var bestL, bestR int
	for v, vec := range posMap {
		if len(vec) < 2 {
			continue
		}
		l := vec[0]
		r := vec[len(vec)-1]
		midSum := int64(0)
		if r > l+1 {
			midSum = posPref[r-1] - posPref[l]
		}
		total := int64(v) + int64(v) + midSum
		if total > bestSum {
			bestSum = total
			bestL = l
			bestR = r
		}
	}
	// compute sum according to bestL,bestR
	sum := a[bestL] + a[bestR]
	for i := bestL + 1; i < bestR; i++ {
		if a[i] > 0 {
			sum += a[i]
		}
	}
	return sum
}

func genCaseA(rng *rand.Rand) testCaseA {
	n := rng.Intn(8) + 2
	arr := make([]int64, n+1)
	for i := 1; i <= n; i++ {
		arr[i] = int64(rng.Intn(21) - 10)
	}
	// ensure at least two equal numbers
	if n >= 2 {
		arr[2] = arr[1]
	}
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", n)
	for i := 1; i <= n; i++ {
		if i > 1 {
			sb.WriteByte(' ')
		}
		fmt.Fprintf(&sb, "%d", arr[i])
	}
	sb.WriteByte('\n')
	sum := solveA(arr)
	return testCaseA{input: sb.String(), a: arr, sum: sum}
}

func runCaseA(bin string, tc testCaseA) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(tc.input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	fields := strings.Fields(out.String())
	if len(fields) < 2 {
		return fmt.Errorf("bad output")
	}
	var gotSum int64
	var k int
	if _, err := fmt.Sscan(fields[0], &gotSum); err != nil {
		return fmt.Errorf("bad output: %v", err)
	}
	if _, err := fmt.Sscan(fields[1], &k); err != nil {
		return fmt.Errorf("bad output: %v", err)
	}
	if len(fields) != 2+k {
		return fmt.Errorf("expected %d indices, got %d", k, len(fields)-2)
	}
	removed := make(map[int]bool)
	for i := 0; i < k; i++ {
		var idx int
		if _, err := fmt.Sscan(fields[2+i], &idx); err != nil {
			return fmt.Errorf("bad index: %v", err)
		}
		if idx < 1 || idx >= len(tc.a) {
			return fmt.Errorf("index out of range: %d", idx)
		}
		if removed[idx] {
			return fmt.Errorf("duplicate index %d", idx)
		}
		removed[idx] = true
	}
	sum := int64(0)
	for i := 1; i < len(tc.a); i++ {
		if !removed[i] {
			sum += tc.a[i]
		}
	}
	if sum != gotSum {
		return fmt.Errorf("reported sum %d does not match computed sum %d", gotSum, sum)
	}
	if gotSum != tc.sum {
		return fmt.Errorf("expected sum %d got %d", tc.sum, gotSum)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		tc := genCaseA(rng)
		if err := runCaseA(bin, tc); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, tc.input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
