package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

const BL = 550

type triple struct {
   nx, qid, idx int
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   var t int
   fmt.Fscan(reader, &t)
   for i := 0; i < t; i++ {
       solve(reader, writer)
   }
}

func solve(reader *bufio.Reader, writer *bufio.Writer) {
   var n, q int
   fmt.Fscan(reader, &n)
   fmt.Fscan(reader, &q)

   a := make([][]int, q)
   for i := 0; i < q; i++ {
       var k int
       fmt.Fscan(reader, &k)
       a[i] = make([]int, k)
       for j := 0; j < k; j++ {
           fmt.Fscan(reader, &a[i][j])
           a[i][j]--
       }
   }

   sm := make([][]triple, n)
   id := make([]int, n)

   for i := 0; i < q; i++ {
       if len(a[i]) > BL {
           for j := 0; j < n; j++ {
               id[j] = -1
           }
           for j, x := range a[i] {
               id[x] = j
           }
           for j := 0; j < q; j++ {
               if j == i {
                   continue
               }
               last := n
               for _, x := range a[j] {
                   y := id[x]
                   if y != -1 {
                       if last < y-1 {
                           fmt.Fprintln(writer, "Human")
                           return
                       }
                       last = y
                   }
               }
           }
       } else {
           for j := 0; j+1 < len(a[i]); j++ {
               x := a[i][j]
               nx := a[i][j+1]
               sm[x] = append(sm[x], triple{nx, i, j})
           }
       }
   }

   for j := 0; j < n; j++ {
       id[j] = 0
   }
   cur := 1

   for i := 0; i < n; i++ {
       cur++
       lg, rg := cur, cur
       trips := sm[i]
       sort.Slice(trips, func(p, q int) bool {
           if trips[p].nx != trips[q].nx {
               return trips[p].nx < trips[q].nx
           }
           if trips[p].qid != trips[q].qid {
               return trips[p].qid < trips[q].qid
           }
           return trips[p].idx < trips[q].idx
       })
       prev := -1
       for _, tr := range trips {
           cur++
           if tr.nx != prev {
               rg = cur
               prev = tr.nx
           }
           ca := a[tr.qid]
           for jj := tr.idx + 1; jj < len(ca); jj++ {
               x := ca[jj]
               if lg <= id[x] && id[x] < rg {
                   fmt.Fprintln(writer, "Human")
                   return
               }
               id[x] = cur
           }
       }
   }
   fmt.Fprintln(writer, "Robot")
}
