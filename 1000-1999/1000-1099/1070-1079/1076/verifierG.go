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

func solve(n, m, q int, arr []int, queries [][]int64) string {
	a := make([]int, n)
	copy(a, arr)
	for _, qu := range queries {
		if qu[0] == 1 {
			l := int(qu[1]) - 1
			r := int(qu[2]) - 1
			d := qu[3]
			if d%2 == 1 {
				for i := l; i <= r; i++ {
					a[i] ^= 1
				}
			}
		} else {
			l := int(qu[1]) - 1
			r := int(qu[2]) - 1
			xor := 0
			for i := l; i <= r; i++ {
				xor ^= a[i]
			}
			if xor == 0 {
				qu[3] = 1
			} else {
				qu[3] = 2
			}
		}
	}
	var sb strings.Builder
	for _, qu := range queries {
		if qu[0] == 2 {
			if sb.Len() > 0 {
				sb.WriteByte('\n')
			}
			sb.WriteString(fmt.Sprintf("%d", qu[3]))
		}
	}
	return sb.String()
}

func generateCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(5) + 1
	m := rng.Intn(5) + 1
	q := rng.Intn(5) + 1
	arr := make([]int, n)
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d %d\n", n, m, q))
	for i := 0; i < n; i++ {
		arr[i] = rng.Intn(10) + 1
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", arr[i]))
	}
	sb.WriteByte('\n')
	queries := make([][]int64, q)
	type2 := false
	for i := 0; i < q; i++ {
		tp := rng.Intn(2) + 1
		if i == q-1 && !type2 {
			tp = 2
		}
		if tp == 1 {
			l := rng.Intn(n) + 1
			r := rng.Intn(n-l+1) + l
			d := int64(rng.Intn(5) + 1)
			queries[i] = []int64{1, int64(l), int64(r), d}
			sb.WriteString(fmt.Sprintf("1 %d %d %d\n", l, r, d))
		} else {
			l := rng.Intn(n) + 1
			r := rng.Intn(n-l+1) + l
			queries[i] = []int64{2, int64(l), int64(r), 0}
			sb.WriteString(fmt.Sprintf("2 %d %d\n", l, r))
			type2 = true
		}
	}
	expect := solve(n, m, q, arr, queries)
	return sb.String(), expect
}

func runCase(bin, in, expect string) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(in)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	got := strings.TrimSpace(out.String())
	expect = strings.TrimSpace(expect)
	if got != expect {
		return fmt.Errorf("expected:\n%s\n\ngot:\n%s", expect, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierG.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, expect := generateCase(rng)
		if err := runCase(bin, in, expect); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
