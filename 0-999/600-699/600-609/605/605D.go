package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

const INF = 2000000001

// Spell holds the parameters of a spell and its original index
type Spell struct {
   a, b, c, d int
   ord        int
}

func main() {
   in := bufio.NewReader(os.Stdin)
   var n int
   if _, err := fmt.Fscan(in, &n); err != nil {
       return
   }
   // read spells
   spells := make([]Spell, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(in, &spells[i].a, &spells[i].b, &spells[i].c, &spells[i].d)
       spells[i].ord = i
   }
   // sort by a
   sort.Slice(spells, func(i, j int) bool {
       return spells[i].a < spells[j].a
   })
   // append dummy spell at end for segment tree
   dummyIdx := n
   spells = append(spells, Spell{a: 0, b: INF, c: 0, d: 0, ord: n})

   // segment tree size
   M := len(spells)
   S := 1
   for S < M {
       S <<= 1
   }
   seg := make([]int, 2*S)
   // initialize leaves
   for i := 0; i < S; i++ {
       if i < M {
           seg[S+i] = i
       } else {
           seg[S+i] = dummyIdx
       }
   }
   // build internal nodes
   for i := S - 1; i > 0; i-- {
       l, r := seg[2*i], seg[2*i+1]
       if spells[l].b <= spells[r].b {
           seg[i] = l
       } else {
           seg[i] = r
       }
   }

   // helper: remove index i from tree
   pop := func(idx int) {
       pos := S + idx
       seg[pos] = dummyIdx
       for pos >>= 1; pos > 0; pos >>= 1 {
           l, r := seg[2*pos], seg[2*pos+1]
           if spells[l].b <= spells[r].b {
               seg[pos] = l
           } else {
               seg[pos] = r
           }
       }
   }
   // helper: query min-b index in [l..r]
   query := func(l, r int) int {
       l += S; r += S
       res := dummyIdx
       for l <= r {
           if l&1 == 1 {
               if spells[seg[l]].b < spells[res].b {
                   res = seg[l]
               }
               l++
           }
           if r&1 == 0 {
               if spells[seg[r]].b < spells[res].b {
                   res = seg[r]
               }
               r--
           }
           l >>= 1; r >>= 1
       }
       return res
   }
   // binary search for max i with spells[i].a <= x
   bsrc := func(x int) int {
       idx := sort.Search(n, func(i int) bool { return spells[i].a > x }) - 1
       return idx
   }

   // prepare BFS
   dist := make([]int, n)
   parent := make([]int, n)
   visited := make([]bool, n)
   for i := 0; i < n; i++ {
       parent[i] = -2
   }
   // initial spells with a==0 && b==0
   queue := make([]int, 0, n)
   for i := 0; i < n; i++ {
       if spells[i].a == 0 && spells[i].b == 0 {
           // mark
           visited[spells[i].ord] = true
           dist[spells[i].ord] = 1
           parent[spells[i].ord] = -1
           queue = append(queue, i)
           pop(i)
       }
   }
   // if no initial
   if len(queue) == 0 {
       fmt.Println(-1)
       return
   }
   // BFS
   for qi := 0; qi < len(queue); qi++ {
       cur := queue[qi]
       curOrd := spells[cur].ord
       // find possible next spells
       r := bsrc(spells[cur].c)
       for r >= 0 {
           nxt := query(0, r)
           if nextB := spells[nxt].b; nextB > spells[cur].d {
               break
           }
           // visit nxt
           pop(nxt)
           ord := spells[nxt].ord
           if !visited[ord] {
               visited[ord] = true
               dist[ord] = dist[curOrd] + 1
               parent[ord] = curOrd
               queue = append(queue, nxt)
           }
       }
   }
   // check target original index n-1
   if !visited[n-1] {
       fmt.Println(-1)
       return
   }
   // reconstruct path
   path := []int{}
   for cur := n - 1; cur != -1; cur = parent[cur] {
       path = append(path, cur+1)
   }
   // reverse
   for i, j := 0, len(path)-1; i < j; i, j = i+1, j-1 {
       path[i], path[j] = path[j], path[i]
   }
   // output
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   fmt.Fprintln(writer, dist[n-1])
   for i, v := range path {
       if i > 0 {
           writer.WriteByte(' ')
       }
       fmt.Fprint(writer, v)
   }
   writer.WriteByte('\n')
}
