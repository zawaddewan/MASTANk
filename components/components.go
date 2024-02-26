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

func MakeNode(X float64, Y float64, isfixedX bool, isfixedY bool) *Node {
	node := Node{X, Y, isfixedX, isfixedY, 0, 0, false, nil}
	NodeList = append(NodeList, &node)
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

type Element struct {
	N1		*Node
	N2		*Node
	E 		float64
	A		float64
	Stiff	*mat.Dense
	P		float64
}

var ElementList []*Element = make([]*Element, 0)

func MakeElement(n1 *Node, n2 *Node, e float64, a float64) *Element {
	element := Element{n1, n2, e, a, mat.NewDense(4, 4, nil), 0}
	element.genStiffness()
	ElementList = append(ElementList, &element)
	return &element
}

func (e Element) genStiffness() {
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

func Solve() *mat.VecDense {
	lastFree := 0
	lastFixed := 0

	ext := make([]float64, 2*len(NodeList))
	
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
			ext[NodeList[i].XDeg] = NodeList[i].Load[0]
			ext[NodeList[i].YDeg] = NodeList[i].Load[1]	
		}
	}

	ext = ext[0:lastFree]

	global := mat.NewDense(2*len(NodeList), 2*len(NodeList), nil)
	gRows, _ := global.Dims()
	
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

	globalFree := global.Slice(0, lastFree, 0, lastFree)
	extFree := mat.NewVecDense(lastFree, ext)
	del := mat.NewVecDense(lastFree, nil)
	del.SolveVec(globalFree, extFree)

	globalFixed := global.Slice(lastFree, gRows, 0, lastFree)
	extFixed := mat.NewVecDense(lastFixed, nil)
	extFixed.MulVec(globalFixed, del)

	for i := 0; i < len(ElementList); i++ {
		local := make([]float64, 4)
		if !ElementList[i].N1.FixedX {
			local[0] = del.AtVec(ElementList[i].N1.XDeg)
		}
		if !ElementList[i].N1.FixedY {
			local[1] = del.AtVec(ElementList[i].N1.YDeg)
		}
		if !ElementList[i].N2.FixedX {
			local[2] = del.AtVec(ElementList[i].N2.XDeg)
		}
		if !ElementList[i].N2.FixedY {
			local[3] = del.AtVec(ElementList[i].N2.YDeg)
		}

		vec := mat.NewVecDense(4, local)
		vec.MulVec(ElementList[i].Stiff.Slice(0, 4, 0, 4), vec)
		member := mat.NewVecDense(2, []float64{ElementList[i].N1.X - ElementList[i].N2.X, ElementList[i].N1.Y - ElementList[i].N2.Y})
		ElementList[i].P = mat.Norm(vec.SliceVec(0, 2), 2) * math.Copysign(1, mat.Dot(vec.SliceVec(0, 2), member))
	}

	return del
}







