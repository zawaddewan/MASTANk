package main

import (
	"MASTANk/components"
	"fmt"
	"gonum.org/v1/gonum/mat"
	"github.com/icza/gog"
	"strings"
	"strconv"
	"bufio"
	"os"
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

	displacements := components.Solve()

	fmt.Println(mat.Formatted(displacements, mat.Prefix(""), mat.Squeeze()))
	for i := 0; i < len(components.ElementList); i++ {
		fmt.Println(components.ElementList[i].P)
	}	
	
}