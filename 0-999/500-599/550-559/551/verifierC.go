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

func can(a []int64, n, m int, t int64) bool {
	rem := make([]int64, n)
	copy(rem, a)
	cur := 0
	for s := 0; s < m && cur < n; s++ {
		op := t
		if op <= 0 {
			break
		}
		op--
		if op <= 0 {
			continue
		}
		if int64(cur) > op {
			continue
		}
		op -= int64(cur)
		for op > 0 && cur < n {
			if rem[cur] <= op {
				op -= rem[cur]
				rem[cur] = 0
				cur++
				if cur < n {
					if op <= 0 {
						break
					}
					op--
				}
			} else {
				rem[cur] -= op
				op = 0
				break
			}
		}
	}
	return cur >= n
}

func expected(a []int64, n, m int) int64 {
	var sum int64
	for _, v := range a {
		sum += v
	}
	lo, hi := int64(0), sum+int64(n)+5
	for lo+1 < hi {
		mid := (lo + hi) / 2
		if can(a, n, m, mid) {
			hi = mid
		} else {
			lo = mid
		}
	}
	return hi
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	f, err := os.Open("testcasesC.txt")
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not open testcases: %v\n", err)
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
		parts := strings.Fields(line)
		if len(parts) < 2 {
			fmt.Printf("case %d: invalid line\n", idx)
			os.Exit(1)
		}
		n, _ := strconv.Atoi(parts[0])
		m, _ := strconv.Atoi(parts[1])
		if len(parts) != 2+n {
			fmt.Printf("case %d: wrong number of piles\n", idx)
			os.Exit(1)
		}
		a := make([]int64, n)
		for i := 0; i < n; i++ {
			val, _ := strconv.Atoi(parts[2+i])
			a[i] = int64(val)
		}
		exp := expected(a, n, m)
		var input strings.Builder
		input.WriteString(fmt.Sprintf("%d %d\n", n, m))
		for i := 0; i < n; i++ {
			if i > 0 {
				input.WriteByte(' ')
			}
			input.WriteString(strconv.FormatInt(a[i], 10))
		}
		input.WriteByte('\n')
		cmd := exec.Command(bin)
		cmd.Stdin = strings.NewReader(input.String())
		var out bytes.Buffer
		var stderr bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = &stderr
		err := cmd.Run()
		if err != nil {
			fmt.Printf("case %d: runtime error: %v\nstderr: %s\n", idx, err, stderr.String())
			os.Exit(1)
		}
		gotStr := strings.TrimSpace(out.String())
		got, err := strconv.ParseInt(gotStr, 10, 64)
		if err != nil || got != exp {
			fmt.Printf("case %d failed: expected %d got %s\n", idx, exp, gotStr)
			os.Exit(1)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "scanner error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", idx)
}
