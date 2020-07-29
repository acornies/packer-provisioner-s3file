package s3file

import (
	"context"

	"github.com/hashicorp/hcl/v2/hcldec"
	"github.com/hashicorp/packer/common"
	"github.com/hashicorp/packer/packer"
	"github.com/hashicorp/packer/template/interpolate"
)

type Config struct {
	common.PackerConfig `mapstructure:",squash"`

	S3Endpoint string
	AccessKey  string
	SecretKey  string
	Bucket     string
	Key        string
	Region     string

	Destination string

	ctx interpolate.Context
}

type Provisioner struct {
	config Config
}

func (p *Provisioner) ConfigSpec() hcldec.ObjectSpec { return p.config.FlatMapstructure().HCL2Spec() }

// Prepare is called with a set of configurations to setup the
// internal state of the provisioner. The multiple configurations
// should be merged in some sane way.
func (p *Provisioner) Prepare(raws ...interface{}) error {
	return nil
}

// Provision is called to actually provision the machine. A context is
// given for cancellation, a UI is given to communicate with the user, and
// a communicator is given that is guaranteed to be connected to some
// machine so that provisioning can be done.
func (p *Provisioner) Provision(ctx context.Context, ui packer.Ui, comm packer.Communicator, generatedData map[string]interface{}) error {
	return nil
}
