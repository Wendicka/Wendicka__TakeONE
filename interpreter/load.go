/*
  load.go
  
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

import(
  //"os"
  "fmt"
  "bytes"
  "errors"
  "trickyunits/qff"
  "trickyunits/qint"
)


func appint32(ori []byte,i int32) []byte {
    ret:=ori
    ai,_:=qint.Int32toBytes(i)
    for _,b := range ai{
      ret=append(ret,b)
    }
    return ret
}

// Convert bytes array into an actual Wendicka VM
func Bytes2VM(l []byte) (*VM,error){
  var e error
  ret:=NewVM(nil)
  bt:=bytes.NewReader(l)
  h:=qff.RawReadString(bt,3)
  readchunk:=""
  if h!="WBC" { return nil,errors.New("Given memory bank is NOT Wendicka Byte Code")}
  bits:=qff.ReadInt64(bt)
  if bits!=8 { return nil,errors.New(fmt.Sprintf("This version of the Wendicka interpreter only supports 8 bit instructions and this bytecode contains %d bit instructions",bits))}
  for bt.Len()>0 { //!qff.EOF(*bt) {
    ins:=qff.ReadByte(bt)
    if ins!=0 && readchunk=="" { return nil,errors.New(fmt.Sprintf("Chunkless code: %X",ins))}
    switch ins {
      case 0:
        namechunk:=qff.ReadString(bt)
        if _,exist:=ret.chunks[namechunk];exist { return nil,errors.New("Duplicate chunk: "+namechunk)}
        ret.chunks[namechunk]=&tchunk{}
        ret.chunks[namechunk].labels = map[string] int{}
        readchunk=namechunk
        chat("Created chunk: ",namechunk)
      case 254:
        namelabel:=qff.ReadString(bt)
        ch:=ret.chunks[readchunk]
        ch.labels[namelabel]=len(ch.instruction)
        chat("= Created label in chunk ",readchunk," named ",namelabel," as position ",fmt.Sprintf("%d",ch.labels[namelabel]))
      default:
        q,ok:=winstructs[ins]
        if !ok { return nil,errors.New(fmt.Sprintf("Unknown instruction: %2X/%3d",ins,ins))}
        a:=[][]byte {}//*tIdentifier{}
        for i,partype:=range q.needparam {
          switch partype {
            case "string":
              s:= qff.ReadString(bt)
              a = append(a,[]byte(s))
            case "byte":
			  b:= qff.ReadByte(bt)
			  c:= []byte{b}
			  a = append(a,c)
            case "identifier":
              b:=[]byte{}
              t:=qff.ReadByte(bt)
              b=append(b,t)
              switch(t){
                case 0,1: // string + identifiername
                  l:=qff.ReadInt32(bt)
                  //b=appint32(b,l)
                  for i:=int32(0);i<l;i++{ b=append(b,qff.ReadByte(bt));}
                case 2,3: // int64 + float64
                  for i:=0;i<8;i++{ b=append(b,qff.ReadByte(bt));}
                default:
                  return nil,errors.New("Invalid identifier type")
              }
              a = append(a,b)
            default:
              return nil,errors.New(fmt.Sprintf("Unknown paramter type for instruction %X, position %d: %s",ins,i,partype))
          }
        }
        cc:=&tchunkins{}
        cc.param=a
        cc.ins=ins
        ch:=ret.chunks[readchunk]
        ch.instruction = append(ch.instruction,cc)

        chat("Handled instruction:",fmt.Sprintf("%X in chunk %s",ins,readchunk))

    }
  }
  return ret,e
}

// Load a file and convert that into a Wendicka VM
func File2VM(filename string) (*VM,error){
  b,e:=qff.EGetFile(filename)
  if e!=nil { return nil,e }
  var r *VM
  r,e=Bytes2VM(b)
  return r,e
}
