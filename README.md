# go-mauth-client

## Introduction
This is a simple client for the Medidata MAuth Authentication Protocol.  It can be used to access Platform Services within the Medidata Clinical Cloud.


## The Command Line Tool
A simple cli tool has been added in the [cmd/go_mauth_client](cmd/go_mauth_client) folder.

It can be installed using:

```bash
$ go install github.com/mdsol/go-mauth-client/cmd/go_mauth_client
``` 

See [README.md](cmd/go_mauth_client/README.md) for examples of usage.

## Developer Notes

### Examples

Two examples have been provided:
* [api_gateway_example.go](examples/api_gateway/api_gateway_example.go) - sample code illustrating use for accessing the Medidata API Gateway
* [imedidata_client.go](examples/imedidata/imedidata_client.go) - sample code illustrating an iMedidata call

### Making a Keypair

Make a keypair with:

```bash
$ mkdir keypair_dir
$ cd keypair_dir
$ openssl genrsa -out yourname_mauth.priv.key 2048
$ chmod 0600 yourname_mauth.priv.key
$ openssl rsa -in yourname_mauth.priv.key -pubout -out yourname_mauth.pub.key
```

Provide PUBLIC key to DevOps via Zendesk ticket with settings:

   * Form: OPS - Service Request
   * Assignee: OPS-Devops Cloud Team
   * Service Catalog: Application support

Keep the Private key to yourself, don't check into GitHub, don't tell your mother.

### Generating a Test String

You may need to generate a test string for use, this snippet will help

* Install ruby
* Use `irb` to do the following:
    ```ruby
    irb(main):001:0> require 'openssl'
    irb(main):002:0> key_file = File.read('private_key.pem')
    irb(main):003:0> key = OpenSSL::PKey::RSA.new(key_file)
    irb(main):004:0> require 'digest'
    irb(main):005:0> hashed = Digest::SHA512.hexdigest("Hello world")
    irb(main):006:0> require 'base64'
    irb(main):007:0> Base64.encode64(key.private_encrypt(hashed)).delete("\n")
    => "IUjQhtH4C9lbCRTyca+...Tvlg=="
    ```
* Copy the string to your test case

