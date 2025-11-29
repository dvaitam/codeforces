package main

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

const solution1375BSource = `package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var T int
   fmt.Fscan(reader, &T)
   for T > 0 {
       T--
       var n, m int
       fmt.Fscan(reader, &n, &m)
       a := make([][]int, n)
       for i := 0; i < n; i++ {
           a[i] = make([]int, m)
           for j := 0; j < m; j++ {
               fmt.Fscan(reader, &a[i][j])
           }
       }
       valid := true
       b := make([][]int, n)
       for i := 0; i < n; i++ {
           b[i] = make([]int, m)
           for j := 0; j < m; j++ {
               cnt := 0
               if i > 0 {
                   cnt++
               }
               if i < n-1 {
                   cnt++
               }
               if j > 0 {
                   cnt++
               }
               if j < m-1 {
                   cnt++
               }
               b[i][j] = cnt
               if a[i][j] > b[i][j] {
                   valid = false
               }
           }
       }
       if !valid {
           fmt.Fprintln(writer, "NO")
           continue
       }
       fmt.Fprintln(writer, "YES")
       for i := 0; i < n; i++ {
           for j := 0; j < m; j++ {
               fmt.Fprint(writer, b[i][j])
               if j+1 < m {
                   fmt.Fprint(writer, " ")
               }
           }
           fmt.Fprintln(writer)
       }
   }
}
`

// Keep the embedded reference solution reachable.
var _ = solution1375BSource

type testCase struct {
	n    int
	m    int
	grid [][]int
}

const testcasesRaw = `2 1 2 0
4 4 3 3 1 0 3 0 3 3 0 3 2 1 0 2 0 0
1 1 3
2 4 0 1 3 3 1 2 1 1
4 3 0 3 0 1 2 0 2 3 1 2 2 3
4 1 3 1 3 3
2 3 2 0 3 0 1 3
3 4 0 3 0 2 3 1 1 1 0 1 1 3
3 3 3 2 0 3 1 1 3 0 3
3 2 3 3 2 3 2 0
3 4 0 1 1 1 0 2 0 0 0 0 3 0
3 2 2 0 1 2 2 0
2 2 2 1 2 2
4 3 3 3 0 0 2 3 2 3 1 2 0 2
2 4 0 1 0 3 1 0 1 3
4 2 3 1 0 3 2 3 0 2
2 2 0 2 0 0
3 3 1 3 2 1 0 0 1 3 1
1 4 1 2 0 1
4 2 3 0 3 2 3 0 2 3
3 1 1 1 2
2 3 3 1 2 0 3 2
4 2 0 0 0 1 1 1 1 2
3 3 2 2 2 0 2 1 3 1 0
3 1 3 0 3
2 2 2 0 3 0
2 1 2 2
3 1 3 2 0
1 3 0 0 0
4 1 0 1 1 3
2 1 3 1
2 2 0 3 3 2
3 4 2 0 1 2 0 0 0 2 2 3 3 2
4 1 0 2 3 0
3 2 3 2 2 1 1 2
2 2 2 0 2 0
4 1 2 1 3 2
1 3 1 2 2
2 3 0 0 1 1 0 1
4 1 2 0 0 0
3 3 1 1 2 0 2 3 1 0
4 3 0 3 3 1 2 1 2 1 3 2 3 2
4 1 2 1 0 1
3 1 1 1 3
4 1 2 3 0 2
2 2 3 0 2 0
3 4 3 2 3 3 0 0 0 0 0 3 3 0
1 3 0 3 1
2 2 2 3 3 1
3 3 2 2 2 3 0 0 3 0 1
1 3 2 3 2
3 1 1 3 1
3 2 1 3 3 2 2 3
3 2 0 0 1 0 1 0
3 3 2 2 0 3 3 0 2 3 3
2 2 3 0 0 0
3 4 1 3 3 3 0 0 3 1 3 3 1 1
1 4 0 2 2 3
4 2 0 2 1 1 1 1 0 2
4 1 3 2 1 1
4 2 3 0 0 3 3 2 0 2
3 2 0 2 3 3 0 0
3 4 2 1 2 0 3 3 0 2 1 3 2 2
1 3 3 1 0
1 4 2 2 1
3 2 0 3 2 1 2 1
4 4 0 0 2 1 1 3 0 1 0 0 1 2 0 3 3 2 3 2
3 2 0 2 3 1 3 0
4 2 3 1 1 3 0 3 0 2
3 4 2 1 2 0 0 3 2 1 3 2 1
2 1 1 0
3 4 2 0 2 2 2 1 0 2 3 0 2 2
4 1 3 3 3 2
1 3 0 2 2
3 3 1 2 3 0 2 1 2 3 2
4 4 1 3 1 1 1 2 2 0 2 0 3 0 3 2 0 2 3
3 4 3 0 2 2 0 0 0 3 2 2 3 1
1 1 1
2 4 1 1 0 3 3 2 3 3
1 1 2
4 4 1 3 3 1 2 2 2 2 0 2 2 3 0 1 0 3 3 3
2 3 3 2 3 0 1 2
2 1 2 1
2 1 2 0
3 4 2 2 2 1 3 3 1 3 2 2 0 2
3 3 1 3 3 0 1 1 2 2 0
3 4 2 1 2 1 0 3 0 0 3 3 0 3
3 3 0 2 0 1 0 1 3 1 1
2 4 0 3 0 1 2 0 0 0
2 2 0 2 1 2
1 4 1 3 2 2
2 2 3 0 1 0
4 4 3 3 2 1 3 3 0 1 0 0 3 3 2 2 0 2 1 2 0
1 1 3
2 1 2 2
2 1 3 3
2 3 2 2 0 1 0 3
3 4 2 2 3 3 1 0 3 2 3 3 3 3
4 3 3 3 0 1 2 0 2 1 3 0 2 2 3 0 0
3 4 3 1 0 0 1 2 1 0 3 3 3 0
3 4 3 3 3 1 1 2 2 3 3 3 1 0
1 1 3
3 1 3 3 1
4 1 2 1 1 2
4 4 3 3 3 0 0 0 2 1 2 3 1 0 2 3 0 2 3 0
4 3 3 0 3 1 3 0 3 2 2 3 3 2 1 2
4 4 0 2 2 3 0 3 0 1 3 0 1 0 3 3 0 3 2 2 1
3 1 2 1 3
2 1 2 3
1 4 0 3 0 3
1 4 3 1 3
2 1 3 0
4 1 1 2 3
2 2 2 0 1 3
3 2 2 0 0 3 3 0
3 1 3 1 3
2 3 0 0 1 3 3 3
2 1 0 1
2 4 2 2 2 1 3 0 1 2
2 3 3 1 3 1 0 1
4 1 1 3 3 1
4 3 3 0 0 2 2 3 3 3 3 3
3 1 3 0 2
2 3 2 2 3 3 1 2
4 3 0 3 1 2 0 0 3 3 1 1 0 0 0
3 1 1 2 0
1 3 2 1 3
4 4 1 1 1 0 1 0 3 1 0 1 3 0 2 2 2 0 0
3 1 3 1 1
2 2 0 3 2 3
4 2 2 0 2 2 1 1 0 0
2 3 3 2 0 3 1
3 3 0 1 3 0 3 3 3 3 2
3 2 2 0 2 1 1 3
1 1 3
`

