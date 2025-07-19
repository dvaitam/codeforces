package main

import (
   "bufio"
   "fmt"
   "os"
   "sort"
)

type bomb struct {
   x, y, d int
}

func abs(a int) int {
   if a < 0 {
       return -a
   }
   return a
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n int
   fmt.Fscan(reader, &n)
   bs := make([]bomb, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &bs[i].x, &bs[i].y)
       bs[i].d = abs(bs[i].x) + abs(bs[i].y)
   }
   sort.Slice(bs, func(i, j int) bool {
       return bs[i].d < bs[j].d
   })
   // count operations
   var k int
   for _, b := range bs {
       if b.x != 0 {
           k++
       }
       if b.y != 0 {
           k++
       }
       // pick
       k++
       // return moves
       if b.y != 0 {
           k++
       }
       if b.x != 0 {
           k++
       }
       // destroy
       k++
   }
   fmt.Fprintln(writer, k)
   // output operations
   for _, b := range bs {
       // go to bomb
       if b.x > 0 {
           fmt.Fprintf(writer, "1 %d R\n", b.x)
       } else if b.x < 0 {
           fmt.Fprintf(writer, "1 %d L\n", -b.x)
       }
       if b.y > 0 {
           fmt.Fprintf(writer, "1 %d U\n", b.y)
       } else if b.y < 0 {
           fmt.Fprintf(writer, "1 %d D\n", -b.y)
       }
       // pick
       fmt.Fprintln(writer, 2)
       // return to origin
       if b.y > 0 {
           fmt.Fprintf(writer, "1 %d D\n", b.y)
       } else if b.y < 0 {
           fmt.Fprintf(writer, "1 %d U\n", -b.y)
       }
       if b.x > 0 {
           fmt.Fprintf(writer, "1 %d L\n", b.x)
       } else if b.x < 0 {
           fmt.Fprintf(writer, "1 %d R\n", -b.x)
       }
       // destroy
       fmt.Fprintln(writer, 3)
   }
}
