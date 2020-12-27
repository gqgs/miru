# miru
Performant content-based image retrieval using generalized hyperplane trees [1].

![miru](/assets/miru.png)

### Install

```
make install
```

#### Docker

```
make docker
docker run -v "$(pwd)":"$(pwd)" miru miru-insert --help
docker run -v "$(pwd)":"$(pwd)" miru miru-search --help
docker run -v "$(pwd)":"$(pwd)" miru miru-plot --help
```

### Usage

#### Index images

```
miru-insert --folder [folder]
```

#### Search for an image

```
miru-search --file [path|url]
```

#### Displaying the tree

```
miru-plot --out digraph.dot
dot -Tsvg digraph.dot -o miru.svg
```

### Performance 

`miru-search` should scale logarithmically according to the size of the database, assuming there is enough variety in the processed set of images.

---

[1] J. K. Uhlmann. Satisfying general proximity/similarity queries
with metric trees.


