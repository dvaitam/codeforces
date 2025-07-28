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

type operation struct {
	typ int
	pos int
	a1  int
	p1  int
	a2  int
	p2  int
}

func solveE(input string) string {
	in := bufio.NewReader(strings.NewReader(input))
	out := new(bytes.Buffer)
	var T int
	fmt.Fscan(in, &T)
	for ; T > 0; T-- {
		var s1, s2 string
		fmt.Fscan(in, &s1)
		fmt.Fscan(in, &s2)
		n := len(s1)
		b1 := []byte(s1)
		b2 := []byte(s2)
		var t, q int
		fmt.Fscan(in, &t, &q)
		events := make([][]int, q+t+5)
		blocked := make([]bool, n)
		mism := 0
		for i := 0; i < n; i++ {
			if b1[i] != b2[i] {
				mism++
			}
		}
		for time := 1; time <= q; time++ {
			for _, pos := range events[time] {
				if blocked[pos] {
					blocked[pos] = false
					if b1[pos] != b2[pos] {
						mism++
					}
				}
			}
			var typ int
			fmt.Fscan(in, &typ)
			switch typ {
			case 1:
				var pos int
				fmt.Fscan(in, &pos)
				pos--
				if !blocked[pos] {
					if b1[pos] != b2[pos] {
						mism--
					}
					blocked[pos] = true
					events[time+t] = append(events[time+t], pos)
				}
			case 2:
				var a1, p1, a2, p2 int
				fmt.Fscan(in, &a1, &p1, &a2, &p2)
				p1--
				p2--
				idxs := map[int]struct{}{p1: {}, p2: {}}
				for idx := range idxs {
					if !blocked[idx] && b1[idx] != b2[idx] {
						mism--
					}
				}
				var c1, c2 byte
				if a1 == 1 {
					c1 = b1[p1]
				} else {
					c1 = b2[p1]
				}
				if a2 == 1 {
					c2 = b1[p2]
				} else {
					c2 = b2[p2]
				}
				if a1 == 1 {
					b1[p1] = c2
				} else {
					b2[p1] = c2
				}
				if a2 == 1 {
					b1[p2] = c1
				} else {
					b2[p2] = c1
				}
				for idx := range idxs {
					if !blocked[idx] && b1[idx] != b2[idx] {
						mism++
					}
				}
			case 3:
				if mism == 0 {
					fmt.Fprintln(out, "YES")
				} else {
					fmt.Fprintln(out, "NO")
				}
			}
		}
	}
	return out.String()
}

func generateCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(5) + 1
	letters := "abcde"
	makeStr := func() string {
		b := make([]byte, n)
		for i := 0; i < n; i++ {
			b[i] = letters[rng.Intn(len(letters))]
		}
		return string(b)
	}
	s1 := makeStr()
	s2 := makeStr()
	tval := rng.Intn(3) + 1
	q := rng.Intn(5) + 1
	blockedUntil := make([]int, n)
	ops := make([]string, 0, q)
	for timeStep := 1; timeStep <= q; timeStep++ {
		for i := 0; i < n; i++ {
			if blockedUntil[i] == timeStep {
				blockedUntil[i] = 0
			}
		}
		typ := rng.Intn(3) + 1
		switch typ {
		case 1:
			avail := []int{}
			for i := 0; i < n; i++ {
				if blockedUntil[i] == 0 {
					avail = append(avail, i)
				}
			}
			if len(avail) == 0 {
				typ = 3
			} else {
				pos := avail[rng.Intn(len(avail))]
				blockedUntil[pos] = timeStep + tval
				ops = append(ops, fmt.Sprintf("1 %d", pos+1))
				continue
			}
			fallthrough
		case 2:
			if typ == 2 {
				avail := []int{}
				for i := 0; i < n; i++ {
					if blockedUntil[i] == 0 {
						avail = append(avail, i)
					}
				}
				if len(avail) == 0 {
					typ = 3
				} else {
					p1 := avail[rng.Intn(len(avail))]
					p2 := avail[rng.Intn(len(avail))]
					a1 := rng.Intn(2) + 1
					a2 := rng.Intn(2) + 1
					ops = append(ops, fmt.Sprintf("2 %d %d %d %d", a1, p1+1, a2, p2+1))
					continue
				}
			}
			fallthrough
		case 3:
			ops = append(ops, "3")
		}
	}
	var sb strings.Builder
	sb.WriteString("1\n")
	sb.WriteString(s1 + "\n")
	sb.WriteString(s2 + "\n")
	sb.WriteString(fmt.Sprintf("%d %d\n", tval, q))
	for _, op := range ops {
		sb.WriteString(op + "\n")
	}
	expected := solveE(sb.String())
	return sb.String(), strings.TrimSpace(expected)
}

func runCase(bin, input, expected string) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errb bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errb
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, errb.String())
	}
	got := strings.TrimSpace(out.String())
	if got != expected {
		return fmt.Errorf("expected %q got %q", expected, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
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
