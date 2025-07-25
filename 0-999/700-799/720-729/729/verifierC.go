package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strings"
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
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func can(v int64, d []int64, t int64, maxD int64) bool {
	if v < maxD {
		return false
	}
	var total int64
	for _, dist := range d {
		if v >= 2*dist {
			total += dist
		} else {
			total += 3*dist - v
		}
		if total > t {
			return false
		}
	}
	return total <= t
}

func expected(n, k int, s, t int64, cars [][2]int64, stations []int64) string {
	all := make([]int64, 0, k+2)
	all = append(all, 0)
	all = append(all, stations...)
	all = append(all, s)
	sort.Slice(all, func(i, j int) bool { return all[i] < all[j] })
	m := len(all) - 1
	d := make([]int64, m)
	var maxD int64
	for i := 0; i < m; i++ {
		dist := all[i+1] - all[i]
		d[i] = dist
		if dist > maxD {
			maxD = dist
		}
	}
	var maxV int64
	for _, cv := range cars {
		if cv[1] > maxV {
			maxV = cv[1]
		}
	}
	l, r := int64(0), maxV+1
	for l < r {
		mid := (l + r) / 2
		if can(mid, d, t, maxD) {
			r = mid
		} else {
			l = mid + 1
		}
	}
	if l > maxV || !can(l, d, t, maxD) {
		return "-1"
	}
	vReq := l
	ans := int64(1<<62 - 1)
	for _, cv := range cars {
		if cv[1] >= vReq && cv[0] < ans {
			ans = cv[0]
		}
	}
	if ans == int64(1<<62-1) {
		return "-1"
	}
	return fmt.Sprintf("%d", ans)
}

func main() {
	if len(os.Args) != 2 && !(len(os.Args) == 3 && os.Args[1] == "--") {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[len(os.Args)-1]
	r := rand.New(rand.NewSource(3))
	for tc := 1; tc <= 100; tc++ {
		n := r.Intn(4) + 1
		k := r.Intn(3) + 1
		sVal := int64(r.Intn(15) + 5)
		tVal := int64(r.Intn(40) + int(sVal) + 1)
		cars := make([][2]int64, n)
		for i := 0; i < n; i++ {
			cars[i][0] = int64(r.Intn(50) + 1)
			cars[i][1] = int64(r.Intn(30) + 1)
		}
		posSet := make(map[int64]bool)
		stations := make([]int64, 0, k)
		for len(stations) < k {
			p := int64(r.Intn(int(sVal-1)) + 1)
			if !posSet[p] {
				posSet[p] = true
				stations = append(stations, p)
			}
		}
		var sb strings.Builder
		fmt.Fprintf(&sb, "%d %d %d %d\n", n, k, sVal, tVal)
		for i := 0; i < n; i++ {
			fmt.Fprintf(&sb, "%d %d\n", cars[i][0], cars[i][1])
		}
		for i := 0; i < k; i++ {
			if i > 0 {
				sb.WriteByte(' ')
			}
			fmt.Fprintf(&sb, "%d", stations[i])
		}
		sb.WriteByte('\n')
		input := sb.String()
		expect := expected(n, k, sVal, tVal, cars, stations)
		out, err := run(bin, input)
		if err != nil {
			fmt.Printf("test %d: %v\n", tc, err)
			os.Exit(1)
		}
		if out != expect {
			fmt.Printf("test %d failed\ninput:\n%sexpected: %s\ngot: %s\n", tc, input, expect, out)
			os.Exit(1)
		}
	}
	fmt.Println("All 100 tests passed")
}
