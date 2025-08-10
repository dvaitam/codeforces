package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"
)

func buildRef() (string, error) {
	_, file, _, _ := runtime.Caller(0)
	dir := filepath.Dir(file)
	ref := filepath.Join(dir, "refA.bin")
	cmd := exec.Command("go", "build", "-o", ref, filepath.Join(dir, "623A.go"))
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("failed to build reference: %v\n%s", err, out)
	}
	return ref, nil
}

func generateCase(rng *rand.Rand) string {
	n := rng.Intn(8) + 1
	var edges [][2]int
	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			if rng.Intn(2) == 0 {
				edges = append(edges, [2]int{i + 1, j + 1})
			}
		}
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", n, len(edges)))
	for _, e := range edges {
		sb.WriteString(fmt.Sprintf("%d %d\n", e[0], e[1]))
	}
	return sb.String()
}

func runCase(bin, ref, input string) error {
	// parse input graph
	sc := bufio.NewScanner(strings.NewReader(input))
	sc.Split(bufio.ScanWords)
	sc.Scan()
	n, _ := strconv.Atoi(sc.Text())
	sc.Scan()
	m, _ := strconv.Atoi(sc.Text())
	adj := make([][]bool, n)
	for i := 0; i < n; i++ {
		adj[i] = make([]bool, n)
	}
	for i := 0; i < m; i++ {
		sc.Scan()
		u, _ := strconv.Atoi(sc.Text())
		sc.Scan()
		v, _ := strconv.Atoi(sc.Text())
		u--
		v--
		adj[u][v] = true
		adj[v][u] = true
	}

	// run reference to know correct verdict
	cmdRef := exec.Command(ref)
	cmdRef.Stdin = strings.NewReader(input)
	var refOut bytes.Buffer
	cmdRef.Stdout = &refOut
	cmdRef.Stderr = &refOut
	if err := cmdRef.Run(); err != nil {
		return fmt.Errorf("reference runtime error: %v\n%s", err, refOut.String())
	}
	refLines := strings.Split(strings.TrimSpace(refOut.String()), "\n")
	refAns := refLines[0]

	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	gotLines := strings.Split(strings.TrimSpace(out.String()), "\n")
	ans := gotLines[0]
	if refAns == "No" {
		if ans != "No" {
			return fmt.Errorf("expected No got %q", strings.Join(gotLines, "\n"))
		}
		return nil
	}
	if ans != "Yes" {
		return fmt.Errorf("expected Yes got %q", strings.Join(gotLines, "\n"))
	}
	if len(gotLines) < 2 {
		return fmt.Errorf("missing labeling line")
	}
	lab := gotLines[1]
	if len(lab) != n {
		return fmt.Errorf("labeling has wrong length")
	}
	for i := 0; i < n; i++ {
		c := lab[i]
		if c != 'a' && c != 'b' && c != 'c' {
			return fmt.Errorf("invalid character %q", c)
		}
	}
	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			diff := int(lab[i]) - int(lab[j])
			if diff < 0 {
				diff = -diff
			}
			if diff > 1 {
				if adj[i][j] {
					return fmt.Errorf("edge between %d and %d not allowed", i+1, j+1)
				}
			} else {
				if !adj[i][j] {
					return fmt.Errorf("missing edge between %d and %d", i+1, j+1)
				}
			}
		}
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	ref, err := buildRef()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer os.Remove(ref)
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in := generateCase(rng)
		if err := runCase(bin, ref, in); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
