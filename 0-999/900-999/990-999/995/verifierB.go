package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

type testCaseB struct {
	n   int
	arr []int
}

func solveB(arr []int) int {
	a := make([]int, len(arr))
	copy(a, arr)
	swaps := 0
	for i := 0; i < len(a); i += 2 {
		if a[i] == a[i+1] {
			continue
		}
		j := i + 1
		for j < len(a) && a[j] != a[i] {
			j++
		}
		for j > i+1 {
			a[j], a[j-1] = a[j-1], a[j]
			j--
			swaps++
		}
	}
	return swaps
}

func genCaseB(rng *rand.Rand) testCaseB {
	n := rng.Intn(20) + 1
	size := 2 * n
	arr := make([]int, size)
	// each id appears twice
	idxs := rng.Perm(size)
	for i := 0; i < n; i++ {
		arr[idxs[2*i]] = i + 1
		arr[idxs[2*i+1]] = i + 1
	}
	return testCaseB{n: n, arr: arr}
}

func runCaseB(bin string, tc testCaseB) error {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", tc.n)
	for i, v := range tc.arr {
		if i > 0 {
			sb.WriteByte(' ')
		}
		fmt.Fprintf(&sb, "%d", v)
	}
	sb.WriteByte('\n')

	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(sb.String())
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	var got int
	if _, err := fmt.Fscan(strings.NewReader(out.String()), &got); err != nil {
		return fmt.Errorf("bad output: %v", err)
	}
	exp := solveB(tc.arr)
	if got != exp {
		return fmt.Errorf("expected %d got %d", exp, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		tc := genCaseB(rng)
		if err := runCaseB(bin, tc); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
