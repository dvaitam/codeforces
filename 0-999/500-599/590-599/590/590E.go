package main

import (
   "bufio"
   "fmt"
   "os"
   "math/bits"
)

// Node for Aho-Corasick
type Node struct {
   nxt  [2]int
   fail int
   end  int
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   for {
       var n int
       if _, err := fmt.Fscan(reader, &n); err != nil {
           break
       }
       strs := make([]string, n)
       for i := 0; i < n; i++ {
           fmt.Fscan(reader, &strs[i])
       }
       res := solve(n, strs)
       fmt.Fprintln(writer, len(res))
       for i, v := range res {
           if i > 0 {
               fmt.Fprint(writer, " ")
           }
           fmt.Fprint(writer, v)
       }
       fmt.Fprintln(writer)
   }
}

func solve(n int, strs []string) []int {
   // Build Aho-Corasick trie
   nodes := make([]Node, 1)
   nodes[0].nxt[0], nodes[0].nxt[1] = -1, -1
   nodes[0].fail = 0
   nodes[0].end = -1
   root := 0

   // Bitset adjacency
   L := (n + 63) >> 6
   g := make([][]uint64, n)
   for i := 0; i < n; i++ {
       g[i] = make([]uint64, L)
   }

   // Insert patterns and mark direct equal overlaps
   for idx, s := range strs {
       now := root
       for i := 0; i < len(s); i++ {
           c := s[i] - 'a'
           if c < 0 || c > 1 {
               continue
           }
           if nodes[now].nxt[c] == -1 {
               nodes = append(nodes, Node{nxt: [2]int{-1, -1}, fail: 0, end: -1})
               nodes[now].nxt[c] = len(nodes) - 1
           }
           now = nodes[now].nxt[c]
       }
       if nodes[now].end != -1 {
           j := nodes[now].end
           g[idx][j>>6] |= 1 << uint(j&63)
           g[j][idx>>6] |= 1 << uint(idx&63)
       }
       nodes[now].end = idx
   }

   // Build failure links
   queue := make([]int, 0)
   for c := 0; c < 2; c++ {
       v := nodes[root].nxt[c]
       if v == -1 {
           nodes[root].nxt[c] = root
       } else {
           nodes[v].fail = root
           queue = append(queue, v)
       }
   }
   for qi := 0; qi < len(queue); qi++ {
       u := queue[qi]
       if nodes[u].end == -1 {
           nodes[u].end = nodes[nodes[u].fail].end
       }
       for c := 0; c < 2; c++ {
           v := nodes[u].nxt[c]
           if v == -1 {
               nodes[u].nxt[c] = nodes[nodes[u].fail].nxt[c]
           } else {
               nodes[v].fail = nodes[nodes[u].fail].nxt[c]
               queue = append(queue, v)
           }
       }
   }

   // Query each string
   for idx, s := range strs {
       now := root
       for i := 0; i < len(s); i++ {
           c := s[i] - 'a'
           if c < 0 || c > 1 {
               continue
           }
           now = nodes[now].nxt[c]
           x := nodes[nodes[now].fail].end
           y := nodes[now].end
           if i == len(s)-1 {
               if x != -1 {
                   g[x][idx>>6] |= 1 << uint(idx&63)
               }
           }
           if i < len(s)-1 && y != -1 {
               g[y][idx>>6] |= 1 << uint(idx&63)
           }
       }
   }

   // Transitive closure (Floyd-Warshall via bitsets)
   for k := 0; k < n; k++ {
       kb := k >> 6
       km := uint(k & 63)
       for i := 0; i < n; i++ {
           if g[i][kb]&(1<<km) != 0 {
               for w := 0; w < L; w++ {
                   g[i][w] |= g[k][w]
               }
           }
       }
   }

   // Build adjacency list
   adj := make([][]int, n)
   for i := 0; i < n; i++ {
       for w := 0; w < L; w++ {
           bits := g[i][w]
           for bits != 0 {
               b := bits & -bits
               r := bitsTrailing(b)
               j := w*64 + r
               if j < n {
                   adj[i] = append(adj[i], j)
               }
               bits ^= b
           }
       }
   }

   // Hopcroft-Karp
   const INF = int(1e9)
   pairU := make([]int, n)
   pairV := make([]int, n)
   dist := make([]int, n)
   for i := 0; i < n; i++ {
       pairU[i], pairV[i] = -1, -1
   }
   
   var bfsHK func() bool
   var dfsHK func(u int) bool
   
   bfsHK = func() bool {
       queue := make([]int, 0)
       for u := 0; u < n; u++ {
           if pairU[u] == -1 {
               dist[u] = 0
               queue = append(queue, u)
           } else {
               dist[u] = INF
           }
       }
       found := false
       for qi := 0; qi < len(queue); qi++ {
           u := queue[qi]
           for _, v := range adj[u] {
               if pairV[v] != -1 {
                   if dist[pairV[v]] == INF {
                       dist[pairV[v]] = dist[u] + 1
                       queue = append(queue, pairV[v])
                   }
               } else {
                   found = true
               }
           }
       }
       return found
   }

   dfsHK = func(u int) bool {
       for _, v := range adj[u] {
           if pairV[v] == -1 || (dist[pairV[v]] == dist[u]+1 && dfsHK(pairV[v])) {
               pairU[u], pairV[v] = v, u
               return true
           }
       }
       dist[u] = INF
       return false
   }

   matching := 0
   for bfsHK() {
       for u := 0; u < n; u++ {
           if pairU[u] == -1 && dfsHK(u) {
               matching++
           }
       }
   }

   // Minimum vertex cover via alternating BFS
   usex := make([]bool, n)
   usey := make([]bool, n)
   for i := 0; i < n; i++ {
       if pairU[i] != -1 {
           usex[i] = true
       }
   }
   vis2 := make([]bool, n)
   for {
       changed := false
       for i := 0; i < n; i++ {
           if !usex[i] && !vis2[i] {
               vis2[i] = true
               for _, j := range adj[i] {
                   if !usey[j] {
                       usey[j] = true
                       if pairV[j] != -1 {
                           usex[pairV[j]] = false
                       }
                       changed = true
                   }
               }
           }
       }
       if !changed {
           break
       }
   }

   // Collect result: nodes not in vertex cover
   res := make([]int, 0)
   for i := 0; i < n; i++ {
       if !usex[i] && !usey[i] {
           res = append(res, i+1)
       }
   }
   return res
}

func bitsTrailing(b uint64) int {
   return bits.TrailingZeros64(b)
}
