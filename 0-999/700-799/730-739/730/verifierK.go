package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

const refSourceK = `package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"runtime/debug"
)

func init() {
	debug.SetMaxStack(1 << 28)
}

type Scanner struct {
	buf []byte
	pos int
}

func NewScanner() *Scanner {
	b, _ := io.ReadAll(os.Stdin)
	return &Scanner{buf: b, pos: 0}
}

func (s *Scanner) NextInt() int {
	for s.pos < len(s.buf) && s.buf[s.pos] <= ' ' {
		s.pos++
	}
	if s.pos >= len(s.buf) {
		return 0
	}
	res := 0
	for s.pos < len(s.buf) && s.buf[s.pos] > ' ' {
		res = res*10 + int(s.buf[s.pos]-'0')
		s.pos++
	}
	return res
}

type Edge struct {
	to, id int
}

func main() {
	sc := NewScanner()
	if sc.pos >= len(sc.buf) {
		return
	}
	T := sc.NextInt()
	if T == 0 {
		return
	}

	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	for tCase := 0; tCase < T; tCase++ {
		n := sc.NextInt()
		m := sc.NextInt()
		s := sc.NextInt()
		t := sc.NextInt()

		U := make([]int, m+1)
		V := make([]int, m+1)
		adj := make([][]Edge, n+1)

		adj[s] = append(adj[s], Edge{t, 0})
		adj[t] = append(adj[t], Edge{s, 0})

		for i := 1; i <= m; i++ {
			U[i] = sc.NextInt()
			V[i] = sc.NextInt()
			adj[U[i]] = append(adj[U[i]], Edge{V[i], i})
			adj[V[i]] = append(adj[V[i]], Edge{U[i], i})
		}

		for i := 0; i < len(adj[s]); i++ {
			if adj[s][i].id == 0 {
				adj[s][0], adj[s][i] = adj[s][i], adj[s][0]
				break
			}
		}

		timer := 0
		dfn := make([]int, n+1)
		low := make([]int, n+1)
		parent := make([]int, n+1)

		var dfs func(u, edgeID int)
		dfs = func(u, edgeID int) {
			timer++
			dfn[u] = timer
			low[u] = timer
			for _, e := range adj[u] {
				if e.id == edgeID {
					continue
				}
				v := e.to
				if dfn[v] == 0 {
					parent[v] = u
					dfs(v, e.id)
					if low[v] < low[u] {
						low[u] = low[v]
					}
				} else {
					if dfn[v] < low[u] {
						low[u] = dfn[v]
					}
				}
			}
		}

		dfs(s, -1)

		isBi := true
		if timer < n {
			isBi = false
		} else {
			childrenOfS := 0
			for _, e := range adj[s] {
				if parent[e.to] == s {
					childrenOfS++
				}
			}
			if childrenOfS > 1 {
				isBi = false
			}

			for i := 1; i <= n; i++ {
				if i == s {
					continue
				}
				for _, e := range adj[i] {
					if parent[e.to] == i {
						if low[e.to] >= dfn[i] {
							isBi = false
						}
					}
				}
			}
		}

		if !isBi {
			fmt.Fprintln(out, "No")
			continue
		}

		nodeAt := make([]int, n+1)
		for i := 1; i <= n; i++ {
			nodeAt[dfn[i]] = i
		}

		prev := make([]int, n+1)
		next := make([]int, n+1)

		prev[t] = s
		next[s] = t
		sign := make([]int, n+1)
		sign[s] = -1

		for i := 3; i <= n; i++ {
			v := nodeAt[i]
			p := parent[v]
			l := nodeAt[low[v]]

			if sign[l] == -1 {
				pr := prev[p]
				prev[v] = pr
				next[v] = p
				if pr != 0 {
					next[pr] = v
				}
				prev[p] = v
				sign[p] = 1
			} else {
				nx := next[p]
				next[v] = nx
				prev[v] = p
				if nx != 0 {
					prev[nx] = v
				}
				next[p] = v
				sign[p] = -1
			}
		}

		rank := make([]int, n+1)
		r := 0
		for curr := s; curr != 0; curr = next[curr] {
			r++
			rank[curr] = r
		}

		fmt.Fprintln(out, "Yes")
		for i := 1; i <= m; i++ {
			u := U[i]
			v := V[i]
			if rank[u] < rank[v] {
				fmt.Fprintln(out, u, v)
			} else {
				fmt.Fprintln(out, v, u)
			}
		}
	}
}
`

func buildRef() (string, error) {
	tmp, err := os.CreateTemp("", "refK_*.go")
	if err != nil {
		return "", err
	}
	if _, err := tmp.WriteString(refSourceK); err != nil {
		tmp.Close()
		return "", err
	}
	tmp.Close()
	defer os.Remove(tmp.Name())
	ref := filepath.Join(os.TempDir(), "refK.bin")
	cmd := exec.Command("go", "build", "-o", ref, tmp.Name())
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("failed to build reference: %v\n%s", err, out)
	}
	return ref, nil
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
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

type Case struct{ input string }

func genCases() []Case {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	cases := make([]Case, 100)
	for i := range cases {
		T := rng.Intn(2) + 1
		var sb strings.Builder
		fmt.Fprintf(&sb, "%d\n", T)
		for t := 0; t < T; t++ {
			n := rng.Intn(4) + 2
			maxEdges := n * (n - 1) / 2
			minEdges := n - 1
			m := minEdges
			if maxEdges > minEdges {
				m = rng.Intn(maxEdges-minEdges+1) + minEdges
			}
			s := rng.Intn(n) + 1
			tdest := rng.Intn(n) + 1
			for tdest == s {
				tdest = rng.Intn(n) + 1
			}
			fmt.Fprintf(&sb, "%d %d %d %d\n", n, m, s, tdest)
			edges := make(map[[2]int]struct{})
			for j := 0; j < m; j++ {
				u := rng.Intn(n) + 1
				v := rng.Intn(n) + 1
				for u == v || edges[[2]int{u, v}] != (struct{}{}) {
					u = rng.Intn(n) + 1
					v = rng.Intn(n) + 1
				}
				edges[[2]int{u, v}] = struct{}{}
				edges[[2]int{v, u}] = struct{}{}
				fmt.Fprintf(&sb, "%d %d\n", u, v)
			}
		}
		cases[i] = Case{sb.String()}
	}
	return cases
}

func runCase(bin, ref string, c Case) error {
	exp, err := runBinary(ref, c.input)
	if err != nil {
		return fmt.Errorf("reference failed: %v", err)
	}
	got, err := runBinary(bin, c.input)
	if err != nil {
		return err
	}
	if strings.TrimSpace(exp) != strings.TrimSpace(got) {
		return fmt.Errorf("expected %s got %s", exp, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierK.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	ref, err := buildRef()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer os.Remove(ref)
	cases := genCases()
	for i, c := range cases {
		if err := runCase(bin, ref, c); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, c.input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
