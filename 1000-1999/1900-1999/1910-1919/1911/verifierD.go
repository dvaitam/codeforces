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

func expectedOutput(arr []int) (string, string, bool) {
	count := make(map[int]int)
	for _, v := range arr {
		count[v]++
		if count[v] > 2 {
			return "NO", "", false
		}
	}
	inc := []int{}
	dec := []int{}
	for v, c := range count {
		inc = append(inc, v)
		if c == 2 {
			dec = append(dec, v)
		}
	}
	sort.Ints(inc)
	sort.Sort(sort.Reverse(sort.IntSlice(dec)))
	var sb1 strings.Builder
	for i, v := range inc {
		if i > 0 {
			sb1.WriteByte(' ')
		}
		sb1.WriteString(strconv.Itoa(v))
	}
	var sb2 strings.Builder
	for i, v := range dec {
		if i > 0 {
			sb2.WriteByte(' ')
		}
		sb2.WriteString(strconv.Itoa(v))
	}
	return "YES", sb1.String() + "\n" + sb2.String(), true
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	f, err := os.Open("testcasesD.txt")
	if err != nil {
		panic(err)
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
		fields := strings.Fields(line)
		n, _ := strconv.Atoi(fields[0])
		arr := make([]int, n)
		for i := 0; i < n; i++ {
			arr[i], _ = strconv.Atoi(fields[i+1])
		}
		status, expectedSeq, ok := expectedOutput(arr)

		var input strings.Builder
		input.WriteString(fmt.Sprintf("%d\n", n))
		for i, v := range arr {
			if i > 0 {
				input.WriteByte(' ')
			}
			input.WriteString(strconv.Itoa(v))
		}
		input.WriteByte('\n')

		cmd := exec.Command(bin)
		cmd.Stdin = strings.NewReader(input.String())
		var out bytes.Buffer
		var errBuf bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = &errBuf
		if err := cmd.Run(); err != nil {
			fmt.Printf("Test %d: runtime error: %v\nstderr: %s\n", idx, err, errBuf.String())
			os.Exit(1)
		}
		outLines := strings.Split(strings.TrimSpace(out.String()), "\n")
		if outLines[0] != status {
			fmt.Printf("Test %d failed: expected status %s got %s\n", idx, status, outLines[0])
			os.Exit(1)
		}
		if status == "NO" {
			continue
		}
		expLines := strings.Split(expectedSeq, "\n")
		if len(outLines) < 3 {
			fmt.Printf("Test %d failed: expected 3 lines output\n", idx)
			os.Exit(1)
		}
		if strings.TrimSpace(outLines[1]) != expLines[0] || strings.TrimSpace(outLines[2]) != expLines[1] {
			fmt.Printf("Test %d failed: expected sequences\n%s\n%s\n got \n%s\n%s\n", idx, expLines[0], expLines[1], strings.TrimSpace(outLines[1]), strings.TrimSpace(outLines[2]))
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", idx)
}
