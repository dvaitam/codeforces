package main

import (
    "bufio"
    "fmt"
    "os"
    "sort"
)

type Group struct {
    vertices []int
}

func popcount(x int) int {
    c := 0
    for x > 0 {
        c += x & 1
        x >>= 1
    }
    return c
}

func main() {
    reader := bufio.NewReader(os.Stdin)
    var n, m int
    fmt.Fscan(reader, &n, &m)
    adj := make([][]bool, n)
    for i := range adj {
        adj[i] = make([]bool, n)
    }
    for i := 0; i < m; i++ {
        var u, v int
        fmt.Fscan(reader, &u, &v)
        u--
        v--
        adj[u][v] = true
        adj[v][u] = true
    }

    bitlen := (n + 7) / 8
    getKey := func(i int) string {
        b := make([]byte, bitlen)
        for j := 0; j < n; j++ {
            if adj[i][j] {
                b[j/8] |= 1 << uint(j%8)
            }
        }
        return string(b)
    }

    groupsMap := make(map[string]int)
    var groups []Group
    belong := make([]int, n)
    for i := 0; i < n; i++ {
        key := getKey(i)
        id, ok := groupsMap[key]
        if !ok {
            id = len(groups)
            groupsMap[key] = id
            groups = append(groups, Group{})
        }
        groups[id].vertices = append(groups[id].vertices, i)
        belong[i] = id
    }
    g := len(groups)

    gAdj := make([][]bool, g)
    for i := 0; i < g; i++ {
        gAdj[i] = make([]bool, g)
    }
    for i := 0; i < g; i++ {
        for j := i + 1; j < g; j++ {
            v1 := groups[i].vertices[0]
            v2 := groups[j].vertices[0]
            if adj[v1][v2] {
                gAdj[i][j] = true
                gAdj[j][i] = true
            }
        }
    }

    cubeAdj := make([][]bool, 16)
    for i := 0; i < 16; i++ {
        cubeAdj[i] = make([]bool, 16)
        for j := 0; j < 16; j++ {
            if popcount(i^j) == 1 {
                cubeAdj[i][j] = true
            }
        }
    }

    degree := make([]int, g)
    for i := 0; i < g; i++ {
        for j := 0; j < g; j++ {
            if gAdj[i][j] {
                degree[i]++
            }
        }
    }
    order := make([]int, g)
    for i := range order {
        order[i] = i
    }
    sort.Slice(order, func(a, b int) bool {
        if degree[order[a]] == degree[order[b]] {
            return order[a] < order[b]
        }
        return degree[order[a]] > degree[order[b]]
    })

    mapping := make([]int, g)
    for i := range mapping {
        mapping[i] = -1
    }
    used := make([]bool, 16)

    mapping[order[0]] = 0
    used[0] = true

    var dfs func(pos int) bool
    dfs = func(pos int) bool {
        if pos == g {
            return true
        }
        idx := order[pos]
        if mapping[idx] != -1 {
            return dfs(pos + 1)
        }
        for num := 0; num < 16; num++ {
            if used[num] {
                continue
            }
            ok := true
            for j := 0; j < g; j++ {
                if mapping[j] != -1 {
                    if gAdj[idx][j] != cubeAdj[num][mapping[j]] {
                        ok = false
                        break
                    }
                }
            }
            if !ok {
                continue
            }
            mapping[idx] = num
            used[num] = true
            if dfs(pos + 1) {
                return true
            }
            mapping[idx] = -1
            used[num] = false
        }
        return false
    }

    if !dfs(1) { // start from pos 1 since order[0] preassigned
        fmt.Println("No solution")
        return
    }

    res := make([]int, n)
    for i := 0; i < n; i++ {
        res[i] = mapping[belong[i]]
    }
    writer := bufio.NewWriter(os.Stdout)
    defer writer.Flush()
    for i, val := range res {
        if i > 0 {
            fmt.Fprint(writer, " ")
        }
        fmt.Fprint(writer, val)
    }
    fmt.Fprintln(writer)
}

