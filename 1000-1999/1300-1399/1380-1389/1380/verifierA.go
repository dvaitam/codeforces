package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

const solution1380ASource = `package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   in := bufio.NewReader(os.Stdin)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()

   var T int
   if _, err := fmt.Fscan(in, &T); err != nil {
       return
   }
   for T > 0 {
       T--
       var n int
       fmt.Fscan(in, &n)
       a := make([]int, n)
       for i := 0; i < n; i++ {
           fmt.Fscan(in, &a[i])
       }
       stack := make([]int, 0, n)
       found := false
       var x, y, z int
       for i := 0; i < n; i++ {
           for len(stack) > 0 && a[stack[len(stack)-1]] > a[i] {
               if len(stack) > 1 && !found {
                   x = stack[len(stack)-2] + 1
                   y = stack[len(stack)-1] + 1
                   z = i + 1
                   found = true
               }
               stack = stack[:len(stack)-1]
           }
           stack = append(stack, i)
       }
       if found {
           fmt.Fprintln(out, "YES")
           fmt.Fprintf(out, "%d %d %d", x, y, z)
       } else {
           fmt.Fprintln(out, "NO")
       }
       if T > 0 {
           fmt.Fprintln(out)
       }
   }
}
`

// Preserve the embedded reference solution.
var _ = solution1380ASource

type testCase struct {
	arr []int
}

const testcasesRaw = `100
9
8 9 2 6 4 5 3 1 7
8
7 8 1 6 3 2 5 4
5
2 4 5 1 3
8
7 6 2 4 3 1 5 8
10
4 2 3 6 7 10 1 5 9 8
3
1 2 3
8
8 5 7 2 1 3 6 4
6
6 3 1 4 5 2
8
2 7 4 6 5 3 1 8
8
8 1 2 3 6 5 7 4
9
7 8 1 5 9 2 3 4 6
10
5 7 8 4 6 1 9 3 10 2
6
1 6 3 4 5 2
10
9 10 7 4 1 5 3 2 6 8
6
6 4 5 3 1 2
8
2 7 6 5 1 4 8 3
3
2 1 3
3
2 3 1
4
3 2 1 4
4
2 4 3 1
6
6 3 5 4 1 2
3
1 2 3
4
2 4 1 3
4
1 4 2 3
5
2 3 5 4 1
9
8 1 7 2 9 6 3 5 4
3
2 1 3
5
4 2 1 5 3
5
5 2 3 4 1
7
2 5 1 7 4 3 6
10
8 9 4 7 10 3 5 1 6 2
7
1 4 7 2 5 3 6
9
6 4 9 3 8 5 1 2 7
6
1 2 3 5 4 6
3
1 3 2
3
3 2 1
4
2 1 4 3
10
6 7 2 5 3 4 10 1 8 9
9
9 3 5 1 8 6 7 2 4
8
7 3 4 5 6 8 2 1
6
5 4 1 2 6 3
9
6 9 5 7 8 4 3 1 2
4
3 1 4 2
3
2 3 1
7
4 2 7 1 6 3 5
10
10 9 1 8 5 4 2 3 6 7
5
1 4 2 5 3
4
2 1 4 3
3
3 1 2
5
1 4 2 3 5
5
5 4 2 1 3
3
3 1 2
4
1 4 2 3
7
5 3 2 7 6 1 4
8
3 8 6 5 7 1 4 2
8
7 1 8 4 5 6 2 3
4
2 3 4 1
3
3 1 2
4
1 2 3 4
10
5 10 9 6 2 3 1 7 4 8
9
2 5 3 9 8 7 1 6 4
4
1 3 2 4
6
4 5 3 2 1 6
4
1 3 4 2
8
7 6 8 2 4 5 1 3
9
2 1 7 5 4 3 6 8 9
6
4 6 2 5 3 1
5
1 5 4 2 3
10
3 4 5 7 6 9 8 1 10 2
10
1 2 5 8 4 10 7 9 3 6
3
2 1 3
4
3 2 1 4
9
6 4 3 1 2 9 5 7 8
5
4 1 2 5 3
10
6 9 3 2 1 8 4 7 10 5
6
5 2 4 1 6 3
9
6 9 4 2 1 5 7 8 3
3
1 2 3
9
5 3 6 7 2 4 9 1 8
8
1 7 4 8 3 5 6 2
4
4 2 3 1
8
4 8 1 3 5 7 6 2
7
3 4 5 2 1 7 6
10
7 2 9 5 1 10 3 8 4 6
7
4 3 7 2 6 5 1
9
6 3 7 8 4 2 1 5 9
9
5 3 6 7 1 8 4 9 2
3
1 3 2
7
6 4 2 5 7 3 1
8
5 4 8 3 7 1 2 6
6
5 4 3 6 2 1
7
4 7 5 2 1 6 3
10
7 5 4 1 10 3 9 6 8 2
8
2 4 1 6 5 7 3 8
4
4 2 1 3
5
3 1 4 5 2
6
4 1 6 5 3 2
5
3 4 2 1 5
4
1 3 4 2
9
5 2 4 1 8 6 9 7 3
`

