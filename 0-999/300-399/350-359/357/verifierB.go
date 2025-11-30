package main

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func solve(n, m int, dances [][3]int) string {
	a := make([]int, n+1)
	for i := 0; i < m; i++ {
		t := dances[i]
		used := [4]bool{}
		for j := 0; j < 3; j++ {
			used[a[t[j]]] = true
		}
		for j := 0; j < 3; j++ {
			if a[t[j]] == 0 {
				for k := 1; k <= 3; k++ {
					if !used[k] {
						a[t[j]] = k
						used[k] = true
						break
					}
				}
			}
		}
	}
	var sb strings.Builder
	for i := 1; i <= n; i++ {
		if i > 1 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(a[i]))
	}
	return sb.String()
}

var testcasesRaw = `6 7
1 3 4
4 3 6
3 5 2
5 2 3
2 1 3
5 6 2
3 1 5

8 6
8 5 1
6 4 3
4 5 8
8 7 5
5 1 8
1 8 6

6 1
5 4 3

4 6
1 4 3
2 1 4
1 4 2
4 1 2
3 4 1
3 4 1

7 5
4 1 5
7 4 3
5 2 3
2 7 6
1 5 3

6 2
1 2 5
1 6 4

8 5
4 7 2
7 5 3
8 4 6
6 1 3
2 4 5

8 6
4 2 1
5 1 6
4 3 2
6 4 1
2 7 8
4 1 5

8 2
1 8 6
4 5 7

3 7
1 2 3
1 3 2
1 3 2
1 3 2
3 2 1
2 1 3
1 2 3

6 3
1 5 4
1 5 6
6 4 2

5 6
4 2 3
2 1 3
2 5 4
5 3 1
5 4 3
2 1 5

8 7
5 6 3
7 8 6
5 2 8
1 4 6
2 3 6
1 5 3
3 2 4

5 5
3 2 5
3 4 5
4 1 5
5 2 3
3 2 1

4 8
4 3 2
1 2 3
1 4 2
1 2 4
4 3 2
1 4 2
3 2 4
1 2 4

7 2
7 6 2
1 4 6

5 1
2 1 3

3 2
1 3 2
2 3 1

3 8
2 1 3
2 3 1
2 1 3
1 3 2
1 3 2
1 2 3
2 1 3
3 2 1

5 3
2 4 3
3 1 4
2 3 5

5 6
1 3 4
1 5 2
2 5 3
3 5 2
5 2 4
1 4 3

4 1
3 1 4

5 7
3 5 2
1 5 3
4 5 2
3 1 2
1 4 2
1 3 2
2 5 3

7 7
7 6 1
1 7 2
6 2 1
4 1 6
4 5 6
3 4 6
7 5 2

6 2
3 2 6
5 2 4

4 6
1 4 3
4 3 1
1 2 3
3 1 4
2 3 1
1 4 2

6 6
5 2 1
5 4 2
5 4 6
5 4 3
4 6 2
5 6 2

3 6
3 2 1
1 3 2
3 1 2
3 2 1
1 3 2
1 3 2

3 5
1 2 3
1 3 2
2 3 1
3 1 2
3 2 1

5 7
4 5 3
5 2 1
1 2 4
5 3 2
1 4 2
3 4 2
4 1 5

8 3
4 3 6
6 1 7
8 4 2

6 2
6 2 3
4 1 6

6 8
1 6 4
2 1 5
5 6 2
6 3 1
6 5 3
2 4 5
1 6 4
5 3 1

8 5
3 4 8
2 5 8
1 7 4
8 3 2
8 3 6

3 1
1 2 3

3 4
1 3 2
3 2 1
1 2 3
1 3 2

3 2
2 1 3
2 3 1

3 8
2 3 1
2 1 3
2 1 3
2 1 3
1 2 3
3 1 2
1 3 2
3 1 2

5 1
5 2 1

4 8
1 2 3
3 1 4
2 4 3
4 2 3
3 1 4
1 3 2
2 4 1
2 1 4

7 4
2 7 3
3 4 1
7 2 5
1 4 7

8 2
3 4 8
7 6 2

7 7
3 6 7
1 2 4
6 3 5
1 4 6
1 4 3
4 2 3
3 4 1

4 2
3 1 4
4 2 1

6 7
2 6 4
3 5 2
3 4 1
4 2 3
1 4 5
1 2 3
1 3 2

6 8
1 4 2
5 4 3
6 1 5
3 1 5
3 4 5
4 1 3
3 6 2
5 1 6

3 5
2 3 1
3 1 2
1 2 3
2 1 3
3 1 2

8 7
5 3 4
6 5 8
3 2 1
2 4 7
8 2 5
5 3 6
8 6 4

4 8
4 3 2
4 3 1
4 2 1
4 1 3
1 2 3
4 1 3
1 3 2
1 2 4

3 1
2 3 1

4 5
2 1 4
4 3 2
4 1 2
4 3 2
2 4 1

7 6
2 7 6
4 3 7
4 7 6
6 2 7
4 2 5
6 1 4

3 4
3 1 2
2 1 3
1 2 3
3 2 1

6 8
1 5 4
5 4 6
3 4 2
3 6 1
1 6 2
3 1 6
6 1 2
1 4 2

7 7
5 2 4
2 3 5
1 5 7
7 3 6
5 4 3
3 1 5
1 2 3

3 4
3 2 1
1 2 3
2 3 1
2 3 1

3 5
3 1 2
2 3 1
1 2 3
2 3 1
1 3 2

3 3
2 1 3
3 1 2
3 2 1

4 4
4 1 3
3 2 4
4 1 2
2 3 4

6 8
5 3 6
2 1 5
5 1 2
6 4 2
2 3 1
6 5 4
1 6 4
6 3 1

8 6
4 6 7
5 2 6
4 1 5
5 6 3
4 3 6
8 3 5

4 3
1 3 2
3 4 1
2 4 1

4 2
2 1 4
2 1 3

4 7
3 4 1
4 1 3
3 2 4
3 1 4
2 3 4
4 3 2
1 2 4

7 7
6 2 3
7 3 4
1 2 7
3 6 2
3 4 2
3 7 2
3 1 7

7 1
6 7 2

5 1
4 1 5

4 1
1 4 2

3 1
2 3 1

4 8
4 1 2
3 1 2
4 2 1
4 2 3
1 3 4
4 2 1
1 3 4
2 1 4

4 3
2 1 3
2 3 1
4 3 1

7 8
6 2 5
5 1 3
3 6 4
1 6 7
4 1 3
3 6 2
4 2 5
3 2 4

5 5
4 5 1
3 5 4
4 1 3
5 3 4
1 4 2

8 2
7 3 4
1 8 3

8 1
5 6 7

7 5
6 2 5
5 3 4
7 3 2
1 5 3
7 2 5

8 6
3 7 2
6 7 8
1 5 8
3 8 7
5 6 3
6 4 7

7 7
2 1 7
5 1 4
2 3 1
6 4 3
6 5 2
1 5 4
3 2 5

4 2
2 3 1
4 2 3

5 5
1 4 3
3 5 4
5 3 2
2 4 1
1 2 3

8 7
6 4 1
7 6 5
4 1 3
3 6 2
6 8 7
8 1 6
4 5 2

5 3
4 1 3
4 2 3
4 1 5

4 8
1 2 3
3 2 4
3 4 2
1 2 4
4 1 2
2 4 3
4 3 2
1 3 2

5 3
5 3 2
1 4 2
5 2 1

3 6
2 1 3
1 3 2
1 2 3
2 3 1
3 2 1
3 1 2

6 5
5 1 2
2 6 4
6 1 5
5 4 2
3 5 1

3 3
2 1 3
2 3 1
3 2 1

6 3
1 5 2
6 5 1
3 1 2

8 8
4 7 8
3 6 8
3 4 5
1 5 8
6 1 7
2 4 5
6 5 4
6 7 4

7 6
7 6 1
2 3 1
2 6 7
1 3 5
2 1 4
6 1 3

6 1
5 2 3

3 7
1 2 3
3 2 1
3 2 1
1 3 2
2 1 3
1 2 3
1 3 2

3 3
2 1 3
2 3 1
1 3 2

5 1
5 1 4

5 6
1 4 5
2 3 1
1 4 2
2 5 1
2 5 3
2 3 1

6 7
1 2 5
6 3 5
2 6 3
2 6 4
4 3 2
4 5 1
2 5 1

8 7
3 2 1
1 7 3
1 8 7
2 5 7
3 7 2
7 8 1
7 4 5

8 6
4 2 3
4 5 8
2 8 4
6 1 4
7 2 4
7 2 4

4 8
4 3 1
3 2 4
1 3 2
2 4 1
3 2 1
1 4 2
1 2 3
4 1 2

7 8
2 7 3
4 7 2
5 3 7
1 3 6
2 6 3
7 4 5
7 1 2
3 2 5`

