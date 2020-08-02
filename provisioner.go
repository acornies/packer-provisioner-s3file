package main

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	getter "github.com/hashicorp/go-getter"
	"github.com/hashicorp/hcl/v2/hcldec"
	"github.com/hashicorp/packer/common"
	"github.com/hashicorp/packer/helper/config"
	"github.com/hashicorp/packer/packer"
	"github.com/hashicorp/packer/template/interpolate"
)

type Config struct {
	common.PackerConfig `mapstructure:",squash"`

	URL         string `mapstructure:"url"`
	S3AccessKey string `mapstructure:"s3_access_key"`
	S3SecretKey string `mapstructure:"s3_secret_key"`
	Destination string `mapstructure:"destination"`

	ctx interpolate.Context
}

type Provisioner struct {
	config Config
}

func (p *Provisioner) ConfigSpec() hcldec.ObjectSpec { return p.config.FlatMapstructure().HCL2Spec() }

func (p *Provisioner) Prepare(raws ...interface{}) error {

	err := config.Decode(&p.config, &config.DecodeOpts{
		Interpolate:        true,
		InterpolateContext: &p.config.ctx,
		InterpolateFilter:  &interpolate.RenderFilter{},
	}, raws...)

	if err != nil {
		return err
	}

	var errs *packer.MultiError

	if len(os.Getenv("AWS_ACCESS_KEY")) == 0 && len(p.config.S3AccessKey) == 0 {
		errs = packer.MultiErrorAppend(errs, errors.New("AWS_ACCESS_KEY environment variable or inline s3_access_key is required"))
	}

	if len(os.Getenv("AWS_SECRET_KEY")) == 0 && len(p.config.S3SecretKey) == 0 {
		errs = packer.MultiErrorAppend(errs, errors.New("AWS_SECRET_KEY environment variable required inline s3_secret_key is required"))
	}

	_, err = getURL(p.config.URL)
	if err != nil {
		errs = packer.MultiErrorAppend(errs, errors.New("Bad url parameter"))
	}

	if len(p.config.Destination) == 0 {
		errs = packer.MultiErrorAppend(errs, errors.New("Destination for S3 download is required"))
	}

	if errs != nil && len(errs.Errors) > 0 {
		return errs
	}

	return nil
}

func (p *Provisioner) Provision(ctx context.Context, ui packer.Ui, comm packer.Communicator, generatedData map[string]interface{}) error {

	if generatedData == nil {
		generatedData = make(map[string]interface{})
	}

	p.config.ctx.Data = generatedData

	s3 := new(getter.S3Getter)
	u, _ := getURL(p.config.URL)

	if len(p.config.S3AccessKey) > 0 {
		u.Query().Add("access_key", p.config.S3AccessKey)
	}

	if len(p.config.S3SecretKey) > 0 {
		u.Query().Add("secret_key", p.config.S3SecretKey)
	}

	local := getTempLocation(p.config.Destination)

	// Download
	err := s3.GetFile(local, u)
	ui.Say(fmt.Sprintf("Downloading %s => %s", p.config.URL, p.config.Destination))
	if err != nil {
		return err
	}

	return p.ProvisionUpload(ui, comm)
}

func (p *Provisioner) ProvisionUpload(ui packer.Ui, comm packer.Communicator) error {

	local := getTempLocation(p.config.Destination)

	src, err := interpolate.Render(local, &p.config.ctx)
	if err != nil {
		return fmt.Errorf("Error interpolating source: %s", err)
	}

	dst, err := interpolate.Render(p.config.Destination, &p.config.ctx)
	if err != nil {
		return fmt.Errorf("Error interpolating destination: %s", err)
	}

	ui.Say(fmt.Sprintf("Uploading %s => %s", src, dst))

	info, err := os.Stat(src)
	if err != nil {
		return err
	}

	f, err := os.Open(src)
	if err != nil {
		return err
	}
	defer f.Close()

	fi, err := f.Stat()
	if err != nil {
		return err
	}

	pf := ui.TrackProgress(filepath.Base(src), 0, info.Size(), f)
	defer pf.Close()

	// Upload the file
	if err = comm.Upload(dst, pf, &fi); err != nil {
		ui.Error(fmt.Sprintf("Upload failed: %s", err))
		return err
	}

	// Remove local file
	err = os.Remove(local)
	if err != nil {
		ui.Error(fmt.Sprintf("Clean up failed: %s", err))
		return err
	}
	return nil
}

func getURL(s string) (*url.URL, error) {
	u, err := url.Parse(s)
	return u, err
}

func getTempLocation(s string) string {
	split := strings.Split(s, "/")
	last := split[len(split)-1]
	return fmt.Sprintf("./%s", last)
}
