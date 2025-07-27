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

func randString(n int) string {
	b := make([]byte, n)
	for i := range b {
		if rand.Intn(2) == 0 {
			b[i] = 'R'
		} else {
			b[i] = 'B'
		}
	}
	return string(b)
}

func genTest() []byte {
	n := rand.Intn(6) + 1
	m := rand.Intn(6) + 1
	q := rand.Intn(6) + 1
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d %d\n", n, m, q))
	sb.WriteString(randString(n) + "\n")
	sb.WriteString(randString(n) + "\n")
	sb.WriteString(randString(m) + "\n")
	sb.WriteString(randString(m) + "\n")
	dirs := []string{"L", "R", "U", "D"}
	for i := 0; i < q; i++ {
		d := dirs[rand.Intn(4)]
		if d == "L" || d == "R" {
			l := rand.Intn(n) + 1
			r := rand.Intn(n-l+1) + l
			sb.WriteString(fmt.Sprintf("%s %d %d\n", d, l, r))
		} else {
			l := rand.Intn(m) + 1
			r := rand.Intn(m-l+1) + l
			sb.WriteString(fmt.Sprintf("%s %d %d\n", d, l, r))
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
		fmt.Println("usage: go run verifierH2.go /path/to/binary")
		return
	}
	ref := "./refH2.bin"
	if err := exec.Command("go", "build", "-o", ref, "1368H2.go").Run(); err != nil {
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
