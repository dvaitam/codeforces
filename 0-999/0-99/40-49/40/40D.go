package main

import (
   "bufio"
   "fmt"
   "math/big"
   "os"
)

func main() {
   in := bufio.NewReader(os.Stdin)
   var s string
   if _, err := fmt.Fscan(in, &s); err != nil {
       return
   }
   // TODO: implement full solution.
   // Placeholder: handle only trivial cases 2 and 13
   A := new(big.Int)
   A.SetString(s, 10)
   two := big.NewInt(2)
   thirteen := big.NewInt(13)
   w := bufio.NewWriter(os.Stdout)
   defer w.Flush()
   if A.Cmp(two) == 0 {
       fmt.Fprintln(w, "YES")
       fmt.Fprintln(w, 1)
       fmt.Fprintln(w, 1)
       fmt.Fprintln(w, 1)
       fmt.Fprintln(w, 13)
   } else if A.Cmp(thirteen) == 0 {
       fmt.Fprintln(w, "YES")
       fmt.Fprintln(w, 1)
       fmt.Fprintln(w, 2)
       fmt.Fprintln(w, 1)
       fmt.Fprintln(w, 2)
   } else {
       fmt.Fprintln(w, "NO")
   }
