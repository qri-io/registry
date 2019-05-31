default: test

define GOPACKAGES 
github.com/libp2p/go-libp2p-crypto \
github.com/jbenet/go-multihash \
github.com/sirupsen/logrus \
github.com/qri-io/apiutil \
github.com/qri-io/dataset
endef

define GX_DEP_PACKAGES 
github.com/qri-io/dag
endef

require-gopath:
ifndef GOPATH
  $(error $GOPATH must be set. plz check: https://github.com/golang/go/wiki/SettingGOPATH)
endif

install-deps: require-gopath
	go get -v -u $(GOPACKAGES)

install-gx:
	go get -v github.com/whyrusleeping/gx github.com/whyrusleeping/gx-go

install-gx-deps:
	gx install

install-gx-dep-packages:
	go get -v $(GX_DEP_PACKAGES)

test:
	go test ./...