package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func runExe(path string, input []byte) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(path, ".go") {
		cmd = exec.Command("go", "run", path)
	} else {
		cmd = exec.Command(path)
	}
	cmd.Stdin = bytes.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

// Embedded correct solver for 1985H2
func solveReference(input []byte) string {
	data := input
	pos := 0
	n := len(data)

	skip := func() {
		for pos < n && data[pos] <= ' ' {
			pos++
		}
	}

	nextInt := func() int {
		skip()
		val := 0
		for pos < n {
			b := data[pos]
			if b <= ' ' {
				break
			}
			val = val*10 + int(b-'0')
			pos++
		}
		return val
	}

	nextBytes := func() []byte {
		skip()
		start := pos
		for pos < n && data[pos] > ' ' {
			pos++
		}
		return data[start:pos]
	}

	t := nextInt()
	var out bytes.Buffer

	for tc := 0; tc < t; tc++ {
		nr := nextInt()
		m := nextInt()
		nm := nr * m

		grid := make([]byte, nm)
		rowDots := make([]int, nr)
		colDots := make([]int, m)

		for i := 0; i < nr; i++ {
			s := nextBytes()
			base := i * m
			copy(grid[base:base+m], s)
			for j := 0; j < m; j++ {
				if s[j] == '.' {
					rowDots[i]++
					colDots[j]++
				}
			}
		}

		rowDiff := make([]int, nr+1)
		colDiff := make([]int, m+1)
		w := m + 1
		diff := make([]int, (nr+1)*w)
		queue := make([]int, nm)

		for i := 0; i < nr; i++ {
			base := i * m
			for j := 0; j < m; j++ {
				start := base + j
				if grid[start] != '#' {
					continue
				}

				grid[start] = '*'
				head, tail := 0, 1
				queue[0] = start

				size := 0
				minR, maxR := i, i
				minC, maxC := j, j

				for head < tail {
					p := queue[head]
					head++
					size++

					r := p / m
					c := p % m

					if r < minR {
						minR = r
					}
					if r > maxR {
						maxR = r
					}
					if c < minC {
						minC = c
					}
					if c > maxC {
						maxC = c
					}

					if r > 0 {
						q := p - m
						if grid[q] == '#' {
							grid[q] = '*'
							queue[tail] = q
							tail++
						}
					}
					if r+1 < nr {
						q := p + m
						if grid[q] == '#' {
							grid[q] = '*'
							queue[tail] = q
							tail++
						}
					}
					if c > 0 {
						q := p - 1
						if grid[q] == '#' {
							grid[q] = '*'
							queue[tail] = q
							tail++
						}
					}
					if c+1 < m {
						q := p + 1
						if grid[q] == '#' {
							grid[q] = '*'
							queue[tail] = q
							tail++
						}
					}
				}

				r1 := minR - 1
				if r1 < 0 {
					r1 = 0
				}
				r2 := maxR + 1
				if r2 >= nr {
					r2 = nr - 1
				}
				c1 := minC - 1
				if c1 < 0 {
					c1 = 0
				}
				c2 := maxC + 1
				if c2 >= m {
					c2 = m - 1
				}

				rowDiff[r1] += size
				rowDiff[r2+1] -= size
				colDiff[c1] += size
				colDiff[c2+1] -= size

				diff[r1*w+c1] += size
				diff[(r2+1)*w+c1] -= size
				diff[r1*w+c2+1] -= size
				diff[(r2+1)*w+c2+1] += size
			}
		}

		rowSum := make([]int, nr)
		cur := 0
		for i := 0; i < nr; i++ {
			cur += rowDiff[i]
			rowSum[i] = cur
		}

		colSum := make([]int, m)
		cur = 0
		for j := 0; j < m; j++ {
			cur += colDiff[j]
			colSum[j] = cur
		}

		ans := 0
		for i := 0; i < nr; i++ {
			rowBase := i * w
			prevRowBase := (i - 1) * w
			gridBase := i * m
			for j := 0; j < m; j++ {
				idx := rowBase + j
				v := diff[idx]
				if i > 0 {
					v += diff[prevRowBase+j]
				}
				if j > 0 {
					v += diff[rowBase+j-1]
				}
				if i > 0 && j > 0 {
					v -= diff[prevRowBase+j-1]
				}
				diff[idx] = v

				score := rowDots[i] + colDots[j]
				if grid[gridBase+j] == '.' {
					score--
				}
				score += rowSum[i] + colSum[j] - v

				if score > ans {
					ans = score
				}
			}
		}

		out.WriteString(strconv.Itoa(ans))
		if tc+1 < t {
			out.WriteByte('\n')
		}
	}

	return out.String()
}

func buildCase(grid []string) []byte {
	n := len(grid)
	m := len(grid[0])
	var sb strings.Builder
	sb.WriteString("1\n")
	sb.WriteString(fmt.Sprintf("%d %d\n", n, m))
	for _, row := range grid {
		sb.WriteString(row)
		sb.WriteByte('\n')
	}
	return []byte(sb.String())
}

func genRandomCase(rng *rand.Rand) []byte {
	n := rng.Intn(4) + 1
	m := rng.Intn(4) + 1
	grid := make([]string, n)
	for i := 0; i < n; i++ {
		b := make([]byte, m)
		for j := 0; j < m; j++ {
			if rng.Intn(2) == 0 {
				b[j] = '.'
			} else {
				b[j] = '#'
			}
		}
		grid[i] = string(b)
	}
	return buildCase(grid)
}

func genTests() [][]byte {
	rng := rand.New(rand.NewSource(9))
	tests := [][]byte{
		buildCase([]string{"#"}),
		buildCase([]string{".", "#"}),
		buildCase([]string{"..", "##"}),
	}
	for len(tests) < 100 {
		tests = append(tests, genRandomCase(rng))
	}
	return tests
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierH2.go /path/to/binary")
		return
	}
	bin := os.Args[1]

	tests := genTests()
	for i, tc := range tests {
		exp := strings.TrimSpace(solveReference(tc))

		got, err := runExe(bin, tc)
		if err != nil {
			fmt.Printf("candidate runtime error on test %d: %v\n%s", i+1, err, got)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != exp {
			fmt.Printf("test %d failed\ninput:\n%sexpected:\n%s\ngot:\n%s\n", i+1, string(tc), exp, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
