package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

type testCase struct {
	n int
	a int
	b int
}

const solution1007CSource = `package main

import (
	"bufio"
	"fmt"
	"os"
)

var (
	reader = bufio.NewReader(os.Stdin)
	writer = bufio.NewWriter(os.Stdout)
)

// query sends a guess (x, y) and returns the interactor's response:
// 1 means x < a, 2 means y < b, 3 means x > a or y > b.
func query(x, y int64) int {
	fmt.Fprintf(writer, "? %d %d\n", x, y)
	writer.Flush()
	var res int
	if _, err := fmt.Fscan(reader, &res); err != nil {
		os.Exit(0)
	}
	return res
}

func main() {
	defer writer.Flush()
	var n int64
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return
	}
	// binary search for a using y = n
	var la, ra int64 = 1, n
	for la < ra {
		mid := (la + ra) / 2
		res := query(mid, n)
		if res == 1 {
			// x < a
			la = mid + 1
		} else {
			// res == 3: x >= a (or y > b)
			ra = mid
		}
	}
	a := la
	// binary search for b using x = a
	var lb, rb int64 = 1, n
	for lb < rb {
		mid := (lb + rb) / 2
		res := query(a, mid)
		if res == 2 {
			// y < b
			lb = mid + 1
		} else {
			// res == 3: y >= b (or x > a)
			rb = mid
		}
	}
	b := lb
	// output the answer
	fmt.Fprintf(writer, "! %d %d\n", a, b)
}
`

// Keep the embedded reference solution reachable so it is preserved in the binary.
var _ = solution1007CSource

var testcases = []testCase{
	{n: 9, a: 2, b: 2},
	{n: 48, a: 11, b: 48},
	{n: 87, a: 40, b: 33},
	{n: 79, a: 28, b: 78},
	{n: 6, a: 5, b: 6},
	{n: 22, a: 14, b: 21},
	{n: 52, a: 52, b: 47},
	{n: 67, a: 48, b: 57},
	{n: 66, a: 35, b: 5},
	{n: 5, a: 3, b: 4},
	{n: 42, a: 25, b: 28},
	{n: 69, a: 22, b: 23},
	{n: 32, a: 15, b: 2},
	{n: 24, a: 11, b: 6},
	{n: 19, a: 17, b: 17},
	{n: 48, a: 33, b: 44},
	{n: 73, a: 24, b: 58},
	{n: 55, a: 48, b: 34},
	{n: 99, a: 47, b: 76},
	{n: 47, a: 24, b: 29},
	{n: 22, a: 13, b: 15},
	{n: 85, a: 68, b: 32},
	{n: 64, a: 36, b: 64},
	{n: 66, a: 66, b: 46},
	{n: 86, a: 59, b: 60},
	{n: 46, a: 37, b: 36},
	{n: 94, a: 59, b: 63},
	{n: 86, a: 29, b: 42},
	{n: 91, a: 22, b: 79},
	{n: 36, a: 31, b: 20},
	{n: 40, a: 33, b: 36},
	{n: 68, a: 65, b: 53},
	{n: 41, a: 14, b: 32},
	{n: 67, a: 47, b: 10},
	{n: 45, a: 1, b: 13},
	{n: 97, a: 14, b: 8},
	{n: 75, a: 7, b: 35},
	{n: 77, a: 30, b: 14},
	{n: 98, a: 67, b: 18},
	{n: 36, a: 16, b: 14},
	{n: 9, a: 7, b: 1},
	{n: 9, a: 6, b: 6},
	{n: 24, a: 8, b: 22},
	{n: 5, a: 1, b: 1},
	{n: 10, a: 1, b: 1},
	{n: 95, a: 3, b: 48},
	{n: 34, a: 9, b: 11},
	{n: 96, a: 24, b: 67},
	{n: 90, a: 1, b: 50},
	{n: 77, a: 6, b: 32},
	{n: 21, a: 2, b: 1},
	{n: 46, a: 40, b: 41},
	{n: 97, a: 96, b: 15},
	{n: 38, a: 22, b: 32},
	{n: 5, a: 3, b: 4},
	{n: 72, a: 6, b: 34},
	{n: 98, a: 52, b: 80},
	{n: 92, a: 20, b: 61},
	{n: 30, a: 3, b: 22},
	{n: 89, a: 41, b: 14},
	{n: 5, a: 4, b: 2},
	{n: 68, a: 51, b: 63},
	{n: 67, a: 42, b: 19},
	{n: 45, a: 17, b: 17},
	{n: 79, a: 54, b: 3},
	{n: 91, a: 72, b: 18},
	{n: 87, a: 8, b: 33},
	{n: 6, a: 2, b: 2},
	{n: 23, a: 4, b: 15},
	{n: 83, a: 30, b: 66},
	{n: 92, a: 5, b: 32},
	{n: 31, a: 23, b: 15},
	{n: 11, a: 5, b: 2},
	{n: 77, a: 30, b: 47},
	{n: 34, a: 28, b: 18},
	{n: 69, a: 1, b: 20},
	{n: 6, a: 4, b: 4},
	{n: 22, a: 4, b: 17},
	{n: 94, a: 12, b: 31},
	{n: 15, a: 2, b: 1},
	{n: 25, a: 25, b: 8},
	{n: 15, a: 4, b: 1},
	{n: 68, a: 60, b: 59},
	{n: 41, a: 35, b: 25},
	{n: 29, a: 22, b: 25},
	{n: 28, a: 24, b: 26},
	{n: 57, a: 28, b: 33},
	{n: 4, a: 1, b: 4},
	{n: 69, a: 24, b: 13},
	{n: 86, a: 62, b: 47},
	{n: 4, a: 1, b: 3},
	{n: 39, a: 24, b: 20},
	{n: 4, a: 4, b: 1},
	{n: 15, a: 5, b: 4},
	{n: 88, a: 3, b: 58},
	{n: 9, a: 7, b: 8},
	{n: 61, a: 14, b: 57},
	{n: 77, a: 10, b: 1},
	{n: 38, a: 2, b: 24},
	{n: 41, a: 5, b: 15},
}

