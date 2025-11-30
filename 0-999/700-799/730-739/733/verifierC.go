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

type op struct {
        pos int
        typ byte
}

// Embedded testcases from testcasesC.txt.
const testcasesRaw = `100
1
5
1
5
1
3
1
3
4
2 4 3 3
3
7 3 2
1
5
1
5
3
5 5 3
1
13
6
3 3 2 1 4 3
6
6 4 2 2 1 1
6
3 1 5 4 2 2
2
9 8
2
3 2
2
2 3
1
2
1
2
4
1 3 2 1
1
7
1
5
1
5
6
3 1 1 4 5 3
5
13 1 1 1 1
5
5 1 5 2 5
2
3 15
1
5
1
5
3
4 1 4
2
8 1
2
5 4
1
9
2
2 5
2
3 4
2
3 4
1
7
2
3 4
1
7
1
3
1
3
5
3 3 1 3 2
5
2 7 1 1 1
1
3
1
3
1
3
1
3
5
4 2 3 1 4
1
14
5
3 1 4 5 2
1
15
4
2 5 3 4
3
11 1 2
5
5 5 2 2 2
5
1 8 5 1 1
3
3 2 1
2
1 5
6
2 2 3 4 5 1
5
4 1 6 2 4
5
1 4 1 4 2
4
9 1 1 1
5
4 2 3 4 4
3
5 3 9
6
3 2 5 4 4 4
3
2 3 17
4
5 4 4 1
1
14
4
3 2 4 2
4
6 1 1 3
6
3 2 5 5 2 1
5
8 6 1 2 1
5
1 1 3 3 3
3
4 6 1
3
1 4 5
2
5 5
5
2 3 1 1 2
1
9
1
5
1
5
3
3 4 3
3
5 1 4
6
4 1 5 1 2 3
3
12 1 3
3
1 2 1
1
4
3
5 4 1
2
7 3
2
5 4
2
4 5
6
4 1 4 1 1 2
5
1 5 1 3 3
5
5 2 5 3 1
1
16
1
3
1
3
4
4 4 3 3
1
14
3
3 2 2
1
7
2
5 4
2
6 3
5
1 5 4 1 1
2
3 9
5
4 1 2 1 5
1
13
2
2 2
2
3 1
4
1 2 4 2
3
4 3 2
5
4 2 4 3 2
3
5 9 1
6
2 5 4 2 4 1
1
18
6
4 2 3 1 4 3
4
10 3 2 2
3
4 4 2
2
3 7
6
3 3 3 4 3 3
1
19
5
2 1 4 4 2
2
9 4
1
2
1
2
3
4 5 3
1
12
4
2 2 5 1
2
9 1
3
5 1 2
3
4 2 2
1
5
1
5
6
5 1 4 2 2 3
2
4 13
5
3 5 1 2 2
2
11 2
3
3 2 1
3
3 1 2
2
1 2
2
1 2
5
1 4 2 4 2
2
10 3
1
5
1
5
4
3 1 4 5
3
1 6 6
3
5 2 2
3
2 1 6
2
5 5
2
4 6
3
5 2 3
2
3 7
4
2 2 2 5
2
1 10
1
2
1
2
6
1 1 3 5 3 2
3
2 4 9
6
2 1 3 3 5 5
3
1 9 9
1
2
1
2
5
2 4 2 1 5
4
7 3 3 1
3
4 2 4
1
10
1
2
1
2
5
2 1 1 4 2
1
10
5
4 5 1 4 1
4
5 6 3 1
1
5
1
5
2
4 2
2
2 4
1
4
1
4
6
5 1 1 4 3 2
4
2 5 7 2
6
2 4 2 3 3 2
6
10 1 2 1 1 1
5
1 5 1 3 3
1
13
1
1
1
1
3
4 4 3
1
11
4
2 4 5 4
2
8 7
5
4 5 2 5 1
2
2 15
5
1 2 3 5 4
3
13 1 1
6
5 4 4 5 2 3
1
23
4
2 3 3 2
1
10
5
1 4 1 3 2
1
11
2
1 2
2
2 1`
// parseInput reads a single test case input and returns the initial and final
// queues.
func parseInput(input string) (n int, a []int64, k int, b []int64, err error) {
        reader := bufio.NewReader(strings.NewReader(input))
        if _, err = fmt.Fscan(reader, &n); err != nil {
                return
        }
        a = make([]int64, n)
        for i := 0; i < n; i++ {
                if _, err = fmt.Fscan(reader, &a[i]); err != nil {
                        return
                }
        }
        if _, err = fmt.Fscan(reader, &k); err != nil {
                return
        }
        b = make([]int64, k)
        for i := 0; i < k; i++ {
                if _, err = fmt.Fscan(reader, &b[i]); err != nil {
                        return
                }
        }
        return
}

