package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math"
	"os"
	"os/exec"
	"strings"
)

func calc(n, k int64) int64 {
	var res int64
	for n > 0 {
		res += n / k
		n /= k
	}
	return res
}

func solveCase(n, b int64) string {
	bb := b
	ans := int64(math.MaxInt64)
	for i := int64(2); i*i <= bb; i++ {
		if bb%i == 0 {
			cnt := int64(0)
			for bb%i == 0 {
				cnt++
				bb /= i
			}
			val := calc(n, i) / cnt
			if val < ans {
				ans = val
			}
		}
	}
	if bb > 1 {
		val := calc(n, bb)
		if val < ans {
			ans = val
		}
	}
	return fmt.Sprint(ans)
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	data, err := os.ReadFile("testcasesC.txt")
	if err != nil {
		fmt.Println("failed to read testcasesC.txt:", err)
		os.Exit(1)
	}
	scanner := bufio.NewScanner(bytes.NewReader(data))
	idx := 0
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		idx++
		var n, b int64
		fmt.Sscan(line, &n, &b)
		expect := solveCase(n, b)
		cmd := exec.Command(bin)
		cmd.Stdin = strings.NewReader(line + "\n")
		var out bytes.Buffer
		var stderr bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = &stderr
		err := cmd.Run()
		if err != nil {
			fmt.Printf("test %d: runtime error: %v\nstderr: %s\n", idx, err, stderr.String())
			os.Exit(1)
		}
		got := strings.TrimSpace(out.String())
		if got != expect {
			fmt.Printf("test %d failed expected %s got %s\n", idx, expect, got)
			os.Exit(1)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Println("scanner error:", err)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", idx)
}
