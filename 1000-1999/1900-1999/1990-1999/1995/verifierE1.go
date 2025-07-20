package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"sort"
	"strings"
)

const INF int = 2000000001

func M(x, n int) int { return (x + 2*n) % (2 * n) }

func expectedE1(n int, v []int) int {
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
							if dp[c1][d1] != INF {
								val := r[M(2*d2-c1, n)] + r[M(2*d2+c2+1, n)]
								if val >= minn {
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

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE1.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	file, err := os.Open("testcasesE1.txt")
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open testcases: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	idx := 0
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		idx++
		fields := strings.Fields(line)
		if len(fields) < 1 {
			fmt.Printf("test %d: invalid line\n", idx)
			os.Exit(1)
		}
		var n int
		fmt.Sscan(fields[0], &n)
		if len(fields) != 1+2*n {
			fmt.Printf("test %d: invalid number of values\n", idx)
			os.Exit(1)
		}
		v := make([]int, 2*n)
		for i := 0; i < 2*n; i++ {
			fmt.Sscan(fields[1+i], &v[i])
		}
		expect := expectedE1(n, v)
		input := fmt.Sprintf("1\n%d\n", n)
		for i := 0; i < 2*n; i++ {
			input += fmt.Sprintf("%d ", v[i])
		}
		input += "\n"
		cmd := exec.Command(bin)
		cmd.Stdin = strings.NewReader(input)
		var out bytes.Buffer
		cmd.Stdout = &out
		var stderr bytes.Buffer
		cmd.Stderr = &stderr
		err := cmd.Run()
		if err != nil {
			fmt.Printf("test %d: runtime error: %v\nstderr: %s\n", idx, err, stderr.String())
			os.Exit(1)
		}
		res := strings.TrimSpace(out.String())
		var got int
		if _, err := fmt.Sscan(res, &got); err != nil {
			fmt.Printf("test %d: failed to parse output %q\n", idx, res)
			os.Exit(1)
		}
		if got != expect {
			fmt.Printf("test %d failed: expected %d got %d\n", idx, expect, got)
			os.Exit(1)
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "scanner error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("All %d tests passed\n", idx)
}
