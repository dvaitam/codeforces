package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

func computeLoneliness(a []int) int {
	n := len(a)
	maxGap := 0
	for bit := 0; bit < 20; bit++ {
		last := 0
		appear := false
		for i := 0; i < n; i++ {
			if (a[i]>>bit)&1 == 1 {
				appear = true
				gap := i + 1 - last - 1
				if gap > maxGap {
					maxGap = gap
				}
				last = i + 1
			}
		}
		if appear {
			gap := n + 1 - last - 1
			if gap > maxGap {
				maxGap = gap
			}
		}
	}
	return maxGap + 1
}

type Case struct {
	a   []int
	ans int
}

func genCases(n int) []Case {
	rand.Seed(time.Now().UnixNano())
	cs := make([]Case, n)
	for i := 0; i < n; i++ {
		size := rand.Intn(15) + 1
		arr := make([]int, size)
		for j := 0; j < size; j++ {
			arr[j] = rand.Intn(1 << 5)
		}
		cs[i] = Case{arr, computeLoneliness(arr)}
	}
	return cs
}

func buildInput(cs []Case) string {
	var sb strings.Builder
	fmt.Fprintln(&sb, len(cs))
	for _, c := range cs {
		fmt.Fprintln(&sb, len(c.a))
		for j, v := range c.a {
			if j > 0 {
				sb.WriteByte(' ')
			}
			fmt.Fprint(&sb, v)
		}
		sb.WriteByte('\n')
	}
	return sb.String()
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	cs := genCases(100)
	input := buildInput(cs)
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		fmt.Println("runtime error:", err)
		os.Exit(1)
	}
	outputs := strings.Fields(strings.TrimSpace(out.String()))
	if len(outputs) != len(cs) {
		fmt.Printf("expected %d outputs got %d\n", len(cs), len(outputs))
		os.Exit(1)
	}
	for i, res := range outputs {
		v, err := strconv.Atoi(res)
		if err != nil || v != cs[i].ans {
			fmt.Printf("mismatch on case %d: expected %d got %s\n", i+1, cs[i].ans, res)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
