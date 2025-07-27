package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

func main() {
   in := bufio.NewReader(os.Stdin)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()

   var t int
   fmt.Fscan(in, &t)
   for ; t > 0; t-- {
       var n, k int
       var s string
       fmt.Fscan(in, &n, &k)
       fmt.Fscan(in, &s)
       // initial score
       score := 0
       wCount := 0
       prevW := false
       for i := 0; i < n; i++ {
           if s[i] == 'W' {
               wCount++
               score += 1
               if prevW {
                   score += 1
               }
               prevW = true
           } else {
               prevW = false
           }
       }
       if wCount == 0 {
           // no initial wins
           if k == 0 {
               fmt.Fprintln(out, 0)
           } else {
               // make k wins in a row: score = k + (k-1)
               if k > n {
                   k = n
               }
               fmt.Fprintln(out, 2*k-1)
           }
           continue
       }
       // collect interior gaps between W's
       gaps := make([]int, 0)
       // find first and last W positions
       firstW, lastW := -1, -1
       for i := 0; i < n; i++ {
           if s[i] == 'W' {
               if firstW < 0 {
                   firstW = i
               }
               lastW = i
           }
       }
       // interior gaps
       cnt := 0
       for i := firstW; i <= lastW; i++ {
           if s[i] == 'L' {
               cnt++
           } else {
               if cnt > 0 {
                   gaps = append(gaps, cnt)
                   cnt = 0
               }
           }
       }
       // sort and fill interior gaps
       sort.Ints(gaps)
       for _, g := range gaps {
           if k >= g {
               // fill entire gap
               k -= g
               // base +g, adjacency +(g+1)
               score += 2*g + 1
           } else {
               // partial fill
               score += 2 * k
               k = 0
               break
           }
       }
       if k > 0 {
           // leading Ls
           lead := 0
           for i := 0; i < firstW && k > 0; i++ {
               lead++
           }
           // trailing Ls
           trail := 0
           for i := n - 1; i > lastW && k > 0; i-- {
               trail++
           }
           // total edge slots
           slots := lead + trail
           use := k
           if use > slots {
               use = slots
           }
           score += 2 * use
           // k -= use // not needed further
       }
       fmt.Fprintln(out, score)
   }
}
