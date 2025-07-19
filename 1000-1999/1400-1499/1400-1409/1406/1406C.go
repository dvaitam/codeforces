package main

import (
   "bufio"
   "fmt"
   "os"
)

var (
   reader = bufio.NewReader(os.Stdin)
   writer = bufio.NewWriter(os.Stdout)
)

func readInt() (int, error) {
   var x int
   _, err := fmt.Fscan(reader, &x)
   return x, err
}

func main() {
   defer writer.Flush()
   T, err := readInt()
   if err != nil {
       return
   }
   for t := 0; t < T; t++ {
       n, _ := readInt()
       // build graph
       g := make([][]int, n+1)
       for i := 1; i < n; i++ {
           u, _ := readInt()
           v, _ := readInt()
           g[u] = append(g[u], v)
           g[v] = append(g[v], u)
       }
       // find centroids
       size := make([]int, n+1)
       best := n + 1
       centroids := []int{}
       var dfs1 func(int, int)
       dfs1 = func(u, p int) {
           size[u] = 1
           maxSub := 0
           for _, v := range g[u] {
               if v == p {
                   continue
               }
               dfs1(v, u)
               size[u] += size[v]
               if size[v] > maxSub {
                   maxSub = size[v]
               }
           }
           others := n - size[u]
           if others > maxSub {
               maxSub = others
           }
           if maxSub < best {
               best = maxSub
               centroids = centroids[:0]
           }
           if maxSub == best {
               centroids = append(centroids, u)
           }
       }
       dfs1(1, 0)
       if len(centroids) == 1 {
           c := centroids[0]
           // any neighbor
           u := g[c][0]
           fmt.Fprintf(writer, "%d %d\n", c, u)
           fmt.Fprintf(writer, "%d %d\n", c, u)
       } else {
           c1, c2 := centroids[0], centroids[1]
           // find leaf in subtree of c1 excluding c2
           var findLeaf func(int, int) (leaf, parent int)
           findLeaf = func(u, p int) (int, int) {
               for _, v := range g[u] {
                   if v == p || v == c2 {
                       continue
                   }
                   return findLeaf(v, u)
               }
               return u, p
           }
           leaf, parent := findLeaf(c1, 0)
           fmt.Fprintf(writer, "%d %d\n", leaf, parent)
           fmt.Fprintf(writer, "%d %d\n", leaf, c2)
       }
   }
}
