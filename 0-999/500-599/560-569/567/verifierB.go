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

func run(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		return out.String() + stderr.String(), err
	}
	return out.String(), nil
}

type Event struct {
	op string
	id int
}

type Test struct {
	events []Event
	input  string
}

func genTest(rng *rand.Rand) Test {
	n := rng.Intn(90) + 10 // 10..99
	inside := make(map[int]bool)
	used := make(map[int]bool)
	nextID := 1
	events := make([]Event, 0, n)
	// initial random inside people
	initial := rng.Intn(5)
	for i := 0; i < initial; i++ {
		id := nextID
		nextID++
		inside[id] = true
	}
	for i := 0; i < n; i++ {
		if rng.Intn(2) == 0 { // add
			id := nextID
			if len(inside) > 0 && rng.Intn(2) == 0 {
				// maybe reuse an old id not inside
				for candidate := range used {
					if !inside[candidate] {
						id = candidate
						break
					}
				}
			}
			used[id] = true
			inside[id] = true
			events = append(events, Event{"+", id})
		} else { // remove
			var id int
			if len(inside) > 0 && rng.Intn(3) != 0 {
				// remove someone inside
				for k := range inside {
					id = k
					break
				}
				delete(inside, id)
			} else {
				id = nextID
				nextID++
				used[id] = true
				// event implies he was inside before log
			}
			events = append(events, Event{"-", id})
		}
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", len(events)))
	for _, e := range events {
		sb.WriteString(fmt.Sprintf("%s %d\n", e.op, e.id))
	}
	return Test{events: events, input: sb.String()}
}

func solve(t Test) string {
	inside := make(map[int]bool)
	cur := 0
	ans := 0
	for _, e := range t.events {
		if e.op == "+" {
			inside[e.id] = true
			cur++
			if cur > ans {
				ans = cur
			}
		} else {
			if inside[e.id] {
				delete(inside, e.id)
				cur--
			} else {
				ans++
			}
		}
	}
	return strconv.Itoa(ans)
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		t := genTest(rng)
		expected := solve(t)
		out, err := run(bin, t.input)
		if err != nil {
			fmt.Printf("test %d runtime error: %v\noutput:\n%s", i+1, err, out)
			os.Exit(1)
		}
		out = strings.TrimSpace(out)
		if out != expected {
			fmt.Printf("test %d failed\ninput:\n%s\nexpected:%s got:%s\n", i+1, t.input, expected, out)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
