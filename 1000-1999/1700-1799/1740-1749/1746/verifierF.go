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

// Embedded correct solver for 1746F.
// Uses hashing: assign random weights to each distinct value,
// maintain prefix sums in a BIT; for query (l,r,k) check that
// the sum of weights in [l..r] is divisible by k for every hash.
const W = 40

func solve(n, q int, a []int, queries [][]int) []string {
	rng := rand.New(rand.NewSource(12345))

	idMap := make(map[int]int)
	var vals [][W]uint32

	getID := func(x int) int {
		if id, ok := idMap[x]; ok {
			return id
		}
		id := len(vals)
		idMap[x] = id
		var v [W]uint32
		for i := 0; i < W; i++ {
			v[i] = rng.Uint32()
		}
		vals = append(vals, v)
		return id
	}

	ids := make([]int, n+1)
	tree := make([][W]uint64, n+1)

	for i := 1; i <= n; i++ {
		id := getID(a[i])
		ids[i] = id
		for w := 0; w < W; w++ {
			tree[i][w] += uint64(vals[id][w])
		}
	}

	for i := 1; i <= n; i++ {
		p := i + (i & -i)
		if p <= n {
			for w := 0; w < W; w++ {
				tree[p][w] += tree[i][w]
			}
		}
	}

	var results []string

	for _, qr := range queries {
		if qr[0] == 1 {
			idx, x := qr[1], qr[2]
			oldID := ids[idx]
			newID := getID(x)
			if oldID != newID {
				ids[idx] = newID
				for j := idx; j <= n; j += j & -j {
					for w := 0; w < W; w++ {
						tree[j][w] += uint64(vals[newID][w])
						tree[j][w] -= uint64(vals[oldID][w])
					}
				}
			}
		} else {
			l, r, k := qr[1], qr[2], qr[3]
			if (r-l+1)%k != 0 {
				results = append(results, "NO")
			} else if k == 1 {
				results = append(results, "YES")
			} else {
				var sum [W]uint64
				for j := r; j > 0; j -= j & -j {
					for w := 0; w < W; w++ {
						sum[w] += tree[j][w]
					}
				}
				for j := l - 1; j > 0; j -= j & -j {
					for w := 0; w < W; w++ {
						sum[w] -= tree[j][w]
					}
				}
				ok := true
				uk := uint64(k)
				for w := 0; w < W; w++ {
					if sum[w]%uk != 0 {
						ok = false
						break
					}
				}
				if ok {
					results = append(results, "YES")
				} else {
					results = append(results, "NO")
				}
			}
		}
	}
	return results
}

func genTest(rng *rand.Rand) (string, int, int, []int, [][]int) {
	n := rng.Intn(5) + 1
	q := rng.Intn(5) + 1
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", n, q))
	a := make([]int, n+1)
	for i := 1; i <= n; i++ {
		if i > 1 {
			sb.WriteByte(' ')
		}
		a[i] = rng.Intn(5) + 1
		sb.WriteString(fmt.Sprintf("%d", a[i]))
	}
	sb.WriteString("\n")
	var queries [][]int
	for i := 0; i < q; i++ {
		typ := rng.Intn(2) + 1
		if typ == 1 {
			idx := rng.Intn(n) + 1
			x := rng.Intn(5) + 1
			sb.WriteString(fmt.Sprintf("1 %d %d\n", idx, x))
			queries = append(queries, []int{1, idx, x})
		} else {
			l := rng.Intn(n) + 1
			r := rng.Intn(n-l+1) + l
			k := rng.Intn(n) + 1
			sb.WriteString(fmt.Sprintf("2 %d %d %d\n", l, r, k))
			queries = append(queries, []int{2, l, r, k})
		}
	}
	return sb.String(), n, q, a, queries
}

func runBinary(bin string, input string) (string, error) {
	cmd := exec.Command(bin)
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
		fmt.Println("usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	for i := 1; i <= 200; i++ {
		input, n, _, a, queries := genTest(rng)
		expected := solve(n, len(queries), a, queries)
		got, err := runBinary(bin, input)
		if err != nil {
			fmt.Printf("runtime error on test %d: %v\n", i, err)
			os.Exit(1)
		}
		expStr := strings.Join(expected, "\n")
		if strings.TrimSpace(got) != strings.TrimSpace(expStr) {
			fmt.Printf("wrong answer on test %d\ninput:\n%sexpected:\n%s\ngot:\n%s\n", i, input, expStr, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
