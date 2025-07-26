package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

const (
   MAXB = 1000001
   INF64 = 9e18
)

type Monster struct {
   x, y int
   z    int64
}

var (
   treeMax []int64
   treeAdd []int64
)

func build(node, l, r int, base []int64) {
   if l == r {
       treeMax[node] = base[l]
       return
   }
   mid := (l + r) >> 1
   left, right := node<<1, node<<1|1
   build(left, l, mid, base)
   build(right, mid+1, r, base)
   treeMax[node] = max(treeMax[left], treeMax[right])
}

func apply(node int, v int64) {
   treeMax[node] += v
   treeAdd[node] += v
}

func push(node int) {
   if treeAdd[node] != 0 {
       apply(node<<1, treeAdd[node])
       apply(node<<1|1, treeAdd[node])
       treeAdd[node] = 0
   }
}

func update(node, l, r, ql, qr int, v int64) {
   if ql > r || qr < l {
       return
   }
   if ql <= l && r <= qr {
       apply(node, v)
       return
   }
   push(node)
   mid := (l + r) >> 1
   update(node<<1, l, mid, ql, qr, v)
   update(node<<1|1, mid+1, r, ql, qr, v)
   treeMax[node] = max(treeMax[node<<1], treeMax[node<<1|1])
}

func max(a, b int64) int64 {
   if a > b {
       return a
   }
   return b
}

func main() {
   in := bufio.NewReader(os.Stdin)
   var n, m, p int
   fmt.Fscan(in, &n, &m, &p)
   wp := make(map[int]int64)
   for i := 0; i < n; i++ {
       var a, c int64
       fmt.Fscan(in, &a, &c)
       prev, ok := wp[int(a)]
       if !ok || c < prev {
           wp[int(a)] = c
       }
   }
   // armor
   bestCost := make([]int64, MAXB+2)
   for i := range bestCost {
       bestCost[i] = INF64
   }
   for i := 0; i < m; i++ {
       var b, c int64
       fmt.Fscan(in, &b, &c)
       if c < bestCost[int(b)] {
           bestCost[int(b)] = c
       }
   }
   // monsters
   monsters := make([]Monster, p)
   for i := 0; i < p; i++ {
       var x, y int
       var z int64
       fmt.Fscan(in, &x, &y, &z)
       monsters[i] = Monster{x, y, z}
   }
   sort.Slice(monsters, func(i, j int) bool {
       return monsters[i].x < monsters[j].x
   })
   // weapons sorted
   type W struct{ a int; c int64 }
   weps := make([]W, 0, len(wp))
   for a, c := range wp {
       weps = append(weps, W{a, c})
   }
   sort.Slice(weps, func(i, j int) bool { return weps[i].a < weps[j].a })
   // build segment tree over armor b: 0..MAXB
   size := MAXB + 1
   treeMax = make([]int64, 4*(size+1))
   treeAdd = make([]int64, 4*(size+1))
   base := make([]int64, size+1)
   for i := 0; i <= size; i++ {
       if bestCost[i] < INF64 {
           base[i] = -bestCost[i]
       } else {
           base[i] = -INF64
       }
   }
   build(1, 0, size, base)
   // sweep
   ans := -INF64
   mi := 0
   for _, w := range weps {
       A := w.a
       // add monsters with x < A
       for mi < len(monsters) && monsters[mi].x < A {
           y := monsters[mi].y
           z := monsters[mi].z
           if y+1 <= size {
               update(1, 0, size, y+1, size, z)
           }
           mi++
       }
       cur := treeMax[1] - w.c
       if cur > ans {
           ans = cur
       }
   }
   fmt.Println(ans)
}
