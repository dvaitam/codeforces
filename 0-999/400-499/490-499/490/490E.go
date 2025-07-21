package main

import (
   "bufio"
   "fmt"
   "os"
   "strconv"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   writer := bufio.NewWriter(os.Stdout)
   defer writer.Flush()

   var n int
   if _, err := fmt.Fscan(reader, &n); err != nil {
       return
   }
   patterns := make([]string, n)
   for i := 0; i < n; i++ {
       fmt.Fscan(reader, &patterns[i])
   }
   result := make([]int64, n)
   var prev int64 = 0
   for i, pat := range patterns {
       k := len(pat)
       lb := prev + 1
       lbStr := strconv.FormatInt(lb, 10)
       var cur []byte
       var ok bool
       if k < len(lbStr) {
           ok = false
       } else if k > len(lbStr) {
           // minimal fill
           cur = make([]byte, k)
           for j := 0; j < k; j++ {
               if pat[j] == '?' {
                   if j == 0 {
                       cur[j] = '1'
                   } else {
                       cur[j] = '0'
                   }
               } else {
                   cur[j] = pat[j]
               }
           }
           ok = true
       } else {
           // k == len(lbStr)
           lbBytes := []byte(lbStr)
           cur = make([]byte, k)
           ok = dfsFill([]byte(pat), lbBytes, cur, 0, true)
       }
       if !ok {
           fmt.Fprintln(writer, "NO")
           return
       }
       // parse cur
       x, err := strconv.ParseInt(string(cur), 10, 64)
       if err != nil {
           fmt.Fprintln(writer, "NO")
           return
       }
       result[i] = x
       prev = x
   }
   fmt.Fprintln(writer, "YES")
   for _, v := range result {
       fmt.Fprintln(writer, v)
   }
}

// dfsFill tries to fill cur for pattern >= lb
func dfsFill(pat, lb, cur []byte, pos int, tight bool) bool {
   k := len(pat)
   if pos == k {
       return true
   }
   if pat[pos] == '?' {
       // try digits
       for d := 0; d <= 9; d++ {
           if pos == 0 && d == 0 {
               continue
           }
           if tight {
               need := int(lb[pos] - '0')
               if d < need {
                   continue
               }
           }
           cur[pos] = byte('0' + d)
           nextTight := tight && (byte('0'+d) == lb[pos])
           if dfsFill(pat, lb, cur, pos+1, nextTight) {
               return true
           }
       }
   } else {
       d := pat[pos]
       if tight && d < lb[pos] {
           return false
       }
       cur[pos] = d
       nextTight := tight && (d == lb[pos])
       if dfsFill(pat, lb, cur, pos+1, nextTight) {
           return true
       }
   }
   return false
}
