package main

import (
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strings"
)

type BIT struct {
	n    int
	tree []int
}

func NewBIT(n int) *BIT {
	return &BIT{n: n, tree: make([]int, n+2)}
}

func (b *BIT) Add(i, delta int) {
	for i <= b.n {
		b.tree[i] += delta
		i += i & -i
	}
}

func (b *BIT) Sum(i int) int {
	if i > b.n {
		i = b.n
	}
	s := 0
	for i > 0 {
		s += b.tree[i]
		i &= i - 1
	}
	return s
}

type Station struct {
	x   int
	r   int
	f   int
	idx int
}

func uniqueInts(a []int) []int {
	if len(a) == 0 {
		return a
	}
	j := 1
	for i := 1; i < len(a); i++ {
		if a[i] != a[i-1] {
			a[j] = a[i]
			j++
		}
	}
	return a[:j]
}

func expected(n, k int, stations []Station) int64 {
	freqCoords := make(map[int][]int)
	for i := 0; i < n; i++ {
		freqCoords[stations[i].f] = append(freqCoords[stations[i].f], stations[i].x)
	}
	bits := make(map[int]*BIT, len(freqCoords))
	coords := make(map[int][]int, len(freqCoords))
	for f, arr := range freqCoords {
		sort.Ints(arr)
		arr = uniqueInts(arr)
		coords[f] = arr
		bits[f] = NewBIT(len(arr))
	}
	for i := range stations {
		arr := coords[stations[i].f]
		stations[i].idx = sort.SearchInts(arr, stations[i].x) + 1
	}
	sort.Slice(stations, func(i, j int) bool { return stations[i].x < stations[j].x })
	ans := int64(0)
	left := 0
	for i, st := range stations {
		for left < i && stations[left].x+stations[left].r < st.x {
			rm := stations[left]
			bits[rm.f].Add(rm.idx, -1)
			left++
		}
		for freq := st.f - k; freq <= st.f+k; freq++ {
			if freq < 1 || freq > 10000 {
				continue
			}
			bit := bits[freq]
			if bit == nil {
				continue
			}
			arr := coords[freq]
			pos := sort.SearchInts(arr, st.x-st.r)
			cnt := bit.Sum(len(arr)) - bit.Sum(pos)
			ans += int64(cnt)
		}
		bits[st.f].Add(st.idx, 1)
	}
	return ans
}

func run(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	out, err := cmd.CombinedOutput()
	return strings.TrimSpace(string(out)), err
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(42))

	for t := 0; t < 100; t++ {
		n := rng.Intn(20) + 1
		k := rng.Intn(5)
		stations := make([]Station, n)
		usedX := map[int]bool{}
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d %d\n", n, k))
		for i := 0; i < n; i++ {
			x := rng.Intn(1000) + 1
			for usedX[x] {
				x = rng.Intn(1000) + 1
			}
			usedX[x] = true
			r := rng.Intn(10) + 1
			f := rng.Intn(10) + 1
			stations[i] = Station{x: x, r: r, f: f}
			sb.WriteString(fmt.Sprintf("%d %d %d\n", x, r, f))
		}
		input := sb.String()
		exp := fmt.Sprintf("%d", expected(n, k, stations))
		got, err := run(bin, input)
		if err != nil {
			fmt.Printf("runtime error on test %d: %v\ninput:\n%s\noutput:\n%s\n", t+1, err, input, got)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != exp {
			fmt.Printf("wrong answer on test %d\ninput:\n%s\nexpected: %s\ngot: %s\n", t+1, input, exp, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed.")
}