func parseTestcases() []testCase {
	fields := strings.Fields(testcasesRaw)
	var res []testCase
	for i := 0; i < len(fields); {
		if i+2 >= len(fields) {
			break
		}
		n, _ := strconv.Atoi(fields[i])
		m, _ := strconv.Atoi(fields[i+1])
		i += 2
		if i+n*m > len(fields) {
			break
		}
		grid := make([][]int, n)
		for r := 0; r < n; r++ {
			grid[r] = make([]int, m)
			for c := 0; c < m; c++ {
				val, _ := strconv.Atoi(fields[i])
				i++
				grid[r][c] = val
			}
		}
		res = append(res, testCase{n: n, m: m, grid: grid})
	}
	return res
}

func computeExpected(tc testCase) (bool, [][]int) {
	n, m := tc.n, tc.m
	b := make([][]int, n)
	valid := true
	for i := 0; i < n; i++ {
		b[i] = make([]int, m)
		for j := 0; j < m; j++ {
			cnt := 0
			if i > 0 {
				cnt++
			}
			if i < n-1 {
				cnt++
			}
			if j > 0 {
				cnt++
			}
			if j < m-1 {
				cnt++
			}
			b[i][j] = cnt
			if tc.grid[i][j] > cnt {
				valid = false
			}
		}
	}
	return valid, b
}

func buildInput(tc testCase) string {
	var sb strings.Builder
	sb.WriteString("1\n")
	sb.WriteString(fmt.Sprintf("%d %d\n", tc.n, tc.m))
	for i := 0; i < tc.n; i++ {
		for j := 0; j < tc.m; j++ {
			if j > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.Itoa(tc.grid[i][j]))
		}
		sb.WriteByte('\n')
	}
	return sb.String()
}

func runCase(bin string, idx int, tc testCase) error {
	valid, expectGrid := computeExpected(tc)
	input := buildInput(tc)
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("case %d failed: %v\nstderr: %s", idx, err, string(out))
	}
	tokens := strings.Fields(strings.TrimSpace(string(out)))
	if len(tokens) == 0 {
		return fmt.Errorf("case %d failed: empty output", idx)
	}
	if !valid {
		if strings.ToUpper(tokens[0]) != "NO" {
			return fmt.Errorf("case %d failed: expected NO got %s", idx, tokens[0])
		}
		return nil
	}
	if strings.ToUpper(tokens[0]) != "YES" {
		return fmt.Errorf("case %d failed: expected YES got %s", idx, tokens[0])
	}
	if len(tokens)-1 != tc.n*tc.m {
		return fmt.Errorf("case %d failed: expected %d numbers got %d", idx, tc.n*tc.m, len(tokens)-1)
	}
	p := 1
	for i := 0; i < tc.n; i++ {
		for j := 0; j < tc.m; j++ {
			val, err := strconv.Atoi(tokens[p])
			if err != nil {
				return fmt.Errorf("case %d failed: invalid number %q", idx, tokens[p])
			}
			if val != expectGrid[i][j] {
				return fmt.Errorf("case %d failed at (%d,%d): expected %d got %d", idx, i, j, expectGrid[i][j], val)
			}
			p++
		}
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	testcases := parseTestcases()
	for i, tc := range testcases {
		if err := runCase(bin, i+1, tc); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(testcases))
}
