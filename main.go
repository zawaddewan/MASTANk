package main

import (
	"github.com/fogleman/gg"
	"github.com/icza/gog"
	"MASTANk/components"
	"strings"
	"strconv"
	"bufio"
	"math"
	"fmt"
	"os"
	"io"
)


func main() {

	var min [2]float64
	var max [2]float64

	nodes, err1 := os.Open("nodes.txt")
	sections, err2 := os.Open("sections.txt")
	elements, err3 := os.Open("elements.txt")
	loads, err4 := os.Open("loads.txt")

	if err1 != nil && err2 != nil && err3 != nil && err4 != nil {
		panic("error reading input files")
	}

	scanner := bufio.NewScanner(nodes)

	for scanner.Scan() {
		data := strings.Fields(scanner.Text())
		coords := [2]float64{gog.Must(strconv.ParseFloat(data[0], 64)), gog.Must(strconv.ParseFloat(data[1], 64))}
		components.MakeNode(coords[0], coords[1],
					gog.Must(strconv.ParseBool(data[2])), gog.Must(strconv.ParseBool(data[3])))
		if coords[0] < min[0] {
			min[0] = coords[0]
		}
		if coords[1] < min[1] {
 			min[1] = coords[1]
 		} 
 		if coords[0] > max[0] {
 			max[0] = coords[0]
 		}
 		if coords[0] > max[1] {
   			max[1] = coords[1]
   		}  
		
	}

	scanner = bufio.NewScanner(sections)

	for scanner.Scan() {
		data := strings.Fields(scanner.Text())
		components.MakeSection(gog.Must(strconv.ParseFloat(data[0], 64)), gog.Must(strconv.ParseFloat(data[1], 64)))
	}

	scanner = bufio.NewScanner(elements)

	for scanner.Scan() {
		data := strings.Fields(scanner.Text())
		node1, _ := strconv.Atoi(data[0])
		node2, _ := strconv.Atoi(data[1])
		section, _ := strconv.Atoi(data[2])
		components.MakeElement(components.NodeList[node1], components.NodeList[node2], components.SectionList[section])
	}

	scanner = bufio.NewScanner(loads)

	for scanner.Scan() {
		data := strings.Fields(scanner.Text())
		node, _ := strconv.Atoi(data[0])
		fx, _ := strconv.ParseFloat(data[1], 64)
		fy, _ := strconv.ParseFloat(data[2], 64)
		components.ApplyPointLoad(components.NodeList[node], fx, fy)
	}

	nodes.Close()
	elements.Close()
	sections.Close()
	loads.Close()

	displacements, supports := components.Solve()

	file, err := os.OpenFile("results.txt", os.O_RDWR | os.O_CREATE | os.O_TRUNC, 0755)

	if err != nil {
		fmt.Println("error writing results file")
	}

	writer := io.MultiWriter(os.Stdout, file)

	fmt.Fprintln(writer, "Truss forces and Stresses")

	for i := 0; i < len(components.ElementList); i++ {
		force := components.ElementList[i].P
		stress := force / components.ElementList[i].A
		fmt.Fprintf(writer, "Element %d: %5E, %5E\n", i, force, stress)
	}	

	fmt.Fprintln(writer, "\nNodal Displacements in X and Y")

	for i := 0; i < len(components.NodeList); i++ {
		dx := 0.0
		dy := 0.0
		if(!components.NodeList[i].FixedX) {
			dx = displacements.AtVec(components.NodeList[i].XDeg)
		}
		if(!components.NodeList[i].FixedY) {
			dy = displacements.AtVec(components.NodeList[i].YDeg)
		}

		fmt.Fprintf(writer, "Node %d: %5E, %5E\n", i, dx, dy)	
	}

	fmt.Fprintln(writer, "\nSupport Reactions")

	for node, i := range components.FixedNodes {
		fmt.Fprintf(writer, "Node %d: %5E, %5E \n", i, supports.AtVec(node.XDeg - displacements.Len()), supports.AtVec(node.YDeg - displacements.Len()))
	}

	dim := [2]int{1000, 1000}
	

	pre := gg.NewContext(dim[0], dim[1])
	post := gg.NewContext(dim[0], dim[1])
	
	pre.InvertY()
	post.InvertY()
	pre.SetRGB(1, 1, 1)
	post.SetRGB(1, 1, 1)
	pre.Clear()
	post.Clear()

	mid := [2]float64{0.5*(min[0]+max[0]), 0.5*(min[1]+max[1])}
	delta := math.Max(max[0] - mid[0], max[1] - mid[1])

	k := (float64(dim[0]) / 2 - 100) / delta

	pre.SetRGB(0, 0, 0)
	post.SetRGB(0, 0, 0)

	post.SetLineWidth(3)
	post.Stroke()

	for i := 0; i < len(components.NodeList); i++ {
		node := components.NodeList[i]
		coords := [2]float64{float64(dim[0])/2 + (node.X - mid[0])*k, float64(dim[0])/2 + (node.Y - mid[1])*k}
		pre.DrawCircle(coords[0], coords[1], 4)
		pre.Fill()
		if !node.FixedX {
			exp := math.Pow10(1+int(math.Round(-math.Log10(math.Abs(displacements.AtVec(node.XDeg))))))
			coords[0] += displacements.AtVec(node.XDeg) * exp
		}
		if !node.FixedY {
			exp := math.Pow10(1+int(math.Round(-math.Log10(math.Abs(displacements.AtVec(node.YDeg))))))
			coords[1] += displacements.AtVec(node.YDeg) * exp
		}
		if node.FixedX {
			post.SetRGB(0, 1, 0)
			if node.FixedY {
				post.DrawRegularPolygon(3, coords[0], coords[1] - 10, 10, -math.Pi)
			}else {
				post.DrawCircle(coords[0] - 10, coords[1], 10)
			}
			post.Stroke()
			post.SetRGB(0, 0, 0)
		}else if node.FixedY {
			post.SetRGB(0, 1, 0)
			post.DrawCircle(coords[0], coords[1] - 10, 10)
			post.Stroke()
			post.SetRGB(0, 0, 0)
		}
		
		if node.Loaded {
			post.MoveTo(coords[0], coords[1])
			post.SetRGB(0, 1, 0)
			angle := math.Atan2(node.Load[1], node.Load[0])
			dx := 50 * math.Cos(angle)
			dy := 50 * math.Sin(angle)
			post.LineTo(coords[0] + dx, coords[1] + dy)
			post.Stroke()
			post.DrawRegularPolygon(3, coords[0] + dx, coords[1] + dy, 15, -angle)
			post.Fill()
			post.MoveTo(0, 0)
			post.SetRGB(0, 0, 0)
		}
		post.DrawCircle(coords[0], coords[1], 4)
		post.Fill()
	}

	pre.SetLineWidth(3)
	pre.Stroke()

	post.SetLineWidth(3)
	post.Stroke()

	for i := 0; i < len(components.ElementList); i++ {
		element := components.ElementList[i]
		coords1 := [2]float64{float64(dim[0])/2 + (element.N1.X - mid[0])*k, float64(dim[0])/2 + (element.N1.Y - mid[1])*k}
		coords2 := [2]float64{float64(dim[0])/2 + (element.N2.X - mid[0])*k, float64(dim[0])/2 + (element.N2.Y - mid[1])*k}

		pre.MoveTo(coords1[0], coords1[1])
		pre.LineTo(coords2[0], coords2[1])
		pre.Stroke()

		post.SetRGB(0, 0, 0)

		if element.P > 0 {
			post.SetRGB(1, 0, 0)
		}

		if element.P < 0 {
			post.SetRGB(0, 0, 1)
		}

		if !element.N1.FixedX {
			exp := math.Pow10(1+int(-math.Round(math.Log10(math.Abs(displacements.AtVec(element.N1.XDeg))))))
			coords1[0] += displacements.AtVec(element.N1.XDeg) * exp
		}
		if !element.N1.FixedY {
			exp := math.Pow10(1+int(-math.Round(math.Log10(math.Abs(displacements.AtVec(element.N1.YDeg))))))
			coords1[1] += displacements.AtVec(element.N1.YDeg) * exp
		}
		if !element.N2.FixedX {
			exp := math.Pow10(1+int(-math.Round(math.Log10(math.Abs(displacements.AtVec(element.N2.XDeg))))))
			coords2[0] += displacements.AtVec(element.N2.XDeg) * exp
		}
		if !element.N2.FixedY {
			exp := math.Pow10(1+int(-math.Round(math.Log10(math.Abs(displacements.AtVec(element.N2.YDeg))))))
			coords2[1] += displacements.AtVec(element.N2.YDeg) * exp
		}

		post.MoveTo(coords1[0], coords1[1])
		post.LineTo(coords2[0], coords2[1])
		post.Stroke()
	}
	
    pre.SavePNG("pre.png")
    post.SavePNG("post.png")

}