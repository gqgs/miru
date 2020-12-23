# miru
Performant content-based image retrieval using generalized hyperplane trees [1].

### Install

```
make install
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

### Performance 

`miru-search` should scale logarithmically according to the size of the database, assuming there is enough variety in the processed set of images.

---

[1] J. K. Uhlmann. Satisfying general proximity/similarity queries
with metric trees.


