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

const refSource = `package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

type Pt struct {
	x int64
	y int64
}

func cross(a, b, c Pt) int64 {
	return (b.x-a.x)*(c.y-a.y) - (b.y-a.y)*(c.x-a.x)
}

func lessAround(pts []Pt, c, p, q int) bool {
	dx1 := pts[p].x - pts[c].x
	dy1 := pts[p].y - pts[c].y
	dx2 := pts[q].x - pts[c].x
	dy2 := pts[q].y - pts[c].y
	u1 := dy1 > 0 || (dy1 == 0 && dx1 > 0)
	u2 := dy2 > 0 || (dy2 == 0 && dx2 > 0)
	if u1 != u2 {
		return u1
	}
	return dx1*dy2-dy1*dx2 > 0
}

func main() {
	in := bufio.NewReaderSize(os.Stdin, 1<<20)
	out := bufio.NewWriterSize(os.Stdout, 1<<20)
	defer out.Flush()

	var N, K int
	fmt.Fscan(in, &N, &K)
	pts := make([]Pt, N)
	for i := 0; i < N; i++ {
		fmt.Fscan(in, &pts[i].x, &pts[i].y)
	}

	W := (N + 63) >> 6
	leftBits := make([]uint64, N*N*W)
	for i := 0; i < N; i++ {
		pi := pts[i]
		for j := 0; j < N; j++ {
			if i == j {
				continue
			}
			pj := pts[j]
			base := (i*N + j) * W
			for k := 0; k < N; k++ {
				pk := pts[k]
				if (pj.x-pi.x)*(pk.y-pi.y)-(pj.y-pi.y)*(pk.x-pi.x) > 0 {
					leftBits[base+(k>>6)] |= uint64(1) << uint(k&63)
				}
			}
		}
	}

	L := N - 1
	order := make([]int, N*L)
	pos := make([]int, N*N)
	for i := range pos {
		pos[i] = -1
	}
	endPos := make([]int, N*L)
	tmp := make([]int, L)
	for c := 0; c < N; c++ {
		idx := 0
		for p := 0; p < N; p++ {
			if p != c {
				tmp[idx] = p
				idx++
			}
		}
		ids := tmp[:L]
		sort.Slice(ids, func(i, j int) bool {
			return lessAround(pts, c, ids[i], ids[j])
		})
		base := c * L
		copy(order[base:base+L], ids)
		for i, p := range ids {
			pos[c*N+p] = i
		}
		r := 0
		for l := 0; l < L; l++ {
			if r < l {
				r = l
			}
			bp := ids[l]
			for r+1 < l+L && cross(pts[c], pts[bp], pts[ids[(r+1)%L]]) > 0 {
				r++
			}
			endPos[base+l] = r
		}
	}

	ans2 := int64(0)
	vals := make([]int64, 2*L)
	bestPos := make([]int64, L)
	deque := make([]int, 2*L)
	idsBuf := make([]int, L)

	for a := 0; a < N; a++ {
		pa := pts[a]
		baseOrd := a * L
		m := 0
		for t := 0; t < L; t++ {
			p := order[baseOrd+t]
			pp := pts[p]
			if pp.y > pa.y || (pp.y == pa.y && pp.x > pa.x) {
				idsBuf[m] = p
				m++
			} else {
				break
			}
		}
		if m < K-1 {
			continue
		}
		ids := idsBuf[:m]
		mm := m * m
		areaLocal := make([]int64, mm)
		triEmpty := make([]bool, mm)
		for i := 0; i < m; i++ {
			b := ids[i]
			baseAB := (a*N + b) * W
			for j := i + 1; j < m; j++ {
				c := ids[j]
				idx := i*m + j
				areaLocal[idx] = cross(pa, pts[b], pts[c])
				baseBC := (b*N + c) * W
				baseCA := (c*N + a) * W
				ok := true
				for w := 0; w < W; w++ {
					if leftBits[baseAB+w]&leftBits[baseBC+w]&leftBits[baseCA+w] != 0 {
						ok = false
						break
					}
				}
				triEmpty[idx] = ok
			}
		}

		dpPrev := make([]int64, mm)
		dpCurr := make([]int64, mm)
		for i := 0; i < mm; i++ {
			dpPrev[i] = -1
			dpCurr[i] = -1
		}

		maxLast := m - 1
		if K > 3 {
			maxLast = m - 1 - (K - 3)
		}
		for i := 0; i < maxLast; i++ {
			row := i * m
			for j := i + 1; j <= maxLast; j++ {
				idx := row + j
				if triEmpty[idx] {
					v := areaLocal[idx]
					dpPrev[idx] = v
					if K == 3 && v > ans2 {
						ans2 = v
					}
				}
			}
		}
		if K == 3 {
			continue
		}

		for length := 4; length <= K; length++ {
			for i := 0; i < mm; i++ {
				dpCurr[i] = -1
			}
			maxLast = m - 1
			if length < K {
				maxLast = m - 1 - (K - length)
			}
			iStart := length - 3
			for i := iStart; i < maxLast; i++ {
				mid := ids[i]
				for t := 0; t < L; t++ {
					vals[t] = -1
				}
				midPosBase := mid * N
				for h := 0; h < i; h++ {
					v := dpPrev[h*m+i]
					if v >= 0 {
						vals[pos[midPosBase+ids[h]]] = v
					}
				}
				for t := 0; t < L; t++ {
					vals[t+L] = vals[t]
				}

				head, tail := 0, 0
				r := 0
				epBase := mid * L
				for p := 0; p < L; p++ {
					if r < p {
						r = p
					}
					target := endPos[epBase+p]
					for r < target {
						r++
						v := vals[r]
						for head < tail && vals[deque[tail-1]] <= v {
							tail--
						}
						deque[tail] = r
						tail++
					}
					left := p + 1
					for head < tail && deque[head] < left {
						head++
					}
					if head < tail {
						bestPos[p] = vals[deque[head]]
					} else {
						bestPos[p] = -1
					}
				}

				row := i * m
				for j := i + 1; j <= maxLast; j++ {
					idx := row + j
					if !triEmpty[idx] {
						continue
					}
					best := bestPos[pos[midPosBase+ids[j]]]
					if best >= 0 {
						v := best + areaLocal[idx]
						dpCurr[idx] = v
						if length == K && v > ans2 {
							ans2 = v
						}
					}
				}
			}
			dpPrev, dpCurr = dpCurr, dpPrev
		}
	}

	fmt.Fprintf(out, "%.2f\n", float64(ans2)/2.0)
}
`

