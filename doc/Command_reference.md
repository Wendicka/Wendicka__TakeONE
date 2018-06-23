# Wendicka

Wendicka is a byte-code engine for scripting.
A builder or compiler can be used for prototyping or to make code translators output in text (may be handier to debug them, I guess.

The syntax has been inspired on assembly as that was the easiest to parse :P
I have something more in mind for that later, but no promises there yet, and that will be a separate project USING Wendicka.

This is a quick command overview.
Now please note, commands are case INSENSITIVE... The identifiers are not.
Variables names MUST be prefixed with a $ and any variable prefixed with $__ is basically reserved for Wendicka's internal workings.


Reserved variables:
- $__ARG
  - Any value assigned to it will be an argument or parameter for the next "CALL" command
  - Arguments/parameters are taken in the same order as they are assigned
  - The variable is WRITE-ONLY. Trying to read from it will NEVER work.
  
  
Commands:
- CALL callableitem
  - Calls a chunk or an registered API function from the underlying engine
- CHECK value1,checktype,value2
  - Boolean check
  - You can use =,!=,<,>,>=,<= and such to make various boolean checks on the checktype field 
  - Commas are required here! This is a primitive parser after all :P
- CHUNK ChunkName
  - Sets up a new Chunk
    - Chunks can be compared to functions/procedures in languages like C, Go, Pascal and such. 
- DEC $var,value
  - Decreases $var by value
  - Only works on nummeric variables and values
- EXIT
  - Quits the current chunk call and goes back to the parent chunk call
  - Can alternatively be replaced by its alias END
  - Whenever a new Chunk is created with the CHUNK command, this command is added to the previous chunk automatically
- INC $var,value
  - Increases $var by value
  - Only works on nummeric variables and values
- JUMP Labelname
  - JUMP to a certain label
  - The variant PJUMP only works if the last CHECK or SOMETHING commands returned true and will otherwise be ignored.
  - The variant NJUMP only works if the last CHECK or SOMETHING commands returned false and will otherwise be ignored.
  - Please note jump instructions must be in the same chunk as the labels they jump to. You cannot jump cross-chunk.
- LABEL Labelname
  - Hookpoint for jump instructions to jump to
  - Please note jump instructions must be in the same chunk as the labels they jump to. You cannot jump cross-chunk.
- MOVE $var,value
  - Assigns value to $var  
  - Assigning parameters to $__ARG and other reserved variables also works through this command.
  - Alternativeley you can also use the alias command MOV
- RETURN value
  - Set the return value(s) for exiting this chunk
  - Please note, you *can* return multiple values, by just adding multiple RETURN commands
  - Also note that unlike languages such a C,C++,Java,php,Go,BlitzMax,Lua and many others, the RETURN command will in Wendicka NOT end the chunk execution. If this has to happen after return you ARE required to use an EXIT command next. That is neither a bug nor an oversight, but done deliberately!!
- SOMETHING $var
  - set the last check to "true" if:
    - The variable contains a nummeric value higher than 0
    - A boolean value set to true
    - Contains a string that is not empty
    - Points to an existing chunk
    - Points to an existing table
