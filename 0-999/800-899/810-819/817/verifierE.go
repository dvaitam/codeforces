package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

type op struct {
	t int
	p int
	l int
}

func simulate(ops []op) []int {
	counts := map[int]int{}
	res := []int{}
	for _, op := range ops {
		switch op.t {
		case 1:
			counts[op.p]++
		case 2:
			counts[op.p]--
			if counts[op.p] <= 0 {
				delete(counts, op.p)
			}
		case 3:
			cnt := 0
			for v, c := range counts {
				if v^op.p < op.l {
					cnt += c
				}
			}
			res = append(res, cnt)
		}
	}
	return res
}

func runBin(bin, input string) (string, error) {
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
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func genCase(rng *rand.Rand) (string, string) {
	q := rng.Intn(20) + 1
	ops := make([]op, 0, q)
	counts := map[int]int{}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", q))
	for i := 0; i < q; i++ {
		t := rng.Intn(3) + 1
		if t == 2 && len(counts) == 0 {
			t = 1
		}
		if t == 1 {
			p := rng.Intn(100) + 1
			ops = append(ops, op{t: 1, p: p})
			counts[p]++
			sb.WriteString(fmt.Sprintf("1 %d\n", p))
		} else if t == 2 {
			var key int
			idx := rng.Intn(len(counts))
			j := 0
			for k := range counts {
				if j == idx {
					key = k
					break
				}
				j++
			}
			ops = append(ops, op{t: 2, p: key})
			counts[key]--
			if counts[key] <= 0 {
				delete(counts, key)
			}
			sb.WriteString(fmt.Sprintf("2 %d\n", key))
		} else {
			p := rng.Intn(100) + 1
			l := rng.Intn(100) + 1
			ops = append(ops, op{t: 3, p: p, l: l})
			sb.WriteString(fmt.Sprintf("3 %d %d\n", p, l))
		}
	}
	results := simulate(ops)
	var exp strings.Builder
	for i, v := range results {
		if i > 0 {
			exp.WriteByte('\n')
		}
		exp.WriteString(strconv.Itoa(v))
	}
	return sb.String(), exp.String()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		input, exp := genCase(rng)
		out, err := runBin(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: runtime error: %v\ninput:\n%s\n", i+1, err, input)
			os.Exit(1)
		}
		if strings.TrimSpace(out) != strings.TrimSpace(exp) {
			fmt.Fprintf(os.Stderr, "case %d failed:\ninput:\n%s\nexpected:\n%s\ngot:\n%s\n", i+1, input, exp, out)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
