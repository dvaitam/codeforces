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

const numTestsE = 100

func prepareBinary(path string) (string, func(), error) {
	if strings.HasSuffix(path, ".go") {
		tmp := filepath.Join(os.TempDir(), "candE")
		cmd := exec.Command("go", "build", "-o", tmp, path)
		if out, err := cmd.CombinedOutput(); err != nil {
			return "", nil, fmt.Errorf("go build failed: %v: %s", err, out)
		}
		return tmp, func() { os.Remove(tmp) }, nil
	}
	return path, nil, nil
}

func buildOracle() (string, func(), error) {
	tmp := filepath.Join(os.TempDir(), "oracleE")
	cmd := exec.Command("go", "build", "-o", tmp, "1494E.go")
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

func genCase(rng *rand.Rand) string {
	n := rng.Intn(5) + 2
	m := rng.Intn(30) + 1
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", n, m))
	type edge struct{ u, v int }
	edges := make(map[edge]byte)
	letters := "abcdefghijklmnopqrstuvwxyz"
	queries := 0
	for i := 0; i < m; i++ {
		t := rng.Intn(3)
		if i == m-1 && queries == 0 {
			t = 2
		}
		switch t {
		case 0: // add
			u := rng.Intn(n) + 1
			v := rng.Intn(n-1) + 1
			if v >= u {
				v++
			}
			e := edge{u, v}
			if _, ok := edges[e]; ok {
				i--
				continue
			}
			c := letters[rng.Intn(len(letters))]
			edges[e] = c
			sb.WriteString(fmt.Sprintf("+ %d %d %c\n", u, v, c))
		case 1: // remove
			if len(edges) == 0 {
				t = 0
				i--
				continue
			}
			k := rng.Intn(len(edges))
			var e edge
			j := 0
			for x := range edges {
				if j == k {
					e = x
					break
				}
				j++
			}
			sb.WriteString(fmt.Sprintf("- %d %d\n", e.u, e.v))
			delete(edges, e)
		default: // query
			k := rng.Intn(8) + 2
			sb.WriteString(fmt.Sprintf("? %d\n", k))
			queries++
		}
	}
	return sb.String()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
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
	rng := rand.New(rand.NewSource(4))
	for i := 0; i < numTestsE; i++ {
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
