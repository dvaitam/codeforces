package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

type testCase struct {
	input  string
	expect string
}

const inf64 = int64(4e18)

func solveRef(input string) string {
	reader := bufio.NewReader(strings.NewReader(input))
	var n int
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return ""
	}
	w := make([]int, n)
	h := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &w[i], &h[i])
	}
	k := n / 2
	candMap := make(map[int]struct{})
	for i := 0; i < n; i++ {
		candMap[w[i]] = struct{}{}
		candMap[h[i]] = struct{}{}
	}
	cand := make([]int, 0, len(candMap))
	for hh := range candMap {
		cand = append(cand, hh)
	}
	sortInts(cand)
	best := inf64
	tmp := make([]int64, 0, n)
	for _, H := range cand {
		baseW := int64(0)
		mandatory := 0
		ok := true
		tmp = tmp[:0]
		for i := 0; i < n; i++ {
			wi, hi := w[i], h[i]
			switch {
			case wi <= H && hi <= H:
				baseW += int64(wi)
				tmp = append(tmp, int64(hi-wi))
			case hi <= H:
				baseW += int64(wi)
			case wi <= H:
				mandatory++
				baseW += int64(hi)
			default:
				ok = false
			}
			if !ok {
				break
			}
		}
		if !ok || mandatory > k {
			continue
		}
		sortInt64s(tmp)
		limit := k - mandatory
		for i := 0; i < len(tmp) && i < limit; i++ {
			if tmp[i] < 0 {
				baseW += tmp[i]
			} else {
				break
			}
		}
		area := baseW * int64(H)
		if area < best {
			best = area
		}
	}
	if best == inf64 {
		return "0"
	}
	return fmt.Sprintf("%d", best)
}

func makeTest(n int, rects [][2]int) testCase {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for _, rc := range rects {
		sb.WriteString(fmt.Sprintf("%d %d\n", rc[0], rc[1]))
	}
	input := sb.String()
	return testCase{input: input, expect: solveRef(input)}
}

func genTests() []testCase {
	rand.Seed(42)
	tests := []testCase{
		makeTest(1, [][2]int{{5, 7}}),
		makeTest(2, [][2]int{{3, 4}, {5, 6}}),
		makeTest(3, [][2]int{{2, 9}, {4, 5}, {6, 1}}),
	}
	for t := 0; t < 200; t++ {
		n := rand.Intn(6) + 1
		rects := make([][2]int, n)
		for i := 0; i < n; i++ {
			rects[i][0] = rand.Intn(10) + 1
			rects[i][1] = rand.Intn(10) + 1
		}
		tests = append(tests, makeTest(n, rects))
	}
	return tests
}

func run(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func checkOutput(expect, output string) error {
	output = strings.TrimSpace(output)
	if output == "" {
		return fmt.Errorf("empty output")
	}
	val, err := strconv.ParseInt(output, 10, 64)
	if err != nil {
		return fmt.Errorf("invalid integer output %q", output)
	}
	exp, _ := strconv.ParseInt(expect, 10, 64)
	if val != exp {
		return fmt.Errorf("expected %s but got %s", expect, output)
	}
	return nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := genTests()
	for i, tc := range tests {
		got, err := run(bin, tc.input)
		if err != nil {
			fmt.Printf("test %d runtime error: %v\ninput:\n%s\n", i+1, err, tc.input)
			os.Exit(1)
		}
		if err := checkOutput(tc.expect, got); err != nil {
			fmt.Printf("test %d failed: %v\ninput:\n%s\nexpected: %s\nactual: %s\n", i+1, err, tc.input, tc.expect, got)
			os.Exit(1)
		}
	}
	fmt.Println("all tests passed")
}

func sortInts(a []int) {
	if len(a) <= 1 {
		return
	}
	quickSortInts(a, 0, len(a)-1)
}

func quickSortInts(a []int, l, r int) {
	i, j := l, r
	p := a[(l+r)/2]
	for i <= j {
		for a[i] < p {
			i++
		}
		for a[j] > p {
			j--
		}
		if i <= j {
			a[i], a[j] = a[j], a[i]
			i++
			j--
		}
	}
	if l < j {
		quickSortInts(a, l, j)
	}
	if i < r {
		quickSortInts(a, i, r)
	}
}

func sortInt64s(a []int64) {
	if len(a) <= 1 {
		return
	}
	quickSortInt64s(a, 0, len(a)-1)
}

func quickSortInt64s(a []int64, l, r int) {
	i, j := l, r
	p := a[(l+r)/2]
	for i <= j {
		for a[i] < p {
			i++
		}
		for a[j] > p {
			j--
		}
		if i <= j {
			a[i], a[j] = a[j], a[i]
			i++
			j--
		}
	}
	if l < j {
		quickSortInt64s(a, l, j)
	}
	if i < r {
		quickSortInt64s(a, i, r)
	}
}
