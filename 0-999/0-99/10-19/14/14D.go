package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   in := bufio.NewReader(os.Stdin)
   var n int
   if _, err := fmt.Fscan(in, &n); err != nil {
       return
   }
   // build adjacency list
   adj := make([][]int, n)
   for i := 0; i < n-1; i++ {
       var a, b int
       fmt.Fscan(in, &a, &b)
       a--
       b--
       adj[a] = append(adj[a], b)
       adj[b] = append(adj[b], a)
   }
   maxProfit := 0
   // temporary slices
   dist := make([]int, n)
   parent := make([]int, n)
   skip := make([]bool, n)
   visited := make([]bool, n)
   // BFS queue
   queue := make([]int, 0, n)
   for u := 0; u < n; u++ {
       // BFS from u to get dist and parent
       for i := 0; i < n; i++ {
           dist[i] = -1
           parent[i] = -1
       }
       queue = queue[:0]
       dist[u] = 0
       queue = append(queue, u)
       for qi := 0; qi < len(queue); qi++ {
           v := queue[qi]
           for _, w := range adj[v] {
               if dist[w] == -1 {
                   dist[w] = dist[v] + 1
                   parent[w] = v
                   queue = append(queue, w)
               }
           }
       }
       // for each v > u, consider path u-v
       for v := u + 1; v < n; v++ {
           pathLen := dist[v]
           // mark path nodes to skip
           for i := 0; i < n; i++ {
               skip[i] = false
               visited[i] = false
           }
           x := v
           for x != -1 {
               skip[x] = true
               x = parent[x]
           }
           // find max diameter in remaining components
           best2 := 0
           // component BFS buffers
           for i := 0; i < n; i++ {
               if skip[i] || visited[i] {
                   continue
               }
               // first BFS to find farthest from i
               // dists for this BFS
               d1 := make([]int, 0, n)
               q1 := make([]int, 0, n)
               idx1 := make(map[int]bool)
               q1 = append(q1, i)
               idx1[i] = true
               visited[i] = true
               dmap := map[int]int{i: 0}
               var far int = i
               for qi := 0; qi < len(q1); qi++ {
                   v1 := q1[qi]
                   for _, w1 := range adj[v1] {
                       if skip[w1] || idx1[w1] {
                           continue
                       }
                       idx1[w1] = true
                       visited[w1] = true
                       dmap[w1] = dmap[v1] + 1
                       q1 = append(q1, w1)
                       if dmap[w1] > dmap[far] {
                           far = w1
                       }
                   }
               }
               // second BFS from far to compute diameter
               q2 := make([]int, 0, n)
               idx2 := make(map[int]bool)
               d2map := map[int]int{far: 0}
               q2 = append(q2, far)
               idx2[far] = true
               var d2max int
               for qi := 0; qi < len(q2); qi++ {
                   v2 := q2[qi]
                   for _, w2 := range adj[v2] {
                       if skip[w2] || idx2[w2] {
                           continue
                       }
                       idx2[w2] = true
                       d2map[w2] = d2map[v2] + 1
                       q2 = append(q2, w2)
                       if d2map[w2] > d2max {
                           d2max = d2map[w2]
                       }
                   }
               }
               if d2max > best2 {
                   best2 = d2max
               }
           }
           profit := pathLen * best2
           if profit > maxProfit {
               maxProfit = profit
           }
       }
   }
   fmt.Println(maxProfit)
}
