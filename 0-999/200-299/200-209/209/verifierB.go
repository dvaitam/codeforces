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

func expected(a, b, c int64) int64 {
	counts := []int64{a, b, c}
	ans := int64(-1)
	for i := 0; i < 3; i++ {
		u := counts[(i+1)%3]
		v := counts[(i+2)%3]
		if (u+v)&1 == 1 {
			continue
		}
		f := u
		if v > u {
			f = v
		}
		if ans == -1 || f < ans {
			ans = f
		}
	}
	return ans
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	f, err := os.Open("testcasesB.txt")
	if err != nil {
		panic(err)
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
		parts := strings.Fields(line)
		if len(parts) < 3 {
			continue
		}
		a, _ := strconv.ParseInt(parts[0], 10, 64)
		b, _ := strconv.ParseInt(parts[1], 10, 64)
		c, _ := strconv.ParseInt(parts[2], 10, 64)
		exp := expected(a, b, c)
		input := fmt.Sprintf("%d %d %d\n", a, b, c)
		cmd := exec.Command(bin)
		cmd.Stdin = bytes.NewBufferString(input)
		out, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Printf("Test %d: runtime error: %v\n", idx, err)
			os.Exit(1)
		}
		gotStr := strings.TrimSpace(string(out))
		got, err2 := strconv.ParseInt(gotStr, 10, 64)
		if err2 != nil || got != exp {
			fmt.Printf("Test %d failed: expected %d got %s\n", idx, exp, gotStr)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", idx)
}
