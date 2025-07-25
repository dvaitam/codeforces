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
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", q)
	ms := make(multiset)
	ms.add(0)
	expected := []int{}
	exists := []int{0}
	for i := 0; i < q; i++ {
		opType := rng.Intn(3)
		if len(exists) == 1 {
			opType = rng.Intn(2) // avoid remove when only 0 present
		}
		switch opType {
		case 0: // add
			x := rng.Intn(1000)
			fmt.Fprintf(&sb, "+ %d\n", x)
			ms.add(x)
			exists = append(exists, x)
		case 1: // remove
			idx := rng.Intn(len(exists))
			x := exists[idx]
			fmt.Fprintf(&sb, "- %d\n", x)
			ms.remove(x)
			exists = append(exists[:idx], exists[idx+1:]...)
		case 2: // query
			x := rng.Intn(1000)
			fmt.Fprintf(&sb, "? %d\n", x)
			expected = append(expected, ms.maxXor(x))
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
