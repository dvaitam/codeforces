package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

const testcases = `6
5 4 2 1 6 3

7
4 5 6 3 7 2 1

5
2 4 3 1 5

2
2 1

5
5 1 4 2 3

8
7 3 1 4 6 5 8 2

2
2 1

3
2 3 1

8
4 6 1 5 2 7 8 3

5
4 5 2 1 3

3
3 2 1

2
1 2

4
4 2 3 1

2
1 2

4
1 4 2 3

3
1 3 2

7
4 7 5 6 3 1 2

8
6 8 2 3 4 5 1 7

5
4 5 3 2 1

2
2 1

8
4 1 2 6 7 5 3 8

6
4 2 1 5 3 6

2
2 1

4
4 2 1 3

4
4 3 2 1

4
4 3 1 2

7
5 1 7 2 4 3 6

2
1 2

2
2 1

7
1 2 7 4 6 3 5

3
2 3 1

8
1 4 7 5 6 2 3 8

2
2 1

5
5 4 1 2 3

4
3 1 4 2

6
1 5 2 3 6 4

4
3 1 2 4

2
1 2

8
5 7 1 2 8 4 6 3

4
4 2 1 3

7
7 4 2 3 1 6 5

8
2 1 4 5 7 6 3 8

6
6 5 2 3 4 1

5
1 4 2 5 3

7
5 2 1 7 4 6 3

4
2 1 4 3

6
1 4 5 3 6 2

7
7 5 3 2 4 6 1

4
2 1 4 3

3
3 2 1

6
4 1 3 5 2 6

7
5 3 7 2 4 1 6

3
3 2 1

3
3 1 2

6
6 3 2 5 4 1

7
3 4 6 2 7 1 5

2
1 2

5
5 2 1 4 3

2
1 2

2
1 2

6
4 1 2 6 3 5

5
5 1 2 3 4

3
1 2 3

6
6 1 3 5 2 4

4
3 4 2 1

4
1 3 4 2

5
1 3 5 2 4

5
3 2 1 5 4

7
5 2 1 6 7 4 3

4
3 1 2 4

5
4 2 3 5 1

3
3 1 2

8
5 2 8 6 1 4 7 3

3
3 1 2

3
1 3 2

8
1 2 5 4 7 8 6 3

6
2 3 4 1 5 6

4
1 2 4 3

6
4 6 2 1 5 3

4
4 3 1 2

4
2 1 3 4

7
7 3 6 4 1 2 5

8
6 1 5 3 2 8 4 7

4
4 2 3 1

6
2 5 1 6 4 3

5
3 1 2 4 5

5
2 3 4 1 5

4
2 3 4 1

4
3 2 1 4

7
2 3 7 5 6 1 4

7
5 6 1 7 4 3 2

5
5 4 3 1 2

2
2 1

7
2 7 5 3 1 6 4

7
3 7 6 1 2 5 4

4
1 3 2 4

7
4 5 6 7 2 1 3

6
3 1 4 2 6 5

6
6 1 4 5 3 2

8
2 6 1 5 8 3 7 4`

func run(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func good(p []int) int {
	n := len(p)
	prefixMax := make([]int, n)
	maxVal := 0
	for i := 0; i < n; i++ {
		if p[i] > maxVal {
			maxVal = p[i]
		}
		prefixMax[i] = maxVal
	}
	suffixMin := make([]int, n)
	minVal := n + 1
	for i := n - 1; i >= 0; i-- {
		if p[i] < minVal {
			minVal = p[i]
		}
		suffixMin[i] = minVal
	}
	cnt := 0
	for i := 0; i < n; i++ {
		leftMax := 0
		if i > 0 {
			leftMax = prefixMax[i-1]
		}
		rightMin := n + 1
		if i+1 < n {
			rightMin = suffixMin[i+1]
		}
		if leftMax < p[i] && p[i] < rightMin {
			cnt++
		}
	}
	return cnt
}

func brute(p []int) int {
	n := len(p)
	best := 0
	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			p[i], p[j] = p[j], p[i]
			if g := good(p); g > best {
				best = g
			}
			p[i], p[j] = p[j], p[i]
		}
	}
	return best
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	blocks := strings.Split(testcases, "\n\n")
	count := 0
	for idx, block := range blocks {
		block = strings.TrimSpace(block)
		if block == "" {
			continue
		}
		lines := strings.Split(block, "\n")
		if len(lines) < 2 {
			fmt.Fprintf(os.Stderr, "invalid test block at %d\n", idx+1)
			os.Exit(1)
		}
		n, _ := strconv.Atoi(strings.TrimSpace(lines[0]))
		fields := strings.Fields(lines[1])
		if len(fields) != n {
			fmt.Fprintf(os.Stderr, "case %d invalid permutation length %d expected %d\n", idx+1, len(fields), n)
			os.Exit(1)
		}
		p := make([]int, n)
		for i, f := range fields {
			v, err := strconv.Atoi(f)
			if err != nil {
				fmt.Fprintf(os.Stderr, "case %d invalid number: %v\n", idx+1, err)
				os.Exit(1)
			}
			p[i] = v
		}
		count++
		want := brute(append([]int(nil), p...))
		var sb strings.Builder
		sb.WriteString("1\n")
		sb.WriteString(fmt.Sprintf("%d\n", n))
		for i, v := range p {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.Itoa(v))
		}
		sb.WriteByte('\n')
		got, err := run(bin, sb.String())
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
		gotVal, err := strconv.Atoi(strings.TrimSpace(got))
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d output parse error: %v\n", idx+1, err)
			os.Exit(1)
		}
		if gotVal != want {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %d got %d\n", idx+1, want, gotVal)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", count)
}
