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
	// map each value in b to its position
	posB := make(map[int]int, lb)
	for i, v := range b {
		posB[v] = i
	}

	// map values of a to positions in b (or -1 if absent)
	pos := make([]int, la)
	for i, v := range a {
		if p, ok := posB[v]; ok {
			pos[i] = p
		} else {
			pos[i] = -1
		}
	}

	// handle rotation of a by doubling the array
	pos2 := make([]int, 2*la)
	copy(pos2, pos)
	copy(pos2[la:], pos)

	const INF int64 = -1
	arr := make([]int64, len(pos2)) // unwrapped positions in b
	var prev int64 = INF
	lb64 := int64(lb)
	for i, p := range pos2 {
		if p == -1 {
			arr[i] = INF
			prev = INF
			continue
		}
		cur := int64(p)
		if prev != INF && cur <= prev {
			cur += ((prev-cur)/lb64 + 1) * lb64
		}
		arr[i] = cur
		prev = cur
	}

	best := 0
	l := 0
	for r := 0; r < len(arr); r++ {
		if arr[r] == INF {
			l = r + 1
			continue
		}
		for l <= r && (arr[r]-arr[l] >= lb64 || r-l+1 > la) {
			l++
		}
		if r-l+1 > best {
			best = r - l + 1
		}
	}

	if best > la {
		best = la
	}
	if best > lb {
		best = lb
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
