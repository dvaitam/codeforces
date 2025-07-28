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

type edge struct{ u, v int }

func randTree(n int) []edge {
	edges := make([]edge, 0, n-1)
	for i := 2; i <= n; i++ {
		p := rand.Intn(i-1) + 1
		edges = append(edges, edge{i, p})
	}
	return edges
}

func genInput() string {
	t := 100
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", t))
	for i := 0; i < t; i++ {
		n := rand.Intn(7) + 2
		sb.WriteString(fmt.Sprintf("%d\n", n))
		for _, e := range randTree(n) {
			sb.WriteString(fmt.Sprintf("%d %d\n", e.u, e.v))
		}
	}
	return sb.String()
}

func runBin(path, in string) ([]byte, error) {
	cmd := exec.Command(path)
	cmd.Stdin = strings.NewReader(in)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	return out.Bytes(), err
}

func main() {
	rand.Seed(time.Now().UnixNano())
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	user := os.Args[1]
	ref := "./refD.bin"
	if err := exec.Command("go", "build", "-o", ref, "1667D.go").Run(); err != nil {
		fmt.Println("failed to build reference:", err)
		os.Exit(1)
	}
	defer os.Remove(ref)
	input := genInput()
	exp, err := runBin(ref, input)
	if err != nil {
		fmt.Println("reference failed:", err)
		os.Exit(1)
	}
	out, err := runBin(user, input)
	if err != nil {
		fmt.Println("runtime error:", err)
		os.Exit(1)
	}
	if strings.TrimSpace(string(out)) != strings.TrimSpace(string(exp)) {
		fmt.Printf("wrong answer\ninput:\n%s\nexpected:\n%s\nfound:\n%s\n", input, exp, out)
		os.Exit(1)
	}
	fmt.Println("ok")
}
