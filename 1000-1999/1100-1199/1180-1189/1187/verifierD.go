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

const INF = int(1e9)
const MAXN = 300000
const MAXSIZE = 1 << 19

// Embedded correct solver for 1187/D
func solveD(input string) string {
	data := []byte(input)
	idx := 0
	nextInt := func() int {
		n := len(data)
		for idx < n {
			c := data[idx]
			if c >= '0' && c <= '9' {
				break
			}
			idx++
		}
		v := 0
		for idx < n {
			c := data[idx]
			if c < '0' || c > '9' {
				break
			}
			v = v*10 + int(c-'0')
			idx++
		}
		return v
	}

	pos := make([][]int, MAXN+1)
	ptr := make([]int, MAXN+1)
	tree := make([]int, 2*MAXSIZE)
	usedBuf := make([]int, 0, MAXN)

	t := nextInt()
	var out bytes.Buffer

	for ; t > 0; t-- {
		n := nextInt()
		used := usedBuf[:0]

		for i := 1; i <= n; i++ {
			x := nextInt()
			if len(pos[x]) == 0 {
				used = append(used, x)
			}
			pos[x] = append(pos[x], i)
		}

		size := 1
		for size < n {
			size <<= 1
		}
		limit := size << 1
		for i := 1; i < limit; i++ {
			tree[i] = INF
		}
		for _, x := range used {
			tree[size+x-1] = pos[x][0]
		}
		for i := size - 1; i > 0; i-- {
			if tree[i<<1] < tree[i<<1|1] {
				tree[i] = tree[i<<1]
			} else {
				tree[i] = tree[i<<1|1]
			}
		}

		ok := true
		for i := 0; i < n; i++ {
			x := nextInt()
			if !ok {
				continue
			}
			if ptr[x] >= len(pos[x]) {
				ok = false
				continue
			}
			p := pos[x][ptr[x]]

			if x > 1 {
				l, r := size, size+x-1
				minv := INF
				for l < r {
					if l&1 == 1 {
						if tree[l] < minv {
							minv = tree[l]
						}
						l++
					}
					if r&1 == 1 {
						r--
						if tree[r] < minv {
							minv = tree[r]
						}
					}
					l >>= 1
					r >>= 1
				}
				if minv < p {
					ok = false
					continue
				}
			}

			ptr[x]++
			nv := INF
			if ptr[x] < len(pos[x]) {
				nv = pos[x][ptr[x]]
			}
			id := size + x - 1
			tree[id] = nv
			for id >>= 1; id > 0; id >>= 1 {
				if tree[id<<1] < tree[id<<1|1] {
					tree[id] = tree[id<<1]
				} else {
					tree[id] = tree[id<<1|1]
				}
			}
		}

		if ok {
			out.WriteString("YES\n")
		} else {
			out.WriteString("NO\n")
		}

		for _, x := range used {
			pos[x] = pos[x][:0]
			ptr[x] = 0
		}
	}

	return strings.TrimSpace(out.String())
}

func runBin(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func generateCase(rng *rand.Rand) string {
	t := rng.Intn(3) + 1
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", t))
	for i := 0; i < t; i++ {
		n := rng.Intn(8) + 1
		sb.WriteString(fmt.Sprintf("%d\n", n))
		for j := 0; j < n; j++ {
			if j > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprintf("%d", rng.Intn(n)+1))
		}
		sb.WriteByte('\n')
		for j := 0; j < n; j++ {
			if j > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprintf("%d", rng.Intn(n)+1))
		}
		sb.WriteByte('\n')
	}
	return sb.String()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	cand := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		input := generateCase(rng)
		exp := solveD(input)
		out, err := runBin(cand, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:%s", i+1, err, input)
			os.Exit(1)
		}
		if out != exp {
			fmt.Fprintf(os.Stderr, "case %d failed: expected\n%s\ngot\n%s\ninput:%s", i+1, exp, out, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
