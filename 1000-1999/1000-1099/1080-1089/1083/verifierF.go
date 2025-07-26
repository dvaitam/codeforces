package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

type test struct {
	input    string
	expected string
}

func runBinary(bin, input string) (string, error) {
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
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

const B = 500

func solveF(input string) string {
	reader := bufio.NewReader(strings.NewReader(input))
	writer := bufio.NewWriter(&bytes.Buffer{})
	var n, k, q int
	fmt.Fscan(reader, &n, &k, &q)
	a := make([]int, n+2)
	b := make([]int, n+2)
	for i := 1; i <= n; i++ {
		fmt.Fscan(reader, &a[i])
	}
	for i := 1; i <= n; i++ {
		fmt.Fscan(reader, &b[i])
	}
	d := make([]int, n+2)
	nch := make([][]int, k)
	for i := 0; i < k; i++ {
		nch[i] = make([]int, 0)
	}
	for i := 0; i <= n; i++ {
		d[i] = (a[i] ^ a[i+1]) ^ (b[i] ^ b[i+1])
		rem := i % k
		nch[rem] = append(nch[rem], i)
	}
	pos := make([]int, n+2)
	epos := make([]int, n+2)
	spos := make([]int, n+2)
	m := 0
	zch := 0
	for i := 0; i < k; i++ {
		s := 0
		for _, u := range nch[i] {
			s ^= d[u]
			spos[m] = s
			pos[u] = m
			m++
		}
		if m > 0 && spos[m-1] == 0 {
			zch++
		}
		for _, u := range nch[i] {
			epos[u] = m
		}
	}
	blocks := (m + B - 1) / B
	off := make([]int, blocks)
	cnt := make([]map[int]int, blocks)
	for i := 0; i < blocks; i++ {
		cnt[i] = make(map[int]int)
	}
	for i := 0; i < m; i++ {
		blk := i / B
		cnt[blk][spos[i]]++
	}
	modify := func(p, v int) {
		if p < 0 || p > n {
			return
		}
		l := pos[p]
		r := epos[p]
		if r > 0 {
			br := (r - 1) / B
			if spos[r-1]^off[br] == 0 {
				zch--
			}
		}
		idl := l / B
		idr := r / B
		if idl == idr {
			for i := l; i < r; i++ {
				blk := idl
				cnt[blk][spos[i]]--
				spos[i] ^= v
				cnt[blk][spos[i]]++
			}
		} else {
			endL := (idl + 1) * B
			for i := l; i < endL; i++ {
				blk := idl
				cnt[blk][spos[i]]--
				spos[i] ^= v
				cnt[blk][spos[i]]++
			}
			for blk := idl + 1; blk < idr; blk++ {
				off[blk] ^= v
			}
			startR := idr * B
			for i := startR; i < r; i++ {
				blk := idr
				cnt[blk][spos[i]]--
				spos[i] ^= v
				cnt[blk][spos[i]]++
			}
		}
		if r > 0 {
			br := (r - 1) / B
			if spos[r-1]^off[br] == 0 {
				zch++
			}
		}
	}
	output := func(w *bufio.Writer) string {
		if zch != k {
			return "-1"
		}
		sumZero := 0
		for blk := range off {
			sumZero += cnt[blk][off[blk]]
		}
		res := m - sumZero
		return fmt.Sprintf("%d", res)
	}
	resBuilder := strings.Builder{}
	resBuilder.WriteString(output(writer))
	resBuilder.WriteByte('\n')
	for qi := 0; qi < q; qi++ {
		var op string
		var p, v int
		fmt.Fscan(reader, &op, &p, &v)
		if op[0] == 'a' {
			a[p] ^= v
			v, a[p] = a[p], v
		} else {
			b[p] ^= v
			v, b[p] = b[p], v
		}
		modify(p, v)
		modify(p-1, v)
		resBuilder.WriteString(output(writer))
		resBuilder.WriteByte('\n')
	}
	return strings.TrimSpace(resBuilder.String())
}

func generateTests() []test {
	rng := rand.New(rand.NewSource(47))
	var tests []test
	for len(tests) < 100 {
		n := rng.Intn(5) + 1
		k := rng.Intn(3) + 1
		q := rng.Intn(3) + 1
		var sb strings.Builder
		fmt.Fprintf(&sb, "%d %d %d\n", n, k, q)
		a := make([]int, n+1)
		b := make([]int, n+1)
		for i := 1; i <= n; i++ {
			a[i] = rng.Intn(8)
			fmt.Fprintf(&sb, "%d ", a[i])
		}
		sb.WriteByte('\n')
		for i := 1; i <= n; i++ {
			b[i] = rng.Intn(8)
			fmt.Fprintf(&sb, "%d ", b[i])
		}
		sb.WriteByte('\n')
		for i := 0; i < q; i++ {
			if rng.Intn(2) == 0 {
				p := rng.Intn(n + 1)
				v := rng.Intn(4)
				fmt.Fprintf(&sb, "a %d %d\n", p, v)
			} else {
				p := rng.Intn(n + 1)
				v := rng.Intn(4)
				fmt.Fprintf(&sb, "b %d %d\n", p, v)
			}
		}
		input := sb.String()
		expected := solveF(input)
		tests = append(tests, test{input, expected})
	}
	return tests
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := generateTests()
	for i, t := range tests {
		out, err := runBinary(bin, t.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d runtime error: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(out) != strings.TrimSpace(t.expected) {
			fmt.Fprintf(os.Stderr, "case %d failed\ninput:\n%sexpected:%s\ngot:%s\n", i+1, t.input, t.expected, out)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed.\n", len(tests))
}
