package triangle

func Delaunay(pts [][2]float64) [][3]int {
	in := newTriangulateIO()
	out := newTriangulateIO()
	defer freeTriangulateIO(in)
	defer freeTriangulateIO(out)

	in.SetPoints(pts)
	in.SetPointMarkers(make([]int, len(pts)))

	triang("Qz", in, out, nil)

	triangles := out.Triangles()

	return triangles
}

func Voronoi(pts [][2]float64) ([][2]float64, [][2]int, []int, [][2]float64) {
	in := newTriangulateIO()
	out := newTriangulateIO()
	vorout := newTriangulateIO()
	defer freeTriangulateIO(in)
	defer freeTriangulateIO(out)
	defer freeTriangulateIO(vorout)

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
