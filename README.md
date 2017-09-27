go-gnuplot
==========

Simple-minded functions to work with ``gnuplot``.
``go-gnuplot`` runs ``gnuplot`` as a subprocess and pushes commands
via the ``STDIN`` of that subprocess.

See http://www.gnuplot.info for more informations on the
exact semantics of these commands.

This is a fork of
[sbinet/go-gnuplot](https://www.github.com/sbinet/go-gnuplot). The fork is
motivated by the lack of maintenance to the original repo. This version will
aim to extend on and fix bugs with the original implementation. The original
API will not change so that this can become a drag and drop replacement of the
original. See the issues for specific planned changes.

Installation
------------

The ``go-gnuplot`` package is ``go get`` installable:

```sh
$ go get github.com/ckitagawa/go-gnuplot
```

Example
--------

```go
package main

import (
  "github.com/ckitagawa/go-gnuplot"
  "fmt"
)

func main() {
	fname := ""
	persist := false
	debug := true

	p,err := gnuplot.NewPlotter(fname, persist, debug)
	if err != nil {
		err_string := fmt.Sprintf("** err: %v\n", err)
		panic(err_string)
	}
	defer p.Close()

	p.PlotX([]float64{0,1,2,3,4,5,6,7,8,9,10}, "some data")
	p.CheckedCmd("set terminal pdf")
	p.CheckedCmd("set output 'plot002.pdf'")
	p.CheckedCmd("replot")


	p.CheckedCmd("q")
	return
}
```

![plot-t-002](https://github.com/ckitagawa/go-gnuplot/raw/master/examples/imgs/plot002.png)


Documentation
-------------

API documentation can be found here:

 http://godoc.org/github.com/ckitagawa/go-gnuplot
