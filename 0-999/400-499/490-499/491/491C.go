package main

import (
   "bufio"
   "fmt"
   "os"
)

const maxn = 70

var (
   cost   [maxn][maxn]int
   u, v   [maxn]int
   dist   [maxn]int
   dad    [maxn]int
   seen   [maxn]bool
   Lmate  [maxn]int
   Rmate  [maxn]int
   n, m   int
)

func getId(ch byte) int {
   if ch >= 'a' && ch <= 'z' {
       return int(ch - 'a')
   }
   return int(ch - 'A' + 26)
}

func getCh(id int) byte {
   if id >= 26 {
       return byte('A' + id - 26)
   }
   return byte('a' + id)
}

// MinCostMatching finds minimum cost perfect matching in a complete bipartite graph
// with cost matrix cost[0..n-1][0..n-1]. Returns total cost and Lmate where
// Lmate[i] is the matched right vertex for left vertex i.
func MinCostMatching() (int, []int) {
   // potentials
   for i := 0; i < n; i++ {
       u[i] = cost[i][0]
       for j := 1; j < n; j++ {
           if cost[i][j] < u[i] {
               u[i] = cost[i][j]
           }
       }
   }
   for j := 0; j < n; j++ {
       v[j] = cost[0][j] - u[0]
       for i := 1; i < n; i++ {
           if cost[i][j]-u[i] < v[j] {
               v[j] = cost[i][j] - u[i]
           }
       }
   }
   for i := 0; i < n; i++ {
       Lmate[i] = -1
   }
   for j := 0; j < n; j++ {
       Rmate[j] = -1
   }
   mated := 0
   // initial greedy match
   for i := 0; i < n; i++ {
       for j := 0; j < n; j++ {
           if Rmate[j] != -1 {
               continue
           }
           if cost[i][j]-u[i]-v[j] == 0 {
               Lmate[i] = j
               Rmate[j] = i
               mated++
               break
           }
       }
   }
   for mated < n {
       // find unmatched left
       var s int
       for s = 0; s < n; s++ {
           if Lmate[s] == -1 {
               break
           }
       }
       // initialize Dijkstra
       for k := 0; k < n; k++ {
           dad[k] = -1
           seen[k] = false
           dist[k] = cost[s][k] - u[s] - v[k]
       }
       var j int
       for {
           // find unseen with minimal dist
           j = -1
           for k := 0; k < n; k++ {
               if seen[k] {
                   continue
               }
               if j == -1 || dist[k] < dist[j] {
                   j = k
               }
           }
           seen[j] = true
           if Rmate[j] == -1 {
               break
           }
           // relax edges
           i := Rmate[j]
           for k := 0; k < n; k++ {
               if seen[k] {
                   continue
               }
               newDist := dist[j] + cost[i][k] - u[i] - v[k]
               if dist[k] > newDist {
                   dist[k] = newDist
                   dad[k] = j
               }
           }
       }
       // update potentials
       for k := 0; k < n; k++ {
           if k == j || !seen[k] {
               continue
           }
           i := Rmate[k]
           v[k] += dist[k] - dist[j]
           u[i] -= dist[k] - dist[j]
       }
       u[s] += dist[j]
       // augment along path
       for pj := j; dad[pj] >= 0; pj = dad[pj] {
           d := dad[pj]
           Rmate[pj] = Rmate[d]
           Lmate[Rmate[pj]] = pj
       }
       Rmate[j] = s
       Lmate[s] = j
       mated++
   }
   // compute result
   value := 0
   for i := 0; i < n; i++ {
       value += cost[i][Lmate[i]]
   }
   // copy match
   match := make([]int, n)
   for i := 0; i < n; i++ {
       match[i] = Lmate[i]
   }
   return value, match
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   if _, err := fmt.Fscan(reader, &m, &n); err != nil {
       return
   }
   var s1, s2 string
   fmt.Fscan(reader, &s1, &s2)
   // initialize cost
   for i := 0; i < maxn; i++ {
       for j := 0; j < maxn; j++ {
           cost[i][j] = 0
       }
   }
   // fill cost
   for i := 0; i < m; i++ {
       a := getId(s1[i])
       b := getId(s2[i])
       cost[a][b]--
   }
   // solve
   value, match := MinCostMatching()
   // output
   fmt.Println(-value)
   for i := 0; i < n; i++ {
       fmt.Printf("%c", getCh(match[i]))
   }
   fmt.Println()
}
