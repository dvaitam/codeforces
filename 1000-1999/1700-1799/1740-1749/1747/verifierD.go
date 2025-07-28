package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
)

// generateTests creates a single test for problem D with q=100 queries
func generateTests() (string, []int) {
	rand.Seed(1)
	n := 10
	q := 100
	arr := make([]int, n)
	for i := range arr {
		arr[i] = rand.Intn(8)
	}
	var buf bytes.Buffer
	fmt.Fprintln(&buf, n, q)
	for i, v := range arr {
		if i > 0 {
			fmt.Fprint(&buf, " ")
		}
		fmt.Fprint(&buf, v)
	}
	fmt.Fprintln(&buf)
	lvals := make([]int, q)
	rvals := make([]int, q)
	for i := 0; i < q; i++ {
		l := rand.Intn(n) + 1
		r := rand.Intn(n-l+1) + l
		lvals[i] = l
		rvals[i] = r
		fmt.Fprintln(&buf, l, r)
	}
	// compute expected answers using solution algorithm
	px := make([]int, n+1)
	nz := make([]int, n+1)
	even := map[int][]int{}
	odd := map[int][]int{}
	for i := 0; i <= n; i++ {
		if i > 0 {
			px[i] = px[i-1] ^ arr[i-1]
			if arr[i-1] != 0 {
				nz[i] = nz[i-1] + 1
			} else {
				nz[i] = nz[i-1]
			}
		}
		if i%2 == 0 {
			even[px[i]] = append(even[px[i]], i)
		} else {
			odd[px[i]] = append(odd[px[i]], i)
		}
	}
	answers := make([]int, q)
	for idx := 0; idx < q; idx++ {
		l := lvals[idx]
		r := rvals[idx]
		xor := px[r] ^ px[l-1]
		if xor != 0 {
			answers[idx] = -1
			continue
		}
		if nz[r]-nz[l-1] == 0 {
			answers[idx] = 0
			continue
		}
		if (r-l+1)%2 == 1 {
			answers[idx] = 1
			continue
		}
		if arr[l-1] == 0 || arr[r-1] == 0 {
			answers[idx] = 1
			continue
		}
		par := (l - 1) % 2
		var arrP []int
		if par == 0 {
			arrP = odd[px[l-1]]
		} else {
			arrP = even[px[l-1]]
		}
		i := sort.SearchInts(arrP, l)
		if i < len(arrP) && arrP[i] <= r {
			answers[idx] = 2
			continue
		}
		par = r % 2
		if par == 0 {
			arrP = odd[px[r]]
		} else {
			arrP = even[px[r]]
		}
		i = sort.SearchInts(arrP, l-1)
		if i < len(arrP) && arrP[i] < r {
			answers[idx] = 2
			continue
		}
		answers[idx] = -1
	}
	return buf.String(), answers
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	input, answers := generateTests()
	cmd := exec.Command(os.Args[1])
	cmd.Stdin = bytes.NewBufferString(input)
	out, err := cmd.Output()
	if err != nil {
		fmt.Println("failed to run binary:", err)
		os.Exit(1)
	}
	scanner := bufio.NewScanner(bytes.NewReader(out))
	for idx, exp := range answers {
		if !scanner.Scan() {
			fmt.Printf("missing output for query %d\n", idx+1)
			os.Exit(1)
		}
		var got int
		fmt.Sscan(scanner.Text(), &got)
		if got != exp {
			fmt.Printf("query %d: expected %d, got %d\n", idx+1, exp, got)
			os.Exit(1)
		}
	}
	if scanner.Scan() {
		fmt.Println("extra output detected")
		os.Exit(1)
	}
	fmt.Println("All tests passed")
}