// checkSolution verifies that the contestant's output is a valid sequence of
// operations producing the desired final queue.
func checkSolution(out string, n int, a []int64, k int, b []int64) error {
        lines := strings.Split(strings.TrimSpace(out), "\n")
        if len(lines) == 0 || strings.TrimSpace(lines[0]) != "YES" {
                return fmt.Errorf("expected YES got %q", out)
        }
        ops := lines[1:]
        if len(ops) != n-k {
                return fmt.Errorf("expected %d operations got %d", n-k, len(ops))
        }
        arr := append([]int64(nil), a...)
        for _, line := range ops {
                fields := strings.Fields(line)
                if len(fields) != 2 {
                        return fmt.Errorf("invalid operation %q", line)
                }
                pos, err := strconv.Atoi(fields[0])
                if err != nil {
                        return fmt.Errorf("invalid index in %q", line)
                }
                if pos < 1 || pos > len(arr) {
                        return fmt.Errorf("index out of range in %q", line)
                }
                dir := fields[1]
                switch dir {
                case "L":
                        if pos == 1 {
                                return fmt.Errorf("cannot eat left from position %d", pos)
                        }
                        if arr[pos-1] <= arr[pos-2] {
                                return fmt.Errorf("eater not heavier in %q", line)
                        }
                        arr[pos-2] += arr[pos-1]
                        arr = append(arr[:pos-1], arr[pos:]...)
                case "R":
                        if pos == len(arr) {
                                return fmt.Errorf("cannot eat right from position %d", pos)
                        }
                        if arr[pos-1] <= arr[pos] {
                                return fmt.Errorf("eater not heavier in %q", line)
                        }
                        arr[pos-1] += arr[pos]
                        arr = append(arr[:pos], arr[pos+1:]...)
                default:
                        return fmt.Errorf("unknown direction in %q", line)
                }
        }
        if len(arr) != k {
                return fmt.Errorf("expected %d monsters got %d", k, len(arr))
        }
        for i := 0; i < k; i++ {
                if arr[i] != b[i] {
                        return fmt.Errorf("final weights mismatch at %d: got %d expected %d", i, arr[i], b[i])
                }
        }
        return nil
}

func parseTestcases() ([]string, error) {
	scan := bufio.NewScanner(strings.NewReader(testcasesRaw))
	scan.Split(bufio.ScanWords)
	if !scan.Scan() {
		return nil, fmt.Errorf("missing test count")
	}
	t, err := strconv.Atoi(scan.Text())
	if err != nil {
		return nil, fmt.Errorf("bad test count: %v", err)
	}
	var cases []string
	for i := 0; i < t; i++ {
		if !scan.Scan() {
			return nil, fmt.Errorf("unexpected EOF reading n for case %d", i+1)
		}
		n, err := strconv.Atoi(scan.Text())
		if err != nil {
			return nil, fmt.Errorf("case %d bad n: %v", i+1, err)
		}
		a := make([]string, n)
		for j := 0; j < n; j++ {
			if !scan.Scan() {
				return nil, fmt.Errorf("case %d incomplete a", i+1)
			}
			a[j] = scan.Text()
		}
		if !scan.Scan() {
			return nil, fmt.Errorf("case %d missing k", i+1)
		}
		k, err := strconv.Atoi(scan.Text())
		if err != nil {
			return nil, fmt.Errorf("case %d bad k: %v", i+1, err)
		}
		b := make([]string, k)
		for j := 0; j < k; j++ {
			if !scan.Scan() {
				return nil, fmt.Errorf("case %d incomplete b", i+1)
			}
			b[j] = scan.Text()
		}
		var sb strings.Builder
		sb.WriteString(strconv.Itoa(n))
		sb.WriteByte('\n')
		for j := 0; j < n; j++ {
			if j > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(a[j])
		}
		sb.WriteByte('\n')
		sb.WriteString(strconv.Itoa(k))
		sb.WriteByte('\n')
		for j := 0; j < k; j++ {
			if j > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(b[j])
		}
		sb.WriteByte('\n')
		cases = append(cases, sb.String())
	}
	return cases, nil
}

