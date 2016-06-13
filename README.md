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

