package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/bits"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

const MODH int64 = 1000000007

func solveH(n, k int, arr []int) string {
	freq := make([]int, 64)
	for _, v := range arr {
		freq[v]++
	}
	pow2 := make([]int64, n+1)
	pow2[0] = 1
	for i := 1; i <= n; i++ {
		pow2[i] = pow2[i-1] * 2 % MODH
	}
	cntSuper := make([]int, 64)
	for mask := 0; mask < 64; mask++ {
		cnt := 0
		for val := 0; val < 64; val++ {
			if val&mask == mask {
				cnt += freq[val]
			}
		}
		cntSuper[mask] = cnt
	}
	dp := make([]int64, 64)
	for mask := 63; mask >= 0; mask-- {
		total := (pow2[cntSuper[mask]] - 1 + MODH) % MODH
		for sup := mask + 1; sup < 64; sup++ {
			if sup&mask == mask {
				total -= dp[sup]
				if total < 0 {
					total += MODH
				}
			}
		}
		dp[mask] = total
	}
	var result int64
	for mask := 0; mask < 64; mask++ {
		if bits.OnesCount(uint(mask)) == k {
			result += dp[mask]
			if result >= MODH {
				result -= MODH
			}
		}
	}
	return strconv.FormatInt(result, 10)
}

func genTestsH() ([]string, string) {
	const t = 100
	rand.Seed(1)
	var input strings.Builder
	fmt.Fprintln(&input, t)
	expected := make([]string, t)
	for i := 0; i < t; i++ {
		n := rand.Intn(20) + 1
		k := rand.Intn(7)
		fmt.Fprintf(&input, "%d %d\n", n, k)
		arr := make([]int, n)
		for j := 0; j < n; j++ {
			arr[j] = rand.Intn(64)
			if j+1 == n {
				fmt.Fprintln(&input, arr[j])
			} else {
				fmt.Fprint(&input, arr[j], " ")
			}
		}
		expected[i] = solveH(n, k, arr)
	}
	return expected, input.String()
}

func runBinary(path, in string) ([]string, error) {
	cmd := exec.Command(path)
	cmd.Stdin = strings.NewReader(in)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return nil, err
	}
	scanner := bufio.NewScanner(&out)
	var lines []string
	for scanner.Scan() {
		lines = append(lines, strings.TrimSpace(scanner.Text()))
	}
	return lines, scanner.Err()
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierH.go /path/to/binary")
		os.Exit(1)
	}
	expected, input := genTestsH()
	lines, err := runBinary(os.Args[1], input)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error running binary:", err)
		os.Exit(1)
	}
	if len(lines) != len(expected) {
		fmt.Fprintf(os.Stderr, "expected %d lines, got %d\n", len(expected), len(lines))
		os.Exit(1)
	}
	for i, exp := range expected {
		if lines[i] != exp {
			fmt.Fprintf(os.Stderr, "test %d failed: expected %s got %s\n", i+1, exp, lines[i])
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed!")
}
