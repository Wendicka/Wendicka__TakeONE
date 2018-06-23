package wi

import(
  "trickyunits/qint"
  "trickyunits/qstr"
  "strconv"
  "errors"
  "fmt"
)

type winstruct struct {
  do_it func(w *VM,args [][]byte) bool
  needparam []string
}

var winstructs = map[byte] *winstruct {}

func igidentifier(w *VM,b[]byte) (*tIdentifier,error){
  i:=b[0]
  p:=b[1:]
  switch i{
    case 0:
      r:=tIdentifier{}
      r.itype="string"
      r.vstring=string(p)
      return &r,nil
    case 1:
      r:=w.identifiers.i
      return r[string(p)],nil
    case 2:
      var e error
      r:=tIdentifier{}
      r.itype="integer"
      r.vint,e=qint.BytesToInt64(p)
      //fmt.Println("Form integer",b,i,p,r.vint)
      return &r,e
    case 3:
      var e error
      r:=tIdentifier{}
      r.itype="float"
      r.vfloat,e=qint.BytesToFloat(p)
      return &r,e
    default:
      fmt.Print(b)
      return nil,errors.New(fmt.Sprintf("Unknown Identifier Code. Either the code is compiled incorrectly, or the code is corrupt, or your version of Wendicka is too old for this particular code. (%3d)",i))
  }
  return nil,errors.New("I'm completely at a loss. This error could only happen to either a bug or an outdated version of Wendicka.")
}

