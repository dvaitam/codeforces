package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strings"
	"time"
)

type rat struct {
	en, de int
}

func gcd(a, b int) int {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

func newRat(en, de int) rat {
	g := gcd(en, de)
	return rat{en / g, de / g}
}

func ratLess(a, b rat) bool {
	if a.de != b.de {
		return a.de < b.de
	}
	if a.en != b.en {
		return a.en < b.en
	}
	return false
}

func rev(x int) int {
	r := 0
	for x > 0 {
		r = r*10 + x%10
		x /= 10
	}
	return r
}

type pair struct {
	r rat
	x int
}

var S []pair

func initPairs() {
	const N = 100500
	S = make([]pair, N)
	for i := 1; i <= N; i++ {
		S[i-1] = pair{newRat(i, rev(i)), i}
	}
	sort.Slice(S, func(i, j int) bool {
		if S[i].r.de != S[j].r.de {
			return S[i].r.de < S[j].r.de
		}
		if S[i].r.en != S[j].r.en {
			return S[i].r.en < S[j].r.en
		}
		return S[i].x < S[j].x
	})
}

func col(x, mx int) int {
	if x == 0 {
		return 0
	}
	r := newRat(rev(x), x)
	n := len(S)
	lo := sort.Search(n, func(i int) bool {
		return !ratLess(S[i].r, r)
	})
	hi := sort.Search(n, func(i int) bool {
		if ratLess(r, S[i].r) {
			return true
		}
		if ratLess(S[i].r, r) {
			return false
		}
		return S[i].x > mx
	})
	return hi - lo
}

const INF int64 = 1 << 60

func solve(input string) string {
	in := bufio.NewScanner(strings.NewReader(input))
	in.Split(bufio.ScanWords)
	var mxx, mxy, w int
	if !in.Scan() {
		return ""
	}
	fmt.Sscan(in.Text(), &mxx)
	in.Scan()
	fmt.Sscan(in.Text(), &mxy)
	in.Scan()
	fmt.Sscan(in.Text(), &w)
	x := mxx + 1
	y := 0
	cur := int64(0)
	bestArea := INF
	bestX, bestY := -1, -1
	for x > 0 {
		cur -= int64(col(x, y))
		x--
		for cur < int64(w) && y <= mxy {
			y++
			cur += int64(col(y, x))
		}
		if y > mxy {
			break
		}
		area := int64(x) * int64(y)
		if area < bestArea {
			bestArea = area
			bestX = x
			bestY = y
		}
	}
	if bestArea >= INF {
		return "-1"
	}
	return fmt.Sprintf("%d %d", bestX, bestY)
}

func genCase(rng *rand.Rand) string {
	maxx := rng.Intn(500) + 1
	maxy := rng.Intn(500) + 1
	w := rng.Intn(100000) + 1
	return fmt.Sprintf("%d %d %d\n", maxx, maxy, w)
}

func runCase(bin string, input, expected string) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(out.String())
	if got != expected {
		return fmt.Errorf("expected %s got %s", expected, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	initPairs()
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		input := genCase(rng)
		expected := solve(input)
		if err := runCase(bin, input, expected); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:%s", i+1, err, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
