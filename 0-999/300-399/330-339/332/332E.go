package main

import (
   "bufio"
   "fmt"
   "os"
   "strconv"
)

func main() {
   reader := bufio.NewReader(os.Stdin)
   // Read p
   p, err := reader.ReadString('\n')
   if err != nil {
       return
   }
   p = p[:len(p)-1]
   // Read s
   s, err := reader.ReadString('\n')
   if err != nil {
       return
   }
   s = s[:len(s)-1]
   // Read k
   line, err := reader.ReadString('\n')
   if err != nil {
       return
   }
   k, err := strconv.Atoi(string(line[:len(line)-1]))
   if err != nil {
       return
   }
   n := len(p)
   sl := len(s)
   // q: '?': unknown, '0' or '1'
   q := make([]byte, k)
   for i := 0; i < k; i++ {
       q[i] = '?'
   }
   sPos := 0
   ok := true
   for i := 0; i < n && sPos < sl; i++ {
       j := i % k
       if q[j] == '1' {
           if p[i] != s[sPos] {
               ok = false
               break
           }
           sPos++
       } else if q[j] == '0' {
           continue
       } else {
           // unknown
           if p[i] == s[sPos] {
               q[j] = '1'
               sPos++
           } else {
               q[j] = '0'
           }
       }
   }
   if !ok || sPos < sl {
       fmt.Println("0")
       return
   }
   // fill unknowns with '0'
   for i := 0; i < k; i++ {
       if q[i] == '?' {
           q[i] = '0'
       }
   }
   // Verify complete extraction matches s (no extra picks)
   var out []byte
   for i := 0; i < n; i++ {
       if q[i%k] == '1' {
           out = append(out, p[i])
       }
   }
   if len(out) != sl || string(out) != s {
       fmt.Println("0")
       return
   }
   fmt.Println(string(q))
}
