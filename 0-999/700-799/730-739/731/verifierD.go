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

func solveD(n, c int, words [][]int) string {
	L, R := 0, c-1
	prev := words[0]
	ok := true
	for i := 1; i < n; i++ {
		curr := words[i]
		minl := len(prev)
		if len(curr) < minl {
			minl = len(curr)
		}
		k := 0
		for k < minl && prev[k] == curr[k] {
			k++
		}
		if k == minl {
			if len(prev) > len(curr) {
				ok = false
				break
			}
		} else {
			a, b := prev[k], curr[k]
			if a > b {
				Li := c - a + 1
				Ri := c - b
				if Li > R || Ri < L {
					ok = false
					break
				}
				if Li > L {
					L = Li
				}
				if Ri < R {
					R = Ri
				}
			}
		}
		prev = curr
	}
	if !ok || L > R {
		return "-1\n"
	}
	return fmt.Sprintf("%d\n", L)
}

func genCaseD(rng *rand.Rand) (string, string) {
	n := rng.Intn(5) + 1
	c := rng.Intn(8) + 2
	words := make([][]int, n)
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d\n", n, c)
	for i := 0; i < n; i++ {
		l := rng.Intn(4) + 1
		words[i] = make([]int, l)
		fmt.Fprintf(&sb, "%d", l)
		for j := 0; j < l; j++ {
			words[i][j] = rng.Intn(c) + 1
			fmt.Fprintf(&sb, " %d", words[i][j])
		}
		sb.WriteByte('\n')
	}
	return sb.String(), solveD(n, c, words)
}

func run(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

func main() {
	rand.Seed(time.Now().UnixNano())
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		return
	}
	bin := os.Args[1]
	for i := 0; i < 100; i++ {
		in, expect := genCaseD(rand.New(rand.NewSource(time.Now().UnixNano() + int64(i))))
		got, err := run(bin, in)
		if err != nil {
			fmt.Printf("test %d: runtime error: %v\ninput:\n%soutput:\n%s", i+1, err, in, got)
			return
		}
		if strings.TrimSpace(got) != strings.TrimSpace(expect) {
			fmt.Printf("test %d failed\ninput:\n%sexpected:\n%sbut got:\n%s", i+1, in, expect, got)
			return
		}
	}
	fmt.Println("All tests passed")
}
