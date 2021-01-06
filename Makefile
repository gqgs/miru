
CGO_CFLAGS="-march=native -Ofast -pipe"

install: search insert plot

search:
	CGO_CFLAGS=${CGO_CFLAGS} go install ./cmd/miru-search/

insert:
	CGO_CFLAGS=${CGO_CFLAGS} go install ./cmd/miru-insert/

plot:
	CGO_CFLAGS=${CGO_CFLAGS} go install ./cmd/miru-plot/

docker:
	docker build -t miru .
