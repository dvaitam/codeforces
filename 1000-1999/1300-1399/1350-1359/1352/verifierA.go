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

func expectedDigits(n int) []int {
	res := []int{}
	base := 1
	for n > 0 {
		d := n % 10
		if d != 0 {
			res = append(res, d*base)
		}
		n /= 10
		base *= 10
	}
	return res
}

func runCase(bin string, n int) error {
	input := fmt.Sprintf("1\n%d\n", n)
	out, err := runBinary(bin, input)
	if err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out)
	}
	fields := strings.Fields(out)
	if len(fields) == 0 {
		return fmt.Errorf("no output")
	}
	k, err := strconv.Atoi(fields[0])
	if err != nil {
		return fmt.Errorf("cannot parse k: %v", err)
	}
	nums := []int{}
	for _, f := range fields[1:] {
		v, err := strconv.Atoi(f)
		if err != nil {
			return fmt.Errorf("invalid number %q", f)
		}
		nums = append(nums, v)
	}
	exp := expectedDigits(n)
	if k != len(exp) || len(nums) != k {
		return fmt.Errorf("expected %d numbers got %d", len(exp), len(nums))
	}
	m := make(map[int]int)
	for _, v := range exp {
		m[v]++
	}
	for _, v := range nums {
		m[v]--
	}
	for _, c := range m {
		if c != 0 {
			return fmt.Errorf("numbers do not match expected")
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
	data, err := os.ReadFile("testcasesA.txt")
	if err != nil {
		fmt.Println("could not read testcasesA.txt:", err)
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
		if !sc.Scan() {
			fmt.Println("not enough test cases")
			os.Exit(1)
		}
		n, _ := strconv.Atoi(sc.Text())
		if err := runCase(bin, n); err != nil {
			fmt.Printf("case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
