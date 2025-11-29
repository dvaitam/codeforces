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

// Embedded source for the reference solution (was 1003A.go).
// Keeping the full text here ensures the verifier is self-contained.
const solutionSource = `package main

import "fmt"

func main() {
	var n int
	if _, err := fmt.Scan(&n); err != nil {
		return
	}
	freq := make(map[int]int)
	var x int
	maxCount := 0
	for i := 0; i < n; i++ {
		fmt.Scan(&x)
		freq[x]++
		if freq[x] > maxCount {
			maxCount = freq[x]
		}
	}
	fmt.Println(maxCount)
}
`

const testcasesRaw = `100
13 4 1 3 5 4 4 3 4 3 5 2 5 2
10 2 1 5 3 5 5 2 3 1 1
11 4 5 1 3 4 3 5 2 5 4 4
17 3 1 5 1 1 4 1 5 4 3 2 3 1 2 5 2 2
5 5 4 1 1 3
17 4 1 3 5 3 1 5 3 5 2 5 5 5 3 4 1 5
13 3 5 2 3 2 2 2 1 5 3 4 1 1
5 2 1 1 5 4
17 3 5 2 2 5 4 5 3 4 4 3 1 3 5 1 4 5
11 2 2 1 3 1 2 3 2 3 4 1
4 2 2 1 5
18 5 1 1 1 2 5 5 1 4 1 3 1 1 5 1 2 2 1
16 2 1 1 5 4 5 1 3 1 2 1 3 3 4 2 1
17 4 1 5 1 4 2 3 3 4 5 2 2 1 2 2 3 5
9 1 5 4 2 1 4 4 5 5
10 3 4 3 2 5 1 4 1 3 1
18 3 2 2 4 3 5 3 3 5 5 2 3 4 4 1 1 5 2
11 2 2 2 4 4 5 4 1 4 5 4
2 2 4
3 3 2 4
17 4 5 5 1 1 4 3 3 4 1 4 2 5 1 2 1 4
14 3 1 2 1 1 5 5 1 2 1 5 2 3 3
6 1 4 4 1 1 3
15 1 3 2 5 3 1 2 3 1 1 1 2 3 5 3
12 5 1 5 4 4 4 3 5 2 2 4 5
10 1 2 2 3 3 3 3 1 3 5
2 1 3
6 2 5 3 3 4 5
5 3 1 4 2 1
10 2 5 1 3 4 3 3 4 1 1
18 4 4 3 3 1 4 1 4 4 1 3 3 2 2 5 4 1 1
3 2 2 1
13 1 1 4 5 5 3 4 4 5 2 4 1 3
8 3 5 2 4 2 3 1 1
1 5
15 2 1 4 4 3 2 1 2 5 2 1 2 4 4 3
18 2 1 5 4 2 5 4 4 5 4 3 4 4 2 5 5 2 1
11 3 3 1 5 2 3 5 2 4 5 3
16 1 1 5 1 1 2 2 1 3 1 4 3 2 2 4 3
17 4 5 5 1 5 1 5 5 1 4 2 3 5 5 4 4 4
20 5 2 1 1 2 3 5 5 3 3 1 4 3 3 4 4 4 1 2 2
8 3 3 1 1 4 4 2 4
20 1 2 3 4 1 5 4 4 4 1 1 4 2 1 1 5 5 2 3 1
18 3 2 4 4 1 1 5 4 5 3 1 5 3 2 4 3 1 5
7 1 4 4 3 2 4 3
3 1 1 4
9 1 5 5 5 2 2 1 5 5
14 5 3 1 2 4 5 4 1 1 4 1 1 4 2
1 4
14 4 1 4 3 3 1 3 1 1 3 1 3 3 2
1 2
12 1 5 2 2 1 2 1 1 3 3 1 5
8 2 2 4 1 4 3 3 2
1 2
12 3 4 3 3 5 3 2 5 1 1 5 5
10 2 4 2 2 2 3 5 2 2 2
10 3 4 1 2 5 1 4 1 1 2
14 3 5 4 2 5 4 3 3 1 2 4 3 5 1
13 4 1 4 3 4 2 3 3 4 1 2 1 3
4 5 5 2 4
13 2 4 4 2 2 4 3 5 2 3 4 1 4
7 3 1 4 5 4 1 2
10 1 3 5 5 2 4 4 1 4 2
18 4 3 1 1 3 1 1 3 4 5 5 4 4 1 3 3 3 2
20 1 1 1 3 3 5 3 1 5 2 2 1 4 3 3 5 2 5 5 2
18 1 4 5 4 3 3 4 3 5 2 2 1 1 4 4 5 4 2
18 3 3 4 4 2 4 4 5 3 4 1 4 3 2 4 1 5 2
1 3
16 4 1 5 1 1 4 1 3 1 1 5 1 3 3 2 2
19 3 2 1 4 4 3 4 2 3 4 4 2 4 2 5 3 2 2 2
15 3 4 4 4 4 2 2 4 2 5 1 4 1 2 1
6 3 1 2 2 5 3
20 1 5 3 3 4 4 1 5 5 4 5 4 4 3 4 2 3 3 1 1
2 2 3
1 3
1 2
3 4 2 5
13 5 2 4 2 3 5 1 5 1 3 3 5 4
11 3 1 5 1 2 3 1 2 5 3 2
7 3 3 3 5 4 3 4
12 2 1 3 5 1 1 4 4 4 1 4 4
15 4 1 1 1 2 1 2 4 2 4 5 1 4 5 4
2 2 2
16 2 2 3 3 3 4 1 5 3 5 5 2 3 4 5 5
15 5 3 3 2 1 1 5 1 2 4 2 2 3 1 5
17 4 1 1 4 3 1 5 3 2 5 3 2 2 1 5 3 3
8 3 4 3 5 2 2 1 5
17 3 3 5 1 2 4 2 2 5 1 2 2 4 5 2 2 2
8 4 3 5 5 2 4 1 5
1 5
20 3 4 4 3 1 2 5 2 4 4 5 3 1 3 2 5 4 2 3 3
13 1 2 1 3 2 3 4 2 3 3 2 3 1
12 5 5 1 2 3 1 4 1 1 2 1 1
8 3 1 1 3 4 2 2 4
14 2 3 3 2 3 4 4 1 4 3 5 5 4 1
19 1 4 4 2 1 5 2 5 5 2 1 3 2 2 2 1 2 5 2
3 4 5 1
20 4 2 5 5 1 3 3 4 1 1 4 1 3 3 2 4 2 5 3 2
13 3 3 4 4 1 3 5 3 5 4 1 5 5
18 3 1 4 4 1 4 3 4 1 1 3 1 3 5 3 2 5 5`

