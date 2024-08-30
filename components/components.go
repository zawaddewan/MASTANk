package components

import (
	"gonum.org/v1/gonum/mat"
	"math"
)

type Node struct {
	X		float64
	Y		float64
	FixedX	bool
	FixedY	bool
	XDeg	int
	YDeg	int
	Loaded	bool
	Load 	[]float64
}

var NodeList []*Node = make([]*Node, 0)
var FixedNodes map[*Node]int = make(map[*Node]int)

func MakeNode(X float64, Y float64, isfixedX bool, isfixedY bool) *Node {
	node := Node{X, Y, isfixedX, isfixedY, 0, 0, false, nil}
	NodeList = append(NodeList, &node)
	if isfixedX || isfixedY {
		FixedNodes[&node] = len(NodeList) - 1
	}
	return &node
}

func Dist(n1 Node, n2 Node) float64 {
	return math.Sqrt((n1.X - n2.X)*(n1.X - n2.X) + (n1.Y - n2.Y)*(n1.Y - n2.Y))
}

func NodeSin(n1 Node, n2 Node) float64 {
	return (n2.Y - n1.Y) / Dist(n1, n2)
}

func NodeCos(n1 Node, n2 Node) float64 {
	return (n2.X - n1.X) / Dist(n1, n2)
}

func ApplyPointLoad(n *Node, fx float64, fy float64) {
	n.Loaded = true
	n.Load = []float64{fx, fy}
}

func OrderDegreesFreedom() []float64 {
	
	lastFree := 0
	lastFixed := 0
	degrees := make([]float64, 2*len(NodeList))
		
	for i := 0; i < len(NodeList); i++ {
		if NodeList[i].FixedX {
			NodeList[i].XDeg = 2*len(NodeList) - 1 - lastFixed
			lastFixed++
		}else {
			NodeList[i].XDeg = lastFree
			lastFree++
		}
		
		if NodeList[i].FixedY {
			NodeList[i].YDeg = 2*len(NodeList) - 1 - lastFixed
			lastFixed++
		}else {
			NodeList[i].YDeg = lastFree
			lastFree++
		}

		if NodeList[i].Loaded {
			degrees[NodeList[i].XDeg] = NodeList[i].Load[0]
			degrees[NodeList[i].YDeg] = NodeList[i].Load[1]	
		}
	}

	return degrees[0:lastFree]
}

type Section struct {
	Modulus float64
	Area float64
}

var SectionList []*Section = make([]*Section, 0)

func MakeSection(modulus float64, area float64) *Section {
	section := Section{modulus, area}
	SectionList = append(SectionList, &section)
	return &section
}

type Element struct {
	N1		*Node
	N2		*Node
	E 		float64
	A		float64
	Stiff	*mat.Dense
	P		float64
}

var ElementList []*Element = make([]*Element, 0)

func MakeElement(n1 *Node, n2 *Node, s ...*Section) *Element {
	element := Element{n1, n2, 0, 0, mat.NewDense(4, 4, nil), 0}
	if len(s) != 0 {
		element.ApplySection(s[0])
	}
	ElementList = append(ElementList, &element)
	return &element
}

func (e *Element) ApplySection(s *Section) {
	e.E = s.Modulus
	e.A = s.Area
	e.genStiffness()
}

func (e *Element) toVector() *mat.VecDense {
	return mat.NewVecDense(2, []float64{e.N1.X - e.N2.X, e.N1.Y - e.N2.Y})
}

func (e *Element) genStiffness() {
	stiff := mat.NewDense(4, 4, nil)
	for i := 0; i < 4; i++ {
		for j := 0; j < 4; j++ {
			entry := 1.0
			if i % 2 == 0 {
				entry *= NodeCos(*e.N1, *e.N2)
			}else {
				entry *= NodeSin(*e.N1, *e.N2)
			}

			if j % 2 == 0 {
				entry *= NodeCos(*e.N1, *e.N2)
			}else {
				entry *= NodeSin(*e.N1, *e.N2)
			}
			
			if (i > 1) != (j > 1) {
				entry *= -1
			}
			
			stiff.Set(i, j, entry)
		}
	}
	
	e.Stiff.Scale(e.E * e.A / Dist(*e.N1, *e.N2), stiff)
}

func (e *Element) CalcForces(del *mat.VecDense) {
	displacement := mat.NewVecDense(4, nil)

	if !e.N1.FixedX {
		displacement.SetVec(0, del.AtVec(e.N1.XDeg))
	}

	if !e.N1.FixedY {
		displacement.SetVec(1, del.AtVec(e.N1.YDeg))
	}

	if !e.N2.FixedX {
		displacement.SetVec(2, del.AtVec(e.N2.XDeg))
	}

	if !e.N2.FixedY {
		displacement.SetVec(3, del.AtVec(e.N2.YDeg))
	}

	displacement.MulVec(e.Stiff, displacement)
	e.P = mat.Norm(displacement.SliceVec(0, 2), 2) * math.Copysign(1, mat.Dot(displacement.SliceVec(0, 2), e.toVector()))	
}

func GenGlobal() *mat.Dense {
	global := mat.NewDense(2*len(NodeList), 2*len(NodeList), nil)
	
	for i := 0; i < len(ElementList); i++ {
		dimr, dimc := ElementList[i].Stiff.Dims()
		
		index := make([]int, 4)

		index[0] = ElementList[i].N1.XDeg
		index[1] = ElementList[i].N1.YDeg
		index[2] = ElementList[i].N2.XDeg
		index[3] = ElementList[i].N2.YDeg

		for j := 0; j < dimr; j++ {
			for k := 0; k < dimc; k++ {
				curr := global.At(index[j], index[k])
				add := ElementList[i].Stiff.At(j, k)
				global.Set(index[j], index[k], curr + add)
			}
		}	
	}
	
	return global
}

func Solve() (*mat.VecDense, *mat.VecDense) {

	ext := OrderDegreesFreedom()

	global := GenGlobal()

	globalFree := global.Slice(0, len(ext), 0, len(ext))

	
	extFree := mat.NewVecDense(len(ext), ext)
	del := mat.NewVecDense(len(ext), nil)
	err := del.SolveVec(globalFree, extFree)

	if(err != nil) {
		panic("Matrix is singular - free-body motion detected")
	}

	globalFixed := global.Slice(len(ext), 2*len(NodeList), 0, len(ext))
	extFixed := mat.NewVecDense(2*len(NodeList) - len(ext), nil)
	extFixed.MulVec(globalFixed, del)

	for i := 0; i < len(ElementList); i++ {
		ElementList[i].CalcForces(del)
	}

	return del, extFixed
}







