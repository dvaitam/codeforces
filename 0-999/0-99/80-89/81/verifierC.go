package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"
)

func expected(n, a, b int, t []int) []int {
	f := make([]int, n)
	if a == b {
		for i := 0; i < n; i++ {
			if i < a {
				f[i] = 1
			} else {
				f[i] = 2
			}
		}
	} else {
		pickLargest := b > a
		u := make([]int, n)
		copy(u, t)
		sort.Ints(u)
		if pickLargest {
			threshold := u[n-a]
			cntGreater := 0
			for _, v := range t {
				if v > threshold {
					cntGreater++
				}
			}
			eqSelect := a - cntGreater
			for i, v := range t {
				if v > threshold {
					f[i] = 1
				} else if v == threshold && eqSelect > 0 {
					f[i] = 1
					eqSelect--
				} else {
					f[i] = 2
				}
			}
		} else {
			threshold := u[a-1]
			cntLess := 0
			for _, v := range t {
				if v < threshold {
					cntLess++
				}
			}
			eqSelect := a - cntLess
			for i, v := range t {
				if v < threshold {
					f[i] = 1
				} else if v == threshold && eqSelect > 0 {
					f[i] = 1
					eqSelect--
				} else {
					f[i] = 2
				}
			}
		}
	}
	return f
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	file, err := os.Open("testcasesC.txt")
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
		parts := strings.Fields(line)
		if len(parts) < 3 {
			fmt.Printf("test %d: invalid line\n", idx)
			os.Exit(1)
		}
		nVal, _ := strconv.Atoi(parts[0])
		aVal, _ := strconv.Atoi(parts[1])
		bVal, _ := strconv.Atoi(parts[2])
		if len(parts) != nVal+3 {
			fmt.Printf("test %d: wrong number of values\n", idx)
			os.Exit(1)
		}
		arr := make([]int, nVal)
		for i := 0; i < nVal; i++ {
			v, _ := strconv.Atoi(parts[3+i])
			arr[i] = v
		}
		expectArr := expected(nVal, aVal, bVal, arr)
		var input strings.Builder
		input.WriteString(fmt.Sprintf("%d %d %d\n", nVal, aVal, bVal))
		for i := 0; i < nVal; i++ {
			if i > 0 {
				input.WriteByte(' ')
			}
			input.WriteString(strconv.Itoa(arr[i]))
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
			fmt.Printf("test %d: runtime error: %v\nstderr: %s\n", idx, err, stderr.String())
			os.Exit(1)
		}
		outVals := strings.Fields(strings.TrimSpace(out.String()))
		if len(outVals) != nVal {
			fmt.Printf("test %d: expected %d numbers got %d\n", idx, nVal, len(outVals))
			os.Exit(1)
		}
		for i := 0; i < nVal; i++ {
			got, _ := strconv.Atoi(outVals[i])
			if got != expectArr[i] {
				fmt.Printf("test %d failed at position %d: expected %d got %d\n", idx, i+1, expectArr[i], got)
				os.Exit(1)
			}
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "scanner error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", idx)
}
