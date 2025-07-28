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

func run(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func good(p []int) int {
	n := len(p)
	prefixMax := make([]int, n)
	maxVal := 0
	for i := 0; i < n; i++ {
		if p[i] > maxVal {
			maxVal = p[i]
		}
		prefixMax[i] = maxVal
	}
	suffixMin := make([]int, n)
	minVal := n + 1
	for i := n - 1; i >= 0; i-- {
		if p[i] < minVal {
			minVal = p[i]
		}
		suffixMin[i] = minVal
	}
	cnt := 0
	for i := 0; i < n; i++ {
		leftMax := 0
		if i > 0 {
			leftMax = prefixMax[i-1]
		}
		rightMin := n + 1
		if i+1 < n {
			rightMin = suffixMin[i+1]
		}
		if leftMax < p[i] && p[i] < rightMin {
			cnt++
		}
	}
	return cnt
}

func brute(p []int) int {
	n := len(p)
	best := 0
	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			p[i], p[j] = p[j], p[i]
			if g := good(p); g > best {
				best = g
			}
			p[i], p[j] = p[j], p[i]
		}
	}
	return best
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	file, err := os.Open("testcasesF.txt")
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open testcases: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	idx := 0
	for {
		if !scanner.Scan() {
			break
		}
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		idx++
		n, _ := strconv.Atoi(line)
		scanner.Scan()
		pFields := strings.Fields(scanner.Text())
		p := make([]int, n)
		for i := 0; i < n; i++ {
			v, _ := strconv.Atoi(pFields[i])
			p[i] = v
		}
		scanner.Scan() // blank line
		want := fmt.Sprintf("%d", brute(append([]int(nil), p...)))
		var sb strings.Builder
		sb.WriteString("1\n")
		sb.WriteString(fmt.Sprintf("%d\n", n))
		for i, v := range p {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.Itoa(v))
		}
		sb.WriteByte('\n')
		input := sb.String()
		got, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", idx, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != want {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\n", idx, want, got)
			os.Exit(1)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "scanner error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", idx)
}
