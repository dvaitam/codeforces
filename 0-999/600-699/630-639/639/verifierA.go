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

type query struct{ typ, id int }

func runCandidate(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func solveA(n, k int, t []int, qs []query) []string {
	displayed := make([]query, 0, k)
	inDisp := make([]bool, n+1)
	res := make([]string, 0, len(qs))
	for _, q := range qs {
		if q.typ == 1 {
			if k == 0 {
				continue
			}
			if len(displayed) < k {
				displayed = append(displayed, query{q.id, t[q.id]})
				for i := len(displayed) - 1; i > 0 && displayed[i].id != 0 && displayed[i].typ > displayed[i-1].typ; i-- {
					displayed[i], displayed[i-1] = displayed[i-1], displayed[i]
				}
				inDisp[q.id] = true
			} else if t[q.id] > displayed[len(displayed)-1].typ {
				rem := displayed[len(displayed)-1].id
				inDisp[rem] = false
				displayed[len(displayed)-1] = query{q.id, t[q.id]}
				for i := len(displayed) - 1; i > 0 && displayed[i].typ > displayed[i-1].typ; i-- {
					displayed[i], displayed[i-1] = displayed[i-1], displayed[i]
				}
				inDisp[q.id] = true
			}
		} else {
			if inDisp[q.id] {
				res = append(res, "YES")
			} else {
				res = append(res, "NO")
			}
		}
	}
	return res
}

func generateCase(rng *rand.Rand) (string, []string) {
	n := rng.Intn(6) + 1
	k := rng.Intn(n + 1)
	q := rng.Intn(10) + 1
	t := make([]int, n+1)
	usedT := make(map[int]bool)
	for i := 1; i <= n; i++ {
		v := rng.Intn(100) + 1
		for usedT[v] {
			v = rng.Intn(100) + 1
		}
		usedT[v] = true
		t[i] = v
	}
	qs := make([]query, 0, q)
	used := make([]bool, n+1)
	haveType2 := false
	for len(qs) < q {
		if !haveType2 || rng.Intn(2) == 0 && len(qs) < q-1 && len(qs) < n {
			// type1
			id := rng.Intn(n) + 1
			if used[id] {
				continue
			}
			used[id] = true
			qs = append(qs, query{1, id})
		} else {
			id := rng.Intn(n) + 1
			qs = append(qs, query{2, id})
			haveType2 = true
		}
	}
	if !haveType2 {
		qs[len(qs)-1].typ = 2
	}
	input := fmt.Sprintf("%d %d %d\n", n, k, q)
	for i := 1; i <= n; i++ {
		if i > 1 {
			input += " "
		}
		input += fmt.Sprintf("%d", t[i])
	}
	input += "\n"
	for _, qu := range qs {
		input += fmt.Sprintf("%d %d\n", qu.typ, qu.id)
	}
	out := solveA(n, k, t, qs)
	return input, out
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		input, outLines := generateCase(rng)
		got, err := runCandidate(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
		gotLines := strings.Split(strings.TrimSpace(got), "\n")
		expLines := outLines
		if len(gotLines) != len(expLines) {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %d lines got %d\ninput:\n%s", i+1, len(expLines), len(gotLines), input)
			os.Exit(1)
		}
		for j := range expLines {
			if strings.TrimSpace(gotLines[j]) != expLines[j] {
				fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, expLines[j], gotLines[j], input)
				os.Exit(1)
			}
		}
	}
	fmt.Println("All tests passed")
}
