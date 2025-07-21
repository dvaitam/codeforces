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

func run(bin, input string) (string, error) {
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
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func solve(n, m int, prices []int, fruits []string) (int, int) {
	count := map[string]int{}
	for _, f := range fruits {
		count[f]++
	}
	counts := make([]int, 0, len(count))
	for _, c := range count {
		counts = append(counts, c)
	}
	sort.Sort(sort.Reverse(sort.IntSlice(counts)))
	sort.Ints(prices)
	k := len(counts)
	minSum := 0
	for i := 0; i < k; i++ {
		minSum += counts[i] * prices[i]
	}
	maxSum := 0
	for i := 0; i < k; i++ {
		maxSum += counts[i] * prices[n-1-i]
	}
	return minSum, maxSum
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	names := []string{"apple", "banana", "orange", "kiwi", "grape", "lemon", "melon", "pear", "plum", "mango"}
	for i := 0; i < 100; i++ {
		n := rng.Intn(8) + 2 // 2..9
		m := rng.Intn(8) + 2
		prices := make([]int, n)
		for j := range prices {
			prices[j] = rng.Intn(100) + 1
		}
		fruits := make([]string, m)
		for j := range fruits {
			fruits[j] = names[rng.Intn(n)]
		}
		input := fmt.Sprintf("%d %d\n", n, m)
		for j, p := range prices {
			if j > 0 {
				input += " "
			}
			input += strconv.Itoa(p)
		}
		input += "\n"
		for _, f := range fruits {
			input += f + "\n"
		}
		minSum, maxSum := solve(n, m, prices, fruits)
		expected := fmt.Sprintf("%d %d", minSum, maxSum)
		out, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
		if out != expected {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, expected, out, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