func expectedCase(input string) string {
	reader := bufio.NewReader(strings.NewReader(input))
	var n int
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return ""
	}
	a := make([]int64, n)
	var suma int64
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &a[i])
		suma += a[i]
	}
	var k int
	fmt.Fscan(reader, &k)
	b := make([]int64, k)
	var sumb int64
	for i := 0; i < k; i++ {
		fmt.Fscan(reader, &b[i])
		sumb += b[i]
	}
	if suma != sumb {
		return "NO\n"
	}
	flag := make([]int, n)
	nxt := make([]int, n)
	var cur, cnt int
	var temp int64
	for i := 0; i < n; i++ {
		temp += a[i]
		if temp == b[cur] {
			flag[cnt] = i
			cnt++
			cur++
			temp = 0
		} else if temp > b[cur] {
			return "NO\n"
		}
	}
	if cur != k {
		return "NO\n"
	}
	var ans []op
	initNxt := func(left, right int) {
		for i := left; i < right; i++ {
			nxt[i] = i + 1
		}
		if right < len(nxt) {
			nxt[right] = -1
		}
	}
	check := func(left, right int) bool {
		return nxt[left] == -1
	}
	var findIndex func(left, index int) int
	findIndex = func(left, index int) int {
		cnt := 0
		for i := left; i != -1; i = nxt[i] {
			cnt++
			if index == i {
				return cnt
			}
		}
		return -1
	}
	solve := func(left, right, pre int) int64 {
		var mx int64 = -1
		var index int
		first := true
		for i := left; i != -1; i = nxt[i] {
			j := nxt[i]
			if j != -1 && a[i] != a[j] {
				sum := a[i] + a[j]
				if first {
					first = false
					mx = sum
					index = i
				} else if sum > mx {
					mx = sum
					index = i
				}
			}
		}
		if mx != -1 {
			pos1 := findIndex(left, index)
			pos2 := findIndex(left, nxt[index])
			if a[index] > a[nxt[index]] {
				ans = append(ans, op{pos: pre + pos1, typ: 'R'})
			} else {
				ans = append(ans, op{pos: pre + pos2, typ: 'L'})
			}
			a[index] += a[nxt[index]]
			nxt[index] = nxt[nxt[index]]
		}
		return mx
	}
	noFail := false
	for i := 0; i < cnt; i++ {
		var left int
		right := flag[i]
		if i == 0 {
			left = 0
		} else {
			left = flag[i-1] + 1
		}
		initNxt(left, right)
		for !check(left, right) {
			p := solve(left, right, i)
			if p == -1 {
				noFail = true
				break
			}
		}
		if noFail {
			break
		}
	}
	if noFail {
		return "NO\n"
	}
	var sb strings.Builder
	sb.WriteString("YES\n")
	for _, t := range ans {
		sb.WriteString(fmt.Sprintf("%d %c\n", t.pos, t.typ))
	}
	return sb.String()
}

func runCase(exe string, input string, expected string) error {
        cmd := exec.Command(exe)
        cmd.Stdin = strings.NewReader(input)
        var out bytes.Buffer
        cmd.Stdout = &out
        cmd.Stderr = &out
        if err := cmd.Run(); err != nil {
                return fmt.Errorf("runtime error: %v\n%s", err, out.String())
        }
        got := out.String()
        // If our reference solver says NO, contestant must also say NO.
        if strings.HasPrefix(strings.TrimSpace(expected), "NO") {
                if strings.TrimSpace(got) != "NO" {
                        return fmt.Errorf("expected NO got %q", strings.TrimSpace(got))
                }
                return nil
        }
        n, a, k, b, err := parseInput(input)
        if err != nil {
                return fmt.Errorf("failed to parse input: %v", err)
        }
        if strings.TrimSpace(got) == "NO" {
                return fmt.Errorf("expected YES got NO")
        }
        if err := checkSolution(got, n, a, k, b); err != nil {
                return err
        }
        return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	exe := os.Args[1]
	cases, err := parseTestcases()
	if err != nil {
		fmt.Println("failed to parse embedded testcases:", err)
		os.Exit(1)
	}
	for i, input := range cases {
		exp := expectedCase(input)
		if err := runCase(exe, input, exp); err != nil {
			fmt.Printf("case %d failed: %v\ninput:%s", i+1, err, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
