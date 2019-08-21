module github.com/qri-io/registry

go 1.12

replace (
	github.com/go-critic/go-critic v0.0.0-20181204210945-c3db6069acc5 => github.com/go-critic/go-critic v0.0.0-20190422201921-c3db6069acc5
	github.com/go-critic/go-critic v0.0.0-20181204210945-ee9bf5809ead => github.com/go-critic/go-critic v0.0.0-20190210220443-ee9bf5809ead
	github.com/golangci/errcheck v0.0.0-20181003203344-ef45e06d44b6 => github.com/golangci/errcheck v0.0.0-20181223084120-ef45e06d44b6
	github.com/golangci/go-tools v0.0.0-20180109140146-af6baa5dc196 => github.com/golangci/go-tools v0.0.0-20190318060251-af6baa5dc196
	github.com/golangci/gofmt v0.0.0-20181105071733-0b8337e80d98 => github.com/golangci/gofmt v0.0.0-20181222123516-0b8337e80d98
	github.com/golangci/gosec v0.0.0-20180901114220-66fb7fc33547 => github.com/golangci/gosec v0.0.0-20190211064107-66fb7fc33547
	github.com/golangci/lint-1 v0.0.0-20180610141402-ee948d087217 => github.com/golangci/lint-1 v0.0.0-20190420132249-ee948d087217
	mvdan.cc/unparam v0.0.0-20190124213536-fbb59629db34 => mvdan.cc/unparam v0.0.0-20190209190245-fbb59629db34
)

require (
	github.com/ipfs/go-ipld-format v0.0.2
	github.com/ipfs/interface-go-ipfs-core v0.0.8
	github.com/libp2p/go-libp2p-crypto v0.0.2
	github.com/mr-tron/base58 v1.1.2
	github.com/multiformats/go-multihash v0.0.5
	github.com/qri-io/apiutil v0.1.0
	github.com/qri-io/dag v0.1.1-0.20190605213518-cb095ea6b6d9
	github.com/qri-io/dataset v0.1.3-0.20190617151150-bd20b1913ba5
	github.com/sirupsen/logrus v1.4.2
	github.com/ugorji/go v1.1.5-pre // indirect
	golang.org/x/crypto v0.0.0-20190701094942-4def268fd1a4 // indirect
	golang.org/x/net v0.0.0-20190522155817-f3200d17e092 // indirect
)
