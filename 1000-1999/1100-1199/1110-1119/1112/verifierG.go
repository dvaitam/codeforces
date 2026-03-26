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

// ─── embedded correct solver ───

const solverMapSize = 1 << 23

var solverKeys []int64
var solverValues []int32

func solverInitMap() {
	solverKeys = make([]int64, solverMapSize)
	solverValues = make([]int32, solverMapSize)
	for i := 0; i < solverMapSize; i++ {
		solverKeys[i] = -1
	}
}

func solverPut(key int64, val int32) {
	idx := (uint64(key) * 11400714819323198485) >> 41
	mask := uint64(solverMapSize - 1)
	idx &= mask
	for solverKeys[idx] != -1 && solverKeys[idx] != key {
		idx = (idx + 1) & mask
	}
	solverKeys[idx] = key
	solverValues[idx] = val
}

func solverGet(key int64) (int32, bool) {
	idx := (uint64(key) * 11400714819323198485) >> 41
	mask := uint64(solverMapSize - 1)
	idx &= mask
	for solverKeys[idx] != -1 {
		if solverKeys[idx] == key {
			return solverValues[idx], true
		}
		idx = (idx + 1) & mask
	}
	return 0, false
}

type solverNode struct {
	parent int32
	diff   int32
	c      int16
	d      int8
}

func solverSumOfDigits(c int16) int32 {
	s := int32(0)
	for c > 0 {
		s += int32(c % 10)
		c /= 10
	}
	return s
}

func solveG(input string) string {
	var a int
	fmt.Sscan(input, &a)

	solverInitMap()

	nodes := make([]solverNode, 0, 1000000)
	currentQ := make([]int32, 0, 100000)
	nextQ := make([]int32, 0, 100000)

	nodes = append(nodes, solverNode{parent: -1, d: -1, c: 0, diff: 0})
	currentQ = append(currentQ, 0)

	startKey := (int64(0) << 32) | int64(uint32(0))
	solverPut(startKey, 0)

	for len(currentQ) > 0 {
		for _, currIdx := range currentQ {
			curr := nodes[currIdx]

			for d := 0; d <= 9; d++ {
				c := curr.c
				diff := curr.diff

				nc := int16((d*a + int(c)) / 10)
				mDigit := int32((d*a + int(c)) % 10)
				ndiff := diff + int32(d) - int32(a)*mDigit

				if ndiff < -5000000 || ndiff > 5000000 {
					continue
				}

				if d != 0 && ndiff == int32(a)*solverSumOfDigits(nc) {
					var buf bytes.Buffer
					buf.WriteByte(byte(d + '0'))
					p := currIdx
					for p != 0 {
						buf.WriteByte(byte(nodes[p].d + '0'))
						p = nodes[p].parent
					}
					return buf.String()
				}

				key := (int64(nc) << 32) | int64(uint32(ndiff))
				if _, ok := solverGet(key); !ok {
					if len(nodes) >= solverMapSize/4*3 {
						return "-1"
					}
					nodes = append(nodes, solverNode{
						parent: currIdx,
						diff:   ndiff,
						c:      nc,
						d:      int8(d),
					})
					idx := int32(len(nodes) - 1)
					solverPut(key, idx)
					nextQ = append(nextQ, idx)
				}
			}
		}
		currentQ, nextQ = nextQ, currentQ
		nextQ = nextQ[:0]
	}
	return "-1"
}

// ─── verifier ───

func run(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out, stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func genCase(r *rand.Rand) string {
	a := r.Intn(999) + 2
	return fmt.Sprintf("%d\n", a)
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
		expect := solveG(input)
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
