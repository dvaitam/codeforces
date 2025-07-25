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

type queryB struct {
	t   int
	d   int
	val int
}

func solveB(n, k, a, b int, queries []queryB) []int64 {
	cnt := make([]int, n+1)
	var res []int64
	for _, q := range queries {
		if q.t == 1 {
			d := q.d
			add := q.val
			cnt[d] += add
		} else {
			p := q.d
			var before, after int64
			for i := 1; i < p; i++ {
				if cnt[i] < b {
					before += int64(cnt[i])
				} else {
					before += int64(b)
				}
			}
			for i := p + k; i <= n; i++ {
				if cnt[i] < a {
					after += int64(cnt[i])
				} else {
					after += int64(a)
				}
			}
			res = append(res, before+after)
		}
	}
	return res
}

func generateB(rng *rand.Rand) (string, []int64) {
	n := rng.Intn(20) + 1
	k := rng.Intn(n) + 1
	a := rng.Intn(5) + 2
	b := rng.Intn(a-1) + 1
	q := rng.Intn(20) + 1
	queries := make([]queryB, q)
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d %d %d %d\n", n, k, a, b, q)
	for i := 0; i < q; i++ {
		if rng.Intn(2) == 0 {
			d := rng.Intn(n) + 1
			add := rng.Intn(5) + 1
			queries[i] = queryB{1, d, add}
			fmt.Fprintf(&sb, "1 %d %d\n", d, add)
		} else {
			p := rng.Intn(n) + 1
			queries[i] = queryB{2, p, 0}
			fmt.Fprintf(&sb, "2 %d\n", p)
		}
	}
	expected := solveB(n, k, a, b, queries)
	return sb.String(), expected
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierB.go /path/to/binary")
		return
	}
	exe := os.Args[1]
	rng := rand.New(rand.NewSource(43))
	for caseNum := 0; caseNum < 100; caseNum++ {
		input, expected := generateB(rng)
		cmd := exec.Command(exe)
		cmd.Stdin = strings.NewReader(input)
		var out bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = &out
		if err := cmd.Run(); err != nil {
			fmt.Printf("case %d runtime error: %v\n%s", caseNum+1, err, out.String())
			return
		}
		scanner := bufio.NewScanner(bytes.NewReader(out.Bytes()))
		scanner.Split(bufio.ScanWords)
		for i, exp := range expected {
			if !scanner.Scan() {
				fmt.Printf("case %d output ended early\ninput:\n%s", caseNum+1, input)
				return
			}
			got := scanner.Text()
			if got != fmt.Sprint(exp) {
				fmt.Printf("case %d query %d expected %d got %s\ninput:\n%s", caseNum+1, i+1, exp, got, input)
				return
			}
		}
		if scanner.Scan() {
			fmt.Printf("case %d extra output\n", caseNum+1)
			return
		}
	}
	fmt.Println("All tests passed")
}
