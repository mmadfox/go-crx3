# go-crx3 [![Coverage Status](https://coveralls.io/repos/github/mediabuyerbot/go-crx3/badge.svg?branch=master&v=2)](https://coveralls.io/github/mediabuyerbot/go-crx3?branch=master)
Provides a sets of tools packing, unpacking, zip, unzip, download, gen id, etc...

## Table of contents
+ [Installation](#installation)
+ [Commands](#commands)
+ [Examples](#examples)
  - [Encode to base64 string](#base64)
  - [Pack a zip file or unzipped directory into a crx extension](#pack)
  - [Unpack chrome extension into current directory](#unpack)
  - [Download a chrome extension from the web store](#download)
  - [Add unpacked extension to archive](#zip)
  - [Unzip an extension to the directory](#unzip)
  - [Keygen](#keygen)
  - [Generate extension id](#gen-id)
  - [IsDir, IsZip, IsCRX3 helpers](#isdir-iszip-iscrx3)
  - [Load or save private key](#newprivatekey-loadprivatekey-saveprivatekey)
+ [License](#license)


### Installation
```ssh
go get -u github.com/mediabuyerbot/go-crx3/crx3
```

### Commands
```shell script
make proto 
make covertest
``` 

### Examples
#### Pack
##### Pack a zip file or unzipped directory into a crx extension 
```go
import crx3 "github.com/mediabuyerbot/go-crx3"

if err := crx3.Extension("/path/to/file.zip").Pack(nil); err != nil {
    panic(err)
}
```

```go
import crx3 "github.com/mediabuyerbot/go-crx3"

pk, err := crx3.LoadPrivateKey("/path/to/key.pem")
if err != nil { 
    panic(err) 
}
if err := crx3.Extension("/path/to/file.zip").Pack(pk); err != nil {
    panic(err)
}
```

```go
import crx3 "github.com/mediabuyerbot/go-crx3"

pk, err := crx3.LoadPrivateKey("/path/to/key.pem")
if err != nil { 
    panic(err) 
}
if err := crx3.Extension("/path/to/file.zip").PackTo("/path/to/ext.crx", pk); err != nil {
    panic(err)
}
```
```shell script
$ crx3 pack /path/to/file.zip 
$ crx3 pack /path/to/file.zip -p /path/to/key.pem 
$ crx3 pack /path/to/file.zip -p /path/to/key.pem -o /path/to/ext.crx 
```

#### Unpack
##### Unpack chrome extension into current directory
```go
import crx3 "github.com/mediabuyerbot/go-crx3"

if err := crx3.Extension("/path/to/ext.crx").Unpack(); err != nil {
   panic(err)
}
```
```shell script
$ crx3 unpack /path/to/ext.crx 
```

#### Base64
##### Encode an extension file to a base64 string
```go
import crx3 "github.com/mediabuyerbot/go-crx3"
import "fmt"

b, err := crx3.Extension("/path/to/ext.crx").Base64()
if err != nil {
   panic(err)
}
fmt.Println(string(b))
```
```shell script
$ crx3 base64 /path/to/ext.crx [-o /path/to/file] 
```

#### Download 
##### Download a chrome extension from the web store
```go
import crx3 "github.com/mediabuyerbot/go-crx3"

extensionID := "blipmdconlkpinefehnmjammfjpmpbjk"
filepath := "/path/to/ext.crx"
if err := crx3.DownloadFromWebStore(extensionID,filepath); err != nil {
    panic(err)
}
```
```shell script
$ crx3 download blipmdconlkpinefehnmjammfjpmpbjk [-o /custom/path]
$ crx3 download https://chrome.google.com/webstore/detail/lighthouse/blipmdconlkpinefehnmjammfjpmpbjk
```

#### Zip
##### Zip add an unpacked extension to the archive
```go
import crx3 "github.com/mediabuyerbot/go-crx3"

if err := crx3.Extension("/path/to/unpacked").Zip(); err != nil {
    panic(err)
}
```
```shell script
$ crx3 zip /path/to/unpacked [-o /custom/path] 
```

#### Unzip
##### Unzip an extension to the current directory
```go
import crx3 "github.com/mediabuyerbot/go-crx3"

if err := crx3.Extension("/path/to/ext.zip").Unzip(); err != nil {
    panic(err)
}
```
```shell script
$ crx3 unzip /path/to/ext.zip [-o /custom/path] 
``` 

#### Gen ID
##### Generate extension id (like dgmchnekcpklnjppdmmjlgpmpohmpmgp)
```go
import crx3 "github.com/mediabuyerbot/go-crx3"

id, err := crx3.Extension("/path/to/ext.crx").ID()
if err != nil {
    panic(err)
}
```
```shell script
$ crx3 id /path/to/ext.crx 
```

#### IsDir, IsZip, IsCRX3
```go
import crx3 "github.com/mediabuyerbot/go-crx3"

crx3.Extension("/path/to/ext.zip").IsZip()
crx3.Extension("/path/to/ext").IsDir()
crx3.Extension("/path/to/ext.crx").IsCRX3()
```

#### NewPrivateKey, LoadPrivateKey, SavePrivateKey 
```go
import crx3 "github.com/mediabuyerbot/go-crx3"

pk, err := crx3.NewPrivateKey()
if err != nil {
    panic(err)
}
if err := crx3.SavePrivateKey("/path/to/key.pem", pk); err != nil {
    panic(err)
}
pk, err = crx3.LoadPrivateKey("/path/to/key.pem")
```

#### Keygen
```shell script
$ crx3 keygen /path/to/key.pem 
``` 

## License
go-crx3 is released under the Apache 2.0 license. See [LICENSE.txt](https://github.com/mediabuyerbot/go-crx3/blob/master/LICENSE)
