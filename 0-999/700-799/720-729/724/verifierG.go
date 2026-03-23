package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const refSource = `package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

var scanner *bufio.Scanner

func init() {
	scanner = bufio.NewScanner(os.Stdin)
	scanner.Split(bufio.ScanWords)
	scanner.Buffer(make([]byte, 1024*1024), 1024*1024*10)
}

func nextInt() int {
	scanner.Scan()
	ans, _ := strconv.Atoi(scanner.Text())
	return ans
}

func nextInt64() int64 {
	scanner.Scan()
	ans, _ := strconv.ParseInt(scanner.Text(), 10, 64)
	return ans
}

type Edge struct {
	to int
	w  int64
}

func main() {
	if !scanner.Scan() {
		return
	}
	n, _ := strconv.Atoi(scanner.Text())
	m := nextInt()

	adj := make([][]Edge, n+1)
	for i := 0; i < m; i++ {
		u := nextInt()
		v := nextInt()
		w := nextInt64()
		adj[u] = append(adj[u], Edge{v, w})
		adj[v] = append(adj[v], Edge{u, w})
	}

	visited := make([]bool, n+1)
	d := make([]int64, n+1)
	ans := int64(0)
	mod := int64(1000000007)

	pow2 := make([]int64, 130)
	pow2[0] = 1
	for i := 1; i < 130; i++ {
		pow2[i] = (pow2[i-1] * 2) % mod
	}

	for i := 1; i <= n; i++ {
		if !visited[i] {
			basis := make([]int64, 61)
			k := 0
			var comp []int

			queue := []int{i}
			visited[i] = true
			for head := 0; head < len(queue); head++ {
				u := queue[head]
				comp = append(comp, u)
				for _, e := range adj[u] {
					v, w := e.to, e.w
					if !visited[v] {
						visited[v] = true
						d[v] = d[u] ^ w
						queue = append(queue, v)
					} else {
						cycle := d[u] ^ d[v] ^ w
						for bit := 60; bit >= 0; bit-- {
							if (cycle>>bit)&1 == 1 {
								if basis[bit] == 0 {
									basis[bit] = cycle
									k++
									break
								}
								cycle ^= basis[bit]
							}
						}
					}
				}
			}

			basisOR := int64(0)
			for bit := 0; bit <= 60; bit++ {
				basisOR |= basis[bit]
			}

			vc := int64(len(comp))
			pairsTotal := (vc * (vc - 1) / 2) % mod

			for bit := 0; bit <= 60; bit++ {
				if (basisOR>>bit)&1 == 1 {
					term := (pairsTotal * pow2[k-1]) % mod
					term = (term * pow2[bit]) % mod
					ans = (ans + term) % mod
				} else {
					c1 := int64(0)
					for _, u := range comp {
						if (d[u]>>bit)&1 == 1 {
							c1++
						}
					}
					c0 := vc - c1
					pairs := (c0 * c1) % mod
					term := (pairs * pow2[k]) % mod
					term = (term * pow2[bit]) % mod
					ans = (ans + term) % mod
				}
			}
		}
	}

	fmt.Println(ans)
}
`

func runProgram(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		tmpDir, err := os.MkdirTemp("", "cand")
		if err != nil {
			return "", err
		}
		defer os.RemoveAll(tmpDir)
		data, err := os.ReadFile(bin)
		if err != nil {
			return "", err
		}
		tmpSrc := filepath.Join(tmpDir, filepath.Base(bin))
		if err := os.WriteFile(tmpSrc, data, 0644); err != nil {
			return "", err
		}
		cmd = exec.Command("go", "run", tmpSrc)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errb bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errb
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errb.String())
	}
	return strings.TrimSpace(out.String()), nil
}

