package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

func run(bin string, input []byte) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = bytes.NewReader(input)
	out, err := cmd.CombinedOutput()
	return string(out), err
}

func genTest() []byte {
	n := rand.Intn(5) + 2 // 2..6 employees
	m := rand.Intn(5) + 2 // 2..6 events
	parents := make([]int, n+1)
	for i := 2; i <= n; i++ {
		parents[i] = rand.Intn(i-1) + 1
	}
	subCnt := make([]int, n+1)
	for i := 2; i <= n; i++ {
		p := parents[i]
		for p != 0 {
			subCnt[p]++
			p = parents[p]
		}
	}
	t := make([]int, n+1)
	for i := 1; i <= n; i++ {
		if subCnt[i] == 0 {
			t[i] = 0
		} else {
			t[i] = rand.Intn(subCnt[i] + 1)
		}
	}
	on := make([]bool, n+1)
	events := make([]int, 0, m)
	for len(events) < m {
		leaveProb := 0
		for i := 1; i <= n; i++ {
			if !on[i] {
				leaveProb++
			}
		}
		takeLeave := rand.Intn(2)
		if takeLeave == 0 && leaveProb > 0 || len(events) == 0 {
			// leave
			var candidates []int
			for i := 1; i <= n; i++ {
				if !on[i] {
					candidates = append(candidates, i)
				}
			}
			x := candidates[rand.Intn(len(candidates))]
			on[x] = true
			events = append(events, x)
		} else {
			var candidates []int
			for i := 1; i <= n; i++ {
				if on[i] {
					candidates = append(candidates, i)
				}
			}
			if len(candidates) == 0 {
				continue
			}
			x := candidates[rand.Intn(len(candidates))]
			on[x] = false
			events = append(events, -x)
		}
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", n, m))
	for i := 2; i <= n; i++ {
		if i > 2 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprint(parents[i]))
	}
	if n > 1 {
		sb.WriteByte('\n')
	} else {
		sb.WriteString("\n")
	}
	for i := 1; i <= n; i++ {
		if i > 1 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprint(t[i]))
	}
	sb.WriteByte('\n')
	for i, e := range events {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprint(e))
	}
	sb.WriteByte('\n')
	return []byte(sb.String())
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	candidate := os.Args[1]
	ref := "./refE.bin"
	if err := exec.Command("go", "build", "-o", ref, "925E.go").Run(); err != nil {
		fmt.Println("failed to build reference solution:", err)
		os.Exit(1)
	}
	defer os.Remove(ref)
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < 100; i++ {
		input := genTest()
		want, err := run(ref, input)
		if err != nil {
			fmt.Println("reference failed:", err)
			os.Exit(1)
		}
		got, err := run(candidate, input)
		if err != nil {
			fmt.Printf("candidate runtime error on test %d: %v\n", i+1, err)
			fmt.Println("input:\n", string(input))
			os.Exit(1)
		}
		if strings.TrimSpace(want) != strings.TrimSpace(got) {
			fmt.Printf("wrong answer on test %d\n", i+1)
			fmt.Println("input:\n", string(input))
			fmt.Println("expected:\n", want)
			fmt.Println("got:\n", got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed.")
}
