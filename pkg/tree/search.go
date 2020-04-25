package tree

import (
	"math"
	"sync"
)

// Search recursively traversals the tree to find the
// images the most ressemble the input image
func (t *Tree) Search(comparer Comparer, accuracy uint) (results, error) {
	return t.search(1, comparer, accuracy)
}

func (t *Tree) search(nodeID int64, comparer Comparer, accuracy uint) (results, error) {
	node, err := t.storage.Get(nodeID)
	if err != nil {
		return nil, err
	}
	if node.LeftChild == nil && node.RightChild == nil {
		var res results
		for _, data := range []*[]byte{node.LeftObject, node.RightObject} {
			if data != nil {
				cmp, filename, err := comparer.Compare(*data)
				if err != nil {
					return nil, err
				}
				res = append(res, result{
					Filename: filename,
					Score:    cmp,
				})
			}
		}
		return res, nil
	}
	// invariant: node has 2 elements here
	cmp0, filename0, err := comparer.Compare(*node.LeftObject)
	if err != nil {
		return nil, err
	}
	cmp1, filename1, err := comparer.Compare(*node.RightObject)
	if err != nil {
		return nil, err
	}
	imagesAreDissimilar := math.Abs(cmp0-cmp1) >= float64(accuracy)
	if imagesAreDissimilar {
		if cmp0 < cmp1 {
			if node.LeftChild == nil {
				return results{result{
					Filename: filename0,
					Score:    cmp0,
				}}, nil
			}
			res, err := t.search(*node.LeftChild, comparer, accuracy)
			if err != nil {
				return nil, err
			}
			res = append(res, result{
				Filename: filename0,
				Score:    cmp0,
			})
			return res, nil
		}
		if node.RightChild == nil {
			return results{result{
				Filename: filename1,
				Score:    cmp1,
			}}, nil
		}
		res, err := t.search(*node.RightChild, comparer, accuracy)
		if err != nil {
			return nil, err
		}
		res = append(res, result{
			Filename: filename1,
			Score:    cmp1,
		})

		return res, nil
	}

	var res = results{result{
		Filename: filename0,
		Score:    cmp0,
	}, result{
		Filename: filename1,
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
		res0, err0 = t.search(*node.LeftChild, comparer, accuracy)
	}()
	go func() {
		defer wg.Done()
		if node.RightChild == nil {
			return
		}
		res1, err1 = t.search(*node.RightChild, comparer, accuracy)
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
