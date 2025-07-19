package main

import (
   "bufio"
   "fmt"
   "os"
)

var rdr = bufio.NewReader(os.Stdin)

// kll reads an int64 from standard input
func kll() int64 {
   var n int64
   fmt.Fscan(rdr, &n)
   return n
}

// murad processes one test case
func murad() {
   // read input (n), unused in this solution
   _ = kll()
   // output result
   fmt.Println(0)
}

func main() {
   // single test case
   murad()
}
