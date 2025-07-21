package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var n int
   if _, err := fmt.Fscan(reader, &n); err != nil {
       return
   }
   c := make([]int, n+1)
   for i := 1; i <= n; i++ {
       fmt.Fscan(reader, &c[i])
   }
   adj := make([][]int, n+1)
   indeg0 := make([]int, n+1)
   for i := 1; i <= n; i++ {
       var k int
       fmt.Fscan(reader, &k)
       for j := 0; j < k; j++ {
           var u int
           fmt.Fscan(reader, &u)
           // u must be before i
           adj[u] = append(adj[u], i)
           indeg0[i]++
       }
   }
   // movement cost matrix 1-based
   moveCost := [4][4]int{
       {},
       {0, 0, 1, 2},
       {0, 2, 0, 1},
       {0, 1, 2, 0},
   }
   const INF = 1<<60
   res := INF
   // simulate for each start computer
   for start := 1; start <= 3; start++ {
       // copy indegree
       indeg := make([]int, n+1)
       copy(indeg, indeg0)
       // avail tasks per computer
       avail := make([][]int, 4)
       for i := 1; i <= n; i++ {
           if indeg[i] == 0 {
               avail[c[i]] = append(avail[c[i]], i)
           }
       }
       curr := start
       cost := 0
       done := 0
       // simulate
       for done < n {
           if len(avail[curr]) > 0 {
               // take next task
               u := avail[curr][0]
               avail[curr] = avail[curr][1:]
               cost += 1
               done++
               for _, v := range adj[u] {
                   indeg[v]--
                   if indeg[v] == 0 {
                       avail[c[v]] = append(avail[c[v]], v)
                   }
               }
           } else {
               // need to move to next computer with available tasks
               bestM := 0
               best := INF
               for m := 1; m <= 3; m++ {
                   if len(avail[m]) > 0 {
                       if moveCost[curr][m] < best {
                           best = moveCost[curr][m]
                           bestM = m
                       }
                   }
               }
               // if no avail tasks anywhere (shouldn't happen), break
               if bestM == 0 {
                   break
               }
               cost += best
               curr = bestM
           }
       }
       if done == n && cost < res {
           res = cost
       }
   }
   // output result
   fmt.Println(res)
}
