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
	v   int
	idx int
}

type testCase struct {
	input    string
	expected string
}

func expectedAnswer(a []int) (int, int) {
	n := len(a)
	arr := make([]pair, n)
	for i, v := range a {
		arr[i] = pair{v, i}
	}
	sort.Slice(arr, func(i, j int) bool {
		if arr[i].v != arr[j].v {
			return arr[i].v < arr[j].v
		}
		return arr[i].idx < arr[j].idx
	})
	maxCount := 1
	retL, retR := 0, 0
	l := 0
	r := 0
	for r < n {
		if arr[l].v == arr[r].v {
			r++
		} else {
			count := r - l
			if count > maxCount || (count == maxCount && arr[r-1].idx-arr[l].idx < arr[retR].idx-arr[retL].idx) {
				maxCount = count
				retL = l
				retR = r - 1
			}
			l = r
		}
	}
	if r-l > maxCount || (r-l == maxCount && arr[r-1].idx-arr[l].idx < arr[retR].idx-arr[retL].idx) {
		retL = l
		retR = r - 1
	}
	return arr[retL].idx + 1, arr[retR].idx + 1
}

func generateCase(rng *rand.Rand) testCase {
	n := rng.Intn(20) + 1
	arr := make([]int, n)
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for i := 0; i < n; i++ {
		arr[i] = rng.Intn(20) + 1
		sb.WriteString(fmt.Sprintf("%d ", arr[i]))
	}
	sb.WriteString("\n")
	l, r := expectedAnswer(arr)
	return testCase{input: sb.String(), expected: fmt.Sprintf("%d %d", l, r)}
}

func runCase(bin string, tc testCase) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(tc.input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(out.String())
	if got != tc.expected {
		return fmt.Errorf("expected %s got %s", tc.expected, got)
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
		tc := generateCase(rng)
		if err := runCase(bin, tc); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, tc.input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
