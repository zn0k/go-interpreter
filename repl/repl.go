package repl

import (
	"bufio"
	"fmt"
	"interpreter/compiler"
	"interpreter/evaluator"
	"interpreter/lexer"
	"interpreter/object"
	"interpreter/parser"
	"interpreter/vm"
	"io"
)

const PROMPT = ">>"

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)
	//env := object.NewEnvironment()
	macroEnv := object.NewEnvironment()

	fmt.Fprintf(out, PROMPT)
	for scanner.Scan() {
		line := scanner.Text()
		l := lexer.New(line)
		p := parser.New(l)

		program := p.ParseProgram()

		if len(p.Errors()) != 0 {
			PrintParserErrors(out, p.Errors())
			continue
		}

		evaluator.DefineMacros(program, macroEnv)
		expanded := evaluator.ExpandMacros(program, macroEnv)

		// use the VM instead, which isn't very functional right now
		/*evaluated := evaluator.Eval(expanded, env)
		if evaluated != nil {
			io.WriteString(out, evaluated.Inspect())
			io.WriteString(out, "\n")
		}

		fmt.Fprintf(out, PROMPT)*/

		comp := compiler.New()
		err := comp.Compile(expanded)
		if err != nil {
			fmt.Fprintf(out, "Compiler error: %s\n", err)
			continue
		}

		machine := vm.New(comp.Bytecode())
		err = machine.Run()
		if err != nil {
			fmt.Fprintf(out, "VM execution error: %s\n", err)
			continue
		}

		stackTop := machine.StackTop()
		io.WriteString(out, stackTop.Inspect())
		io.WriteString(out, "\n")
	}
}

func PrintParserErrors(out io.Writer, errors []string) {
	for _, msg := range errors {
		io.WriteString(out, "\t"+msg+"\n")
	}
}
