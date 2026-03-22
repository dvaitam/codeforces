package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strings"
	"time"
)

// ---- Embedded solver for 243E ----

func solve243E(input string) string {
	reader := bufio.NewReader(strings.NewReader(input))
	var n int
	fmt.Fscanf(reader, "%d\n", &n)

	mat := make([]string, n)
	for i := 0; i < n; i++ {
		fmt.Fscanf(reader, "%s\n", &mat[i])
	}

	w := make([][]int, n)
	for i := 0; i < n; i++ {
		w[i] = make([]int, n)
	}
	S := make([]int, n)

	for i := 0; i < n; i++ {
		for j := i; j < n; j++ {
			c := 0
			for r := 0; r < n; r++ {
				if mat[r][i] == '1' && mat[r][j] == '1' {
					c++
				}
			}
			w[i][j] = c
			w[j][i] = c
		}
	}

	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			if i != j {
				S[i] += w[i][j]
			}
		}
	}

	visited := make([]bool, n)
	var components [][]int

	for i := 0; i < n; i++ {
		if !visited[i] {
			comp := []int{}
			q := []int{i}
			visited[i] = true
			for len(q) > 0 {
				u := q[0]
				q = q[1:]
				comp = append(comp, u)
				for v := 0; v < n; v++ {
					if !visited[v] && w[u][v] > 0 {
						visited[v] = true
						q = append(q, v)
					}
				}
			}
			components = append(components, comp)
		}
	}

	finalOrder := []int{}

	for _, comp := range components {
		if len(comp) <= 2 {
			finalOrder = append(finalOrder, comp...)
			continue
		}

		found := false
	search:
		for _, E1 := range comp {
			minW := int(1e9)
			for _, v := range comp {
				if v != E1 && w[E1][v] < minW {
					minW = w[E1][v]
				}
			}

			candidates := []int{}
			for _, v := range comp {
				if v != E1 && w[E1][v] == minW {
					candidates = append(candidates, v)
				}
			}

			sort.Slice(candidates, func(i, j int) bool {
				return S[candidates[i]] < S[candidates[j]]
			})

			limit := 5
			if len(candidates) < limit {
				limit = len(candidates)
			}

			for _, E2 := range candidates[:limit] {
				cCopy := make([]int, len(comp))
				copy(cCopy, comp)

				sort.Slice(cCopy, func(i, j int) bool {
					u, v := cCopy[i], cCopy[j]
					d1 := w[E1][u] - w[E2][u]
					d2 := w[E1][v] - w[E2][v]
					if d1 != d2 {
						return d1 > d2
					}
					if S[u] != S[v] {
						return S[u] < S[v]
					}
					for k := 0; k < n; k++ {
						if w[u][k] != w[v][k] {
							return w[u][k] < w[v][k]
						}
					}
					return u < v
				})

				valid := true
				for i := 0; i < n; i++ {
					started := false
					ended := false
					for _, c := range cCopy {
						if mat[i][c] == '1' {
							if ended {
								valid = false
								break
							}
							started = true
						} else {
							if started {
								ended = true
							}
						}
					}
					if !valid {
						break
					}
				}

				if valid {
					finalOrder = append(finalOrder, cCopy...)
					found = true
					break search
				}
			}
		}

		if !found {
			return "NO"
		}
	}

	for i := 0; i < n; i++ {
		started := false
		ended := false
		for _, c := range finalOrder {
			if mat[i][c] == '1' {
				if ended {
					return "NO"
				}
				started = true
			} else {
				if started {
					ended = true
				}
			}
		}
	}

	var sb strings.Builder
	sb.WriteString("YES\n")
	for i := 0; i < n; i++ {
		for _, c := range finalOrder {
			sb.WriteByte(mat[i][c])
		}
		sb.WriteByte('\n')
	}
	return strings.TrimSpace(sb.String())
}

// ---- Verifier harness ----

func runBinary(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

func genTest() string {
	n := rand.Intn(5) + 1
	var b strings.Builder
	fmt.Fprintf(&b, "%d\n", n)
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			if rand.Intn(2) == 0 {
				b.WriteByte('0')
			} else {
				b.WriteByte('1')
			}
		}
		b.WriteByte('\n')
	}
	return b.String()
}

func buildIfGo(path string) (string, func(), error) {
	if strings.HasSuffix(path, ".go") {
		tmp, err := os.CreateTemp("", "solbin*")
		if err != nil {
			return "", nil, err
		}
		tmp.Close()
		out, err := exec.Command("go", "build", "-o", tmp.Name(), path).CombinedOutput()
		if err != nil {
			os.Remove(tmp.Name())
			return "", nil, fmt.Errorf("build failed: %v\n%s", err, out)
		}
		return tmp.Name(), func() { os.Remove(tmp.Name()) }, nil
	}
	return path, func() {}, nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	candidate, _ := os.Getwd()
	_ = candidate
	rand.Seed(time.Now().UnixNano())

	bin, cleanup, err := buildIfGo(os.Args[1])
	if err != nil {
		fmt.Println("failed to build candidate:", err)
		os.Exit(1)
	}
	defer cleanup()

	for i := 0; i < 100; i++ {
		input := genTest()
		expected := strings.TrimSpace(solve243E(input))
		actOut, err := runBinary(bin, input)
		if err != nil {
			fmt.Printf("test %d: runtime error: %v\n", i+1, err)
			os.Exit(1)
		}
		actual := strings.TrimSpace(actOut)
		if actual != expected {
			fmt.Printf("test %d failed:\ninput:\n%s\nexpected:%s\ngot:%s\n", i+1, input, expected, actual)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed.")
}
