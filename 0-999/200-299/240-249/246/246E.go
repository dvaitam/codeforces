package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

// Fenwick tree for sum queries and point updates
type Fenwick struct {
   n int
   t []int
}

func NewFenwick(n int) *Fenwick {
   return &Fenwick{n: n, t: make([]int, n+1)}
}

func (f *Fenwick) Add(i, v int) {
   for ; i <= f.n; i += i & -i {
       f.t[i] += v
   }
}

func (f *Fenwick) Sum(i int) int {
   s := 0
   for ; i > 0; i -= i & -i {
       s += f.t[i]
   }
   return s
}

func (f *Fenwick) RangeSum(l, r int) int {
   if r < l {
       return 0
   }
   return f.Sum(r) - f.Sum(l-1)
}

type Query struct {
   l, r, id int
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n int
   fmt.Fscan(reader, &n)
   names := make([]string, n+1)
   parent := make([]int, n+1)
   nameID := make([]int, n+1)
   nameMap := make(map[string]int, n)
   nextID := 1

   children := make([][]int, n+1)
   for i := 1; i <= n; i++ {
       var s string
       var p int
       fmt.Fscan(reader, &s, &p)
       names[i] = s
       parent[i] = p
       if p != 0 {
           children[p] = append(children[p], i)
       }
       id, ok := nameMap[s]
       if !ok {
           id = nextID
           nameMap[s] = id
           nextID++
       }
       nameID[i] = id
   }

   depth := make([]int, n+1)
   tin := make([]int, n+1)
   tout := make([]int, n+1)
   // prepare depth lists
   tins := make([][]int, n+1)
   namesAtDepth := make([][]int, n+1)
   time := 1
   // iterative DFS for each root
   type Frame struct { u, idx int }
   stack := make([]Frame, 0, n)
   for i := 1; i <= n; i++ {
       if parent[i] != 0 {
           continue
       }
       depth[i] = 0
       stack = append(stack, Frame{i, -1})
       for len(stack) > 0 {
           top := &stack[len(stack)-1]
           u := top.u
           if top.idx == -1 {
               // enter u
               tin[u] = time
               time++
               d := depth[u]
               tins[d] = append(tins[d], tin[u])
               namesAtDepth[d] = append(namesAtDepth[d], nameID[u])
               top.idx = 0
           } else if top.idx < len(children[u]) {
               v := children[u][top.idx]
               top.idx++
               depth[v] = depth[u] + 1
               stack = append(stack, Frame{v, -1})
           } else {
               // exit u
               tout[u] = time - 1
               stack = stack[:len(stack)-1]
           }
       }
   }

   var m int
   fmt.Fscan(reader, &m)
   ans := make([]int, m)
   queriesByDepth := make([][]Query, n+1)

   for qi := 0; qi < m; qi++ {
       var v, k int
       fmt.Fscan(reader, &v, &k)
       td := depth[v] + k
       if td >= len(tins) || len(tins[td]) == 0 {
           ans[qi] = 0
           continue
       }
       // find range of tins in subtree
       arr := tins[td]
       // lower_bound for tin[v]
       l := lowerBound(arr, tin[v])
       // upper_bound for tout[v]
       r := upperBound(arr, tout[v]) - 1
       if l > r {
           ans[qi] = 0
       } else {
           // convert to 1-based indices
           queriesByDepth[td] = append(queriesByDepth[td], Query{l + 1, r + 1, qi})
       }
   }

   // process per depth
   maxNameID := nextID
   for d, qs := range queriesByDepth {
       if len(qs) == 0 {
           continue
       }
       // sort queries by r
       sort.Slice(qs, func(i, j int) bool { return qs[i].r < qs[j].r })
       // BIT over namesAtDepth[d]
       cnt := len(namesAtDepth[d])
       bit := NewFenwick(cnt)
       last := make([]int, maxNameID)
       qi := 0
       for i := 1; i <= cnt; i++ {
           id := namesAtDepth[d][i-1]
           if last[id] != 0 {
               bit.Add(last[id], -1)
           }
           bit.Add(i, 1)
           last[id] = i
           // answer queries with r == i
           for qi < len(qs) && qs[qi].r == i {
               l := qs[qi].l
               r := qs[qi].r
               ans[qs[qi].id] = bit.RangeSum(l, r)
               qi++
           }
       }
   }

   // output answers
   for i, v := range ans {
       if i > 0 {
           writer.WriteByte(' ')
       }
       fmt.Fprint(writer, v)
   }
}

// lowerBound finds first i such that a[i] >= x
func lowerBound(a []int, x int) int {
   l, r := 0, len(a)
   for l < r {
       m := l + (r-l)/2
       if a[m] < x {
           l = m + 1
       } else {
           r = m
       }
   }
   return l
}

// upperBound finds first i such that a[i] > x
func upperBound(a []int, x int) int {
   l, r := 0, len(a)
   for l < r {
       m := l + (r-l)/2
       if a[m] <= x {
           l = m + 1
       } else {
           r = m
       }
   }
   return l
}
