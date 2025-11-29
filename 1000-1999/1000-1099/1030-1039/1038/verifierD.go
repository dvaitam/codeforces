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

// Embedded source for the reference solution (was 1038D.go).
const solutionSource = `package main

import (
	"bufio"
	"fmt"
	"os"
)

var reader *bufio.Reader
var writer *bufio.Writer

func init() {
	reader = bufio.NewReader(os.Stdin)
	writer = bufio.NewWriter(os.Stdout)
}

func readInt() int64 {
	var sign int64 = 1
	var b byte
	var err error
	// skip non-numbers
	for {
		b, err = reader.ReadByte()
		if err != nil {
			return 0
		}
		if b == '-' {
			sign = -1
			b, _ = reader.ReadByte()
			break
		}
		if b >= '0' && b <= '9' {
			break
		}
	}
	var x int64
	for {
		if b < '0' || b > '9' {
			break
		}
		x = x*10 + int64(b-'0')
		b, _ = reader.ReadByte()
	}
	return x * sign
}

func abs(x int64) int64 {
	if x < 0 {
		return -x
	}
	return x
}

func main() {
	defer writer.Flush()
	n64 := readInt()
	n := int(n64)
	a := make([]int64, n)
	var ans int64
	for i := 0; i < n; i++ {
		a[i] = readInt()
		ans += abs(a[i])
	}
	if n == 1 {
		fmt.Fprint(writer, a[0])
		return
	}
	allPos := true
	for i := 0; i < n; i++ {
		if a[i] < 0 {
			allPos = false
			break
		}
	}
	if allPos {
		mn := a[0]
		for i := 1; i < n; i++ {
			if a[i] < mn {
				mn = a[i]
			}
		}
		ans -= 2 * mn
	} else {
		allNeg := true
		for i := 0; i < n; i++ {
			if a[i] > 0 {
				allNeg = false
				break
			}
		}
		if allNeg {
			mn := abs(a[0])
			for i := 1; i < n; i++ {
				ai := abs(a[i])
				if ai < mn {
					mn = ai
				}
			}
			ans -= 2 * mn
		}
	}
	fmt.Fprint(writer, ans)
}
`

const testcasesRaw = `100
4
17 14 -12 3
10
10 20 17 -16 18 -20 10 -4 15 -6
4
10 14 15 10
7
20 -11 -6 20 -11 13 4
1
-16
3
17 -18 -1
1
-3
8
18 4 7 5 16 8 -12 3
2
-18 -12
8
-7 -4 7 20 -1 6 12 4
10
2 14 17 6 17 -6 1 -19 -3 18
3
0 14 16
10
-14 -7 20 16 -3 -2 -13 -16 10 20
8
-15 2 -16 6 -11 -19 -2 7
7
-13 -18 18 19 -18 4 17
6
15 -3 12 -5 -18 -1
1
-16
2
18 14
1
-8
7
-2 19 -4 -11 -18 1 0
6
-12 4 4 9 13 4
10
15 -14 19 12 -3 7 20 -5 -1 7
5
13 -1 15 1 -20
7
17 0 -19 4 19 17 20
3
-17 20 20
6
9 2 2 18 -3 11
1
17
1
-19
6
-4 20 9 -1 17 18
6
-9 3 -9 0 3 18
5
-1 4 -14 -19 16
3
-1 12 -6
5
-5 0 -9 7 -14
2
18 0
6
-6 8 -10 -15 1 -7
10
8 -3 -6 -13 -18 13 -8 0 16 -9
5
1 -15 19 2 17
3
6 -2 13
5
9 2 20 6 -2
7
16 6 -18 6 -11 -8 -20
8
19 12 7 15 -6 -18 9 13
5
14 1 -6 -16 17
5
-13 -5 -18 -18 12
4
7 16 -17 -20
8
-13 -10 12 -1 -5 -19 13 14
7
-17 19 -13 1 -12 -4 14
8
-17 2 -6 -8 -13 14 -13 -10
4
-3 -12 -20 11
10
5 -17 -3 -5 -3 19 13 13 7 -17
8
0 -20 -17 -12 -18 -13 -17 -16
8
-18 -15 12 12 11 0 -10 0
2
2 4
7
17 -1 3 -4 -8 1 7
2
-12 15
1
4
2
16 -9
1
3
8
18 14 4 20 -18 19 7 -17
6
20 11 0 6 6 9
1
-5
4
14 -3 17 -16
7
-6 7 -12 -19 0 3 15
5
-13 9 -13 13 4
2
0 16
9
-14 17 -20 10 -11 -5 4 -18 13
2
16 -14
7
-9 -19 1 -13 -19 -13 10
5
17 -1 -15 -18 16
9
13 -5 -14 15 -14 15 -17 15 0
10
-9 -16 -5 -9 -5 9 19 5 -4 3
10
5 2 15 6 -15 4 12 -5 6 -10
7
16 17 13 10 -11 5 -11
3
-14 11 10
9
8 17 -9 -12 -3 -8 -11 17 12
6
-6 14 -2 6 18 17
10
-3 -7 -1 -19 -3 10 4 -8 -9 16
6
-5 0 10 -11 6 10
10
-7 9 17 15 -19 10 -16 5 -18 9
4
-5 -16 -7 -4
4
-8 -4 -12 -9
10
-18 -4 -10 -18 0 -9 7 -15 -15 -13
2
-4 -2
1
2
8
17 1 -20 -19 1 1 7 4
8
-16 -7 17 11 5 -12 14 0
2
-3 -16
7
-13 8 13 -4 -14 13 3
6
8 -2 -4 -14 1 16
9
13 -13 11 12 2 -17 -2 16 -9
3
-9 3 9
2
-14 15
3
1 18 6
9
-1 -9 9 10 -1 -9 -16 -14 -9
9
14 16 5 2 -14 -3 -3 4 -17
3
-18 10 12
5
-5 12 2 1 5
8
14 -16 2 11 -13 -11 -3 17
2
-13 16
2
-9 -8
10
6 5 -12 17 18 -11 5 -8 14 13`

