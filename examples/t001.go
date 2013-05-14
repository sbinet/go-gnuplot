package main

import "fmt"

//import "io/ioutil"
import "github.com/sbinet/go-gnuplot"

func main() {
	fname := ""
	persist := false
	debug := true

	p, err := gnuplot.NewPlotter(fname, persist, debug)
	if err != nil {
		err_string := fmt.Sprintf("** err: %v\n", err)
		panic(err_string)
	}
	defer p.Close()

	p.CheckedCmd("plot %f*x", 23.0)
	p.CheckedCmd("plot %f * cos(%f * x)", 32.0, -3.0)
	//p.CheckedCmd("save foo.ps")
	p.CheckedCmd("set terminal pdf")
	p.CheckedCmd("set output 'plot001.pdf'")
	p.CheckedCmd("replot")

	p.CheckedCmd("q")
	//p.proc.Wait(0)

	return
}

/* EOF */
