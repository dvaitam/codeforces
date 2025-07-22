package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strings"
	"time"
)

type testCaseA struct {
	input  string
	output string
}

func findValueA(pos int64, tp []int, x, l, lenArr []int64) int64 {
	for {
		idx := sort.Search(len(lenArr), func(i int) bool { return lenArr[i] >= pos })
		if tp[idx] == 1 {
			return x[idx]
		}
		prevLen := int64(0)
		if idx > 0 {
			prevLen = lenArr[idx-1]
		}
		pos = (pos-prevLen-1)%l[idx] + 1
	}
}

func solveA(r io.Reader) string {
	reader := bufio.NewReader(r)
	var m int
	if _, err := fmt.Fscan(reader, &m); err != nil {
		return ""
	}
	tp := make([]int, m)
	x := make([]int64, m)
	l := make([]int64, m)
	c := make([]int64, m)
	lenArr := make([]int64, m)
	var curLen int64
	for i := 0; i < m; i++ {
		fmt.Fscan(reader, &tp[i])
		if tp[i] == 1 {
			fmt.Fscan(reader, &x[i])
			curLen++
		} else {
			fmt.Fscan(reader, &l[i], &c[i])
			curLen += l[i] * c[i]
		}
		lenArr[i] = curLen
	}
	var q int
	fmt.Fscan(reader, &q)
	queries := make([]int64, q)
	for i := 0; i < q; i++ {
		fmt.Fscan(reader, &queries[i])
	}
	var sb strings.Builder
	for i, pos := range queries {
		if i > 0 {
			sb.WriteByte(' ')
		}
		fmt.Fprint(&sb, findValueA(pos, tp, x, l, lenArr))
	}
	sb.WriteByte('\n')
	return sb.String()
}

func runCaseA(bin string, in string, expect string) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(in)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	if strings.TrimSpace(out.String()) != strings.TrimSpace(expect) {
		return fmt.Errorf("expected %q got %q", strings.TrimSpace(expect), strings.TrimSpace(out.String()))
	}
	return nil
}

func genCaseA(rng *rand.Rand) string {
	m := rng.Intn(8) + 1
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", m))
	seqLen := 0
	for i := 0; i < m; i++ {
		tp := rng.Intn(2) + 1
		if tp == 1 || seqLen == 0 {
			val := rng.Intn(100) + 1
			sb.WriteString(fmt.Sprintf("1 %d\n", val))
			seqLen++
		} else {
			l := rng.Intn(seqLen) + 1
			c := rng.Intn(2) + 1
			if seqLen+l*c > 1000 {
				if l > 0 {
					c = (1000 - seqLen) / l
				}
				if c <= 0 {
					val := rng.Intn(100) + 1
					sb.WriteString(fmt.Sprintf("1 %d\n", val))
					seqLen++
					continue
				}
			}
			sb.WriteString(fmt.Sprintf("2 %d %d\n", l, c))
			seqLen += l * c
		}
	}
	q := rng.Intn(min(seqLen, 10)) + 1
	sb.WriteString(fmt.Sprintf("%d\n", q))
	positions := make([]int, q)
	for i := 0; i < q; i++ {
		positions[i] = rng.Intn(seqLen) + 1
	}
	sort.Ints(positions)
	for i, p := range positions {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprint(p))
	}
	sb.WriteByte('\n')
	return sb.String()
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in := genCaseA(rng)
		expect := solveA(strings.NewReader(in))
		if err := runCaseA(bin, in, expect); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
