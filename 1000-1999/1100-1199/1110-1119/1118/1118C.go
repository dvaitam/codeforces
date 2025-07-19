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
   total := n * n
   cnt := make(map[int]int)
   for i := 0; i < total; i++ {
       var x int
       fmt.Fscan(reader, &x)
       cnt[x]++
   }
   k := n / 2
   // collect counts
   vals := make([]int, 0, len(cnt))
   for v := range cnt {
       vals = append(vals, v)
   }
   sort.Ints(vals)
   // initial list of quadruples, pairs, singles
   quadList := make([]int, 0)
   pairList := make([]int, 0)
   singleList := make([]int, 0)
   for _, v := range vals {
       c := cnt[v]
       quads := c / 4
       for i := 0; i < quads; i++ {
           quadList = append(quadList, v)
       }
       rem := c % 4
       if rem/2 > 0 {
           for i := 0; i < rem/2; i++ {
               pairList = append(pairList, v)
           }
       }
       if rem%2 == 1 {
           singleList = append(singleList, v)
       }
   }
   // even n
   if n%2 == 0 {
       needQuads := k * k
       if len(singleList) != 0 || len(pairList) != 0 || len(quadList) < needQuads {
           fmt.Fprintln(writer, "NO")
           return
       }
       // fill matrix
       b := make([][]int, n)
       for i := range b {
           b[i] = make([]int, n)
       }
       idx := 0
       for i := 0; i < k; i++ {
           for j := 0; j < k; j++ {
               v := quadList[idx]
               idx++
               b[i][j] = v
               b[i][n-1-j] = v
               b[n-1-i][j] = v
               b[n-1-i][n-1-j] = v
           }
       }
       fmt.Fprintln(writer, "YES")
       for i := 0; i < n; i++ {
           for j := 0; j < n; j++ {
               fmt.Fprint(writer, b[i][j], " ")
           }
           fmt.Fprintln(writer)
       }
       return
   }
   // odd n
   needQuads := k * k
   if len(singleList) != 1 || len(quadList) < needQuads {
       fmt.Fprintln(writer, "NO")
       return
   }
   // split quads
   usedQuads := quadList[:needQuads]
   leftoverQuads := quadList[needQuads:]
   // augment pairs with leftovers
   for _, v := range leftoverQuads {
       pairList = append(pairList, v, v)
   }
   needPairs := 2 * k
   if len(pairList) < needPairs {
       fmt.Fprintln(writer, "NO")
       return
   }
   // take required pairs
   pairList = pairList[:needPairs]
   // fill matrix
   b := make([][]int, n)
   for i := range b {
       b[i] = make([]int, n)
   }
   // place quadruples
   qidx := 0
   for i := 0; i < k; i++ {
       for j := 0; j < k; j++ {
           v := usedQuads[qidx]
           qidx++
           b[i][j] = v
           b[i][n-1-j] = v
           b[n-1-i][j] = v
           b[n-1-i][n-1-j] = v
       }
   }
   // place pairs: first for center row, then center column
   pidx := 0
   // center row
   for j := 0; j < k; j++ {
       v := pairList[pidx]
       pidx++
       b[k][j] = v
       b[k][n-1-j] = v
   }
   // center column
   for i := 0; i < k; i++ {
       v := pairList[pidx]
       pidx++
       b[i][k] = v
       b[n-1-i][k] = v
   }
   // center cell
   b[k][k] = singleList[0]
   fmt.Fprintln(writer, "YES")
   for i := 0; i < n; i++ {
       for j := 0; j < n; j++ {
           fmt.Fprint(writer, b[i][j], " ")
       }
       fmt.Fprintln(writer)
   }
}
