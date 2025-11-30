package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"
)

const testcases = `18 31 3 21 2 21 26 9 49 41 18 27 43 10 39 10 26 20 33
2 11 9
6 31 46 42 46 49 47
2 47 34
2 36 45
14 12 23 38 6 6 36 12 17 13 17 21 45 46 17
10 34 30 10 49 29 36 10 3 41 38
6 42 33 3 49 21 5
8 42 30 40 16 30 34 11 46
12 42 9 31 50 36 4 35 6 34 22 1 50
4 7 28 39 23
16 22 25 33 24 41 8 9 21 2 12 47 9 2 22 39 13
2 27 42
2 45 20
14 4 39 50 46 11 23 5 27 4 29 23 39 40 49
10 44 20 37 30 27 12 2 30 17 13
14 5 23 7 8 2 23 2 12 26 40 45 42 1 21
16 36 48 45 32 31 6 4 35 26 17 2 42 34 7 6 22
12 7 31 3 10 34 41 19 3 1 25 22 11
18 45 10 11 12 50 11 41 44 16 40 22 2 31 41 44 26 3 15
8 41 19 22 11 16 23 15 11
14 30 24 37 9 25 37 50 1 11 38 1 44 25 46
6 10 2 2 21 33 1
2 4 50
4 37 40 10 50
6 44 25 2 27 28 37
12 46 16 9 24 33 14 35 26 5 9 27 37
12 7 28 28 16 31 25 15 26 16 42 31 26
4 17 18 34 24
18 2 39 40 50 31 16 18 3 40 21 26 41 7 35 4 10 46 26
2 50 27
14 28 7 46 30 39 30 11 11 22 31 27 11 38 19
18 8 24 23 10 41 23 50 31 41 34 49 4 13 17 12 47 38 21
10 25 41 3 19 36 28 3 44 27 18
14 47 14 23 9 9 8 40 23 11 2 28 37 26 30
4 42 46 46 44
4 50 28 35 47
18 9 11 10 14 11 15 2 34 9 32 23 40 48 19 46 22 44 8
14 50 17 11 22 41 33 22 35 45 10 25 48 48 36
10 16 25 23 25 31 33 20 27 27 7
6 10 1 38 39 41 46
18 7 45 50 42 45 14 39 42 39 33 8 18 45 47 40 11 25 6
2 1 8
12 47 31 21 7 44 29 24 38 45 17 43 32
8 11 36 36 5 50 33 11 2
6 34 27 39 43 14 29
14 17 2 38 9 25 11 29 37 4 25 6 42 38 26
12 15 33 30 3 31 40 7 18 33 32 36 42
14 31 17 12 15 35 24 11 20 39 10 30 5 5 31
14 37 27 36 6 17 31 15 8 19 10 24 7 9 43
2 43 9
18 13 1 3 26 36 45 49 39 32 46 45 7 31 36 23 50 22 7
2 16 15
16 20 18 15 1 32 23 33 22 6 5 20 37 28 15 48 24
14 49 10 15 19 13 48 50 31 43 23 19 25 40 9
4 26 23 33 31
8 42 46 24 41 23 28 18 23
14 50 50 46 19 7 31 19 8 29 10 23 16 48 12
12 32 15 8 44 25 25 30 33 30 37 40 15
14 33 20 32 15 21 34 44 1 6 31 21 26 15 28
2 37 50
2 27 6
10 13 47 21 12 8 12 45 24 2 15
2 1 25
18 1 9 8 39 39 13 6 48 30 13 1 34 40 27 5 35 12 15
8 27 25 31 1 28 14 25 49
2 40 18
2 38 23
12 22 44 30 42 9 39 34 6 17 7 46 7
10 2 45 10 40 50 43 9 25 14 37
12 13 27 33 33 8 36 7 45 31 8 33 29
16 12 30 36 22 9 27 17 25 6 37 33 22 15 30 16 23
16 50 27 2 29 47 1 36 26 49 50 29 15 28 16 17 31
16 10 15 29 19 24 42 32 39 10 34 44 6 13 20 34 41
4 4 10 22 3
12 40 11 50 30 25 44 13 27 32 13 11 26
2 50 22
18 13 47 21 12 34 35 46 39 38 28 10 41 45 40 32 43 37 14
16 3 16 32 39 21 36 47 13 1 3 4 9 16 29 44 40
8 43 44 7 49 34 28 42 41
10 44 22 32 44 13 11 24 36 23 46
16 26 29 43 21 30 5 10 15 49 8 49 10 26 30 46 35
8 24 46 3 2 26 14 6 27
14 35 14 41 1 43 9 34 42 32 20 25 35 35 5
4 39 47 32 49
2 14 27
12 20 6 48 2 50 16 40 2 38 33 24 10
16 17 8 38 29 18 15 25 25 48 33 22 15 38 17 6 3
6 21 41 38 37 32 5
18 31 28 21 25 47 2 8 28 39 43 9 49 10 2 8 46 28 16
2 41 18
10 25 8 23 19 24 40 41 13 48 35
18 4 22 33 36 1 11 32 2 41 17 22 9 23 7 40 36 27 49
4 17 37 7 34
4 38 9 32 3
18 15 23 7 36 46 14 12 45 43 15 26 13 39 45 48 29 23 11
10 47 23 36 13 14 30 21 9 18 36
2 6 43
12 20 46 21 3 6 40 1 19 8 11 46 18
8 33 36 23 29 6 27 33 35`

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
	sb1.WriteString(strconv.Itoa(len(inc)))
	for _, v := range inc {
		sb1.WriteByte(' ')
		sb1.WriteString(strconv.Itoa(v))
	}
	var sb2 strings.Builder
	sb2.WriteString(strconv.Itoa(len(dec)))
	for _, v := range dec {
		sb2.WriteByte(' ')
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
	lines := strings.Split(testcases, "\n")
	count := 0
	for idx, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		count++
		fields := strings.Fields(line)
		n, _ := strconv.Atoi(fields[0])
		arr := make([]int, n)
		for i := 0; i < n; i++ {
			arr[i], _ = strconv.Atoi(fields[i+1])
		}
		status, expectedSeq, _ := expectedOutput(arr)

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
			fmt.Printf("Test %d failed: expected status %s got %s\n", idx+1, status, outLines[0])
			os.Exit(1)
		}
		if status == "NO" {
			continue
		}
		expLines := strings.Split(expectedSeq, "\n")
		if len(outLines) < 3 {
			fmt.Printf("Test %d failed: expected 3 lines output\n", idx+1)
			os.Exit(1)
		}
		if strings.TrimSpace(outLines[1]) != expLines[0] || strings.TrimSpace(outLines[2]) != expLines[1] {
			fmt.Printf("Test %d failed: expected sequences\n%s\n%s\n got \n%s\n%s\n", idx+1, expLines[0], expLines[1], strings.TrimSpace(outLines[1]), strings.TrimSpace(outLines[2]))
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", count)
}
