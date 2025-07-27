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
       // 1-based array with sentinel at 0
       a := make([]int, n+2)
       a[0] = 0
       for i := 1; i <= n; i++ {
           ai, _ := readInt()
           a[i] = ai
       }
       // initial answer: sum of positive deltas
       var ans int64
       for i := 1; i <= n; i++ {
           delta := a[i] - a[i-1]
           if delta > 0 {
               ans += int64(delta)
           }
       }
       // print initial
       wrt.WriteString(strconv.FormatInt(ans, 10))
       wrt.WriteByte('\n')
       // process swaps
       for i := 0; i < q; i++ {
           l, _ := readInt()
           r, _ := readInt()
           // indices affecting deltas
           idxs := [4]int{l, l + 1, r, r + 1}
           // remove old contributions
           for j := 0; j < 4; j++ {
               idx := idxs[j]
               if idx < 1 || idx > n {
                   continue
               }
               // skip duplicates
               dup := false
               for k := 0; k < j; k++ {
                   if idxs[k] == idx {
                       dup = true
                       break
                   }
               }
               if dup {
                   continue
               }
               d := a[idx] - a[idx-1]
               if d > 0 {
                   ans -= int64(d)
               }
           }
           // perform swap
           a[l], a[r] = a[r], a[l]
           // add new contributions
           for j := 0; j < 4; j++ {
               idx := idxs[j]
               if idx < 1 || idx > n {
                   continue
               }
               dup := false
               for k := 0; k < j; k++ {
                   if idxs[k] == idx {
                       dup = true
                       break
                   }
               }
               if dup {
                   continue
               }
               d := a[idx] - a[idx-1]
               if d > 0 {
                   ans += int64(d)
               }
           }
           // output current answer
           wrt.WriteString(strconv.FormatInt(ans, 10))
           wrt.WriteByte('\n')
       }
   }
}
