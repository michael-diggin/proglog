module github.com/michael-diggin/proglog

go 1.14

require (
	github.com/casbin/casbin v1.9.1
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/golang/protobuf v1.5.2
	github.com/grpc-ecosystem/go-grpc-middleware v1.3.0
	github.com/hashicorp/raft v1.1.1
	github.com/hashicorp/raft-boltdb v0.0.0-20210422161416-485fa74b0b01
	github.com/hashicorp/serf v0.9.5
	github.com/kr/pretty v0.1.0 // indirect
	github.com/soheilhy/cmux v0.1.5
	github.com/spf13/cobra v1.1.3
	github.com/spf13/viper v1.8.1
	github.com/stretchr/testify v1.7.0
	github.com/tysontate/gommap v0.0.0-20210506040252-ef38c88b18e1
	go.opencensus.io v0.23.0
	go.uber.org/zap v1.17.0
	google.golang.org/genproto v0.0.0-20210602131652-f16073e35f0c
	google.golang.org/grpc v1.38.0
	google.golang.org/protobuf v1.26.0
	gopkg.in/check.v1 v1.0.0-20180628173108-788fd7840127 // indirect
	launchpad.net/gocheck v0.0.0-20140225173054-000000000087 // indirect
)

replace github.com/hasicorp/raft-boltdb => github.com/travisjeffery/raft-boltdb v1.0.0
