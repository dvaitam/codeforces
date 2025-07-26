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

const MOD = 1000000007

func solveCase(n int, arr []int) string {
	freq := make([]int, 71)
	for _, v := range arr {
		if v >= 1 && v <= 70 {
			freq[v]++
		}
	}
	primes := []int{2, 3, 5, 7, 11, 13, 17, 19, 23, 29, 31, 37, 41, 43, 47, 53, 59, 61, 67}
	m := len(primes)
	mask := make([]int, 71)
	for v := 1; v <= 70; v++ {
		cur := v
		mm := 0
		for i, p := range primes {
			cnt := 0
			for cur%p == 0 {
				cur /= p
				cnt ^= 1
			}
			if cnt == 1 {
				mm |= 1 << i
			}
		}
		mask[v] = mm
	}
	pow2 := make([]int, n+1)
	pow2[0] = 1
	for i := 1; i <= n; i++ {
		pow2[i] = (pow2[i-1] * 2) % MOD
	}
	size := 1 << m
	dp := make([]int, size)
	dp[0] = 1
	for v := 1; v <= 70; v++ {
		c := freq[v]
		if c == 0 {
			continue
		}
		even := pow2[c-1]
		odd := pow2[c-1]
		msk := mask[v]
		newDP := make([]int, size)
		for j := 0; j < size; j++ {
			if dp[j] == 0 {
				continue
			}
			val := dp[j]
			ne := (val * even) % MOD
			newDP[j] = (newDP[j] + ne) % MOD
			no := (val * odd) % MOD
			j2 := j ^ msk
			newDP[j2] = (newDP[j2] + no) % MOD
		}
		dp = newDP
	}
	ans := dp[0] - 1
	if ans < 0 {
		ans += MOD
	}
	return fmt.Sprintf("%d", ans)
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	data, err := os.ReadFile("testcasesC.txt")
	if err != nil {
		fmt.Println("could not read testcasesC.txt:", err)
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
		fields := strings.Fields(line)
		if len(fields) < 1 {
			fmt.Println("bad test case")
			os.Exit(1)
		}
		n, _ := strconv.Atoi(fields[0])
		if len(fields)-1 != n {
			fmt.Println("bad test case length")
			os.Exit(1)
		}
		arr := make([]int, n)
		for i := 0; i < n; i++ {
			v, _ := strconv.Atoi(fields[1+i])
			arr[i] = v
		}
		expected := solveCase(n, arr)
		input := fmt.Sprintf("%d\n%s\n", n, strings.Join(fields[1:], " "))
		cmd := exec.Command(bin)
		cmd.Stdin = strings.NewReader(input)
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
		if got != expected {
			fmt.Printf("test %d failed: expected %s got %s\n", idx, expected, got)
			os.Exit(1)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Println("scanner error:", err)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", idx)
}
