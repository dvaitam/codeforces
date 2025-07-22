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

type opB struct {
	op string
	x  int
}

func solveB(n int, ops []opB) []string {
	spf := make([]int, n+1)
	for i := 2; i <= n; i++ {
		if spf[i] == 0 {
			for j := i; j <= n; j += i {
				if spf[j] == 0 {
					spf[j] = i
				}
			}
		}
	}
	owner := make([]int, n+1)
	on := make([]bool, n+1)
	res := make([]string, len(ops))
	for idx, o := range ops {
		if o.op == "+" {
			if on[o.x] {
				res[idx] = "Already on"
				continue
			}
			conflict := false
			conflictWith := 0
			j := o.x
			for j > 1 {
				p := spf[j]
				if owner[p] != 0 {
					conflict = true
					conflictWith = owner[p]
					break
				}
				for j%p == 0 {
					j /= p
				}
			}
			if conflict {
				res[idx] = fmt.Sprintf("Conflict with %d", conflictWith)
			} else {
				res[idx] = "Success"
				on[o.x] = true
				j = o.x
				for j > 1 {
					p := spf[j]
					owner[p] = o.x
					for j%p == 0 {
						j /= p
					}
				}
			}
		} else {
			if !on[o.x] {
				res[idx] = "Already off"
				continue
			}
			j := o.x
			for j > 1 {
				p := spf[j]
				owner[p] = 0
				for j%p == 0 {
					j /= p
				}
			}
			on[o.x] = false
			res[idx] = "Success"
		}
	}
	return res
}

func genCaseB(rng *rand.Rand) (int, []opB) {
	n := rng.Intn(50) + 2
	m := rng.Intn(30) + 1
	ops := make([]opB, m)
	for i := 0; i < m; i++ {
		if rng.Intn(2) == 0 {
			ops[i].op = "+"
		} else {
			ops[i].op = "-"
		}
		ops[i].x = rng.Intn(n-1) + 1
	}
	return n, ops
}

func runCaseB(bin string, n int, ops []opB) error {
	var input bytes.Buffer
	fmt.Fprintf(&input, "%d %d\n", n, len(ops))
	for _, o := range ops {
		fmt.Fprintf(&input, "%s %d\n", o.op, o.x)
	}
	cmd := exec.Command(bin)
	cmd.Stdin = bytes.NewReader(input.Bytes())
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("runtime error: %v", err)
	}
	scanner := bufio.NewScanner(bytes.NewReader(out))
	exp := solveB(n, ops)
	for i, e := range exp {
		if !scanner.Scan() {
			return fmt.Errorf("missing output line %d", i+1)
		}
		got := strings.TrimSpace(scanner.Text())
		if got != e {
			return fmt.Errorf("line %d expected %q got %q", i+1, e, got)
		}
	}
	if scanner.Scan() {
		return fmt.Errorf("extra output: %s", scanner.Text())
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		return
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for t := 0; t < 100; t++ {
		n, ops := genCaseB(rng)
		if err := runCaseB(bin, n, ops); err != nil {
			fmt.Printf("case %d failed: %v\n", t+1, err)
			return
		}
	}
	fmt.Println("All tests passed")
}
