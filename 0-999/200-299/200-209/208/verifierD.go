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

func expectedAnswer(n int, p []int64, costs [5]int64) (string, string) {
	counts := make([]int64, 5)
	var points int64
	for _, v := range p {
		points += v
		for {
			switch {
			case points >= costs[4]:
				points -= costs[4]
				counts[4]++
			case points >= costs[3]:
				points -= costs[3]
				counts[3]++
			case points >= costs[2]:
				points -= costs[2]
				counts[2]++
			case points >= costs[1]:
				points -= costs[1]
				counts[1]++
			case points >= costs[0]:
				points -= costs[0]
				counts[0]++
			default:
				goto done
			}
		}
	done:
	}
	var cntStr strings.Builder
	for i, c := range counts {
		if i > 0 {
			cntStr.WriteByte(' ')
		}
		cntStr.WriteString(fmt.Sprint(c))
	}
	return cntStr.String(), fmt.Sprint(points)
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
		fields := strings.Fields(line)
		if len(fields) < 6 {
			fmt.Printf("test %d: invalid line\n", idx)
			os.Exit(1)
		}
		n, _ := strconv.Atoi(fields[0])
		if len(fields) != 1+n+5 {
			fmt.Printf("test %d: wrong field count\n", idx)
			os.Exit(1)
		}
		p := make([]int64, n)
		for i := 0; i < n; i++ {
			v, _ := strconv.ParseInt(fields[1+i], 10, 64)
			p[i] = v
		}
		var costs [5]int64
		for i := 0; i < 5; i++ {
			v, _ := strconv.ParseInt(fields[1+n+i], 10, 64)
			costs[i] = v
		}
		expectCnt, expectPoints := expectedAnswer(n, p, costs)
		var sb strings.Builder
		sb.WriteString(fmt.Sprint(n))
		sb.WriteByte('\n')
		for i, v := range p {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprint(v))
		}
		sb.WriteByte('\n')
		for i, c := range costs {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprint(c))
		}
		sb.WriteByte('\n')
		cmd := exec.Command(bin)
		cmd.Stdin = strings.NewReader(sb.String())
		var out bytes.Buffer
		var stderr bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = &stderr
		err := cmd.Run()
		if err != nil {
			fmt.Printf("test %d: runtime error: %v\nstderr: %s\n", idx, err, stderr.String())
			os.Exit(1)
		}
		outputs := strings.Split(strings.TrimSpace(out.String()), "\n")
		if len(outputs) != 2 {
			fmt.Printf("test %d: expected two output lines got %d\n", idx, len(outputs))
			os.Exit(1)
		}
		gotCnt := strings.TrimSpace(outputs[0])
		gotPts := strings.TrimSpace(outputs[1])
		if gotCnt != expectCnt || gotPts != expectPoints {
			fmt.Printf("test %d failed\nexpected:\n%s\n%s\n\ngot:\n%s\n%s\n", idx, expectCnt, expectPoints, gotCnt, gotPts)
			os.Exit(1)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "scanner error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", idx)
}
