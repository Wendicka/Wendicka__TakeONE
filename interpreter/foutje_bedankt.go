/*
  foutje_bedankt.go
  
  version: 18.06.23
  Copyright (C) 2018 Jeroen P. Broks
  This software is provided 'as-is', without any express or implied
  warranty.  In no event will the authors be held liable for any damages
  arising from the use of this software.
  Permission is granted to anyone to use this software for any purpose,
  including commercial applications, and to alter it and redistribute it
  freely, subject to the following restrictions:
  1. The origin of this software must not be misrepresented; you must not
     claim that you wrote the original software. If you use this software
     in a product, an acknowledgment in the product documentation would be
     appreciated but is not required.
  2. Altered source versions must be plainly marked as such, and must not be
     misrepresented as being the original software.
  3. This notice may not be removed or altered from any source distribution.
*/
package wi

import (
  "os"
  "fmt"
)

// Constains the last error as a string.
// This value resets every time a new call to a Wendicka VM is made
// If no errors occurs, this is just an empty string.
var Werror string


// What to do when an error occurs in a Wendicka script.
// panic   = Will use the Go Panic routine to crash out the program. (Most dirty and undesirable method)
// crash   = Error message will show, and the program will end
// warn    = Error message will show, but your program will continue
// nothing = Will show nothing
// both "panic" and "crash" will end your program.
// both "warn" and "nothing" will (mostly) terminate the script execution and store the error in WError
var Werhand = "panic"

type thwerror struct {
    f func(er string)
}

var hwerror = map[string] thwerror{}

func init(){
  hwerror["panic"] = thwerror{
    func(er string) {
      panic ( "ERROR!\n"+er )
      os.Exit(1)
    },
  }
  hwerror["crash"] = thwerror{
    func (er string) {
      fmt.Println("ERROR!\n"+er)
      os.Exit(1)
    },
  }
  hwerror["warn"] = thwerror{
    func (er string) {
      fmt.Println("ERROR!\n"+er)
      Werror = er
    },
  }
  hwerror["nothing"] = thwerror{
    func (er string) {
      Werror = er
    },
  }

}

func hError() string{
	return Werror
}

func TBError() (bool,string){
	return Werror=="",Werror
}

func wError(e string){
  hwerror[Werhand].f(e)
}
