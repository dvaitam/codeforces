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

func run(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errb bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errb
	err := cmd.Run()
	if err != nil {
		return out.String() + errb.String(), err
	}
	return out.String(), nil
}

type plane struct {
	x int
	v int
}

type Test struct {
	n      int
	w      int
	planes []plane
	input  string
}

func genTest(rng *rand.Rand) Test {
	n := rng.Intn(7) + 1
	w := rng.Intn(10)
	planes := make([]plane, n)
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", n, w))
	for i := 0; i < n; i++ {
		x := rng.Intn(20) + 1
		if rng.Intn(2) == 0 {
			x = -x
		}
		v := rng.Intn(20) + w + 1
		if rng.Intn(2) == 0 {
			v = -v
		}
		if x*v > 0 {
			v = -v
		}
		planes[i] = plane{x: x, v: v}
		sb.WriteString(fmt.Sprintf("%d %d\n", x, v))
	}
	return Test{n: n, w: w, planes: planes, input: sb.String()}
}

func solve(t Test) string {
	type pair struct{ a, b int64 }
	prec := int64(1e9)
	planes := make([]pair, len(t.planes))
	for i, p := range t.planes {
		a := -float64(p.x) / float64(p.v-t.w)
		b := -float64(p.x) / float64(p.v+t.w)
		planes[i] = pair{int64(a*float64(prec) + 0.5), int64(b*float64(prec) + 0.5)}
	}
	sort.Slice(planes, func(i, j int) bool { return planes[i].a < planes[j].a })
	bvals := make([]int64, len(planes))
	for i := range planes {
		bvals[i] = planes[i].b
	}
	uniqMap := map[int64]struct{}{}
	uniq := make([]int64, 0, len(bvals))
	for _, v := range bvals {
		if _, ok := uniqMap[v]; !ok {
			uniqMap[v] = struct{}{}
			uniq = append(uniq, v)
		}
	}
	sort.Slice(uniq, func(i, j int) bool { return uniq[i] < uniq[j] })
	rank := make(map[int64]int, len(uniq))
	for i, v := range uniq {
		rank[v] = i + 1
	}
	bit := make([]int64, len(uniq)+3)
	add := func(i int64, val int64) {
		for i < int64(len(bit)) {
			bit[i] += val
			i += i & -i
		}
	}
	sum := func(i int64) int64 {
		s := int64(0)
		for i > 0 {
			s += bit[i]
			i -= i & -i
		}
		return s
	}
	total := int64(0)
	ans := int64(0)
	i := 0
	for i < len(planes) {
		j := i
		for j < len(planes) && planes[j].a == planes[i].a {
			j++
		}
		for k := i; k < j; k++ {
			r := int64(rank[planes[k].b])
			ans += total - sum(r-1)
		}
		for k := i; k < j; k++ {
			r := int64(rank[planes[k].b])
			add(r, 1)
			total++
		}
		m := int64(j - i)
		ans += m * (m - 1) / 2
		i = j
	}
	return fmt.Sprintf("%d", ans)
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		t := genTest(rng)
		expected := solve(t)
		out, err := run(bin, t.input)
		if err != nil {
			fmt.Printf("test %d runtime error: %v\n%s", i+1, err, out)
			os.Exit(1)
		}
		out = strings.TrimSpace(out)
		if out != expected {
			fmt.Printf("test %d failed\ninput:\n%s\nexpected:%s got:%s\n", i+1, t.input, expected, out)
			os.Exit(1)
		}
	}
	fmt.Println("ok 100 tests")
}
