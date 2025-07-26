package main

import (
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

func buildRef() (string, error) {
	ref := "refF2.bin"
	cmd := exec.Command("go", "build", "-o", ref, "1118F2.go")
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("compile ref: %v %s", err, string(out))
	}
	return ref, nil
}

func runBinary(bin string, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	out, err := cmd.CombinedOutput()
	return string(out), err
}

func randomTree(n int) [][2]int {
	edges := make([][2]int, 0, n-1)
	for i := 2; i <= n; i++ {
		p := rand.Intn(i-1) + 1
		edges = append(edges, [2]int{i, p})
	}
	return edges
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierF2.go /path/to/binary")
		return
	}
	bin := os.Args[1]
	ref, err := buildRef()
	if err != nil {
		fmt.Println(err)
		return
	}
	defer os.Remove(ref)
	rand.Seed(8)
	for t := 1; t <= 100; t++ {
		n := rand.Intn(8) + 2
		k := rand.Intn(n-1) + 2
		colors := make([]int, n)
		for i := 0; i < k; i++ {
			colors[i] = i + 1
		}
		for i := k; i < n; i++ {
			colors[i] = rand.Intn(k + 1)
		}
		edges := randomTree(n)
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d %d\n", n, k))
		for i := 0; i < n; i++ {
			sb.WriteString(fmt.Sprintf("%d ", colors[i]))
		}
		sb.WriteString("\n")
		for _, e := range edges {
			sb.WriteString(fmt.Sprintf("%d %d\n", e[0], e[1]))
		}
		expOut, err := runBinary(ref, sb.String())
		if err != nil {
			fmt.Printf("test %d reference runtime error: %v\n", t, err)
			fmt.Println(expOut)
			return
		}
		out, err := runBinary(bin, sb.String())
		if err != nil {
			fmt.Printf("test %d runtime error: %v\n", t, err)
			fmt.Println(out)
			return
		}
		if strings.TrimSpace(out) != strings.TrimSpace(expOut) {
			fmt.Printf("test %d failed\ninput:\n%sexpected %s got %s\n", t, sb.String(), strings.TrimSpace(expOut), strings.TrimSpace(out))
			return
		}
	}
	fmt.Println("all tests passed")
}
