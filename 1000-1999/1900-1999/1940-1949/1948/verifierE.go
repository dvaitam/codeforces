package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
)

func generateCase() (int, int) {
	n := rand.Intn(20) + 2
	k := rand.Intn(2*n) + 1
	return n, k
}

func checkSolution(n, k int, a []int, q int, c []int) bool {
	if len(a) != n || len(c) != n {
		return false
	}
	seen := make(map[int]bool)
	for _, v := range a {
		if v < 1 || v > n || seen[v] {
			return false
		}
		seen[v] = true
	}
	expectedQ := n
	if k > 1 {
		expectedQ = (n + k - 1) / k
	}
	if q != expectedQ {
		return false
	}
	for _, x := range c {
		if x < 1 || x > q {
			return false
		}
	}
	// check cliques
	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			if c[i] == c[j] {
				if abs(i-j)+abs(a[i]-a[j]) > k {
					return false
				}
			}
		}
	}
	return true
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		return
	}
	bin := os.Args[1]

	rand.Seed(42)
	const t = 100
	ns := make([]int, t)
	ks := make([]int, t)
	for i := 0; i < t; i++ {
		ns[i], ks[i] = generateCase()
	}

	var input bytes.Buffer
	fmt.Fprintf(&input, "%d\n", t)
	for i := 0; i < t; i++ {
		fmt.Fprintf(&input, "%d %d\n", ns[i], ks[i])
	}

	cmd := exec.Command(bin)
	cmd.Stdin = bytes.NewReader(input.Bytes())
	var out bytes.Buffer
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		fmt.Println("failed to run binary:", err)
		os.Exit(1)
	}

	scanner := bufio.NewScanner(bytes.NewReader(out.Bytes()))
	for idx := 0; idx < t; idx++ {
		n := ns[idx]
		k := ks[idx]
		// read first line: permutation
		if !scanner.Scan() {
			fmt.Printf("missing output for case %d\n", idx+1)
			os.Exit(1)
		}
		line1 := scanner.Text()
		a := make([]int, 0, n)
		r := bufio.NewScanner(bytes.NewReader([]byte(line1)))
		r.Split(bufio.ScanWords)
		for r.Scan() {
			var v int
			fmt.Sscan(r.Text(), &v)
			a = append(a, v)
		}
		if len(a) != n {
			fmt.Printf("case %d: expected %d numbers in first line\n", idx+1, n)
			os.Exit(1)
		}
		// second line q
		if !scanner.Scan() {
			fmt.Printf("case %d: missing q\n", idx+1)
			os.Exit(1)
		}
		var q int
		fmt.Sscan(scanner.Text(), &q)
		// third line clique ids
		if !scanner.Scan() {
			fmt.Printf("case %d: missing third line\n", idx+1)
			os.Exit(1)
		}
		line3 := scanner.Text()
		cids := make([]int, 0, n)
		r2 := bufio.NewScanner(bytes.NewReader([]byte(line3)))
		r2.Split(bufio.ScanWords)
		for r2.Scan() {
			var v int
			fmt.Sscan(r2.Text(), &v)
			cids = append(cids, v)
		}
		if len(cids) != n {
			fmt.Printf("case %d: expected %d clique ids\n", idx+1, n)
			os.Exit(1)
		}
		if !checkSolution(n, k, a, q, cids) {
			fmt.Printf("case %d: invalid solution\n", idx+1)
			os.Exit(1)
		}
	}
	if scanner.Scan() {
		fmt.Println("warning: extra output detected")
	}
	fmt.Println("All tests passed!")
}
