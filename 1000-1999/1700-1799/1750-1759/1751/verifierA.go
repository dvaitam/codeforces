package main

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

type Task struct {
	size     int
	data     int
	machines []int
	dataDeps []int
	taskDeps []int
}

type Case struct {
	tasks    []Task
	n, m     int
	power    []int
	speed    []int
	capacity []int
}

func divCeil(x, y int) int { return (x + y - 1) / y }

func genCase(rng *rand.Rand) Case {
	l := rng.Intn(5) + 3 // 3..7 tasks
	n := rng.Intn(3) + 1 // 1..3 machines
	m := rng.Intn(3) + 2 // 2..4 disks

	tasks := make([]Task, l+1)
	totalData := 0
	for i := 1; i <= l; i++ {
		size := rng.Intn(11) + 10 // 10..20
		data := rng.Intn(6)       // 0..5
		totalData += data
		k := rng.Intn(n) + 1
		perm := rng.Perm(n)
		machines := make([]int, k)
		for j := 0; j < k; j++ {
			machines[j] = perm[j] + 1
		}
		tasks[i] = Task{size: size, data: data, machines: machines}
	}

	power := make([]int, n+1)
	for i := 1; i <= n; i++ {
		power[i] = rng.Intn(5) + 1
	}

	speed := make([]int, m+1)
	capacity := make([]int, m+1)
	for k := 1; k <= m; k++ {
		speed[k] = rng.Intn(3) + 1
		capacity[k] = totalData + rng.Intn(10)
	}

	// dependencies
	for i := 1; i <= l; i++ {
		for j := i + 1; j <= l; j++ {
			r := rng.Intn(4)
			if r == 0 {
				tasks[j].dataDeps = append(tasks[j].dataDeps, i)
			} else if r == 1 {
				tasks[j].taskDeps = append(tasks[j].taskDeps, i)
			}
		}
	}

	return Case{tasks: tasks, n: n, m: m, power: power, speed: speed, capacity: capacity}
}

func buildInput(c Case) string {
	l := len(c.tasks) - 1
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", l)
	for i := 1; i <= l; i++ {
		t := c.tasks[i]
		fmt.Fprintf(&sb, "%d %d %d %d", i, t.size, t.data, len(t.machines))
		for _, m := range t.machines {
			fmt.Fprintf(&sb, " %d", m)
		}
		sb.WriteByte('\n')
	}
	fmt.Fprintf(&sb, "%d\n", c.n)
	for i := 1; i <= c.n; i++ {
		fmt.Fprintf(&sb, "%d %d\n", i, c.power[i])
	}
	fmt.Fprintf(&sb, "%d\n", c.m)
	for i := 1; i <= c.m; i++ {
		fmt.Fprintf(&sb, "%d %d %d\n", i, c.speed[i], c.capacity[i])
	}

	// collect deps
	var dataEdges [][2]int
	var taskEdges [][2]int
	for j := 1; j <= l; j++ {
		for _, v := range c.tasks[j].dataDeps {
			dataEdges = append(dataEdges, [2]int{v, j})
		}
		for _, v := range c.tasks[j].taskDeps {
			taskEdges = append(taskEdges, [2]int{v, j})
		}
	}
	fmt.Fprintf(&sb, "%d\n", len(dataEdges))
	for _, e := range dataEdges {
		fmt.Fprintf(&sb, "%d %d\n", e[0], e[1])
	}
	fmt.Fprintf(&sb, "%d\n", len(taskEdges))
	for _, e := range taskEdges {
		fmt.Fprintf(&sb, "%d %d\n", e[0], e[1])
	}
	return sb.String()
}

func run(bin string, input string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return "", fmt.Errorf("timeout")
		}
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return out.String(), nil
}

