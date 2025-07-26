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

func runProgram(bin string, input []byte) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = bytes.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func expectedD1(s string) string {
	n := len(s)
	id, hmin, h, count := 0, 0, 0, 0
	for i := 0; i < n; i++ {
		if h < hmin {
			id = i
			hmin = h
			count = 0
		}
		if h == hmin {
			count++
		}
		if s[i] == '(' {
			h++
		} else {
			h--
		}
	}
	if h != 0 {
		return "0\n1 1"
	}
	s2 := s[id:] + s[:id]
	best := count
	curr1, curr2 := 0, 0
	a, b := 0, 0
	a1, a2 := 0, 0
	h = 0
	for i := 0; i < n; i++ {
		if s2[i] == '(' {
			h++
		} else {
			h--
		}
		switch {
		case h == 0:
			if curr1 > best {
				best = curr1
				a = a1
				b = i
			}
			curr1 = 0
			a1 = i + 1
		case h == 1:
			curr1++
			if curr2+count > best {
				best = curr2 + count
				a = a2
				b = i
			}
			curr2 = 0
			a2 = i + 1
		case h == 2:
			curr2++
		}
	}
	start := (a + id) % n
	end := (b + id) % n
	return fmt.Sprintf("%d\n%d %d", best, start+1, end+1)
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierD1.go /path/to/binary")
		return
	}
	rand.Seed(time.Now().UnixNano())
	bin := os.Args[1]
	if bin == "--" && len(os.Args) >= 3 {
		bin = os.Args[2]
	}

	for t := 0; t < 100; t++ {
		n := rand.Intn(20) + 1
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d\n", n))
		var s strings.Builder
		for i := 0; i < n; i++ {
			if rand.Intn(2) == 0 {
				s.WriteByte('(')
			} else {
				s.WriteByte(')')
			}
		}
		str := s.String()
		sb.WriteString(str)
		sb.WriteByte('\n')
		expected := expectedD1(str)
		out, err := runProgram(bin, []byte(sb.String()))
		if err != nil || strings.TrimSpace(out) != expected {
			fmt.Printf("Test %d failed\n", t+1)
			fmt.Println("Input:\n", sb.String())
			fmt.Println("Expected:\n", expected)
			fmt.Println("Output:\n", out)
			if err != nil {
				fmt.Println("Error:", err)
			}
			return
		}
	}
	fmt.Println("All tests passed")
}
