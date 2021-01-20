
CGO_CFLAGS="-march=native -Ofast -pipe"

install: search index plot

search:
	CGO_CFLAGS=${CGO_CFLAGS} go install ./cmd/miru-search/

index:
	CGO_CFLAGS=${CGO_CFLAGS} go install ./cmd/miru-index/

plot:
	CGO_CFLAGS=${CGO_CFLAGS} go install ./cmd/miru-plot/

docker:
	docker build -t miru .
