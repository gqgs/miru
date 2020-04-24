package tree

import (
	"math"
	"miru/pkg/image"
	"sync"
)

// Search recursively traversals the tree to find the
// images the most ressemble the input image
func (t *Tree) Search(img *image.Image, accuracy uint) (results, error) {
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

func (t *Tree) search(nodeID int, img *image.Image, accuracy uint) (results, error) {
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
		var res results
		for _, imgData := range []*[]byte{image0, image1} {
			if imgData != nil {
				if err := t.serializer.Unmarshal(*imgData, &dbImage); err != nil {
					return nil, err
				}
				cmp := image.Compare(img, &dbImage)
				res = append(res, result{
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
				return results{result{
					Filename: dbImage0.Filename,
					Score:    cmp0,
				}}, nil
			}
			res, err := t.search(*left, img, accuracy)
			if err != nil {
				return nil, err
			}
			res = append(res, result{
				Filename: dbImage0.Filename,
				Score:    cmp0,
			})
			return res, nil
		}
		if right == nil {
			return results{result{
				Filename: dbImage1.Filename,
				Score:    cmp1,
			}}, nil
		}
		res, err := t.search(*right, img, accuracy)
		if err != nil {
			return nil, err
		}
		res = append(res, result{
			Filename: dbImage1.Filename,
			Score:    cmp1,
		})

		return res, nil
	}

	var res = results{result{
		Filename: dbImage0.Filename,
		Score:    cmp0,
	}, result{
		Filename: dbImage1.Filename,
		Score:    cmp1,
	}}

	// TODO: check if access needs to be synced
	var res0, res1 results
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
