package main

import (
	"MASTANk/components"
	"fmt"
	"gonum.org/v1/gonum/mat"
)

func main() {
	E := 200e9
	A1 := 1.767145867e-4
	// A2 := 3.141592653e-4

	a := components.MakeNode(0,0, true, true)
	b := components.MakeNode(0.5, 0.5, false, false)
	c := components.MakeNode(0.5, 0, false, false)
	d := components.MakeNode(0, 0.5, true, false)

	components.MakeElement(a, b, E, A1)
	components.MakeElement(a, c, E, A1)
	components.MakeElement(a, d, E, A1)
	components.MakeElement(b, c, E, A1)
	components.MakeElement(b, d, E, A1)
	components.MakeElement(c, d, E, A1)

	components.ApplyPointLoad(b, 0, -1)

	
	displacements := components.Solve()
	
	fmt.Println(mat.Formatted(displacements, mat.Prefix(""), mat.Squeeze()))
	for i := 0; i < len(components.ElementList); i++ {
		fmt.Println(components.ElementList[i].P)
	}	
}