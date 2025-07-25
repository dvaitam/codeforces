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
)

func runBinary(binPath, input string) (string, error) {
	cmd := exec.Command(binPath)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

func genTree(n int, rng *rand.Rand) [][2]int {
	edges := make([][2]int, 0, n-1)
	for i := 2; i <= n; i++ {
		p := rng.Intn(i-1) + 1
		edges = append(edges, [2]int{p, i})
	}
	return edges
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	candidate := os.Args[1]

	_, file, _, _ := runtime.Caller(0)
	dir := filepath.Dir(file)
	refSrc := filepath.Join(dir, "796C.go")
	refBin := filepath.Join(os.TempDir(), "796C_ref_bin")
	cmd := exec.Command("go", "build", "-o", refBin, refSrc)
	if out, err := cmd.CombinedOutput(); err != nil {
		fmt.Println("Failed to build reference:", err)
		fmt.Println(string(out))
		os.Exit(1)
	}
	defer os.Remove(refBin)

	rng := rand.New(rand.NewSource(44))

	for tc := 1; tc <= 100; tc++ {
		n := rng.Intn(14) + 2
		a := make([]int, n)
		for i := range a {
			a[i] = rng.Intn(21) - 10
		}
		edges := genTree(n, rng)

		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d\n", n))
		for i, v := range a {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprintf("%d", v))
		}
		sb.WriteByte('\n')
		for _, e := range edges {
			sb.WriteString(fmt.Sprintf("%d %d\n", e[0], e[1]))
		}
		input := sb.String()

		exp, err := runBinary(refBin, input)
		if err != nil {
			fmt.Println("Reference solution error on test", tc)
			fmt.Println(err)
			os.Exit(1)
		}
		got, err := runBinary(candidate, input)
		if err != nil {
			fmt.Printf("Candidate runtime error on test %d: %v\n", tc, err)
			os.Exit(1)
		}
		if strings.TrimSpace(exp) != strings.TrimSpace(got) {
			fmt.Printf("Wrong answer on test %d\nInput:\n%sExpected:%sGot:%s\n", tc, input, exp, got)
			os.Exit(1)
		}
	}

	fmt.Println("All tests passed.")
}