var _ = solutionSource

// Embedded reference logic from the solution source above.
func maxFrequency(arr []int) int {
	freq := make(map[int]int)
	maxCount := 0
	for _, v := range arr {
		freq[v]++
		if freq[v] > maxCount {
			maxCount = freq[v]
		}
	}
	return maxCount
}

func runCase(bin string, arr []int) error {
	var sb strings.Builder
	sb.WriteString(strconv.Itoa(len(arr)))
	sb.WriteByte('\n')
	for i, v := range arr {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(v))
	}
	sb.WriteByte('\n')
	input := sb.String()
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
		return fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	gotStr := strings.TrimSpace(out.String())
	got, err := strconv.Atoi(gotStr)
	if err != nil {
		return fmt.Errorf("invalid output %q", gotStr)
	}
	want := maxFrequency(arr)
	if got != want {
		return fmt.Errorf("expected %d got %d", want, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	scan := bufio.NewScanner(strings.NewReader(testcasesRaw))
	scan.Split(bufio.ScanWords)
	if !scan.Scan() {
		fmt.Println("invalid test file")
		os.Exit(1)
	}
	t, _ := strconv.Atoi(scan.Text())
	for i := 0; i < t; i++ {
		if !scan.Scan() {
			fmt.Println("invalid test file")
			os.Exit(1)
		}
		n, _ := strconv.Atoi(scan.Text())
		arr := make([]int, n)
		for j := 0; j < n; j++ {
			if !scan.Scan() {
				fmt.Println("invalid test file")
				os.Exit(1)
			}
			v, _ := strconv.Atoi(scan.Text())
			arr[j] = v
		}
		if err := runCase(bin, arr); err != nil {
			fmt.Printf("case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
