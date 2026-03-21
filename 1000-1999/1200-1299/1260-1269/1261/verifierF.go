package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

const testCount = 120
const mod int64 = 998244353

// ──────────────────────────────────────────────────────────────────────
// Embedded correct solver: trie-based XOR convolution.
// Mirrors the accepted solution logic.
// ──────────────────────────────────────────────────────────────────────

var sInv2 int64 = 499122177

func power(base, exp int64) int64 {
	res := int64(1)
	base %= mod
	for exp > 0 {
		if exp%2 == 1 {
			res = (res * base) % mod
		}
		base = (base * base) % mod
		exp /= 2
	}
	return res
}

type sNode struct {
	id     int32
	isFull bool
	left   *sNode
	right  *sNode
}

var sFullNode = &sNode{id: 1, isFull: true}
var sNextNodeId int32 = 2

type sNodeKey struct {
	left, right int32
}

var sInternMap = make(map[sNodeKey]*sNode)
var sMemoXor = make(map[sNodeKey]*sNode)
var sMemoMerge = make(map[sNodeKey]*sNode)

func sGetNode(left, right *sNode) *sNode {
	if left == sFullNode && right == sFullNode {
		return sFullNode
	}
	if left == nil && right == nil {
		return nil
	}
	var lid, rid int32 = 0, 0
	if left != nil {
		lid = left.id
	}
	if right != nil {
		rid = right.id
	}
	k := sNodeKey{lid, rid}
	if n, ok := sInternMap[k]; ok {
		return n
	}
	n := &sNode{id: sNextNodeId, isFull: false, left: left, right: right}
	sNextNodeId++
	sInternMap[k] = n
	return n
}

func sMerge(a, b *sNode) *sNode {
	if a == nil {
		return b
	}
	if b == nil {
		return a
	}
	if a == b {
		return a
	}
	if a.isFull || b.isFull {
		return sFullNode
	}
	lid, rid := a.id, b.id
	if lid > rid {
		lid, rid = rid, lid
	}
	k := sNodeKey{lid, rid}
	if res, ok := sMemoMerge[k]; ok {
		return res
	}
	left := sMerge(a.left, b.left)
	right := sMerge(a.right, b.right)
	res := sGetNode(left, right)
	sMemoMerge[k] = res
	return res
}

func sXorTrie(a, b *sNode) *sNode {
	if a == nil || b == nil {
		return nil
	}
	if a.isFull || b.isFull {
		return sFullNode
	}
	lid, rid := a.id, b.id
	if lid > rid {
		lid, rid = rid, lid
	}
	k := sNodeKey{lid, rid}
	if res, ok := sMemoXor[k]; ok {
		return res
	}
	ll := sXorTrie(a.left, b.left)
	rr := sXorTrie(a.right, b.right)
	lr := sXorTrie(a.left, b.right)
	rl := sXorTrie(a.right, b.left)
	left := sMerge(ll, rr)
	right := sMerge(lr, rl)
	res := sGetNode(left, right)
	sMemoXor[k] = res
	return res
}

func sBuild(d int, L, R, valL, valR int64) *sNode {
	if L > valR || R < valL {
		return nil
	}
	if L <= valL && valR <= R {
		return sFullNode
	}
	mid := valL + (valR-valL)/2
	left := sBuild(d-1, L, R, valL, mid)
	right := sBuild(d-1, L, R, mid+1, valR)
	return sGetNode(left, right)
}

type sSumKey struct {
	id int32
	d  int32
}

var sMemoSum = make(map[sSumKey]struct{ count, sum int64 })

