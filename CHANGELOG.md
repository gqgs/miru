# Changelog

All notable changes to this project will be documented in this file. See [standard-version](https://github.com/conventional-changelog/standard-version) for commit guidelines.

## [1.2.0](https://github.com/gqgs/miru/compare/v1.1.0...v1.2.0) (2021-09-11)


### Features

* **image:** AVX/AVX2 instructions for hist cmp ([65fe0b8](https://github.com/gqgs/miru/commit/65fe0b8aeba0f26b05525d155947aab26e769312))


### Bug Fixes

* **tree:** remove storage error check from hot path ([b96a402](https://github.com/gqgs/miru/commit/b96a402bbab7ab3ce6e7f23f1fe4b6b054568466))

## [1.1.0](https://github.com/gqgs/miru/compare/v1.0.0...v1.1.0) (2021-06-10)


### Features

* **all:** accept positional arguments ([3af26e4](https://github.com/gqgs/miru/commit/3af26e41238c9347c5bd854abf74145e9fc3f37f))
* **search:** allow exporting results as JSON ([bda7806](https://github.com/gqgs/miru/commit/bda7806110ae6beb0ff4e1fef7b6febcf612bc8d))
* **storage:** keep recent nodes in LRU cache ([aa6eb52](https://github.com/gqgs/miru/commit/aa6eb52fa74a4d64706b2c985d91e7a0a977ac6c))


### Bug Fixes

* **lint:** gocritic checks ([ed26673](https://github.com/gqgs/miru/commit/ed266737b2ec197d7a3a8f2fef3dec0eea59bc73))
* **perf:** decrease allocs in search ([ff7bfde](https://github.com/gqgs/miru/commit/ff7bfdebbdf7c4fdf1aa0020ed58e87cae5966af))
