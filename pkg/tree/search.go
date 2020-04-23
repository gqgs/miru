package tree

import (
	"math"
	"miru/pkg/image"
	"sync"
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
func (t *Tree) Search(img *image.Image, accuracy int) (Results, error) {
	var err error
	t.stmt, err = t.db.Prepare(
		`SELECT *
		FROM tree
		WHERE id = ?
		`)
	if err != nil {
		return nil, err
	}
	defer t.stmt.Close()

	return t.search(1, img, accuracy)
}

func (t *Tree) search(nodeID int, img *image.Image, accuracy int) (Results, error) {
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
				if err := t.serializer.Unmarshal(*imgData, &dbImage); err != nil {
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
	if err := t.serializer.Unmarshal(*image0, &dbImage0); err != nil {
		return nil, err
	}
	if err := t.serializer.Unmarshal(*image1, &dbImage1); err != nil {
		return nil, err
	}
	cmp0 := image.Compare(img, &dbImage0)
	cmp1 := image.Compare(img, &dbImage1)
	imagesAreDissimilar := math.Abs(cmp0-cmp1) >= float64(accuracy)
	if imagesAreDissimilar {
		if cmp0 < cmp1 {
			if left == nil {
				return Results{Result{
					Filename: dbImage0.Filename,
					Score:    cmp0,
				}}, nil
			}
			res, err := t.search(*left, img, accuracy)
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
		res, err := t.search(*right, img, accuracy)
		if err != nil {
			return nil, err
		}
		res = append(res, Result{
			Filename: dbImage1.Filename,
			Score:    cmp1,
		})

		return res, nil
	}

	var res = Results{Result{
		Filename: dbImage0.Filename,
		Score:    cmp0,
	}, Result{
		Filename: dbImage1.Filename,
		Score:    cmp1,
	}}

	// TODO: check if access needs to be synced
	var res0, res1 Results
	var err0, err1 error
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		if left == nil {
			return
		}
		res0, err0 = t.search(*left, img, accuracy)
	}()
	go func() {
		defer wg.Done()
		if right == nil {
			return
		}
		res1, err1 = t.search(*right, img, accuracy)
	}()
	wg.Wait()
	if err0 != nil {
		return nil, err0
	}
	if err1 != nil {
		return nil, err1
	}
	res = append(res, res0...)
	res = append(res, res1...)

	return res, nil
}
