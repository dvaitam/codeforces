package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strconv"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Split(bufio.ScanWords)

	buf := make([]byte, 0, 10*1024*1024)
	scanner.Buffer(buf, 10*1024*1024)

	read := func() int {
		scanner.Scan()
		num, _ := strconv.Atoi(scanner.Text())
		return num
	}

	if !scanner.Scan() {
		return
	}
	N, _ := strconv.Atoi(scanner.Text())
	Q := read()

	adj := make([][]int, N+1)
	orig_degree := make([]int, N+1)
	for i := 0; i < N-1; i++ {
		u, v := read(), read()
		adj[u] = append(adj[u], v)
		adj[v] = append(adj[v], u)
		orig_degree[u]++
		orig_degree[v]++
	}

	L_total := 0
	for i := 1; i <= N; i++ {
		if orig_degree[i] == 1 {
			L_total++
		}
	}

	tin := make([]int, N+1)
	tout := make([]int, N+1)
	up := make([][20]int, N+1)
	depth := make([]int, N+1)
	l_orig := make([]int, N+1)
	timer := 0

	var dfs1 func(u, p int)
	dfs1 = func(u, p int) {
		timer++
		tin[u] = timer
		up[u][0] = p
		for i := 1; i < 20; i++ {
			up[u][i] = up[up[u][i-1]][i-1]
		}

		if orig_degree[u] == 1 {
			l_orig[u] = 1
		}
		for _, v := range adj[u] {
			if v != p {
				depth[v] = depth[u] + 1
				dfs1(v, u)
				l_orig[u] += l_orig[v]
			}
		}
		tout[u] = timer
	}

	depth[1] = 0
	dfs1(1, 1)

	Dist_even := make([]int, N+1)
	Dist_odd := make([]int, N+1)
	TotalOriginalEvenCount := 0

	var dfs2 func(u, p int)
	dfs2 = func(u, p int) {
		for _, v := range adj[u] {
			if v != p {
				if l_orig[v]%2 == 0 {
					Dist_even[v] = Dist_even[u] + 1
					Dist_odd[v] = Dist_odd[u]
					TotalOriginalEvenCount++
				} else {
					Dist_even[v] = Dist_even[u]
					Dist_odd[v] = Dist_odd[u] + 1
				}
				dfs2(v, u)
			}
		}
	}
	dfs2(1, 1)

	get_lca := func(u, v int) int {
		if depth[u] < depth[v] {
			u, v = v, u
		}
		diff := depth[u] - depth[v]
		for i := 0; i < 20; i++ {
			if (diff & (1 << i)) != 0 {
				u = up[u][i]
			}
		}
		if u == v {
			return u
		}
		for i := 19; i >= 0; i-- {
			if up[u][i] != up[v][i] {
				u = up[u][i]
				v = up[v][i]
			}
		}
		return up[u][0]
	}

	count := make([]int, N+1)
	C := make([]int, 0, N)
	virt_adj := make([][]int, N+1)
	is_F := make([]bool, N+1)

	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	for q := 0; q < Q; q++ {
		D_i := read()
		C = C[:0]
		for j := 0; j < D_i; j++ {
			u := read()
			if count[u] == 0 {
				C = append(C, u)
			}
			count[u]++
		}

		newLeaves := L_total + D_i
		for _, u := range C {
			if orig_degree[u] == 1 {
				newLeaves--
			}
		}

		if newLeaves%2 != 0 {
			fmt.Fprintln(out, "-1")
			for _, u := range C {
				count[u] = 0
			}
			continue
		}

		S_initial := make([]int, 0, len(C))
		for _, u := range C {
			c_u := count[u]
			f_u := 0
			if orig_degree[u] == 1 {
				if c_u%2 == 0 {
					f_u = 1
				}
			} else {
				if c_u%2 != 0 {
					f_u = 1
				}
			}
			if f_u == 1 {
				S_initial = append(S_initial, u)
			}
		}

		S := make([]int, len(S_initial))
		copy(S, S_initial)
		has_root := false
		for _, u := range S {
			if u == 1 {
				has_root = true
				break
			}
		}
		if !has_root {
			S = append(S, 1)
		}

		for _, u := range S_initial {
			is_F[u] = true
		}

		sort.Slice(S, func(i, j int) bool { return tin[S[i]] < tin[S[j]] })
		m := len(S)
		for i := 0; i < m-1; i++ {
			lca := get_lca(S[i], S[i+1])
			S = append(S, lca)
		}

		sort.Slice(S, func(i, j int) bool { return tin[S[i]] < tin[S[j]] })
		uniqueS := []int{S[0]}
		for i := 1; i < len(S); i++ {
			if S[i] != S[i-1] {
				uniqueS = append(uniqueS, S[i])
			}
		}
		S = uniqueS

		for _, u := range S {
			virt_adj[u] = virt_adj[u][:0]
		}

		stack := []int{S[0]}
		for i := 1; i < len(S); i++ {
			v := S[i]
			for len(stack) > 0 {
				top := stack[len(stack)-1]
				if tin[top] <= tin[v] && tout[v] <= tout[top] {
					break
				}
				stack = stack[:len(stack)-1]
			}
			p := stack[len(stack)-1]
			virt_adj[p] = append(virt_adj[p], v)
			stack = append(stack, v)
		}

		flip_count := 0
		var dp func(u int) int
		dp = func(u int) int {
			w := 0
			if is_F[u] {
				w = 1
			}
			for _, v := range virt_adj[u] {
				child_w := dp(v)
				w = (w + child_w) % 2
				if child_w == 1 {
					e_even := Dist_even[v] - Dist_even[u]
					e_odd := Dist_odd[v] - Dist_odd[u]
					flip_count += e_odd - e_even
				}
			}
			return w
		}

		dp(S[0])

		ans := (N - 1 + D_i) + TotalOriginalEvenCount + flip_count
		fmt.Fprintln(out, ans)

		for _, u := range S_initial {
			is_F[u] = false
		}
		for _, u := range C {
			count[u] = 0
		}
	}
}
