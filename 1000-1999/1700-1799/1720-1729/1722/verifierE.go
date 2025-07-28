package main

import (
	"bytes"
	"context"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

const maxHW = 1000

type rectangle struct{ h, w int }
type query struct{ hs, ws, hb, wb int }

type testCaseE struct {
	n       int
	rects   []rectangle
	q       int
	queries []query
}

func solveE(n int, rects []rectangle, q int, queries []query) []int64 {
	rect := make([][]int64, maxHW+1)
	for i := range rect {
		rect[i] = make([]int64, maxHW+1)
	}
	for _, r := range rects {
		if r.h <= maxHW && r.w <= maxHW {
			rect[r.h][r.w] += int64(r.h * r.w)
		}
	}
	prefix := make([][]int64, maxHW+1)
	for i := range prefix {
		prefix[i] = make([]int64, maxHW+1)
	}
	for i := 1; i <= maxHW; i++ {
		var rowSum int64
		for j := 1; j <= maxHW; j++ {
			rowSum += rect[i][j]
			prefix[i][j] = prefix[i-1][j] + rowSum
		}
	}
	ans := make([]int64, q)
	for idx, qu := range queries {
		if qu.hs+1 > qu.hb-1 || qu.ws+1 > qu.wb-1 {
			ans[idx] = 0
			continue
		}
		h2, w2 := qu.hb-1, qu.wb-1
		ans[idx] = prefix[h2][w2] - prefix[qu.hs][w2] - prefix[h2][qu.ws] + prefix[qu.hs][qu.ws]
	}
	return ans
}

func runBinary(bin, input string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func generateTests() []testCaseE {
	rand.Seed(46)
	tests := make([]testCaseE, 100)
	for i := range tests {
		n := rand.Intn(5) + 1
		rects := make([]rectangle, n)
		for j := 0; j < n; j++ {
			rects[j] = rectangle{rand.Intn(5) + 1, rand.Intn(5) + 1}
		}
		q := rand.Intn(5) + 1
		queries := make([]query, q)
		for j := 0; j < q; j++ {
			hs := rand.Intn(5)
			ws := rand.Intn(5)
			hb := hs + rand.Intn(5-hs) + 1
			wb := ws + rand.Intn(5-ws) + 1
			queries[j] = query{hs, ws, hb, wb}
		}
		tests[i] = testCaseE{n, rects, q, queries}
	}
	return tests
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		return
	}
	bin := os.Args[1]
	tests := generateTests()
	for i, tc := range tests {
		var sb strings.Builder
		sb.WriteString("1\n")
		sb.WriteString(fmt.Sprintf("%d %d\n", tc.n, tc.q))
		for _, r := range tc.rects {
			sb.WriteString(fmt.Sprintf("%d %d\n", r.h, r.w))
		}
		for _, qu := range tc.queries {
			sb.WriteString(fmt.Sprintf("%d %d %d %d\n", qu.hs, qu.ws, qu.hb, qu.wb))
		}
		input := sb.String()
		expectedAns := solveE(tc.n, tc.rects, tc.q, tc.queries)
		expected := make([]string, len(expectedAns))
		for j, v := range expectedAns {
			expected[j] = fmt.Sprint(v)
		}
		expectedStr := strings.Join(expected, "\n")
		output, err := runBinary(bin, input)
		if err != nil {
			fmt.Printf("test %d: execution error: %v\n", i+1, err)
			return
		}
		if strings.TrimSpace(output) != expectedStr {
			fmt.Printf("test %d failed:\ninput:%sexpected %s got %s\n", i+1, input, expectedStr, output)
			return
		}
	}
	fmt.Printf("all %d tests passed\n", len(tests))
}
