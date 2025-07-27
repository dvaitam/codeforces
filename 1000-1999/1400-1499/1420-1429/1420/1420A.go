package main

import (
   "bufio"
   "io"
   "os"
)

var (
   rdr = bufio.NewReader(os.Stdin)
   wrt = bufio.NewWriter(os.Stdout)
)

// readInt reads the next integer from standard input.
func readInt() (int, error) {
   var b byte
   var err error
   // skip non-numeric bytes
   for {
       b, err = rdr.ReadByte()
       if err != nil {
           return 0, err
       }
       if (b >= '0' && b <= '9') || b == '-' {
           break
       }
   }
   neg := false
   if b == '-' {
       neg = true
       b, err = rdr.ReadByte()
       if err != nil {
           return 0, err
       }
   }
   n := int(b - '0')
   for {
       b, err = rdr.ReadByte()
       if err != nil {
           if err == io.EOF {
               break
           }
           return 0, err
       }
       if b < '0' || b > '9' {
           break
       }
       n = n*10 + int(b-'0')
   }
   if neg {
       n = -n
   }
   return n, nil
}

func main() {
   defer wrt.Flush()
   t, err := readInt()
   if err != nil {
       return
   }
   for ; t > 0; t-- {
       n, _ := readInt()
       prev, _ := readInt()
       ok := false
       for i := 1; i < n; i++ {
           curr, _ := readInt()
           if curr >= prev {
               ok = true
           }
           prev = curr
       }
       if ok {
           wrt.WriteString("YES")
       } else {
           wrt.WriteString("NO")
       }
       wrt.WriteByte('\n')
   }
}
