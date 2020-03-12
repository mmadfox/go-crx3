# go-crx3 [![Coverage Status](https://coveralls.io/repos/github/mediabuyerbot/go-crx3/badge.svg?branch=master)](https://coveralls.io/github/mediabuyerbot/go-crx3?branch=master)
Provides a set of tools packing, unpacking, zip, unzip, download, etc.

## Contents
+ [Installation](#installation)
+ [Commands](#commands)
+ [Examples](#examples)
  - [Pack](#pack)
  - [Unpack](#unpack)
  - [Encode to base64](#encode-to-base64)
  - [Download](#download)
  - [Zip](#zip)
  - [Unzip](#unzip)
  - [Keygen](#keygen)
  - [IsDir, IsZip, IsCRX3](#isdir-iszip-iscrx3)
  - [Private key](#newprivatekey-loadprivatekey-saveprivatekey)
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
```go
import crx3 "github.com/mediabuyerbot/go-crx3"

if err := crx3.Extension("/path/to/ext.crx").Unpack(); err != nil {
   panic(err)
}
```
```shell script
$ crx3 unpack /path/to/ext.crx 
```

#### Encode to base64
```go
import crx3 "github.com/mediabuyerbot/go-crx3"
import "fmt"

b, err := crx3.Extension("/path/to/ext.crx").ToBase64()
if err != nil {
   panic(err)
}
fmt.Println(string(b))
```
```shell script
$ crx3 encode /path/to/ext.crx [-o /path/to/file] 
```

#### Download 
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
```go
import crx3 "github.com/mediabuyerbot/go-crx3"

if err := crx3.Extension("/path/to/ext.zip").Unzip(); err != nil {
    panic(err)
}
```
```shell script
$ crx3 unzip /path/to/ext.zip [-o /custom/path] 
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
