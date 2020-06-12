module github.com/Akachain/hstx-go-sdk/handler

go 1.13

require (
	github.com/Akachain/akc-go-sdk v1.1.1
	github.com/VictoriaMetrics/fastcache v1.5.7 // indirect
	github.com/hyperledger/fabric v1.4.4
	github.com/hyperledger/fabric-protos-go v0.0.0-20200124220212-e9cfc186ba7b // indirect
	github.com/mitchellh/mapstructure v1.1.2

	github.com/Akachain/hstx-go-sdk/model v0.0.0
)

replace github.com/satori/go.uuid v1.2.0 => github.com/satori/go.uuid v1.2.1-0.20181028125025-b2ce2384e17b

replace github.com/Akachain/hstx-go-sdk/models v0.0.0 => ../model
