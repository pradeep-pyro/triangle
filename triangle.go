package triangle

import (
	"fmt"
	"strconv"
)

type SegmentSplitting uint8

const (
	SplittingAllowed SegmentSplitting = 0 + iota
	NoSplittingInBoundary
	NoSplitting
)

// options is a struct to hold parameters for controlling
// constraints in Constrained Delaunay Triangulation
type options struct {
	ConformingDelaunay bool
	EncloseConvexHull  bool
	SegmentSplitting   SegmentSplitting
	Area, Angle        float64
	MaxSteinerPoints   int
}

// NewOptions returns a new options struct with default parameters
func NewOptions() *options {
	return &options{ConformingDelaunay: true, EncloseConvexHull: false,
		SegmentSplitting: SplittingAllowed, Area: 15.0, Angle: 20.0,
		MaxSteinerPoints: -1}
}

// optsToString generates a string from the constraint and quality options
// in a format that can be passed to triangle
func optsToString(opts *options) string {
	str := fmt.Sprintf("zq%ga%g", opts.Angle, opts.Area)
	if opts.ConformingDelaunay {
		str += "D"
	}
	if opts.EncloseConvexHull {
		str += "c"
	}
	if opts.MaxSteinerPoints > -1 {
		str += "S" + strconv.Itoa(opts.MaxSteinerPoints)
	}

	if opts.SegmentSplitting == NoSplittingInBoundary {
		str += "Y"
	} else if opts.SegmentSplitting == NoSplitting {
		str += "YY"
	}
	return str
}

// Delaunay computes the unconstrained Delaunay triangulation of a given
// set of points
func Delaunay(pts [][2]float64) [][3]int {
	in := NewTriangulateIO()
	out := NewTriangulateIO()
	defer FreeTriangulateIO(in)
	defer FreeTriangulateIO(out)

	in.SetPoints(pts)
	in.SetPointMarkers(make([]int, len(pts)))

	triang("Qz", in, out, nil)

	triangles := out.Triangles()

	return triangles
}

// Voronoi computes the Voronoi diagram of a given set of points
// It returns a set of Voronoi vertices, a set of edges between the points, as well
// as infinite which can occur around the boundary (defined by rayOrigins which indexes
// into the vertices, and rayDirs which provides the direction).
func Voronoi(pts [][2]float64) ([][2]float64, [][2]int, []int, [][2]float64) {
	in := NewTriangulateIO()
	out := NewTriangulateIO()
	vorout := NewTriangulateIO()
	defer FreeTriangulateIO(in)
	defer FreeTriangulateIO(out)
	defer FreeTriangulateIO(vorout)

	in.SetPoints(pts)
	in.SetPointMarkers(make([]int, len(pts)))

	triang("Qzv", in, out, vorout)

	verts := vorout.Points()
	e := vorout.Edges()
	dir := vorout.Normals()

	edges := make([][2]int, 0, len(e))
	rayOrigins := make([]int, 0, len(e))
	rayDirs := make([][2]float64, 0, len(e))

	for i, val := range e {
		if val[1] != -1 {
			edges = append(edges, val)
		} else {
			rayOrigins = append(rayOrigins, val[0])
			rayDirs = append(rayDirs, dir[i])
		}
	}

	return verts, edges, rayOrigins, rayDirs
}

// Triangulate performs constrained Delaunay triangulation.
// Constraints and quality options can be set using the second argument.
// Holes and segments that must appear in the triangulation can be set using methods (SetSegments()
// and SetHoles()) in the input triangulateIO struct.
//
// Note that FreeTriangulateIO() has to called explicitly on the in and out to release the memory.
func Triangulate(in *triangulateIO, opts *options) *triangulateIO {
	out := NewTriangulateIO()

	optsStr := optsToString(opts)
	if in.NumberOfSegments() > 0 || in.NumberOfHoles() > 0 {
		optsStr += "p"
	}

	triang(optsStr, in, out, nil)

	return out
}
