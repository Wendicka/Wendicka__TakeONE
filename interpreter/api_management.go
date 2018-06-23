package wi

import (
  "fmt"
  "trickyunits/qstr"
  "strconv"
  "strings"
)

// Args
func (self *VM) Arg_Count() int{
  /*
  i:=0
  for _,ok:=self.identifiers.i[fmt.Sprintf("$__ARG%d",i)];ok;ok=ok {
    i++;
  }
  return i
  */
  return len(self.calls[self.ccall].params.i)
}

func (self *VM) ID_ConvString(id *tIdentifier) string{
  /*
  ak:=v //fmt.Sprintf("$__ARG%d",i)
  id,ok:=self.identifiers.i[ak]
  if !ok { return "nil" }
  */
  switch id.itype {
  case "nil":
    return "nil"
  case "string":
    return id.vstring
  case "int","integer":
    return fmt.Sprintf("%d",id.vint)
  case "bool","boolean":
    if id.vint<=0 { return "false"; } else { return "true"; }
  case "float":
      return fmt.Sprintf("%f",id.vfloat)
  case "chunk","function","procedure","api":
    return "<<"+id.itype+":"+id.vstring+">>"
  case "table","array":
    return "<<"+id.itype+":"+fmt.Sprintf("%d",id.vint)+">>"  
  default:
    return "Unknown type! ("+id.itype+")"
  }
}

func (self *VM) Arg_ConvString(i int) string{
  //ak:=fmt.Sprintf("$__ARG%d",i)
  ak:=self.calls[self.ccall].params.i[i]
  return self.ID_ConvString(ak)
}

func (self *VM) Var_ConvString(v string) string{
  id,ok:=self.identifiers.i[v]
  if !ok { return "nil" } else { return self.ID_ConvString(id)}
}

func (self *VM) Arg_GetString(i int) string {
	ac:=self.Arg_Count()
	if i>=ac { wError(fmt.Sprintf("String expected as argument #%d, while only %d arguments are present",i+1,ac)); return "" }
	id:=self.calls[len(self.calls)-1].params.i[i]
	//id,ok:=igidentifier( //self.identifiers.i[i]
	//if !ok { wError(fmt.Sprintf("String expected for argument #%d, but what I got is nothing at all",i+1)); return "" }
	switch id.itype{
		case "string": return id.vstring
		case "integer": return fmt.Sprintf("%d",id.vint)
		case "bool","boolean":
			if id.vint<=0 { return "false"; } else { return "true"; }
		case "float":
			return fmt.Sprintf("%f",id.vfloat)
		default:
			wError(fmt.Sprintf("String expected for argument #%d, but what I got is %s!",i+1,id.itype))
			return ""
	}
}

func (self *VM) Arg_OptString(i int, defaultstring string) string{
	ac:=self.Arg_Count()
	if i>=ac {  return defaultstring }
	id:=self.calls[self.ccall].params.i[i]
	//id,ok:=self.identifiers.i[ak]
	//if !ok {  return defaultstring }
	if id.itype=="nil" {  return defaultstring }
	return self.Arg_GetString(i)
}

func (self *VM) Arg_GetInt(i int) int64 {
	ac:=self.Arg_Count()
	if i>=ac { wError(fmt.Sprintf("Integer expected as argument #%d, while only %d arguments are present",i+1,ac)); return 0 }
	id:=self.calls[self.ccall].params.i[i]
	//id,ok:=self.identifiers.i[ak]
	//if !ok { wError(fmt.Sprintf("Integer expected for argument #%d, but what I got is nothing at all",i+1)); return 0 }
	switch id.itype{
		case "integer": return id.vint
		case "float":   return int64(id.vint)
		default:
			wError(fmt.Sprintf("Integer expected for argument #%d, but what I got is %s!",i+1,id.itype))
			return 0
	}
}

func (self *VM) Arg_OptInt(i int, defaultint int64) int64{
	ac:=self.Arg_Count()
	if i>=ac {  return defaultint }
	id:=self.calls[self.ccall].params.i[i]
	//id,ok:=self.identifiers.i[ak]
	//if !ok {  return defaultint }
	if id.itype=="nil" {  return defaultint }
	return self.Arg_GetInt(i)
}