func init(){
  winstructs[1] = &winstruct{
    // MOV
    func (w *VM,args [][]byte) bool{
      v:=string(args[0])
      if args[0][0]!='$' {
        wError("Invalid MOV! "+v)
        return false
      }
      if len(args)<=1 { wError("No values received to appoint to: "+string(args[0])); return false}
      id,e:=igidentifier(w,args[1])
      if e!=nil { wError(e.Error()); return false;}
      switch v {
        case "$__ARG":
          w.nextcallparam.i = append(w.nextcallparam.i,id)
          chat("Argument ",id.itype,string(args[1])," added for next call")
          return true
        default:
          kut:=&w.identifiers.i
          (*kut)[v]=id
          chat("Value ",string(args[1]),"("+id.itype+")"," assigned to variable ",v,fmt.Sprintf("%d",id.vint))
          return true
        }
    },
    []string{"string","identifier"},
  }

  winstructs[2] = &winstruct{
    // INC
    func (w *VM,args [][]byte) bool{
      v:=string(args[0])
      if args[0][0]!='$' {
        wError("Invalid MOV! "+v)
        return false
      }
      if len(args)<=1 { wError("No values received to alter: "+string(args[0])); return false}
      vr,ok:=w.identifiers.i[string(args[0])]
      if !ok { wError("Non-existent variable "+string(args[0])+" cannot be altered"); return false}
      ch,e:=igidentifier(w,args[1])
      if e!=nil { wError(e.Error()); return false;}
      switch(vr.itype){
		case "integer":
			switch(ch.itype){
				case "integer":
					vr.vint += ch.vint
				case "float":
					vr.vint += int64(ch.vfloat)
				default:
					wError("I cannot alter a "+vr.itype+" with a "+ch.itype)
					return false
			}
		case "float":
			switch(ch.itype){
				case "integer":
					vr.vfloat += float64(ch.vint)
				case "float":
					vr.vfloat += ch.vfloat
				default:
					wError("I cannot alter a "+vr.itype+" with a "+ch.itype)
					return false
			}
		default:
			wError("Illegal identifier for INC "+vr.itype)
			return false
		}
      return true
	  },
    []string{"string","identifier"},
  }
  winstructs[3] = &winstruct{
    // DEC
    func (w *VM,args [][]byte) bool{
      v:=string(args[0])
      if args[0][0]!='$' {
        wError("Invalid MOV! "+v)
        return false
      }
      if len(args)<=1 { wError("No values received to alter: "+string(args[0])); return false}
      vr,ok:=w.identifiers.i[string(args[0])]
      if !ok { wError("Non-existent variable "+string(args[0])+" cannot be altered"); return false}
      ch,e:=igidentifier(w,args[1])
      if e!=nil { wError(e.Error()); return false;}
      switch(vr.itype){
		case "integer":
			switch(ch.itype){
				case "integer":
					vr.vint -= ch.vint
				case "float":
					vr.vint -= int64(ch.vfloat)
				default:
					wError("I cannot alter a "+vr.itype+" with a "+ch.itype)
					return false
			}
		case "float":
			switch(ch.itype){
				case "integer":
					vr.vfloat -= float64(ch.vint)
				case "float":
					vr.vfloat -= ch.vfloat
				default:
					wError("I cannot alter a "+vr.itype+" with a "+ch.itype)
					return false
			}
		default:
			wError("Illegal identifier for INC "+vr.itype)
			return false
		}
      return true
    },
    []string{"string","identifier"},
  }
  winstructs[6] = &winstruct{
	  func (w *VM,args[][]byte) bool{
		  id1,e1:=igidentifier(w,args[1])
		  id2,e2:=igidentifier(w,args[2])
		  if e1!=nil { wError("ID1,MUL:\t"+e1.Error()); return false }
		  if e2!=nil { wError("ID2,MUL:\t"+e2.Error()); return false }
		  if id1==nil { wError("MUL.ID1==nil"); return false }
		  if id2==nil { wError("MUL.ID2==nil"); return false }
		  uitkomst:=&tIdentifier{}
		  if args[0][0]!='$' { wError("MULTIPLY requires variable to store the result!"); return false}
		  sarg:=string(args[0])
		  if id1.itype=="integer" && id2.itype=="integer" {
			uitkomst.itype="integer"
			uitkomst.vint=id1.vint*id2.vint
		  } else {
			p1:=float64(0)
			p2:=float64(0)
			if id1.itype=="float" { p1=id1.vfloat
			} else if id1.itype=="integer" { p1=float64(id1.vint)
			} else { wError("1st value for multiplication is invalid: "+id1.itype); return false}
			if id2.itype=="float" { p2=id2.vfloat
			} else if id2.itype=="integer" { p2=float64(id2.vint)
			} else { wError("2nd value for multiplication is invalid: "+id2.itype); return false}
			uitkomst.itype="float"
			uitkomst.vfloat=p1*p2
		  }
		  store:=&w.identifiers.i
		  (*store)[sarg]=uitkomst
		  return true
	  },[]string{"string","identifier","identifier"},
  }
  
  winstructs[8] = &winstruct{
	  func (w *VM,args[][]byte) bool {
		arg:=args[0]
		sarg:=string(arg)
		if arg[0]!='$' { wError("Raw input needs a VARIABLE to store the input in!"); return false }
		answer:=qstr.RawInput("")
		store:=&w.identifiers.i
		id:=tIdentifier{}
		id.itype="string"
		id.vstring=answer
		(*store)[sarg]=&id
		return true
	  },
	  []string{"string"},
  }

  winstructs[9] = &winstruct{
	  func (w *VM,args[][]byte) bool {
		arg:=args[0]
		sarg:=string(arg)
		if arg[0]!='$' { wError("Raw int input needs a VARIABLE to store the input in!"); return false }
		ok:=false
		ia:=int64(0)
		for !ok{
			answer:=qstr.RawInput("")
			a,e:=strconv.ParseInt(answer, 10, 64)
			if e!=nil {
				fmt.Print("Redo from start: "+e.Error())
			} else {
				ok=true
				ia=a
			}
		}
		store:=&w.identifiers.i
		id:=tIdentifier{}
		id.itype="integer"
		id.vint=ia
		(*store)[sarg]=&id
		return true
	  },
	  []string{"string"},
  }


  // RETURN
  winstructs[11] = &winstruct{
		func(w *VM,args[][]byte) bool {
			cl:=w.calls[len(w.calls)-1]
			r:=cl.returns
			i,_:=igidentifier(w,args[0])
			r.i=append(r.i,i)
			return true
		},[]string{"identifier"},
  }
  // GETRET
  winstructs[12] = &winstruct{
		func(w *VM,args[][]byte) bool {      
			v:=string(args[0])
			if args[0][0]!='$' {
				wError("Invalid GETRET! "+v)
				return false
			}
		if len(args)<2 { wError("GETRET error!  "+string(args[1])); return false}
		//id,e:=igidentifier(w,args[1])
		//if e!=nil { wError(e.Error()); return false;}
		idnum,e:=igidentifier(w,args[1])
		if e!=nil { wError(e.Error()); return false}
		if idnum.itype!="integer" { wError("Return index must be integer"); return false}
		inum:=int(idnum.vint)
		if inum>=len(w.lastcallreturn.i) { wError("Return data out of range"); return false}
		id:=w.lastcallreturn.i[inum]
		if qstr.Prefixed(v,"$__") { wError("GETRET may not be used with potential reserved variables!"); return false }
		kut:=&w.identifiers.i
		(*kut)[v]=id
		//chat("Value ",string(args[1]),"("+id.itype+")"," assigned to variable ",v,fmt.Sprintf("%d",id.vint))
		return true
        },
		[]string{"string","identifier"},
  }



  // CALL
  winstructs[10] = &winstruct{
    func (w *VM,args [][]byte) bool{
          ch:=string(args[0])
          if args[0][0]=='$' {
            id:=w.identifiers.i[ch]
            if id==nil {wError(ch+" appears to be nothing at all. What must I call?"); return false}
            if id.itype=="chunk" || id.itype=="api" { ch=id.vstring } else {wError("Call identifiers must refer to apis or chunks"); return false}
          }
          afs:=w.vapi.fs
          if _,ok:=w.chunks[ch];ok {
            w.callChunk(ch)
          } else if _,aok:=(*afs)[ch];aok {
            s,e:=w.callAPI(ch)
            if !s {wError("API Call error: "+e); return false}
          } else {
			  wError("No chunk nor api called "+ch+" has been found!")
			  return false
		  }
      return true
    },
    []string{"string"},
  }
  // Something
  winstructs[249] = &winstruct{
	  func (w *VM, args[][]byte) bool{
		  id,e:=igidentifier(w,args[0])
		  if e!=nil { wError(e.Error()); return false }
		  switch id.itype {
			  case "nil":
				w.lastcompare = false
			  case "int","integer","bool","boolean":
				w.lastcompare = id.vint>0
			  case "float":
				w.lastcompare = id.vfloat>0
			  case "string":
				w.lastcompare = id.vstring!=""
			  default:
			    fmt.Print("WARNING! SOMETHING type unknown: "+id.itype)
			}
		  return true
		  },
		  []string{"identifier"},
  }
  
  // Check
  winstructs[250] = &winstruct{
	  func(w *VM, args[][]byte) bool{
		  id1,e1:=igidentifier(w,args[1])
		  id2,e2:=igidentifier(w,args[2])
		  d:=args[0][0]
		  if e1!=nil { wError("CHECK#1"+e1.Error()); return false; }
		  if e2!=nil { wError("CHECK#2"+e2.Error()); return false; }
		  switch d{
			  case 0,1:
				if id1.itype==id2.itype && id1.itype=="string" { 
					w.lastcompare=id1.vstring == id2.vstring
				} else if id1.itype==id2.itype && id1.itype=="integer" {
					w.lastcompare=id1.vint == id2.vint
				} else if id1.itype==id2.itype && id1.itype=="float" {
					w.lastcompare=id1.vfloat == id2.vfloat
				} else {
					fmt.Println("WARNING! Unsupported compare ("+id1.itype+"<=>"+id2.itype+")")
					w.lastcompare = false
				}
				if d==1 { w.lastcompare=!w.lastcompare }
			  case 3:
				if id1.itype==id2.itype && id1.itype=="integer" {
					w.lastcompare = id1.vint < id2.vint
					chat(fmt.Sprintf("%s %d < %s %d ",id1.itype,id1.vint,id2.itype,id2.vint))
				} else if id1.itype==id2.itype && id1.itype=="float" {
					w.lastcompare = id1.vfloat < id2.vfloat
				} else if id1.itype=="float" && id2.itype=="integer" {
					w.lastcompare = id1.vfloat < float64(id2.vint)
				} else if id1.itype=="float" && id2.itype=="integer" {
					w.lastcompare = float64(id1.vint) < id2.vfloat
				} else {
					fmt.Println("WARNING! Unsupported smaller compare ("+id1.itype+"<=>"+id2.itype+")")
					w.lastcompare = false
				}
			  case 4:
				if id1.itype==id2.itype && id1.itype=="integer" {
					w.lastcompare = id1.vint > id2.vint
				} else if id1.itype==id2.itype && id1.itype=="float" {
					w.lastcompare = id1.vfloat > id2.vfloat
				} else if id1.itype=="float" && id2.itype=="integer" {
					w.lastcompare = id1.vfloat > float64(id2.vint)
				} else if id1.itype=="float" && id2.itype=="integer" {
					w.lastcompare = float64(id1.vint) > id2.vfloat
				} else {
					fmt.Println("WARNING! Unsupported greater compare ("+id1.itype+"<=>"+id2.itype+")")
					w.lastcompare = false
				}
			  case 5:
				if id1.itype==id2.itype && id1.itype=="integer" {
					w.lastcompare = id1.vint >= id2.vint
				} else if id1.itype==id2.itype && id1.itype=="float" {
					w.lastcompare = id1.vfloat >= id2.vfloat
				} else if id1.itype=="float" && id2.itype=="integer" {
					w.lastcompare = id1.vfloat >= float64(id2.vint)
				} else if id1.itype=="float" && id2.itype=="integer" {
					w.lastcompare = float64(id1.vint) >= id2.vfloat
				} else {
					fmt.Println("WARNING! Unsupported greater equal compare ("+id1.itype+"<=>"+id2.itype+")")
					w.lastcompare = false
				}
			  case 6:
				if id1.itype==id2.itype && id1.itype=="integer" {
					w.lastcompare = id1.vint <= id2.vint
				} else if id1.itype==id2.itype && id1.itype=="float" {
					w.lastcompare = id1.vfloat <= id2.vfloat
				} else if id1.itype=="float" && id2.itype=="integer" {
					w.lastcompare = id1.vfloat <= float64(id2.vint)
				} else if id1.itype=="float" && id2.itype=="integer" {
					w.lastcompare = float64(id1.vint) <= id2.vfloat
				} else {
					fmt.Println("WARNING! Unsupported smaller equal compare ("+id1.itype+"<=>"+id2.itype+")")
					w.lastcompare = false
				}
			  default:
				wError(fmt.Sprintf("Unknown comparing code: %d!",d))
				return false
		  }
		  chat(fmt.Sprintf("d = %d #%d",d,len(args[0])))
		  if w.lastcompare { chat("CHECK = TRUE"); } else { chat("CHECK = FALSE"); }
		  return true
	  },
	  []string{"byte","identifier","identifier"},
  }
  
  // Jump
  winstructs[0xfd] = &winstruct{
	  func(w *VM,args [][]byte) bool{
		  l:=string(args[0])
		  cll:=len(w.calls)-1
		  cl:=w.calls[cll]
		  ch:=cl.achunk
		  pos,ok:=ch.labels[l]
		  if !ok { wError("Jump to non-existent label: "+l); return false }
		  cl.pos=pos-1 // -1 has to be taken in order for the pos++ instruction that always comes next!
		  return true
	  },
	  []string{"string"},  
  }
  
  // Conditional jumping
  winstructs[252] = &winstruct{ // false
	  func(w *VM,args [][]byte) bool{
		  if !w.lastcompare {
			  d:=winstructs[0xfd]
			  return d.do_it(w,args)
		  } else { return true }
	  },
	  []string{"string"},  
  }
  winstructs[251] = &winstruct{ // true
	  func(w *VM,args [][]byte) bool{
		  if w.lastcompare {
			  d:=winstructs[0xfd]
			  return d.do_it(w,args)
		  } else { return true }
	  },
	  []string{"string"},  
  }

  // END CALL
  winstructs[0xff] = &winstruct{
    func (w *VM,args [][]byte) bool{
		cll:=len(w.calls)-1
		cl:=w.calls[cll]
		cl.ended=true
		w.calls=w.calls[:cll]
      return true
    },
    []string{},

  }
}
