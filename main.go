package main

import (
	"MASTANk/components"
	"fmt"
	"math"
	"gonum.org/v1/gonum/mat"
)

func main() {
	a := components.MakeNode(0, 0, true, true)
	b := components.MakeNode(5, 0, false, true)
	c := components.MakeNode(5, 5 * math.Tan(60 * math.Pi / 180), false, false)
	d := components.MakeNode(10, 5 * math.Tan(60 * math.Pi / 180), false, false)

	components.MakeElement(a, b, 200000e6, 10e-3)
	components.MakeElement(a, c, 200000e6, 15e-3)
	components.MakeElement(c, d, 200000e6, 10e-3)
	components.MakeElement(b, d, 200000e6, 15e-3)
	components.MakeElement(b, c, 200000e6, 15e-3)

	components.ApplyPointLoad(d, 400e3 * math.Cos(45 * math.Pi / 180), -400e3 * math.Cos(45 * math.Pi / 180))

	
	displacements := components.Solve()
	
	fmt.Println(mat.Formatted(displacements, mat.Prefix(""), mat.Squeeze()))
	fmt.Println(components.ElementList[0].P)
	
}