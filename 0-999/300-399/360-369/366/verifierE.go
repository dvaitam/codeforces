package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

func solve(input string) string {
	r := bufio.NewReader(strings.NewReader(input))
	var n, m, k, s int
	if _, err := fmt.Fscan(r, &n, &m, &k, &s); err != nil {
		return ""
	}
	inf := int(1e9)
	maxS := make([]int, k+1)
	minS := make([]int, k+1)
	maxD := make([]int, k+1)
	minD := make([]int, k+1)
	for t := 1; t <= k; t++ {
		maxS[t] = -inf
		minS[t] = inf
		maxD[t] = -inf
		minD[t] = inf
	}
	for i := 1; i <= n; i++ {
		for j := 1; j <= m; j++ {
			var t int
			fmt.Fscan(r, &t)
			S := i + j
			D := i - j
			if S > maxS[t] {
				maxS[t] = S
			}
			if S < minS[t] {
				minS[t] = S
			}
			if D > maxD[t] {
				maxD[t] = D
			}
			if D < minD[t] {
				minD[t] = D
			}
		}
	}
	Q := make([]int, s)
	for i := 0; i < s; i++ {
		fmt.Fscan(r, &Q[i])
	}
	dist := make([][]int, k+1)
	for t := 0; t <= k; t++ {
		dist[t] = make([]int, k+1)
	}
	for t1 := 1; t1 <= k; t1++ {
		for t2 := 1; t2 <= k; t2++ {
			var d int
			if t1 == t2 {
				d1 := maxS[t1] - minS[t1]
				d2 := maxD[t1] - minD[t1]
				if d2 > d1 {
					d1 = d2
				}
				d = d1
			} else {
				d1 := maxS[t1] - minS[t2]
				if maxS[t2]-minS[t1] > d1 {
					d1 = maxS[t2] - minS[t1]
				}
				d2 := maxD[t1] - minD[t2]
				if maxD[t2]-minD[t1] > d2 {
					d2 = maxD[t2] - minD[t1]
				}
				if d2 > d1 {
					d = d2
				} else {
					d = d1
				}
			}
			dist[t1][t2] = d
		}
	}
	ans := 0
	for i := 0; i+1 < s; i++ {
		t1 := Q[i]
		t2 := Q[i+1]
		if dist[t1][t2] > ans {
			ans = dist[t1][t2]
		}
	}
	return fmt.Sprintf("%d\n", ans)
}

func genTest() (string, string) {
	n := rand.Intn(4) + 1
	m := rand.Intn(4) + 1
	k := rand.Intn(3) + 1
	s := rand.Intn(8) + 2
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d %d %d\n", n, m, k, s)
	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			fmt.Fprintf(&sb, "%d ", rand.Intn(k)+1)
		}
		sb.WriteByte('\n')
	}
	for i := 0; i < s; i++ {
		fmt.Fprintf(&sb, "%d ", rand.Intn(k)+1)
	}
	sb.WriteByte('\n')
	inp := sb.String()
	out := solve(inp)
	return inp, out
}

func runBinary(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func main() {
	rand.Seed(time.Now().UnixNano())
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierE.go /path/to/binary")
		return
	}
	bin := os.Args[1]
	for i := 0; i < 100; i++ {
		in, exp := genTest()
		got, err := runBinary(bin, in)
		if err != nil {
			fmt.Printf("Test %d: runtime error: %v\nInput:\n%sOutput:\n%s\n", i+1, err, in, got)
			return
		}
		if got != strings.TrimSpace(exp) {
			fmt.Printf("Test %d failed\nInput:\n%sExpected:\n%sGot:\n%s\n", i+1, in, exp, got)
			return
		}
	}
	fmt.Println("All tests passed")
}
