Not maintained project
Active development has moved to https://git.webix.io/Servers/dsfs

Dead Simple File Storage
----

File storage for temporary files.

### Usage

Create config.yml file ( check sample.config.yml file )

- folder - folder where files will be stored ( current dir by default )
- secret - cookies secret, change to some random junk
- google - google login API credentials


All config options can be redefined through ENV variables

```bash
#run and use current folder for file storage
dsfs

#run with custom config file
dsfs config.yaml

#run with custom settings
env CONFIGOR_FOLDER=/share dsfs
``` 

### License 

MIT
