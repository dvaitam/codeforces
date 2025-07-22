package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func isPrime(n int64) bool {
	if n < 2 {
		return false
	}
	for i := int64(2); i*i <= n; i++ {
		if n%i == 0 {
			return false
		}
	}
	return true
}

func solveB(nums []int64) string {
	var sb strings.Builder
	for _, x := range nums {
		root := int64(math.Round(math.Sqrt(float64(x))))
		if root*root == x && isPrime(root) {
			sb.WriteString("YES\n")
		} else {
			sb.WriteString("NO\n")
		}
	}
	return strings.TrimRight(sb.String(), "\n")
}

func runCase(bin string, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	f, err := os.Open("testcasesB.txt")
	if err != nil {
		fmt.Println("could not open testcasesB.txt:", err)
		os.Exit(1)
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	idx := 0
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		idx++
		fields := strings.Fields(line)
		n, _ := strconv.Atoi(fields[0])
		if len(fields) != 1+n {
			fmt.Printf("case %d: expected %d numbers, got %d\n", idx, n, len(fields)-1)
			os.Exit(1)
		}
		nums := make([]int64, n)
		for i := 0; i < n; i++ {
			v, _ := strconv.ParseInt(fields[1+i], 10, 64)
			nums[i] = v
		}
		input := fmt.Sprintf("%d\n", n)
		for i, v := range nums {
			if i > 0 {
				input += " "
			}
			input += fmt.Sprintf("%d", v)
		}
		input += "\n"
		expect := solveB(nums)
		got, err := runCase(bin, input)
		if err != nil {
			fmt.Printf("case %d: runtime error: %v\n", idx, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != expect {
			fmt.Printf("case %d failed: expected %q got %q\n", idx, expect, got)
			os.Exit(1)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Println("read error:", err)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", idx)
}
