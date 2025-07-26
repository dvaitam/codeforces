package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

func runBinary(bin string, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	out, err := cmd.CombinedOutput()
	return string(out), err
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		return
	}
	bin := os.Args[1]
	rand.Seed(1)
	for t := 1; t <= 100; t++ {
		q := rand.Intn(5) + 1
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d\n", q))
		exp := make([]string, q)
		for i := 0; i < q; i++ {
			n := rand.Int63n(1000) + 1
			a := rand.Int63n(1000) + 1
			b := rand.Int63n(1000) + 1
			sb.WriteString(fmt.Sprintf("%d %d %d\n", n, a, b))
			cost1 := n * a
			cost2 := (n/2)*b + (n%2)*a
			if cost1 < cost2 {
				exp[i] = fmt.Sprintf("%d", cost1)
			} else {
				exp[i] = fmt.Sprintf("%d", cost2)
			}
		}
		out, err := runBinary(bin, sb.String())
		if err != nil {
			fmt.Printf("test %d runtime error: %v\n", t, err)
			fmt.Println(string(out))
			return
		}
		scanner := bufio.NewScanner(strings.NewReader(out))
		got := []string{}
		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())
			if line != "" {
				got = append(got, line)
			}
		}
		if len(got) != q {
			fmt.Printf("test %d wrong number of lines: got %d expected %d\ninput:\n%soutput:\n%s", t, len(got), q, sb.String(), out)
			return
		}
		for i := 0; i < q; i++ {
			if got[i] != exp[i] {
				fmt.Printf("test %d case %d failed\ninput:\n%sexpected %s got %s\n", t, i+1, sb.String(), exp[i], got[i])
				return
			}
		}
	}
	fmt.Println("all tests passed")
}
