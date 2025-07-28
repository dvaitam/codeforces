package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// generateTests creates t test cases and returns them as a slice of strings in input format
func generateTests() []string {
	rand.Seed(1)
	t := 100
	tests := make([]string, t)
	for i := 0; i < t; i++ {
		n := rand.Intn(4) + 1 // 1..4
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d\n", n))
		for j := 0; j < n; j++ {
			k := rand.Intn(3) + 1 // 1..3
			used := map[int]bool{}
			sb.WriteString(fmt.Sprintf("%d", k))
			for l := 0; l < k; l++ {
				p := rand.Intn(6) + 1 // bits 1..6
				for used[p] {
					p = rand.Intn(6) + 1
				}
				used[p] = true
				sb.WriteString(fmt.Sprintf(" %d", p))
			}
			sb.WriteByte('\n')
		}
		tests[i] = sb.String()
	}
	return tests
}

func parseInput(in string) []int {
	scanner := bufio.NewScanner(strings.NewReader(in))
	scanner.Split(bufio.ScanWords)
	scanner.Scan()
	n := atoi(scanner.Text())
	arr := make([]int, n)
	for i := 0; i < n; i++ {
		scanner.Scan()
		k := atoi(scanner.Text())
		val := 0
		for j := 0; j < k; j++ {
			scanner.Scan()
			p := atoi(scanner.Text())
			val |= 1 << p
		}
		arr[i] = val
	}
	return arr
}

func atoi(s string) int {
	v, _ := strconv.Atoi(s)
	return v
}

func hasEqualSubseq(arr []int) bool {
	n := len(arr)
	seen := map[int]bool{}
	for mask := 1; mask < (1 << n); mask++ {
		or := 0
		for i := 0; i < n; i++ {
			if mask&(1<<i) != 0 {
				or |= arr[i]
			}
		}
		if seen[or] {
			return true
		}
		seen[or] = true
	}
	return false
}

func verifyCase(bin string, tc string) error {
	arr := parseInput(tc)
	exp := "No"
	if hasEqualSubseq(arr) {
		exp = "Yes"
	}
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader("1\n" + tc)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("execution error: %v", err)
	}
	res := strings.TrimSpace(string(out))
	res = strings.ToLower(res)
	if res != strings.ToLower(exp) {
		return fmt.Errorf("expected %s got %s", exp, res)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := generateTests()
	for i, tc := range tests {
		if err := verifyCase(bin, tc); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, tc)
			os.Exit(1)
		}
	}
	fmt.Println("all tests passed")
}
