package main

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

const solution1101GSource = `package main

import (
	"bufio"
	"fmt"
	"os"
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
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &a[i])
	}
	// Compute total xor
	x := 0
	for _, v := range a {
		x ^= v
	}
	if x == 0 {
		fmt.Fprintln(writer, -1)
		return
	}
	ans := 0
	// Gaussian elimination over GF(2) to find basis size
	for i := 29; i >= 0; i-- {
		// find vector with bit i set
		id := -1
		mask := 1 << i
		for j, v := range a {
			if v&mask != 0 {
				id = j
				break
			}
		}
		if id == -1 {
			continue
		}
		ans++
		// move pivot to end
		last := len(a) - 1
		a[id], a[last] = a[last], a[id]
		// eliminate bit i from all other vectors
		pivot := a[last]
		for j := 0; j < last; j++ {
			if a[j]&mask != 0 {
				a[j] ^= pivot
			}
		}
		// remove pivot
		a = a[:last]
	}
	fmt.Fprintln(writer, ans)
}
`

// Keep the embedded reference solution reachable so it is preserved in the binary.
var _ = solution1101GSource

type testCase struct {
	n   int
	arr []int
}

var testcases = []testCase{
	{n: 5, arr: []int{8, 6, 12, 19, 14}},
	{n: 8, arr: []int{9, 10, 0, 18, 19, 19, 9, 3}},
	{n: 6, arr: []int{2, 4, 12, 14, 8, 11}},
	{n: 5, arr: []int{19, 7, 5, 17, 14}},
	{n: 8, arr: []int{8, 15, 3, 2, 5, 17, 1, 8}},
	{n: 6, arr: []int{10, 11, 1, 12, 10, 5}},
	{n: 5, arr: []int{15, 1, 17, 16, 9}},
	{n: 6, arr: []int{8, 19, 20, 10, 2, 11}},
	{n: 1, arr: []int{14}},
	{n: 3, arr: []int{11, 19, 9}},
	{n: 5, arr: []int{15, 0, 14, 5, 7}},
	{n: 6, arr: []int{7, 3, 20, 1, 12, 14}},
	{n: 7, arr: []int{20, 19, 1, 10, 0, 0, 14}},
	{n: 5, arr: []int{6, 5, 14, 3, 12}},
	{n: 1, arr: []int{8}},
	{n: 2, arr: []int{4, 1}},
	{n: 1, arr: []int{8}},
	{n: 7, arr: []int{16, 18, 9, 2, 10, 15, 18}},
	{n: 6, arr: []int{8, 10, 7, 1, 16, 7}},
	{n: 8, arr: []int{16, 0, 12, 10, 1, 16, 6, 3}},
	{n: 4, arr: []int{0, 17, 10, 4}},
	{n: 4, arr: []int{18, 15, 17, 15}},
	{n: 7, arr: []int{5, 14, 8, 5, 1, 11, 9}},
	{n: 6, arr: []int{20, 8, 7, 12, 0, 1}},
	{n: 2, arr: []int{4, 19}},
	{n: 6, arr: []int{18, 5, 0, 1, 17, 5}},
	{n: 3, arr: []int{2, 4, 4}},
	{n: 3, arr: []int{13, 12, 0}},
	{n: 6, arr: []int{9, 11, 18, 0, 7, 18}},
	{n: 1, arr: []int{11}},
	{n: 6, arr: []int{10, 2, 18, 11, 11, 2}},
	{n: 7, arr: []int{6, 11, 18, 8, 4, 6, 16}},
	{n: 3, arr: []int{19, 4, 16}},
	{n: 2, arr: []int{10, 18}},
	{n: 2, arr: []int{19, 5}},
	{n: 6, arr: []int{8, 14, 8, 14, 5, 12}},
	{n: 7, arr: []int{11, 3, 10, 4, 8, 5, 10}},
	{n: 8, arr: []int{2, 11, 11, 8, 11, 4, 13, 6}},
	{n: 4, arr: []int{5, 18, 14, 19}},
	{n: 7, arr: []int{18, 7, 0, 11, 4, 0, 14}},
	{n: 7, arr: []int{18, 3, 1, 12, 9, 19, 6}},
	{n: 7, arr: []int{0, 19, 7, 12, 19, 15, 15}},
	{n: 7, arr: []int{12, 12, 1, 9, 5, 5, 10}},
	{n: 2, arr: []int{16, 19}},
	{n: 8, arr: []int{5, 18, 12, 7, 4, 8, 0, 14}},
	{n: 6, arr: []int{15, 11, 16, 7, 9, 18}},
	{n: 1, arr: []int{15}},
	{n: 3, arr: []int{8, 6, 11}},
	{n: 1, arr: []int{9}},
	{n: 7, arr: []int{7, 0, 7, 19, 12, 4, 8}},
	{n: 4, arr: []int{4, 12, 18, 15}},
	{n: 6, arr: []int{15, 19, 4, 0, 19, 12}},
	{n: 7, arr: []int{11, 9, 4, 9, 15, 8, 12}},
	{n: 4, arr: []int{14, 6, 5, 15}},
	{n: 8, arr: []int{1, 13, 6, 16, 9, 13, 3, 9}},
	{n: 5, arr: []int{5, 4, 6, 9, 4}},
	{n: 3, arr: []int{11, 6, 18}},
	{n: 3, arr: []int{1, 3, 4}},
	{n: 5, arr: []int{18, 9, 9, 5, 1}},
	{n: 5, arr: []int{18, 11, 3, 13, 8}},
	{n: 6, arr: []int{19, 17, 16, 20, 8, 3}},
	{n: 1, arr: []int{0}},
	{n: 6, arr: []int{13, 19, 15, 10, 18, 7}},
	{n: 2, arr: []int{13, 13}},
	{n: 2, arr: []int{18, 8}},
	{n: 6, arr: []int{4, 1, 0, 4, 12, 11}},
	{n: 1, arr: []int{11}},
	{n: 6, arr: []int{16, 18, 18, 12, 17, 3}},
	{n: 2, arr: []int{19, 17}},
	{n: 4, arr: []int{6, 10, 20, 17}},
	{n: 7, arr: []int{4, 1, 2, 1, 20, 11, 6}},
	{n: 7, arr: []int{18, 20, 3, 13, 6, 10, 2}},
	{n: 2, arr: []int{4, 8}},
	{n: 4, arr: []int{5, 6, 12, 13}},
	{n: 4, arr: []int{6, 0, 16, 8}},
	{n: 7, arr: []int{6, 16, 3, 19, 11, 20, 10}},
	{n: 5, arr: []int{16, 7, 15, 18, 6}},
	{n: 1, arr: []int{18}},
	{n: 2, arr: []int{8, 13}},
	{n: 2, arr: []int{1, 3}},
	{n: 6, arr: []int{10, 3, 19, 9, 3, 7}},
	{n: 5, arr: []int{18, 2, 3, 12, 15}},
	{n: 3, arr: []int{10, 7, 4}},
	{n: 4, arr: []int{4, 1, 18, 1}},
	{n: 8, arr: []int{15, 17, 17, 7, 13, 7, 6, 19}},
	{n: 2, arr: []int{10, 5}},
	{n: 5, arr: []int{3, 14, 14, 6, 20}},
	{n: 8, arr: []int{8, 15, 6, 14, 5, 19, 6, 18}},
	{n: 7, arr: []int{11, 13, 8, 12, 5, 4, 19}},
	{n: 8, arr: []int{3, 7, 19, 11, 18, 8, 0, 9}},
	{n: 8, arr: []int{6, 12, 7, 9, 18, 10, 0, 11}},
	{n: 1, arr: []int{13}},
	{n: 1, arr: []int{18}},
	{n: 1, arr: []int{1}},
	{n: 1, arr: []int{13}},
	{n: 5, arr: []int{14, 10, 20, 0, 17}},
	{n: 7, arr: []int{3, 13, 8, 1, 10, 1, 1}},
	{n: 6, arr: []int{20, 3, 19, 0, 20, 16}},
	{n: 2, arr: []int{18, 17}},
	{n: 1, arr: []int{6}},
}

