# go-mauth-client

## Installation

It can be installed using:

```bash
$ go install github.com/mdsol/go-mauth-client/cmd/go_mauth_client
``` 

### Usage
```bash
$ go_mauth_client -help
Usage of go_mauth_client:
  -app-uuid string
        Specify the App UUID
  -config string
        Specify the configuration file
  -data string
        Specify the data
  -headers
        Print the Response Headers
  -method string
        Specify the method (GET, POST, PUT, DELETE) (default "GET")
  -pretty
        Prettify the Output
  -private-key string
        Specify the private key file
  -verbose
        Print out more information
  -version
        Print out the version
```

### The configuration file
The configuration file is a simple JSON file with the following structure:
```json
{"app_uuid": "1990d36d-4444-4105-b42e-223355334499",
"private_key_file": "innovate_private_key.pem"}
```
As an alternative the content of the private key can be included using a `private_key_text` attribute.

### Example Usages
* Non-prettified output
    ```bash
    $ ./go-mauth-client -config innovate.json https://innovate.imedidata.com/api/v2/studies/55555555-5508-45c6-3333-1234512345.json
    {"study":{"name":"Mediflex (DEV)","uuid":"55555555-5508-45c6-3333-1234512345", ... ,"study_environment_type":"Development"}}
    ```
* Prettified output (applies to JSON and XML based on `Content-Type` Header)
    ```bash
    $ go-mauth-client -pretty -config credentials_1.json https://innovate.imedidata.com/api/v2/users/c123a678-79e5-11e1-7789-123138140309/studies
    {
        "studies": [
            {
                "name": "test_for_innovate (DEV)",
                "uuid": "3241245e-b2ae-1123-98f7-145bf03bbbee",
                "href": "https://innovate.imedidata.com/api/v2/studies/3241245e-b2ae-1123-98f7-145bf03bbbee",
                "parent_uuid": "9718efe1-4311-11e0-8747-1231390e6521",
                "created_at": "2013/02/13 20:01:04 +0000",
                "updated_at": "2016/05/19 21:12:30 +0000"
            },
            ....
    }
    ```
