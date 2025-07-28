package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

const (
	maxV = 105
	maxE = 10000
	inf  = 0x3f3f3f3f
)

var (
	S, T, e int
	fir     [maxV]int
	to, nxt [maxE]int
	w       [maxE]int
	cost    [maxE]int
	dis     [maxV]int64
	q       [maxE]int
	vis     [maxV]bool
	cur     [maxV]int
	fdem    [maxV]int
	odd     []bool
	ansEdge []int
	sumNode []int
	anss    int64
)

func adde(x, y, z, cst int) {
	e++
	to[e] = y
	nxt[e] = fir[x]
	fir[x] = e
	w[e] = z
	cost[e] = cst
	e++
	to[e] = x
	nxt[e] = fir[y]
	fir[y] = e
	w[e] = 0
	cost[e] = -cst
}

func spfa() bool {
	for i := 1; i <= T; i++ {
		dis[i] = 1 << 60
		vis[i] = false
	}
	head, tail := 0, 0
	q[tail] = T
	dis[T] = 0
	vis[T] = true
	tail++
	for head < tail {
		u := q[head]
		head++
		vis[u] = false
		for i := fir[u]; i > 0; i = nxt[i] {
			rev := i ^ 1
			v := to[i]
			if w[rev] == 0 {
				continue
			}
			nd := dis[u] + int64(cost[rev])
			if dis[v] > nd {
				dis[v] = nd
				if !vis[v] {
					vis[v] = true
					q[tail] = v
					tail++
				}
			}
		}
	}
	return dis[S] < 0
}

func dfs(u, flow int) int {
	if u == T || flow == 0 {
		return flow
	}
	vis[u] = true
	used := flow
	for i := cur[u]; i > 0; i = nxt[i] {
		cur[u] = i
		v := to[i]
		if vis[v] || dis[v]+int64(cost[i]) != dis[u] || w[i] == 0 {
			continue
		}
		can := flow
		if w[i] < can {
			can = w[i]
		}
		if can > used {
			can = used
		}
		f := dfs(v, can)
		if f > 0 {
			w[i] -= f
			w[i^1] += f
			used -= f
			if used == 0 {
				break
			}
		}
	}
	return flow - used
}

func MCMF() int64 {
	var flow int
	var res int64
	for spfa() {
		for i := 1; i <= T; i++ {
			cur[i] = fir[i]
			vis[i] = false
		}
		f := dfs(S, inf)
		flow += f
		res += dis[S] * int64(f)
	}
	return res
}

func solveF(input string) string {
	in := bufio.NewReader(strings.NewReader(input))
	var n, m int
	if _, err := fmt.Fscan(in, &n, &m); err != nil {
		return ""
	}
	S, T = 1, n
	e = 1
	odd = make([]bool, m+1)
	ansEdge = make([]int, m+1)
	sumNode = make([]int, n+2)
	for i := 1; i <= m; i++ {
		var x, y, z, t int
		fmt.Fscan(in, &x, &y, &z, &t)
		adde(x, y, z>>1, t<<1)
		if z&1 == 1 {
			fdem[x]--
			fdem[y]++
			odd[i] = true
			anss += int64(t)
		}
	}
	for i := 2; i < n; i++ {
		if fdem[i]&1 != 0 {
			return "Impossible"
		}
		if fdem[i] > 0 {
			adde(S, i, fdem[i]>>1, -inf)
		} else if fdem[i] < 0 {
			adde(i, T, (-fdem[i])>>1, -inf)
		}
	}
	anss += MCMF()
	for i := 1; i <= m; i++ {
		flowUsed := w[2*i+1]
		ans := flowUsed << 1
		if odd[i] {
			ans |= 1
		}
		ansEdge[i] = ans
		u := to[2*i]
		v := to[2*i+1]
		sumNode[u] += ans
		sumNode[v] -= ans
	}
	for i := 2; i < n; i++ {
		if sumNode[i] != 0 {
			return "Impossible"
		}
	}
	var buf strings.Builder
	buf.WriteString("Possible\n")
	for i := 1; i <= m; i++ {
		fmt.Fprintf(&buf, "%d ", ansEdge[i])
	}
	return strings.TrimSpace(buf.String())
}

func genTestF(rng *rand.Rand) string {
	n := rng.Intn(4) + 2
	m := rng.Intn(6) + 1
	var buf strings.Builder
	fmt.Fprintf(&buf, "%d %d\n", n, m)
	for i := 0; i < m; i++ {
		x := rng.Intn(n-1) + 1
		y := rng.Intn(n-x) + x + 1
		if y > n {
			y = n
		}
		if y == 1 {
			y = n
		}
		c := rng.Intn(5) + 1
		w := rng.Intn(11) - 5
		fmt.Fprintf(&buf, "%d %d %d %d\n", x, y, c, w)
	}
	return buf.String()
}

func run(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
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

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	exe := os.Args[1]
	rng := rand.New(rand.NewSource(1))
	for i := 1; i <= 100; i++ {
		in := genTestF(rng)
		expect := solveF(in)
		got, err := run(exe, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: %v\n", i, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != expect {
			fmt.Printf("case %d failed\ninput:\n%sexpected:\n%s\ngot:\n%s\n", i, in, expect, got)
			os.Exit(1)
		}
	}
	fmt.Println("all tests passed")
}
