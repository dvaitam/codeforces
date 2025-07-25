package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"
	"time"
)

type caseB struct {
	n, k, p, x, y int
	arr           []int
	input         string
	exist         bool
}

func computeSolution(n, k, p, x, y int, arr []int) (bool, []int) {
	sum := 0
	for _, v := range arr {
		sum += v
	}
	nNew := n - k
	if sum+nNew > x {
		return false, nil
	}
	cur := x - (sum + nNew)
	b := make([]int, nNew)
	for i := 0; i < nNew; i++ {
		if cur >= y-1 {
			b[i] = y
			cur -= y - 1
		} else {
			b[i] = 1
		}
	}
	all := append(append([]int(nil), arr...), b...)
	sort.Ints(all)
	if all[n/2] < y || all[n-1] > p {
		return false, nil
	}
	return true, b
}

func genCase(rng *rand.Rand) caseB {
	n := rng.Intn(49)*2 + 1 // odd between 1 and 99
	k := rng.Intn(n)
	p := rng.Intn(20) + 1
	x := rng.Intn(n*p-n+1) + n
	y := rng.Intn(p) + 1
	arr := make([]int, k)
	for i := 0; i < k; i++ {
		arr[i] = rng.Intn(p) + 1
	}
	var bldr strings.Builder
	fmt.Fprintf(&bldr, "%d %d %d %d %d\n", n, k, p, x, y)
	for i := 0; i < k; i++ {
		if i > 0 {
			bldr.WriteByte(' ')
		}
		fmt.Fprintf(&bldr, "%d", arr[i])
	}
	if k > 0 {
		bldr.WriteByte('\n')
	}
	input := bldr.String()
	exist, _ := computeSolution(n, k, p, x, y, arr)
	return caseB{n, k, p, x, y, arr, input, exist}
}

func run(bin, input string) (string, error) {
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

func validate(c caseB, out string) error {
	out = strings.TrimSpace(out)
	if out == "-1" {
		if c.exist {
			return fmt.Errorf("solution exists but candidate printed -1")
		}
		return nil
	}
	if !c.exist {
		return fmt.Errorf("no solution exists but candidate produced one")
	}
	tokens := strings.Fields(out)
	if len(tokens) != c.n-c.k {
		return fmt.Errorf("expected %d numbers, got %d", c.n-c.k, len(tokens))
	}
	b := make([]int, c.n-c.k)
	for i, t := range tokens {
		v, err := strconv.Atoi(t)
		if err != nil {
			return fmt.Errorf("invalid integer: %v", err)
		}
		if v < 1 || v > c.p {
			return fmt.Errorf("mark out of range")
		}
		b[i] = v
	}
	all := append(append([]int(nil), c.arr...), b...)
	total := 0
	for _, v := range all {
		total += v
	}
	sort.Ints(all)
	if total > c.x {
		return fmt.Errorf("total sum exceeds limit")
	}
	if all[c.n/2] < c.y {
		return fmt.Errorf("median below required")
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	const cases = 100
	for i := 0; i < cases; i++ {
		c := genCase(rng)
		got, err := run(bin, c.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: %v\n", i+1, err)
			os.Exit(1)
		}
		if err := validate(c, got); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, c.input)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", cases)
}
