package main

import (
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

func solve(n, k int64) (bool, [][2]int64) {
	if n > k*(k-1) {
		return false, nil
	}
	cnt := int64(0)
	add := int64(0)
	fi := int64(0)
	se := int64(0)
	res := make([][2]int64, 0, n)
	for cnt < n {
		if cnt%k == 0 {
			add++
		}
		cnt++
		fi = fi%k + 1
		se = (fi+add-1)%k + 1
		res = append(res, [2]int64{fi, se})
	}
	return true, res
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		return
	}
	bin := os.Args[1]
	rand.Seed(6)
	for t := 1; t <= 100; t++ {
		k := rand.Int63n(5) + 2
		n := rand.Int63n(k*(k-1) + 5)
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d %d\n", n, k))
		expectedOK, _ := solve(n, k)
		out, err := runBinary(bin, sb.String())
		if err != nil {
			fmt.Printf("test %d runtime error: %v\n", t, err)
			fmt.Println(out)
			return
		}
		out = strings.TrimSpace(out)
		lines := strings.Split(out, "\n")
		if strings.EqualFold(lines[0], "NO") {
			if expectedOK {
				fmt.Printf("test %d expected YES but got NO\ninput:\n%s\n", t, sb.String())
				return
			}
			continue
		}
		if !strings.EqualFold(lines[0], "YES") {
			fmt.Printf("test %d invalid output\n", t)
			return
		}
		if !expectedOK {
			fmt.Printf("test %d expected NO but got YES\ninput:\n%s\n", t, sb.String())
			return
		}
		if int64(len(lines)-1) != n {
			fmt.Printf("test %d wrong number of lines\n", t)
			return
		}
		for i := int64(0); i < n; i++ {
			var x, y int64
			if _, err := fmt.Sscan(lines[i+1], &x, &y); err != nil {
				fmt.Printf("test %d cannot parse line\n", t)
				return
			}
			if x < 1 || x > k || y < 1 || y > k || x == y {
				fmt.Printf("test %d invalid pair\n", t)
				return
			}
		}
	}
	fmt.Println("all tests passed")
}
