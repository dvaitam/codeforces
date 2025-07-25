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
   const V = 4
   // parity of degrees
   var bMask int
   // total sum of values
   var total int64
   // initialize distance matrix
   const INF = int64(4e18)
   dist := make([][]int64, V)
   for i := 0; i < V; i++ {
       dist[i] = make([]int64, V)
       for j := 0; j < V; j++ {
           if i == j {
               dist[i][j] = 0
           } else {
               dist[i][j] = INF
           }
       }
   }
   // read blocks
   for i := 0; i < n; i++ {
       var c1, c2 int
       var v int64
       fmt.Fscan(in, &c1, &v, &c2)
       // make zero-based
       u := c1 - 1
       w := c2 - 1
       total += v
       // toggle parity
       bMask ^= 1 << u
       bMask ^= 1 << w
       // undirected edge, consider minimal for direct
       if v < dist[u][w] {
           dist[u][w] = v
           dist[w][u] = v
       }
   }
   // Floyd-Warshall for all-pairs shortest paths
   for k := 0; k < V; k++ {
       for i := 0; i < V; i++ {
           for j := 0; j < V; j++ {
               if dist[i][k] + dist[k][j] < dist[i][j] {
                   dist[i][j] = dist[i][k] + dist[k][j]
               }
           }
       }
   }
   // consider target parity fMask with <=2 bits
   best := int64(0)
   for fMask := 0; fMask < (1 << V); fMask++ {
       bc := bitsCount(fMask)
       if bc != 0 && bc != 2 {
           continue
       }
       xMask := bMask ^ fMask
       // T-join nodes: bits in xMask
       tcnt := bitsCount(xMask)
       var cost int64
       switch tcnt {
       case 0:
           cost = 0
       case 2:
           // find the two nodes
           var u, v int
           for i := 0; i < V; i++ {
               if xMask&(1<<i) != 0 {
                   if u == 0 && v == 0 {
                       u = i + 1
                   } else {
                       v = i + 1
                   }
               }
           }
           // adjust to zero-based
           u--
           v--
           cost = dist[u][v]
       case 4:
           // all four nodes in T
           // pairings: (0,1)+(2,3), (0,2)+(1,3), (0,3)+(1,2)
           c1 := dist[0][1] + dist[2][3]
           c2 := dist[0][2] + dist[1][3]
           c3 := dist[0][3] + dist[1][2]
           cost = c1
           if c2 < cost {
               cost = c2
           }
           if c3 < cost {
               cost = c3
           }
       default:
           continue
       }
       if cost >= INF/2 {
           continue
       }
       cur := total - cost
       if cur > best {
           best = cur
       }
   }
   // Edge case: best could remain 0, but at least one block has positive value
   // best covers single-block sequences too since with one block, parity of used is two ones -> fMask should be parity matches.
   fmt.Println(best)
}

// bitsCount returns number of set bits in mask
func bitsCount(x int) int {
   // builtin PopCount for int64; here x small
   cnt := 0
   for x != 0 {
       cnt++
       x &= x - 1
   }
   return cnt
}
