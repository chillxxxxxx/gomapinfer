package common

import (
	"github.com/dhconnelly/rtreego"
	"math"
)

// RtreegoRect 将 Rectangle 转换为 rtreego.Rect（值类型）
func RtreegoRect(r Rectangle) rtreego.Rect {
	dx := math.Max(0.00000001, r.Max.X-r.Min.X)
	dy := math.Max(0.00000001, r.Max.Y-r.Min.Y)
	rect, err := rtreego.NewRect(rtreego.Point{r.Min.X, r.Min.Y}, []float64{dx, dy})
	if err != nil {
		panic(err)
	}
	return rect // 修改为返回 rtreego.Rect 而不是 *rtreego.Rect
}

type edgeSpatial struct {
	edge *Edge
	rect rtreego.Rect // 修改为值类型存储
}

func (e *edgeSpatial) Bounds() rtreego.Rect { // 修改返回类型为 rtreego.Rect
	if e.rect == (rtreego.Rect{}) { // 检查是否为空矩形
		r := e.edge.Src.Point.Rectangle()
		r = r.Extend(e.edge.Dst.Point)
		e.rect = RtreegoRect(r)
	}
	return e.rect // 返回值类型
}

type Rtree struct {
	tree *rtreego.Rtree
}

func (rtree Rtree) Search(rect Rectangle) []*Edge {
	spatials := rtree.tree.SearchIntersect(RtreegoRect(rect)) // 直接传入值类型
	edges := make([]*Edge, len(spatials))
	for i := range spatials {
		edges[i] = spatials[i].(*edgeSpatial).edge
	}
	return edges
}

func (graph *Graph) Rtree() Rtree {
	rtree := rtreego.NewTree(2, 25, 50)
	for _, edge := range graph.Edges {
		rtree.Insert(&edgeSpatial{edge: edge})
	}
	return Rtree{rtree}
}
