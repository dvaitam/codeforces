package main

import (
	"bytes"
	"fmt"
	"math/bits"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

// Embedded solver from 1569F.go
func solveF(input string) string {
	reader := strings.NewReader(input)
	var n, m, k int
	fmt.Fscan(reader, &n, &m, &k)

	var adj [12][12]bool
	for i := 0; i < m; i++ {
		var u, v int
		fmt.Fscan(reader, &u, &v)
		u--
		v--
		adj[u][v] = true
		adj[v][u] = true
	}

	K := n / 2
	type Path struct {
		seq [6]int
	}
	var paths [4096][12][]Path

	var genPaths func(u int, mask int, depth int, currentSeq *[6]int)
	genPaths = func(u int, mask int, depth int, currentSeq *[6]int) {
		if depth == K {
			paths[mask][u] = append(paths[mask][u], Path{*currentSeq})
			return
		}
		for v := 0; v < n; v++ {
			if (mask&(1<<v)) == 0 && adj[u][v] {
				currentSeq[depth] = v
				genPaths(v, mask|(1<<v), depth+1, currentSeq)
			}
		}
	}

	for i := 0; i < n; i++ {
		var seq [6]int
		seq[0] = i
		genPaths(i, 1<<i, 1, &seq)
	}

	goodMatchings := make(map[uint64]bool)
	fullMask := (1 << n) - 1

	for mask1 := 0; mask1 < (1 << (n - 1)); mask1++ {
		if bits.OnesCount(uint(mask1)) != K {
			continue
		}
		mask2 := fullMask ^ mask1

		for u1 := 0; u1 < n; u1++ {
			if len(paths[mask1][u1]) == 0 {
				continue
			}
			for u2 := 0; u2 < n; u2++ {
				if !adj[u1][u2] || len(paths[mask2][u2]) == 0 {
					continue
				}
				for _, p1 := range paths[mask1][u1] {
					for _, p2 := range paths[mask2][u2] {
						var edges [6]uint64
						for i := 0; i < K; i++ {
							u, v := p1.seq[i], p2.seq[i]
							if u > v {
								u, v = v, u
							}
							edges[i] = uint64((u << 4) | v)
						}
						for i := 1; i < K; i++ {
							for j := i; j > 0 && edges[j-1] > edges[j]; j-- {
								edges[j-1], edges[j] = edges[j], edges[j-1]
							}
						}
						var id uint64
						for i := 0; i < K; i++ {
							id = (id << 8) | edges[i]
						}
						goodMatchings[id] = true
					}
				}
			}
		}
	}

	goodPartitions := make(map[uint64]bool)
	for id := range goodMatchings {
		var mate [12]int
		for i := 0; i < K; i++ {
			edge := (id >> (8 * (K - 1 - i))) & 0xFF
			u, v := int(edge>>4), int(edge&0xF)
			mate[u] = v
			mate[v] = u
		}
		var rgs uint64
		var nextID uint64 = 0
		var blockID [12]int
		for i := 0; i < 12; i++ {
			blockID[i] = -1
		}
		for i := 0; i < n; i++ {
			if blockID[i] == -1 {
				blockID[i] = int(nextID)
				blockID[mate[i]] = int(nextID)
				nextID++
			}
			rgs |= uint64(blockID[i]) << (4 * i)
		}
		goodPartitions[rgs] = true
	}

	queue := make([]uint64, 0, len(goodPartitions))
	for rgs := range goodPartitions {
		queue = append(queue, rgs)
	}

	for head := 0; head < len(queue); head++ {
		curr := queue[head]
		var maxBlock uint64 = 0
		var blocks [12]uint64
		for i := 0; i < n; i++ {
			b := (curr >> (4 * i)) & 0xF
			blocks[i] = b
			if b > maxBlock {
				maxBlock = b
			}
		}

		numBlocks := int(maxBlock + 1)
		for b1 := 0; b1 < numBlocks; b1++ {
			for b2 := b1 + 1; b2 < numBlocks; b2++ {
				var nextRGS uint64
				for i := 0; i < n; i++ {
					b := blocks[i]
					if b == uint64(b2) {
						b = uint64(b1)
					} else if b > uint64(b2) {
						b--
					}
					nextRGS |= b << (4 * i)
				}
				if !goodPartitions[nextRGS] {
					goodPartitions[nextRGS] = true
					queue = append(queue, nextRGS)
				}
			}
		}
	}

	var ans int64 = 0
	for rgs := range goodPartitions {
		var maxBlock uint64 = 0
		for i := 0; i < n; i++ {
			b := (rgs >> (4 * i)) & 0xF
			if b > maxBlock {
				maxBlock = b
			}
		}
		c := int(maxBlock + 1)

		if k >= c {
			ways := int64(1)
			for i := 0; i < c; i++ {
				ways *= int64(k - i)
			}
			ans += ways
		}
	}

	return fmt.Sprintf("%d", ans)
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

func genCase(r *rand.Rand) string {
	n := r.Intn(4)*2 + 2 // even between 2 and 8
	maxEdges := n * (n - 1) / 2
	m := r.Intn(maxEdges + 1)
	k := r.Intn(5) + 1
	type edge struct{ u, v int }
	edges := make(map[edge]struct{})
	for len(edges) < m {
		u := r.Intn(n) + 1
		v := r.Intn(n) + 1
		if u == v {
			continue
		}
		if u > v {
			u, v = v, u
		}
		e := edge{u, v}
		if _, ok := edges[e]; !ok {
			edges[e] = struct{}{}
		}
	}
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d %d\n", n, len(edges), k)
	for e := range edges {
		fmt.Fprintf(&sb, "%d %d\n", e.u, e.v)
	}
	return sb.String()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	userBin := os.Args[1]
	r := rand.New(rand.NewSource(1))
	for i := 0; i < 100; i++ {
		input := genCase(r)
		want := solveF(input)
		got, err := run(userBin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "test %d: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
		if want != got {
			fmt.Printf("test %d failed\ninput:\n%sexpected: %s\ngot: %s\n", i+1, input, want, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", 100)
}
