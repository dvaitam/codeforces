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

type testC struct {
	n      int
	c      int
	d      int
	beauty []int
	cost   []int
	typ    []byte
}

func solveC(tc testC) int {
	best := 0
	for i := 0; i < tc.n; i++ {
		for j := i + 1; j < tc.n; j++ {
			b := tc.beauty[i] + tc.beauty[j]
			t1 := tc.typ[i]
			t2 := tc.typ[j]
			if t1 == 'C' && t2 == 'C' {
				if tc.cost[i]+tc.cost[j] <= tc.c && b > best {
					best = b
				}
			} else if t1 == 'D' && t2 == 'D' {
				if tc.cost[i]+tc.cost[j] <= tc.d && b > best {
					best = b
				}
			} else if t1 == 'C' && t2 == 'D' {
				if tc.cost[i] <= tc.c && tc.cost[j] <= tc.d && b > best {
					best = b
				}
			} else if t1 == 'D' && t2 == 'C' {
				if tc.cost[i] <= tc.d && tc.cost[j] <= tc.c && b > best {
					best = b
				}
			}
		}
	}
	return best
}

func genC() (string, int) {
	n := rand.Intn(7) + 2
	c := rand.Intn(20)
	d := rand.Intn(20)
	beauty := make([]int, n)
	cost := make([]int, n)
	typ := make([]byte, n)
	for i := 0; i < n; i++ {
		beauty[i] = rand.Intn(30) + 1
		cost[i] = rand.Intn(20) + 1
		if rand.Intn(2) == 0 {
			typ[i] = 'C'
		} else {
			typ[i] = 'D'
		}
	}
	tc := testC{n, c, d, beauty, cost, typ}
	ans := solveC(tc)
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d %d\n", n, c, d))
	for i := 0; i < n; i++ {
		sb.WriteString(fmt.Sprintf("%d %d %c\n", beauty[i], cost[i], typ[i]))
	}
	input := sb.String()
	return input, ans
}

func runBinary(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var buf bytes.Buffer
	cmd.Stdout = &buf
	cmd.Stderr = &buf
	err := cmd.Run()
	return buf.String(), err
}

func main() {
	rand.Seed(time.Now().UnixNano())
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierC.go <binary>")
		return
	}
	bin := os.Args[1]
	for i := 0; i < 100; i++ {
		in, exp := genC()
		got, err := runBinary(bin, in)
		if err != nil {
			fmt.Printf("Test %d: runtime error: %v\nInput:\n%sOutput:\n%s", i+1, err, in, got)
			return
		}
		got = strings.TrimSpace(got)
		var val int
		if _, err := fmt.Sscan(got, &val); err != nil || val != exp {
			fmt.Printf("Test %d failed\nInput:\n%sExpected:\n%d\nGot:\n%s", i+1, in, exp, got)
			return
		}
	}
	fmt.Println("All tests passed")
}
