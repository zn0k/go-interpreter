package main

import (
	"fmt"
	"interpreter/compiler"
	"interpreter/evaluator"
	"interpreter/lexer"
	"interpreter/object"
	"interpreter/parser"
	"interpreter/repl"
	"interpreter/vm"
	"io/ioutil"
	"os"
	"os/user"
)

func main() {
	// check if a parameter was given. if so, treat it as an input file
	if len(os.Args) > 1 {
		buf, err := ioutil.ReadFile(os.Args[1])
		if err != nil {
			panic(fmt.Sprintf("unable to open %s as input file", os.Args[1]))
		}
		input := string(buf)

		//env := object.NewEnvironment()
		macroEnv := object.NewEnvironment()

		l := lexer.New(input)
		p := parser.New(l)

		program := p.ParseProgram()

		if len(p.Errors()) != 0 {
			repl.PrintParserErrors(os.Stdout, p.Errors())
			panic("errors parsing input file")
		}

		evaluator.DefineMacros(program, macroEnv)
		expanded := evaluator.ExpandMacros(program, macroEnv)

		// switch to VM mode
		/*evaluated := evaluator.Eval(expanded, env)

		if evaluated != nil {
			fmt.Println(evaluated.Inspect())
		}*/

		comp := compiler.New()
		err = comp.Compile(expanded)
		if err != nil {
			fmt.Printf("Compiler error: %s]n", err)
		}

		machine := vm.New(comp.Bytecode())
		err = machine.Run()
		if err != nil {
			fmt.Printf("VM execution error: %s\n", err)
		}

		lastPopped := machine.LastPoppedStackElem()
		fmt.Println(lastPopped.Inspect())
	} else {
		// open REPL
		user, err := user.Current()
		if err != nil {
			panic(err)
		}

		fmt.Printf("Hello, %s! Interpreter REPL\n\n", user.Username)
		repl.Start(os.Stdin, os.Stdout)
	}
}
