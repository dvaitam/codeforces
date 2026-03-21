package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"
)

type Test struct {
	input string
}

// Embedded reference solver for 850D

type solveNode struct {
	deg int
	id  int
}

func comb2(x int) int {
	return x * (x - 1) / 2
}

func hasBit(bs []uint64, pos int) bool {
	return ((bs[pos>>6] >> uint(pos&63)) & 1) != 0
}

func shiftOr(src, dst []uint64, shift int) {
	ws := shift >> 6
	bs := uint(shift & 63)
	if bs == 0 {
		for i, v := range src {
			if v == 0 {
				continue
			}
			j := i + ws
			if j < len(dst) {
				dst[j] |= v
			}
		}
	} else {
		rb := 64 - bs
		for i, v := range src {
			if v == 0 {
				continue
			}
			j := i + ws
			if j < len(dst) {
				dst[j] |= v << bs
			}
			if j+1 < len(dst) {
				dst[j+1] |= v >> rb
			}
		}
	}
}

func keepRange(bs []uint64, l, h int) {
	if l < 0 {
		l = 0
	}
	if h < l {
		for i := range bs {
			bs[i] = 0
		}
		return
	}
	for i := 0; i < len(bs); i++ {
		lo := i * 64
		hi := lo + 63
		if hi < l || lo > h {
			bs[i] = 0
			continue
		}
		mask := ^uint64(0)
		if l > lo {
			start := uint(l - lo)
			mask &= ^((uint64(1) << start) - 1)
		}
		if h < hi {
			end := uint(h - lo + 1)
			mask &= (uint64(1) << end) - 1
		}
		bs[i] &= mask
	}
}

func buildGraph(nodes []solveNode, mat [][]byte) {
	if len(nodes) <= 1 {
		return
	}
	sort.Slice(nodes, func(i, j int) bool {
		if nodes[i].deg != nodes[j].deg {
			return nodes[i].deg < nodes[j].deg
		}
		return nodes[i].id < nodes[j].id
	})
	v := nodes[0]
	d := v.deg
	rem := make([]solveNode, len(nodes)-1)
	for i := 1; i < len(nodes); i++ {
		rem[i-1] = nodes[i]
		if i-1 < d {
			mat[v.id][nodes[i].id] = '1'
		} else {
			mat[nodes[i].id][v.id] = '1'
			rem[i-1].deg--
		}
	}
	buildGraph(rem, mat)
}

