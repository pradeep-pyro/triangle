package triangle

/*
#include "triangle.h"
#include <stdlib.h>
*/
// #cgo LDFLAGS: -lm
import "C"
import (
	"unsafe"
)

// The allocations here are using C.malloc and C.free over trimalloc and
// trifree, as C.malloc works the same but gives a better error report in cases
// if we go out of memory.

func trimalloc(size C.ulong) unsafe.Pointer {
	return C.malloc(size)
}

func trifree(ptr unsafe.Pointer) {
	C.free(ptr)
}

type triangulateIO struct {
	ct *C.struct_triangulateio
}

func NewTriangulateIO() *triangulateIO {
	t := triangulateIO{}
	t.ct = (*C.struct_triangulateio)(trimalloc(C.sizeof_struct_triangulateio))
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
	trifree(unsafe.Pointer(t.ct.edgelist))
	trifree(unsafe.Pointer(t.ct.edgemarkerlist))

	// trifree(unsafe.Pointer(t.ct.holelist))
	// hole list is reused in both in and out of ConstrainedDelaunay

	trifree(unsafe.Pointer(t.ct.neighborlist))
	trifree(unsafe.Pointer(t.ct.normlist))
	trifree(unsafe.Pointer(t.ct.pointattributelist))
	trifree(unsafe.Pointer(t.ct.pointlist))
	trifree(unsafe.Pointer(t.ct.pointmarkerlist))
	trifree(unsafe.Pointer(t.ct.regionlist))
	trifree(unsafe.Pointer(t.ct.segmentlist))
	trifree(unsafe.Pointer(t.ct.segmentmarkerlist))
	trifree(unsafe.Pointer(t.ct.trianglearealist))
	trifree(unsafe.Pointer(t.ct.triangleattributelist))
	trifree(unsafe.Pointer(t.ct.trianglelist))
	trifree(unsafe.Pointer(t.ct))
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
	trifree(unsafe.Pointer(t.ct.edgelist))
	t.ct.edgelist, t.ct.numberofedges = intSlice2DToCArr(edges)
}

func (t *triangulateIO) SetPoints(pts [][2]float64) {
	trifree(unsafe.Pointer(t.ct.pointlist))
	t.ct.pointlist, t.ct.numberofpoints = flt64Slice2DToCArr(pts)
}

func (t *triangulateIO) SetPointMarkers(markers []int32) {
	trifree(unsafe.Pointer(t.ct.pointmarkerlist))
	t.ct.pointmarkerlist, _ = intSliceToCArr(markers)
}

func (t *triangulateIO) SetSegments(segments [][2]int32) {
	trifree(unsafe.Pointer(t.ct.segmentlist))
	t.ct.segmentlist, t.ct.numberofsegments = intSlice2DToCArr(segments)
}

func (t *triangulateIO) SetSegmentMarkers(markers []int32) {
	trifree(unsafe.Pointer(t.ct.segmentmarkerlist))
	t.ct.segmentmarkerlist, _ = intSliceToCArr(markers)
}

func (t *triangulateIO) SetTriangles(tri [][3]int32) {
	trifree(unsafe.Pointer(t.ct.trianglelist))
	t.ct.trianglelist, t.ct.numberoftriangles = intSlice3DToCArr(tri)
}

func (t *triangulateIO) SetTriangleAreas(areas []float64) {
	trifree(unsafe.Pointer(t.ct.trianglearealist))
	t.ct.trianglearealist, _ = flt64SliceToCArr(areas)
}

func (t *triangulateIO) SetHoles(holes [][2]float64) {
	if t.ct.holelist != nil {
		trifree(unsafe.Pointer(t.ct.holelist))
	}
	t.ct.holelist, t.ct.numberofholes = flt64Slice2DToCArr(holes)
}

func (t *triangulateIO) Triangles() [][3]int32 {
	return cArrToIntSlice3D(t.ct.trianglelist, t.NumberOfTriangles())
}

func triang(opt string, in, out, vorout *triangulateIO) {
	copt := C.CString(opt)
	defer trifree(unsafe.Pointer(copt))
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

func intSliceToCArr(slice []int32) (*C.int, C.int) {
	sz := len(slice)
	ptr := (*C.int)(trimalloc(C.sizeof_int * (C.ulong)(sz)))
	cArr := (*[1 << 30]C.int)(unsafe.Pointer(ptr))[:sz:sz]
	for i := range slice {
		cArr[i] = C.int(slice[i])
	}
	return ptr, C.int(len(slice))
}

func intSlice2DToCArr(slice [][2]int32) (*C.int, C.int) {
	sz := 2 * len(slice)
	ptr := (*C.int)(trimalloc(C.sizeof_int * (C.ulong)(sz)))
	cArr := (*[1 << 30]C.int)(unsafe.Pointer(ptr))[:sz:sz]
	for i := range slice {
		j := 2 * i
		cArr[j] = C.int(slice[i][0])
		cArr[j+1] = C.int(slice[i][1])
	}
	return ptr, C.int(len(slice))
}

func intSlice3DToCArr(slice [][3]int32) (*C.int, C.int) {
	sz := 3 * len(slice)
	ptr := (*C.int)(trimalloc(C.sizeof_int * (C.ulong)(sz)))
	cArr := (*[1 << 30]C.int)(unsafe.Pointer(ptr))[:sz:sz]
	for i := range slice {
		j := 3 * i
		cArr[j] = C.int(slice[i][0])
		cArr[j+1] = C.int(slice[i][1])
		cArr[j+2] = C.int(slice[i][2])
	}
	return ptr, C.int(len(slice))
}

func flt64SliceToCArr(slice []float64) (*C.double, C.int) {
	sz := len(slice)
	ptr := (*C.double)(trimalloc(C.sizeof_double * (C.ulong)(sz)))
	cArr := (*[1 << 30]C.double)(unsafe.Pointer(ptr))[:sz:sz]
	for i := range slice {
		cArr[i] = C.double(slice[i])
	}
	return ptr, C.int(len(slice))
}

func flt64Slice2DToCArr(slice [][2]float64) (*C.double, C.int) {
	sz := 2 * len(slice)
	ptr := (*C.double)(trimalloc(C.sizeof_double * (C.ulong)(sz)))
	cArr := (*[1 << 30]C.double)(unsafe.Pointer(ptr))[:sz:sz]
	for i := range slice {
		j := i * 2
		cArr[j] = C.double(slice[i][0])
		cArr[j+1] = C.double(slice[i][1])
	}
	return ptr, C.int(len(slice))
}
