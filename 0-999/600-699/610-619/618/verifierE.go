package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

// ---------- embedded solver (from cf_t25_618_E.go) ----------

type Node struct {
	x, y float64
	lazy int
}

var (
	tree    []Node
	cos_deg [360]float64
	sin_deg [360]float64
)

func build(u, l, r int) {
	if l == r {
		tree[u].x = 1.0
		tree[u].y = 0.0
		tree[u].lazy = 0
		return
	}
	mid := (l + r) / 2
	build(u*2, l, mid)
	build(u*2+1, mid+1, r)
	push_up(u)
}

func push_up(u int) {
	tree[u].x = tree[u*2].x + tree[u*2+1].x
	tree[u].y = tree[u*2].y + tree[u*2+1].y
}

func apply_rot(u int, d int) {
	nx := tree[u].x*cos_deg[d] - tree[u].y*sin_deg[d]
	ny := tree[u].x*sin_deg[d] + tree[u].y*cos_deg[d]
	tree[u].x = nx
	tree[u].y = ny
	tree[u].lazy = (tree[u].lazy + d) % 360
}

func push_down(u int) {
	if tree[u].lazy != 0 {
		apply_rot(u*2, tree[u].lazy)
		apply_rot(u*2+1, tree[u].lazy)
		tree[u].lazy = 0
	}
}

func update_extend(u, l, r, idx int, val float64) {
	if l == r {
		ln := math.Hypot(tree[u].x, tree[u].y)
		tree[u].x += tree[u].x / ln * val
		tree[u].y += tree[u].y / ln * val
		return
	}
	push_down(u)
	mid := (l + r) / 2
	if idx <= mid {
		update_extend(u*2, l, mid, idx, val)
	} else {
		update_extend(u*2+1, mid+1, r, idx, val)
	}
	push_up(u)
}

func update_rotate(u, l, r, ql, qr, d int) {
	if ql <= l && r <= qr {
		apply_rot(u, d)
		return
	}
	push_down(u)
	mid := (l + r) / 2
	if ql <= mid {
		update_rotate(u*2, l, mid, ql, qr, d)
	}
	if qr > mid {
		update_rotate(u*2+1, mid+1, r, ql, qr, d)
	}
	push_up(u)
}

func readIntBuf(in *bufio.Reader) int {
	var res int
	var c byte
	var err error
	for {
		c, err = in.ReadByte()
		if err != nil {
			return 0
		}
		if c >= '0' && c <= '9' {
			break
		}
	}
	res = int(c - '0')
	for {
		c, err = in.ReadByte()
		if err != nil || c < '0' || c > '9' {
			break
		}
		res = res*10 + int(c-'0')
	}
	return res
}

func solveE(input string) string {
	for i := 0; i < 360; i++ {
		rad := float64(i) * math.Pi / 180.0
		cos_deg[i] = math.Cos(rad)
		sin_deg[i] = math.Sin(rad)
	}

	in := bufio.NewReaderSize(strings.NewReader(input), 1<<20)
	var out bytes.Buffer

	n := readIntBuf(in)
	m := readIntBuf(in)

	tree = make([]Node, 4*n+1)
	if n > 0 {
		build(1, 1, n)
	}

	for i := 0; i < m; i++ {
		typ := readIntBuf(in)
		idx := readIntBuf(in)
		val := readIntBuf(in)

		if typ == 1 {
			update_extend(1, 1, n, idx, float64(val))
		} else {
			d := (360 - (val % 360)) % 360
			update_rotate(1, 1, n, idx, n, d)
		}
		fmt.Fprintf(&out, "%.10f %.10f\n", tree[1].x, tree[1].y)
	}
	return strings.TrimSpace(out.String())
}

// ---------- end embedded solver ----------

func run(prog, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(prog, ".go") {
		cmd = exec.Command("go", "run", prog)
	} else {
		cmd = exec.Command(prog)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func genCase(rng *rand.Rand, n, m int) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", n, m))
	for i := 0; i < m; i++ {
		typ := rng.Intn(2) + 1
		seg := rng.Intn(n) + 1
		val := rng.Intn(359) + 1
		sb.WriteString(fmt.Sprintf("%d %d %d\n", typ, seg, val))
	}
	return sb.String()
}

func parseFloatPair(s string) (float64, float64, error) {
	parts := strings.Fields(s)
	if len(parts) != 2 {
		return 0, 0, fmt.Errorf("expected two floats, got %d tokens", len(parts))
	}
	a, err := strconv.ParseFloat(parts[0], 64)
	if err != nil {
		return 0, 0, err
	}
	b, err := strconv.ParseFloat(parts[1], 64)
	if err != nil {
		return 0, 0, err
	}
	return a, b, nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for t := 0; t < 120; t++ {
		n := rng.Intn(10) + 1
		m := rng.Intn(10) + 1
		input := genCase(rng, n, m)
		expectedStr := solveE(input)
		gotStr, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", t+1, err, input)
			os.Exit(1)
		}
		expLines := strings.Split(expectedStr, "\n")
		gotLines := strings.Split(gotStr, "\n")
		if len(gotLines) != len(expLines) {
			fmt.Fprintf(os.Stderr, "case %d: expected %d lines got %d lines\ninput:\n%s", t+1, len(expLines), len(gotLines), input)
			os.Exit(1)
		}
		for li := 0; li < len(expLines); li++ {
			exX, exY, err := parseFloatPair(expLines[li])
			if err != nil {
				fmt.Fprintf(os.Stderr, "internal parse error on line %d: %v\n", li+1, err)
				os.Exit(1)
			}
			gX, gY, err := parseFloatPair(gotLines[li])
			if err != nil {
				fmt.Fprintf(os.Stderr, "case %d line %d: expected two numbers got %q\n", t+1, li+1, gotLines[li])
				os.Exit(1)
			}
			if math.Abs(gX-exX) > 1e-4*math.Max(1, math.Abs(exX)) || math.Abs(gY-exY) > 1e-4*math.Max(1, math.Abs(exY)) {
				fmt.Fprintf(os.Stderr, "case %d line %d: expected %.5f %.5f got %.5f %.5f\ninput:\n%s", t+1, li+1, exX, exY, gX, gY, input)
				os.Exit(1)
			}
		}
	}
	fmt.Println("All tests passed")
}
