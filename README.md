# e-factura-go #

Package efactura provides a client for using the ANAF e-factura API.

## Installation ##

e-factura-go is compatible with modern Go releases in module mode, with Go installed:

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
API via the Client object and for generating an UBL Invoice XML.

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
    OAuth2ConfigCredentials(anafAppClientID, anafApplientSecret),
    OAuth2ConfigRedirectURL(anafAppRedirectURL),
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
msg := efactura.MessageRASP{
    UploadIndex: 5008787839,
    Message: "test",
}
uploadRes, err := client.UploadRASPMessage(ctx, msg, "123456789")
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

## Tasks ##

- [ ] Support full OAuth2 authentication flow for the client, not just passing
  the initial token. This however will be tricky to implement properly since
  the OAuth2 app registered in the ANAF developer profile must have a fixed
  list of HTTPS redirect URLs and the redirect URL used for creating the OAuth2
  config must exactly matche one of the URLs.
- [ ] Add tests for all REST API calls and more tests for validating generated
  XML (maybe checking with the tools provided by mfinante).
- [ ] Implement CreditNote.
- [ ] Support parsing the ZIP file from the DownloadInvoice.
- [ ] Extend the InvoiceBuilder to add all Invoice fields
- [ ] Implement all business terms.
- [ ] Check and test API limits.
- [ ] Godoc and more code examples.
- [ ] Test coverage

## Contributing ##

TODO

## License ##

This library is distributed under the Apache License version 2.0 found in the
[LICENSE](./LICENSE) file.