func topoOrder(tasks []Task) []int {
	l := len(tasks) - 1
	indeg := make([]int, l+1)
	adj := make([][]int, l+1)
	for i := 1; i <= l; i++ {
		for _, j := range tasks[i].dataDeps {
			adj[j] = append(adj[j], i)
			indeg[i]++
		}
		for _, j := range tasks[i].taskDeps {
			adj[j] = append(adj[j], i)
			indeg[i]++
		}
	}
	q := make([]int, 0, l)
	for i := 1; i <= l; i++ {
		if indeg[i] == 0 {
			q = append(q, i)
		}
	}
	order := make([]int, 0, l)
	for len(q) > 0 {
		v := q[0]
		q = q[1:]
		order = append(order, v)
		for _, u := range adj[v] {
			indeg[u]--
			if indeg[u] == 0 {
				q = append(q, u)
			}
		}
	}
	return order
}

func verify(c Case, output string) error {
	l := len(c.tasks) - 1
	reader := bufio.NewReader(strings.NewReader(output))
	start := make([]int, l+1)
	mach := make([]int, l+1)
	disk := make([]int, l+1)
	seen := make([]bool, l+1)

	for i := 0; i < l; i++ {
		var id, x, y, z int
		if _, err := fmt.Fscan(reader, &id, &x, &y, &z); err != nil {
			if err == io.EOF {
				return fmt.Errorf("missing output lines")
			}
			return fmt.Errorf("bad line %d: %v", i+1, err)
		}
		if id < 1 || id > l {
			return fmt.Errorf("invalid id %d", id)
		}
		if seen[id] {
			return fmt.Errorf("duplicate line for task %d", id)
		}
		seen[id] = true
		start[id] = x
		mach[id] = y
		disk[id] = z
	}
	if _, err := fmt.Fscan(reader, new(int)); err != io.EOF {
		return fmt.Errorf("extraneous output")
	}
	for i := 1; i <= l; i++ {
		if !seen[i] {
			return fmt.Errorf("missing schedule for task %d", i)
		}
	}

	// verify schedule
	order := topoOrder(c.tasks)
	cTime := make([]int, l+1)
	dTime := make([]int, l+1)
	diskUsed := make([]int, c.m+1)
	machineInt := make([][][2]int, c.n+1)

	for _, id := range order {
		t := c.tasks[id]
		y := mach[id]
		z := disk[id]
		if y < 1 || y > c.n {
			return fmt.Errorf("task %d invalid machine %d", id, y)
		}
		if z < 1 || z > c.m {
			return fmt.Errorf("task %d invalid disk %d", id, z)
		}
		allowed := false
		for _, v := range t.machines {
			if v == y {
				allowed = true
				break
			}
		}
		if !allowed {
			return fmt.Errorf("task %d not allowed on machine %d", id, y)
		}
		a := start[id]
		if a < 0 {
			return fmt.Errorf("task %d negative start", id)
		}
		for _, j := range t.dataDeps {
			if a < dTime[j] {
				return fmt.Errorf("task %d starts before data dep %d", id, j)
			}
		}
		for _, j := range t.taskDeps {
			if a < cTime[j] {
				return fmt.Errorf("task %d starts before task dep %d", id, j)
			}
		}
		readT := 0
		for _, j := range t.dataDeps {
			readT += divCeil(c.tasks[j].data, c.speed[disk[j]])
		}
		execT := divCeil(t.size, c.power[y])
		storeT := divCeil(t.data, c.speed[z])
		b := a + readT
		cFinish := b + execT
		d := cFinish + storeT

		for _, iv := range machineInt[y] {
			if a < iv[1] && iv[0] < d {
				return fmt.Errorf("machine %d overlap", y)
			}
		}
		machineInt[y] = append(machineInt[y], [2]int{a, d})

		diskUsed[z] += t.data
		cTime[id] = cFinish
		dTime[id] = d
	}
	for k := 1; k <= c.m; k++ {
		if diskUsed[k] > c.capacity[k] {
			return fmt.Errorf("disk %d capacity exceeded", k)
		}
	}
	return nil
}

func runCase(bin string, c Case) error {
	in := buildInput(c)
	out, err := run(bin, in)
	if err != nil {
		return err
	}
	return verify(c, out)
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 1; i <= 100; i++ {
		c := genCase(rng)
		if err := runCase(bin, c); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
