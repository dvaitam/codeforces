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

const refSource = `package main

import (
	"bufio"
	"fmt"
	"os"
)

func abs(a int) int {
	if a < 0 {
		return -a
	}
	return a
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Split(bufio.ScanWords)
	buf := make([]byte, 0, 1024*1024)
	scanner.Buffer(buf, 1024*1024)

	nextInt := func() int {
		scanner.Scan()
		res := 0
		for _, b := range scanner.Bytes() {
			res = res*10 + int(b-'0')
		}
		return res
	}

	n := nextInt()
	m := nextInt()
	k := nextInt()

	grid := make([]int8, n*m)
	color_cells := make([][]int, k+1)

	for i := 0; i < n*m; i++ {
		c := int8(nextInt())
		grid[i] = c
		color_cells[c] = append(color_cells[c], i)
	}

	D := make([][]int32, k+1)
	Q := make([]int, n*m)

	for c := int8(1); c <= int8(k); c++ {
		D[c] = make([]int32, n*m)
		for i := 0; i < n*m; i++ {
			D[c][i] = -1
		}
		visited_color := make([]bool, k+1)
		head, tail := 0, 0

		for _, v := range color_cells[c] {
			D[c][v] = 0
			Q[tail] = v
			tail++
		}
		visited_color[c] = true

		for head < tail {
			u := Q[head]
			head++
			d := D[c][u]

			u_col := grid[u]
			if !visited_color[u_col] {
				visited_color[u_col] = true
				for _, v := range color_cells[u_col] {
					if D[c][v] == -1 {
						D[c][v] = d + 1
						Q[tail] = v
						tail++
					}
				}
			}

			r, c_idx := u/m, u%m

			if r > 0 {
				v := u - m
				if D[c][v] == -1 {
					D[c][v] = d + 1
					Q[tail] = v
					tail++
				}
			}
			if r < n-1 {
				v := u + m
				if D[c][v] == -1 {
					D[c][v] = d + 1
					Q[tail] = v
					tail++
				}
			}
			if c_idx > 0 {
				v := u - 1
				if D[c][v] == -1 {
					D[c][v] = d + 1
					Q[tail] = v
					tail++
				}
			}
			if c_idx < m-1 {
				v := u + 1
				if D[c][v] == -1 {
					D[c][v] = d + 1
					Q[tail] = v
					tail++
				}
			}
		}
	}

	q := nextInt()
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	for i := 0; i < q; i++ {
		r1 := nextInt() - 1
		c1 := nextInt() - 1
		r2 := nextInt() - 1
		c2 := nextInt() - 1

		u := r1*m + c1
		v := r2*m + c2
		ans := int32(abs(r1-r2) + abs(c1-c2))

		for c := int8(1); c <= int8(k); c++ {
			cost := D[c][u] + D[c][v] + 1
			if cost < ans {
				ans = cost
			}
		}
		fmt.Fprintln(writer, ans)
	}
}
`

func genCase(r *rand.Rand) string {
	n := r.Intn(4) + 1
	m := r.Intn(4) + 1
	maxColors := n * m
	k := r.Intn(min(5, maxColors)) + 1
	grid := make([][]int, n)
	for i := range grid {
		grid[i] = make([]int, m)
	}
	// ensure each color appears at least once
	for c := 1; c <= k; c++ {
		i := r.Intn(n)
		j := r.Intn(m)
		grid[i][j] = c
	}
	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			if grid[i][j] == 0 {
				grid[i][j] = r.Intn(k) + 1
			}
		}
	}
	q := r.Intn(5) + 1
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d %d\n", n, m, k))
	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			if j > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprintf("%d", grid[i][j]))
		}
		sb.WriteByte('\n')
	}
	sb.WriteString(fmt.Sprintf("%d\n", q))
	for i := 0; i < q; i++ {
		r1 := r.Intn(n) + 1
		c1 := r.Intn(m) + 1
		r2 := r.Intn(n) + 1
		c2 := r.Intn(m) + 1
		sb.WriteString(fmt.Sprintf("%d %d %d %d\n", r1, c1, r2, c2))
	}
	return sb.String()
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func run(cmdPath, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(cmdPath, ".go") {
		cmd = exec.Command("go", "run", cmdPath)
	} else if strings.HasSuffix(cmdPath, ".py") {
		cmd = exec.Command("python3", cmdPath)
	} else {
		cmd = exec.Command(cmdPath)
	}
	cmd.Stdin = strings.NewReader(input)
	var out, stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func buildReferenceBinary() (string, error) {
	srcFile, err := os.CreateTemp("", "cf-1301F-src-*.go")
	if err != nil {
		return "", err
	}
	if _, err := srcFile.WriteString(refSource); err != nil {
		srcFile.Close()
		os.Remove(srcFile.Name())
		return "", err
	}
	srcFile.Close()
	defer os.Remove(srcFile.Name())

	tmp, err := os.CreateTemp("", "cf-1301F-ref-*")
	if err != nil {
		return "", err
	}
	tmp.Close()
	os.Remove(tmp.Name())

	cmd := exec.Command("go", "build", "-o", tmp.Name(), srcFile.Name())
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		os.Remove(tmp.Name())
		return "", fmt.Errorf("go build error: %v\n%s", err, stderr.String())
	}
	return tmp.Name(), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	oracle, err := buildReferenceBinary()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer os.Remove(oracle)

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 1; i <= 100; i++ {
		input := genCase(rng)
		expect, err := run(oracle, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "oracle error on case %d: %v\n", i, err)
			os.Exit(1)
		}
		got, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i, err, input)
			os.Exit(1)
		}
		if got != expect {
			fmt.Fprintf(os.Stderr, "case %d failed: expected\n%s\ngot\n%s\ninput:\n%s", i, expect, got, input)
			os.Exit(1)
		}
	}
	fmt.Println("All 100 tests passed")
}
