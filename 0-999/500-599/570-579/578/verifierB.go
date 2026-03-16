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

func solveB(n, k int, x int64, arr []int64) int64 {
	prefix := make([]int64, n+1)
	for i := 0; i < n; i++ {
		prefix[i+1] = prefix[i] | arr[i]
	}
	suffix := make([]int64, n+1)
	for i := n - 1; i >= 0; i-- {
		suffix[i] = suffix[i+1] | arr[i]
	}
	pow := int64(1)
	for i := 0; i < k; i++ {
		pow *= x
	}
	best := int64(0)
	for i := 0; i < n; i++ {
		val := prefix[i] | (arr[i] * pow) | suffix[i+1]
		if val > best {
			best = val
		}
	}
	return best
}

func run(bin string, input string) (string, error) {
	cmd := exec.Command(bin)
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

const testcasesBRaw = `1 1 2
23
3 3 4
38 13 38
1 5 3
27
7 5 4
34 28 32 17 2 1 23
8 3 5
27 33 10 35 11 15 14 1
3 3 3
8 32 32
6 5 3
28 50 26 47 33 48
6 5 4
23 28 10 48 25 45
8 5 3
31 17 31 32 32 50 22 42
8 4 4
36 46 35 46 29 31 42 14
6 2 4
49 30 19 19 45 32
7 3 3
31 32 23 43 39 4 50
6 1 3
47 6 3 36 41 3
5 5 3
43 6 48 33 8
5 2 3
3 27 45 48 2
1 3 4
11
4 1 2
7 4 1 2
1 3 4
8
3 2 2
24 37 2
4 2 2
0 22 39 40
2 3 4
31 1
5 4 2
16 48 25 39 45
3 4 3
5 42 43
6 1 2
28 50 8 33 37 49
7 4 4
9 21 16 16 38 26 41
1 5 3
42
1 3 2
8
3 2 2
29 40 14
1 2 3
45
8 1 4
5 37 14 39 50 39 45 23
5 4 4
33 48 0 9 2
7 4 3
7 32 46 5 15 6 6
1 2 3
6
4 1 5
29 19 34 41
7 2 3
46 27 27 32 1 37 37
1 4 3
6
8 3 2
33 7 39 23 18 44 23 19
1 4 2
6
5 2 2
28 3 26 40 31
8 2 2
0 18 1 23 19 46 4 14
8 2 2
36 23 25 45 29 8 48 22
7 1 4
7 7 5 39 21 41 25
4 1 2
39 42 30 49
1 4 4
22
8 2 4
17 30 33 30 46 46 26 31
5 4 3
10 31 38 16 35
7 1 2
4 22 11 34 9 26 4
2 1 3
18 24
4 3 5
11 33 18 7
3 5 5
6 21 33
4 5 4
10 10 29 45
4 4 4
50 48 36 46
3 4 5
46 1 38
7 2 5
32 3 30 17 25 16 45
7 4 4
35 21 45 47 42 5 48
4 5 3
25 42 24 40
1 3 5
33
8 2 2
1 25 13 46 36 38 24 13
2 4 3
17 47
4 4 3
0 39 43 27
8 3 3
29 45 13 48 4 22 0 31
2 5 5
43 21
8 3 5
1 5 39 48 22 11 48 48
7 3 3
3 10 31 24 29 43 18
3 1 4
35 29 0
6 1 5
36 28 13 43 19 31
3 4 4
4 16 20
5 3 4
41 41 25 33 5
4 4 3
32 40 5 19
1 2 5
35
4 5 4
3 7 7 43
7 3 3
20 22 4 21 29 23 10
8 4 4
29 8 45 28 40 13 17 20
3 1 3
30 12 48
6 2 4
8 50 8 14 17 35
7 4 4
17 46 38 32 37 44 46
6 4 4
34 39 40 42 4 23
5 4 5
11 16 22 28 30
2 2 4
24 8
1 1 4
10
6 1 5
0 34 20 15 38 24
5 4 3
23 20 12 31 6
3 2 4
16 9 26
6 3 2
21 12 15 45 15 46
1 3 4
41
1 2 3
4
7 4 4
8 20 33 36 7 21 41
7 2 2
25 49 30 31 39 20 34
2 5 5
25 44
8 2 5
24 33 28 2 6 28 37 8
2 5 3
4 25
5 4 2
16 6 42 22 14
3 1 3
27 42 5
6 4 2
30 15 4 30 8 35
1 2 2
3
4 5 2
33 21 43 33
4 2 4
31 0 8 34
2 2 2
29 13
1 5 3
40
7 3 5
45 33 32 49 43 10 32
`

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	scanner := bufio.NewScanner(strings.NewReader(testcasesBRaw))
	idx := 0
	for {
		if !scanner.Scan() {
			break
		}
		line1 := strings.TrimSpace(scanner.Text())
		if line1 == "" {
			continue
		}
		idx++
		parts := strings.Fields(line1)
		if len(parts) != 3 {
			fmt.Fprintf(os.Stderr, "case %d bad header\n", idx)
			os.Exit(1)
		}
		n, _ := strconv.Atoi(parts[0])
		k, _ := strconv.Atoi(parts[1])
		x64, _ := strconv.ParseInt(parts[2], 10, 64)
		if !scanner.Scan() {
			fmt.Fprintf(os.Stderr, "case %d missing array line\n", idx)
			os.Exit(1)
		}
		line2 := strings.TrimSpace(scanner.Text())
		nums := strings.Fields(line2)
		if len(nums) != n {
			fmt.Fprintf(os.Stderr, "case %d expected %d numbers got %d\n", idx, n, len(nums))
			os.Exit(1)
		}
		arr := make([]int64, n)
		for i, s := range nums {
			v, _ := strconv.ParseInt(s, 10, 64)
			arr[i] = v
		}
		expected := solveB(n, k, x64, arr)
		var input strings.Builder
		fmt.Fprintf(&input, "%d %d %d\n", n, k, x64)
		for i, v := range arr {
			if i > 0 {
				input.WriteByte(' ')
			}
			fmt.Fprintf(&input, "%d", v)
		}
		input.WriteByte('\n')
		got, err := run(bin, input.String())
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", idx, err)
			os.Exit(1)
		}
		val, err := strconv.ParseInt(strings.TrimSpace(got), 10, 64)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: non-integer output %s\n", idx, got)
			os.Exit(1)
		}
		if val != expected {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %d got %d\n", idx, expected, val)
			os.Exit(1)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "scanner error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", idx)
}
