build:
	GOOS=darwin GOARCH=arm64 go build -gcflags="all=-N -l" -o ~/.terraform.d/plugins/registry.terraform.io/matthisholleville/pritunl/0.1.6/darwin_arm64/terraform-provider-pritunl_v0.1.6 main.go

build-current-dir:
	GOOS=darwin GOARCH=arm64 go build -gcflags="all=-N -l" main.go

docs:
	tfplugindocs generate --provider-name pritunl
