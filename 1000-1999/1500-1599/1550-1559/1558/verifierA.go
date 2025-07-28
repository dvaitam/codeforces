package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"
)

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func possible(a, b int) []int {
	n := a + b
	set := make(map[int]struct{})
	for start := 0; start < 2; start++ {
		var serveA, serveB int
		if start == 0 {
			serveA = (n + 1) / 2
			serveB = n / 2
		} else {
			serveA = n / 2
			serveB = (n + 1) / 2
		}
		diff := a - serveA
		L := max(0, -diff)
		U := min(serveA, serveB-diff)
		if L > U {
			continue
		}
		for w := L; w <= U; w++ {
			k := 2*w + diff
			set[k] = struct{}{}
		}
	}
	res := make([]int, 0, len(set))
	for k := range set {
		res = append(res, k)
	}
	sort.Ints(res)
	return res
}

func runCase(bin string, a, b int) error {
	input := fmt.Sprintf("1\n%d %d\n", a, b)
	cmd := exec.Command(bin)
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	tokens := strings.Fields(out.String())
	exp := possible(a, b)
	if len(tokens) != 1+len(exp) {
		return fmt.Errorf("expected %d numbers, got %d", 1+len(exp), len(tokens))
	}
	m, err := strconv.Atoi(tokens[0])
	if err != nil {
		return fmt.Errorf("invalid integer %q", tokens[0])
	}
	if m != len(exp) {
		return fmt.Errorf("expected count %d got %d", len(exp), m)
	}
	for i, v := range exp {
		got, err := strconv.Atoi(tokens[1+i])
		if err != nil || got != v {
			return fmt.Errorf("expected %d at position %d got %q", v, i+1, tokens[1+i])
		}
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	f, err := os.Open("testcasesA.txt")
	if err != nil {
		fmt.Println("could not open testcasesA.txt:", err)
		os.Exit(1)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	scanner.Split(bufio.ScanWords)
	if !scanner.Scan() {
		fmt.Println("invalid test file")
		os.Exit(1)
	}
	t, _ := strconv.Atoi(scanner.Text())
	for i := 0; i < t; i++ {
		if !scanner.Scan() {
			fmt.Println("invalid test file")
			os.Exit(1)
		}
		a, _ := strconv.Atoi(scanner.Text())
		if !scanner.Scan() {
			fmt.Println("invalid test file")
			os.Exit(1)
		}
		b, _ := strconv.Atoi(scanner.Text())
		if err := runCase(bin, a, b); err != nil {
			fmt.Printf("case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
