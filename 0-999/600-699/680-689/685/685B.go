package main

import (
   "bufio"
   "os"
   "strconv"
)

var (
   reader = bufio.NewReader(os.Stdin)
   writer = bufio.NewWriter(os.Stdout)
)

func readInt() int {
   x := 0
   sign := 1
   b, err := reader.ReadByte()
   for err == nil && (b < '0' || b > '9') && b != '-' {
       b, err = reader.ReadByte()
   }
   if err != nil {
       return 0
   }
   if b == '-' {
       sign = -1
       b, _ = reader.ReadByte()
   }
   for ; err == nil && b >= '0' && b <= '9'; b, err = reader.ReadByte() {
       x = x*10 + int(b-'0')
   }
   return x * sign
}

func main() {
   defer writer.Flush()
   n := readInt()
   q := readInt()
   children := make([][]int, n+1)
   for i := 2; i <= n; i++ {
       p := readInt()
       children[p] = append(children[p], i)
   }
   size := make([]int, n+1)
   hson := make([]int, n+1)
   centroid := make([]int, n+1)
   fa := make([]int, n+1)

   // get maximum subtree size when cutting at x in tree rooted at u
   getMax := func(u, x int) int {
       a := size[u] - size[x]
       b := 0
       if hson[x] != 0 {
           b = size[hson[x]]
       }
       if a > b {
           return a
       }
       return b
   }

   var dfs func(u int)
   dfs = func(u int) {
       size[u] = 1
       hson[u] = 0
       for _, v := range children[u] {
           fa[v] = u
           dfs(v)
           size[u] += size[v]
           if size[v] > size[hson[u]] {
               hson[u] = v
           }
       }
       // initialize centroid as itself
       centroid[u] = u
       if hson[u] != 0 {
           p := centroid[hson[u]]
           // climb up while parent can be a better centroid
           for fa[p] != u && getMax(u, fa[p]) <= getMax(u, p) {
               p = fa[p]
           }
           if size[hson[u]] > getMax(u, p) {
               centroid[u] = p
           }
       }
   }

   dfs(1)
   for i := 0; i < q; i++ {
       u := readInt()
       writer.WriteString(strconv.Itoa(centroid[u]))
       writer.WriteByte('\n')
   }
}
