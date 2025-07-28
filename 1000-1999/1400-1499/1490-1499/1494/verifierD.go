package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const numTestsD = 100

func prepareBinary(path string) (string, func(), error) {
	if strings.HasSuffix(path, ".go") {
		tmp := filepath.Join(os.TempDir(), "candD")
		cmd := exec.Command("go", "build", "-o", tmp, path)
		if out, err := cmd.CombinedOutput(); err != nil {
			return "", nil, fmt.Errorf("go build failed: %v: %s", err, out)
		}
		return tmp, func() { os.Remove(tmp) }, nil
	}
	return path, nil, nil
}

func buildOracle() (string, func(), error) {
	tmp := filepath.Join(os.TempDir(), "oracleD")
	cmd := exec.Command("go", "build", "-o", tmp, "1494D.go")
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", nil, fmt.Errorf("build oracle failed: %v: %s", err, out)
	}
	return tmp, func() { os.Remove(tmp) }, nil
}

func run(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var buf bytes.Buffer
	cmd.Stdout = &buf
	cmd.Stderr = &buf
	err := cmd.Run()
	return strings.TrimSpace(buf.String()), err
}

type node struct {
	id     int
	salary int
	parent *node
}

func buildTree(rng *rand.Rand, n int) ([]*node, *node) {
	nodes := make([]*node, n)
	maxSalary := 0
	for i := 0; i < n; i++ {
		salary := rng.Intn(50) + 1
		if salary > maxSalary {
			maxSalary = salary
		}
		nodes[i] = &node{id: i, salary: salary}
	}
	curID := n
	for len(nodes) > 1 {
		i := rng.Intn(len(nodes))
		j := rng.Intn(len(nodes) - 1)
		if j >= i {
			j++
		}
		a := nodes[i]
		b := nodes[j]
		salary := maxSalary + rng.Intn(50) + 1
		if salary > maxSalary {
			maxSalary = salary
		}
		p := &node{id: curID, salary: salary}
		curID++
		a.parent = p
		b.parent = p
		nodes[i] = p
		nodes[j] = nodes[len(nodes)-1]
		nodes = nodes[:len(nodes)-1]
	}
	return nodes[0:curID], nodes[0]
}

func lca(a, b *node) *node {
	seen := map[*node]struct{}{}
	for a != nil {
		seen[a] = struct{}{}
		a = a.parent
	}
	for b != nil {
		if _, ok := seen[b]; ok {
			return b
		}
		b = b.parent
	}
	return nil
}

func genCase(rng *rand.Rand) string {
	n := rng.Intn(5) + 2
	leaves := make([]*node, n)
	nodes, root := buildTree(rng, n)
	// collect leaves
	for _, nd := range nodes {
		if nd.id < n {
			leaves[nd.id] = nd
		}
	}
	// compute matrix
	mat := make([][]int, n)
	for i := 0; i < n; i++ {
		mat[i] = make([]int, n)
		for j := 0; j < n; j++ {
			if i == j {
				mat[i][j] = leaves[i].salary
			} else {
				mat[i][j] = lca(leaves[i], leaves[j]).salary
			}
		}
	}
	// build input
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			if j > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprintf("%d", mat[i][j]))
		}
		sb.WriteByte('\n')
	}
	_ = root
	return sb.String()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		return
	}
	bin, clean, err := prepareBinary(os.Args[1])
	if err != nil {
		fmt.Println("compile error:", err)
		return
	}
	if clean != nil {
		defer clean()
	}
	oracle, c2, err := buildOracle()
	if err != nil {
		fmt.Println(err)
		return
	}
	defer c2()
	rng := rand.New(rand.NewSource(3))
	for i := 0; i < numTestsD; i++ {
		input := genCase(rng)
		want, err := run(oracle, input)
		if err != nil {
			fmt.Printf("oracle runtime error on case %d: %v\n", i+1, err)
			return
		}
		got, err := run(bin, input)
		if err != nil {
			fmt.Printf("candidate runtime error on case %d: %v\ninput:\n%s", i+1, err, input)
			return
		}
		if want != got {
			fmt.Printf("case %d failed\ninput:\n%sexpected:\n%s\ngot:\n%s\n", i+1, input, want, got)
			return
		}
	}
	fmt.Println("All tests passed")
}
