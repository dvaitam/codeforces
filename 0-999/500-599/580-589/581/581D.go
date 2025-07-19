package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   in := bufio.NewReader(os.Stdin)
   out := bufio.NewWriter(os.Stdout)
   defer out.Flush()
   var h0, w0, h1, w1, h2, w2 int
   if _, err := fmt.Fscan(in, &h0, &w0, &h1, &w1, &h2, &w2); err != nil {
       return
   }
   dims := [3][2]int{{h0, w0}, {h1, w1}, {h2, w2}}
   labels := [3]byte{'A', 'B', 'C'}
   // total area and side length
   area := dims[0][0]*dims[0][1] + dims[1][0]*dims[1][1] + dims[2][0]*dims[2][1]
   mx := dims[0][0]
   for i := 0; i < 3; i++ {
       if dims[i][0] > mx {
           mx = dims[i][0]
       }
       if dims[i][1] > mx {
           mx = dims[i][1]
       }
   }
   if mx*mx != area {
       fmt.Fprintln(out, -1)
       return
   }
   // try placements
   for pi := 1; pi <= 3; pi++ {
       // apply permutation swaps
       if pi == 2 {
           dims[1], dims[2] = dims[2], dims[1]
           labels[1], labels[2] = labels[2], labels[1]
       }
       if pi == 3 {
           dims[0], dims[2] = dims[2], dims[0]
           labels[0], labels[2] = labels[2], labels[0]
       }
       // rotations
       for t := 1; t <= 8; t++ {
           // pattern 1: two on bottom, one on top
           if dims[0][0] == dims[1][0] && dims[0][1]+dims[1][1] == mx && dims[2][0]+dims[0][0] == mx && dims[2][1] == mx {
               // output
               fmt.Fprintln(out, mx)
               // top block
               for i := 0; i < dims[2][0]; i++ {
                   for j := 0; j < mx; j++ {
                       out.WriteByte(labels[2])
                   }
                   out.WriteByte('\n')
               }
               // bottom blocks
               for i := 0; i < dims[0][0]; i++ {
                   for j := 0; j < dims[0][1]; j++ {
                       out.WriteByte(labels[0])
                   }
                   for j := 0; j < dims[1][1]; j++ {
                       out.WriteByte(labels[1])
                   }
                   out.WriteByte('\n')
               }
               out.Flush()
               os.Exit(0)
           }
           // pattern 2: three stacked vertically
           if dims[0][0]+dims[1][0]+dims[2][0] == mx && dims[0][1] == mx && dims[1][1] == mx && dims[2][1] == mx {
               fmt.Fprintln(out, mx)
               for i := 0; i < mx; i++ {
                   var label byte
                   if i < dims[0][0] {
                       label = labels[0]
                   } else if i < dims[0][0]+dims[1][0] {
                       label = labels[1]
                   } else {
                       label = labels[2]
                   }
                   for j := 0; j < mx; j++ {
                       out.WriteByte(label)
                   }
                   out.WriteByte('\n')
               }
               out.Flush()
               os.Exit(0)
           }
           // apply rotation swaps for next
           if t%2 == 1 {
               dims[2][0], dims[2][1] = dims[2][1], dims[2][0]
           }
           if t%4 == 2 {
               dims[1][0], dims[1][1] = dims[1][1], dims[1][0]
           }
           if t%8 == 4 {
               dims[0][0], dims[0][1] = dims[0][1], dims[0][0]
           }
       }
   }
   fmt.Fprintln(out, -1)
}
