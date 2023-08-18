package tree

import (
	"math"
	"sync"
)

// Search recursively traversals the tree to find the
// elements the most ressemble the input element
func (t *Tree) Search(comparer Comparer, accuracy uint) (results, error) {
	return t.search(1, comparer, accuracy)
}

func (t *Tree) search(nodeID int64, comparer Comparer, accuracy uint) (results, error) {
	node, err := t.storage.Get(nodeID)
	if err != nil {
		return nil, err
	}
	if !node.LeftChild.Valid && !node.RightChild.Valid {
		var res results
		for _, data := range [][]byte{node.LeftObject, node.RightObject} {
			if data != nil {
				cmp, filename, err := comparer.Compare(data)
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
	cmp0, filename0, err := comparer.Compare(node.LeftObject)
	if err != nil {
		return nil, err
	}
	cmp1, filename1, err := comparer.Compare(node.RightObject)
	if err != nil {
		return nil, err
	}
	imagesAreDissimilar := math.Abs(cmp0-cmp1) >= float64(accuracy)
	if imagesAreDissimilar {
		if cmp0 < cmp1 {
			if !node.LeftChild.Valid {
				return results{result{
					Filename: filename0,
					Score:    cmp0,
				}}, nil
			}
			res, err := t.search(node.LeftChild.Int64, comparer, accuracy)
			if err != nil {
				return nil, err
			}
			res = append(res, result{
				Filename: filename0,
				Score:    cmp0,
			})
			return res, nil
		}
		if !node.RightChild.Valid {
			return results{result{
				Filename: filename1,
				Score:    cmp1,
			}}, nil
		}
		res, err := t.search(node.RightChild.Int64, comparer, accuracy)
		if err != nil {
			return nil, err
		}
		res = append(res, result{
			Filename: filename1,
			Score:    cmp1,
		})

		return res, nil
	}

	var res0, res1 results
	var err0, err1 error
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		if !node.LeftChild.Valid {
			return
		}
		res0, err0 = t.search(node.LeftChild.Int64, comparer, accuracy)
	}()
	go func() {
		defer wg.Done()
		if !node.RightChild.Valid {
			return
		}
		res1, err1 = t.search(node.RightChild.Int64, comparer, accuracy)
	}()
	wg.Wait()
	if err0 != nil {
		return nil, err0
	}
	if err1 != nil {
		return nil, err1
	}

	res := make(results, 2+len(res0)+len(res1))
	i := copy(res, results{result{
		Filename: filename0,
		Score:    cmp0,
	}, result{
		Filename: filename1,
		Score:    cmp1,
	}})
	i += copy(res[i:], res0)
	copy(res[i:], res1)

	return res, nil
}
