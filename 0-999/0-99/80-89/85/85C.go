package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n int
   if _, err := fmt.Fscan(reader, &n); err != nil {
       return
   }
   parent := make([]int, n)
   lc := make([]int, n)
   rc := make([]int, n)
   val := make([]int, n)
   lmax := make([]int, n)
   rmin := make([]int, n)
   ma := make([]int, n)
   mi := make([]int, n)
   for i := range lc {
       lc[i], rc[i], lmax[i], rmin[i], ma[i], mi[i] = -1, -1, -1, -1, -1, -1
   }
   var root int
   for i := 0; i < n; i++ {
       var p, v int
       fmt.Fscan(reader, &p, &v)
       if p != -1 {
           p--
       }
       parent[i] = p
       val[i] = v
       if p != -1 {
           if lc[p] != -1 {
               rc[p] = i
           } else {
               lc[p] = i
           }
       } else {
           root = i
       }
   }
   // ensure left child has smaller value
   for i := 0; i < n; i++ {
       if rc[i] != -1 && val[lc[i]] > val[rc[i]] {
           lc[i], rc[i] = rc[i], lc[i]
       }
   }
   // compute metrics bottom-up
   stack := []int{root}
   for len(stack) > 0 {
       cur := stack[len(stack)-1]
       if lc[cur] == -1 {
           ma[cur] = val[cur]
           mi[cur] = val[cur]
           lmax[cur] = val[cur]
           rmin[cur] = val[cur]
           stack = stack[:len(stack)-1]
           continue
       }
       if lmax[lc[cur]] != -1 {
           mi[cur] = mi[lc[cur]]
           ma[cur] = ma[rc[cur]]
           lmax[cur] = ma[lc[cur]]
           rmin[cur] = mi[rc[cur]]
           stack = stack[:len(stack)-1]
           continue
       }
       stack = append(stack, lc[cur], rc[cur])
   }
   // compute answer for leaves
   ans := make([]float64, n)
   type State struct { pos, count int; sum int64 }
   st := []State{{root, 0, 0}}
   for len(st) > 0 {
       cur := st[len(st)-1]
       st = st[:len(st)-1]
       if lc[cur.pos] == -1 {
           ans[cur.pos] = float64(cur.sum) / float64(cur.count)
           continue
       }
       st = append(st, State{lc[cur.pos], cur.count + 1, cur.sum + int64(rmin[cur.pos])})
       st = append(st, State{rc[cur.pos], cur.count + 1, cur.sum + int64(lmax[cur.pos])})
   }
   // collect leaves and inner nodes
   type Leaf struct { val int; ans float64 }
   var leafs []Leaf
   var inner []int
   for i := 0; i < n; i++ {
       if lc[i] == -1 {
           leafs = append(leafs, Leaf{val[i], ans[i]})
       } else {
           inner = append(inner, val[i])
       }
   }
   sort.Slice(leafs, func(i, j int) bool { return leafs[i].val < leafs[j].val })
   sort.Ints(inner)
   var qcount int
   fmt.Fscan(reader, &qcount)
   for i := 0; i < qcount; i++ {
       var q int
       fmt.Fscan(reader, &q)
       pos := sort.Search(len(inner), func(i int) bool { return inner[i] > q })
       fmt.Printf("%.9f\n", leafs[pos].ans)
   }
}