func abs(x int64) int64 {
	if x < 0 {
		return -x
	}
	return x
}

func expected(arr []int64) string {
	n := len(arr)
	if n == 1 {
		return fmt.Sprintf("%d", arr[0])
	}
	var sum int64
	allPos := true
	allNeg := true
	minPos := int64(1<<63 - 1)
	minAbs := int64(1<<63 - 1)
	for _, v := range arr {
		if v < 0 {
			allPos = false
		} else {
			if v < minPos {
				minPos = v
			}
		}
		if v > 0 {
			allNeg = false
		}
		av := abs(v)
		if av < minAbs {
			minAbs = av
		}
		sum += av
	}
	if allPos {
		sum -= 2 * minPos
	} else if allNeg {
		sum -= 2 * minAbs
	}
	return fmt.Sprintf("%d", sum)
}

func runCase(exe, input, exp string) error {
	var cmd *exec.Cmd
	if strings.HasSuffix(exe, ".go") {
		cmd = exec.Command("go", "run", exe)
	} else {
		cmd = exec.Command(exe)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(out.String())
	if got != strings.TrimSpace(exp) {
		return fmt.Errorf("expected %q got %q", exp, got)
	}
	return nil
}

func main() {
	var _ = solutionSource
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	exe := os.Args[1]
	scan := bufio.NewScanner(strings.NewReader(testcasesRaw))
	scan.Split(bufio.ScanWords)
	if !scan.Scan() {
		fmt.Println("invalid test file")
		os.Exit(1)
	}
	t, _ := strconv.Atoi(scan.Text())
	for caseIdx := 0; caseIdx < t; caseIdx++ {
		scan.Scan()
		n, _ := strconv.Atoi(scan.Text())
		arr := make([]int64, n)
		for i := 0; i < n; i++ {
			scan.Scan()
			v, _ := strconv.Atoi(scan.Text())
			arr[i] = int64(v)
		}
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d\n", n))
		for i := 0; i < n; i++ {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.FormatInt(arr[i], 10))
		}
		sb.WriteByte('\n')
		exp := expected(arr)
		if err := runCase(exe, sb.String(), exp); err != nil {
			fmt.Printf("case %d failed: %v\n", caseIdx+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