func runCase(bin string, n, a, b int) error {
	cmd := exec.Command(bin)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return err
	}
	cmd.Stderr = os.Stderr
	if err := cmd.Start(); err != nil {
		return err
	}

	reader := bufio.NewReader(stdout)
	writer := bufio.NewWriter(stdin)

	fmt.Fprintln(writer, n)
	writer.Flush()

	queries := 0
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			cmd.Process.Kill()
			return fmt.Errorf("read error: %v", err)
		}
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "?") {
			queries++
			if queries > 600 {
				cmd.Process.Kill()
				return fmt.Errorf("too many queries")
			}
			parts := strings.Fields(line)
			if len(parts) != 3 {
				cmd.Process.Kill()
				return fmt.Errorf("invalid query: %s", line)
			}
			x, _ := strconv.Atoi(parts[1])
			y, _ := strconv.Atoi(parts[2])
			var res int
			if x < a {
				res = 1
			} else if y < b {
				res = 2
			} else if x > a || y > b {
				res = 3
			} else {
				res = 3
			}
			fmt.Fprintln(writer, res)
			writer.Flush()
		} else if strings.HasPrefix(line, "!") {
			parts := strings.Fields(line)
			if len(parts) != 3 {
				cmd.Process.Kill()
				return fmt.Errorf("invalid answer: %s", line)
			}
			x, _ := strconv.Atoi(parts[1])
			y, _ := strconv.Atoi(parts[2])
			if x == a && y == b {
				stdin.Close()
				cmd.Wait()
				return nil
			}
			cmd.Process.Kill()
			return fmt.Errorf("wrong answer %d %d expected %d %d", x, y, a, b)
		} else {
			cmd.Process.Kill()
			return fmt.Errorf("invalid output: %s", line)
		}
	}
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	for i, tc := range testcases {
		if err := runCase(bin, tc.n, tc.a, tc.b); err != nil {
			fmt.Printf("case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(testcases))
}
