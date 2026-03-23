package main

import (
	"bytes"
	"fmt"
	"io"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

// ---------- embedded solver (from cf_t23_1361_E.go) ----------

type FastScanner struct {
	data []byte
	idx  int
	n    int
}

func NewFastScanner(input []byte) *FastScanner {
	return &FastScanner{data: input, n: len(input)}
}

func (fs *FastScanner) NextInt() int {
	for fs.idx < fs.n {
		c := fs.data[fs.idx]
		if c >= '0' && c <= '9' {
			break
		}
		fs.idx++
	}
	val := 0
	for fs.idx < fs.n {
		c := fs.data[fs.idx]
		if c < '0' || c > '9' {
			break
		}
		val = val*10 + int(c-'0')
		fs.idx++
	}
	return val
}

func writeIntBuf(w *bytes.Buffer, x int) {
	if x == 0 {
		w.WriteByte('0')
		return
	}
	if x < 0 {
		w.WriteByte('-')
		x = -x
	}
	var buf [20]byte
	i := len(buf)
	for x > 0 {
		i--
		buf[i] = byte('0' + x%10)
		x /= 10
	}
	w.Write(buf[i:])
}

func solveEmbedded(input []byte) string {
	in := NewFastScanner(input)
	out := &bytes.Buffer{}

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	t := in.NextInt()

	for ; t > 0; t-- {
		n := in.NextInt()
		m := in.NextInt()

		adj := make([][]int, n+1)
		for i := 0; i < m; i++ {
			u := in.NextInt()
			v := in.NextInt()
			adj[u] = append(adj[u], v)
		}

		parent := make([]int, n+1)
		tin := make([]int, n+1)
		tout := make([]int, n+1)
		it := make([]int, n+1)
		orderArr := make([]int, n)
		stackArr := make([]int, n)

		runCheck := func(root int, collect bool) (bool, []int, []int, []int, []int) {
			for i := 1; i <= n; i++ {
				parent[i] = 0
				tin[i] = 0
				tout[i] = 0
				it[i] = 0
			}

			orderLen := 0
			stackLen := 1
			stackArr[0] = root
			timer := 0

			for stackLen > 0 {
				u := stackArr[stackLen-1]
				if tin[u] == 0 {
					timer++
					tin[u] = timer
				}
				if it[u] < len(adj[u]) {
					v := adj[u][it[u]]
					it[u]++
					if tin[v] == 0 {
						parent[v] = u
						stackArr[stackLen] = v
						stackLen++
					}
				} else {
					timer++
					tout[u] = timer
					orderArr[orderLen] = u
					orderLen++
					stackLen--
				}
			}

			if orderLen != n {
				return false, nil, nil, nil, nil
			}

			var backT, backH []int
			if collect {
				backT = make([]int, 0, m)
				backH = make([]int, 0, m)
			}

			for u := 1; u <= n; u++ {
				for _, v := range adj[u] {
					if parent[v] == u {
						continue
					}
					if !(tin[v] <= tin[u] && tout[u] <= tout[v]) {
						return false, nil, nil, nil, nil
					}
					if collect {
						backT = append(backT, u)
						backH = append(backH, v)
					}
				}
			}

			return true, parent, orderArr[:orderLen], backT, backH
		}

		const sampleLimit = 120
		root := 0

		if n <= sampleLimit {
			for v := 1; v <= n; v++ {
				ok, _, _, _, _ := runCheck(v, false)
				if ok {
					root = v
					break
				}
			}
		} else {
			perm := make([]int, n)
			for i := 0; i < n; i++ {
				perm[i] = i + 1
			}
			for i := 0; i < sampleLimit; i++ {
				j := i + rng.Intn(n-i)
				perm[i], perm[j] = perm[j], perm[i]
			}
			for i := 0; i < sampleLimit; i++ {
				v := perm[i]
				ok, _, _, _, _ := runCheck(v, false)
				if ok {
					root = v
					break
				}
			}
		}

		if root == 0 {
			writeIntBuf(out, -1)
			out.WriteByte('\n')
			continue
		}

		ok, parentRes, order, backT, backH := runCheck(root, true)
		if !ok {
			writeIntBuf(out, -1)
			out.WriteByte('\n')
			continue
		}

		cnt := make([]int, n+1)
		xv := make([]int, n+1)
		headOf := make([]int, len(backH)+1)

		for i := 0; i < len(backH); i++ {
			id := i + 1
			u := backT[i]
			a := backH[i]
			cnt[u]++
			cnt[a]--
			xv[u] ^= id
			xv[a] ^= id
			headOf[id] = a
		}

		for _, u := range order {
			p := parentRes[u]
			if p != 0 {
				cnt[p] += cnt[u]
				xv[p] ^= xv[u]
			}
		}

		good := make([]bool, n+1)
		for i := len(order) - 1; i >= 0; i-- {
			u := order[i]
			if u == root {
				good[u] = true
			} else if cnt[u] == 1 {
				good[u] = good[headOf[xv[u]]]
			}
		}

		k := 0
		for i := 1; i <= n; i++ {
			if good[i] {
				k++
			}
		}

		if k*5 < n {
			writeIntBuf(out, -1)
			out.WriteByte('\n')
			continue
		}

		first := true
		for i := 1; i <= n; i++ {
			if good[i] {
				if !first {
					out.WriteByte(' ')
				}
				first = false
				writeIntBuf(out, i)
			}
		}
		out.WriteByte('\n')
	}
	return strings.TrimSpace(out.String())
}

// ---------- verifier infrastructure ----------

func run(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(stdout.String()), nil
}

func genCase(rng *rand.Rand) string {
	n := rng.Intn(5) + 1
	maxEdges := n * (n - 1)
	m := rng.Intn(maxEdges + 1)
	edges := make([][2]int, 0, m)
	for len(edges) < m {
		u := rng.Intn(n) + 1
		v := rng.Intn(n) + 1
		if u == v {
			continue
		}
		edges = append(edges, [2]int{u, v})
	}
	var sb strings.Builder
	sb.WriteString("1\n")
	sb.WriteString(fmt.Sprintf("%d %d\n", n, m))
	for _, e := range edges {
		sb.WriteString(fmt.Sprintf("%d %d\n", e[0], e[1]))
	}
	return sb.String()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	_ = io.Discard // suppress unused import

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	const cases = 100
	for i := 1; i <= cases; i++ {
		input := genCase(rng)
		expect := solveEmbedded([]byte(input))
		got, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: %v\n", i, err)
			os.Exit(1)
		}
		if got != expect {
			fmt.Printf("case %d failed\nexpected:\n%s\n\ngot:\n%s\n", i, expect, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", cases)
}
