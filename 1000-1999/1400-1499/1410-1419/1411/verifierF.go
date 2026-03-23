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

const oracleSource = `package main

import (
	"bufio"
	"fmt"
	"os"
)

func powMod(a, b int64) int64 {
	res := int64(1)
	a %= 1000000007
	for b > 0 {
		if b%2 == 1 {
			res = (res * a) % 1000000007
		}
		a = (a * a) % 1000000007
		b /= 2
	}
	return res
}

func f(x, y int) int {
	if x < y {
		return x + (y-x)/3
	}
	return y + (x-y)/3
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var t int
	if _, err := fmt.Fscan(reader, &t); err != nil {
		return
	}

	for i := 0; i < t; i++ {
		var n int
		fmt.Fscan(reader, &n)
		p := make([]int, n+1)
		for j := 1; j <= n; j++ {
			fmt.Fscan(reader, &p[j])
		}

		visited := make([]bool, n+1)
		c0, c1, cg1, c2, cg2 := 0, 0, 0, 0, 0
		A_len := 0

		for j := 1; j <= n; j++ {
			if !visited[j] {
				length := 0
				curr := j
				for !visited[curr] {
					visited[curr] = true
					curr = p[curr]
					length++
				}
				A_len++
				if length%3 == 0 {
					c0++
				} else if length%3 == 1 {
					if length == 1 {
						c1++
					} else {
						cg1++
					}
				} else if length%3 == 2 {
					if length == 2 {
						c2++
					} else {
						cg2++
					}
				}
			}
		}

		tot1 := c1 + cg1
		tot2 := c2 + cg2

		MOD := int64(1000000007)
		var days int64
		var minSwaps int

		if n%3 == 0 {
			days = powMod(3, int64(n/3))
			C := c0 + f(tot1, tot2)
			minSwaps = A_len + n/3 - 2*C
		} else if n%3 == 2 {
			days = (2 * powMod(3, int64((n-2)/3))) % MOD
			maxC := -1
			if tot2 >= 1 {
				cand := c0 + 1 + f(tot1, tot2-1)
				if cand > maxC {
					maxC = cand
				}
			}
			if tot1 >= 2 {
				cand := c0 + 1 + f(tot1-2, tot2)
				if cand > maxC {
					maxC = cand
				}
			}
			minSwaps = A_len + (n+1)/3 - 2*maxC
		} else {
			days = (4 * powMod(3, int64((n-4)/3))) % MOD
			swaps1 := int(1e9)
			swaps2 := int(1e9)

			if tot2 >= 2 {
				C := c0 + 2 + f(tot1, tot2-2)
				cand := A_len + (n+2)/3 - 2*C
				if cand < swaps1 {
					swaps1 = cand
				}
			}
			if tot1 >= 2 && tot2 >= 1 {
				C := c0 + 2 + f(tot1-2, tot2-1)
				cand := A_len + (n+2)/3 - 2*C
				if cand < swaps1 {
					swaps1 = cand
				}
			}
			if tot1 >= 4 {
				C := c0 + 2 + f(tot1-4, tot2)
				cand := A_len + (n+2)/3 - 2*C
				if cand < swaps1 {
					swaps1 = cand
				}
			}

			if cg1 >= 1 {
				C := c0 + 1 + f(tot1-1, tot2)
				cand := A_len + (n-1)/3 - 2*C
				if cand < swaps2 {
					swaps2 = cand
				}
			}
			if tot2 >= 2 {
				C := c0 + 1 + f(tot1, tot2-2)
				cand := A_len + (n-1)/3 - 2*C
				if cand < swaps2 {
					swaps2 = cand
				}
			}
			if c1 >= 1 {
				C := c0 + f(tot1-1, tot2)
				cand := A_len + (n-1)/3 - 2*C
				if cand < swaps2 {
					swaps2 = cand
				}
			}

			if swaps1 < swaps2 {
				minSwaps = swaps1
			} else {
				minSwaps = swaps2
			}
		}

		fmt.Fprintf(writer, "%d %d\n", days, minSwaps)
	}
}
`

func buildOracle() (string, error) {
	dir := os.TempDir()
	src := filepath.Join(dir, fmt.Sprintf("oracle1411F_%d.go", time.Now().UnixNano()))
	if err := os.WriteFile(src, []byte(oracleSource), 0644); err != nil {
		return "", fmt.Errorf("write oracle source: %v", err)
	}
	defer os.Remove(src)
	oracle := src[:len(src)-3]
	cmd := exec.Command("go", "build", "-o", oracle, src)
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("build oracle failed: %v\n%s", err, out)
	}
	return oracle, nil
}

func runBinary(path, input string) (string, error) {
	cmd := exec.Command(path)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

func genCase(rng *rand.Rand) string {
	t := rng.Intn(3) + 1
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", t))
	for caseNum := 0; caseNum < t; caseNum++ {
		n := rng.Intn(8) + 1
		sb.WriteString(fmt.Sprintf("%d\n", n))
		p := rng.Perm(n)
		for i := 0; i < n; i++ {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprintf("%d", p[i]+1))
		}
		sb.WriteByte('\n')
	}
	return sb.String()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	candPath := os.Args[1]

	oracle, err := buildOracle()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer os.Remove(oracle)

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		input := genCase(rng)
		exp, err := runBinary(oracle, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference error on case %d: %v\n", i+1, err)
			os.Exit(1)
		}
		got, err := runBinary(candPath, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate runtime error on case %d: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
		if strings.TrimSpace(exp) != strings.TrimSpace(got) {
			fmt.Printf("case %d failed\ninput:\n%sexpected:\n%s\ngot:\n%s", i+1, input, exp, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
