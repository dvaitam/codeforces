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

const solverSrc = `package main

import (
	"bufio"
	"fmt"
	"math/bits"
	"os"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	buf := make([]byte, 1024*1024)
	scanner.Buffer(buf, 10*1024*1024)
	scanner.Split(bufio.ScanWords)

	scanInt := func() int {
		scanner.Scan()
		res := 0
		for _, b := range scanner.Bytes() {
			res = res*10 + int(b-'0')
		}
		return res
	}

	if !scanner.Scan() {
		return
	}
	nBytes := scanner.Bytes()
	n := 0
	for _, b := range nBytes {
		n = n*10 + int(b-'0')
	}
	m := scanInt()

	adjList := make([][]int, n)
	adjBits := make([][]uint64, n)
	words := (n + 63) / 64
	for i := 0; i < n; i++ {
		adjBits[i] = make([]uint64, words)
	}

	for i := 0; i < m; i++ {
		u := scanInt() - 1
		v := scanInt() - 1
		adjList[u] = append(adjList[u], v)
		adjList[v] = append(adjList[v], u)
		adjBits[u][v/64] |= 1 << (v % 64)
		adjBits[v][u/64] |= 1 << (u % 64)
	}

	A2 := make([]int32, n*n)
	for i := 0; i < n; i++ {
		for j := i; j < n; j++ {
			c := 0
			for w := 0; w < words; w++ {
				c += bits.OnesCount64(adjBits[i][w] & adjBits[j][w])
			}
			A2[i*n+j] = int32(c)
			A2[j*n+i] = int32(c)
		}
	}

	var C3 int64 = 0
	t := make([]int64, n)
	for i := 0; i < n; i++ {
		for _, k := range adjList[i] {
			t[i] += int64(A2[i*n+k])
		}
		t[i] /= 2
		C3 += t[i]
	}
	C3 /= 3

	var Tr int64 = 0
	for i := 0; i < n; i++ {
		idxI := i * n
		a2i := A2[idxI : idxI+n]
		for _, k := range adjList[i] {
			idxK := k * n
			a2k := A2[idxK : idxK+n]
			var sum int64 = 0
			for j := 0; j < n; j++ {
				sum += int64(a2i[j]) * int64(a2k[j])
			}
			Tr += sum
		}
	}

	var S int64 = 0
	for i := 0; i < n; i++ {
		S += int64(len(adjList[i])) * t[i]
	}

	ans := (Tr / 10) + 3*C3 - S
	fmt.Println(ans)
}
`

func buildSolver() (string, func(), error) {
	dir, err := os.MkdirTemp("", "verE51")
	if err != nil {
		return "", nil, err
	}
	cleanup := func() { os.RemoveAll(dir) }
	src := filepath.Join(dir, "solver.go")
	if err := os.WriteFile(src, []byte(solverSrc), 0644); err != nil {
		cleanup()
		return "", nil, err
	}
	bin := filepath.Join(dir, "solver")
	cmd := exec.Command("go", "build", "-o", bin, src)
	if out, err := cmd.CombinedOutput(); err != nil {
		cleanup()
		return "", nil, fmt.Errorf("build solver: %v\n%s", err, out)
	}
	return bin, cleanup, nil
}

func generateCase(rng *rand.Rand) string {
	n := rng.Intn(4) + 5
	edges := make([][2]int, 0)
	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			if rng.Float32() < 0.3 {
				edges = append(edges, [2]int{i, j})
			}
		}
	}
	m := len(edges)
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", n, m))
	for _, e := range edges {
		sb.WriteString(fmt.Sprintf("%d %d\n", e[0]+1, e[1]+1))
	}
	return sb.String()
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

func main() {
	if len(os.Args) != 2 && !(len(os.Args) == 3 && os.Args[1] == "--") {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[len(os.Args)-1]

	ref, cleanup, err := buildSolver()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer cleanup()

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in := generateCase(rng)
		expect, err := run(ref, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d ref failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
		got, err := run(bin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
		if got != expect {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, expect, got, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
