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