func parseTestcases() ([]struct {
	n int
	m int
	d [][3]int
}, error) {
	reader := strings.NewReader(testcasesRaw)
	var res []struct {
		n int
		m int
		d [][3]int
	}
	for {
		var n, m int
		if _, err := fmt.Fscan(reader, &n, &m); err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		d := make([][3]int, m)
		for i := 0; i < m; i++ {
			fmt.Fscan(reader, &d[i][0], &d[i][1], &d[i][2])
		}
		res = append(res, struct {
			n int
			m int
			d [][3]int
		}{n: n, m: m, d: d})
	}
	return res, nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	cases, err := parseTestcases()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse embedded testcases: %v\n", err)
		os.Exit(1)
	}

	for idx, tc := range cases {
		expected := solve(tc.n, tc.m, tc.d)

		var sb strings.Builder
		fmt.Fprintf(&sb, "%d %d\n", tc.n, tc.m)
		for _, t := range tc.d {
			fmt.Fprintf(&sb, "%d %d %d\n", t[0], t[1], t[2])
		}
		sb.WriteByte('\n')

		cmd := exec.Command(bin)
		cmd.Stdin = strings.NewReader(sb.String())
		out, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Fprintf(os.Stderr, "test %d: runtime error: %v\n%s", idx+1, err, string(out))
			os.Exit(1)
		}

		got := strings.TrimSpace(string(out))
		if got != expected {
			fmt.Fprintf(os.Stderr, "test %d failed\nexpected: %s\n got: %s\n", idx+1, expected, got)
			os.Exit(1)
		}
	}

	fmt.Printf("All %d tests passed\n", len(cases))
}