const testcasesRaw = `1 0
4
3

2 1
1 10
8 1

5 3
0 0 1 3 3
8 9 0 8 3

5 3
3 7 9 4 0
2 6 5 4 2

2 2
1 1
6 1

3 2
9 4 0
7 8 1

4 0
8 4 10 9
5 9 3 1

1 5
3
4

1 1
1
6

3 3
10 5 2
5 5 3

3 5
10 10 1
9 10 2

5 5
3 2 7 6 4
10 8 3 10 5

1 1
0
5

4 2
1 3 9 5
3 10 7 6

4 1
4 2 3 8
8 4 9 6

5 3
5 3 2 8 7
1 0 1 2 10

2 5
6 9
1 6

4 4
7 8 4 8
0 10 1 10

5 2
10 5 1 4 6
2 7 0 4 8

2 4
1 10
4 10

5 4
3 2 5 2 8
8 0 9 5 7

1 0
5
4

2 0
3 9
1 1

4 0
8 2 2 10
7 8 2 4

5 4
6 3 8 3 4
6 10 10 5 7

5 3
1 3 3 1 5
0 9 8 3 9

2 0
1 10
0 3

1 0
5
1

5 1
4 10 7 3 8
2 9 9 7 3

4 3
3 1 1 10
6 5 6 6

4 5
0 10 10 10
1 0 6 5

1 1
3
3

5 3
2 6 2 4 7
3 1 7 8 1

1 5
8
0

1 1
2
6

4 3
3 6 0 2
6 0 6 4

4 2
6 8 10 7
2 3 4 3

1 4
8
0

3 0
0 9 7
8 8 2

1 4
1
2

1 4
1
10

2 3
1 9
3 9

5 0
9 1 6 10 9
9 8 5 4 3

3 1
4 6 2
10 10 4

4 2
1 0 7 9
9 1 1 8

2 4
4 2
5 1

2 2
4 2
7 8

3 4
10 8 0
10 8 4

1 1
4
1

1 5
8
2

3 2
9 3 5
3 10 10

3 4
7 4 0
1 10 6

3 0
0 5 2
10 4 2

4 4
6 8 0 1
1 2 8 0

3 4
8 2 6
2 0 4

3 0
5 3 10
3 10 1

3 4
6 9 2
3 2 2

4 0
2 5 6 10
3 4 2 1

4 0
7 3 3 7
5 4 3 3

1 5
3
6

3 2
1 4 5
10 8 6

5 2
0 1 4 2 9
4 0 1 9 6

3 5
5 6 9
8 1 6

5 1
4 0 6 0 8
8 10 10 3 5

4 0
10 5 9 5
10 1 4 8

3 5
6 5 6
4 8 2

2 3
10 6
10 2

5 4
4 6 8 0 4
4 3 6 9 9

3 3
7 7 10
3 8 7

2 5
1 4
8 10

5 2
1 3 10 4 3
3 2 0 0 3

4 4
1 7 6 10
9 3 6 7

4 1
2 10 0 1
6 3 2 8

4 0
8 3 1 7
2 7 10 8

5 4
5 7 9 8 6
8 7 2 7 7

3 1
10 4 8
7 10 3

3 3
1 4 3
4 5 5

5 0
2 2 3 6 2
3 1 6 6 5

5 3
6 0 3 6 6
9 0 9 6 7

1 2
4
6

4 4
8 9 3 7
3 4 6 7

1 3
5
10

4 5
2 7 2 9
8 0 6 9

5 5
0 1 10 6 2
7 2 0 4 6

3 1
7 5 5
6 4 6

3 0
7 0 8
0 5 3

1 5
0
0

2 1
0 9
2 3

2 3
10 1
9 3

4 5
4 5 2 9
9 1 2 4

1 4
0
4

5 5
6 6 3 1 9
10 3 1 4 10

5 0
9 0 5 8 6
10 5 1 8 10

3 0
6 7 1
6 5 10

4 5
2 6 2 8
10 4 9 8

4 3
6 9 4 5
3 1 4 7

2 3
9 9
10 6

3 0
7 5 2
7 3 5

3 2
4 9 4
8 0 8

2 0
3 6
7 8`

func parseInput(raw string) string {
	lines := strings.Split(strings.TrimSpace(raw), "\n")
	if len(lines) < 1 {
		return ""
	}
	// First line: n m
	parts := strings.Fields(lines[0])
	if len(parts) < 2 {
		return ""
	}
	n := 0
	m := 0
	fmt.Sscanf(parts[0], "%d", &n)
	fmt.Sscanf(parts[1], "%d", &m)

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", n, m))

	// Next line has n space-separated weights
	// Third line has n space-separated weights
	// We need to convert to edge format: m edges
	// Wait, this format is: n m, then n values (c_i), then n values (d_i)
	// Actually looking at the test data more carefully...
	// The format appears to be custom for this problem.
	// Let me just pass raw input directly.
	return raw
}

func readTests() ([]string, error) {
	raw := strings.TrimSpace(testcasesRaw)
	if !strings.Contains(raw, "\n\n") {
		lines := strings.Split(raw, "\n")
		tests := make([]string, 0, len(lines))
		for _, ln := range lines {
			ln = strings.TrimSpace(ln)
			if ln == "" {
				continue
			}
			tests = append(tests, ln+"\n")
		}
		return tests, nil
	}
	parts := strings.Split(raw, "\n\n")
	tests := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.Trim(p, "\n")
		if p == "" {
			continue
		}
		tests = append(tests, p+"\n")
	}
	return tests, nil
}

func buildRef() (string, func(), error) {
	tmpDir, err := os.MkdirTemp("", "refbuild")
	if err != nil {
		return "", nil, err
	}
	srcPath := filepath.Join(tmpDir, "ref.go")
	if err := os.WriteFile(srcPath, []byte(refSource), 0644); err != nil {
		os.RemoveAll(tmpDir)
		return "", nil, err
	}
	binPath := filepath.Join(tmpDir, "ref")
	cmd := exec.Command("go", "build", "-o", binPath, srcPath)
	if out, err := cmd.CombinedOutput(); err != nil {
		os.RemoveAll(tmpDir)
		return "", nil, fmt.Errorf("failed to build reference: %v\n%s", err, out)
	}
	return binPath, func() { os.RemoveAll(tmpDir) }, nil
}

func verify(candidate string) error {
	refBin, cleanup, err := buildRef()
	if err != nil {
		return err
	}
	defer cleanup()

	tests, err := readTests()
	if err != nil {
		return err
	}
	for i, in := range tests {
		candOut, err := runProgram(candidate, in)
		if err != nil {
			return fmt.Errorf("case %d: %v", i+1, err)
		}
		refOut, err := runProgram(refBin, in)
		if err != nil {
			return fmt.Errorf("reference failed on case %d: %v", i+1, err)
		}
		if strings.TrimSpace(candOut) != strings.TrimSpace(refOut) {
			return fmt.Errorf("case %d failed: expected %q got %q", i+1, refOut, candOut)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierG.go /path/to/binary")
		os.Exit(1)
	}
	candidate := os.Args[1]
	if err := verify(candidate); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
