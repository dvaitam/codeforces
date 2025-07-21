package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   in := bufio.NewReader(os.Stdin)
   var r, n int
   if _, err := fmt.Fscan(in, &r, &n); err != nil {
       return
   }
   // initial availability: true = free
   avail := make([][3]bool, r+2)
   for i := 1; i <= r; i++ {
       avail[i][1], avail[i][2] = true, true
   }
   // mark reclaimed and forbidden
   rec := make([][3]bool, r+2)
   for k := 0; k < n; k++ {
       var ri, ci int
       fmt.Fscan(in, &ri, &ci)
       rec[ri][ci] = true
   }
   // apply recs and forbidden
   for i := 1; i <= r; i++ {
       for c := 1; c <= 2; c++ {
           if rec[i][c] {
               // remove reclaimed
               avail[i][c] = false
               // block neighbors in other column
               oc := 3 - c
               for di := -1; di <= 1; di++ {
                   j := i + di
                   if j >= 1 && j <= r {
                       avail[j][oc] = false
                   }
               }
           }
       }
   }
   // precompute Grundy g[len][topBlock][botBlock]
   // topBlock/botBlock: 0=none,1=col1 blocked,2=col2 blocked
   maxr := r
   g := make([][][3]int, maxr+1)
   for i := 0; i <= maxr; i++ {
       g[i] = make([][3]int, 3)
   }
   // length 0: all zero
   for tb := 0; tb < 3; tb++ {
       for bb := 0; bb < 3; bb++ {
           g[0][tb][bb] = 0
       }
   }
   for length := 1; length <= maxr; length++ {
       for tb := 0; tb < 3; tb++ {
           for bb := 0; bb < 3; bb++ {
               used := make([]bool, length*2+5)
               for i := 1; i <= length; i++ {
                   for c := 1; c <= 2; c++ {
                       // check if blocked by boundary
                       if (i == 1 && tb == c) || (i == length && bb == c) {
                           continue
                       }
                       // move at (i,c)
                       // top segment 1..i-1
                       var g1 int
                       if i > 1 {
                           // bottom block for top segment = block at row i-1 from rec at i
                           bb1 := 3 - c
                           g1 = g[i-1][tb][bb1]
                       }
                       // bottom segment i+1..length
                       var g2 int
                       if i < length {
                           tb2 := 3 - c
                           g2 = g[length-i][tb2][bb]
                       }
                       x := g1 ^ g2
                       if x < len(used) {
                           used[x] = true
                       }
                   }
               }
               // mex
               mex := 0
               for mex < len(used) && used[mex] {
                   mex++
               }
               g[length][tb][bb] = mex
           }
       }
   }
   // compute XOR over segments
   total := 0
   // scan rows
   var segStart int
   var segTB int
   inSeg := false
   for i := 1; i <= r; i++ {
       if !avail[i][1] && !avail[i][2] {
           // boundary
           if inSeg {
               // close segment at i-1
               length := i - segStart
               // bottomBlock at row i-1
               bb := 0
               if !avail[i-1][1] {
                   bb = 1
               }
               if !avail[i-1][2] {
                   bb = 2
               }
               total ^= g[length][segTB][bb]
           }
           inSeg = false
       } else {
           if !inSeg {
               // start new
               inSeg = true
               segStart = i
               // topBlock at row i
               segTB = 0
               if !avail[i][1] {
                   segTB = 1
               }
               if !avail[i][2] {
                   segTB = 2
               }
           }
       }
   }
   if inSeg {
       length := r - segStart + 1
       bb := 0
       if !avail[r][1] {
           bb = 1
       }
       if !avail[r][2] {
           bb = 2
       }
       total ^= g[length][segTB][bb]
   }
   if total != 0 {
       fmt.Println("WIN")
   } else {
       fmt.Println("LOSE")
   }
}
