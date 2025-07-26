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

func solveA(nums []int) string {
	ones, twos := 0, 0
	for _, v := range nums {
		if v == 1 {
			ones++
		} else {
			twos++
		}
	}
	var res []int
	if ones == 0 {
		for i := 0; i < twos; i++ {
			res = append(res, 2)
		}
	} else if twos == 0 {
		for i := 0; i < ones; i++ {
			res = append(res, 1)
		}
	} else {
		res = append(res, 2)
		twos--
		res = append(res, 1)
		ones--
		for i := 0; i < twos; i++ {
			res = append(res, 2)
		}
		for i := 0; i < ones; i++ {
			res = append(res, 1)
		}
	}
	var sb strings.Builder
	for i, v := range res {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", v))
	}
	sb.WriteByte('\n')
	return sb.String()
}

func genCaseA(rng *rand.Rand) (string, string) {
	n := rng.Intn(20) + 1
	nums := make([]int, n)
	for i := range nums {
		if rng.Intn(2) == 0 {
			nums[i] = 1
		} else {
			nums[i] = 2
		}
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for i, v := range nums {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", v))
	}
	sb.WriteByte('\n')
	return sb.String(), solveA(nums)
}

func run(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		return
	}
	bin := os.Args[1]
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < 100; i++ {
		in, expect := genCaseA(rand.New(rand.NewSource(time.Now().UnixNano() + int64(i))))
		got, err := run(bin, in)
		if err != nil {
			fmt.Printf("test %d: runtime error: %v\ninput:\n%soutput:\n%s", i+1, err, in, got)
			return
		}
		if strings.TrimSpace(got) != strings.TrimSpace(expect) {
			fmt.Printf("test %d failed\ninput:\n%sexpected:\n%sbut got:\n%s", i+1, in, expect, got)
			return
		}
	}
	fmt.Println("All tests passed")
}