func parseTestcases() []testCase {
	reader := bufio.NewReader(strings.NewReader(testcasesRaw))
	var t int
	if _, err := fmt.Fscan(reader, &t); err != nil {
		return nil
	}
	res := make([]testCase, t)
	for i := 0; i < t; i++ {
		var n int
		fmt.Fscan(reader, &n)
		arr := make([]int, n)
		for j := 0; j < n; j++ {
			fmt.Fscan(reader, &arr[j])
		}
		res[i] = testCase{arr: arr}
	}
	return res
}

func solveExpected(arr []int) (bool, [3]int) {
	stack := []int{}
	found := false
	var x, y, z int
	for i, v := range arr {
		for len(stack) > 0 && arr[stack[len(stack)-1]] > v {
			if len(stack) > 1 && !found {
				x = stack[len(stack)-2] + 1
				y = stack[len(stack)-1] + 1
				z = i + 1
				found = true
			}
			stack = stack[:len(stack)-1]
		}
		stack = append(stack, i)
	}
	if found {
		return true, [3]int{x, y, z}
	}
	return false, [3]int{}
}

func buildInput(tc testCase) string {
	var sb strings.Builder
	sb.WriteString("1\n")
	sb.WriteString(fmt.Sprintf("%d\n", len(tc.arr)))
	for i, v := range tc.arr {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(v))
	}
	sb.WriteByte('\n')
	return sb.String()
}

func runCase(bin string, idx int, tc testCase) error {
	ok, idxs := solveExpected(tc.arr)
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
	if strings.ToUpper(tokens[0]) == "NO" {
		if ok {
			return fmt.Errorf("case %d failed: expected YES got NO", idx)
		}
		return nil
	}
	if strings.ToUpper(tokens[0]) != "YES" {
		return fmt.Errorf("case %d failed: first token must be YES or NO", idx)
	}
	if !ok {
		return fmt.Errorf("case %d failed: expected NO got YES", idx)
	}
	if len(tokens) < 4 {
		return fmt.Errorf("case %d failed: missing indices", idx)
	}
	i1, err1 := strconv.Atoi(tokens[1])
	i2, err2 := strconv.Atoi(tokens[2])
	i3, err3 := strconv.Atoi(tokens[3])
	if err1 != nil || err2 != nil || err3 != nil {
		return fmt.Errorf("case %d failed: invalid indices", idx)
	}
	if i1 != idxs[0] || i2 != idxs[1] || i3 != idxs[2] {
		return fmt.Errorf("case %d failed: expected %d %d %d got %d %d %d", idx, idxs[0], idxs[1], idxs[2], i1, i2, i3)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
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
