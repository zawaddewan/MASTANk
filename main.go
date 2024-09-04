package main

import (
	"github.com/icza/gog"
	"MASTANk/components"
	"strings"
	"strconv"
	"bufio"
	"fmt"
	"os"
	"io"
)


func main() {
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
		components.MakeNode(gog.Must(strconv.ParseFloat(data[0], 64)), gog.Must(strconv.ParseFloat(data[1], 64)), 
		gog.Must(strconv.ParseBool(data[2])), gog.Must(strconv.ParseBool(data[3])))
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
}