func solveD(scores []int) (int, []int, bool) {
	m := len(scores)
	minScore := scores[0]
	maxScore := scores[m-1]

	prefix := make([]int, m+1)
	for i := 0; i < m; i++ {
		prefix[i+1] = prefix[i] + scores[i]
	}
	baseSum := prefix[m]

	minN := m
	if maxScore+1 > minN {
		minN = maxScore + 1
	}
	if 2*minScore+1 > minN {
		minN = 2*minScore + 1
	}
	maxN := 2*maxScore + 1
	if maxN > 61 {
		maxN = 61
	}
	if minN > maxN {
		return 0, nil, false
	}

	for n := minN; n <= maxN; n++ {
		r := n - m
		total := comb2(n)
		targetExtra := total - baseSum
		if targetExtra < 0 {
			continue
		}
		if targetExtra < r*minScore || targetExtra > r*maxScore {
			continue
		}

		w := targetExtra/64 + 1
		stride := (r + 1) * w
		data := make([]uint64, (m+1)*stride)
		data[0] = 1

		for p := 0; p < m; p++ {
			curBase := p * stride
			nextBase := (p + 1) * stride
			for i := 0; i < stride; i++ {
				data[nextBase+i] = 0
			}

			score := scores[p]
			for used := 0; used <= r; used++ {
				srcStart := curBase + used*w
				src := data[srcStart : srcStart+w]
				nonZero := false
				for _, v := range src {
					if v != 0 {
						nonZero = true
						break
					}
				}
				if !nonZero {
					continue
				}
				for add := 0; add <= r-used; add++ {
					shift := add * score
					if score > 0 && shift > targetExtra {
						break
					}
					dstStart := nextBase + (used+add)*w
					shiftOr(src, data[dstStart:dstStart+w], shift)
				}
			}

			for used := 0; used <= r; used++ {
				bs := data[nextBase+used*w : nextBase+(used+1)*w]
				if p+1 == m {
					if used == r {
						keepRange(bs, targetExtra, targetExtra)
					} else {
						for i := range bs {
							bs[i] = 0
						}
					}
				} else {
					k := (p + 1) + used
					l := comb2(k) - prefix[p+1]
					rem := r - used
					low := targetExtra - rem*maxScore
					if low > l {
						l = low
					}
					h := targetExtra - rem*scores[p+1]
					if h > targetExtra {
						h = targetExtra
					}
					keepRange(bs, l, h)
				}
			}
		}

		finalState := data[m*stride+r*w : m*stride+(r+1)*w]
		if !hasBit(finalState, targetExtra) {
			continue
		}

		extra := make([]int, m)
		used := r
		cur := targetExtra
		ok := true
		for p := m; p >= 1; p-- {
			score := scores[p-1]
			prevBase := (p - 1) * stride
			found := false
			for add := 0; add <= used; add++ {
				prev := cur - add*score
				if prev < 0 {
					if score > 0 {
						break
					}
					continue
				}
				bs := data[prevBase+(used-add)*w : prevBase+(used-add+1)*w]
				if hasBit(bs, prev) {
					extra[p-1] = add
					used -= add
					cur = prev
					found = true
					break
				}
			}
			if !found {
				ok = false
				break
			}
		}
		if !ok || used != 0 || cur != 0 {
			continue
		}

		counts := make([]int, m)
		for i := 0; i < m; i++ {
			counts[i] = 1 + extra[i]
		}
		return n, counts, true
	}

	return 0, nil, false
}

func refSolve(input string) string {
	reader := strings.NewReader(input)
	var m int
	fmt.Fscan(reader, &m)
	scores := make([]int, m)
	for i := 0; i < m; i++ {
		fmt.Fscan(reader, &scores[i])
	}
	sort.Ints(scores)

	n, counts, ok := solveD(scores)
	if !ok {
		return "=("
	}

	degrees := make([]int, 0, n)
	for i, s := range scores {
		for j := 0; j < counts[i]; j++ {
			degrees = append(degrees, s)
		}
	}

	mat := make([][]byte, n)
	for i := 0; i < n; i++ {
		row := make([]byte, n)
		for j := 0; j < n; j++ {
			row[j] = '0'
		}
		mat[i] = row
	}

	nodes := make([]solveNode, n)
	for i, d := range degrees {
		nodes[i] = solveNode{deg: d, id: i}
	}
	buildGraph(nodes, mat)

	var sb strings.Builder
	fmt.Fprintln(&sb, n)
	for i := 0; i < n; i++ {
		sb.Write(mat[i])
		sb.WriteByte('\n')
	}
	return strings.TrimSpace(sb.String())
}

func runExe(path string, input string) (string, error) {
	cmd := exec.Command(path)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

func genTests() []Test {
	rand.Seed(0)
	tests := make([]Test, 0, 101)
	for i := 0; i < 100; i++ {
		m := rand.Intn(5) + 1
		values := rand.Perm(31)[:m]
		var sb strings.Builder
		sb.WriteString(strconv.Itoa(m) + "\n")
		for j, v := range values {
			if j > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.Itoa(v))
		}
		sb.WriteByte('\n')
		tests = append(tests, Test{sb.String()})
	}
	tests = append(tests, Test{"1\n0\n"})
	return tests
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		return
	}
	bin := os.Args[1]

	tests := genTests()
	for i, tc := range tests {
		exp := refSolve(tc.input)
		got, err := runExe(bin, tc.input)
		if err != nil {
			fmt.Printf("candidate runtime error on test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(exp) != strings.TrimSpace(got) {
			fmt.Printf("Test %d failed\nInput:\n%sExpected:\n%s\nGot:\n%s\n", i+1, tc.input, exp, got)
			os.Exit(1)
		}
	}
	fmt.Println("all tests passed")
}
