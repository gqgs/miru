package tree

import (
	"database/sql"
	"miru/pkg/image"
)

// Add recursively traversals the tree to find the
// correct insert position for the image
func (t *Tree) Add(img *image.Image) error {
	t.mu.Lock()
	defer t.mu.Unlock()

	var err error
	t.stmt, err = t.db.Prepare(
		`SELECT *
		FROM tree
		WHERE id = ?
		`)
	if err != nil {
		return err
	}
	defer t.stmt.Close()

	return t.add(1, img)
}

func (t *Tree) add(nodeID int, img *image.Image) error {
	var (
		id     int
		image0 *[]byte
		image1 *[]byte
		left   *int
		right  *int
	)

	err := t.stmt.QueryRow(nodeID).Scan(&id, &image0, &image1, &left, &right)
	switch err {
	case sql.ErrNoRows:
		_, err = t.createNode(img)
		return err
	case nil:
		if image0 == nil {
			data, err := t.serializer.Marshal(img)
			if err != nil {
				return err
			}
			_, err = t.db.Exec(`UPDATE tree SET image0 = ? WHERE id = ?`, data, nodeID)
			return err
		}
		if image1 == nil {
			data, err := t.serializer.Marshal(img)
			if err != nil {
				return err
			}
			_, err = t.db.Exec(`UPDATE tree SET image1 = ? WHERE id = ?`, data, nodeID)
			return err
		}
		var dbImage0, dbImage1 image.Image
		if err = t.serializer.Unmarshal(*image0, &dbImage0); err != nil {
			return err
		}
		if err = t.serializer.Unmarshal(*image1, &dbImage1); err != nil {
			return err
		}
		cmp0 := image.Compare(img, &dbImage0)
		cmp1 := image.Compare(img, &dbImage1)
		if cmp0 < cmp1 {
			if left == nil {
				lastID, err := t.createNode(img)
				if err != nil {
					return err
				}
				_, err = t.db.Exec("UPDATE tree SET left = ? WHERE id = ?", lastID, nodeID)
				return err
			}
			return t.add(*left, img)
		}
		if right == nil {
			lastID, err := t.createNode(img)
			if err != nil {
				return err
			}
			_, err = t.db.Exec("UPDATE tree SET right = ? WHERE id = ?", lastID, nodeID)
			return err
		}
		return t.add(*right, img)
	}
	return err
}

func (t *Tree) createNode(img *image.Image) (int64, error) {
	data, err := t.serializer.Marshal(img)
	if err != nil {
		return 0, err
	}
	result, err := t.db.Exec(`INSERT INTO tree (image0) VALUES (?)`, data)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}
