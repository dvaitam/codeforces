package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"
	"time"
)

type interval struct{ l, r int }

type testCase struct {
	n, k int
	segs []interval
}

func expected(tc testCase) (int, []int) {
	type node struct{ l, r, id int }
	arr := make([]node, tc.n)
	for i := 0; i < tc.n; i++ {
		arr[i] = node{tc.segs[i].l, tc.segs[i].r, i + 1}
	}
	sort.Slice(arr, func(i, j int) bool { return arr[i].l < arr[j].l })
	h := make([]int, 0)
	rmin := make([]int, tc.n)
	const INF = int(1e9)
	heapPush := func(x int) {
		h = append(h, x)
		for i := len(h)/2 - 1; i >= 0; i-- {
			for child := 2*i + 1; child < len(h); child = 2*child + 1 {
				if child+1 < len(h) && h[child+1] < h[child] {
					child++
				}
				if h[i] <= h[child] {
					break
				}
				h[i], h[child] = h[child], h[i]
			}
		}
	}
	heapPop := func() int {
		v := h[0]
		h[0] = h[len(h)-1]
		h = h[:len(h)-1]
		for i := 0; ; {
			l := 2*i + 1
			if l >= len(h) {
				break
			}
			r := l + 1
			if r < len(h) && h[r] < h[l] {
				l = r
			}
			if h[i] <= h[l] {
				break
			}
			h[i], h[l] = h[l], h[i]
			i = l
		}
		return v
	}
	heapTop := func() int { return h[0] }
	tot := 0
	for i := 0; i < tc.n; i++ {
		if tot == tc.k {
			t := heapTop()
			if t <= arr[i].r {
				heapPop()
				heapPush(arr[i].r)
			}
			rmin[i] = heapTop()
		} else {
			heapPush(arr[i].r)
			tot++
			if tot != tc.k {
				rmin[i] = -INF
			} else {
				rmin[i] = heapTop()
			}
		}
	}
	ans := 0
	ansi := 0
	for i := 0; i < tc.n; i++ {
		cur := rmin[i] - arr[i].l + 1
		if cur < 0 {
			cur = 0
		}
		if cur > ans {
			ans = cur
			ansi = i
		}
	}
	selected := make([]node, ansi+1)
	copy(selected, arr[:ansi+1])
	sort.Slice(selected, func(i, j int) bool { return selected[i].r > selected[j].r })
	ids := make([]int, 0, tc.k)
	for i := 0; i < tc.k && i < len(selected); i++ {
		ids = append(ids, selected[i].id)
	}
	sort.Ints(ids)
	return ans, ids
}

func runCase(bin string, tc testCase) error {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d\n", tc.n, tc.k)
	for i := 0; i < tc.n; i++ {
		fmt.Fprintf(&sb, "%d %d\n", tc.segs[i].l, tc.segs[i].r)
	}
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(sb.String())
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	scanner := bufio.NewScanner(bytes.NewReader(out.Bytes()))
	scanner.Split(bufio.ScanWords)
	if !scanner.Scan() {
		return fmt.Errorf("no output")
	}
	ansStr := scanner.Text()
	ansVal, err := strconv.Atoi(ansStr)
	if err != nil {
		return fmt.Errorf("bad ans")
	}
	ids := []int{}
	for scanner.Scan() {
		v, err := strconv.Atoi(scanner.Text())
		if err != nil {
			return fmt.Errorf("bad id")
		}
		ids = append(ids, v)
	}
	if len(ids) != tc.k {
		return fmt.Errorf("need %d ids", tc.k)
	}
	sort.Ints(ids)
	expAns, _ := expected(tc)
	if ansVal != expAns {
		return fmt.Errorf("expected ans %d got %d", expAns, ansVal)
	}
	interL := -1 << 30
	interR := 1 << 30
	for _, id := range ids {
		if id <= 0 || id > tc.n {
			return fmt.Errorf("invalid id")
		}
		seg := tc.segs[id-1]
		if seg.l > interL {
			interL = seg.l
		}
		if seg.r < interR {
			interR = seg.r
		}
	}
	length := interR - interL + 1
	if length < 0 {
		length = 0
	}
	if length != ansVal {
		return fmt.Errorf("ids intersection length %d mismatch", length)
	}
	return nil
}

func randomCase(rng *rand.Rand) testCase {
	n := rng.Intn(5) + 1
	k := rng.Intn(n) + 1
	segs := make([]interval, n)
	for i := 0; i < n; i++ {
		l := rng.Intn(20)
		r := l + rng.Intn(20)
		segs[i] = interval{l, r}
	}
	return testCase{n, k, segs}
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	cases := []testCase{
		{2, 1, []interval{{0, 5}, {3, 8}}},
		{3, 2, []interval{{0, 2}, {2, 4}, {1, 3}}},
	}
	for i := 0; i < 100; i++ {
		cases = append(cases, randomCase(rng))
	}
	for i, tc := range cases {
		if err := runCase(bin, tc); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
