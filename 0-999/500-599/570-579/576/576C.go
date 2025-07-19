package main

import (
   "bufio"
   "io"
   "os"
   "sort"
   "strconv"
)

type pt struct {
   x, y, idx int
}

// readInt reads the next integer from the reader.
func readInt(r *bufio.Reader) (int, error) {
   var c byte
   var err error
   // skip non-numeric characters
   for {
       if c, err = r.ReadByte(); err != nil {
           return 0, err
       }
       if (c >= '0' && c <= '9') || c == '-' {
           break
       }
   }
   neg := false
   if c == '-' {
       neg = true
       if c, err = r.ReadByte(); err != nil {
           return 0, err
       }
   }
   x := 0
   for ; c >= '0' && c <= '9'; {
       x = x*10 + int(c-'0')
       if c, err = r.ReadByte(); err != nil {
           if err == io.EOF {
               break
           }
           return 0, err
       }
   }
   if neg {
       return -x, nil
   }
   return x, nil
}

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   n, err := readInt(reader)
   if err != nil {
       return
   }
   pts := make([]pt, n)
   for i := 0; i < n; i++ {
       a, _ := readInt(reader)
       b, _ := readInt(reader)
       pts[i] = pt{x: a, y: b, idx: i + 1}
   }
   sort.Slice(pts, func(i, j int) bool {
       bi := pts[i].x / 1000
       bj := pts[j].x / 1000
       if bi != bj {
           return bi < bj
       }
       if bi%2 == 1 {
           return pts[i].y < pts[j].y
       }
       return pts[i].y > pts[j].y
   })
   for i, p := range pts {
       writer.WriteString(strconv.Itoa(p.idx))
       if i+1 < n {
           writer.WriteByte(' ')
       }
   }
   writer.WriteByte('\n')
}
