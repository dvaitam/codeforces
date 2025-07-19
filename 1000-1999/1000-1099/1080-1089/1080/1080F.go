package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

// Que holds interval data
type Que struct { l, r, p int }

// Node for persistent segment tree
type Node struct { lson, rson, mn int }

var (
   n, m, K   int
   savQue     []Que
   dataArr    []int
   rt         []int
   tree       []Node
   cnt        int
)

func min(a, b int) int { if a < b { return a }; return b }
func max(a, b int) int { if a > b { return a }; return b }

// update returns new root after inserting/updating position pos with value val
func update(v, l, r, pos, val int) int {
   cnt++
   u := cnt
   // copy node
   tree[u] = tree[v]
   if l == r {
       tree[u].mn = max(tree[u].mn, val)
       return u
   }
   mid := (l + r) >> 1
   if pos <= mid {
       tree[u].lson = update(tree[v].lson, l, mid, pos, val)
   } else {
       tree[u].rson = update(tree[v].rson, mid+1, r, pos, val)
   }
   tree[u].mn = min(tree[tree[u].lson].mn, tree[tree[u].rson].mn)
   return u
}

// query returns min mn in [L, R]
func query(v, l, r, L, R int) int {
   const INF = int(1e18)
   if v == 0 || l > R || r < L {
       return INF
   }
   if l >= L && r <= R {
       return tree[v].mn
   }
   mid := (l + r) >> 1
   res := INF
   if L <= mid {
       res = min(res, query(tree[v].lson, l, mid, L, R))
   }
   if R > mid {
       res = min(res, query(tree[v].rson, mid+1, r, L, R))
   }
   return res
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   fmt.Fscan(reader, &n, &m, &K)
   savQue = make([]Que, K+1)
   dataArr = make([]int, K+1)
   rt = make([]int, K+1)
   for i := 1; i <= K; i++ {
       fmt.Fscan(reader, &savQue[i].l, &savQue[i].r, &savQue[i].p)
       dataArr[i] = savQue[i].r
   }
   sort.Ints(dataArr[1:])
   sort.Slice(savQue[1:], func(i, j int) bool {
       return savQue[i+1].r < savQue[j+1].r
   })
   // unique dataArr
   cc := 0
   for i := 1; i <= K; i++ {
       if i == 1 || dataArr[i] != dataArr[i-1] {
           cc++
           dataArr[cc] = dataArr[i]
       }
   }
   // prepare tree nodes
   maxNodes := (K + 5) * 25
   tree = make([]Node, maxNodes)
   cnt = 0
   prevR := -1
   // build persistent versions
   for i := 1; i <= K; i++ {
       rVal := savQue[i].r
       // find index x of rVal in dataArr[1..cc]
       lo, hi := 1, cc
       for lo < hi {
           mid := (lo + hi) >> 1
           if dataArr[mid] < rVal {
               lo = mid + 1
           } else {
               hi = mid
           }
       }
       x := lo
       if i > 1 && rVal == prevR {
           rt[x] = update(rt[x], 1, n, savQue[i].p, savQue[i].l)
       } else {
           rt[x] = update(rt[x-1], 1, n, savQue[i].p, savQue[i].l)
       }
       prevR = rVal
   }
   // answer queries
   for i := 0; i < m; i++ {
       var a, b, xq, y int
       fmt.Fscan(reader, &a, &b, &xq, &y)
       // upper_bound for y
       lo, hi := 1, cc+1
       for lo < hi {
           mid := (lo + hi) >> 1
           if mid <= cc && dataArr[mid] <= y {
               lo = mid + 1
           } else {
               hi = mid
           }
       }
       p := lo - 1
       if p < 1 {
           fmt.Fprintln(writer, "no")
       } else {
           res := query(rt[p], 1, n, a, b)
           if res >= xq {
               fmt.Fprintln(writer, "yes")
           } else {
               fmt.Fprintln(writer, "no")
           }
       }
   }
}
