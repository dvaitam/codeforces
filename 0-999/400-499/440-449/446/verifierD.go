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

func buildReference() (string, error) {
	srcPath := os.Getenv("REFERENCE_SOURCE_PATH")
	if srcPath == "" {
		_, file, _, ok := runtime.Caller(0)
		if !ok {
			return "", fmt.Errorf("cannot determine verifier directory and REFERENCE_SOURCE_PATH not set")
		}
		srcPath = filepath.Join(filepath.Dir(file), "446D.go")
	}
	content, err := os.ReadFile(srcPath)
	if err != nil {
		return "", fmt.Errorf("read reference source: %v", err)
	}
	tmp, err := os.CreateTemp("", "446D-ref-*")
	if err != nil {
		return "", err
	}
	tmp.Close()
	if strings.Contains(string(content), "#include") {
		cppPath := tmp.Name() + ".cpp"
		if err := os.WriteFile(cppPath, content, 0644); err != nil {
			os.Remove(tmp.Name())
			return "", fmt.Errorf("write cpp: %v", err)
		}
		defer os.Remove(cppPath)
		cmd := exec.Command("g++", "-O2", "-o", tmp.Name(), cppPath)
		if out, err := cmd.CombinedOutput(); err != nil {
			os.Remove(tmp.Name())
			return "", fmt.Errorf("%v\n%s", err, string(out))
		}
	} else {
		cmd := exec.Command("go", "build", "-o", tmp.Name(), srcPath)
		if out, err := cmd.CombinedOutput(); err != nil {
			os.Remove(tmp.Name())
			return "", fmt.Errorf("%v\n%s", err, string(out))
		}
	}
	return tmp.Name(), nil
}

func generateCase(rng *rand.Rand) string {
	n := rng.Intn(4) + 2
	maxM := n * (n - 1) / 2
	m := n - 1 + rng.Intn(maxM-(n-1)+1)
	// create connected graph via tree
	edges := make([][2]int, 0, m)
	for i := 2; i <= n; i++ {
		u := rng.Intn(i-1) + 1
		edges = append(edges, [2]int{u, i})
	}
	used := make(map[[2]int]bool)
	for _, e := range edges {
		if e[0] > e[1] {
			e[0], e[1] = e[1], e[0]
		}
		used[[2]int{e[0], e[1]}] = true
	}
	for len(edges) < m {
		u := rng.Intn(n) + 1
		v := rng.Intn(n) + 1
		if u == v {
			continue
		}
		a := u
		b := v
		if a > b {
			a, b = b, a
		}
		key := [2]int{a, b}
		if used[key] {
			continue
		}
		used[key] = true
		edges = append(edges, [2]int{u, v})
	}
	k := int64(rng.Intn(4) + 2)
	a := make([]int, n+1)
	trapCnt := 1
	for i := 2; i < n; i++ {
		if rng.Intn(2) == 0 && trapCnt < 3 {
			a[i] = 1
			trapCnt++
		}
	}
	a[n] = 1
	a[1] = 0
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d %d\n", n, len(edges), k))
	for i := 1; i <= n; i++ {
		if i > 1 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprint(a[i]))
	}
	sb.WriteByte('\n')
	for _, e := range edges {
		sb.WriteString(fmt.Sprintf("%d %d\n", e[0], e[1]))
	}
	return sb.String()
}

func run(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	refBin, err := buildReference()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to build reference: %v\n", err)
		os.Exit(1)
	}
	defer os.Remove(refBin)

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in := generateCase(rng)
		refOut, err := run(refBin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference failed on case %d: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
		var exp float64
		if _, err := fmt.Sscan(refOut, &exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d: cannot parse reference output %q: %v\n", i+1, refOut, err)
			os.Exit(1)
		}
		out, err := run(bin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
		var got float64
		if _, err := fmt.Sscan(out, &got); err != nil {
			fmt.Fprintf(os.Stderr, "case %d: cannot parse output: %v\n", i+1, err)
			os.Exit(1)
		}
		if diff := got - exp; diff < -1e-4 || diff > 1e-4 {
			fmt.Fprintf(os.Stderr, "case %d: expected %.6f got %.6f\ninput:\n%s", i+1, exp, got, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
