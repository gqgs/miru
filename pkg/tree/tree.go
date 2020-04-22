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
		data, err := msgpack.Marshal(img)
		if err != nil {
			return err
		}
		_, err = t.db.Exec(`INSERT INTO tree (image0) VALUES (?)`, data)
		return err
	case nil:
		if image0 == nil {
			data, err := msgpack.Marshal(img)
			if err != nil {
				return err
			}
			_, err = t.db.Exec(`UPDATE tree SET image0 = ?`, data)
			return err
		}
		if image1 == nil {
			data, err := msgpack.Marshal(img)
			if err != nil {
				return err
			}
			_, err = t.db.Exec(`UPDATE tree SET image1 = ?`, data)
			return err
		}
		var dbImage0, dbImage1 image.Image
		if err = msgpack.Unmarshal(*image0, &dbImage0); err != nil {
			return err
		}
		if err = msgpack.Unmarshal(*image1, &dbImage1); err != nil {
			return err
		}
		// cmp1 := image.Compare(img, &dbImage0)
		// cmp2 := image.Compare(img, &dbImage1)
		// fmt.Println("cmp1", cmp1, img, dbImage0)
		// fmt.Println("cmp2", cmp2, img, dbImage1)
	}
	return err
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