func (self *VM) Arg_GetFloat(i int) float64 {
	ac:=self.Arg_Count()
	if i>=ac { wError(fmt.Sprintf("Float expected as argument #%d, while only %d arguments are present",i+1,ac)); return 0 }
	id:=self.calls[self.ccall].params.i[i]
	//id,ok:=self.identifiers.i[ak]
	//if !ok { wError(fmt.Sprintf("Float expected for argument #%d, but what I got is nothing at all",i+1)); return 0 }
	switch id.itype{
		case "integer": return float64(id.vint)
		case "float":   return id.vfloat
		default:
			wError(fmt.Sprintf("Float expected for argument #%d, but what I got is %s.",i+1,id.itype))
			return 0
	}
}

func (self *VM) Arg_OptFloat(i int, defaultfloat float64) float64{
	ac:=self.Arg_Count()
	if i>=ac {  return defaultfloat }
	id:=self.calls[self.ccall].params.i[i]
	//id,ok:=self.identifiers.i[ak]
	//if !ok {  return defaultfloat }
	if id.itype=="nil" {  return defaultfloat }
	return self.Arg_GetFloat(i)
}


func (self *VM) Return_String(s string) {
	rtc:=self.calls[len(self.calls)-1]
	rtv:=&tIdentifier{}
	rtv.itype="string"
	rtv.vstring=s
	rtc.returns.i = append(rtc.returns.i,rtv)
	//fmt.Println("Returning: ",rtc.returns.i) // debug
}

func (self *VM) Return_Int(s int64) {
	rtc:=self.calls[len(self.calls)-1]
	rtv:=&tIdentifier{}
	rtv.itype="integer"
	rtv.vint=s
	rtc.returns.i = append(rtc.returns.i,rtv)
}

func (self *VM) Return_Float(s float64) {
	rtc:=self.calls[len(self.calls)-1]
	rtv:=&tIdentifier{}
	rtv.itype="float"
	rtv.vfloat=s
	rtc.returns.i = append(rtc.returns.i,rtv)
}

func (self *VM) Return_Bool(s bool) {
	rtc:=self.calls[len(self.calls)-1]
	rtv:=&tIdentifier{}
	rtv.itype="boolean"
	if s {
		rtv.vint=1
	}else{
		rtv.vint=0
	}
	rtc.returns.i = append(rtc.returns.i,rtv)
}


func (self *VM) Var_Get(v string) *tIdentifier{
    if qstr.Prefixed(v,"$__GETARG") {
      d:=v[len("$__GETARG"):]
      dv, err := strconv.ParseInt(d, 10, 64)
      if err!=nil { wError(err.Error()); return nil}
      if dv>=int64(self.Arg_Count()) {return nil}
      return self.calls[self.ccall].params.i[dv]
    }
    // If there's nothing else, just return the var if possible
    value,ok:=self.identifiers.i[v]
    if !ok {return nil}
    return value
}

// Minimal APIs
func wi_api_print(w *VM) (bool,string){
  // declare
  r:=[]string{}
  // form this dline
  l:=w.Arg_Count()
  for j:=0;j<l;j++{
    r=append(r,w.Arg_ConvString(j))
  }
  // execute
  for i,v := range r {
    if i!=0 { fmt.Print("\t"); }
    fmt.Print(v)
  }

  // Closure
  fmt.Println()
  return true,""
}

func wi_api_replace(w *VM) (bool,string){
	haystack:=w.Arg_GetString(0)
	needle:=w.Arg_GetString(1)
	replacewith:=w.Arg_GetString(2)
	cnt:=w.Arg_OptInt(3,-1)
	w.Return_String(strings.Replace(haystack,needle,replacewith,int(cnt)))
	return TBError()
}

// register
func (w *API) minapi(){
  w.Register("PRINT",wi_api_print)
  w.Register("REPLACE",wi_api_replace)
  wim_init(w)
}
