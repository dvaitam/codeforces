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

func runBinary(bin, input string) (string, error) {
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
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func solve(n, m int, W int64, w, c, a []int64) string {
	type comp struct {
		price float64
		w     int64
	}
	comps := make([]comp, m)
	var answer float64
	for day := 0; day < n; day++ {
		for i := 0; i < m; i++ {
			price := (float64(c[i]) - float64(a[i])*float64(day)) / float64(w[i])
			comps[i] = comp{price: price, w: w[i]}
		}
		sort.Slice(comps, func(i, j int) bool { return comps[i].price < comps[j].price })
		rem := W
		for i := 0; i < m && rem > 0; i++ {
			if comps[i].w <= rem {
				answer += float64(comps[i].w) * comps[i].price
				rem -= comps[i].w
			} else {
				answer += float64(rem) * comps[i].price
				rem = 0
			}
		}
	}
	return fmt.Sprintf("%.9f", answer)
}

func generateCase(r *rand.Rand) (string, string) {
	n := r.Intn(5) + 1
	m := r.Intn(5) + 1
	W := int64(r.Intn(20) + 1)
	w := make([]int64, m)
	c := make([]int64, m)
	a := make([]int64, m)
	for i := 0; i < m; i++ {
		w[i] = int64(r.Intn(5) + 1)
		c[i] = int64(r.Intn(10) + 1)
		a[i] = int64(r.Intn(5))
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d %d\n", n, m, W))
	for i := 0; i < m; i++ {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprint(w[i]))
	}
	sb.WriteByte('\n')
	for i := 0; i < m; i++ {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprint(c[i]))
	}
	sb.WriteByte('\n')
	for i := 0; i < m; i++ {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprint(a[i]))
	}
	sb.WriteByte('\n')
	return sb.String(), solve(n, m, W, w, c, a)
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := generateCase(r)
		out, err := runBinary(bin, in)
		if err != nil {
			fmt.Printf("Test %d runtime error: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(out) != exp {
			fmt.Printf("Test %d failed.\nInput:\n%sExpected: %s\nGot: %s\n", i+1, in, exp, out)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
