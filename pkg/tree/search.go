package tree

import (
	"math"
	"miru/pkg/image"
	"sync"
)

// Search recursively traversals the tree to find the
// images the most ressemble the input image
func (t *Tree) Search(img *image.Image, accuracy uint) (results, error) {
	return t.search(1, img, accuracy)
}

func (t *Tree) search(nodeID int64, img *image.Image, accuracy uint) (results, error) {
	node, err := t.storage.Get(nodeID)
	if err != nil {
		return nil, err
	}
	if node.LeftChild == nil && node.RightChild == nil {
		var dbImage image.Image
		var res results
		for _, imgData := range []*[]byte{node.LeftObject, node.RightObject} {
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
	if err := t.serializer.Unmarshal(*node.LeftObject, &dbImage0); err != nil {
		return nil, err
	}
	if err := t.serializer.Unmarshal(*node.RightObject, &dbImage1); err != nil {
		return nil, err
	}
	cmp0 := image.Compare(img, &dbImage0)
	cmp1 := image.Compare(img, &dbImage1)
	imagesAreDissimilar := math.Abs(cmp0-cmp1) >= float64(accuracy)
	if imagesAreDissimilar {
		if cmp0 < cmp1 {
			if node.LeftChild == nil {
				return results{result{
					Filename: dbImage0.Filename,
					Score:    cmp0,
				}}, nil
			}
			res, err := t.search(*node.LeftChild, img, accuracy)
			if err != nil {
				return nil, err
			}
			res = append(res, result{
				Filename: dbImage0.Filename,
				Score:    cmp0,
			})
			return res, nil
		}
		if node.RightChild == nil {
			return results{result{
				Filename: dbImage1.Filename,
				Score:    cmp1,
			}}, nil
		}
		res, err := t.search(*node.RightChild, img, accuracy)
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

	var res0, res1 results
	var err0, err1 error
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		if node.LeftChild == nil {
			return
		}
		res0, err0 = t.search(*node.LeftChild, img, accuracy)
	}()
	go func() {
		defer wg.Done()
		if node.RightChild == nil {
			return
		}
		res1, err1 = t.search(*node.RightChild, img, accuracy)
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
