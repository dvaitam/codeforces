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

const (
	MAX  = 200200
	MEXN = 1000
)

func solveH(N, F int, rows [][3]int) string {
	isp := make([]bool, MAX)
	pr := make([]int, MAX)
	cnt := make([]int, MAX)
	for i := 2; i < MAX; i++ {
		pr[i] = i
	}
	var moves []int
	for i := 2; i*i < MAX; i++ {
		if !isp[i] {
			for j := i * i; j < MAX; j += i {
				isp[j] = true
			}
			moves = append(moves, i*i)
		}
	}
	for i := 2; i < MAX; i++ {
		if !isp[i] {
			moves = append(moves, i)
			for j := i; j < MAX; j += i {
				pr[j] /= i
				cnt[j]++
			}
		}
	}
	for i := 2; i < MAX; i++ {
		if pr[i] == 1 && cnt[i] == 2 && i != F {
			moves = append(moves, i)
		}
	}
	sort.Ints(moves)
	chk := make([]bool, MAX)
	for _, x := range moves {
		if x < MAX {
			chk[x] = true
		}
	}
	mex := make([]int, MAX)
	in := make([][]int, MEXN)
	in[0] = append(in[0], 0)
	for i := 1; i < MAX; i++ {
		for {
			g := false
			mi := mex[i]
			for _, p := range in[mi] {
				if i-p >= 0 && chk[i-p] {
					g = true
					break
				}
			}
			if !g {
				break
			}
			mex[i]++
		}
		mi := mex[i]
		if mi < MEXN {
			in[mi] = append(in[mi], i)
		}
	}
	xor := 0
	for _, row := range rows {
		a, b, c := row[0], row[1], row[2]
		xor ^= mex[b-a-1]
		xor ^= mex[c-b-1]
	}
	if xor != 0 {
		return "Alice\nBob"
	}
	return "Bob\nAlice"
}

func run(bin string, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierH.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rand.Seed(time.Now().UnixNano())
	for t := 0; t < 100; t++ {
		N := rand.Intn(5) + 1
		F := rand.Intn(20) + 2
		rows := make([][3]int, N)
		var input strings.Builder
		input.WriteString(fmt.Sprintf("%d %d\n", N, F))
		for i := 0; i < N; i++ {
			a := rand.Intn(10)
			b := a + rand.Intn(10) + 1
			c := b + rand.Intn(10) + 1
			rows[i] = [3]int{a, b, c}
			input.WriteString(fmt.Sprintf("%d %d %d\n", a, b, c))
		}
		expect := solveH(N, F, rows)
		got, err := run(bin, input.String())
		if err != nil {
			fmt.Printf("runtime error on test %d: %v\n", t+1, err)
			fmt.Println("input:\n", input.String())
			os.Exit(1)
		}
		if strings.TrimSpace(got) != strings.TrimSpace(expect) {
			fmt.Printf("wrong answer on test %d\n", t+1)
			fmt.Println("input:\n", input.String())
			fmt.Printf("expected:\n%s\n got:\n%s\n", expect, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
