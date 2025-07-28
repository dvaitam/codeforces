package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n int
   if _, err := fmt.Fscan(reader, &n); err != nil {
       return
   }
   p := make([]int, n+1)
   m := make([]int64, n+1)
   // Node 0 is king
   children := make([][]int, n+1)
   for i := 1; i <= n; i++ {
       var oi int
       var mi int64
       fmt.Fscan(reader, &oi, &mi)
       p[i] = oi
       m[i] = mi
       // ensure children slice has oi
       children[oi] = append(children[oi], i)
   }
   // Compute subtree sums T
   T := make([]int64, n+1)
   for i := 1; i <= n; i++ {
       T[i] = m[i]
   }
   // Input guarantee o_i < i => children indices > parent
   for i := n; i >= 1; i-- {
       T[p[i]] += T[i]
   }
   // Build for each node u the sorted list of child subtree sums
   V := make([][]int64, n+1)
   PS := make([][]int64, n+1)
   for u := 0; u <= n; u++ {
       ch := children[u]
       if len(ch) == 0 {
           V[u] = nil
           PS[u] = nil
           continue
       }
       arr := make([]int64, len(ch))
       for i, v := range ch {
           arr[i] = T[v]
       }
       sort.Slice(arr, func(i, j int) bool { return arr[i] < arr[j] })
       V[u] = arr
       ps := make([]int64, len(arr)+1)
       for i := 0; i < len(arr); i++ {
           ps[i+1] = ps[i] + arr[i]
       }
       PS[u] = ps
   }
   // For each node, compute minimal X
   // Results
   res := make([]int64, n+1)
   for i := 1; i <= n; i++ {
       cur := T[i]
       u := p[i]
       for {
           // apply f_u(cur)
           arr := V[u]
           if len(arr) > 0 {
               // find first idx where arr[idx] >= cur
               idx := sort.Search(len(arr), func(j int) bool { return arr[j] >= cur })
               sumSmall := PS[u][idx]
               numLarge := int64(len(arr) - idx)
               cur = sumSmall + numLarge*cur
           }
           if u == 0 {
               break
           }
           u = p[u]
       }
       res[i] = cur
   }
   // Output
   for i := 1; i <= n; i++ {
       if i > 1 {
           writer.WriteByte(' ')
       }
       writer.WriteString(fmt.Sprintf("%d", res[i]))
   }
   writer.WriteByte('\n')
}
