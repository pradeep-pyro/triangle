# Go wrapper for triangle
This package is a wrapper for the [Triangle library] originally written in C by Jonathan Shewchuk. Features include

 - constrained & conforming Delaunay triangulation
 - Voronoi tesselation
 - triangle mesh generation with area and angle constraints

for points in 2D.

### Usage

#### Delaunay triangulation and Voronoi tesselation of points

Functions `Delaunay()` and `Voronoi()` accept a slice of 2D points represented as `[][2]float64`.

```go
// Points forming a spiral shape (https://www.cs.cmu.edu/~quake/spiral.node)
var pts = [][2]float64{{0, 0}, {-0.416, 0.909}, {-1.35, 0.436}, {-1.64, 0.549},
            {-1.31, -1.51}, {-0.532, -2.17}, {0.454, -2.41}, {1.45, -2.21},
            {2.29, -1.66}, {2.88, -0.838}, {3.16, 0.131}, {3.12, 1.14},
            {2.77, 2.08}, {2.16, 2.89}, {1.36, 3.49},
        }
```

`Delaunay()` returns a slice of triangle indices `[][3]int32` (left image below).
``` go
triangles := triangle.Delaunay(pts)
fmt.Println(triangles)
// Triangle indices
// [[0 4 5] [0 5 6] ... [7 8 0]]
```

`Voronoi()` returns four slices: the vertices `[][2]float64` and line segments (edges) `[][2]int32` of the Voronoi diagram, and the origins `[]int32` and directions `[][2]float64` of the edges that are infinite rays (right image below).
```go
vertices, edges, rayOrigins, rayDirections := triangle.Voronoi(pts)

// Vertices
// [[-0.2780 -1.0820] [0.2250 -1.2054] ... [1.1458 -0.8289]]
// Edges (pair of indices into the vertices slice)
// [[0 5] [0 1] ... [11 12]]
// Ray origins (index into vertices slice)
// [0 1 2 5 6 7 8 9 10 11 12 13 14]
// Direction vectors for rays in rayOrigins
// [[-0.660 -0.778] [-0.240 -0.986] ... [0.550 -0.840]]
```
![](http://i.imgur.com/WD7MO2l.png)

#### Constrained and conforming Delaunay triangulation of PSLGs
Functions `ConformingDelaunay()` and `ConstrainedDelaunay()` accept a PSLG as input.
PSLGs are planar straight line graphs defined by a set of points (`[][2]float64`), segments (pair of indices into the points slice `[][2]int32`), and holes ( single point somewhere within each hole `[][2]float64`).

```go
// Points forming the shape of letter "A"
var pts = [][2]float64{{0.200000, -0.776400}, {0.220000, -0.773200},
    {0.245600, -0.756400}, {0.277600, -0.702000}, {0.488800, -0.207600}, {0.504800, -0.207600}, {0.740800, -0.7396}, {0.756000, -0.761200},
    {0.774400, -0.7724}, {0.800000, -0.776400}, {0.800000, -0.792400}, {0.579200, -0.792400}, {0.579200, -0.776400}, {0.621600, -0.771600},
    {0.633600, -0.762800}, {0.639200, -0.744400}, {0.620800, -0.684400}, {0.587200, -0.604400}, {0.360800, -0.604400}, {0.319200, -0.706800},
    {0.312000, -0.739600}, {0.318400, -0.761200}, {0.334400, -0.771600}, {0.371200, -0.776400}, {0.371200, -0.792400}, {0.374400, -0.570000},
    {0.574400, -0.5700}, {0.473600, -0.330800}, {0.200000, -0.792400},
}
// Segments connecting the points
var segs = [][2]int32{{28, 0}, {0, 1}, {1, 2}, {2, 3}, {3, 4}, {4, 5}, {5, 6}, {6, 7}, {7, 8}, {8, 9}, {9, 10}, {10, 11}, {11, 12}, {12, 13}, {13, 14}, {14, 15}, {15, 16}, {16, 17}, {17, 18}, {18, 19}, {19, 20}, {20, 21}, {21, 22}, {22, 23}, {23, 24}, {24, 28}, {25, 26}, {26, 27}, {27, 25},
}
// Hole represented by a point lying inside it 
var holes = [][2]float64{
    {0.47, -0.5},
}
```

`ConstrainedDelaunay()` computes a constrained Delaunay triangulation of a PSLG, where the given segments are retained as such in the resulting triangulation. As a result, not all triangles are Delaunay. (left image below)
```go
verts, faces := triangle.ConstrainedDelaunay(pts, segs, holes)
```

`ConformingDelaunay()` computes a conforming Delaunay triangulation of a PSLG, where each triangle in the result is Delaunay. This is acheived by inserting vertices inbetween the given segments. (right image below)
```go
verts, faces := triangle.ConformingDelaunay(pts, segs, holes)
```

![](http://i.imgur.com/vmSHI2U.png)
![](http://i.imgur.com/9fh6cbW.png)

The function `Triangulate()` can be used for more fine grained control.
Example usage for per-triangle angle and area constraints:
```go
in := triangle.NewTriangulateIO()
in.SetPoints(pts)
opt := triangle.NewOptions()
opt.Angle = 20
opt.Area = 15
out := triangle.Triangulate(in, opt, false)
// Vertices and faces can be obtained by calling out.Points()
// and out.Triangles()
// Remember to free memory on the C side manually
triangle.FreeTriangulateIO(in)
triangle.FreeTriangulateIO(out)
```
![](http://i.imgur.com/Nb2XRPX.png)

[Triangle library]: https://www.cs.cmu.edu/~quake/triangle.html
