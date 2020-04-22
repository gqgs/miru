package tree

import (
	"miru/pkg/image"

	"github.com/vmihailenco/msgpack/v4"
)

type Result struct {
	Filename string
	Score    float64
}

type Results []Result

func (r Results) Len() int           { return len(r) }
func (r Results) Less(i, j int) bool { return r[i].Score < r[j].Score }
func (r Results) Swap(i, j int)      { r[i], r[j] = r[j], r[i] }

func (r *Results) Push(x interface{}) {
	*r = append(*r, x.(Result))
}

func (r *Results) Pop() interface{} {
	old := *r
	n := len(old)
	x := old[n-1]
	*r = old[0 : n-1]
	return x
}

// Search recursively traversals the tree to find the
// images the most ressemble the input image
func (t *Tree) Search(path string) (Results, error) {
	img, err := image.Load(path)
	if err != nil {
		return nil, err
	}

	t.stmt, err = t.db.Prepare(
		`SELECT *
		FROM tree
		WHERE id = ?
		`)
	if err != nil {
		return nil, err
	}
	defer t.stmt.Close()

	return t.search(1, img)
}

func (t *Tree) search(nodeID int, img *image.Image) (Results, error) {
	var (
		id     int
		image0 *[]byte
		image1 *[]byte
		left   *int
		right  *int
	)

	if err := t.stmt.QueryRow(nodeID).Scan(&id, &image0, &image1, &left, &right); err != nil {
		return nil, err
	}
	if left == nil && right == nil {
		var dbImage image.Image
		var res Results
		for _, imgData := range []*[]byte{image0, image1} {
			if imgData != nil {
				if err := msgpack.Unmarshal(*imgData, &dbImage); err != nil {
					return nil, err
				}
				cmp := image.Compare(img, &dbImage)
				res = append(res, Result{
					Filename: dbImage.Filename,
					Score:    cmp,
				})
			}
		}
		return res, nil
	}
	// invariant: node has 2 elements here
	var dbImage0, dbImage1 image.Image
	if err := msgpack.Unmarshal(*image0, &dbImage0); err != nil {
		return nil, err
	}
	if err := msgpack.Unmarshal(*image1, &dbImage1); err != nil {
		return nil, err
	}
	cmp0 := image.Compare(img, &dbImage0)
	cmp1 := image.Compare(img, &dbImage1)
	// TODO: handle accuracy: traverse both paths
	if cmp0 < cmp1 {
		if left == nil {
			return Results{Result{
				Filename: dbImage0.Filename,
				Score:    cmp0,
			}}, nil
		}
		res, err := t.search(*left, img)
		if err != nil {
			return nil, err
		}
		res = append(res, Result{
			Filename: dbImage0.Filename,
			Score:    cmp0,
		})
		return res, nil
	}
	if right == nil {
		return Results{Result{
			Filename: dbImage1.Filename,
			Score:    cmp1,
		}}, nil
	}
	res, err := t.search(*right, img)
	if err != nil {
		return nil, err
	}
	res = append(res, Result{
		Filename: dbImage1.Filename,
		Score:    cmp1,
	})

	return res, nil
}
