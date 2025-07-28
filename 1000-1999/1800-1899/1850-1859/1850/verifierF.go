package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
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

func solveCase(arr []int) int {
	n := len(arr)
	freq := make([]int, n+1)
	for _, v := range arr {
		if v <= n {
			freq[v]++
		}
	}
	count := make([]int, n+1)
	maxCatch := 0
	for d := 1; d <= n; d++ {
		if freq[d] == 0 {
			continue
		}
		for m := d; m <= n; m += d {
			count[m] += freq[d]
			if count[m] > maxCatch {
				maxCatch = count[m]
			}
		}
	}
	return maxCatch
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierF.go /path/to/binary")
		return
	}
	bin := os.Args[1]

	file, err := os.Open("testcasesF.txt")
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
		var n int
		fmt.Sscan(parts[0], &n)
		arr := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Sscan(parts[i+1], &arr[i])
		}
		var sb strings.Builder
		fmt.Fprintf(&sb, "1\n%d\n", n)
		for i := 0; i < n; i++ {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprintf("%d", arr[i]))
		}
		sb.WriteByte('\n')
		inputs = append(inputs, sb.String())
		exps = append(exps, solveCase(arr))
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
