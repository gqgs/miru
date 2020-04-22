package tree

import (
	"database/sql"
	"miru/pkg/image"
	"miru/pkg/tree/internal/database"

	"github.com/vmihailenco/msgpack/v4"
)

type Tree struct {
	db   *sql.DB
	stmt *sql.Stmt
}

func (t *Tree) Close() error {
	return t.db.Close()
}

func (t *Tree) Add(path string) error {
	img, err := image.Load(path)
	if err != nil {
		return err
	}

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

// recursively traversals the tree to find the
// correct position for the image
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
			data, err := msgpack.Marshal(img)
			if err != nil {
				return err
			}
			_, err = t.db.Exec(`UPDATE tree SET image0 = ? WHERE id = ?`, data, nodeID)
			return err
		}
		if image1 == nil {
			data, err := msgpack.Marshal(img)
			if err != nil {
				return err
			}
			_, err = t.db.Exec(`UPDATE tree SET image1 = ? WHERE id = ?`, data, nodeID)
			return err
		}
		var dbImage0, dbImage1 image.Image
		if err = msgpack.Unmarshal(*image0, &dbImage0); err != nil {
			return err
		}
		if err = msgpack.Unmarshal(*image1, &dbImage1); err != nil {
			return err
		}
		cmp1 := image.Compare(img, &dbImage0)
		cmp2 := image.Compare(img, &dbImage1)
		if cmp1 < cmp2 {
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
	data, err := msgpack.Marshal(img)
	if err != nil {
		return 0, err
	}
	result, err := t.db.Exec(`INSERT INTO tree (image0) VALUES (?)`, data)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

func New(dbName string) (*Tree, error) {
	db, err := database.Open(dbName)
	if err != nil {
		return nil, err
	}
	return &Tree{
		db: db,
	}, nil
}
