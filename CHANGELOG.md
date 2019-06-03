# 0.1.0 (2019-06-03)

Registry defines primitives for keeping centralized repositories of qri types (peers, datasets, etc). It uses classical client/server patterns, arranging types into cannonical stores.

This is the first proper release of `registry`. In preparation for go 1.13, in which `go.mod` files and go modules are the primary way to handle go dependencies, we are going to do an official release of all our modules. This will be version v0.1.0 of `registry`. We'll be working on adding details & documentation in the near future.


### Bug Fixes

* fix breaks from new dependency versions ([8afb295](https://github.com/qri-io/registry/commit/8afb295))
* **circleci:** circle ci failing on a `go get`, changing import path ([f8b8dfb](https://github.com/qri-io/registry/commit/f8b8dfb))
* **client pin polling:** fix send on closed channel ([159e404](https://github.com/qri-io/registry/commit/159e404))
* **dataset:** support getting dataset info ([f8d9691](https://github.com/qri-io/registry/commit/f8d9691))
* **datasets:** fix dataset listing ([77b3b53](https://github.com/qri-io/registry/commit/77b3b53))
* **mockserver:** May add datasets to mock registry server ([47f3438](https://github.com/qri-io/registry/commit/47f3438))
* **mockserver:** restore search for mock-server ([1a4c7dd](https://github.com/qri-io/registry/commit/1a4c7dd))
* **pin status polling:** fix defer/goroutine deadlock ([200e1fc](https://github.com/qri-io/registry/commit/200e1fc))
* **pinset:** Ensure pins are stored in lexographical orders ([381f2f2](https://github.com/qri-io/registry/commit/381f2f2))
* **regclient:** should use `page` and `pageSize` when requesting from server ([ae7a78a](https://github.com/qri-io/registry/commit/ae7a78a))
* **root endpoint:** root should return 200 status for health checks ([5b39558](https://github.com/qri-io/registry/commit/5b39558))
* **search:** If json parameters don't give a limit, assume the default ([f9202da](https://github.com/qri-io/registry/commit/f9202da))
* moved dockerile to project root ([ecd8082](https://github.com/qri-io/registry/commit/ecd8082))


### Features

* **adminKey:** added admin post /profiles to re-hydrate ([5cf6781](https://github.com/qri-io/registry/commit/5cf6781))
* **client.GetDataset:** added client dataset getting, tests ([f44bd3e](https://github.com/qri-io/registry/commit/f44bd3e))
* **dataset:** add dataset tracking to registries ([5310939](https://github.com/qri-io/registry/commit/5310939))
* **deregister:** remove profile with signed deregister ([a520827](https://github.com/qri-io/registry/commit/a520827))
* **DsyncFetch:** add dsync dataset fetch method to registry client ([ca61bbc](https://github.com/qri-io/registry/commit/ca61bbc))
* **errs:** added ErrNoRegistry for no connection tracking ([8f15f7f](https://github.com/qri-io/registry/commit/8f15f7f))
* **expiration:** Added ttl to reputation response ([18948a1](https://github.com/qri-io/registry/commit/18948a1))
* **Indexer:** add/remove datasets form search index on publication ([f2f7813](https://github.com/qri-io/registry/commit/f2f7813))
* **interfaces:** abstract profiles & datasets into interfaces ([ce0d78c](https://github.com/qri-io/registry/commit/ce0d78c))
* **listDatasets:** `regclient.ListDatasets` requests a list of datasets from the registry ([e13873b](https://github.com/qri-io/registry/commit/e13873b))
* **mock server, route protector:** added server mock package, protector iface ([8fb0e3f](https://github.com/qri-io/registry/commit/8fb0e3f))
* **Pinset:** add support for pinset mgmt ([a738c7a](https://github.com/qri-io/registry/commit/a738c7a))
* **Pinset:** support async pinning with refactored interface ([de6054d](https://github.com/qri-io/registry/commit/de6054d))
* **regclient:** add initial dataset support for regclient ([815f696](https://github.com/qri-io/registry/commit/815f696))
* **regclient:** added client package for working with registries ([8159463](https://github.com/qri-io/registry/commit/8159463))
* **regclient.reputation:** GetReputation request ([a5a5dc5](https://github.com/qri-io/registry/commit/a5a5dc5))
* **registry struct, pinset:** unified registry struct, pinset iface ([4f33cd9](https://github.com/qri-io/registry/commit/4f33cd9))
* **regserver.reputation:** add reputation handler ([7b78c9d](https://github.com/qri-io/registry/commit/7b78c9d))
* **reputation:** most basic handler added for `/reputation` endpoint ([9b30037](https://github.com/qri-io/registry/commit/9b30037))
* **Reputation, Reputations:** add the concept of reputation to registry ([d239d7d](https://github.com/qri-io/registry/commit/d239d7d))
* added regclient search ([3e51490](https://github.com/qri-io/registry/commit/3e51490))
* added Searchable interface and NewSearchHandler ([e077c3f](https://github.com/qri-io/registry/commit/e077c3f))
* bring package into existence ([1c91c61](https://github.com/qri-io/registry/commit/1c91c61))



