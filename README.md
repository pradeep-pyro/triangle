# Go wrapper for triangle
This package is a wrapper for the [Triangle library] originally written in C by Jonathan Shewchuk. Features include

 - constrained & conforming Delaunay triangulation
 - Voronoi tesselation
 - triangle mesh generation with area and angle constraints

for points in 2D.

### Usage

Functions `Delaunay()` and `Voronoi()` accept a slice of 2D points represented as `[][2]float64`.

```go
// Points forming a spiral shape (https://www.cs.cmu.edu/~quake/spiral.node)
var pts = [][2]float64{{0, 0}, {-0.416, 0.909},
            {-1.35, 0.436}, {-1.64, -0.549},
            {-1.31, -1.51}, {-0.532, -2.17},
            {0.454, -2.41}, {1.45, -2.21},
            {2.29, -1.66}, {2.88, -0.838},
            {3.16, 0.131}, {3.12, 1.14},
            {2.77, 2.08}, {2.16, 2.89},
            {1.36, 3.49},
        }
```

`Delaunay()` returns a slice of triangle indices `[][3]int32`.

``` go
triangles := triangle.Delaunay(pts)
fmt.Println(triangles)
// Triangle indices
// [[0 4 5] [0 5 6] ... [7 8 0]]
```

`Voronoi()` returns four slices: the vertices `[][2]float64` and line segments (edges) `[][2]int32` of the Voronoi diagram, and the origins `[]int32` and directions `[][2]float64` of the edges that are infinite rays.

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

![spiralDelVor](http://i.imgur.com/WD7MO2l.png)

Constrained triangulation can be performed using the `Triangulate()` function.
For example, per-triangle angle and area constraints can be set as follows:

```go
in := triangle.NewTriangulateIO()
in.SetPoints(pts)
opt := triangle.NewOptions()
opt.Angle = 20
opt.Area = 15
out := triangle.Triangulate(in, opt, false)
// Vertices and faces can be obtained by calling out.Points()
// and out.Triangles()
```

![spiralConstrained](http://i.imgur.com/Nb2XRPX.png)

[Triangle library]: https://www.cs.cmu.edu/~quake/triangle.html
