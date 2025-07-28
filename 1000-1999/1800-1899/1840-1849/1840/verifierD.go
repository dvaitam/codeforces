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

func canCover(arr []int64, d int64) bool {
	n := len(arr)
	i := 0
	for k := 0; k < 3 && i < n; k++ {
		limit := arr[i] + 2*d
		i++
		for i < n && arr[i] <= limit {
			i++
		}
	}
	return i == n
}

func solve(arr []int64) int64 {
	sort.Slice(arr, func(i, j int) bool { return arr[i] < arr[j] })
	low, high := int64(0), int64(1_000_000_000)
	for low < high {
		mid := (low + high) / 2
		if canCover(arr, mid) {
			high = mid
		} else {
			low = mid + 1
		}
	}
	return low
}

func generateCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(10) + 1
	arr := make([]int64, n)
	for i := 0; i < n; i++ {
		arr[i] = rng.Int63n(1000) + 1
	}
	var sb strings.Builder
	sb.WriteString("1\n")
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for i := 0; i < n; i++ {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", arr[i]))
	}
	sb.WriteByte('\n')
	expected := fmt.Sprintf("%d", solve(arr))
	return sb.String(), expected
}

func runCase(bin, input, expected string) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errb bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errb
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, errb.String())
	}
	got := strings.TrimSpace(out.String())
	if got != expected {
		return fmt.Errorf("expected %s got %s", expected, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := generateCase(rng)
		if err := runCase(bin, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
