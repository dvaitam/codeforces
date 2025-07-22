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

type Test struct {
	n, m int
	x, y int64
	a, b []int64
}

func generateTest() Test {
	n := rand.Intn(20) + 1
	m := rand.Intn(20) + 1
	x := int64(rand.Intn(5))
	y := int64(rand.Intn(5))
	a := make([]int64, n)
	b := make([]int64, m)
	cur := rand.Intn(3) + 1
	for i := 0; i < n; i++ {
		cur += rand.Intn(3)
		a[i] = int64(cur)
	}
	cur = rand.Intn(3) + 1
	for j := 0; j < m; j++ {
		cur += rand.Intn(3)
		b[j] = int64(cur)
	}
	return Test{n, m, x, y, a, b}
}

func (t Test) Input() string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d %d %d\n", t.n, t.m, t.x, t.y)
	for i, v := range t.a {
		if i > 0 {
			sb.WriteByte(' ')
		}
		fmt.Fprintf(&sb, "%d", v)
	}
	sb.WriteByte('\n')
	for i, v := range t.b {
		if i > 0 {
			sb.WriteByte(' ')
		}
		fmt.Fprintf(&sb, "%d", v)
	}
	sb.WriteByte('\n')
	return sb.String()
}

func solve(t Test) (int, [][2]int) {
	i, j := 0, 0
	res := make([][2]int, 0)
	for i < t.n && j < t.m {
		low := t.a[i] - t.x
		high := t.a[i] + t.y
		if t.b[j] < low {
			j++
		} else if t.b[j] > high {
			i++
		} else {
			res = append(res, [2]int{i + 1, j + 1})
			i++
			j++
		}
	}
	return len(res), res
}

func runBinary(path string, input string) (string, error) {
	cmd := exec.Command(path)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	return out.String(), err
}

func checkOutput(output string, t Test, expectedK int) error {
	r := strings.NewReader(output)
	var k int
	if _, err := fmt.Fscan(r, &k); err != nil {
		return fmt.Errorf("failed to read k: %v", err)
	}
	if k != expectedK {
		return fmt.Errorf("expected %d pairs, got %d", expectedK, k)
	}
	usedS := make(map[int]bool)
	usedV := make(map[int]bool)
	for i := 0; i < k; i++ {
		var u, v int
		if _, err := fmt.Fscan(r, &u, &v); err != nil {
			return fmt.Errorf("failed to read pair %d: %v", i+1, err)
		}
		if u < 1 || u > t.n || v < 1 || v > t.m {
			return fmt.Errorf("invalid indices in pair %d", i+1)
		}
		if usedS[u] || usedV[v] {
			return fmt.Errorf("duplicate indices in pair %d", i+1)
		}
		usedS[u] = true
		usedV[v] = true
		if t.b[v-1] < t.a[u-1]-t.x || t.b[v-1] > t.a[u-1]+t.y {
			return fmt.Errorf("pair %d does not satisfy constraints", i+1)
		}
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go <binary>")
		os.Exit(1)
	}
	binary := os.Args[1]
	rand.Seed(time.Now().UnixNano())
	for tcase := 0; tcase < 100; tcase++ {
		t := generateTest()
		inp := t.Input()
		expK, _ := solve(t)
		out, err := runBinary(binary, inp)
		if err != nil {
			fmt.Fprintf(os.Stderr, "test %d: runtime error: %v\n", tcase+1, err)
			os.Exit(1)
		}
		if err := checkOutput(out, t, expK); err != nil {
			fmt.Fprintf(os.Stderr, "test %d failed: %v\ninput:\n%s\noutput:\n%s\n", tcase+1, err, inp, out)
			os.Exit(1)
		}
	}
	fmt.Println("all tests passed")
}
