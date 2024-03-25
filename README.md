# e-factura-go [![Tests](https://github.com/printesoi/e-factura-go/actions/workflows/tests.yml/badge.svg)](https://github.com/printesoi/e-factura-go/actions/workflows/test.yml) [![Coverage Status](https://coveralls.io/repos/github/printesoi/e-factura-go/badge.svg)](https://coveralls.io/github/printesoi/e-factura-go) [![Go Report Card](https://goreportcard.com/badge/github.com/printesoi/e-factura-go)](https://goreportcard.com/report/github.com/printesoi/e-factura-go)


Package efactura provides a client for using the ANAF e-factura API.

## NOTICE ##

**!!! This project is still in alpha stage, use it at you own risk !!!**

## Installation ##

e-factura-go requires Go version >= 1.21. With Go installed:

```bash
go get github.com/printesoi/e-factura-go
```

will resolve and add the package to the current development module, along with its dependencies.

Alternatively the same can be achieved if you use import in a package:

```go
import "github.com/printesoi/e-factura-go"
```

and run `go get` without parameters. The exported package name is `efactura`,
so you **don't** need to alias the import like:

```go
import efactura "github.com/printesoi/e-factura-go"
```

Finally, to use the top-of-trunk version of this repo, use the following command:

```bash
go get github.com/printesoi/e-factura-go@main
```

## Usage ##

This package can be use both for interacting with (calling) the ANAF e-factura
API via the Client object and for generating an Invoice UBL XML.

```go
import "github.com/printesoi/e-factura-go"
```

Construct a new client:

```go
client, err := efactura.NewClient(
    context.Background(),
    efactura.ClientOAuth2Config(oauth2Cfg),
    efactura.ClientOAuth2InitialToken(initialToken),
    efactura.ClientProductionEnvironment(false),
)
if err != nil {
    // Handle error
}
```

Construct the required OAuth2 config needed for the Client:

```go
oauth2Cfg, err := efactura.MakeOAuth2Config(
    efactura.OAuth2ConfigCredentials(anafAppClientID, anafApplientSecret),
    efactura.OAuth2ConfigRedirectURL(anafAppRedirectURL),
)
if err != nil {
    // Handle error
}
```

Parse the initial token from JSON:

```go
initialToken, err := efactura.TokenFromJSON([]byte(tokenJSON))
if err != nil {
    // Handle error
}
```

Getting a token from an authorization code (the parameter `code` sent via GET
to the redirect URL):

```go
// Assuming the oauth2Cfg is built as above
initialToken, err := oauth2Cfg.Exchange(ctx, authorizationCode)
if err != nil {
    // Handle error
}
```

### Upload invoice ###

```go
var invoice Invoice
// Build invoice (manually, or with the InvoiceBuilder)

uploadRes, err := client.UploadInvoice(ctx, invoice, "123456789")
if err != nil {
    // Handle error
}
if uploadRes.IsOk() {
    fmt.Printf("Upload index: %d\n", uploadRes.GetUploadIndex())
} else {
    // The upload was not successful, check uploadRes.Errors
}
```

For self-billed invoices, and/or if the buyer in not a Romanian entity, you can use
the `UploadOptionSelfBilled()`, `UploadOptionForeign()` upload options:

```go
uploadRes, err := client.UploadInvoice(ctx, invoice, "123456789",
        efactura.UploadOptionSelfBilled(), efactura.UploadOptionForeign())
```

If you have already the raw XML to upload (maybe you generated it by other means),
you can use the UploadXML method.

To upload an Invoice XML:

```go
uploadRes, err := client.UploadXML(ctx, xml, UploadStandardUBL, "123456789")
```

### Upload message ###

```go
msg := efactura.RaspMessage{
    UploadIndex: 5008787839,
    Message: "test",
}
uploadRes, err := client.UploadRaspMessage(ctx, msg, "123456789")
if err != nil {
    // Handle error
}
```

### Get message state ###

```go
resp, err := client.GetMessageState(ctx, uploadIndex)
if err != nil {
    // Handle error
}
switch {
case resp.IsOk():
    // Uploaded invoice was processed
    fmt.Printf("Download ID: %d\n", resp.GetDownloadID())
case resp.IsNok():
    // Processing failed
case resp.IsProcessing():
    // The message/invoice is still processing
case resp.IsInvalidXML():
    // The uploaded XML is invalid
```

