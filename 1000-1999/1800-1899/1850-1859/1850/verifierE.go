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

func calc(arr []int64, w int64) int64 {
	var sum int64
	for _, v := range arr {
		d := v + 2*w
		sum += d * d
		if sum > 1<<63-1 {
			return sum
		}
	}
	return sum
}

func findW(arr []int64, c int64) int64 {
	low := int64(0)
	high := int64(1)
	for calc(arr, high) < c {
		high <<= 1
	}
	for low+1 < high {
		mid := (low + high) / 2
		if calc(arr, mid) >= c {
			high = mid
		} else {
			low = mid
		}
	}
	return high
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		return
	}
	bin := os.Args[1]

	file, err := os.Open("testcasesE.txt")
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open testcases: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var inputs []string
	var exps []int64
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		parts := strings.Fields(line)
		var n int
		fmt.Sscan(parts[0], &n)
		var c int64
		fmt.Sscan(parts[1], &c)
		arr := make([]int64, n)
		for i := 0; i < n; i++ {
			fmt.Sscan(parts[i+2], &arr[i])
		}
		input := fmt.Sprintf("1\n%d %d\n", n, c)
		for i := 0; i < n; i++ {
			if i > 0 {
				input += " "
			}
			input += fmt.Sprintf("%d", arr[i])
		}
		input += "\n"
		inputs = append(inputs, input)
		exps = append(exps, findW(arr, c))
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
		var got int64
		fmt.Fscan(strings.NewReader(out), &got)
		if got != exps[idx] {
			fmt.Fprintf(os.Stderr, "case %d failed\ninput:\n%sexpected %d got %d\n", idx+1, input, exps[idx], got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
