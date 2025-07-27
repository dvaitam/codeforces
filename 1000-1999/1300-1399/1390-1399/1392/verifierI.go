package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

func score(n, m int, a, b []int, x int) int {
	grid := make([][]int, n)
	for i := 0; i < n; i++ {
		grid[i] = make([]int, m)
		for j := 0; j < m; j++ {
			grid[i][j] = a[i] + b[j]
		}
	}
	visited := make([][]bool, n)
	for i := range visited {
		visited[i] = make([]bool, m)
	}
	goodBig, badBig, goodSmall, badSmall := 0, 0, 0, 0
	dirs := [][2]int{{1, 0}, {-1, 0}, {0, 1}, {0, -1}}
	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			if visited[i][j] {
				continue
			}
			t := grid[i][j] >= x
			queue := [][2]int{{i, j}}
			visited[i][j] = true
			border := i == 0 || j == 0 || i == n-1 || j == m-1
			for q := 0; q < len(queue); q++ {
				cx, cy := queue[q][0], queue[q][1]
				for _, d := range dirs {
					nx, ny := cx+d[0], cy+d[1]
					if nx < 0 || ny < 0 || nx >= n || ny >= m || visited[nx][ny] {
						continue
					}
					if (grid[nx][ny] >= x) == t {
						visited[nx][ny] = true
						if nx == 0 || ny == 0 || nx == n-1 || ny == m-1 {
							border = true
						}
						queue = append(queue, [2]int{nx, ny})
					}
				}
			}
			if t {
				if border {
					goodBig++
				} else {
					badBig += 2
				}
			} else {
				if border {
					goodSmall++
				} else {
					badSmall += 2
				}
			}
		}
	}
	return goodBig + badBig - (goodSmall + badSmall)
}

func generateTest() (string, string) {
	n := rand.Intn(3) + 2
	m := rand.Intn(3) + 2
	q := 3
	a := make([]int, n)
	bArr := make([]int, m)
	for i := range a {
		a[i] = rand.Intn(10)
	}
	for i := range bArr {
		bArr[i] = rand.Intn(10)
	}
	var in strings.Builder
	fmt.Fprintf(&in, "%d %d %d\n", n, m, q)
	for i, v := range a {
		if i > 0 {
			in.WriteByte(' ')
		}
		fmt.Fprintf(&in, "%d", v)
	}
	in.WriteByte('\n')
	for i, v := range bArr {
		if i > 0 {
			in.WriteByte(' ')
		}
		fmt.Fprintf(&in, "%d", v)
	}
	in.WriteByte('\n')
	var out strings.Builder
	for i := 0; i < q; i++ {
		x := rand.Intn(20)
		fmt.Fprintf(&in, "%d\n", x)
		fmt.Fprintf(&out, "%d\n", score(n, m, a, bArr, x))
	}
	return in.String(), out.String()
}

func referenceIO(t int) (string, string) {
	var in strings.Builder
	var out strings.Builder
	for i := 0; i < t; i++ {
		ti, to := generateTest()
		in.WriteString(ti)
		out.WriteString(to)
	}
	return in.String(), out.String()
}

func runBinary(path, input string) (string, error) {
	cmd := exec.Command(path)
	cmd.Stdin = strings.NewReader(input)
	var buf bytes.Buffer
	cmd.Stdout = &buf
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	return buf.String(), err
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierI.go /path/to/binary")
		return
	}
	rand.Seed(9)
	in, exp := referenceIO(100)
	out, err := runBinary(os.Args[1], in)
	if err != nil {
		fmt.Println("Runtime error:", err)
		os.Exit(1)
	}
	if strings.TrimSpace(out) != strings.TrimSpace(exp) {
		fmt.Println("Wrong Answer")
		fmt.Println("Expected:\n" + exp)
		fmt.Println("Got:\n" + out)
		os.Exit(1)
	}
	fmt.Println("All tests passed")
}
