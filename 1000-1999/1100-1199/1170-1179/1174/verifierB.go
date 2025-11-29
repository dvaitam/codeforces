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

// Embedded source for the reference solution (was 1174B.go).
const solutionSource = `package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n int
   if _, err := fmt.Fscan(reader, &n); err != nil {
       return
   }
   a := make([]int, n)
   oddCount := 0
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &a[i])
       if a[i]&1 != 0 {
           oddCount++
       }
   }
   // If both even and odd present, we can sort arbitrarily
   if oddCount > 0 && oddCount < n {
       sort.Ints(a)
   }
   // Output result
   for i, v := range a {
       if i > 0 {
           writer.WriteByte(' ')
       }
       fmt.Fprint(writer, v)
   }
   writer.WriteByte('\n')
}
`

const testcasesRaw = `9 8 8 25 6 10 12 14 22 2
3 20 1 13
2 23 3
3 14 10 18
7 24 30 5 19 14 10 21
6 3 8 15 21 12 21
9 2 13 14 1 14 24 29 11 15
4 12 10 16 3
3 26 4 9
2 18 20
3 26 23 15
7 6 25 14 14 6 8 15
6 30 17 5 12 15 21
2 16 25
4 10 1 27 23
8 20 15 1 7 10 4 25 21
5 18 20 5 14 23
8 3 22 16 25 8 18 25 13
5 21 1 4 9 29
1 1
5 13 17 29 19 23
7 15 4 24 9 12 10 28
4 20 3 2 3
5 10 18 11 4 17
4 29 30 25 6
2 14 28
5 10 17 5 19 17
4 18 4 14 21
9 13 24 25 29 26 9 10 15 12
10 21 5 6 4 23 4 13 13 19 15
3 18 22 10
6 21 16 24 14 7 16
8 23 17 11 16 21 2 15 10
3 24 16 2
10 7 1 12 16 13 29 1 28 17 3
2 22 24
7 1 12 2 4 20 1 9
5 24 29 8 5 25
10 10 7 4 14 15 23 11 13 6 11
7 21 29 22 14 5 15 30
3 17 11 5
4 30 6 15 12
7 14 26 16 13 24 8 26
4 15 7 19 23
1 29
7 2 8 21 3 28 6 12
1 24
3 8 20 10
10 3 23 28 17 25 10 25 29 29 12
7 15 2 21 23 17 22 21
9 30 24 30 14 19 15 16 9 23
8 7 11 9 2 2 2 6 12
1 10
1 5
2 26 14
4 20 13 18 30
4 15 7 11 20
2 20 29
2 26 11
6 18 15 29 11 9 1
9 2 7 12 3 7 28 17 12 7
4 9 22 24 24
5 10 17 28 13 9
8 12 28 23 8 2 10 30 18
2 1 15
8 24 15 2 30 26 14 16 15
8 4 3 3 8 4 27 25 5
7 29 30 7 15 20 3 27
7 18 25 29 27 13 2 6
4 16 8 5 27
5 30 12 11 14 4
9 29 10 20 18 26 7 23 10 25
8 17 20 15 18 21 9 9 8
1 4
10 26 23 4 6 24 14 8 7 10 29
1 24
9 17 14 28 2 4 13 21 9 4
10 12 8 22 23 23 18 22 10 8 24
4 3 17 10 22
6 8 12 21 16 10 19
3 5 26 1
9 17 11 12 19 21 1 26 5 29
7 5 6 17 3 5 25 7
8 19 25 23 7 8 24 5 27
4 25 13 12 20
10 5 21 16 29 30 4 20 27 1 17
10 12 16 15 10 1 8 18 21 6 22
8 26 24 16 18 11 23 28 3
5 5 20 13 23 7
6 26 10 13 2 7 2
6 24 24 8 11 28 15
4 9 12 22 6
5 1 12 19 18 2
3 12 1 16
1 1
4 2 1 8 21
6 3 27 2 12 22 14
3 29 30 7
8 14 5 12 10 6 21 11 24
7 13 1 14 29 9 18 18`

var _ = solutionSource

func expected(nums []int) []int {
	n := len(nums)
	odd := 0
	for _, v := range nums {
		if v&1 != 0 {
			odd++
		}
	}
	if odd > 0 && odd < n {
		tmp := append([]int(nil), nums...)
		sort.Ints(tmp)
		return tmp
	}
	return append([]int(nil), nums...)
}

func run(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	scanner := bufio.NewScanner(strings.NewReader(testcasesRaw))
	idx := 0
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		idx++
		fields := strings.Fields(line)
		n, err := strconv.Atoi(fields[0])
		if err != nil || len(fields) != n+1 {
			fmt.Fprintf(os.Stderr, "case %d: invalid line\n", idx)
			os.Exit(1)
		}
		nums := make([]int, n)
		for i := 0; i < n; i++ {
			nums[i], _ = strconv.Atoi(fields[i+1])
		}
		wantSlice := expected(nums)
		wantStrs := make([]string, len(wantSlice))
		for i, v := range wantSlice {
			wantStrs[i] = strconv.Itoa(v)
		}
		want := strings.Join(wantStrs, " ")

		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d\n", n))
		for i, v := range nums {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.Itoa(v))
		}
		sb.WriteByte('\n')
		input := sb.String()

		got, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "test %d: %v\n", idx, err)
			os.Exit(1)
		}
		if got != want {
			fmt.Fprintf(os.Stderr, "test %d failed\nexpected: %s\ngot: %s\n", idx, want, got)
			os.Exit(1)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "scanner error:", err)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", idx)
}
