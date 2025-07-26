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

const numTestsF2 = 100

func can(D int, k []int, offers [][]int, totalK int) bool {
	cheap := 0
	for i := 0; i < len(k); i++ {
		if k[i] == 0 {
			continue
		}
		c := sort.SearchInts(offers[i], D+1)
		if c > k[i] {
			c = k[i]
		}
		cheap += c
	}
	need := 2*totalK - cheap
	return need <= D
}

func solveF2(k []int, offers [][]int) int {
	totalK := 0
	for _, v := range k {
		totalK += v
	}
	maxDay := 0
	for i := 0; i < len(offers); i++ {
		for _, d := range offers[i] {
			if d > maxDay {
				maxDay = d
			}
		}
	}
	left, right := 0, max(maxDay, 2*totalK)
	for left+1 < right {
		mid := (left + right) / 2
		if can(mid, k, offers, totalK) {
			right = mid
		} else {
			left = mid
		}
	}
	return right
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
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
		fmt.Println("usage: go run verifierF2.go /path/to/binary")
		os.Exit(1)
	}
	binary := os.Args[1]
	rand.Seed(7)
	for t := 1; t <= numTestsF2; t++ {
		n := rand.Intn(3) + 1
		m := rand.Intn(5) + 1
		k := make([]int, n)
		total := 0
		for i := range k {
			k[i] = rand.Intn(3) + 1
			total += k[i]
		}
		offers := make([][]int, n)
		events := make([][2]int, m)
		for j := 0; j < m; j++ {
			d := rand.Intn(10) + 1
			tpe := rand.Intn(n) + 1
			events[j] = [2]int{d, tpe}
			offers[tpe-1] = append(offers[tpe-1], d)
		}
		for i := 0; i < n; i++ {
			sort.Ints(offers[i])
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
		expect := solveF2(k, offers)
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
