package main

import (
   "bufio"
   "fmt"
   "io"
   "os"
   "sort"
)

type segment struct {
   x, y, id int
}

var reader = bufio.NewReader(os.Stdin)
var writer = bufio.NewWriter(os.Stdout)

func readInt() int {
   var c byte
   var err error
   // skip non-numeric characters
   for {
       c, err = reader.ReadByte()
       if err != nil {
           if err == io.EOF {
               return 0
           }
           panic(err)
       }
       if (c >= '0' && c <= '9') || c == '-' {
           break
       }
   }
   sign := 1
   if c == '-' {
       sign = -1
       c, _ = reader.ReadByte()
   }
   x := int(c - '0')
   for {
       c, err = reader.ReadByte()
       if err != nil || c < '0' || c > '9' {
           break
       }
       x = x*10 + int(c-'0')
   }
   return x * sign
}

func main() {
   defer writer.Flush()
   T := readInt()
   for ; T > 0; T-- {
       n := readInt()
       segs := make([]segment, n)
       for i := 0; i < n; i++ {
           l := readInt()
           r := readInt()
           segs[i] = segment{l, r, i}
       }
       sort.Slice(segs, func(i, j int) bool {
           if segs[i].x != segs[j].x {
               return segs[i].x < segs[j].x
           }
           return segs[i].y < segs[j].y
       })
       // count clusters
       clusters := 1
       curR := segs[0].y
       for i := 1; i < n; i++ {
           if segs[i].x <= curR {
               if segs[i].y > curR {
                   curR = segs[i].y
               }
           } else {
               clusters++
               if segs[i].y > curR {
                   curR = segs[i].y
               }
           }
       }
       if clusters < 2 {
           writer.WriteString("-1\n")
           continue
       }
       // assign groups
       out := make([]int, n)
       group := 1
       curR = segs[0].y
       out[segs[0].id] = group
       for i := 1; i < n; i++ {
           if segs[i].x <= curR {
               out[segs[i].id] = group
               if segs[i].y > curR {
                   curR = segs[i].y
               }
           } else {
               group++
               if group > 2 {
                   group = 2
               }
               out[segs[i].id] = group
               if segs[i].y > curR {
                   curR = segs[i].y
               }
           }
       }
       // output
       for i := 0; i < n; i++ {
           fmt.Fprintf(writer, "%d", out[i])
           if i+1 < n {
               writer.WriteByte(' ')
           }
       }
       writer.WriteByte('\n')
   }
}
