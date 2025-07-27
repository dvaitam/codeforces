package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

func runBinary(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	return out.String(), err
}

func solveCaseB(arr []int) int {
	n := len(arr)
	freq := make([]int, n+2)
	idx := make([]int, n+2)
	for i, a := range arr {
		if freq[a] == 0 {
			freq[a] = 1
			idx[a] = i + 1
		} else if freq[a] == 1 {
			freq[a] = 2
		}
	}
	ans := -1
	for v := 1; v <= n; v++ {
		if freq[v] == 1 {
			ans = idx[v]
			break
		}
	}
	return ans
}

func generateTests() ([][]int, string) {
	const t = 100
	r := rand.New(rand.NewSource(2))
	arrays := make([][]int, t)
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", t)
	for i := 0; i < t; i++ {
		n := r.Intn(20) + 1
		fmt.Fprintf(&sb, "%d\n", n)
		arr := make([]int, n)
		for j := 0; j < n; j++ {
			arr[j] = r.Intn(n) + 1
			fmt.Fprintf(&sb, "%d ", arr[j])
		}
		fmt.Fprintln(&sb)
		arrays[i] = arr
	}
	return arrays, sb.String()
}

func verify(arrays [][]int, output string) error {
	scanner := bufio.NewScanner(strings.NewReader(output))
	scanner.Split(bufio.ScanWords)
	for idx, arr := range arrays {
		if !scanner.Scan() {
			return fmt.Errorf("case %d: missing output", idx+1)
		}
		var ans int
		fmt.Sscan(scanner.Text(), &ans)
		expected := solveCaseB(arr)
		if ans != expected {
			return fmt.Errorf("case %d: expected %d got %d", idx+1, expected, ans)
		}
	}
	if scanner.Scan() {
		return fmt.Errorf("extra output: %s", scanner.Text())
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go <binary>")
		os.Exit(1)
	}
	arrays, input := generateTests()
	out, err := runBinary(os.Args[1], input)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error running binary:", err)
		os.Exit(1)
	}
	if err := verify(arrays, out); err != nil {
		fmt.Fprintln(os.Stderr, "verification failed:", err)
		os.Exit(1)
	}
	fmt.Println("All tests passed for problem B")
}
