package main

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"time"
)

const mod = 998244353

type Test struct {
	n      int
	x      int64
	d      []int64
	length int
	count  int
}

func computeLIS(x int64, d []int64) (int, int) {
	p := []int64{x}
	cur := x
	for _, v := range d {
		if v > 0 {
			for i := int64(0); i < v; i++ {
				cur++
				p = append(p, cur)
			}
		} else if v < 0 {
			for i := int64(0); i < -v; i++ {
				cur--
				p = append(p, cur)
			}
		}
	}
	n := len(p)
	length := make([]int, n)
	count := make([]int, n)
	for i := 0; i < n; i++ {
		length[i] = 1
		count[i] = 1
		for j := 0; j < i; j++ {
			if p[j] < p[i] {
				if length[j]+1 > length[i] {
					length[i] = length[j] + 1
					count[i] = count[j]
				} else if length[j]+1 == length[i] {
					count[i] = (count[i] + count[j]) % mod
				}
			}
		}
	}
	best := 0
	total := 0
	for i := 0; i < n; i++ {
		if length[i] > best {
			best = length[i]
			total = count[i]
		} else if length[i] == best {
			total = (total + count[i]) % mod
		}
	}
	return best, total
}

func genTests() []Test {
	r := rand.New(rand.NewSource(42))
	tests := make([]Test, 0, 100)
	for len(tests) < 100 {
		n := r.Intn(4) + 1
		x := int64(r.Intn(11) - 5)
		d := make([]int64, n)
		for i := range d {
			d[i] = int64(r.Intn(7) - 3)
		}
		l, c := computeLIS(x, d)
		tests = append(tests, Test{n: n, x: x, d: d, length: l, count: c})
	}
	return tests
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierF.go /path/to/binary")
		return
	}
	binary := os.Args[1]
	tests := genTests()

	var input bytes.Buffer
	fmt.Fprintln(&input, len(tests))
	for _, t := range tests {
		fmt.Fprintln(&input, t.n)
		fmt.Fprintln(&input, t.x)
		for i, v := range t.d {
			if i > 0 {
				fmt.Fprint(&input, " ")
			}
			fmt.Fprint(&input, v)
		}
		fmt.Fprintln(&input)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, binary)
	cmd.Stdin = bytes.NewReader(input.Bytes())
	out, err := cmd.Output()
	if ctx.Err() == context.DeadlineExceeded {
		fmt.Println("time limit exceeded")
		return
	}
	if err != nil {
		fmt.Println("execution error:", err)
		return
	}

	scanner := bufio.NewScanner(bytes.NewReader(out))
	for i, t := range tests {
		if !scanner.Scan() {
			fmt.Printf("missing output for case %d\n", i+1)
			return
		}
		var l, c int
		if _, err := fmt.Sscan(scanner.Text(), &l, &c); err != nil {
			fmt.Printf("bad output on case %d\n", i+1)
			return
		}
		if l != t.length || c%mod != t.count%mod {
			fmt.Printf("wrong answer on case %d\n", i+1)
			return
		}
	}
	if scanner.Scan() {
		fmt.Println("extra output detected")
		return
	}
	fmt.Println("OK")
}
