package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strings"
	"time"
)

// ---------- embedded reference solver for 1740G ----------

type gPortal struct {
	r, c int32
	s    int64
}

func refSolve(input string) string {
	data := []byte(input)
	offset := 0
	nextInt := func() int {
		for offset < len(data) && (data[offset] < '0' || data[offset] > '9') {
			offset++
		}
		if offset >= len(data) {
			return 0
		}
		res := 0
		for offset < len(data) && data[offset] >= '0' && data[offset] <= '9' {
			res = res*10 + int(data[offset]-'0')
			offset++
		}
		return res
	}

	N := nextInt()
	if N == 0 {
		return ""
	}
	M := nextInt()

	totalFaces := int32(4 * N * M)
	parent := make([]int32, totalFaces)
	ends := make([][2]int32, totalFaces)
	endsLen := make([]int8, totalFaces)
	inter := make([]int32, totalFaces)
	rateArr := make([]int8, totalFaces)
	lastTime := make([]int64, totalFaces)
	eArr := make([]int8, totalFaces)

	var gFind func(i int32) int32
	gFind = func(i int32) int32 {
		root := i
		for root != parent[root] {
			root = parent[root]
		}
		curr := i
		for curr != root {
			nxt := parent[curr]
			parent[curr] = root
			curr = nxt
		}
		return root
	}

	mergeInternal := func(fA, fB int32, S int64) {
		rA := gFind(fA)
		rB := gFind(fB)

		if rA == rB {
			for i := int8(0); i < endsLen[rA]; i++ {
				u := ends[rA][i]
				eArr[u] = int8((int64(eArr[u]) + int64(rateArr[u])*(S-lastTime[u])) % 2)
				lastTime[u] = S
			}
			endsLen[rA] = 0
			return
		}

		for i := int8(0); i < endsLen[rA]; i++ {
			u := ends[rA][i]
			eArr[u] = int8((int64(eArr[u]) + int64(rateArr[u])*(S-lastTime[u])) % 2)
			lastTime[u] = S
		}
		for i := int8(0); i < endsLen[rB]; i++ {
			u := ends[rB][i]
			eArr[u] = int8((int64(eArr[u]) + int64(rateArr[u])*(S-lastTime[u])) % 2)
			lastTime[u] = S
		}

		parent[rB] = rA
		inter[rA] += inter[rB] + 1

		newLen := int8(0)
		for i := int8(0); i < endsLen[rA]; i++ {
			u := ends[rA][i]
			if u != fA && u != fB {
				ends[rA][newLen] = u
				newLen++
			}
		}
		for i := int8(0); i < endsLen[rB]; i++ {
			u := ends[rB][i]
			if u != fA && u != fB {
				ends[rA][newLen] = u
				newLen++
			}
		}
		endsLen[rA] = newLen

		newRate := int8(inter[rA] % 2)
		for i := int8(0); i < endsLen[rA]; i++ {
			u := ends[rA][i]
			rateArr[u] = newRate
			lastTime[u] = S
		}
	}

	portals := make([]gPortal, 0, N*M)

	for r := int32(0); r < int32(N); r++ {
		for c := int32(0); c < int32(M); c++ {
			s := int64(nextInt())
			portals = append(portals, gPortal{r, c, s})
			for k := int32(0); k < 4; k++ {
				f := 4*(r*int32(M)+c) + k
				parent[f] = f
				ends[f][0] = f
				endsLen[f] = 1
				lastTime[f] = 1
			}
		}
	}

	for r := int32(0); r < int32(N); r++ {
		for c := int32(0); c < int32(M); c++ {
			f := 4 * (r*int32(M) + c)
			if r > 0 {
				f2 := 4*((r-1)*int32(M)+c) + 2
				rA := f
				rB := f2
				parent[rB] = rA
				ends[rA][0] = f
				ends[rA][1] = f2
				endsLen[rA] = 2
			}
			if c > 0 {
				f3 := 4*(r*int32(M)+c) + 3
				f4 := 4*(r*int32(M)+c-1) + 1
				rA := f3
				rB := f4
				parent[rB] = rA
				ends[rA][0] = f3
				ends[rA][1] = f4
				endsLen[rA] = 2
			}
		}
	}

	sort.Slice(portals, func(i, j int) bool {
		return portals[i].s < portals[j].s
	})

	ans := make([][]byte, N)
	for i := 0; i < N; i++ {
		ans[i] = make([]byte, M)
	}

	i := 0
	for i < len(portals) {
		j := i
		for j < len(portals) && portals[j].s == portals[i].s {
			j++
		}

		group := portals[i:j]
		S := group[0].s

		for _, p := range group {
			for k := int32(0); k < 4; k++ {
				f := 4*(p.r*int32(M)+p.c) + k
				eArr[f] = int8((int64(eArr[f]) + int64(rateArr[f])*(S-lastTime[f])) % 2)
				lastTime[f] = S
			}
		}

		for _, p := range group {
			f0 := 4*(p.r*int32(M)+p.c) + 0
			f1 := 4*(p.r*int32(M)+p.c) + 1
			f2 := 4*(p.r*int32(M)+p.c) + 2
			f3 := 4*(p.r*int32(M)+p.c) + 3

			sum := eArr[f0] + eArr[f1] + eArr[f2] + eArr[f3]
			t := sum % 2
			ans[p.r][p.c] = byte('0' + t)

			if t == 0 {
				mergeInternal(f0, f2, S)
				mergeInternal(f1, f3, S)
			} else {
				mergeInternal(f0, f3, S)
				mergeInternal(f2, f1, S)
			}
		}

		i = j
	}

	var out bytes.Buffer
	for r := 0; r < N; r++ {
		out.Write(ans[r])
		out.WriteByte('\n')
	}
	return strings.TrimSpace(out.String())
}

// ---------- verifier harness ----------

func run(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out, errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func genCase(rng *rand.Rand) string {
	n := rng.Intn(3) + 1
	m := rng.Intn(3) + 1
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", n, m))
	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			if j > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprintf("%d", rng.Intn(5)))
		}
		sb.WriteByte('\n')
	}
	return sb.String()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierG.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 1; i <= 100; i++ {
		input := genCase(rng)
		expect := refSolve(input)
		got, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i, err, input)
			os.Exit(1)
		}
		if got != expect {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i, expect, got, input)
			os.Exit(1)
		}
	}
	fmt.Println("All 100 tests passed")
}
