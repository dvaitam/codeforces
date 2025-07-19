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

   var ys, xs int
   fmt.Fscan(in, &ys, &xs)
   jed := make([]int, ys)
   zer := make([]int, ys)
   chce1 := make([]int, ys)
   chce0 := make([]int, ys)

   type op struct{ y1, x1, y2, x2 int }
   R := make([]op, 0)
   K := make([]op, 0)

   for y := 0; y < ys; y++ {
       for x := 0; x < xs; x++ {
           var buf string
           fmt.Fscan(in, &buf)
           for z := len(buf) - 1; z >= 0; z-- {
               if buf[z] == '1' {
                   if x == 1 {
                       jed[1-y]++
                       R = append(R, op{y, x, 1 - y, x})
                   } else {
                       jed[y]++
                       R = append(R, op{y, x, y, 1})
                   }
               } else {
                   if x == 0 {
                       zer[1-y]++
                       R = append(R, op{y, x, 1 - y, x})
                   } else {
                       zer[y]++
                       R = append(R, op{y, x, y, 0})
                   }
               }
           }
       }
   }
   for y := 0; y < ys; y++ {
       for x := 0; x < xs; x++ {
           var buf string
           fmt.Fscan(in, &buf)
           for z := len(buf) - 1; z >= 0; z-- {
               if buf[z] == '1' {
                   if x == 1 {
                       chce1[1-y]++
                       K = append(K, op{1 - y, x, y, x})
                   } else {
                       chce1[y]++
                       K = append(K, op{y, 1, y, x})
                   }
               } else {
                   if x == 0 {
                       chce0[1-y]++
                       K = append(K, op{1 - y, x, y, x})
                   } else {
                       chce0[y]++
                       K = append(K, op{y, 0, y, x})
                   }
               }
           }
       }
   }
   // balance ones
   for ym, yd := 0, 0; ; {
       for ym < ys && jed[ym] >= chce1[ym] {
           ym++
       }
       if ym == ys {
           break
       }
       for jed[yd] <= chce1[yd] {
           yd++
       }
       jed[yd]--
       jed[ym]++
       R = append(R, op{yd, 1, ym, 1})
   }
   // balance zeros
   for ym, yd := 0, 0; ; {
       for ym < ys && zer[ym] >= chce0[ym] {
           ym++
       }
       if ym == ys {
           break
       }
       for zer[yd] <= chce0[yd] {
           yd++
       }
       zer[yd]--
       zer[ym]++
       R = append(R, op{yd, 0, ym, 0})
   }
   total := len(R) + len(K)
   fmt.Fprintln(out, total)
   for _, v := range R {
       // convert to 1-based
       fmt.Fprintf(out, "%d %d %d %d\n", v.y1+1, v.x1+1, v.y2+1, v.x2+1)
   }
   for _, v := range K {
       fmt.Fprintf(out, "%d %d %d %d\n", v.y1+1, v.x1+1, v.y2+1, v.x2+1)
   }
}
