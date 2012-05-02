go-gnuplot
==========

Simple-minded functions to work with ``gnuplot``.
``go-gnuplot`` runs ``gnuplot`` as a subprocess and pushes commands
via the ``STDIN`` of that subprocess.

See http://www.gnuplot.info for more informations on the
exact semantics of these commands.

Installation
------------

The ``go-gnuplot`` package is ``go get`` installable::

   $ go get bitbucket.org/binet/go-gnuplot/pkg/gnuplot


Example
--------

::

    package main
    
    import "bitbucket.org/binet/go-gnuplot/pkg/gnuplot"
    import "fmt"
    
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

.. image:: https://bitbucket.org/binet/go-gnuplot/src/tip/examples/imgs/plot002.pdf


Documentation
-------------

API documentation can be found here:

 http://gopkgdoc.appspot.com/pkg/bitbucket.org/binet/go-gnuplot/pkg/gnuplot

