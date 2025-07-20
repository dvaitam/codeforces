package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/bits"
	"os"
	"os/exec"
	"strings"
)

func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func expectedD(n, c, k int, s string) int {
	tag := make([]bool, 1<<uint(c))
	sum := make([][]int, c)
	for j := 0; j < c; j++ {
		sum[j] = make([]int, n+1)
		for i := 1; i <= n; i++ {
			val := 0
			if int(s[i-1]-'A') == j {
				val = 1
			}
			sum[j][i] = sum[j][i-1] + val
		}
	}
	for i := 1; i <= n-k+1; i++ {
		t := 0
		for j := 0; j < c; j++ {
			if sum[j][i+k-1]-sum[j][i-1] == 0 {
				t |= 1 << uint(j)
			}
		}
		tag[t] = true
	}
	for i := (1 << uint(c)) - 1; i > 0; i-- {
		if tag[i] {
			for j := 0; j < c; j++ {
				if (i>>uint(j))&1 != 0 {
					tag[i^(1<<uint(j))] = true
				}
			}
		}
	}
	ans := int(1e9)
	last := int(s[len(s)-1] - 'A')
	for i := 0; i < (1 << uint(c)); i++ {
		if !tag[i] && ((i>>uint(last))&1) != 0 {
			cnt := bits.OnesCount(uint(i))
			ans = minInt(ans, cnt)
		}
	}
	return ans
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	file, err := os.Open("testcasesD.txt")
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open testcases: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	idx := 0
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		idx++
		var n, c, k int
		var s string
		fmt.Sscan(line, &n, &c, &k, &s)
		expect := expectedD(n, c, k, s)
		input := fmt.Sprintf("1\n%d %d %d %s\n", n, c, k, s)
		cmd := exec.Command(bin)
		cmd.Stdin = strings.NewReader(input)
		var out bytes.Buffer
		cmd.Stdout = &out
		var stderr bytes.Buffer
		cmd.Stderr = &stderr
		err := cmd.Run()
		if err != nil {
			fmt.Printf("test %d: runtime error: %v\nstderr: %s\n", idx, err, stderr.String())
			os.Exit(1)
		}
		res := strings.TrimSpace(out.String())
		var got int
		if _, err := fmt.Sscan(res, &got); err != nil {
			fmt.Printf("test %d: failed to parse output %q\n", idx, res)
			os.Exit(1)
		}
		if got != expect {
			fmt.Printf("test %d failed: expected %d got %d\n", idx, expect, got)
			os.Exit(1)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "scanner error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", idx)
}
