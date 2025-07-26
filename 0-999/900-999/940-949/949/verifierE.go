package main

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"
	"time"
)

type test struct {
	n   int
	arr []int
}

func genTests() []test {
	rand.Seed(5)
	tests := make([]test, 100)
	for i := range tests {
		n := rand.Intn(10) + 1
		arr := make([]int, n)
		for j := 0; j < n; j++ {
			arr[j] = rand.Intn(21) - 10
		}
		tests[i] = test{n, arr}
	}
	return tests
}

func runBinary(path, input string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var cmd *exec.Cmd
	if strings.HasSuffix(path, ".go") {
		cmd = exec.CommandContext(ctx, "go", "run", path)
	} else {
		cmd = exec.CommandContext(ctx, path)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

// reference from 949E.go
var (
	pw  [20]int
	b   [20][]bool
	c   [20]int
	ans int
	v   []int
)

func dfs(k, ds int) {
	if k == 17 {
		if b[k][0] {
			c[17] = -1
			ds++
		} else if b[k][2] {
			c[17] = 1
			ds++
		} else {
			c[17] = 0
		}
		if ds < ans {
			ans = ds
			v = v[:0]
			for i := 0; i <= 17; i++ {
				if c[i] == 1 {
					v = append(v, 1<<i)
				} else if c[i] == -1 {
					v = append(v, -(1 << i))
				}
			}
		}
		c[17] = 0
		return
	}
	flag := false
	for i := 1; i <= pw[18-k]; i += 2 {
		if b[k][i] {
			flag = true
			break
		}
	}
	if !flag {
		for i := 0; i <= pw[17-k]; i++ {
			b[k+1][i] = b[k][i<<1]
		}
		c[k] = 0
		dfs(k+1, ds)
		return
	}
	for i, j := pw[17-k], 0; i <= pw[18-k]; i, j = i+2, j+1 {
		b[k+1][pw[16-k]+j] = b[k][i] || b[k][i+1]
	}
	for i, j := pw[17-k]-1, 1; i >= 0; i, j = i-2, j+1 {
		b[k+1][pw[16-k]-j] = b[k][i] || b[k][i-1]
	}
	c[k] = 1
	dfs(k+1, ds+1)
	for i, j := pw[17-k], 0; i >= 0; i, j = i-2, j+1 {
		b[k+1][pw[16-k]-j] = b[k][i] || b[k][i-1]
	}
	for i, j := pw[17-k]+1, 1; i <= pw[18-k]; i, j = i+2, j+1 {
		b[k+1][pw[16-k]+j] = b[k][i] || b[k][i+1]
	}
	c[k] = -1
	dfs(k+1, ds+1)
}

func solveRef(t test) (int, []int) {
	pw[0] = 1
	for i := 1; i < 20; i++ {
		pw[i] = pw[i-1] << 1
	}
	maxSz := pw[18] + 1
	for i := range b {
		b[i] = make([]bool, maxSz)
	}
	for _, val := range t.arr {
		b[0][val+pw[17]] = true
	}
	ans = 30
	dfs(0, 0)
	sort.Ints(v)
	return ans, append([]int(nil), v...)
}

func buildInput(t test) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", t.n))
	for i, v := range t.arr {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(v))
	}
	sb.WriteByte('\n')
	return sb.String()
}

func parseOutput(out string) (int, []int, bool) {
	scanner := bufio.NewScanner(strings.NewReader(strings.TrimSpace(out)))
	scanner.Split(bufio.ScanWords)
	if !scanner.Scan() {
		return 0, nil, false
	}
	k, err := strconv.Atoi(scanner.Text())
	if err != nil || k < 0 {
		return 0, nil, false
	}
	vals := make([]int, 0, k)
	for scanner.Scan() {
		v, err := strconv.Atoi(scanner.Text())
		if err != nil {
			return 0, nil, false
		}
		vals = append(vals, v)
	}
	if len(vals) != k {
		return 0, nil, false
	}
	sort.Ints(vals)
	return k, vals, true
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	candidate := os.Args[1]
	tests := genTests()
	for i, t := range tests {
		input := buildInput(t)
		out, err := runBinary(candidate, input)
		if err != nil {
			fmt.Printf("test %d: runtime error: %v\n", i+1, err)
			os.Exit(1)
		}
		gotK, gotArr, ok := parseOutput(out)
		if !ok {
			fmt.Printf("test %d: bad output\n", i+1)
			os.Exit(1)
		}
		expK, expArr := solveRef(t)
		if gotK != expK || len(gotArr) != len(expArr) {
			fmt.Printf("test %d failed\ninput:\n%s\nexpected:%v %v\noutput:%d %v\n", i+1, input, expK, expArr, gotK, gotArr)
			os.Exit(1)
		}
		for j := range gotArr {
			if gotArr[j] != expArr[j] {
				fmt.Printf("test %d failed\ninput:\n%s\nexpected:%v %v\noutput:%d %v\n", i+1, input, expK, expArr, gotK, gotArr)
				os.Exit(1)
			}
		}
	}
	fmt.Println("All tests passed")
}