func solveCaseG(cs testCase) int {
	a := append([]int{}, cs.arr...)
	x := 0
	for _, v := range a {
		x ^= v
	}
	if x == 0 {
		return -1
	}
	ans := 0
	for i := 29; i >= 0; i-- {
		id := -1
		mask := 1 << uint(i)
		for j, v := range a {
			if v&mask != 0 {
				id = j
				break
			}
		}
		if id == -1 {
			continue
		}
		ans++
		last := len(a) - 1
		a[id], a[last] = a[last], a[id]
		pivot := a[last]
		for j := 0; j < last; j++ {
			if a[j]&mask != 0 {
				a[j] ^= pivot
			}
		}
		a = a[:last]
	}
	return ans
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierG.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	if bin == "--" && len(os.Args) > 2 {
		bin = os.Args[2]
	}
	for idx, cs := range testcases {
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d\n", cs.n))
		for i, v := range cs.arr {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.Itoa(v))
		}
		sb.WriteByte('\n')
		input := sb.String()
		expected := solveCaseG(cs)
		cmd := exec.Command(bin)
		cmd.Stdin = strings.NewReader(input)
		out, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", idx+1, err, input)
			os.Exit(1)
		}
		got, err := strconv.Atoi(strings.TrimSpace(string(out)))
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d bad output: %v", idx+1, err)
			os.Exit(1)
		}
		if got != expected {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %d got %d\ninput:\n%s", idx+1, expected, got, input)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(testcases))
}
