package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"sort"
	"strings"
)

func runBinary(bin string, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	out, err := cmd.CombinedOutput()
	return string(out), err
}

func solveCase(n, k int, arr []int) int {
	sort.Ints(arr)
	maxLen, cur := 1, 1
	for i := 1; i < n; i++ {
		if arr[i]-arr[i-1] <= k {
			cur++
		} else {
			if cur > maxLen {
				maxLen = cur
			}
			cur = 1
		}
	}
	if cur > maxLen {
		maxLen = cur
	}
	return n - maxLen
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		return
	}
	bin := os.Args[1]

	file, err := os.Open("testcasesD.txt")
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open testcases: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var inputs []string
	var exps []int
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		parts := strings.Fields(line)
		if len(parts) < 2 {
			fmt.Fprintf(os.Stderr, "invalid line: %s\n", line)
			os.Exit(1)
		}
		var n, k int
		fmt.Sscan(parts[0], &n)
		fmt.Sscan(parts[1], &k)
		arr := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Sscan(parts[i+2], &arr[i])
		}
		input := fmt.Sprintf("1\n%d %d\n", n, k)
		for i := 0; i < n; i++ {
			if i > 0 {
				input += " "
			}
			input += fmt.Sprintf("%d", arr[i])
		}
		input += "\n"
		inputs = append(inputs, input)
		exps = append(exps, solveCase(n, k, arr))
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "scanner error: %v\n", err)
		os.Exit(1)
	}

	for idx, input := range inputs {
		out, err := runBinary(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d runtime error: %v\n%s", idx+1, err, out)
			os.Exit(1)
		}
		var got int
		fmt.Fscan(strings.NewReader(out), &got)
		if got != exps[idx] {
			fmt.Fprintf(os.Stderr, "case %d failed\ninput:\n%sexpected %d got %d\n", idx+1, input, exps[idx], got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
