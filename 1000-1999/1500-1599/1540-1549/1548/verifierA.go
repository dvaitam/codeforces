package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

func runProg(path, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(path, ".go") {
		cmd = exec.Command("go", "run", path)
	} else {
		cmd = exec.Command(path)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func runRef(input string) (string, error) {
	_, self, _, _ := runtime.Caller(0)
	dir := filepath.Dir(self)
	ref := filepath.Join(dir, "1548A.go")
	return runProg(ref, input)
}

type edge struct{ u, v int }

func genCase(rng *rand.Rand) string {
	n := rng.Intn(5) + 1
	maxEdges := n * (n - 1) / 2
	if maxEdges > 5 {
		maxEdges = 5
	}
	m := rng.Intn(maxEdges + 1)
	edges := make(map[edge]bool)
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d\n", n, m)
	for len(edges) < m {
		u := rng.Intn(n) + 1
		v := rng.Intn(n) + 1
		if u == v {
			continue
		}
		if u > v {
			u, v = v, u
		}
		e := edge{u, v}
		if edges[e] {
			continue
		}
		edges[e] = true
		fmt.Fprintf(&sb, "%d %d\n", u, v)
	}
	q := rng.Intn(10) + 1
	sb.WriteString(fmt.Sprintf("%d\n", q))
	hasType3 := false
	for i := 0; i < q; i++ {
		t := rng.Intn(3) + 1
		if i == q-1 && !hasType3 {
			t = 3
		}
		if t == 1 && len(edges) == n*(n-1)/2 {
			t = 3
		}
		if t == 2 && len(edges) == 0 {
			t = 3
		}
		switch t {
		case 1:
			for {
				u := rng.Intn(n) + 1
				v := rng.Intn(n) + 1
				if u == v {
					continue
				}
				if u > v {
					u, v = v, u
				}
				e := edge{u, v}
				if edges[e] {
					continue
				}
				edges[e] = true
				fmt.Fprintf(&sb, "1 %d %d\n", u, v)
				break
			}
		case 2:
			idx := rng.Intn(len(edges))
			var chosen edge
			j := 0
			for e := range edges {
				if j == idx {
					chosen = e
					break
				}
				j++
			}
			delete(edges, chosen)
			fmt.Fprintf(&sb, "2 %d %d\n", chosen.u, chosen.v)
		case 3:
			fmt.Fprintln(&sb, "3")
			hasType3 = true
		}
	}
	return sb.String()
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	if bin == "--" && len(os.Args) >= 3 {
		bin = os.Args[2]
	}
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in := genCase(rng)
		expect, err := runRef(in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference failed on test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		got, err := runProg(bin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate failed on test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != strings.TrimSpace(expect) {
			fmt.Printf("test %d failed\ninput:\n%sexpected:%s\nactual:%s\n", i+1, in, expect, got)
			os.Exit(1)
		}
	}
	fmt.Println("all tests passed")
}
