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

func gcd(a, b int) int {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

// generateCase builds a random array and returns the input string and
// the sorted gcd table for validation along with n.
func generateCase(rng *rand.Rand) (string, []int, int) {
	n := rng.Intn(4) + 1
	arr := make([]int, n)
	for i := range arr {
		arr[i] = rng.Intn(20) + 1
	}
	table := make([]int, 0, n*n)
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			table = append(table, gcd(arr[i], arr[j]))
		}
	}
	rng.Shuffle(len(table), func(i, j int) { table[i], table[j] = table[j], table[i] })

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for i, v := range table {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(v))
	}
	sb.WriteByte('\n')

	sorted := append([]int(nil), table...)
	sort.Ints(sorted)
	return sb.String(), sorted, n
}

func runCase(bin string, input string, expect []int, n int) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	fields := strings.Fields(out.String())
	if len(fields) != n {
		return fmt.Errorf("expected %d numbers got %d", n, len(fields))
	}
	arr := make([]int, n)
	for i, f := range fields {
		v, err := strconv.Atoi(f)
		if err != nil {
			return fmt.Errorf("invalid integer %q", f)
		}
		arr[i] = v
	}
	table := make([]int, 0, n*n)
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			table = append(table, gcd(arr[i], arr[j]))
		}
	}
	sort.Ints(table)
	for i := range table {
		if table[i] != expect[i] {
			return fmt.Errorf("wrong answer")
		}
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
		input, expect, n := generateCase(rng)
		if err := runCase(bin, input, expect, n); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
