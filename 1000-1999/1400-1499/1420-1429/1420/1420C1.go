package main

import (
   "bufio"
   "io"
   "os"
   "strconv"
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
       q, _ := readInt()
       var ans int64
       prev := 0
       for i := 0; i < n; i++ {
           ai, _ := readInt()
           diff := ai - prev
           if diff > 0 {
               ans += int64(diff)
           }
           prev = ai
       }
       // output initial answer
       wrt.WriteString(strconv.FormatInt(ans, 10))
       wrt.WriteByte('\n')
       // skip swap operations (q = 0 in this version)
       for i := 0; i < q; i++ {
           // read and ignore l, r
           _, _ = readInt()
           _, _ = readInt()
           wrt.WriteString(strconv.FormatInt(ans, 10))
           wrt.WriteByte('\n')
       }
   }
}
