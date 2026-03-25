package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

// Embedded CF-accepted solver for 571D

type solveBIT struct {
	tree []int64
	n    int
}

func newSolveBIT(n int) *solveBIT {
	return &solveBIT{tree: make([]int64, n+2), n: n}
}

func (b *solveBIT) add(i int, v int64) {
	for i <= b.n {
		b.tree[i] += v
		i += i & -i
	}
}

func (b *solveBIT) point(i int) int64 {
	var res int64
	for i > 0 {
		res += b.tree[i]
		i -= i & -i
	}
	return res
}

func solveBuildIntervals(n, tot int, left, right []int32, roots []int32) ([]int32, []int32, []int32) {
	l := make([]int32, tot+1)
	r := make([]int32, tot+1)
	pos := make([]int32, n+1)
	stackU := make([]int32, 0, tot*2)
	stackS := make([]byte, 0, tot*2)
	var cnt int32
	n32 := int32(n)

	for _, root := range roots {
		stackU = append(stackU, root)
		stackS = append(stackS, 0)
		for len(stackU) > 0 {
			u := stackU[len(stackU)-1]
			st := stackS[len(stackS)-1]
			stackU = stackU[:len(stackU)-1]
			stackS = stackS[:len(stackS)-1]

			if st == 0 {
				stackU = append(stackU, u)
				stackS = append(stackS, 1)
				if u > n32 {
					ru := right[u]
					lu := left[u]
					if ru != 0 {
						stackU = append(stackU, ru)
						stackS = append(stackS, 0)
					}
					if lu != 0 {
						stackU = append(stackU, lu)
						stackS = append(stackS, 0)
					}
				}
			} else {
				if u <= n32 {
					cnt++
					l[u] = cnt
					r[u] = cnt
					pos[u] = cnt
				} else {
					l[u] = l[left[u]]
					r[u] = r[right[u]]
				}
			}
		}
	}

	return l, r, pos
}

func solveUpdateMax(seg []int32, size int, l, r, val int32) {
	li := int(l) - 1 + size
	ri := int(r) - 1 + size
	for li <= ri {
		if li&1 == 1 {
			if seg[li] < val {
				seg[li] = val
			}
			li++
		}
		if ri&1 == 0 {
			if seg[ri] < val {
				seg[ri] = val
			}
			ri--
		}
		li >>= 1
		ri >>= 1
	}
}

func solveQueryMax(seg []int32, size int, p int32) int32 {
	i := int(p) - 1 + size
	var res int32
	for i > 0 {
		if res < seg[i] {
			res = seg[i]
		}
		i >>= 1
	}
	return res
}

func solveD(input string) string {
	tokens := strings.Fields(input)
	pos := 0
	nextInt := func() int {
		v, _ := strconv.Atoi(tokens[pos])
		pos++
		return v
	}
	nextByte := func() byte {
		b := tokens[pos][0]
		pos++
		return b
	}

	n := nextInt()
	m := nextInt()

	maxNodes := 2*n + 5

	leftUni := make([]int32, maxNodes)
	rightUni := make([]int32, maxNodes)
	leftOff := make([]int32, maxNodes)
	rightOff := make([]int32, maxNodes)
	sizeUni := make([]int32, maxNodes)

	currentUni := make([]int32, n+1)
	currentOff := make([]int32, n+1)

	for i := 1; i <= n; i++ {
		currentUni[i] = int32(i)
		currentOff[i] = int32(i)
		sizeUni[i] = 1
	}

	uniAddNode := make([]int32, m+1)
	uniAddVal := make([]int64, m+1)
	raidNode := make([]int32, m+1)
	queryIdAtTime := make([]int32, m+1)
	queryDorm := make([]int32, m+1)

	uniTot := n
	offTot := n
	qcnt := 0

	for t := 1; t <= m; t++ {
		op := nextByte()
		switch op {
		case 'U':
			a := nextInt()
			b := nextInt()
			uniTot++
			na := currentUni[a]
			nb := currentUni[b]
			leftUni[uniTot] = na
			rightUni[uniTot] = nb
			sizeUni[uniTot] = sizeUni[na] + sizeUni[nb]
			currentUni[a] = int32(uniTot)
			currentUni[b] = 0
		case 'M':
			c := nextInt()
			d := nextInt()
			offTot++
			nc := currentOff[c]
			nd := currentOff[d]
			leftOff[offTot] = nc
			rightOff[offTot] = nd
			currentOff[c] = int32(offTot)
			currentOff[d] = 0
		case 'A':
			x := nextInt()
			node := currentUni[x]
			uniAddNode[t] = node
			uniAddVal[t] = int64(sizeUni[node])
		case 'Z':
			y := nextInt()
			raidNode[t] = currentOff[y]
		case 'Q':
			q := nextInt()
			qcnt++
			queryIdAtTime[t] = int32(qcnt)
			queryDorm[qcnt] = int32(q)
		}
	}

	rootsOff := make([]int32, 0, n)
	for i := 1; i <= n; i++ {
		if currentOff[i] != 0 {
			rootsOff = append(rootsOff, currentOff[i])
		}
	}

	lOff, rOff, posOff := solveBuildIntervals(n, offTot, leftOff, rightOff, rootsOff)

	sizeSeg := 1
	for sizeSeg < n {
		sizeSeg <<= 1
	}
	seg := make([]int32, sizeSeg<<1)

	lastRaid := make([]int32, qcnt+1)

	for t := 1; t <= m; t++ {
		if node := raidNode[t]; node != 0 {
			solveUpdateMax(seg, sizeSeg, lOff[node], rOff[node], int32(t))
		}
		if qid := queryIdAtTime[t]; qid != 0 {
			lastRaid[qid] = solveQueryMax(seg, sizeSeg, posOff[queryDorm[qid]])
		}
	}

	rootsUni := make([]int32, 0, n)
	for i := 1; i <= n; i++ {
		if currentUni[i] != 0 {
			rootsUni = append(rootsUni, currentUni[i])
		}
	}

	lUni, rUni, posUni := solveBuildIntervals(n, uniTot, leftUni, rightUni, rootsUni)

	headReq := make([]int32, m+1)
	reqCap := 2*qcnt + 5
	reqNext := make([]int32, reqCap)
	reqQid := make([]int32, reqCap)
	reqCoef := make([]int8, reqCap)
	reqCnt := 0

	for t := 1; t <= m; t++ {
		if qid := queryIdAtTime[t]; qid != 0 {
			reqCnt++
			reqNext[reqCnt] = headReq[t]
			headReq[t] = int32(reqCnt)
			reqQid[reqCnt] = qid
			reqCoef[reqCnt] = 1

			if r := lastRaid[qid]; r != 0 {
				reqCnt++
				reqNext[reqCnt] = headReq[r]
				headReq[r] = int32(reqCnt)
				reqQid[reqCnt] = qid
				reqCoef[reqCnt] = -1
			}
		}
	}

	bit := newSolveBIT(n)
	ans := make([]int64, qcnt+1)

	for t := 1; t <= m; t++ {
		if node := uniAddNode[t]; node != 0 {
			l := int(lUni[node])
			r := int(rUni[node])
			v := uniAddVal[t]
			bit.add(l, v)
			if r < n {
				bit.add(r+1, -v)
			}
		}
		for e := headReq[t]; e != 0; e = reqNext[e] {
			qid := reqQid[e]
			val := bit.point(int(posUni[queryDorm[qid]]))
			ans[qid] += int64(reqCoef[e]) * val
		}
	}

	var out strings.Builder
	for i := 1; i <= qcnt; i++ {
		out.WriteString(strconv.FormatInt(ans[i], 10))
		out.WriteByte('\n')
	}
	return strings.TrimSpace(out.String())
}

