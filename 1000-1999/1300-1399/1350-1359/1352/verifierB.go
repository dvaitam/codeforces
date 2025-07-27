package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func runBinary(path, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(path, ".go") {
		cmd = exec.Command("go", "run", path)
	} else {
		cmd = exec.Command(path)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

func validRepresentation(n, k int, nums []int) bool {
	if len(nums) != k {
		return false
	}
	sum := 0
	parity := nums[0] % 2
	for _, v := range nums {
		if v <= 0 || v%2 != parity {
			return false
		}
		sum += v
	}
	return sum == n
}

func possible(n, k int) bool {
	if n%2 == k%2 && n >= k {
		return true
	}
	if n%2 == 0 && n >= 2*k {
		return true
	}
	return false
}

func runCase(bin string, n, k int) error {
	input := fmt.Sprintf("1\n%d %d\n", n, k)
	out, err := runBinary(bin, input)
	if err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out)
	}
	fields := strings.Fields(out)
	if len(fields) == 0 {
		return fmt.Errorf("no output")
	}
	if strings.ToLower(fields[0]) == "no" {
		if possible(n, k) {
			return fmt.Errorf("expected YES but got NO")
		}
		return nil
	}
	if strings.ToLower(fields[0]) != "yes" {
		return fmt.Errorf("first word should be YES or NO")
	}
	if !possible(n, k) {
		return fmt.Errorf("expected NO but got YES")
	}
	nums := make([]int, 0, k)
	for _, f := range fields[1:] {
		v, err := strconv.Atoi(f)
		if err != nil {
			return fmt.Errorf("invalid number %q", f)
		}
		nums = append(nums, v)
	}
	if !validRepresentation(n, k, nums) {
		return fmt.Errorf("invalid representation")
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	data, err := os.ReadFile("testcasesB.txt")
	if err != nil {
		fmt.Println("could not read testcasesB.txt:", err)
		os.Exit(1)
	}
	sc := bufio.NewScanner(bytes.NewReader(data))
	sc.Split(bufio.ScanWords)
	if !sc.Scan() {
		fmt.Println("invalid test file")
		os.Exit(1)
	}
	t, _ := strconv.Atoi(sc.Text())
	for i := 0; i < t; i++ {
		sc.Scan()
		n, _ := strconv.Atoi(sc.Text())
		sc.Scan()
		k, _ := strconv.Atoi(sc.Text())
		if err := runCase(bin, n, k); err != nil {
			fmt.Printf("case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
