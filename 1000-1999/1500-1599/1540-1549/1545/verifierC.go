package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

type matrix [][]int

func latinSquare(n int) matrix {
	m := make(matrix, n)
	for i := 0; i < n; i++ {
		m[i] = make([]int, n)
		for j := 0; j < n; j++ {
			m[i][j] = (i+j)%n + 1
		}
	}
	return m
}

func generateCase(rng *rand.Rand) string {
	n := rng.Intn(6) + 5 // 5..10
	base := latinSquare(n)
	used := make(map[string]bool)
	var extras matrix
	for len(extras) < n {
		idx := len(extras)
		arr := append([]int(nil), base[idx]...)
		p := rng.Intn(n)
		q := rng.Intn(n)
		for q == p {
			q = rng.Intn(n)
		}
		arr[p], arr[q] = arr[q], arr[p]
		key := fmt.Sprint(arr)
		if key == fmt.Sprint(base[idx]) || used[key] {
			continue
		}
		used[key] = true
		extras = append(extras, arr)
	}
	var sb strings.Builder
	sb.WriteString("1\n")
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			if j > 0 {
				sb.WriteByte(' ')
			}
			fmt.Fprintf(&sb, "%d", base[i][j])
		}
		sb.WriteByte('\n')
	}
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			if j > 0 {
				sb.WriteByte(' ')
			}
			fmt.Fprintf(&sb, "%d", extras[i][j])
		}
		sb.WriteByte('\n')
	}
	return sb.String()
}

func runCandidate(path, input string) (string, error) {
	cmd := exec.Command(path)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

// solve is an embedded reference solver for 1545C.
func solve(input string) string {
	const MOD = 998244353
	r := strings.NewReader(input)
	var t int
	fmt.Fscan(r, &t)

	var outBuf strings.Builder

	for tc := 0; tc < t; tc++ {
		var n int
		fmt.Fscan(r, &n)

		A := make([][]int, 2*n+1)
		for i := 1; i <= 2*n; i++ {
			A[i] = make([]int, n+1)
			for j := 1; j <= n; j++ {
				fmt.Fscan(r, &A[i][j])
			}
		}

		available := make([]bool, 2*n+1)
		for i := 1; i <= 2*n; i++ {
			available[i] = true
		}

		covered := make([][]bool, n+1)
		deg := make([][]int, n+1)
		arraysWith := make([][][]int, n+1)
		for i := 1; i <= n; i++ {
			covered[i] = make([]bool, n+1)
			deg[i] = make([]int, n+1)
			arraysWith[i] = make([][]int, n+1)
		}

		for i := 1; i <= 2*n; i++ {
			for c := 1; c <= n; c++ {
				v := A[i][c]
				deg[c][v]++
				arraysWith[c][v] = append(arraysWith[c][v], i)
			}
		}

		q := make([][2]int, 0)
		for c := 1; c <= n; c++ {
			for v := 1; v <= n; v++ {
				if deg[c][v] == 1 {
					q = append(q, [2]int{c, v})
				}
			}
		}

		chosen := make([]int, 0, n)

		for len(q) > 0 {
			item := q[0]
			q = q[1:]
			c, v := item[0], item[1]
			if covered[c][v] {
				continue
			}
			chosenArr := -1
			for _, arr := range arraysWith[c][v] {
				if available[arr] {
					chosenArr = arr
					break
				}
			}
			if chosenArr == -1 {
				continue
			}
			chosen = append(chosen, chosenArr)
			available[chosenArr] = false
			for c1 := 1; c1 <= n; c1++ {
				v1 := A[chosenArr][c1]
				if covered[c1][v1] {
					continue
				}
				covered[c1][v1] = true
				for _, B := range arraysWith[c1][v1] {
					if available[B] {
						available[B] = false
						for c2 := 1; c2 <= n; c2++ {
							v2 := A[B][c2]
							if !covered[c2][v2] {
								deg[c2][v2]--
								if deg[c2][v2] == 1 {
									q = append(q, [2]int{c2, v2})
								}
							}
						}
					}
				}
			}
		}

		ans := 1
		color := make([]int, 2*n+1)
		for j := 1; j <= 2*n; j++ {
			color[j] = -1
		}

		for i := 1; i <= 2*n; i++ {
			if available[i] && color[i] == -1 {
				ans = (ans * 2) % MOD
				qArr := []int{i}
				color[i] = 0
				for len(qArr) > 0 {
					curr := qArr[0]
					qArr = qArr[1:]
					if color[curr] == 0 {
						chosen = append(chosen, curr)
					}
					for c := 1; c <= n; c++ {
						v := A[curr][c]
						if !covered[c][v] {
							for _, neighbor := range arraysWith[c][v] {
								if available[neighbor] && color[neighbor] == -1 {
									color[neighbor] = 1 - color[curr]
									qArr = append(qArr, neighbor)
								}
							}
						}
					}
				}
			}
		}

		fmt.Fprintln(&outBuf, ans)
		for i, idx := range chosen {
			if i > 0 {
				fmt.Fprint(&outBuf, " ")
			}
			fmt.Fprint(&outBuf, idx)
		}
		fmt.Fprintln(&outBuf)
	}
	return strings.TrimSpace(outBuf.String())
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go <binary>")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		input := generateCase(rng)
		expect := solve(input)
		got, err := runCandidate(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate failed on test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		if got != expect {
			fmt.Printf("test %d failed\ninput:\n%sexpected:%s\nactual:%s\n", i+1, input, expect, got)
			os.Exit(1)
		}
	}
	fmt.Println("all tests passed")
}
