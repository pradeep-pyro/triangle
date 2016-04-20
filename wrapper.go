package triangle

/*
#cgo CFLAGS: -I.
#cgo LDFLAGS: ${SRCDIR}/libtriangle.a
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
	return cArrToFlt64Slice2D(unsafe.Pointer(t.ct.normlist), t.NumberOfEdges())
}

func (t *triangulateIO) Edges() [][2]int {
	return cArrToIntSlice2D(unsafe.Pointer(t.ct.edgelist), t.NumberOfEdges())
}

func (t *triangulateIO) Points() [][2]float64 {
	return cArrToFlt64Slice2D(unsafe.Pointer(t.ct.pointlist), t.NumberOfPoints())
}

func (t *triangulateIO) PointMarkers() []int {
	return cArrToIntSlice(unsafe.Pointer(t.ct.pointmarkerlist), t.NumberOfPoints())
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
	return cArrToIntSlice3D(unsafe.Pointer(t.ct.trianglelist), t.NumberOfTriangles())
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

func cArrToIntSlice(ptr unsafe.Pointer, length int) []int {
	slice := (*[1 << 30]C.int)(ptr)[:length:length]
	result := make([]int, length)
	for i := 0; i < length; i++ {
		result[i] = int(slice[i])
	}
	return result
}

func cArrToIntSlice2D(ptr unsafe.Pointer, length int) [][2]int {
	sz := length * 2
	slice := (*[1 << 30]C.int)(ptr)[:sz:sz]
	result := make([][2]int, length)
	for i := 0; i < length; i++ {
		j := i * 2
		result[i] = [2]int{int(slice[j]), int(slice[j+1])}
	}
	return result
}

func cArrToIntSlice3D(ptr unsafe.Pointer, length int) [][3]int {
	sz := length * 3
	slice := (*[1 << 30]C.int)(ptr)[:sz:sz]
	result := make([][3]int, length)
	for i := 0; i < length; i++ {
		j := i * 3
		result[i] = [3]int{int(slice[j]), int(slice[j+1]), int(slice[j+2])}
	}
	return result
}

func cArrToFlt64Slice(ptr unsafe.Pointer, length int) []float64 {
	slice := (*[1 << 30]C.double)(ptr)[:length:length]
	result := make([]float64, length)
	for i := 0; i < length; i++ {
		result[i] = float64(slice[i])
	}
	return result
}

func cArrToFlt64Slice2D(ptr unsafe.Pointer, length int) [][2]float64 {
	sz := length * 2
	slice := (*[1 << 30]C.double)(ptr)[:sz:sz]
	result := make([][2]float64, length)
	for i := 0; i < length; i++ {
		j := i * 2
		result[i] = [2]float64{float64(slice[j]), float64(slice[j+1])}
	}
	return result
}
