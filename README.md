## MASTANk: 2D Planar Truss Analysis in Go

A small project inspired by [MASTAN2](https://www.mastan2.com/), a program for analysis truss and frame structures. 

### Dependencies:
[Gonum](https://www.gonum.org/), [gog](https://github.com/icza/gog)

### Usage:
MASTANk takes in input files `nodes.txt`, `elements.txt`, `loads.txt`and `sections.txt`

`nodes.txt` defines the nodes of the structure and is formatted like so:

    0.5 0.5 false false
    0 0 true true
    0 0.5 true true
    0.5 0 false false
The two entries on each line are the (x, y) coordinates of the node and the last two define the freedom of the node (fixed in the x and/or y direction)

`elements.txt` defines truss elements of the structure:

    0 1 0
    0 2 0
    0 3 0
    1 2 0
    1 3 0
    2 3 0
Where the first two entries of the line are the nodes of the truss (0 is the first node defined in `nodes.txt`, 1 is the second, etc.) and the third entry defines the section properties of the truss (0 is the first section defined in `sections.txt`, 1 is the second, etc.)

`sections.txt` defines the section properties used in the structure:

    200e9 1.746e-4
Where the first entry is the modulus of the material and the second is the cross-sectional area of the section.

`loads.txt` defines the loading of the structure:

    0 1 -1
Where the first entry is a node defined in `nodes.txt` and the following entries are the (x,y) components of the force applied on the node.

After making the setup files, running the program will output the forces within each truss, nodal displacements, and the support forces:

    Element 0: -6.254251E-01, -3.582045E+03
    Element 1: 1.442242E+00, 8.260265E+03
    Element 2: -5.577577E-01, -3.194489E+03
    Element 3: 0.000000E+00, 0.000000E+00
    Element 4: -5.577577E-01, -3.194489E+03
    Element 5: 7.887885E-01, 4.517689E+03
    
    Nodal Displacements in X and Y
    Node 0: 2.065066E-08, -3.856089E-08
    Node 1: 0.000000E+00, 0.000000E+00
    Node 2: 0.000000E+00, 0.000000E+00
    Node 3: -7.986221E-09, -3.057467E-08
    
    Support Reactions
    Node 1: 1.000000E+00, 4.422423E-01
    Node 2: -2.000000E+00, 5.577577E-01
   
If the structure is unstable, the matrix solver underlying the program will fail and the message `Matrix is singular - free body motion detected` will be printed. 

### Todo 

 - Draw images of truss before and after loading, labelled with the forces in each truss
 - Add support for distributed loads
 - Venture into the third dimension
 - Add support for frames in two and three dimensions
