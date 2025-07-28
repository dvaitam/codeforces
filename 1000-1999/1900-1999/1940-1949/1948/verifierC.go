package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
)

type Node struct{ r, c int }

func reachable(s1, s2 string) bool {
	n := len(s1)
	grid := [][]byte{[]byte(s1), []byte(s2)}
	visited := make([][]bool, 2)
	for i := 0; i < 2; i++ {
		visited[i] = make([]bool, n)
	}
	q := []Node{{0, 0}}
	visited[0][0] = true
	dirs := [][2]int{{0, 1}, {0, -1}, {1, 0}, {-1, 0}}
	for len(q) > 0 {
		cur := q[0]
		q = q[1:]
		if cur.r == 1 && cur.c == n-1 {
			return true
		}
		for _, d := range dirs {
			nr, nc := cur.r+d[0], cur.c+d[1]
			if nr < 0 || nr >= 2 || nc < 0 || nc >= n {
				continue
			}
			tr, tc := nr, nc
			if grid[nr][nc] == '>' {
				tc++
			} else {
				tc--
			}
			if tr < 0 || tr >= 2 || tc < 0 || tc >= n {
				continue
			}
			if !visited[tr][tc] {
				visited[tr][tc] = true
				q = append(q, Node{tr, tc})
			}
		}
	}
	return false
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		return
	}
	bin := os.Args[1]

	rand.Seed(42)
	const t = 100
	ns := make([]int, t)
	rows1 := make([]string, t)
	rows2 := make([]string, t)
	for i := 0; i < t; i++ {
		n := rand.Intn(19) + 2 // 2..20
		ns[i] = n
		b1 := make([]byte, n)
		b2 := make([]byte, n)
		for j := 0; j < n; j++ {
			if rand.Intn(2) == 0 {
				b1[j] = '<'
			} else {
				b1[j] = '>'
			}
			if rand.Intn(2) == 0 {
				b2[j] = '<'
			} else {
				b2[j] = '>'
			}
		}
		rows1[i] = string(b1)
		rows2[i] = string(b2)
	}

	var input bytes.Buffer
	fmt.Fprintf(&input, "%d\n", t)
	for i := 0; i < t; i++ {
		fmt.Fprintf(&input, "%d\n%s\n%s\n", ns[i], rows1[i], rows2[i])
	}

	cmd := exec.Command(bin)
	cmd.Stdin = bytes.NewReader(input.Bytes())
	var out bytes.Buffer
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		fmt.Println("failed to run binary:", err)
		os.Exit(1)
	}

	scanner := bufio.NewScanner(bytes.NewReader(out.Bytes()))
	for i := 0; i < t; i++ {
		if !scanner.Scan() {
			fmt.Printf("missing output for case %d\n", i+1)
			os.Exit(1)
		}
		got := scanner.Text()
		want := "NO"
		if reachable(rows1[i], rows2[i]) {
			want = "YES"
		}
		if got != want {
			fmt.Printf("case %d: expected %s, got %s\n", i+1, want, got)
			os.Exit(1)
		}
	}
	if scanner.Scan() {
		fmt.Println("warning: extra output detected")
	}
	fmt.Println("All tests passed!")
}
