package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"
)

// Embedded source for the reference solution (was 1185D.go).
const solutionSource = `package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	var n int
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return
	}
	a := make([]pair, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &a[i].val)
		a[i].idx = i + 1
	}
	if n <= 3 {
		fmt.Println(1)
		return
	}
	sort.Slice(a, func(i, j int) bool { return a[i].val < a[j].val })
	// Scenario 1: remove one element in middle
	pos := -1
	ok := true
	if (a[n-1].val-a[0].val)%(n-2) == 0 {
		d := (a[n-1].val - a[0].val) / (n - 2)
		cnt := 0
		first := true
		for i := 0; i < n-1; i++ {
			expected := a[0].val + cnt*d
			if a[i].val != expected {
				if first {
					first = false
					pos = a[i].idx
					continue
				}
				ok = false
				break
			}
			cnt++
		}
	} else {
		ok = false
	}
	if ok && pos != -1 {
		fmt.Println(pos)
		return
	}
	// Scenario 2: remove last element
	ok = true
	if (a[n-2].val-a[0].val)%(n-2) == 0 {
		d := (a[n-2].val - a[0].val) / (n - 2)
		for i := 0; i < n-1; i++ {
			if a[i].val != a[0].val+i*d {
				ok = false
				break
			}
		}
	} else {
		ok = false
	}
	if ok {
		// removed element is the last one
		fmt.Println(a[n-1].idx)
		return
	}
	// Scenario 3: remove first element
	ok = true
	if (a[n-1].val-a[1].val)%(n-2) == 0 {
		d := (a[n-1].val - a[1].val) / (n - 2)
		for i := 0; i < n-2; i++ {
			if a[i+1].val != a[1].val+i*d {
				ok = false
				break
			}
		}
	} else {
		ok = false
	}
	if ok {
		fmt.Println(a[0].idx)
	} else {
		fmt.Println(-1)
	}
}

type pair struct {
	val int
	idx int
}
`

const testcasesRaw = `6 20 7 47 26 31 10
4 5 2 26 36
7 49 4 15 34 35 24 18
5 7 17 14 2 42
7 18 13 11 20 19 41 47
8 6 39 22 43 25 33 16 12
6 31 18 6 36 20 1
7 37 46 20 49 33 13 27
9 39 19 28 29 11 15 20 17 3
4 3 30 41 18
10 45 22 10 44 13 5 27 13 41 41
10 18 12 23 28 48 38 21 41 36 13
8 7 4 46 15 18 49 38 40
6 8 22 12 19 30 2
3 23 45 6
7 48 44 21 2 21 19 21
5 50 42 27 40 44
4 19 40 13 29
7 9 17 25 39 11 22 37
3 24 3 30
6 36 11 32 5 1 31
8 17 5 34 11 30 32 9 28
5 10 49 27 1 24
8 43 17 9 15 34 25 32 15
5 33 25 35 4 41
3 38 28 40
4 35 29 12 16
4 40 15 41 29
8 4 21 47 14 22 42 3 9
8 32 38 35 31 21 10 15 32
6 16 43 12 43 16 20
9 28 32 28 24 49 3 25 7 36
9 37 1 42 40 26 25 17 20 41
8 34 2 46 4 19 24 3 6
10 6 13 33 3 45 21 47 28 38 11
4 17 41 22 8
6 28 13 45 24 35 50
7 30 3 7 16 27 11 36
10 28 10 9 43 11 17 23 11 42 36
9 28 12 25 38 22 24 12 2 27
1 38
6 1 18 36 47 4 38
9 21 22 32 1 6 9 38 5 1
6 23 46 24 2 30 45
9 49 3 29 2 12 13 1 42 38
6 16 28 2 18 47 27
7 48 33 33 32 23 9 47
4 50 25 19 25
7 32 46 17 3 41 46 45
2 6 1
5 26 6 20 36 32
1 41
2 43 39
9 21 40 3 4 9 6 29 7 50
8 47 50 10 18 20 43 30 47
8 3 36 48 40 23 38 20 41
1 45
7 20 16 20 41 39 41 49
7 21 24 16 41 19 35 37
6 36 38 20 2 5 10
7 1 5 8 40 32 34 34
2 31 28
2 17 1
10 17 35 36 1 42 29 46 40 16 24
6 36 9 16 46 3 37
1 26
4 3 25 33 20
7 2 38 3 20 37 12 28
10 32 38 6 33 36 3 9 2 14 17
3 13 7 3
7 16 11 5 24 33 14 9
6 12 2 5 21 14 30
1 41
2 9 8
5 29 20 7 2 20
7 45 19 35 37 23 38 14
3 22 17 32
9 50 39 23 4 12 4 6 22 50
6 36 27 43 5 20 18
3 10 5 8
10 35 25 28 9 2 2 6 39 34 14
1 24
2 37 44
10 45 4 43 25 20 18 23 17 44 6
6 12 34 5 22 4 11
6 49 16 6 37 15 20
2 4 14
7 12 30 41 2 10 29 27
6 47 28 6 5 47 25
9 16 35 20 12 21 30 22 9 30
7 14 19 14 13 37 5 45
3 22 31 3
6 19 35 23 9 36 16
1 5
4 22 50 33 46
7 41 9 48 5 21 47 50
7 20 25 35 47 45 12 49
5 11 20 18 2 7
10 16 34 1 15 13 15 10 20 11 49
4 20 39 15 30
4 23 18 16 4
7 2 35 42 21 38 17 1
8 8 6 49 48 39 18 3 14
4 14 28 5 28
10 16 13 43 37 31 18 49 36 5 40
8 10 5 44 39 17 4 33 28
10 45 21 34 3 25 16 43 45 25 48
4 37 8 37 42
6 7 25 49 13 20 36`

