package triangle

/*
#include "triangle.h"
#include <stdlib.h>
*/
// #cgo LDFLAGS: -lm
import "C"
import "unsafe"

type triangulateIO struct {
	ct *C.struct_triangulateio
}

func NewTriangulateIO() *triangulateIO {
	t := triangulateIO{}
	t.ct = (*C.struct_triangulateio)(C.malloc(C.sizeof_struct_triangulateio))
	if t.ct == nil {
		panic("Unable to allocate memory")
	}
	t.ct.edgelist = nil
	t.ct.edgemarkerlist = nil
	t.ct.holelist = nil
	t.ct.neighborlist = nil
	t.ct.normlist = nil
	t.ct.numberofcorners = 0
	t.ct.numberofedges = 0
	t.ct.numberofholes = 0
	t.ct.numberofpointattributes = 0
	t.ct.numberofpoints = 0
	t.ct.numberofregions = 0
	t.ct.numberofsegments = 0
	t.ct.numberoftriangleattributes = 0
	t.ct.numberoftriangles = 0
	t.ct.pointattributelist = nil
	t.ct.pointlist = nil
	t.ct.pointmarkerlist = nil
	t.ct.regionlist = nil
	t.ct.segmentlist = nil
	t.ct.segmentmarkerlist = nil
	t.ct.trianglearealist = nil
	t.ct.triangleattributelist = nil
	t.ct.trianglelist = nil
	return &t
}

func FreeTriangulateIO(t *triangulateIO) {
	C.free(unsafe.Pointer(t.ct))
}

func (t *triangulateIO) NumberOfEdges() int {
	return int(t.ct.numberofedges)
}

func (t *triangulateIO) NumberOfHoles() int {
	return int(t.ct.numberofholes)
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
	return cArrToFlt64Slice2D(t.ct.normlist, t.NumberOfEdges())
}

func (t *triangulateIO) Edges() [][2]int32 {
	return cArrToIntSlice2D(t.ct.edgelist, t.NumberOfEdges())
}

func (t *triangulateIO) Points() [][2]float64 {
	return cArrToFlt64Slice2D(t.ct.pointlist, t.NumberOfPoints())
}

func (t *triangulateIO) PointMarkers() []int32 {
	return cArrToIntSlice(t.ct.pointmarkerlist, t.NumberOfPoints())
}

func (t *triangulateIO) Segments() [][2]int32 {
	return cArrToIntSlice2D(t.ct.segmentlist, t.NumberOfSegments())
}

func (t *triangulateIO) SetEdges(edges [][2]int32) {
	t.ct.edgelist = (*C.int)(unsafe.Pointer(&edges[0][0]))
	t.ct.numberofedges = C.int(len(edges))
}

func (t *triangulateIO) SetPoints(pts [][2]float64) {
	t.ct.pointlist = (*C.double)(unsafe.Pointer(&pts[0]))
	t.ct.numberofpoints = C.int(len(pts))
}

func (t *triangulateIO) SetPointMarkers(markers []int32) {
	t.ct.pointmarkerlist = (*C.int)(unsafe.Pointer(&markers[0]))
}

func (t *triangulateIO) SetSegments(segments [][2]int32) {
	t.ct.segmentlist = (*C.int)(unsafe.Pointer(&segments[0][0]))
	t.ct.numberofsegments = C.int(len(segments))
}

func (t *triangulateIO) SetSegmentMarkers(markers []int32) {
	t.ct.segmentmarkerlist = (*C.int)(unsafe.Pointer(&markers[0]))
}

func (t *triangulateIO) SetTriangles(tri [][3]int32) {
	t.ct.trianglelist = (*C.int)(unsafe.Pointer(&tri[0][0]))
}

func (t *triangulateIO) SetTriangleAreas(areas []float64) {
	t.ct.trianglearealist = (*C.double)(unsafe.Pointer(&areas[0]))
}

func (t *triangulateIO) SetHoles(holes [][2]float64) {
	t.ct.holelist = (*C.double)(unsafe.Pointer(&holes[0][0]))
	t.ct.numberofholes = C.int(len(holes))
}

func (t *triangulateIO) Triangles() [][3]int32 {
	return cArrToIntSlice3D(t.ct.trianglelist, t.NumberOfTriangles())
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

func cArrToIntSlice(ptr *C.int, length int) []int32 {
	slice := (*[1 << 30]C.int)(unsafe.Pointer(ptr))[:length:length]
	result := make([]int32, length)
	for i := 0; i < length; i++ {
		result[i] = int32(slice[i])
	}
	return result
}

func cArrToIntSlice2D(ptr *C.int, length int) [][2]int32 {
	sz := length * 2
	slice := (*[1 << 30]C.int)(unsafe.Pointer(ptr))[:sz:sz]
	result := make([][2]int32, length)
	for i := 0; i < length; i++ {
		j := i * 2
		result[i] = [2]int32{int32(slice[j]), int32(slice[j+1])}
	}
	return result
}

func cArrToIntSlice3D(ptr *C.int, length int) [][3]int32 {
	sz := length * 3
	slice := (*[1 << 30]C.int)(unsafe.Pointer(ptr))[:sz:sz]
	result := make([][3]int32, length)
	for i := 0; i < length; i++ {
		j := i * 3
		result[i] = [3]int32{int32(slice[j]), int32(slice[j+1]), int32(slice[j+2])}
	}
	return result
}

func cArrToFlt64Slice2D(ptr *C.double, length int) [][2]float64 {
	sz := length * 2
	slice := (*[1 << 30]C.double)(unsafe.Pointer(ptr))[:sz:sz]
	result := make([][2]float64, length)
	for i := 0; i < length; i++ {
		j := i * 2
		result[i] = [2]float64{float64(slice[j]), float64(slice[j+1])}
	}
	return result
}
