package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

var (
   n, cnt, curu int
   a      [][]int
   w, fa  []int
   reader *bufio.Reader
   writer *bufio.Writer
)

func readInt() int {
   var x int
   fmt.Fscan(reader, &x)
   return x
}

func solve(vc []int) int {
   if len(vc) == 1 {
       return vc[0]
   }
   u := vc[0]
   curu = u
   sort.Slice(vc, func(i, j int) bool {
       return a[vc[i]][curu] < a[vc[j]][curu]
   })
   cur := make([]int, 0, len(vc))
   top := u
   for i := 1; i < len(vc); i++ {
       v := vc[i]
       if a[u][v] != a[u][vc[i-1]] {
           if i != 1 {
               g := solve(cur)
               expw := a[u][vc[i-1]]
               if w[g] == expw {
                   fa[top] = g
                   top = g
               } else {
                   cnt++
                   w[cnt] = expw
                   fa[top] = cnt
                   fa[g] = cnt
                   top = cnt
               }
               cur = cur[:0]
           }
       }
       cur = append(cur, v)
   }
   // process last group
   g := solve(cur)
   expw := a[u][vc[len(vc)-1]]
   if w[g] == expw {
       fa[top] = g
       top = g
   } else {
       cnt++
       w[cnt] = expw
       fa[top] = cnt
       fa[g] = cnt
       top = cnt
   }
   return top
}

func main() {
   reader = bufio.NewReader(os.Stdin)
   writer = bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   n = readInt()
   // initialize
   a = make([][]int, n+1)
   for i := 1; i <= n; i++ {
       a[i] = make([]int, n+1)
   }
   for i := 1; i <= n; i++ {
       for j := 1; j <= n; j++ {
           a[i][j] = readInt()
       }
   }
   // prepare
   w = make([]int, 2*n+5)
   fa = make([]int, 2*n+5)
   vc := make([]int, 0, n)
   for i := 1; i <= n; i++ {
       w[i] = a[i][i]
       vc = append(vc, i)
   }
   cnt = n
   root := solve(vc)
   // output
   fmt.Fprintln(writer, cnt)
   for i := 1; i <= cnt; i++ {
       if i > 1 {
           writer.WriteByte(' ')
       }
       fmt.Fprint(writer, w[i])
   }
   writer.WriteByte('\n')
   fmt.Fprintln(writer, root)
   for i := 1; i <= cnt; i++ {
       if fa[i] != 0 {
           fmt.Fprintln(writer, i, fa[i])
       }
   }
}