var _ = solutionSource

type pair struct {
	val int
	idx int
}

func expected(nums []int) int {
	n := len(nums)
	a := make([]pair, n)
	for i, v := range nums {
		a[i] = pair{val: v, idx: i + 1}
	}
	if n <= 3 {
		return 1
	}
	sort.Slice(a, func(i, j int) bool { return a[i].val < a[j].val })
	pos := -1
	ok := true
	if (a[n-1].val-a[0].val)%(n-2) == 0 {
		d := (a[n-1].val - a[0].val) / (n - 2)
		cnt := 0
		first := true
		for i := 0; i < n-1; i++ {
			expected := a[0].val + cnt*d
			if a[i].val != expected {
				if first {
					first = false
					pos = a[i].idx
					continue
				}
				ok = false
				break
			}
			cnt++
		}
	} else {
		ok = false
	}
	if ok && pos != -1 {
		return pos
	}
	ok = true
	if (a[n-2].val-a[0].val)%(n-2) == 0 {
		d := (a[n-2].val - a[0].val) / (n - 2)
		for i := 0; i < n-1; i++ {
			if a[i].val != a[0].val+i*d {
				ok = false
				break
			}
		}
	} else {
		ok = false
	}
	if ok {
		return a[n-1].idx
	}
	ok = true
	if (a[n-1].val-a[1].val)%(n-2) == 0 {
		d := (a[n-1].val - a[1].val) / (n - 2)
		for i := 0; i < n-2; i++ {
			if a[i+1].val != a[1].val+i*d {
				ok = false
				break
			}
		}
	} else {
		ok = false
	}
	if ok {
		return a[0].idx
	}
	return -1
}

func run(bin, input string) (string, error) {
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
		return "", fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierD.go /path/to/binary")
		return
	}
	bin := os.Args[1]

	scanner := bufio.NewScanner(strings.NewReader(testcasesRaw))
	idx := 0
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		idx++
		fields := strings.Fields(line)
		n, err := strconv.Atoi(fields[0])
		if err != nil || len(fields) != n+1 {
			fmt.Printf("invalid line %d\n", idx)
			os.Exit(1)
		}
		nums := make([]int, n)
		for i := 0; i < n; i++ {
			nums[i], _ = strconv.Atoi(fields[i+1])
		}
		want := strconv.Itoa(expected(nums))
		var sb strings.Builder
		sb.WriteString(fields[0])
		for i := 1; i < len(fields); i++ {
			sb.WriteByte(' ')
			sb.WriteString(fields[i])
		}
		sb.WriteByte('\n')
		input := sb.String()
		got, err := run(bin, input)
		if err != nil {
			fmt.Printf("candidate runtime error on case %d: %v\n", idx, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != want {
			fmt.Printf("case %d failed\ninput:%sexpected:%s got:%s\n", idx, input, want, got)
			os.Exit(1)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Println("scanner error:", err)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", idx)
}
