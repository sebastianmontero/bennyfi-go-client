module github.com/sebastianmontero/bennyfi-go-client

go 1.16

require (
	github.com/sebastianmontero/eos-go v0.10.5-0.20250701034638-c1d2a67e01a3
	github.com/sebastianmontero/eos-go-toolbox v0.0.0-20250717142254-4d4e82a05320
	gotest.tools v2.2.0+incompatible
)

// replace github.com/sebastianmontero/eos-go-toolbox => ../eos-go-toolbox

replace github.com/sebastianmontero/eos-go => ../eos-go
