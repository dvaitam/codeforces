package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

type event struct {
	t  int
	op byte
}

func runBinary(path, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(path, ".go") {
		cmd = exec.Command("go", "run", path)
	} else {
		cmd = exec.Command(path)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func hasBad(pref []int, l, r int) bool {
	if l == 0 {
		return pref[r] > 0
	}
	return pref[r]-pref[l-1] > 0
}

func solveB(n, m int, ops []event, ids []int) string {
	events := make([][]event, n+1)
	first := make([]byte, n+1)
	opByte := make([]byte, m+1)
	for i := 1; i <= m; i++ {
		b := ops[i-1].op
		id := ids[i-1]
		opByte[i] = b
		if first[id] == 0 {
			first[id] = b
		}
		events[id] = append(events[id], event{t: i, op: b})
	}
	initial := make([]bool, n+1)
	initCnt := 0
	for id := 1; id <= n; id++ {
		if first[id] == '-' {
			initial[id] = true
			initCnt++
		}
	}
	cur := make([]int, m+1)
	bad := make([]int, m+1)
	pref := make([]int, m+1)
	cur[0] = initCnt
	if cur[0] > 0 {
		bad[0] = 1
	}
	pref[0] = bad[0]
	for i := 1; i <= m; i++ {
		if opByte[i] == '+' {
			cur[i] = cur[i-1] + 1
		} else {
			cur[i] = cur[i-1] - 1
		}
		if cur[i] > 0 {
			bad[i] = 1
		}
		pref[i] = pref[i-1] + bad[i]
	}
	res := make([]int, 0)
	for id := 1; id <= n; id++ {
		userEvents := events[id]
		last := 0
		online := initial[id]
		possible := true
		for _, e := range userEvents {
			if !online {
				if last <= e.t-1 && hasBad(pref, last, e.t-1) {
					possible = false
					break
				}
			}
			if e.op == '+' {
				online = true
			} else {
				online = false
			}
			last = e.t
		}
		if possible && !online {
			if last <= m && hasBad(pref, last, m) {
				possible = false
			}
		}
		if possible {
			res = append(res, id)
		}
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d", len(res)))
	if len(res) > 0 {
		sb.WriteByte('\n')
		for i, id := range res {
			if i > 0 {
				sb.WriteByte(' ')
			}
			fmt.Fprintf(&sb, "%d", id)
		}
	}
	return sb.String()
}

func generateCaseB(rng *rand.Rand) string {
	n := rng.Intn(5) + 1
	m := rng.Intn(10) + 1
	states := make([]bool, n+1)
	for i := 1; i <= n; i++ {
		states[i] = rng.Intn(2) == 1
	}
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d\n", n, m)
	for i := 0; i < m; i++ {
		if rng.Intn(2) == 0 {
			off := []int{}
			for id := 1; id <= n; id++ {
				if !states[id] {
					off = append(off, id)
				}
			}
			if len(off) == 0 {
				on := []int{}
				for id := 1; id <= n; id++ {
					if states[id] {
						on = append(on, id)
					}
				}
				id := on[rng.Intn(len(on))]
				states[id] = false
				fmt.Fprintf(&sb, "- %d\n", id)
				continue
			}
			id := off[rng.Intn(len(off))]
			states[id] = true
			fmt.Fprintf(&sb, "+ %d\n", id)
		} else {
			on := []int{}
			for id := 1; id <= n; id++ {
				if states[id] {
					on = append(on, id)
				}
			}
			if len(on) == 0 {
				off := []int{}
				for id := 1; id <= n; id++ {
					if !states[id] {
						off = append(off, id)
					}
				}
				id := off[rng.Intn(len(off))]
				states[id] = true
				fmt.Fprintf(&sb, "+ %d\n", id)
				continue
			}
			id := on[rng.Intn(len(on))]
			states[id] = false
			fmt.Fprintf(&sb, "- %d\n", id)
		}
	}
	return sb.String()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		tc := generateCaseB(rng)
		scanner := bufio.NewScanner(strings.NewReader(tc))
		scanner.Split(bufio.ScanWords)
		var nums []string
		for scanner.Scan() {
			nums = append(nums, scanner.Text())
		}
		n, _ := strconv.Atoi(nums[0])
		m, _ := strconv.Atoi(nums[1])
		ops := make([]event, m)
		ids := make([]int, m)
		idx := 2
		for j := 0; j < m; j++ {
			op := nums[idx]
			id, _ := strconv.Atoi(nums[idx+1])
			ops[j] = event{t: j + 1, op: op[0]}
			ids[j] = id
			idx += 2
		}
		expect := solveB(n, m, ops, ids)
		got, err := runBinary(bin, tc)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, tc)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != strings.TrimSpace(expect) {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, expect, got, tc)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
