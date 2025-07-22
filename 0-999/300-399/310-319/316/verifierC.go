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

func solveCase(n, m int, grid [][]int) int {
	N := n * m
	K := N / 2
	pos := make([][2]int, K)
	cnt := make([]int, K)
	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			x := grid[i][j] - 1
			idx := i*m + j
			if cnt[x] < 2 {
				pos[x][cnt[x]] = idx
			}
			cnt[x]++
		}
	}
	edges := make([][2]int, 0)
	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			u := i*m + j
			if j+1 < m {
				edges = append(edges, [2]int{u, u + 1})
			}
			if i+1 < n {
				edges = append(edges, [2]int{u, u + m})
			}
		}
	}
	E := len(edges)
	costMat := make([][]int, K)
	cand := make([][]int, K)
	for k := 0; k < K; k++ {
		costMat[k] = make([]int, E)
		for e := 0; e < E; e++ {
			u, v := edges[e][0], edges[e][1]
			p0, p1 := pos[k][0], pos[k][1]
			c1 := 0
			if p0 != u {
				c1++
			}
			if p1 != v {
				c1++
			}
			c2 := 0
			if p0 != v {
				c2++
			}
			if p1 != u {
				c2++
			}
			if c2 < c1 {
				costMat[k][e] = c2
			} else {
				costMat[k][e] = c1
			}
		}
		idxs := make([]int, E)
		for e := 0; e < E; e++ {
			idxs[e] = e
		}
		for i := 1; i < E; i++ {
			v := idxs[i]
			j := i
			for j > 0 && costMat[k][idxs[j-1]] > costMat[k][v] {
				idxs[j] = idxs[j-1]
				j--
			}
			idxs[j] = v
		}
		cand[k] = idxs
	}
	order := make([]int, K)
	for i := 0; i < K; i++ {
		order[i] = i
	}
	type od struct{ k, c0, c1 int }
	ods := make([]od, K)
	for _, k := range order {
		c0, c1 := 0, 0
		for _, e := range cand[k] {
			if costMat[k][e] == 0 {
				c0++
			} else if costMat[k][e] == 1 {
				c1++
			} else {
				break
			}
		}
		ods[k] = od{k: k, c0: c0, c1: c1}
	}
	for i := 1; i < K; i++ {
		tmp := ods[i]
		j := i
		for j > 0 {
			o := ods[j-1]
			if tmp.c0 < o.c0 || (tmp.c0 == o.c0 && tmp.c1 < o.c1) {
				ods[j] = o
				j--
			} else {
				break
			}
		}
		ods[j] = tmp
	}
	for i := 0; i < K; i++ {
		order[i] = ods[i].k
	}
	lbSum := make([]int, K+1)
	lb := make([]int, K)
	for i := 0; i < K; i++ {
		k := order[i]
		mc := costMat[k][cand[k][0]]
		lb[i] = mc
	}
	lbSum[K] = 0
	for i := K - 1; i >= 0; i-- {
		lbSum[i] = lbSum[i+1] + lb[i]
	}
	best := 2 * K
	var dfs func(idx int, used uint64, cur int)
	dfs = func(idx int, used uint64, cur int) {
		if cur+lbSum[idx] >= best {
			return
		}
		if idx == K {
			if cur < best {
				best = cur
			}
			return
		}
		k := order[idx]
		for _, e := range cand[k] {
			c := costMat[k][e]
			if cur+c >= best {
				continue
			}
			u, v := edges[e][0], edges[e][1]
			bit := (uint64(1) << u) | (uint64(1) << v)
			if used&bit != 0 {
				continue
			}
			dfs(idx+1, used|bit, cur+c)
		}
	}
	dfs(0, 0, 0)
	return best
}

func genCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(3) + 2
	m := rng.Intn(3) + 2
	K := n * m / 2
	nums := make([]int, 2*K)
	for i := 0; i < K; i++ {
		nums[2*i] = i + 1
		nums[2*i+1] = i + 1
	}
	rng.Shuffle(len(nums), func(i, j int) { nums[i], nums[j] = nums[j], nums[i] })
	grid := make([][]int, n)
	idx := 0
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d\n", n, m)
	for i := 0; i < n; i++ {
		grid[i] = make([]int, m)
		for j := 0; j < m; j++ {
			grid[i][j] = nums[idx]
			idx++
			if j > 0 {
				sb.WriteByte(' ')
			}
			fmt.Fprintf(&sb, "%d", grid[i][j])
		}
		sb.WriteByte('\n')
	}
	ans := solveCase(n, m, grid)
	return sb.String(), fmt.Sprintf("%d\n", ans)
}

func runCase(bin, input, expected string) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	if strings.TrimSpace(out.String()) != strings.TrimSpace(expected) {
		return fmt.Errorf("expected %s got %s", expected, out.String())
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := genCase(rng)
		if err := runCase(bin, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
