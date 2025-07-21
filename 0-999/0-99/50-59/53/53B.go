package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   var h, w int64
   if _, err := fmt.Fscan(reader, &h, &w); err != nil {
       return
   }
   var bestH, bestW int64
   var bestArea int64

   // Try fixing h' as power of 2
   for x := 0; x < 62; x++ {
       h0 := int64(1) << x
       if h0 > h {
           break
       }
       // w' constraints: ceil(4*h0/5) <= w' <= min(floor(5*h0/4), w)
       wMin := (4*h0 + 5 - 1) / 5
       wMax1 := (5 * h0) / 4
       if wMin < 1 {
           wMin = 1
       }
       wMax := wMax1
       if wMax > w {
           wMax = w
       }
       if wMin > wMax {
           continue
       }
       w0 := wMax
       area := h0 * w0
       if area > bestArea || (area == bestArea && h0 > bestH) {
           bestArea = area
           bestH = h0
           bestW = w0
       }
   }
   // Try fixing w' as power of 2
   for x := 0; x < 62; x++ {
       w0 := int64(1) << x
       if w0 > w {
           break
       }
       // h' constraints: ceil(4*w0/5) <= h' <= min(floor(5*w0/4), h)
       hMin := (4*w0 + 5 - 1) / 5
       hMax1 := (5 * w0) / 4
       if hMin < 1 {
           hMin = 1
       }
       hMax := hMax1
       if hMax > h {
           hMax = h
       }
       if hMin > hMax {
           continue
       }
       h0 := hMax
       area := h0 * w0
       if area > bestArea || (area == bestArea && h0 > bestH) {
           bestArea = area
           bestH = h0
           bestW = w0
       }
   }
   // Output result
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()
   fmt.Fprintf(writer, "%d %d", bestH, bestW)
}
