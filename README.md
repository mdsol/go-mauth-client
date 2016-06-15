# go-mauth-client

Making a Keypair
=================

Make a keypair with:

    $ mkdir keypair_dir
    $ cd keypair_dir
    $ openssl genrsa -out yourname_mauth.priv.key 2048
    $ chmod 0600 yourname_mauth.priv.key
    $ openssl rsa -in yourname_mauth.priv.key -pubout -out yourname_mauth.pub.key

Provide PUBLIC key to DevOps via Zendesk ticket with settings:

   * Form: OPS - Service Request
   * Assignee: OPS-Devops Cloud Team
   * Service Catalog: Application support

Keep the Private key to yourself, don't check into GitHub, don't tell your mother.

Generating a Test String
========================

* Install ruby
* Use `irb` to do the following:
```ruby

irb(main):001:0> require 'openssl'
irb(main):002:0> key_file = File.read('private_key.pem')
irb(main):003:0> key = OpenSSL::PKey::RSA.new(key_file)
irb(main):004:0> require 'base64'
irb(main):004:0> Base64.encode64(key.private_encrypt("Hello world")).delete('\n')
=> "IUjQhtH4C9lbCRTyca+...Tvlg=="

```
* Copy the string to your test case