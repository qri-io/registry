default: test

require-gopath:
ifndef GOPATH
  $(error $GOPATH must be set. plz check: https://github.com/golang/go/wiki/SettingGOPATH)
endif

install-deps: require-gopath
	go get -v github.com/libp2p/go-libp2p-crypto github.com/jbenet/go-multihash github.com/sirupsen/logrus github.com/datatogether/api/apiutil
	