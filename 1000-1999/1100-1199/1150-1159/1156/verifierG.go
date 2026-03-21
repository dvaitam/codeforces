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

// Embedded solver from the accepted solution

type SNode struct {
	leaf bool
	name string
	op   byte
	l    int
	r    int
}

type SOpKey struct {
	op byte
	l  int
	r  int
}

func oracleSolve(input string) string {
	in := bufio.NewReader(strings.NewReader(input))
	var out bytes.Buffer
	w := bufio.NewWriter(&out)

	var n int
	fmt.Fscan(in, &n)

	nodes := make([]SNode, 0, 2*n+10)
	leafID := make(map[string]int)
	opID := make(map[SOpKey]int)

	var getLeaf func(string) int
	getLeaf = func(name string) int {
		if id, ok := leafID[name]; ok {
			return id
		}
		id := len(nodes)
		nodes = append(nodes, SNode{leaf: true, name: name})
		leafID[name] = id
		return id
	}

	getOp := func(op byte, l, r int) int {
		key := SOpKey{op: op, l: l, r: r}
		if id, ok := opID[key]; ok {
			return id
		}
		id := len(nodes)
		nodes = append(nodes, SNode{op: op, l: l, r: r})
		opID[key] = id
		return id
	}

	curTerm := make(map[string]int)
	var cur func(string) int
	cur = func(name string) int {
		if id, ok := curTerm[name]; ok {
			return id
		}
		id := getLeaf(name)
		curTerm[name] = id
		return id
	}

	used := make(map[string]bool)
	used["res"] = true

	for i := 0; i < n; i++ {
		var line string
		fmt.Fscan(in, &line)

		eq := strings.IndexByte(line, '=')
		lhs := line[:eq]
		rhs := line[eq+1:]

		used[lhs] = true

		pos := strings.IndexAny(rhs, "$^#&")
		if pos == -1 {
			used[rhs] = true
			curTerm[lhs] = cur(rhs)
		} else {
			a := rhs[:pos]
			op := rhs[pos]
			b := rhs[pos+1:]
			used[a] = true
			used[b] = true
			curTerm[lhs] = getOp(op, cur(a), cur(b))
		}
	}

	root := cur("res")

	if nodes[root].leaf {
		if nodes[root].name == "res" {
			fmt.Fprintln(w, 0)
		} else {
			fmt.Fprintln(w, 1)
			fmt.Fprintf(w, "res=%s\n", nodes[root].name)
		}
		w.Flush()
		return strings.TrimSpace(out.String())
	}

	chars := []byte("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	tempCount := 0
	nextTemp := func() string {
		for {
			x := tempCount
			tempCount++
			a := chars[(x/(62*62))%62]
			b := chars[(x/62)%62]
			c := chars[x%62]
			name := string([]byte{'Z', a, b, c})
			if !used[name] {
				used[name] = true
				return name
			}
		}
	}

	dest := make([]string, len(nodes))
	visited := make([]bool, len(nodes))
	lines := make([]string, 0, len(nodes))

	dest[root] = "res"

	var nameOf func(int) string
	nameOf = func(id int) string {
		if nodes[id].leaf {
			return nodes[id].name
		}
		return dest[id]
	}

	var dfs func(int)
	dfs = func(id int) {
		if nodes[id].leaf || visited[id] {
			return
		}
		dfs(nodes[id].l)
		dfs(nodes[id].r)
		if dest[id] == "" {
			dest[id] = nextTemp()
		}
		lines = append(lines, fmt.Sprintf("%s=%s%c%s", dest[id], nameOf(nodes[id].l), nodes[id].op, nameOf(nodes[id].r)))
		visited[id] = true
	}

	dfs(root)

	fmt.Fprintln(w, len(lines))
	for _, line := range lines {
		fmt.Fprintln(w, line)
	}
	w.Flush()
	return strings.TrimSpace(out.String())
}

func run(bin, input string) (string, error) {
	cmd := exec.Command(bin)
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

var names = []string{"a", "b", "c", "d", "e", "f", "g", "h", "res"}
var ops = []byte{'$', '^', '#', '&'}

func randName(r *rand.Rand) string { return names[r.Intn(len(names))] }

func genCase(r *rand.Rand) string {
	n := r.Intn(5) + 1
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", n)
	for i := 0; i < n; i++ {
		l := randName(r)
		if r.Intn(2) == 0 {
			r1 := randName(r)
			fmt.Fprintf(&sb, "%s=%s\n", l, r1)
		} else {
			r1 := randName(r)
			r2 := randName(r)
			op := ops[r.Intn(len(ops))]
			fmt.Fprintf(&sb, "%s=%s%c%s\n", l, r1, op, r2)
		}
	}
	return sb.String()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierG.go /path/to/binary")
		os.Exit(1)
	}
	userBin := os.Args[1]
	r := rand.New(rand.NewSource(1))
	for i := 0; i < 100; i++ {
		input := genCase(r)
		want := oracleSolve(input)
		got, err := run(userBin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "test %d: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
		if want != got {
			fmt.Printf("test %d failed\ninput:\n%sexpected: %s\ngot: %s\n", i+1, input, want, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", 100)
}
