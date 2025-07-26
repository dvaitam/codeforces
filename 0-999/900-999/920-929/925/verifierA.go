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

func run(bin string, input []byte) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = bytes.NewReader(input)
	out, err := cmd.CombinedOutput()
	return string(out), err
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func genTest() []byte {
	n := rand.Intn(9) + 2 // floors 2..10
	m := rand.Intn(6) + 5 // sections 5..10
	cl := rand.Intn(min(2, m-2)) + 1
	ce := rand.Intn(min(2, m-cl-1)) + 1

	positions := rand.Perm(m)[:cl+ce]
	stairs := make([]int, cl)
	elevators := make([]int, ce)
	for i := 0; i < cl; i++ {
		stairs[i] = positions[i] + 1
	}
	for i := 0; i < ce; i++ {
		elevators[i] = positions[cl+i] + 1
	}
	sort.Ints(stairs)
	sort.Ints(elevators)
	v := rand.Intn(n-1) + 1
	q := rand.Intn(3) + 1
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d %d %d %d\n", n, m, cl, ce, v))
	for i, x := range stairs {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprint(x))
	}
	sb.WriteByte('\n')
	for i, x := range elevators {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprint(x))
	}
	sb.WriteByte('\n')
	sb.WriteString(fmt.Sprintf("%d\n", q))
	allowed := make([]int, 0)
	for i := 1; i <= m; i++ {
		ok := true
		for _, x := range stairs {
			if i == x {
				ok = false
				break
			}
		}
		if ok {
			for _, x := range elevators {
				if i == x {
					ok = false
					break
				}
			}
		}
		if ok {
			allowed = append(allowed, i)
		}
	}
	if len(allowed) < 2 {
		allowed = append(allowed, 1, 2)
	}
	for i := 0; i < q; i++ {
		x1 := rand.Intn(n) + 1
		x2 := rand.Intn(n) + 1
		y1 := allowed[rand.Intn(len(allowed))]
		y2 := allowed[rand.Intn(len(allowed))]
		for x1 == x2 && y1 == y2 {
			x2 = rand.Intn(n) + 1
			y2 = allowed[rand.Intn(len(allowed))]
		}
		sb.WriteString(fmt.Sprintf("%d %d %d %d\n", x1, y1, x2, y2))
	}
	return []byte(sb.String())
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	candidate := os.Args[1]
	ref := "./refA.bin"
	if err := exec.Command("go", "build", "-o", ref, "925A.go").Run(); err != nil {
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