### Get messages list ###

```go
numDays := 7
resp, err := client.GetMessagesList(ctx, "123456789", numDays, MessageFilterAll)
if err != nil {
    // Handle error
}
if resp.IsOk() {
    for _, message := range resp.Messages {
        switch {
        case message.IsError():
            // The message is an error for an upload
        case message.IsSentInvoice():
            // The message is a sent invoice
        case message.IsReceivedInvoice():
            // The message is a received invoice
        case message.IsBuyerMessage():
            // The message is a message from the buyer
        }
    }
}
```

### Download invoice ###

```go
downloadID := 3013004158
resp, err := client.DownloadInvoice(ctx, downloadID)
if err != nil {
    // Handle error
}
if resp.IsOk() {
    // The contents of the ZIP file is found in the resp.Zip byte slice.
}
```

### Validate invoice ###

```go
var invoice Invoice
// Build invoice (manually, or with the InvoiceBuilder)

validateRes, err := client.ValidateInvoice(ctx, invoice)
if err != nil {
    // Handle error
}
if validateRes.IsOk() {
    // Validation successful
}
```

### Errors ###

This library tries its best to overcome the not so clever API implementation
and to detect limits exceeded errors. To check if the error is cause by a
limit:

```go
resp, err := client.GetMessageState(ctx, uploadIndex)
if err != nil {
    var limitsErr *efactura.LimitExceededError
    var responseErr *efactura.ErrorResponse
    if errors.As(err, &limitsErr) {
        // The limits were exceeded. limitsErr.ErrorResponse contains more
        // information about the HTTP response, and the limitsErr.Limit field
        // contains the limit for the day.
    } else if errors.As(err, &responseErr) {
        // ErrorResponse means we got the HTTP response but we failed to parse
        // it or some other error like invalid response content type.
    }
}
```

## Generating an Invoice ##

TODO: See TestInvoiceBuilder() from builders_test.go for an example of using
InvoiceBuilder for creating an Invoice.

### Getting the raw XML of the invoice ##

In case you need to get the XML encoding of the invoice (eg. you need to store
it somewhere before the upload):

```go
var invoice Invoice
// Build invoice (manually, or with the InvoiceBuilder)

xmlData, err := invoice.XML()
if err != nil {
    // Handle error
}
```

To get the XML with indentation:

```go
xmlData, err := invoice.XMLIndent("", " ")
if err != nil {
    // Handle error
}
```


**NOTE** Don't use the standard `encoding/xml` package for generating the XML
encoding, since it does not produce Canonical XML [XML-C14N]!

### Unmarshal XML to invoice ##

```go
var invoice efactura.Invoice
if err := efactura.UnmarshalInvoice(data, &invoice); err != nil {
    // Handle error
}
```

**NOTE** Only use efactura.UnmarshalInvoice, because `encoding/xml` package
cannot unmarshal a struct like efactura.Invoice due to namespace prefixes!

## Tasks ##

- [ ] Implement all business terms.
- [ ] Support parsing the ZIP file from the DownloadInvoice.
- [ ] Extend the InvoiceBuilder to add all Invoice fields
- [ ] Implement CreditNote.
- [ ] Add tests for all REST API calls and more tests for validating generated
  XML (maybe checking with the tools provided by mfinante).
- [ ] Godoc and more code examples.
- [ ] Test coverage
- [ ] Support full OAuth2 authentication flow for the client, not just passing
  the initial token. This however will be tricky to implement properly since
  the OAuth2 app registered in the ANAF developer profile must have a fixed
  list of HTTPS redirect URLs and the redirect URL used for creating the OAuth2
  config must exactly matche one of the URLs.

## Contributing ##

Pull requests are more than welcome :)

## License ##

This library is distributed under the Apache License version 2.0 found in the
[LICENSE](./LICENSE) file.

## Commercial support ##

If you need help integrating this library in your software or you need
consulting services regarding e-factura APIs contact me (contact email in my
[Github profile](https://github.com/printesoi)).
