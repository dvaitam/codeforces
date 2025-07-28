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

func solveCase(n int, nums []int) int {
	bestIdx := 1
	bestQual := -1
	for i := 0; i < n; i++ {
		a := nums[2*i]
		b := nums[2*i+1]
		if a <= 10 && b > bestQual {
			bestQual = b
			bestIdx = i + 1
		}
	}
	return bestIdx
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		return
	}
	bin := os.Args[1]

	file, err := os.Open("testcasesB.txt")
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
		if len(parts) < 1 {
			continue
		}
		n := 0
		fmt.Sscan(parts[0], &n)
		nums := make([]int, 2*n)
		for i := 0; i < 2*n; i++ {
			fmt.Sscan(parts[i+1], &nums[i])
		}
		input := fmt.Sprintf("1\n%d\n", n)
		for i := 0; i < n; i++ {
			input += fmt.Sprintf("%d %d\n", nums[2*i], nums[2*i+1])
		}
		inputs = append(inputs, input)
		exps = append(exps, solveCase(n, nums))
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
