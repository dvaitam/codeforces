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

func solve(n, k int, a []int) int {
	r := a[k-1]
	prefix := make([]int, n+1)
	for i := 1; i <= n; i++ {
		if a[i-1] > prefix[i-1] {
			prefix[i] = a[i-1]
		} else {
			prefix[i] = prefix[i-1]
		}
	}
	prefixWithout := make([]int, n+1)
	mx := 0
	for i := 1; i <= n; i++ {
		if i == k {
			prefixWithout[i] = mx
		} else {
			if a[i-1] > mx {
				mx = a[i-1]
			}
			prefixWithout[i] = mx
		}
	}
	ng := make([]int, n+2)
	nextIdx := n + 1
	for i := n; i >= 1; i-- {
		if a[i-1] > r {
			nextIdx = i
		}
		ng[i] = nextIdx
	}
	ans := 0
	for p := 1; p <= n; p++ {
		b := a[p-1]
		var pref int
		if p <= k {
			pref = prefix[p-1]
		} else {
			pref = prefixWithout[p-1]
			if k <= p-1 && b > pref {
				pref = b
			}
		}
		q := ng[p+1]
		if p < k && b > r && k < q {
			q = k
		}
		var wins int
		if p == 1 {
			if q == n+1 {
				wins = n - 1
			} else {
				wins = q - 2
			}
		} else {
			if r <= pref {
				wins = 0
			} else {
				if q == n+1 {
					wins = n - p + 1
				} else {
					wins = q - p
				}
			}
		}
		if wins > ans {
			ans = wins
		}
	}
	return ans
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	binary := os.Args[1]

	file, err := os.Open("testcasesB.txt")
	if err != nil {
		panic(err)
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
		fields := strings.Fields(line)
		n, _ := strconv.Atoi(fields[0])
		k, _ := strconv.Atoi(fields[1])
		arr := make([]int, n)
		for i := 0; i < n; i++ {
			v, _ := strconv.Atoi(fields[2+i])
			arr[i] = v
		}
		exp := solve(n, k, arr)

		var input strings.Builder
		input.WriteString("1\n")
		input.WriteString(fmt.Sprintf("%d %d\n", n, k))
		for i := 0; i < n; i++ {
			if i > 0 {
				input.WriteByte(' ')
			}
			input.WriteString(fmt.Sprintf("%d", arr[i]))
		}
		input.WriteByte('\n')

		cmd := exec.Command(binary)
		cmd.Stdin = strings.NewReader(input.String())
		var outBuf bytes.Buffer
		var errBuf bytes.Buffer
		cmd.Stdout = &outBuf
		cmd.Stderr = &errBuf
		err = cmd.Run()
		if err != nil {
			fmt.Printf("Test %d: runtime error: %v\nstderr: %s\n", idx, err, errBuf.String())
			os.Exit(1)
		}
		outStr := strings.TrimSpace(outBuf.String())
		if outStr != fmt.Sprintf("%d", exp) {
			fmt.Printf("Test %d failed: expected %d got %s\n", idx, exp, outStr)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", idx)
}
