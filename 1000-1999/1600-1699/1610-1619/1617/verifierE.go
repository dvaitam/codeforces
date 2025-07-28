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

func getInv(v int) int {
	lg := 0
	for (1 << (lg + 1)) <= v {
		lg++
	}
	inv := (1 << (lg + 1)) - v
	if inv == v {
		return 0
	}
	return inv
}

func distance(p1, p2 []int) int {
	pf := 0
	m := len(p1)
	if len(p2) < m {
		m = len(p2)
	}
	for pf < m && p1[pf] == p2[pf] {
		pf++
	}
	return len(p1) + len(p2) - 2*pf
}

func solveCase(a []int) (int, int, int) {
	n := len(a)
	paths := make([][]int, n)
	for i := 0; i < n; i++ {
		path := []int{a[i]}
		for path[len(path)-1] != 0 {
			path = append(path, getInv(path[len(path)-1]))
		}
		for l, r := 0, len(path)-1; l < r; l, r = l+1, r-1 {
			path[l], path[r] = path[r], path[l]
		}
		paths[i] = path
	}
	from := 0
	best := -1
	far := -1
	for i := 0; i < n; i++ {
		if i == from {
			continue
		}
		d := distance(paths[from], paths[i])
		if d > best {
			best = d
			far = i
		}
	}
	from = far
	best = -1
	far2 := -1
	for i := 0; i < n; i++ {
		if i == from {
			continue
		}
		d := distance(paths[from], paths[i])
		if d > best {
			best = d
			far2 = i
		}
	}
	return from + 1, far2 + 1, best
}

func genTest(rng *rand.Rand) []int {
	n := rng.Intn(9) + 2
	used := make(map[int]bool)
	arr := make([]int, n)
	for i := 0; i < n; i++ {
		for {
			v := rng.Intn(1000)
			if !used[v] {
				used[v] = true
				arr[i] = v
				break
			}
		}
	}
	return arr
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	const tests = 100
	for t := 0; t < tests; t++ {
		arr := genTest(rng)
		u, v, d := solveCase(arr)
		var input bytes.Buffer
		input.WriteString(fmt.Sprintf("%d\n", len(arr)))
		for i, x := range arr {
			if i > 0 {
				input.WriteByte(' ')
			}
			input.WriteString(fmt.Sprintf("%d", x))
		}
		input.WriteByte('\n')
		out, err := run(bin, input.String())
		if err != nil {
			fmt.Printf("test %d runtime error: %v\n%s", t+1, err, out)
			os.Exit(1)
		}
		fields := strings.Fields(strings.TrimSpace(out))
		if len(fields) != 3 {
			fmt.Printf("test %d expected three numbers got %d\n", t+1, len(fields))
			os.Exit(1)
		}
		var gotU, gotV, gotD int
		if _, err := fmt.Sscan(fields[0], &gotU); err != nil {
			fmt.Printf("test %d bad output: %v\n", t+1, err)
			os.Exit(1)
		}
		if _, err := fmt.Sscan(fields[1], &gotV); err != nil {
			fmt.Printf("test %d bad output: %v\n", t+1, err)
			os.Exit(1)
		}
		if _, err := fmt.Sscan(fields[2], &gotD); err != nil {
			fmt.Printf("test %d bad output: %v\n", t+1, err)
			os.Exit(1)
		}
		if gotU != u || gotV != v || gotD != d {
			fmt.Printf("test %d failed expected:%d %d %d got:%d %d %d\n", t+1, u, v, d, gotU, gotV, gotD)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
