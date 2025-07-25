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

func runCandidate(bin, input string) (string, error) {
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

type pair struct{ tp, base int64 }

func solve(n, k int, b, c int64, t []int64) int64 {
	per5 := b
	if per5 > 5*c {
		per5 = 5 * c
	}
	const inf int64 = 1<<63 - 1
	ans := inf
	for r := int64(0); r < 5; r++ {
		arr := make([]pair, n)
		for i := 0; i < n; i++ {
			rem := ((t[i] % 5) + 5) % 5
			diff := (r - rem + 5) % 5
			tp := t[i] + diff
			w := tp / 5
			base := diff*c - w*per5
			arr[i] = pair{tp, base}
		}
		sort.Slice(arr, func(i, j int) bool { return arr[i].tp < arr[j].tp })
		heap := make([]int64, 0, k)
		var sum int64
		push := func(x int64) {
			heap = append(heap, x)
			i := len(heap) - 1
			for i > 0 {
				p := (i - 1) / 2
				if heap[p] >= heap[i] {
					break
				}
				heap[p], heap[i] = heap[i], heap[p]
				i = p
			}
			sum += x
			if len(heap) > k {
				sum -= heap[0]
				heap[0] = heap[len(heap)-1]
				heap = heap[:len(heap)-1]
				i = 0
				for {
					l := 2*i + 1
					if l >= len(heap) {
						break
					}
					r := l + 1
					if r < len(heap) && heap[r] > heap[l] {
						l = r
					}
					if heap[i] >= heap[l] {
						break
					}
					heap[i], heap[l] = heap[l], heap[i]
					i = l
				}
			}
		}
		i := 0
		for i < n {
			curW := arr[i].tp / 5
			for i < n && arr[i].tp/5 == curW {
				push(arr[i].base)
				i++
			}
			if len(heap) == k {
				cost := sum + int64(k)*curW*per5
				if cost < ans {
					ans = cost
				}
			}
		}
	}
	return ans
}

func generateCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(6) + 1
	k := rng.Intn(n) + 1
	b := int64(rng.Intn(5) + 1)
	c := int64(rng.Intn(5) + 1)
	t := make([]int64, n)
	for i := 0; i < n; i++ {
		t[i] = int64(rng.Intn(10))
	}
	input := fmt.Sprintf("%d %d %d %d\n", n, k, b, c)
	for i := 0; i < n; i++ {
		if i > 0 {
			input += " "
		}
		input += fmt.Sprintf("%d", t[i])
	}
	input += "\n"
	exp := fmt.Sprintf("%d", solve(n, k, b, c, t))
	return input, exp
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := generateCase(rng)
		out, err := runCandidate(bin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:%s", i+1, err, in)
			os.Exit(1)
		}
		if strings.TrimSpace(out) != strings.TrimSpace(exp) {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:%s", i+1, exp, out, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
