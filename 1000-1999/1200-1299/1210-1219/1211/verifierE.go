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

func runProgram(bin string, input []byte) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = bytes.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func expectedE(arr []int) string {
	const limit = 200000
	first := make([]int, limit+1)
	second := make([]int, limit+1)
	count := make([]int, limit+1)
	for i, v := range arr {
		count[v]++
		if count[v] == 1 {
			first[v] = i + 1
		} else if count[v] == 2 {
			second[v] = i + 1
		}
	}
	k0 := 0
	for v := 1; v <= limit; v++ {
		if count[v] == 2 {
			k0 = v
		} else {
			break
		}
	}
	type pair struct{ f, s, val int }
	pairs := make([]pair, 0, k0)
	for v := 1; v <= k0; v++ {
		pairs = append(pairs, pair{first[v], second[v], v})
	}
	sort.Slice(pairs, func(i, j int) bool { return pairs[i].f < pairs[j].f })

	curSecond := 0
	seen := 0
	maxValSeen := 0
	k := 0
	ok := true
	for _, p := range pairs {
		if p.s <= curSecond {
			ok = false
		}
		curSecond = p.s
		seen++
		if p.val > maxValSeen {
			maxValSeen = p.val
		}
		if ok && seen == maxValSeen {
			k = maxValSeen
		}
	}
	ans := make([]byte, len(arr))
	for i := range ans {
		ans[i] = 'B'
	}
	for v := 1; v <= k; v++ {
		ans[first[v]-1] = 'R'
		ans[second[v]-1] = 'G'
	}
	return string(ans)
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		return
	}
	rand.Seed(time.Now().UnixNano())
	bin := os.Args[1]

	for t := 0; t < 100; t++ {
		n := rand.Intn(30) + 1
		arr := make([]int, n)
		var input strings.Builder
		input.WriteString(fmt.Sprintf("%d\n", n))
		for i := 0; i < n; i++ {
			arr[i] = rand.Intn(20) + 1
			if i > 0 {
				input.WriteByte(' ')
			}
			input.WriteString(fmt.Sprintf("%d", arr[i]))
		}
		input.WriteByte('\n')
		expected := expectedE(append([]int(nil), arr...))
		out, err := runProgram(bin, []byte(input.String()))
		if err != nil || strings.TrimSpace(out) != expected {
			fmt.Println("Test", t+1, "failed")
			fmt.Println("Input:\n", input.String())
			fmt.Println("Expected:", expected)
			fmt.Println("Output:", out)
			if err != nil {
				fmt.Println("Error:", err)
			}
			return
		}
	}
	fmt.Println("All tests passed")
}
