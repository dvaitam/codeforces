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

func atoi(s string) int {
	v, _ := strconv.Atoi(s)
	return v
}

func solve(la, lb int, a, b []int) int {
	// Build map from element value to its position in b
	posB := make(map[int]int, lb)
	for i, v := range b {
		posB[v] = i
	}

	// Duplicate a to handle its circular nature
	aPrime := make([]int, 2*la)
	copy(aPrime, a)
	copy(aPrime[la:], a)

	m := lb
	var currentSum int
	l := 0
	best := 0

	for r := 0; r < len(aPrime); r++ {
		valR := aPrime[r]
		posR, ok := posB[valR]
		if !ok {
			// element not present in b, reset window
			l = r + 1
			currentSum = 0
			continue
		}

		if r > l {
			prevVal := aPrime[r-1]
			posPrev := posB[prevVal]
			dist := posR - posPrev
			dist %= m
			if dist < 0 {
				dist += m
			}
			currentSum += dist
		}

		for currentSum >= m || r-l+1 > la {
			if l < r {
				valL := aPrime[l]
				valNext := aPrime[l+1]
				posL := posB[valL]
				posNext := posB[valNext]
				dist := posNext - posL
				dist %= m
				if dist < 0 {
					dist += m
				}
				currentSum -= dist
			} else {
				currentSum = 0
			}
			l++
			if l > r {
				currentSum = 0
				break
			}
		}

		if r-l+1 > best {
			best = r - l + 1
		}
	}

	return best
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	file, err := os.Open("testcasesB.txt")
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
		fields := strings.Fields(line)
		if len(fields) < 2 {
			fmt.Printf("bad case %d\n", idx)
			os.Exit(1)
		}
		la := atoi(fields[0])
		lb := atoi(fields[1])
		needed := 2 + la + lb
		if len(fields) != needed {
			fmt.Printf("bad case %d\n", idx)
			os.Exit(1)
		}
		a := make([]int, la)
		for i := 0; i < la; i++ {
			a[i] = atoi(fields[2+i])
		}
		b := make([]int, lb)
		for i := 0; i < lb; i++ {
			b[i] = atoi(fields[2+la+i])
		}
		expected := solve(la, lb, a, b)

		var input strings.Builder
		input.WriteString(fields[0])
		input.WriteByte(' ')
		input.WriteString(fields[1])
		input.WriteByte('\n')
		for i := 0; i < la; i++ {
			if i > 0 {
				input.WriteByte(' ')
			}
			input.WriteString(fields[2+i])
		}
		input.WriteByte('\n')
		for i := 0; i < lb; i++ {
			if i > 0 {
				input.WriteByte(' ')
			}
			input.WriteString(fields[2+la+i])
		}
		input.WriteByte('\n')
		inputStr := input.String()

		cmd := exec.Command(bin)
		cmd.Stdin = strings.NewReader(inputStr)
		var out bytes.Buffer
		var stderr bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = &stderr
		if err := cmd.Run(); err != nil {
			fmt.Printf("test %d: runtime error: %v\nstderr: %s\n", idx, err, stderr.String())
			os.Exit(1)
		}
		got := strings.TrimSpace(out.String())
		if got != fmt.Sprintf("%d", expected) {
			fmt.Printf("test %d failed\nexpected:\n%d\n\ngot:\n%s\n", idx, expected, got)
			os.Exit(1)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "scanner error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", idx)
}
