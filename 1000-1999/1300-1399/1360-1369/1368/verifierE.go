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

func runCmd(path string, input []byte) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(path, ".go") {
		cmd = exec.Command("go", "run", path)
	} else {
		cmd = exec.Command(path)
	}
	cmd.Stdin = bytes.NewReader(input)
	out, err := cmd.CombinedOutput()
	return string(out), err
}

func genTest() []byte {
	t := rand.Intn(3) + 1
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", t))
	for i := 0; i < t; i++ {
		n := rand.Intn(5) + 1
		maxEdges := n * (n - 1) / 2
		m := rand.Intn(maxEdges + 1)
		if m > 5 {
			m = 5
		}
		sb.WriteString(fmt.Sprintf("%d %d\n", n, m))
		used := make(map[[2]int]bool)
		for j := 0; j < m; j++ {
			var x, y int
			for {
				x = rand.Intn(n-1) + 1
				y = rand.Intn(n-x) + x + 1
				if !used[[2]int{x, y}] {
					used[[2]int{x, y}] = true
					break
				}
			}
			sb.WriteString(fmt.Sprintf("%d %d\n", x, y))
		}
	}
	return []byte(sb.String())
}

func main() {
	var cand string
	if len(os.Args) == 2 {
		cand = os.Args[1]
	} else if len(os.Args) == 3 && os.Args[1] == "--" {
		cand = os.Args[2]
	} else {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		return
	}
	ref := "./refE.bin"
	if err := exec.Command("go", "build", "-o", ref, "1368E.go").Run(); err != nil {
		fmt.Println("failed to build reference solution:", err)
		return
	}
	defer os.Remove(ref)
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < 100; i++ {
		input := genTest()
		want, err := runCmd(ref, input)
		if err != nil {
			fmt.Println("reference failed:", err)
			os.Exit(1)
		}
		got, err := runCmd(cand, input)
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
