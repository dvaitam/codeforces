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

func runCandidate(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errb bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errb
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errb.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func expected(a []int) string {
	n := len(a)
	visited := make([]bool, n)
	pieces := 0
	count := 0
	pos := 0
	dir := 1
	changes := 0
	for count < n {
		if !visited[pos] && pieces >= a[pos] {
			visited[pos] = true
			pieces++
			count++
			if count == n {
				break
			}
		}
		next := pos + dir
		if next < 0 || next >= n {
			dir = -dir
			changes++
			next = pos + dir
		}
		pos = next
	}
	return strconv.Itoa(changes)
}

const testcasesBRaw = `100
1
0
9
1 0 1 4 3 5 1 0 2
5
1 0 1 1 2
10
5 5 3 1 2 0 2 0 6 6
4
0 2 0 1
1
0
9
4 1 2 3 0 0 1 1 0
5
1 4 2 0 1
8
3 5 0 1 0 2 1 1
1
0
6
3 2 1 0 4 1
10
1 0 4 0 9 1 4 0 2 0
2
0 0
8
7 5 1 0 1 0 5 3
3
1 1 0
6
0 4 0 0 1 2
3
1 1 0
7
1 0 0 0 2 1 3
1
0
2
0 0
6
3 1 1 2 4 0
1
0
5
1 0 1 1 0
7
3 0 0 1 4 1 2
6
0 0 0 0 5 3
5
0 1 4 1 0
5
3 1 0 1 0
7
0 0 3 0 2 2 1
9
0 0 1 0 4 2 0 0 3
9
2 1 0 1 5 4 4 3 0
2
1 0
6
1 0 0 4 2 0
7
0 3 2 0 6 3 0
5
2 0 1 0 2
2
0 0
5
0 1 1 4 0
6
3 3 2 0 2 0
6
0 3 3 0 1 0
7
2 0 0 2 4 0 0
1
0
4
0 1 0 0
4
1 2 0 1
10
4 6 1 1 0 0 1 0 6 3
3
0 0 2
9
2 2 0 0 2 0 2 4 4
7
0 2 3 0 1 4 3
10
2 0 3 2 2 0 3 4 2 7
2
1 0
5
0 0 1 0 1
7
0 1 5 3 0 1 1
9
4 6 1 1 6 5 0 0 3
1
0
8
0 3 1 2 1 4 0 1
1
0
4
0 1 0 0
5
0 0 2 0 0
2
1 0
7
5 0 2 5 0 1 2
7
1 1 0 4 0 2 2
5
1 2 0 1 2
1
0
10
4 2 1 3 0 4 1 4 2 0
5
0 4 3 0 0
8
5 6 0 2 4 1 1 0
3
0 0 0
5
0 0 0 0 0
5
2 0 1 3 1
5
0 1 0 3 0
1
0
3
0 1 0
3
0 0 0
2
0 0
6
2 0 1 1 0 1
10
1 1 1 0 5 4 0 5 6 0
6
0 2 3 2 1 0
4
0 2 0 0
5
2 0 4 3 0
3
0 2 1
2
0 0
7
2 2 1 4 0 4 0
4
3 2 0 1
4
2 0 0 0
4
1 2 1 0
6
0 3 0 0 2 5
5
1 0 1 4 0
4
0 0 3 0
4
0 3 1 1
7
0 4 1 0 2 0 3
7
0 2 0 0 0 2 3
4
1 3 0 0
7
4 1 2 5 0 0 1
5
0 3 1 0 0
4
0 2 0 3
5
1 0 2 0 0
9
2 0 3 6 5 0 1 6 0
8
0 2 2 2 1 3 0 0
9
1 4 0 1 6 0 6 3 0
1
0
5
1 0 2 1 2
1
0
`

func loadTests() ([]string, []string, error) {
	scan := bufio.NewScanner(strings.NewReader(testcasesBRaw))
	scan.Split(bufio.ScanWords)
	if !scan.Scan() {
		return nil, nil, fmt.Errorf("empty test file")
	}
	t, _ := strconv.Atoi(scan.Text())
	inputs := make([]string, t)
	expects := make([]string, t)
	for caseNum := 0; caseNum < t; caseNum++ {
		if !scan.Scan() {
			return nil, nil, fmt.Errorf("missing n for case %d", caseNum+1)
		}
		n, _ := strconv.Atoi(scan.Text())
		arr := make([]int, n)
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d\n", n))
		for i := 0; i < n; i++ {
			scan.Scan()
			arr[i], _ = strconv.Atoi(scan.Text())
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.Itoa(arr[i]))
		}
		sb.WriteByte('\n')
		inputs[caseNum] = sb.String()
		expects[caseNum] = expected(arr)
	}
	return inputs, expects, nil
}

func main() {
	if len(os.Args) == 3 && os.Args[1] == "--" {
		os.Args = append([]string{os.Args[0]}, os.Args[2])
	}
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	inputs, expects, err := loadTests()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load testcases: %v\n", err)
		os.Exit(1)
	}
	for i := range inputs {
		got, err := runCandidate(bin, inputs[i])
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != expects[i] {
			fmt.Fprintf(os.Stderr, "case %d failed\nexpected: %s\ngot: %s\ninput:\n%s", i+1, expects[i], got, inputs[i])
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(inputs))
}
