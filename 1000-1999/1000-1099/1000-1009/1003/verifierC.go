package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

const testcasesRaw = `100
1 1
-4
6 2
5 -1 -1 4 -2 4
1 1
1
7 7
3 0 3 2 3 -1 -5
1 1
2
6 4
1 3 -3 3 -3 -2
4 1
-3 0 -3 -3
9 9
0 3 5 3 -3 2 1 3 0
10 6
0 2 -3 1 2 5 3 -2 2 -1
8 6
5 2 2 0 4 3 2 2
4 3
-3 4 -1 2
5 3
3 3 3 3 5
10 10
1 -1 -2 2 3 0 5 4 -4 0
1 1
-4
1 1
-1
10 4
5 -4 3 -3 -1 -2 -2 -5 1 -5
1 1
0
3 1
5 -5 -4
2 1
-5 -5
1 1
-1
3 1
-3 3 -5
7 5
-5 -2 -3 -5 -5 0 4
2 2
0 2
1 1
2
9 1
-1 1 4 -3 2 -2 -4 5 5
6 1
-5 2 -3 3 4 1
8 6
-3 0 -1 -1 4 1 5 -5
9 3
5 -5 -1 -5 -3 -3 -3 -4 2
4 1
-2 -2 2 -4
5 1
4 -2 4 4 0
5 4
-1 3 -5 -3 -5
7 4
-3 -4 3 -4 -2 -4 -4
1 1
-2
2 1
-5 3
8 8
-1 3 5 1 -2 5 -2 1
7 5
-5 4 4 -5 1 3 4
3 1
5 2 0
1 1
4
6 3
0 -1 -5 5 1 -4
2 2
-2 5
1 1
-5
7 6
2 2 -2 4 4 -4 -5
5 1
0 -1 -4 -2 2
4 1
4 0 1 2
3 2
1 -4 -1
2 1
-4 4
6 6
1 -2 -4 -5 4 5
8 1
2 -1 0 2 -3 0 -1 2
9 8
1 2 5 -1 1 -2 -3 2 4
5 5
1 5 -4 4 4
2 1
0 -3
9 3
1 -4 -4 5 5 -5 -3 -1 1
4 3
2 -3 3 -1
2 1
3 1
2 2
3 -2
9 5
-3 -3 2 -2 1 0 4 -3 2
8 1
4 1 -3 1 3 -5 2 -1
7 3
1 5 2 0 3 0 5
2 1
3 4
4 4
5 1 5 -5
6 4
3 2 5 -3 -4 -5
7 2
4 4 1 -2 -4 1 3
4 3
4 4 -2 2
10 3
-5 4 5 1 2 -1 3 4 -3 2
4 1
0 -5 2 3
2 2
5 0
8 5
3 2 -5 -4 4 0 -3 1
5 2
-5 -3 2 1 2
5 2
-5 -1 3 2 -5
6 1
3 1 4 2 -2 5
5 4
5 -3 2 3 -1
2 2
0 -1
6 6
-1 5 5 1 3 -4
9 4
1 4 3 -3 3 5 -4 -1 -5
4 4
3 -2 3 -1
1 1
-4
7 7
0 -2 0 0 -4 0 2
6 2
2 2 -1 2 -3 2
4 3
0 -3 -4 -2
8 4
5 0 -3 0 -3 -3 -2 -1
9 7
1 0 -1 4 3 4 0 1 5
5 5
4 5 5 -4 0
5 4
2 -3 -1 0 2
8 2
-3 0 1 -3 -5 -4 0 -3
6 1
5 1 -5 3 0 -2
10 7
3 -1 2 5 -3 0 0 -2 2 -4
3 1
0 -1 -3
7 3
-1 -4 0 -2 -2 -2 4
1 1
0
10 1
-3 -3 -4 1 2 -1 -3 0 3 4
2 2
5 4
7 2
-5 1 2 2 4 0 3
10 10
-4 4 3 3 5 2 1 2 -3 1
7 5
2 -5 -4 2 4 -3 -4
9 3
-4 1 -1 2 -5 -1 -4 5 0
4 2
-5 -3 1 5
2 2
5 2
1 1
-2
2 2
-3 3
1 1
3`

// Reference implementation embedded from 1003C.go to avoid external oracle.
func referenceAverage(n, k int, a []int) float64 {
	sum := make([]int, n)
	if n > 0 {
		sum[0] = a[0]
		for i := 1; i < n; i++ {
			sum[i] = sum[i-1] + a[i]
		}
	}
	var avg float64
	for L := k; L <= n; L++ {
		mx := sum[L-1]
		for j := L; j < n; j++ {
			s := sum[j] - sum[j-L]
			if s > mx {
				mx = s
			}
		}
		cur := float64(mx) / float64(L)
		if cur > avg {
			avg = cur
		}
	}
	return avg
}

func runCase(bin string, n, k int, arr []int) error {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", n, k))
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
	got, err := strconv.ParseFloat(gotStr, 64)
	if err != nil {
		return fmt.Errorf("invalid output %q", gotStr)
	}
	want := referenceAverage(n, k, arr)
	if math.Abs(got-want) > 1e-6 {
		return fmt.Errorf("expected %.6f got %.6f", want, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierC.go /path/to/binary")
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
		if !scan.Scan() {
			fmt.Println("invalid test file")
			os.Exit(1)
		}
		k, _ := strconv.Atoi(scan.Text())
		arr := make([]int, n)
		for j := 0; j < n; j++ {
			if !scan.Scan() {
				fmt.Println("invalid test file")
				os.Exit(1)
			}
			v, _ := strconv.Atoi(scan.Text())
			arr[j] = v
		}
		if err := runCase(bin, n, k, arr); err != nil {
			fmt.Printf("case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
