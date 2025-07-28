package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

func check(start int, row0, row1 string) bool {
	n := len(row0)
	visited := [2][]bool{make([]bool, n), make([]bool, n)}
	r := start
	c := 0
	for {
		if r == 0 {
			if row0[c] != 'B' || visited[0][c] {
				return false
			}
			visited[0][c] = true
		} else {
			if row1[c] != 'B' || visited[1][c] {
				return false
			}
			visited[1][c] = true
		}
		other := 1 - r
		if other == 0 {
			if row0[c] == 'B' && !visited[0][c] {
				r = other
				continue
			}
		} else {
			if row1[c] == 'B' && !visited[1][c] {
				r = other
				continue
			}
		}
		if c+1 < n {
			if r == 0 && row0[c+1] == 'B' {
				c++
				continue
			}
			if r == 1 && row1[c+1] == 'B' {
				c++
				continue
			}
		}
		break
	}
	for i := 0; i < n; i++ {
		if row0[i] == 'B' && !visited[0][i] {
			return false
		}
		if row1[i] == 'B' && !visited[1][i] {
			return false
		}
	}
	return true
}

func solveCase(row0, row1 string) bool {
	starts := []int{}
	if row0[0] == 'B' {
		starts = append(starts, 0)
	}
	if row1[0] == 'B' {
		starts = append(starts, 1)
	}
	for _, st := range starts {
		if check(st, row0, row1) {
			return true
		}
	}
	return false
}

func solveC(input string) string {
	reader := bufio.NewReader(strings.NewReader(input))
	var t int
	fmt.Fscan(reader, &t)
	var out strings.Builder
	for ; t > 0; t-- {
		var m int
		fmt.Fscan(reader, &m)
		var row0, row1 string
		fmt.Fscan(reader, &row0)
		fmt.Fscan(reader, &row1)
		if solveCase(row0, row1) {
			fmt.Fprintln(&out, "YES")
		} else {
			fmt.Fprintln(&out, "NO")
		}
	}
	return strings.TrimSpace(out.String())
}

func genTestC(rng *rand.Rand) string {
	t := rng.Intn(10) + 1
	var buf strings.Builder
	fmt.Fprintf(&buf, "%d\n", t)
	for ; t > 0; t-- {
		m := rng.Intn(20) + 1
		fmt.Fprintf(&buf, "%d\n", m)
		var r0, r1 strings.Builder
		for i := 0; i < m; i++ {
			if rng.Intn(2) == 0 {
				r0.WriteByte('B')
			} else {
				r0.WriteByte('W')
			}
			if rng.Intn(2) == 0 {
				r1.WriteByte('B')
			} else {
				r1.WriteByte('W')
			}
			// ensure at least one B per column
			if r0.Bytes()[i] != 'B' && r1.Bytes()[i] != 'B' {
				if rng.Intn(2) == 0 {
					r0.Bytes()[i] = 'B'
				} else {
					r1.Bytes()[i] = 'B'
				}
			}
		}
		fmt.Fprintf(&buf, "%s\n%s\n", r0.String(), r1.String())
	}
	return buf.String()
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
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	exe := os.Args[1]
	rng := rand.New(rand.NewSource(1))
	for i := 1; i <= 100; i++ {
		in := genTestC(rng)
		expect := solveC(in)
		got, err := run(exe, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: %v\n", i, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != expect {
			fmt.Printf("case %d failed\ninput:\n%sexpected:\n%s\ngot:\n%s\n", i, in, expect, got)
			os.Exit(1)
		}
	}
	fmt.Println("all tests passed")
}