func buildRef() (string, func(), error) {
	tmpDir, err := os.MkdirTemp("", "refbuild")
	if err != nil {
		return "", nil, err
	}
	srcPath := filepath.Join(tmpDir, "ref.go")
	if err := os.WriteFile(srcPath, []byte(refSource), 0644); err != nil {
		os.RemoveAll(tmpDir)
		return "", nil, err
	}
	binPath := filepath.Join(tmpDir, "ref")
	cmd := exec.Command("go", "build", "-o", binPath, srcPath)
	if out, err := cmd.CombinedOutput(); err != nil {
		os.RemoveAll(tmpDir)
		return "", nil, fmt.Errorf("failed to build reference: %v\n%s", err, out)
	}
	return binPath, func() { os.RemoveAll(tmpDir) }, nil
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
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func deterministicCases() []string {
	return []string{"3 3\n0 0\n1 0\n0 1\n"}
}

func randomCase(rng *rand.Rand) string {
	n := rng.Intn(4) + 3
	K := rng.Intn(n-2) + 3
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", n, K))
	for i := 0; i < n; i++ {
		x := rng.Intn(5)
		y := rng.Intn(5)
		sb.WriteString(fmt.Sprintf("%d %d\n", x, y))
	}
	return sb.String()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierH.go /path/to/binary")
		os.Exit(1)
	}
	candidate := os.Args[1]
	ref, cleanup, err := buildRef()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer cleanup()
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	cases := deterministicCases()
	for len(cases) < 100 {
		cases = append(cases, randomCase(rng))
	}
	for i, in := range cases {
		exp, err := run(ref, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference failed on case %d: %v\n", i+1, err)
			os.Exit(1)
		}
		out, err := run(candidate, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
		if out != exp {
			fmt.Fprintf(os.Stderr, "case %d wrong answer\nexpected:\n%s\ngot:\n%s\ninput:\n%s", i+1, exp, out, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
