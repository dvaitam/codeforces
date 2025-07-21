package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

type edge struct { to, w int }
type pair struct { d, w int }

// BIT with versioning to allow fast reset
type BIT struct {
   n    int
   tree []int
   ver  []int
   time int
}

func NewBIT(n int) *BIT {
   return &BIT{n: n, tree: make([]int, n+1), ver: make([]int, n+1), time: 1}
}

func (b *BIT) reset() {
   b.time++
}

func (b *BIT) update(i, v int) {
   for ; i <= b.n; i += i & -i {
       if b.ver[i] != b.time {
           b.ver[i] = b.time
           b.tree[i] = 0
       }
       b.tree[i] += v
   }
}

func (b *BIT) query(i int) int {
   if i > b.n {
       i = b.n
   }
   s := 0
   for ; i > 0; i -= i & -i {
       if b.ver[i] == b.time {
           s += b.tree[i]
       }
   }
   return s
}

var (
   n, L, W int
   adj      [][]edge
   removed  []bool
   subSize  []int
   total    int64
   bit      *BIT
)

func dfsSize(u, p int) int {
   subSize[u] = 1
   for _, e := range adj[u] {
       v := e.to
       if v != p && !removed[v] {
           subSize[u] += dfsSize(v, u)
       }
   }
   return subSize[u]
}

func findCentroid(u, p, sz int) int {
   for _, e := range adj[u] {
       v := e.to
       if v != p && !removed[v] && subSize[v] > sz/2 {
           return findCentroid(v, u, sz)
       }
   }
   return u
}

// collect all nodes in subtree of u, iterative
func collect(u, p int, maxD, maxW int) []pair {
   var res []pair
   // stack elements: u, parent, depth, weight
   type st struct{ u, p, d, w int }
   stack := []st{{u, p, 1, 0}}
   // We need initial weight: weight of edge from centroid to u was passed in w field of first element
   for len(stack) > 0 {
       cur := stack[len(stack)-1]
       stack = stack[:len(stack)-1]
       if cur.d > L || cur.w > W {
           continue
       }
       res = append(res, pair{cur.d, cur.w})
       for _, e := range adj[cur.u] {
           if e.to != cur.p && !removed[e.to] {
               stack = append(stack, st{e.to, cur.u, cur.d + 1, cur.w + e.w})
           }
       }
   }
   return res
}

func decompose(u int) {
   sz := dfsSize(u, -1)
   c := findCentroid(u, -1, sz)
   removed[c] = true
   // vec: sorted by weight
   vec := []pair{{0, 0}}
   // process each subtree
   for _, e := range adj[c] {
       v := e.to
       if removed[v] {
           continue
       }
       sub := collect(v, c, L, W)
       // sort both lists by weight
       sort.Slice(sub, func(i, j int) bool { return sub[i].w < sub[j].w })
       sort.Slice(vec, func(i, j int) bool { return vec[i].w < vec[j].w })
       // count pairs between vec and sub
       bit.reset()
       ptr := 0
       for _, p := range sub {
           limW := W - p.w
           // add vec points with weight <= limW
           for ptr < len(vec) && vec[ptr].w <= limW {
               if vec[ptr].d <= L {
                   bit.update(vec[ptr].d+1, 1)
               }
               ptr++
           }
           remD := L - p.d
           if remD >= 0 {
               // query depths <= remD, shifted by +1
               cnt := bit.query(remD + 1)
               total += int64(cnt)
           }
       }
       // merge vec and sub into vec
       merged := make([]pair, 0, len(vec)+len(sub))
       i, j := 0, 0
       for i < len(vec) && j < len(sub) {
           if vec[i].w < sub[j].w {
               merged = append(merged, vec[i]); i++
           } else {
               merged = append(merged, sub[j]); j++
           }
       }
       for i < len(vec) {
           merged = append(merged, vec[i]); i++
       }
       for j < len(sub) {
           merged = append(merged, sub[j]); j++
       }
       vec = merged
   }
   // count pairs where both endpoints are centroid only counted above
   // recurse
   for _, e := range adj[c] {
       if !removed[e.to] {
           decompose(e.to)
       }
   }
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   fmt.Fscan(reader, &n, &L, &W)
   adj = make([][]edge, n+1)
   for i := 2; i <= n; i++ {
       var p, w int
       fmt.Fscan(reader, &p, &w)
       adj[i] = append(adj[i], edge{p, w})
       adj[p] = append(adj[p], edge{i, w})
   }
   removed = make([]bool, n+1)
   subSize = make([]int, n+1)
   bit = NewBIT(L + 1)
   total = 0
   decompose(1)
   fmt.Fprint(writer, total)
}
