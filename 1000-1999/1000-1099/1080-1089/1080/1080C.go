package main

import (
   "bufio"
   "fmt"
   "os"
)

var reader = bufio.NewReader(os.Stdin)
var writer = bufio.NewWriter(os.Stdout)

func fg(ax, ay, bx, by int64) (int64, int64, bool) {
   if ay < bx || by < ax {
       return -1, -1, false
   }
   if ax <= bx && bx <= ay && ay <= by {
       return bx, ay, true
   }
   if ax <= bx && by <= ay {
       return bx, by, true
   }
   if bx <= ax && ay <= by {
       return ax, ay, true
   }
   if bx <= ax && ax <= by && by <= ay {
       return ax, by, true
   }
   return -1, -1, false
}

// gx computes number of black and white cells in rectangle [x..z] x [y..w]
func gx(x, y, z, w int64) (black, white int64) {
   h := z - x + 1
   l := w - y + 1
   black = h * l / 2
   white = h * l / 2
   if (h&1) == 1 && (l&1) == 1 {
       if (x+y)&1 == 1 {
           black++
       } else {
           white++
       }
   }
   return
}

func S(x, y, z, w int64) int64 {
   return (z - x + 1) * (w - y + 1)
}

func main() {
   defer writer.Flush()
   var T int
   fmt.Fscan(reader, &T)
   for T > 0 {
       T--
       var n, m int64
       fmt.Fscan(reader, &n, &m)
       bla, whi := gx(1, 1, n, m)
       var ax1, ay1, ax2, ay2 int64
       var bx1, by1, bx2, by2 int64
       fmt.Fscan(reader, &ax1, &ay1, &ax2, &ay2)
       fmt.Fscan(reader, &bx1, &by1, &bx2, &by2)
       cx1, cx2, okx := fg(ax1, ax2, bx1, bx2)
       cy1, cy2, oky := fg(ay1, ay2, by1, by2)
       abla, awhi := gx(ax1, ay1, ax2, ay2)
       bbla, bwhi := gx(bx1, by1, bx2, by2)
       bla -= abla
       whi -= awhi
       bla -= bbla
       whi -= bwhi
       if !okx || !oky {
           bla += S(bx1, by1, bx2, by2)
           whi += S(ax1, ay1, ax2, ay2)
           fmt.Fprintf(writer, "%d %d\n", whi, bla)
           continue
       }
       cbla, cwhi := gx(cx1, cy1, cx2, cy2)
       bla += cbla
       whi += cwhi
       bla += S(bx1, by1, bx2, by2)
       whi += S(ax1, ay1, ax2, ay2) - S(cx1, cy1, cx2, cy2)
       fmt.Fprintf(writer, "%d %d\n", whi, bla)
   }
}
