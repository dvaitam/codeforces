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

const (
	INF           int    = 2000000001
	embeddedCases string = `3 9 7 3 4 2 10
2 8 4 1 4
2 10 3 3 10
2 2 2 5 2
3 10 8 10 5 10 7
1 2 3
4 1 6 6 5 7 7 1 5
2 4 2 10 3
2 9 4 8 6
1 5 7
2 6 10 9 7
2 5 2 8 5
2 5 5 2 6
2 8 1 7 8
1 5 9
4 8 10 10 9 4 3 1 5
1 7 5
2 8 8 5 9
4 3 10 8 5 3 1 4 6
3 2 1 7 5 6 10
2 6 7 3 2
2 6 8 8 3
3 7 7 7 6 10 2
1 9 8
2 2 1 9 7
4 4 7 2 7 1 5 10 3
3 5 9 2 9 6 1
2 1 2 6 6
4 4 4 6 10 8 5 9 3
4 7 1 2 1 2 6 4 5
1 9 5
1 2 9
3 1 4 2 3 8 3
2 8 5 3 8
1 9 7
1 2 2
3 2 7 4 5 7 4
4 10 9 9 6 3 7 1 10
4 10 4 5 6 9 6 8 8
3 6 5 3 3 4 8
3 1 6 8 10 5 4
1 2 8
2 4 7 10 6
2 1 8 7 8
4 6 1 4 10 4 2 8 3
1 9 9
4 10 5 4 2 7 1 9 8
3 1 4 10 5 7 8
3 2 1 2 6 7 7
4 4 5 8 7 9 1 7 5
4 3 7 7 7 10 10 5 8
3 5 7 10 3 6 1
4 8 8 5 2 4 6 4 7
1 10 6
2 3 8 2 2
2 10 1 4 3
1 9 7
4 3 5 7 8 5 9 1 3
4 3 2 2 8 5 5 5 7
1 10 5
3 8 2 1 5 4 8
3 9 4 7 2 7 1
1 2 6
1 5 5
1 1 1
1 3 10
2 8 2 10 3
1 9 3
1 2 2
1 4 2
4 7 2 2 8 8 9 3 9
1 3 6
4 9 3 2 9 1 9 4 2
4 6 9 6 9 10 9 7 5
2 10 2 3 10
4 6 5 6 9 6 1 1 1
2 5 4 3 3
1 8 8
2 5 5 3 1
2 6 6 4 2
4 3 5 6 4 2 10 4 6
4 2 5 7 9 6 1 6 4
4 1 7 9 5 3 5 1 5
1 6 7
1 5 4
2 5 6 9 2
2 8 4 10 10
1 7 5
1 5 1
1 10 3
3 3 7 4 10 9 4
1 9 4
4 9 8 10 10 6 6 6 2
1 7 6
2 8 4 7 10
1 4 10
4 7 7 3 3 7 7 8 4
3 2 3 10 1 5 3
1 3 5
2 7 5 9 7`
)

type testCase struct {
	n      int
	values []int
}

func M(x, n int) int { return (x + 2*n) % (2 * n) }

// solveCase mirrors the original reference solution and returns the answer for one test.
func solveCase(n int, v []int) int {
	if n%2 == 0 {
		maxx := 0
		minn := INF
		for i := 0; i < n/2; i++ {
			s := []int{v[2*i] + v[2*i+1], v[2*i] + v[2*i+n+1], v[2*i+n] + v[2*i+n+1], v[2*i+n] + v[2*i+1]}
			sort.Ints(s)
			if s[2] > maxx {
				maxx = s[2]
			}
			if s[1] < minn {
				minn = s[1]
			}
		}
		return maxx - minn
	}
	if n == 1 {
		return 0
	}
	r := make([]int, 0, 2*n)
	cnt := 0
	for i := 0; i < n; i++ {
		r = append(r, v[cnt])
		cnt ^= 1
		r = append(r, v[cnt])
		cnt = M(cnt+n, n)
	}
	ans := INF
	for id := 0; id < n; id++ {
		for m1 := 0; m1 < 2; m1++ {
			for m2 := 0; m2 < 2; m2++ {
				minn := r[M(2*id-m1, n)] + r[M(2*id+m2+1, n)]
				dp := [2][]int{make([]int, n), make([]int, n)}
				for t := 0; t < 2; t++ {
					for j := 0; j < n; j++ {
						dp[t][j] = INF
					}
				}
				dp[m2][id] = minn
				for j := 1; j < n; j++ {
					d2 := (id + j) % n
					d1 := (id + j - 1) % n
					for c1 := 0; c1 < 2; c1++ {
						for c2 := 0; c2 < 2; c2++ {
							if dp[c1][d1] == INF {
								continue
							}
							val := r[M(2*d2-c1, n)] + r[M(2*d2+c2+1, n)]
							if val < minn {
								continue
							}
							cur := dp[c1][d1]
							if val > cur {
								cur = val
							}
							if cur < dp[c2][d2] {
								dp[c2][d2] = cur
							}
						}
					}
				}
				p := (id + n - 1) % n
				if dp[m1][p] != INF {
					diff := dp[m1][p] - minn
					if diff < ans {
						ans = diff
					}
				}
			}
		}
	}
	return ans
}

func loadTestcases() ([]testCase, error) {
	scanner := bufio.NewScanner(strings.NewReader(embeddedCases))
	idx := 0
	var cases []testCase
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		idx++
		fields := strings.Fields(line)
		n, err := strconv.Atoi(fields[0])
		if err != nil {
			return nil, fmt.Errorf("test %d: parse n: %w", idx, err)
		}
		if len(fields) != 1+2*n {
			return nil, fmt.Errorf("test %d: expected %d values, got %d", idx, 2*n, len(fields)-1)
		}
		vals := make([]int, 2*n)
		for i := 0; i < 2*n; i++ {
			vals[i], err = strconv.Atoi(fields[1+i])
			if err != nil {
				return nil, fmt.Errorf("test %d: parse value %d: %w", idx, i, err)
			}
		}
		cases = append(cases, testCase{n: n, values: vals})
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("scanner error: %w", err)
	}
	if len(cases) == 0 {
		return nil, fmt.Errorf("no embedded testcases")
	}
	return cases, nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE1.go /path/to/binary")
		os.Exit(1)
	}
	testcases, err := loadTestcases()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load testcases: %v\n", err)
		os.Exit(1)
	}
	bin := os.Args[1]
	for idx, tc := range testcases {
		expect := solveCase(tc.n, tc.values)
		var input strings.Builder
		input.WriteString("1\n")
		fmt.Fprintf(&input, "%d\n", tc.n)
		for _, val := range tc.values {
			fmt.Fprintf(&input, "%d ", val)
		}
		input.WriteByte('\n')
		cmd := exec.Command(bin)
		cmd.Stdin = strings.NewReader(input.String())
		var out bytes.Buffer
		cmd.Stdout = &out
		var stderr bytes.Buffer
		cmd.Stderr = &stderr
		err := cmd.Run()
		if err != nil {
			fmt.Printf("test %d: runtime error: %v\nstderr: %s\n", idx+1, err, stderr.String())
			os.Exit(1)
		}
		res := strings.TrimSpace(out.String())
		var got int
		if _, err := fmt.Sscan(res, &got); err != nil {
			fmt.Printf("test %d: failed to parse output %q\n", idx+1, res)
			os.Exit(1)
		}
		if got != expect {
			fmt.Printf("test %d failed: expected %d got %d\n", idx+1, expect, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(testcases))
}
