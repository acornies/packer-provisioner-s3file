package main

import (
	"github.com/acornies/packer-provisioner-getter/s3file"
	"github.com/hashicorp/packer/packer/plugin"
)

func main() {
	server, err := plugin.Server()
	if err != nil {
		panic(err)
	}
	server.RegisterProvisioner(new(s3file.Provisioner))
	server.Serve()
}
