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
func Delaunay(pts [][2]float64) [][3]int32 {
	in := NewTriangulateIO()
	out := NewTriangulateIO()
	defer FreeTriangulateIO(in)
	defer FreeTriangulateIO(out)

	in.SetPoints(pts)
	in.SetPointMarkers(make([]int32, len(pts)))

	triang("Qz", in, out, nil)

	triangles := out.Triangles()

	return triangles
}

// Voronoi computes the Voronoi diagram of a given set of points
// It returns a set of Voronoi vertices, a set of edges between the points, as well
// as infinite which can occur around the boundary (defined by rayOrigins which indexes
// into the vertices, and rayDirs which provides the direction).
func Voronoi(pts [][2]float64) ([][2]float64, [][2]int32, []int32, [][2]float64) {
	in := NewTriangulateIO()
	out := NewTriangulateIO()
	vorout := NewTriangulateIO()
	defer FreeTriangulateIO(in)
	defer FreeTriangulateIO(out)
	defer FreeTriangulateIO(vorout)

	in.SetPoints(pts)
	in.SetPointMarkers(make([]int32, len(pts)))

	triang("Qzv", in, out, vorout)

	verts := vorout.Points()
	e := vorout.Edges()
	dir := vorout.Normals()

	edges := make([][2]int32, 0, len(e))
	rayOrigins := make([]int32, 0, len(e))
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

// ConformingDelaunay computes the true Delaunay triangulation of a planar
// straight line graph with the given vertices, edges and holes.
// New vertices (Steiner points) may be inserted to ensure that the resulting
// triangles are all Delaunay.
func ConformingDelaunay(pts [][2]float64, segs [][2]int32,
	holes [][2]float64) ([][2]float64, [][3]int32) {
	in := NewTriangulateIO()
	out := NewTriangulateIO()
	defer FreeTriangulateIO(in)
	defer FreeTriangulateIO(out)

	in.SetPoints(pts)
	in.SetPointMarkers(make([]int32, len(pts)))
	in.SetSegments(segs)
	in.SetSegmentMarkers(make([]int32, len(segs)))
	in.SetHoles(holes)
	triang("QzpD", in, out, nil)

	verts := out.Points()
	triangles := out.Triangles()
	return verts, triangles
}

// Triangulate is the closest wrapper to the C code and can be used for
// flexible needs. Flags, constraints and quality options can be set using
// the second argument.
// Holes and segments that must appear in the triangulation can be set
// using methods (SetSegments() and SetHoles()) in the input triangulateIO
// struct.
//
// Note that FreeTriangulateIO() has to be called explicitly on the in and out
// to release the memory.
func Triangulate(in *triangulateIO, opts *options, verbose bool) *triangulateIO {
	out := NewTriangulateIO()

	optsStr := optsToString(opts)
	if !verbose {
		optsStr += "Q"
	}
	if in.NumberOfSegments() > 0 || in.NumberOfHoles() > 0 {
		optsStr += "p"
	}

	triang(optsStr, in, out, nil)

	return out
}
