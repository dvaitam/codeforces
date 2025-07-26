package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

type state struct {
	vols []int
}

func encode(v []int) string {
	b := make([]byte, len(v)*4)
	p := 0
	for _, x := range v {
		b[p] = byte(x)
		b[p+1] = byte(x >> 8)
		b[p+2] = byte(x >> 16)
		b[p+3] = byte(x >> 24)
		p += 4
	}
	return string(b)
}

func possible(n, k, v int, a []int) bool {
	queue := [][]int{append([]int(nil), a...)}
	visited := map[string]bool{encode(a): true}
	for len(queue) > 0 {
		cur := queue[0]
		queue = queue[1:]
		for _, x := range cur {
			if x == v {
				return true
			}
		}
		for i := 0; i < n; i++ {
			for j := 0; j < n; j++ {
				if i == j || cur[i] == 0 {
					continue
				}
				take := cur[i]
				if take > k {
					take = k
				}
				next := append([]int(nil), cur...)
				next[i] -= take
				next[j] += take
				key := encode(next)
				if !visited[key] {
					visited[key] = true
					queue = append(queue, next)
				}
			}
		}
	}
	return false
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierD.go /path/to/binary")
		return
	}
	bin := os.Args[1]
	rand.Seed(4)
	const tests = 100
	for t := 0; t < tests; t++ {
		n := rand.Intn(3) + 2
		k := rand.Intn(5) + 1
		a := make([]int, n)
		for i := range a {
			a[i] = rand.Intn(6)
		}
		v := rand.Intn(10)
		expected := "NO"
		if possible(n, k, v, a) {
			expected = "YES"
		}
		var input bytes.Buffer
		fmt.Fprintf(&input, "%d %d %d\n", n, k, v)
		for i, x := range a {
			if i > 0 {
				fmt.Fprint(&input, " ")
			}
			fmt.Fprint(&input, x)
		}
		fmt.Fprintln(&input)
		cmd := exec.Command(bin)
		cmd.Stdin = bytes.NewReader(input.Bytes())
		out, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Println("Failed to run binary:", err)
			os.Exit(1)
		}
		answer := strings.Fields(string(out))
		if len(answer) == 0 {
			fmt.Printf("Test %d produced no output\n", t+1)
			os.Exit(1)
		}
		if strings.ToUpper(answer[0]) != expected {
			fmt.Printf("Test %d failed: expected %s got %s\n", t+1, expected, answer[0])
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
