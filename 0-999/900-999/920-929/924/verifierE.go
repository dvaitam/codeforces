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

type Test struct {
	n     int
	l     int
	r     int
	a     []int
	b     []int
	input string
}

func genTest(rng *rand.Rand) Test {
	n := rng.Intn(8) + 1
	sum := 0
	a := make([]int, n)
	for i := 0; i < n; i++ {
		a[i] = rng.Intn(5) + 1
		sum += a[i]
	}
	l := rng.Intn(sum + 1)
	r := l + rng.Intn(sum-l+1)
	b := make([]int, n)
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d %d\n", n, l, r))
	for i, v := range a {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", v))
	}
	sb.WriteByte('\n')
	for i := 0; i < n; i++ {
		if rng.Intn(2) == 0 {
			b[i] = 0
		} else {
			b[i] = 1
		}
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", b[i]))
	}
	sb.WriteByte('\n')
	return Test{n: n, l: l, r: r, a: a, b: b, input: sb.String()}
}

func solve(t Test) string {
	H := 0
	for _, v := range t.a {
		H += v
	}
	L := H - t.r
	R := H - t.l
	type box struct{ h, b int }
	boxes := make([]box, t.n)
	for i := 0; i < t.n; i++ {
		boxes[i] = box{t.a[i], t.b[i]}
	}
	sort.Slice(boxes, func(i, j int) bool {
		if boxes[i].b != boxes[j].b {
			return boxes[i].b < boxes[j].b
		}
		return boxes[i].h > boxes[j].h
	})
	const negInf = -1 << 30
	dp := make([]int, H+1)
	for i := range dp {
		dp[i] = negInf
	}
	dp[0] = 0
	for _, bx := range boxes {
		for j := H; j >= bx.h; j-- {
			if dp[j-bx.h] == negInf {
				continue
			}
			gain := 0
			if bx.b == 1 && j >= L && j <= R {
				gain = 1
			}
			if v := dp[j-bx.h] + gain; v > dp[j] {
				dp[j] = v
			}
		}
	}
	ans := 0
	for _, v := range dp {
		if v > ans {
			ans = v
		}
	}
	return fmt.Sprintf("%d", ans)
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
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
	fmt.Println("All tests passed")
}
