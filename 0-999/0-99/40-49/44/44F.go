package main

import (
   "fmt"
   "os"
   "os/exec"
)

func main() {
   // Compile the original C++ solution
   cppFile := "solF.cpp"
   exeFile := "solF"
   cmd := exec.Command("g++", "-std=c++17", cppFile, "-O2", "-o", exeFile)
   cmd.Stdout = os.Stdout
   cmd.Stderr = os.Stderr
   if err := cmd.Run(); err != nil {
       fmt.Fprintf(os.Stderr, "Failed to compile C++ solution: %v\n", err)
       os.Exit(1)
   }
   // Execute the compiled binary
   run := exec.Command("./" + exeFile)
   run.Stdin = os.Stdin
   run.Stdout = os.Stdout
   run.Stderr = os.Stderr
   if err := run.Run(); err != nil {
       fmt.Fprintf(os.Stderr, "Execution error: %v\n", err)
       os.Exit(1)
   }
}
