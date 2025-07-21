package main

import (
   "bufio"
   "fmt"
   "os"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var w, h int
   if _, err := fmt.Fscan(reader, &w, &h); err != nil {
       return
   }
   // Read original image: h rows of length w
   img := make([][]byte, h)
   for i := 0; i < h; i++ {
       var line string
       fmt.Fscan(reader, &line)
       img[i] = []byte(line)
   }
   // After rotate+flip => transpose: C has size rows=w, cols=h, C[r][c] = img[c][r]
   // Then zoom by 2: output size rows=2*w, cols=2*h
   outRows := 2 * w
   outCols := 2 * h
   buf := make([]byte, outCols)
   for r2 := 0; r2 < outRows; r2++ {
       r := r2 / 2
       // build row
       for c2 := 0; c2 < outCols; c2++ {
           c := c2 / 2
           // from transposed: img[c][r]
           buf[c2] = img[c][r]
       }
       writer.Write(buf)
       writer.WriteByte('\n')
   }
}
