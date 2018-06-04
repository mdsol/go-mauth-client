# go-mauth-client

## Introduction
This is a simple client for the Medidata MAuth Authentication Protocol.  It can be used to access Platform Services within the Medidata Clinical Cloud.

### What is MAuth?
The MAuth protocol provides a fault-tolerant, service-to-service authentication scheme for Medidata and third-party applications that use web services to communicate. The Authentication Service and integrity algorithm is based on digital signatures encrypted and decrypted with a private/public key pair.

The Authentication Service has two responsibilities. It provides message integrity and provenance validation by verifying a message sender's signature; its other task is to manage public keys. Each public key is associated with an application and is used to authenticate message signatures. The private key corresponding to the public key in the Authentication Service is stored by the application making a signed request; the request is encrypted with this private key. The Authentication Service has no knowledge of the application's private key, only its public key.


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
* [medidata_api_example.go](examples/medidata_apis/medidata_api_example.go) - sample code illustrating use for accessing the Medidata API 
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

