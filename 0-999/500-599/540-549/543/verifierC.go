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

func solveC(n, m int, s []string, a [][]int) int {
	cols := m * 26
	cost := make([][]int, n)
	for i := 0; i < n; i++ {
		cost[i] = make([]int, cols)
	}
	for j := 0; j < m; j++ {
		for c := 0; c < 26; c++ {
			fid := j*26 + c
			sum := 0
			cc := byte('a' + c)
			for i := 0; i < n; i++ {
				if s[i][j] == cc {
					sum += a[i][j]
				}
			}
			for i := 0; i < n; i++ {
				cur := sum
				if s[i][j] == cc {
					cur -= a[i][j]
				} else {
					cur += a[i][j]
				}
				cost[i][fid] = cur
			}
		}
	}
	const INF = int(1 << 60)
	u := make([]int, n+1)
	v := make([]int, cols+1)
	p := make([]int, cols+1)
	way := make([]int, cols+1)
	for i := 1; i <= n; i++ {
		p[0] = i
		j0 := 0
		minv := make([]int, cols+1)
		used := make([]bool, cols+1)
		for j := 0; j <= cols; j++ {
			minv[j] = INF
		}
		for {
			used[j0] = true
			i0 := p[j0]
			delta := INF
			j1 := 0
			for j := 1; j <= cols; j++ {
				if !used[j] {
					cur := cost[i0-1][j-1] - u[i0] - v[j]
					if cur < minv[j] {
						minv[j] = cur
						way[j] = j0
					}
					if minv[j] < delta {
						delta = minv[j]
						j1 = j
					}
				}
			}
			for j := 0; j <= cols; j++ {
				if used[j] {
					u[p[j]] += delta
					v[j] -= delta
				} else {
					minv[j] -= delta
				}
			}
			j0 = j1
			if p[j0] == 0 {
				break
			}
		}
		for {
			j1 := way[j0]
			p[j0] = p[j1]
			j0 = j1
			if j0 == 0 {
				break
			}
		}
	}
	assignment := make([]int, n)
	for j := 1; j <= cols; j++ {
		if p[j] > 0 {
			assignment[p[j]-1] = j - 1
		}
	}
	ans := 0
	for i := 0; i < n; i++ {
		ans += cost[i][assignment[i]]
	}
	return ans
}

func genRandomString(m int) string {
	b := make([]byte, m)
	for i := 0; i < m; i++ {
		b[i] = byte('a' + rand.Intn(3))
	}
	return string(b)
}

func genTestC() (string, int) {
	n := rand.Intn(3) + 1
	m := rand.Intn(3) + 1
	s := make([]string, n)
	for i := 0; i < n; i++ {
		s[i] = genRandomString(m)
	}
	a := make([][]int, n)
	for i := 0; i < n; i++ {
		a[i] = make([]int, m)
		for j := 0; j < m; j++ {
			a[i][j] = rand.Intn(10)
		}
	}
	input := fmt.Sprintf("%d %d\n", n, m)
	for i := 0; i < n; i++ {
		input += s[i] + "\n"
	}
	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			if j > 0 {
				input += " "
			}
			input += fmt.Sprintf("%d", a[i][j])
		}
		input += "\n"
	}
	expected := solveC(n, m, s, a)
	return input, expected
}

func runBinary(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run verifierC.go <binary>")
		os.Exit(1)
	}
	bin := os.Args[1]
	rand.Seed(time.Now().UnixNano())
	for t := 1; t <= 100; t++ {
		input, expected := genTestC()
		output, err := runBinary(bin, input)
		if err != nil {
			fmt.Printf("Test %d: runtime error: %v\n", t, err)
			os.Exit(1)
		}
		if strings.TrimSpace(output) != fmt.Sprintf("%d", expected) {
			fmt.Printf("Test %d failed\nInput:\n%sExpected: %d\nGot: %s\n", t, input, expected, output)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed!")
}
