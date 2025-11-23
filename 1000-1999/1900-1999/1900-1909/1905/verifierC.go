package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"sort"
	"strings"
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
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func solve(s string) int {
	n := len(s)
	
	// Find the lexicographically largest subsequence
	// We can iterate from right to left to find this.
	// The elements included are those >= max suffix so far.
	type pair struct {
		val byte
		idx int
	}
	var seq []pair
	suffixMax := byte(0)
	for i := n - 1; i >= 0; i-- {
		if s[i] >= suffixMax {
			seq = append(seq, pair{s[i], i})
			suffixMax = s[i]
		}
	}
	// seq is in reverse order (right to left).
	// Let's reverse it to be left to right.
	for i, j := 0, len(seq)-1; i < j; i, j = i+1, j-1 {
		seq[i], seq[j] = seq[j], seq[i]
	}
	
	// Simulate the cyclic shift on the chosen indices.
	// The elements at indices seq[0].idx, seq[1].idx, ..., seq[k].idx
	// will be replaced by seq[k].val, seq[0].val, ..., seq[k-1].val respectively.
	// Effectively, the new char at seq[i].idx is seq[(i-1)%k].val if we index 0..k-1.
	// BUT: cyclic shift to RIGHT.
	// t_1 t_2 ... t_m -> t_m t_1 ... t_{m-1}
	// So position seq[0].idx gets seq[m-1].val
	// Position seq[1].idx gets seq[0].val
	// ...
	// Position seq[i].idx gets seq[i-1].val
	
	t := []byte(s)
	m := len(seq)
	for i := 0; i < m; i++ {
		prevIdx := (i - 1 + m) % m
		t[seq[i].idx] = seq[prevIdx].val
	}
	
	// Check if sorted
	if !sort.SliceIsSorted(t, func(i, j int) bool { return t[i] < t[j] }) {
		return -1
	}
	
	// If sorted, the number of operations is m - count(max_char in seq)
	// Actually, wait. The logic is simpler.
	// The subsequence is already non-increasing. 
	// Example: z z z y x c b a
	// Shifting right makes: a z z z y x c b (not sorted if rest is interlaced badly)
	// 
	// Actually, if the resulting string is sorted, then the cost is:
	// (number of elements in the subsequence) - (number of occurrences of the maximum character in the subsequence)
	
	maxVal := seq[0].val
	countMax := 0
	for _, p := range seq {
		if p.val == maxVal {
			countMax++
		}
	}
	
	return m - countMax
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	file, err := os.Open("testcasesC.txt")
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open testcases: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	idx := 0
	for scanner.Scan() {
		s := strings.TrimSpace(scanner.Text())
		if s == "" {
			continue
		}
		idx++
		
		// Calculate expected result using our correct solve function
		expected := solve(s)
		want := fmt.Sprintf("%d", expected)
		
		input := fmt.Sprintf("1\n%d\n%s\n", len(s), s)
		got, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", idx, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != want {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\nInput: %s\n", idx, want, got, s)
			os.Exit(1)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "scanner error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", idx)
}