build:
	GOOS=darwin GOARCH=arm64 go build -gcflags="all=-N -l" -o ~/.terraform.d/plugins/registry.terraform.io/matthisholleville/pritunl/0.0.1/darwin_arm64/terraform-provider-pritunl_v0.0.1 main.go

docs:
	tfplugindocs generate --provider-name pritunl
