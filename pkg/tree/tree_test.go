package tree

import (
	"io/ioutil"
	"math/rand"
	"path/filepath"
	"testing"

	"github.com/gqgs/miru/pkg/cache"
	"github.com/gqgs/miru/pkg/compress"
	"github.com/gqgs/miru/pkg/image"
	"github.com/gqgs/miru/pkg/storage"
)

func Test_Tree(t *testing.T) {
	compressor, err := compress.NewCompressor("nop")
	if err != nil {
		t.Fatal(err)
	}
	storage, err := storage.NewSqliteStorage("file::memory:", compressor, cache.New(0))
	if err != nil {
		t.Fatal(err)
	}
	defer storage.Close()

	files, err := ioutil.ReadDir("./testdata/thumbs")
	if err != nil {
		t.Fatal(err)
	}

	tree := New(storage)
	for _, file := range files {
		img, err := image.Load(filepath.Join("./testdata/thumbs", file.Name()))
		if err != nil {
			t.Fatal(err)
		}
		if err := tree.Add(img); err != nil {
			t.Fatal(err)
		}
	}

	i := rand.Intn(len(files)) - 1
	query, err := image.Load(filepath.Join("./testdata/thumbs", files[i].Name()))
	if err != nil {
		t.Fatal(err)
	}
	results, err := tree.Search(query, 2)
	if err != nil {
		t.Fatal(err)
	}
	result := results.Top(1)[0]
	if query.Filename != result.Filename {
		t.Errorf("unexpected result: want: %s: got: %s", query.Filename, result.Filename)
	}
}
