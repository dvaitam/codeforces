package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

const MOD int64 = 1000000007

func insertBasis(b *[5]uint8, x uint8) bool {
	for i := 4; i >= 0; i-- {
		if x&(uint8(1)<<uint(i)) == 0 {
			continue
		}
		if b[i] == 0 {
			b[i] = x
			return true
		}
		x ^= b[i]
	}
	return false
}

func spanAddBitset(u uint32, x int) uint32 {
	if x == 0 || (u&(uint32(1)<<uint(x))) != 0 {
		return u
	}
	w := u
	for v := 0; v < 32; v++ {
		if (u & (uint32(1) << uint(v))) != 0 {
			w |= uint32(1) << uint(v^x)
		}
	}
	return w
}

func basisToBitset(b [5]uint8) uint32 {
	u := uint32(1)
	for i := 0; i < 5; i++ {
		if b[i] != 0 {
			u = spanAddBitset(u, int(b[i]))
		}
	}
	return u
}

func generateStates() ([]uint32, map[uint32]int) {
	states := []uint32{1}
	id := map[uint32]int{1: 0}
	for i := 0; i < len(states); i++ {
		u := states[i]
		for x := 1; x < 32; x++ {
			if (u & (uint32(1) << uint(x))) != 0 {
				continue
			}
			w := spanAddBitset(u, x)
			if _, ok := id[w]; !ok {
				id[w] = len(states)
				states = append(states, w)
			}
		}
	}
	return states, id
}

func addOption(ids *[3]int, ways *[3]int64, cnt *int, id int, w int64) {
	if w == 0 {
		return
	}
	for i := 0; i < *cnt; i++ {
		if ids[i] == id {
			ways[i] += w
			return
		}
	}
	ids[*cnt] = id
	ways[*cnt] = w
	*cnt++
}

