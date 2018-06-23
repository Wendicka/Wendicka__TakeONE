/*
  libs_math.go
  
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

import "math"

// Just a bunch of 
func wi_math_ff(v *VM,mf func(i float64) float64) (bool,string){
	i:=v.Arg_GetFloat(0)
	r:=mf(i)
	v.Return_Float(r)
	return TBError()
}



func wim_sin(v *VM) (bool,string) { return wi_math_ff(v,math.Sin) }
func wim_cos(v *VM) (bool,string) { return wi_math_ff(v,math.Cos) }
func wim_tan(v *VM) (bool,string) { return wi_math_ff(v,math.Tan) }
func wim_sqr(v *VM) (bool,string) { return wi_math_ff(v,math.Sqrt) }

func wim_round(v *VM) (bool,string) { return wi_math_ff(v,math.Round) }
func wim_gamma(v *VM) (bool,string) { return wi_math_ff(v,math.Gamma) }

func wim_floor(v *VM) (bool,string) { return wi_math_ff(v,math.Floor) }
func wim_ceil (v *VM) (bool,string) { return wi_math_ff(v,math.Ceil)  }

func wim_pi   (v *VM) (bool,string) { v.Return_Float(math.Pi); return true,"" }


func wim_init(w *API){
	w.Register("SIN",wim_sin)
	w.Register("COS",wim_cos)
	w.Register("TAN",wim_tan)
	w.Register("SQR",wim_sqr)
	w.Register("ROUND",wim_round)
	w.Register("GAMMA",wim_gamma)
	w.Register("FLOOR",wim_floor)
	w.Register("CEIL",wim_sin)
	w.Register("PI",wim_pi)
}