// genCase generates a valid random test case for 571D
func genCase(rng *rand.Rand) string {
	n := rng.Intn(4) + 1
	m := rng.Intn(10) + 1

	// Track which university/dorm leaders are still active
	uniLeader := make([]bool, n+1)
	dormLeader := make([]bool, n+1)
	for i := 1; i <= n; i++ {
		uniLeader[i] = true
		dormLeader[i] = true
	}

	var sb strings.Builder
	ops := make([]string, 0, m)
	hasQ := false

	for i := 0; i < m; i++ {
		opType := rng.Intn(5)
		switch opType {
		case 0: // U
			// Need two different active uni leaders
			active := []int{}
			for j := 1; j <= n; j++ {
				if uniLeader[j] {
					active = append(active, j)
				}
			}
			if len(active) < 2 {
				// Fall through to Q
				q := rng.Intn(n) + 1
				ops = append(ops, fmt.Sprintf("Q %d", q))
				hasQ = true
				continue
			}
			rng.Shuffle(len(active), func(a, b int) { active[a], active[b] = active[b], active[a] })
			a, b := active[0], active[1]
			uniLeader[b] = false
			ops = append(ops, fmt.Sprintf("U %d %d", a, b))
		case 1: // M
			active := []int{}
			for j := 1; j <= n; j++ {
				if dormLeader[j] {
					active = append(active, j)
				}
			}
			if len(active) < 2 {
				q := rng.Intn(n) + 1
				ops = append(ops, fmt.Sprintf("Q %d", q))
				hasQ = true
				continue
			}
			rng.Shuffle(len(active), func(a, b int) { active[a], active[b] = active[b], active[a] })
			c, d := active[0], active[1]
			dormLeader[d] = false
			ops = append(ops, fmt.Sprintf("M %d %d", c, d))
		case 2: // A
			active := []int{}
			for j := 1; j <= n; j++ {
				if uniLeader[j] {
					active = append(active, j)
				}
			}
			if len(active) == 0 {
				q := rng.Intn(n) + 1
				ops = append(ops, fmt.Sprintf("Q %d", q))
				hasQ = true
				continue
			}
			x := active[rng.Intn(len(active))]
			ops = append(ops, fmt.Sprintf("A %d", x))
		case 3: // Z
			active := []int{}
			for j := 1; j <= n; j++ {
				if dormLeader[j] {
					active = append(active, j)
				}
			}
			if len(active) == 0 {
				q := rng.Intn(n) + 1
				ops = append(ops, fmt.Sprintf("Q %d", q))
				hasQ = true
				continue
			}
			y := active[rng.Intn(len(active))]
			ops = append(ops, fmt.Sprintf("Z %d", y))
		case 4: // Q
			q := rng.Intn(n) + 1
			ops = append(ops, fmt.Sprintf("Q %d", q))
			hasQ = true
		}
	}

	if !hasQ {
		q := rng.Intn(n) + 1
		ops = append(ops, fmt.Sprintf("Q %d", q))
	}

	fmt.Fprintf(&sb, "%d %d\n", n, len(ops))
	for _, op := range ops {
		sb.WriteString(op)
		sb.WriteByte('\n')
	}
	return sb.String()
}

func runCase(bin, input string) error {
	expected := solveD(input)

	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errb bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errb
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\nstderr: %s", err, errb.String())
	}
	got := strings.TrimSpace(out.String())
	if got != expected {
		return fmt.Errorf("expected:\n%s\n\ngot:\n%s", expected, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	for i := 0; i < 100; i++ {
		input := genCase(rng)
		if err := runCase(bin, input); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
