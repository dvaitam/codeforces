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

func solveE(a []int) string {
	n := len(a) - 1
	var sb strings.Builder
	l, r := n, 0
	for i := 1; i <= n; i++ {
		if a[i] > 0 {
			if i < l {
				l = i
			}
			if i > r {
				r = i
			}
		}
	}
	p := -1
	for r > 0 {
		i := p + 2
		for i <= r && a[i] == 0 {
			sb.WriteString("AR")
			p++
			i++
		}
		value := 1
		d := 0
		first := true
		for i = p + 2; i <= r; i++ {
			if a[i] > 0 {
				value += 4
				if first {
					d++
				}
			} else {
				value--
				first = false
			}
			if value <= 0 {
				tar := d
				for j := 0; j < tar; j++ {
					sb.WriteString("AR")
				}
				sb.WriteByte('A')
				for j := 0; j < tar; j++ {
					sb.WriteByte('L')
				}
				sb.WriteByte('A')
				for j := p + 2; j < p+d+2; j++ {
					if a[j] > 0 {
						a[j]--
					}
				}
				break
			}
		}
		if value > 0 {
			tar := r - p - 1
			for i = 0; i < tar; i++ {
				sb.WriteString("AR")
			}
			sb.WriteString("AL")
			p = r - 2
			if a[r] == 1 {
				for p+1 >= l && a[p+1] <= 1 {
					p--
					sb.WriteByte('L')
				}
			}
			for p+1 >= l && a[p+1] > 1 {
				p--
				sb.WriteByte('L')
			}
			sb.WriteByte('A')
			for i = p + 2; i <= r; i++ {
				if a[i] > 0 {
					a[i]--
				}
			}
		}
		l = n
		r = 0
		for i := 1; i <= n; i++ {
			if a[i] > 0 {
				if i < l {
					l = i
				}
				if i > r {
					r = i
				}
			}
		}
	}
	return sb.String()
}

type caseE struct {
	n int
	a []int
}

func genCaseE(rng *rand.Rand) caseE {
	n := rng.Intn(10) + 1
	a := make([]int, n+1)
	pos := false
	for i := 1; i <= n; i++ {
		if rng.Intn(2) == 0 {
			a[i] = rng.Intn(4)
		} else {
			a[i] = 0
		}
		if a[i] > 0 {
			pos = true
		}
	}
	if !pos {
		a[rng.Intn(n)+1] = rng.Intn(3) + 1
	}
	return caseE{n, a}
}

func runCaseE(bin string, c caseE) error {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", c.n))
	for i := 1; i <= c.n; i++ {
		if i > 1 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", c.a[i]))
	}
	sb.WriteByte('\n')
	input := sb.String()
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(out.String())
	expected := solveE(append([]int(nil), c.a...))
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
		c := genCaseE(rng)
		if err := runCaseE(bin, c); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
