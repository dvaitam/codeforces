package main

import (
	"fmt"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"
)

const solution1312BSource = `package main

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

   var t int
   if _, err := fmt.Fscan(reader, &t); err != nil {
       return
   }
   for tc := 0; tc < t; tc++ {
       var n int
       fmt.Fscan(reader, &n)
       a := make([]int, n)
       for i := 0; i < n; i++ {
           fmt.Fscan(reader, &a[i])
       }
       // sort in descending order to ensure i - a[i] values are distinct
       sort.Sort(sort.Reverse(sort.IntSlice(a)))
       // output
       for i, v := range a {
           if i > 0 {
               writer.WriteByte(' ')
           }
           fmt.Fprint(writer, v)
       }
       writer.WriteByte('\n')
   }
}
`

// Keep the embedded reference solution reachable.
var _ = solution1312BSource

type testCase struct {
	n   int
	arr []int
}

const testcasesRaw = `3 73 98 9
5 16 64 98 58 61
7 27 13 63 4 50 56 78
1 90
8 35 93 30 76 14 41 4 3
1 84
9 2 49 88 28 55 93 4 68 29
8 64 71 30 45 30 87 29 98
8 38 3 54 72 83 13 24 81
5 16 96 43 93 92
9 55 65 86 25 39 37 76 64 65
7 76 5 62 32 96 52 54
3 47 71 90
6 12 57 85 66 14 100
3 67 51 48
8 94 4 61 6 40 91 79 76
10 51 83 22 22 65 30 2 99 26 70
9 30 52 66 45 74 46 59 35 85
9 78 94 1 50 95 66 17 67 100
9 27 55 8 62 47 73 71 26 65
9 77 58 35 1 37 97 43 12 92
8 61 95 41 69 78 17 56 49
4 21 61 43 8
5 73 2 41 52 40
9 61 17 88 40 19 59 31 87 10
3 36 86 9
8 46 80 65 37 87 91 76 97
10 55 23 36 87 69 3 38 2 74 37
7 50 13 46 44 37 16 49
2 59 44
6 60 96 90 60 77 68
3 16 5 39
1 79
7 10 79 41 50 2 39 22
5 18 70 41 12 49
5 6 2 15 37 51
9 39 5 25 13 12 74 56 36 92
1 22
9 20 24 29 13 87 50 70 82 19
6 23 1 14 59 38 29
8 47 81 33 33 67 71 56 2
9 74 71 25 45 20 37 97 97 60
10 11 20 50 17 3 5 13 48 42 84
3 87 25 87
4 29 9 79 69
9 60 22 92 52 27 37 91 58 61
7 71 48 75 96 28 43 71
8 63 92 81 38 80 45 17 95
2 5 77
9 83 23 20 99 21 46 23 15 72
1 74
7 48 70 25 35 73 18 41
10 24 45 83 34 68 38 71 11 45 25
2 95 24
3 52 7 34
1 14
6 91 67 70 80 72 13
3 3 69 36
9 20 90 66 16 84 67 77 37 57
4 73 11 73 95
1 66
9 14 28 96 8 33 23 28 62 72
10 1 55 82 90 68 67 39 90 70 51
8 35 39 29 33 100 96 19 58
8 70 37 47 92 26 20 19 26
8 25 14 78 10 22 63 96 50
8 73 6 26 10 25 67 7 51
7 48 61 34 94 16 69 27
5 35 88 39 52 96
4 34 22 1 14
7 44 67 21 55 52 93 71
8 36 49 19 61 63 81 54 30
4 14 60 75 8
3 12 73 92
7 19 88 16 66 92 53 30
8 44 59 20 67 27 61 60 68
7 64 42 100 14 97 82 49
8 85 96 9 36 38 20 95 8
10 60 62 54 52 87 71 19 77 55 100
1 20
8 2 62 39 19 4 71 6 13
9 11 20 95 42 12 95 52 40 93
10 40 80 20 95 32 92 26 11 84 5
10 53 21 6 7 10 42 48 84 43 81
5 46 7 23 6 36
10 80 26 12 65 4 60 57 21 9 99
8 20 12 27 88 27 2 18 17
5 95 92 11 32 94
10 52 23 58 62 70 19 43 11 25 55
4 72 37 47 18
6 36 10 78 67 63 28
3 77 84 55
4 13 84 12 13
2 38 63
4 10 34 61 42
5 12 18 100 67 94
4 89 68 22 88
3 57 60 53
5 13 16 58 88 11
`

func parseTestcases() []testCase {
	fields := strings.Fields(testcasesRaw)
	var res []testCase
	for i := 0; i < len(fields); {
		n, _ := strconv.Atoi(fields[i])
		i++
		if i+n > len(fields) {
			break
		}
		arr := make([]int, n)
		for j := 0; j < n; j++ {
			arr[j], _ = strconv.Atoi(fields[i+j])
		}
		i += n
		res = append(res, testCase{n: n, arr: arr})
	}
	return res
}

func solveExpected(tc testCase) []int {
	arr := append([]int(nil), tc.arr...)
	sort.Sort(sort.Reverse(sort.IntSlice(arr)))
	return arr
}

func buildInput(tc testCase) string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "1\n%d\n", tc.n)
	for i, v := range tc.arr {
		if i > 0 {
			sb.WriteByte(' ')
		}
		fmt.Fprintf(&sb, "%d", v)
	}
	sb.WriteByte('\n')
	return sb.String()
}

func runCase(bin string, idx int, tc testCase) error {
	input := buildInput(tc)
	expectArr := solveExpected(tc)
	expect := make([]string, len(expectArr))
	for i, v := range expectArr {
		expect[i] = strconv.Itoa(v)
	}
	expectStr := strings.Join(expect, " ")

	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("case %d failed: %v\nstderr: %s\ninput:\n%s", idx, err, string(out), input)
	}
	got := strings.TrimSpace(string(out))
	if got != expectStr {
		return fmt.Errorf("case %d failed\nexpected: %s\ngot: %s\ninput:\n%s", idx, expectStr, got, input)
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
