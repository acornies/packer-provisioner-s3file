# packer-provisioner-s3file

A simple HashiCorp [Packer](https://packer.io) provisioner, utilizing [github.com/hashicorp/go-getter](https://github.com/hashicorp/go-getter) to download files from Amazon S3 or S3 compatible endpoints.

## Usage

Using `AWS_ACCESS_KEY` and `AWS_SECRET_KEY` environment variables:

```json
{
  "builders": [],
  "provisioners": [
    {
      "type": "s3file",
      "s3_url": "https://s3-region.amazonaws.com/bucket/path/key",
      "destination": "./example"
    }
  ]
}
```

Using inline S3 credentials and user-defined input variables:

```json
{
  "variables": {
    "aws_access_key": "{{ env `AWS_ACCESS_KEY` }}",
    "aws_secret_key": "{{ env `AWS_SECRET_KEY` }}"
  },
  "builders": [],
  "provisioners": [
    {
      "type": "s3file",
      "s3_url": "https://s3-region.amazonaws.com/bucket/path/key",
      "s3_access_key": "{{user `aws_access_key`}}",
      "s3_secret_key": "{{user `aws_secret_key`}}",
      "destination": "./example"
    }
  ]
}
```
