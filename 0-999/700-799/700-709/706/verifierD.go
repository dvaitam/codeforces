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

type multiset map[int]int

func (m multiset) add(x int) {
	m[x]++
}

func (m multiset) remove(x int) {
	m[x]--
	if m[x] == 0 {
		delete(m, x)
	}
}

func (m multiset) maxXor(x int) int {
	best := 0
	for v := range m {
		if cur := x ^ v; cur > best {
			best = cur
		}
	}
	return best
}

func generateCase(rng *rand.Rand) (string, []int) {
	q := rng.Intn(50) + 1
	ms := make(multiset)
	ms.add(0) // 0 is always present; xi >= 1 per problem spec so it's never removed

	type op struct {
		kind byte
		x    int
	}
	ops := make([]op, q)
	added := []int{} // tracks user-added elements available for removal
	hasQuery := false

	for i := 0; i < q; i++ {
		// Force a query on the last op if none generated yet
		if !hasQuery && i == q-1 {
			x := rng.Intn(1000) + 1
			ops[i] = op{'?', x}
			hasQuery = true
			continue
		}
		opType := rng.Intn(3)
		if len(added) == 0 {
			opType = rng.Intn(2) // only add or query; no elements to remove
			if opType == 1 {
				opType = 2
			}
		}
		switch opType {
		case 0: // add
			x := rng.Intn(1000) + 1
			ops[i] = op{'+', x}
			ms.add(x)
			added = append(added, x)
		case 1: // remove a previously added element
			idx := rng.Intn(len(added))
			x := added[idx]
			ops[i] = op{'-', x}
			ms.remove(x)
			added = append(added[:idx], added[idx+1:]...)
		case 2: // query
			x := rng.Intn(1000) + 1
			ops[i] = op{'?', x}
			hasQuery = true
		}
	}

	// Build input string and expected answers by replaying ops
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", q)
	ms2 := make(multiset)
	ms2.add(0)
	expected := []int{}
	for _, o := range ops {
		fmt.Fprintf(&sb, "%c %d\n", o.kind, o.x)
		switch o.kind {
		case '+':
			ms2.add(o.x)
		case '-':
			ms2.remove(o.x)
		case '?':
			expected = append(expected, ms2.maxXor(o.x))
		}
	}
	return sb.String(), expected
}

func runCase(bin, input string, expected []int) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	scan := bufio.NewScanner(strings.NewReader(out.String()))
	scan.Split(bufio.ScanWords)
	for i, exp := range expected {
		if !scan.Scan() {
			return fmt.Errorf("missing output for query %d", i+1)
		}
		var got int
		fmt.Sscan(scan.Text(), &got)
		if got != exp {
			return fmt.Errorf("query %d: expected %d got %d", i+1, exp, got)
		}
	}
	if scan.Scan() {
		return fmt.Errorf("extra output detected")
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := generateCase(rng)
		if err := runCase(bin, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
