package triangle

import "testing"

// To test if this library adheres to the passing pointers rules as mentioned in
// https://pkg.go.dev/cmd/cgo#hdr-Passing_pointers, run the test suite with the
// command
//
//     GODEBUG=cgocheck=2 go test -gcflags=all=-d=checkptr

// Points forming a spiral shape (https://www.cs.cmu.edu/~quake/spiral.node)
var delaunayVoronoiTestPts = [][2]float64{
	{0, 0}, {-0.416, 0.909}, {-1.35, 0.436}, {-1.64, 0.549},
	{-1.31, -1.51}, {-0.532, -2.17}, {0.454, -2.41}, {1.45, -2.21},
	{2.29, -1.66}, {2.88, -0.838}, {3.16, 0.131}, {3.12, 1.14},
	{2.77, 2.08}, {2.16, 2.89}, {1.36, 3.49},
}

func TestDelaunay(t *testing.T) {
	triangles := Delaunay(delaunayVoronoiTestPts)
	const numExpectedTriangles = 16
	assertEq(t, numExpectedTriangles, len(triangles),
		"Did not get expected number of triangles")
}

func TestVoronoi(t *testing.T) {
	vertices, edges, rayOrigins, rayDirections := Voronoi(delaunayVoronoiTestPts)

	const (
		numExpectedVertices      = 16
		numExpectedEdges         = numExpectedVertices + 2
		numExpectedRayOrigins    = 12
		numExpectedRayDirections = numExpectedRayOrigins
	)

	assertEq(t, numExpectedVertices, len(vertices),
		"Did not get expected number of vertices")
	assertEq(t, numExpectedEdges, len(edges),
		"Did not get expected number of edges")
	assertEq(t, numExpectedRayOrigins, len(rayOrigins),
		"Did not get expected number of ray origins")
	assertEq(t, numExpectedRayDirections, len(rayDirections),
		"Did not get expected number of ray directions")
}

func TestConstrainedConformingDelaunay(t *testing.T) {
	// Points forming the shape of letter "A"
	var pts = [][2]float64{{0.200000, -0.776400}, {0.220000, -0.773200},
		{0.245600, -0.756400}, {0.277600, -0.702000}, {0.488800, -0.207600}, {0.504800, -0.207600}, {0.740800, -0.7396}, {0.756000, -0.761200},
		{0.774400, -0.7724}, {0.800000, -0.776400}, {0.800000, -0.792400}, {0.579200, -0.792400}, {0.579200, -0.776400}, {0.621600, -0.771600},
		{0.633600, -0.762800}, {0.639200, -0.744400}, {0.620800, -0.684400}, {0.587200, -0.604400}, {0.360800, -0.604400}, {0.319200, -0.706800},
		{0.312000, -0.739600}, {0.318400, -0.761200}, {0.334400, -0.771600}, {0.371200, -0.776400}, {0.371200, -0.792400}, {0.374400, -0.570000},
		{0.574400, -0.5700}, {0.473600, -0.330800}, {0.200000, -0.792400},
	}
	// Segments connecting the points
	var segs = [][2]int32{{28, 0}, {0, 1}, {1, 2}, {2, 3}, {3, 4}, {4, 5}, {5, 6}, {6, 7}, {7, 8}, {8, 9}, {9, 10}, {10, 11}, {11, 12}, {12, 13}, {13, 14}, {14, 15}, {15, 16}, {16, 17}, {17, 18}, {18, 19}, {19, 20}, {20, 21}, {21, 22}, {22, 23}, {23, 24}, {24, 28}, {25, 26}, {26, 27}, {27, 25}}
	// Hole represented by a point lying inside it
	var holes = [][2]float64{
		{0.47, -0.5},
	}
	verts, faces := ConstrainedDelaunay(pts, segs, holes)

	const (
		numExpectedConstrainedVertices = 29
		numExpectedConstrainedFaces    = numExpectedConstrainedVertices
	)

	assertEq(t, numExpectedConstrainedVertices, len(verts),
		"ConstrainedDelaunay: Did not get expected number of vertices")
	assertEq(t, numExpectedConstrainedFaces, len(faces),
		"ConstrainedDelaunay: Did not get expected number of faces")

	const (
		numExpectedConformingVertices = 61
		numExpectedConformingFaces    = numExpectedConformingVertices
	)

	verts, faces = ConformingDelaunay(pts, segs, holes)

	assertEq(t, numExpectedConformingVertices, len(verts),
		"ConformingDelaunay: Did not get expected number of vertices")
	assertEq(t, numExpectedConformingFaces, len(faces),
		"ConformingDelaunay: Did not get expected number of faces")

}

func assertEq(t *testing.T, expected, actual int, msg string) {
	t.Helper()
	if expected != actual {
		t.Log(msg)
		t.Fatalf("Expected %d, but was %d", expected, actual)
	}
}
