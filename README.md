# go-mauth-client

## Introduction
This is a simple client for the Medidata MAuth Authentication Protocol.  It can be used to access Platform Services within the Medidata Clinical Cloud.

##The Command Line Tool
As an example a simple cli tool has been added.  It can be built using `go build` and installed using `go install`

### Usage
```sh
Usage of ./go-mauth-client:
  -app-uuid string
    	Specify the App UUID
  -config string
    	Specify the configuration file
  -data string
    	Specify the data
  -method string
    	Specify the method (GET, POST) (default "GET")
  -private-key string
    	Specify the private key file
```

### The configuration file
The configuration file is a simple JSON file with the following structure:
```JSON
{"app_uuid": "1990d36d-4444-4105-b42e-223355334499",
"private_key_file": "innovate_private_key.pem"}
```
As an alternative the content of the private key can be included using a `private_key_text` attribute.

### Example
```sh
go-mauth-client git:(develop) $ ./go-mauth-client -config innovate.json https://innovate.imedidata.com/api/v2/studies/55555555-5508-45c6-3333-1234512345.json
Created MAuth App with App UUID: 12345678-c109-11e1-84f6-5432112345
Status Code: 200
Response Body:
{"study":{"name":"Mediflex (DEV)","uuid":"55555555-5508-45c6-3333-1234512345","mcc_study_uuid":"","oid":"Mediflex (DEV)","live_date":"","title":"","summary":"","drug_device":"","compound_code":"","number":"","program":"","protocol":"Mediflex","indication":"","responsible_party":"","enrollment_target":"","investigator":"","full_description":"EHR Integration testing study. Ian Sparks, Geoff Low","therapeutic_area":"Unknown","phase":"III","close_date":"","is_production":"false","parent_uuid":"55555555-1111-2222-3333-4444444444","status":"active","client_division_uuid":"","study_environment_type":"Development"}}
```

## Developer Notes

### Making a Keypair

Make a keypair with:
```sh
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