func refSolve(input string) string {
	data := []byte(input)
	pos := 0
	readInt := func() int {
		for pos < len(data) && (data[pos] < '0' || data[pos] > '9') {
			pos++
		}
		v := 0
		for pos < len(data) && data[pos] >= '0' && data[pos] <= '9' {
			v = v*10 + int(data[pos]-'0')
			pos++
		}
		return v
	}

	n := readInt()
	m := readInt()

	head := make([]int, n+1)
	for i := range head {
		head[i] = -1
	}

	to := make([]int, 0, 2*m)
	nxt := make([]int, 0, 2*m)
	wt := make([]uint8, 0, 2*m)

	addEdge := func(u, v int, w uint8) {
		to = append(to, v)
		wt = append(wt, w)
		nxt = append(nxt, head[u])
		head[u] = len(to) - 1
	}

	has1 := make([]bool, n+1)
	w1 := make([]uint8, n+1)

	eu := make([]int, 0, m)
	ev := make([]int, 0, m)
	ew := make([]uint8, 0, m)

	for i := 0; i < m; i++ {
		a := readInt()
		b := readInt()
		w := uint8(readInt())
		if a == 1 || b == 1 {
			x := a
			if x == 1 {
				x = b
			}
			has1[x] = true
			w1[x] = w
		} else {
			eu = append(eu, a)
			ev = append(ev, b)
			ew = append(ew, w)
			addEdge(a, b, w)
			addEdge(b, a, w)
		}
	}

	comp := make([]int, n+1)
	dist := make([]uint8, n+1)
	verts := []int{0}
	stack := make([]int, 0, n)

	for s := 2; s <= n; s++ {
		if comp[s] != 0 {
			continue
		}
		cid := len(verts)
		verts = append(verts, 0)
		comp[s] = cid
		dist[s] = 0
		stack = append(stack, s)
		for len(stack) > 0 {
			v := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			verts[cid]++
			for e := head[v]; e != -1; e = nxt[e] {
				u := to[e]
				if comp[u] == 0 {
					comp[u] = cid
					dist[u] = dist[v] ^ wt[e]
					stack = append(stack, u)
				}
			}
		}
	}

	k := len(verts) - 1
	edgesCnt := make([]int, k+1)
	rank := make([]int, k+1)
	bases := make([][5]uint8, k+1)

	for i := 0; i < len(eu); i++ {
		cid := comp[eu[i]]
		edgesCnt[cid]++
		val := dist[eu[i]] ^ dist[ev[i]] ^ ew[i]
		if val != 0 {
			if insertBasis(&bases[cid], val) {
				rank[cid]++
			}
		}
	}

	neighCnt := make([]int, k+1)
	na := make([]int, k+1)
	nb := make([]int, k+1)

	for v := 2; v <= n; v++ {
		if has1[v] {
			cid := comp[v]
			if neighCnt[cid] == 0 {
				na[cid] = v
			} else {
				nb[cid] = v
			}
			neighCnt[cid]++
		}
	}

	states, idMap := generateStates()
	numStates := len(states)

	stateBasis := make([][5]uint8, numStates)
	for i, u := range states {
		var b [5]uint8
		for x := 1; x < 32; x++ {
			if (u & (uint32(1) << uint(x))) != 0 {
				insertBasis(&b, uint8(x))
			}
		}
		stateBasis[i] = b
	}

	addVec := make([][32]int16, numStates)
	for i, u := range states {
		for x := 0; x < 32; x++ {
			w := spanAddBitset(u, x)
			addVec[i][x] = int16(idMap[w])
		}
	}

	comb := make([]int16, numStates*numStates)
	for i := range comb {
		comb[i] = -1
	}
	for i := 0; i < numStates; i++ {
		row := i * numStates
		for j := 0; j < numStates; j++ {
			if (states[i] & states[j]) != 1 {
				continue
			}
			cur := i
			b := stateBasis[j]
			for t := 0; t < 5; t++ {
				if b[t] != 0 {
					cur = int(addVec[cur][b[t]])
				}
			}
			comb[row+j] = int16(cur)
		}
	}

	dp := make([]int64, numStates)
	ndp := make([]int64, numStates)
	dp[0] = 1

	for cid := 1; cid <= k; cid++ {
		for i := 0; i < numStates; i++ {
			ndp[i] = 0
		}

		var optIds [3]int
		var optWays [3]int64
		optCnt := 0

		addOption(&optIds, &optWays, &optCnt, 0, 1)

		d := edgesCnt[cid] - verts[cid] + 1
		if d == rank[cid] {
			sbit := basisToBitset(bases[cid])
			sid := idMap[sbit]
			if neighCnt[cid] == 1 {
				addOption(&optIds, &optWays, &optCnt, sid, 1)
			} else if neighCnt[cid] == 2 {
				addOption(&optIds, &optWays, &optCnt, sid, 2)
				x := int(dist[na[cid]] ^ dist[nb[cid]] ^ w1[na[cid]] ^ w1[nb[cid]])
				if (sbit & (uint32(1) << uint(x))) == 0 {
					tid := int(addVec[sid][x])
					addOption(&optIds, &optWays, &optCnt, tid, 1)
				}
			}
		}

		for s, val := range dp {
			if val == 0 {
				continue
			}
			base := s * numStates
			for i := 0; i < optCnt; i++ {
				r := comb[base+optIds[i]]
				if r < 0 {
					continue
				}
				idx := int(r)
				nv := ndp[idx] + val*optWays[i]
				for nv >= MOD {
					nv -= MOD
				}
				ndp[idx] = nv
			}
		}

		dp, ndp = ndp, dp
	}

	ans := int64(0)
	for _, v := range dp {
		ans += v
		if ans >= MOD {
			ans -= MOD
		}
	}

	return strconv.FormatInt(ans, 10)
}

// ---------- test harness ----------

func generateCase(rng *rand.Rand) string {
	n := rng.Intn(5) + 1
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", n, n-1))
	for i := 2; i <= n; i++ {
		parent := rng.Intn(i-1) + 1
		w := rng.Intn(32)
		sb.WriteString(fmt.Sprintf("%d %d %d\n", parent, i, w))
	}
	return sb.String()
}

func runBin(bin, input string) (string, error) {
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

// suppress unused import warnings
var _ = bufio.NewReader
var _ = io.ReadAll

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in := generateCase(rng)
		exp := refSolve(in)
		got, err := runBin(bin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:%s", i+1, err, in)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != strings.TrimSpace(exp) {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:%s", i+1, exp, got, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
