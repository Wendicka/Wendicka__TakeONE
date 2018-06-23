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


