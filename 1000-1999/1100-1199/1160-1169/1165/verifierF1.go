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

const numTestsF1 = 100

func canFinish(D int, k []int, sales [][]int, totalK int) bool {
	cheap := 0
	for i, ki := range k {
		if ki <= 0 {
			continue
		}
		ci := sort.SearchInts(sales[i], D+1)
		if ci > ki {
			ci = ki
		}
		cheap += ci
	}
	expensive := totalK - cheap
	cost := cheap + expensive*2
	return cost <= D
}

func solveF1(k []int, sales [][]int) int {
	totalK := 0
	for _, v := range k {
		totalK += v
	}
	lo, hi := 0, 3000
	for lo < hi {
		mid := (lo + hi) / 2
		if canFinish(mid, k, sales, totalK) {
			hi = mid
		} else {
			lo = mid + 1
		}
	}
	return lo
}

func run(binary, input string) (string, error) {
	cmd := exec.Command(binary)
	cmd.Stdin = strings.NewReader(input)
	var buf bytes.Buffer
	cmd.Stdout = &buf
	cmd.Stderr = &buf
	err := cmd.Run()
	return buf.String(), err
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierF1.go /path/to/binary")
		os.Exit(1)
	}
	binary := os.Args[1]
	rand.Seed(6)
	for t := 1; t <= numTestsF1; t++ {
		n := rand.Intn(3) + 1
		m := rand.Intn(5) + 1
		k := make([]int, n)
		for i := range k {
			k[i] = rand.Intn(3) + 1
		}
		sales := make([][]int, n)
		events := make([][2]int, m)
		for j := 0; j < m; j++ {
			d := rand.Intn(10) + 1
			tpe := rand.Intn(n) + 1
			events[j] = [2]int{d, tpe}
			sales[tpe-1] = append(sales[tpe-1], d)
		}
		for i := 0; i < n; i++ {
			sort.Ints(sales[i])
		}
		input := fmt.Sprintf("%d %d\n", n, m)
		for i, v := range k {
			if i > 0 {
				input += " "
			}
			input += fmt.Sprintf("%d", v)
		}
		input += "\n"
		for _, e := range events {
			input += fmt.Sprintf("%d %d\n", e[0], e[1])
		}
		expect := solveF1(k, sales)
		out, err := run(binary, input)
		if err != nil {
			fmt.Printf("test %d failed to run: %v\noutput:%s\n", t, err, out)
			os.Exit(1)
		}
		fields := strings.Fields(out)
		if len(fields) == 0 {
			fmt.Printf("test %d: no output\n", t)
			os.Exit(1)
		}
		var got int
		fmt.Sscanf(fields[0], "%d", &got)
		if got != expect {
			fmt.Printf("test %d failed\ninput:%sexpected:%d got:%d\n", t, input, expect, got)
			os.Exit(1)
		}
	}
	fmt.Println("OK")
}
