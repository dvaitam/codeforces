package main

import (
    "fmt"
)

const mod = 998244353

const (
    genA = 0
    genB = 1
    genb = 2
)

var inv = []int{0,2,1}

type Enumerator struct {
    next   [][]int
    parent []int
    size   []int
    queue  []int
}

func NewEnumerator() *Enumerator {
    e := &Enumerator{}
    e.makeCoset()
    return e
}

func (e *Enumerator) makeCoset() int {
    id := len(e.next)
    e.next = append(e.next, []int{-1, -1, -1})
    e.parent = append(e.parent, id)
    e.size = append(e.size, 1)
    e.queue = append(e.queue, id)
    return id
}

func (e *Enumerator) find(x int) int {
    if e.parent[x] != x {
        e.parent[x] = e.find(e.parent[x])
    }
    return e.parent[x]
}

func (e *Enumerator) link(a, g, b int) {
    ra := e.find(a)
    rb := e.find(b)
    ta := e.next[ra][g]
    if ta != -1 {
        e.union(ta, rb)
        return
    }
    e.next[ra][g] = rb
    ig := inv[g]
    tb := e.next[rb][ig]
    if tb != -1 {
        e.union(tb, ra)
    } else {
        e.next[rb][ig] = ra
    }
}

func (e *Enumerator) union(a, b int) int {
    ra := e.find(a)
    rb := e.find(b)
    if ra == rb {
        return ra
    }
    if e.size[ra] < e.size[rb] {
        ra, rb = rb, ra
    }
    e.parent[rb] = ra
    e.size[ra] += e.size[rb]
    for g := 0; g < 3; g++ {
        ta := e.next[ra][g]
        if ta != -1 {
            ta = e.find(ta)
            e.next[ra][g] = ta
        }
        tb := e.next[rb][g]
        if tb == -1 {
            continue
        }
        tb = e.find(tb)
        if ta == -1 {
            e.link(ra, g, tb)
        } else {
            e.union(ta, tb)
        }
    }
    e.queue = append(e.queue, ra)
    return ra
}

func (e *Enumerator) step(a, g int) int {
    ra := e.find(a)
    tgt := e.next[ra][g]
    if tgt != -1 {
        rt := e.find(tgt)
        e.next[ra][g] = rt
        return rt
    }
    d := e.makeCoset()
    e.link(ra, g, d)
    return e.find(e.next[ra][g])
}

func (e *Enumerator) impose(start int, rel []int) {
    cur := e.find(start)
    for _, g := range rel {
        cur = e.step(cur, g)
    }
    e.union(cur, start)
}

func main() {
    s := "ABAB"
    relS := make([]int, len(s))
    for i, ch := range s {
        if ch == 'A' {
            relS[i] = genA
        } else {
            relS[i] = genB
        }
    }
    relations := [][]int{
        {genA, genA},
        {genB, genB, genB},
        {genB, genb},
        {genb, genB},
        relS,
    }
    e := NewEnumerator()
    seen := make(map[int]bool)
    for len(e.queue) > 0 {
        idx := e.queue[0]
        e.queue = e.queue[1:]
        root := e.find(idx)
        if seen[root] {
            continue
        }
        seen[root] = true
        for _, rel := range relations {
            e.impose(root, rel)
        }
        for _, g := range []int{genA, genB} {
            nxt := e.step(root, g)
            e.queue = append(e.queue, nxt)
        }
    }
    reps := map[int]bool{}
    for i := range e.next {
        reps[e.find(i)] = true
    }
    fmt.Println("num reps", len(reps))
    for rep := range reps {
        fmt.Println(rep, e.next[rep])
    }
}
