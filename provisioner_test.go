package main

import (
	"bytes"
	"context"
	"testing"

	"github.com/hashicorp/packer/packer"
)

func testConfig() map[string]interface{} {
	return map[string]interface{}{
		"destination":   "something",
		"s3_access_key": "test",
		"s3_secret_key": "test",
	}
}

func TestProvisioner_Impl(t *testing.T) {
	var raw interface{}
	raw = &Provisioner{}
	if _, ok := raw.(packer.Provisioner); !ok {
		t.Fatalf("must be a provisioner")
	}
}

func TestProvisionerPrepare_EmptyConfig(t *testing.T) {
	var p Provisioner
	config := testConfig()

	delete(config, "destination")
	delete(config, "s3_access_key")
	delete(config, "s3_secret_key")

	err := p.Prepare(config)
	if err == nil {
		t.Fatal("should have error")
	}
}

func TestProvisionerPrepare_BadS3URL(t *testing.T) {
	var p Provisioner
	config := testConfig()

	config["url"] = "bad!@#$%^&*()_{address}"

	err := p.Prepare(config)
	if err == nil {
		t.Fatal("should have error")
	}
}

func TestProvisionerProvision_S3BadCredentials(t *testing.T) {
	var p Provisioner
	config := testConfig()

	config["url"] = "https://s3.amazonaws.com/hc-oss-test/go-getter/folder/main.tf"

	if err := p.Prepare(config); err != nil {
		t.Fatalf("err: %s", err)
	}

	b := bytes.NewBuffer(nil)
	ui := &packer.BasicUi{
		Writer: b,
	}
	comm := &packer.MockCommunicator{}

	err := p.Provision(context.Background(), ui, comm, make(map[string]interface{}))

	if err == nil {
		t.Fatalf("Expected bad credentials failure %s", err)
	}

}
