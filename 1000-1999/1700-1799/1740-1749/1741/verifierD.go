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

type Case struct {
	n    int
	perm []int
}

func genCases() []Case {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	cases := make([]Case, 100)
	for i := range cases {
		k := rng.Intn(4) + 1 // n=1..4 -> size up to 16
		n := 1 << uint(k)
		perm := rand.Perm(n)
		for j := range perm {
			perm[j]++
		}
		cases[i] = Case{n: n, perm: perm}
	}
	return cases
}

func solve(a []int, base int) int {
	n := len(a)
	if n == 1 {
		if a[0] == base {
			return 0
		}
		return -1
	}
	mid := n / 2
	left := a[:mid]
	right := a[mid:]
	check := func(seg []int, l, r int) bool {
		for _, v := range seg {
			if v < l || v > r {
				return false
			}
		}
		return true
	}
	best := -1
	if check(left, base, base+mid-1) && check(right, base+mid, base+n-1) {
		op1 := solve(left, base)
		if op1 != -1 {
			op2 := solve(right, base+mid)
			if op2 != -1 {
				best = op1 + op2
			}
		}
	}
	if check(left, base+mid, base+n-1) && check(right, base, base+mid-1) {
		op1 := solve(left, base+mid)
		if op1 != -1 {
			op2 := solve(right, base)
			if op2 != -1 {
				val := op1 + op2 + 1
				if best == -1 || val < best {
					best = val
				}
			}
		}
	}
	return best
}

func runCandidate(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func runCase(bin string, c Case) error {
	exp := solve(append([]int(nil), c.perm...), 1)
	var sb strings.Builder
	sb.WriteString("1\n")
	sb.WriteString(fmt.Sprintf("%d\n", c.n))
	for i, v := range c.perm {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(v))
	}
	sb.WriteByte('\n')
	out, err := runCandidate(bin, sb.String())
	if err != nil {
		return err
	}
	got, err := strconv.Atoi(strings.TrimSpace(out))
	if err != nil || got != exp {
		return fmt.Errorf("expected %d got %s", exp, out)
	}
	return nil
}

func main() {
	if len(os.Args) == 3 && os.Args[1] == "--" {
		os.Args = append([]string{os.Args[0]}, os.Args[2])
	}
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	cases := genCases()
	for i, c := range cases {
		if err := runCase(bin, c); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
