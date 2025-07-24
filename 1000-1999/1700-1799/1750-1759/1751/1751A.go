package main

import (
	"bufio"
	"fmt"
	"os"
)

type Task struct {
	id       int
	size     int
	data     int
	machines []int
	dataDeps []int
	taskDeps []int
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var l int
	if _, err := fmt.Fscan(reader, &l); err != nil {
		return
	}
	tasks := make([]Task, l+1)
	for i := 0; i < l; i++ {
		var id, size, data, k int
		fmt.Fscan(reader, &id, &size, &data, &k)
		machines := make([]int, k)
		for j := 0; j < k; j++ {
			fmt.Fscan(reader, &machines[j])
		}
		tasks[id] = Task{id: id, size: size, data: data, machines: machines}
	}

	var n int
	fmt.Fscan(reader, &n)
	power := make([]int, n+1)
	for i := 0; i < n; i++ {
		var id, p int
		fmt.Fscan(reader, &id, &p)
		power[id] = p
	}

	var m int
	fmt.Fscan(reader, &m)
	speed := make([]int, m+1)
	capacity := make([]int, m+1)
	for i := 0; i < m; i++ {
		var id, spd, cap int
		fmt.Fscan(reader, &id, &spd, &cap)
		speed[id] = spd
		capacity[id] = cap
	}

	var N int
	fmt.Fscan(reader, &N)
	adj := make([][]int, l+1)
	indeg := make([]int, l+1)
	for i := 0; i < N; i++ {
		var a, b int
		fmt.Fscan(reader, &a, &b)
		tasks[b].dataDeps = append(tasks[b].dataDeps, a)
		adj[a] = append(adj[a], b)
		indeg[b]++
	}

	var M int
	fmt.Fscan(reader, &M)
	for i := 0; i < M; i++ {
		var a, b int
		fmt.Fscan(reader, &a, &b)
		tasks[b].taskDeps = append(tasks[b].taskDeps, a)
		adj[a] = append(adj[a], b)
		indeg[b]++
	}

	// Topological order
	order := make([]int, 0, l)
	queue := make([]int, 0, l)
	for i := 1; i <= l; i++ {
		if indeg[i] == 0 {
			queue = append(queue, i)
		}
	}
	for len(queue) > 0 {
		v := queue[0]
		queue = queue[1:]
		order = append(order, v)
		for _, u := range adj[v] {
			indeg[u]--
			if indeg[u] == 0 {
				queue = append(queue, u)
			}
		}
	}

	diskUsed := make([]int, m+1)
	machineAvail := make([]int, n+1)

	start := make([]int, l+1)
	cTime := make([]int, l+1)
	dTime := make([]int, l+1)
	disk := make([]int, l+1)
	machine := make([]int, l+1)

	for _, id := range order {
		t := tasks[id]
		if t.id == 0 {
			continue
		}
		// choose machine
		y := t.machines[0]
		// choose disk with enough capacity
		z := 1
		for z <= m && diskUsed[z]+t.data > capacity[z] {
			z++
		}
		if z > m {
			z = 1
		}

		st := machineAvail[y]
		for _, j := range t.dataDeps {
			if dTime[j] > st {
				st = dTime[j]
			}
		}
		for _, j := range t.taskDeps {
			if cTime[j] > st {
				st = cTime[j]
			}
		}

		readTime := 0
		for _, j := range t.dataDeps {
			sj := speed[disk[j]]
			readTime += (tasks[j].data + sj - 1) / sj
		}
		execTime := (t.size + power[y] - 1) / power[y]
		storeTime := (t.data + speed[z] - 1) / speed[z]

		a := st
		b := a + readTime
		c := b + execTime
		d := c + storeTime

		start[id] = a
		cTime[id] = c
		dTime[id] = d
		machine[id] = y
		disk[id] = z

		machineAvail[y] = d
		diskUsed[z] += t.data
	}

	for i := 1; i <= l; i++ {
		fmt.Fprintf(writer, "%d %d %d %d\n", i, start[i], machine[i], disk[i])
	}
}
