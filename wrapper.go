package triangle

/*
#cgo CFLAGS: -I. -w
#cgo LDFLAGS: -L${SRCDIR} -ltriangle
#include "triangle.h"
#include <stdlib.h>
*/
import "C"
import "unsafe"

type triangulateIO struct {
	ct *C.struct_triangulateio
}

func newTriangulateIO() *triangulateIO {
	t := triangulateIO{}
	t.ct = (*C.struct_triangulateio)(C.malloc(C.sizeof_struct_triangulateio))
	if t.ct == nil {
		panic("Unable to allocate memory")
	}
	return &t
}

func freeTriangulateIO(t *triangulateIO) {
	C.free(unsafe.Pointer(t.ct))
}

func (t *triangulateIO) NumberOfEdges() int {
	return int(t.ct.numberofedges)
}

func (t *triangulateIO) NumberOfPoints() int {
	return int(t.ct.numberofpoints)
}

func (t *triangulateIO) NumberOfSegments() int {
	return int(t.ct.numberofsegments)
}

func (t *triangulateIO) NumberOfTriangles() int {
	return int(t.ct.numberoftriangles)
}

func (t *triangulateIO) Normals() [][2]float64 {
	num := t.NumberOfEdges()
	sz := num * 2
	slice := (*[1 << 30]C.double)(unsafe.Pointer(t.ct.normlist))[:sz:sz]

	result := make([][2]float64, num)
	for i := 0; i < num; i++ {
		j := i * 2
		result[i] = [2]float64{float64(slice[j]), float64(slice[j+1])}
	}
	return result
}

func (t *triangulateIO) Edges() [][2]int {
	numEdges := t.NumberOfEdges()
	sz := numEdges * 2
	slice := (*[1 << 30]C.int)(unsafe.Pointer(t.ct.edgelist))[:sz:sz]

	result := make([][2]int, numEdges)
	for i := 0; i < numEdges; i++ {
		j := i * 2
		result[i] = [2]int{int(slice[j]), int(slice[j+1])}
	}
	return result
}

func (t *triangulateIO) Points() [][2]float64 {
	numPoints := t.NumberOfPoints()
	sz := numPoints * 2
	slice := (*[1 << 30]C.double)(unsafe.Pointer(t.ct.pointlist))[:sz:sz]

	result := make([][2]float64, numPoints)
	for i := 0; i < numPoints; i++ {
		j := i * 2
		result[i] = [2]float64{float64(slice[j]), float64(slice[j+1])}
	}
	return result
}

func (t *triangulateIO) PointMarkers() []int {
	sz := t.NumberOfPoints()
	slice := (*[1 << 30]C.int)(unsafe.Pointer(t.ct.pointmarkerlist))[:sz:sz]

	result := make([]int, sz)
	for i := 0; i < sz; i++ {
		result[i] = int(slice[i])
	}
	return result
}

func (t *triangulateIO) SetEdges(edges [][2]int) {
	t.ct.edgelist = (*C.int)(unsafe.Pointer(&edges[0][0]))
	t.ct.numberofedges = C.int(len(edges))
}

func (t *triangulateIO) SetPoints(pts [][2]float64) {
	t.ct.pointlist = (*C.double)(unsafe.Pointer(&pts[0][0]))
	t.ct.numberofpoints = C.int(len(pts))
}

func (t *triangulateIO) SetPointMarkers(markers []int) {
	t.ct.pointmarkerlist = (*C.int)(unsafe.Pointer(&markers[0]))
}

func (t *triangulateIO) SetSegments(segments [][2]int) {
	t.ct.segmentlist = (*C.int)(unsafe.Pointer(&segments[0]))
}

func (t *triangulateIO) SetSegmentMarkers(markers [][2]int) {
	t.ct.segmentmarkerlist = (*C.int)(unsafe.Pointer(&markers[0]))
}

func (t *triangulateIO) SetTriangles(tri [][3]int) {
	t.ct.trianglelist = (*C.int)(unsafe.Pointer(&tri[0][0]))
}

func (t *triangulateIO) SetTriangleAreas(areas []float64) {
	t.ct.trianglearealist = (*C.double)(unsafe.Pointer(&areas[0]))
}

func (t *triangulateIO) Triangles() [][3]int {
	numTriangles := t.NumberOfTriangles()
	numCorners := int(t.ct.numberofcorners)
	sz := numTriangles * numCorners
	slice := (*[1 << 30]C.int)(unsafe.Pointer(t.ct.trianglelist))[:sz:sz]

	result := make([][3]int, numTriangles)
	for i := 0; i < numTriangles; i++ {
		j := i * numCorners
		result[i] = [3]int{int(slice[j]), int(slice[j+1]), int(slice[j+2])}
	}
	return result
}

func triang(opt string, in, out, vorout *triangulateIO) {
	copt := C.CString(opt)
	defer C.free(unsafe.Pointer(copt))
	if vorout == nil {
		C.triangulate(copt, in.ct, out.ct, nil)
	} else {
		C.triangulate(copt, in.ct, out.ct, vorout.ct)
	}
}