func sCompute(node *sNode, d int) (int64, int64) {
	if node == nil {
		return 0, 0
	}
	if node == sFullNode {
		cnt := power(2, int64(d+1))
		sum := (cnt * (cnt - 1 + mod)) % mod
		sum = (sum * sInv2) % mod
		return cnt, sum
	}
	k := sSumKey{node.id, int32(d)}
	if res, ok := sMemoSum[k]; ok {
		return res.count, res.sum
	}
	cntL, sumL := sCompute(node.left, d-1)
	cntR, sumR := sCompute(node.right, d-1)
	cnt := (cntL + cntR) % mod
	bitVal := power(2, int64(d))
	sumR_add := (cntR * bitVal) % mod
	sumR_total := (sumR + sumR_add) % mod
	sum := (sumL + sumR_total) % mod
	sMemoSum[k] = struct{ count, sum int64 }{cnt, sum}
	return cnt, sum
}

func resetSolverState() {
	sNextNodeId = 2
	sInternMap = make(map[sNodeKey]*sNode)
	sMemoXor = make(map[sNodeKey]*sNode)
	sMemoMerge = make(map[sNodeKey]*sNode)
	sMemoSum = make(map[sSumKey]struct{ count, sum int64 })
}

func solveEmbedded(segA, segB [][2]uint64) int64 {
	resetSolverState()

	var rootA *sNode
	for _, seg := range segA {
		node := sBuild(59, int64(seg[0]), int64(seg[1]), 0, (int64(1)<<60)-1)
		rootA = sMerge(rootA, node)
	}

	var rootB *sNode
	for _, seg := range segB {
		node := sBuild(59, int64(seg[0]), int64(seg[1]), 0, (int64(1)<<60)-1)
		rootB = sMerge(rootB, node)
	}

	rootC := sXorTrie(rootA, rootB)
	_, sumC := sCompute(rootC, 59)
	return sumC
}

// ──────────────────────────────────────────────────────────────────────

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

func genSegments(r *rand.Rand, maxSeg int, maxVal uint64) ([][2]uint64, int) {
	count := 1 + r.Intn(maxSeg)
	segs := make([][2]uint64, count)
	for i := 0; i < count; i++ {
		l := r.Uint64()%(maxVal/2+1) + 1
		length := r.Uint64()%uint64(1+int(maxVal/10)) + 1
		rVal := l + length
		if rVal > maxVal {
			rVal = maxVal
		}
		if rVal < l {
			l, rVal = rVal, l
		}
		segs[i] = [2]uint64{l, rVal}
	}
	return segs, count
}

func formatInput(segA [][2]uint64, segB [][2]uint64) string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", len(segA))
	for _, seg := range segA {
		fmt.Fprintf(&sb, "%d %d\n", seg[0], seg[1])
	}
	fmt.Fprintf(&sb, "%d\n", len(segB))
	for _, seg := range segB {
		fmt.Fprintf(&sb, "%d %d\n", seg[0], seg[1])
	}
	return sb.String()
}

func parseOutput(out string) (int64, error) {
	fields := strings.Fields(out)
	if len(fields) != 1 {
		return 0, fmt.Errorf("expected single integer, got %d fields", len(fields))
	}
	val, err := strconv.ParseInt(fields[0], 10, 64)
	if err != nil {
		return 0, err
	}
	val %= int64(mod)
	if val < 0 {
		val += int64(mod)
	}
	return val, nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	userBin := os.Args[1]

	r := rand.New(rand.NewSource(1))
	for t := 0; t < testCount; t++ {
		segA, _ := genSegments(r, 5, uint64(1e6))
		segB, _ := genSegments(r, 5, uint64(1e6))
		input := formatInput(segA, segB)
		expectVal := solveEmbedded(segA, segB)
		gotStr, err := run(userBin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "test %d: %v\n", t+1, err)
			os.Exit(1)
		}
		gotVal, err := parseOutput(gotStr)
		if err != nil {
			fmt.Printf("test %d failed\ninput:\n%s\nerror: %v\n", t+1, input, err)
			os.Exit(1)
		}
		if expectVal != gotVal {
			fmt.Printf("test %d failed\ninput:\n%s\nexpected: %d\ngot: %d\n", t+1, input, expectVal, gotVal)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", testCount)
}
