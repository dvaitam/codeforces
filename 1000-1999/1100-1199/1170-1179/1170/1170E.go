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

   var n, m int
   if _, err := fmt.Fscan(reader, &n, &m); err != nil {
       return
   }
   a := make([]int, n)
   pref := make([]int, n+1)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &a[i])
       pref[i+1] = pref[i] + a[i]
   }
   // total slack: empty compartments
   slack := m - pref[n]

   var q int
   fmt.Fscan(reader, &q)
   for qi := 0; qi < q; qi++ {
       var ck int
       fmt.Fscan(reader, &ck)
       ws := make([]int, ck)
       for i := 0; i < ck; i++ {
           fmt.Fscan(reader, &ws[i])
       }
       sumMove := 0
       ok := true
       // process employees grouped by door index
       for i := 0; i < ck; {
           w := ws[i]
           // if in the trailing gap, no door covers
           if w > pref[n] {
               i++
               continue
           }
           // find door j: smallest j such that pref[j] >= w
           j := sort.Search(len(pref), func(k int) bool { return pref[k] >= w })
           // group all ws in this door j
           minW, maxW := w, w
           i++
           for i < ck && ws[i] <= pref[n] {
               // check if same door
               if ws[i] > pref[j] {
                   break
               }
               // ws[i] >= pref[j-1]+1 automatically since ws sorted and >prev w
               maxW = ws[i]
               i++
           }
           // compute movement for door j: left or right
           // left move: move so its end < minW => shift left by pref[j] - (minW-1)
           leftMove := pref[j] - minW + 1
           // right move: move so its start > maxW => shift right by (maxW+1) - (pref[j-1]+1) => maxW - pref[j-1] + 1
           rightMove := maxW - pref[j-1]
           need := leftMove
           if rightMove < need {
               need = rightMove
           }
           sumMove += need
           if sumMove > slack {
               ok = false
               break
           }
       }
       if ok {
           fmt.Fprintln(writer, "YES")
       } else {
           fmt.Fprintln(writer, "NO")
       }
   }
}
