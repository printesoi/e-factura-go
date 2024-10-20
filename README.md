# e-factura-go [![Go Reference](https://pkg.go.dev/badge/github.com/printesoi/e-factura-go@main.svg)](https://pkg.go.dev/github.com/printesoi/e-factura-go@main) [![Tests](https://github.com/printesoi/e-factura-go/actions/workflows/tests.yml/badge.svg)](https://github.com/printesoi/e-factura-go/actions/workflows/test.yml) [![Coverage Status](https://coveralls.io/repos/github/printesoi/e-factura-go/badge.svg)](https://coveralls.io/github/printesoi/e-factura-go) [![Go Report Card](https://goreportcard.com/badge/github.com/printesoi/e-factura-go)](https://goreportcard.com/report/github.com/printesoi/e-factura-go)


Package e-factura-go provides a client for using the RO e-factura and RO e-transport APIs.

## NOTICE ##

**!!! This project is still in alpha stage, use it at you own risk !!!**

## Installation ##

e-factura-go requires Go version >= 1.21. With Go installed:

```bash
go get github.com/printesoi/e-factura-go
```

will resolve and add the package to the current development module, along with its dependencies.

## oauth2 ##

Construct the required OAuth2 config needed for an e-factura or e-transport Client:

```go
import (
    efactura_oauth2 "github.com/printesoi/e-factura-go/pkg/oauth2"
)

oauth2Cfg, err := efactura_oauth2.MakeConfig(
    efactura_oauth2.ConfigCredentials(anafAppClientID, anafAppClientSecret),
    efactura_oauth2.ConfigRedirectURL(anafAppRedirectURL),
)
if err != nil {
    // Handle error
}
```

Generate an authorization link for certificate authorization:

```go
authorizeURL := oauth2Cfg.AuthCodeURL(state)
// Redirect the user to authorizeURL
```

Getting a token from an authorization code (the parameter `code` sent via GET
to the redirect URL):

```go
// Assuming the oauth2Cfg is built as above
token, err := oauth2Cfg.Exchange(ctx, authorizationCode)
if err != nil {
    // Handle error
}
```

If you specified a non-empty state when building the authorization URL, you
will also receive the `state` parameter with `code`.

Parse the initial token from JSON:

```go
token, err := efactura_oauth2.TokenFromJSON([]byte(tokenJSON))
if err != nil {
    // Handle error
}
```

## e-factura ##

This package can be use both for interacting with (calling) the
[RO e-factura API](https://mfinante.gov.ro/ro/web/efactura/informatii-tehnice)
via the Client object and for generating an Invoice-2 UBL XML.

Construct a new simple client for production environment:

```go
import (
    "github.com/printesoi/e-factura-go/pkg/efactura"
)

ctx := context.TODO()
client, err := efactura.NewProductionClient(ctx, efactura_oauth2.TokenSource(ctx, token))
if err != nil {
    // Handle error
}
```

Construct a new simple client for sandbox (test) environment:

```go
ctx := context.TODO()
client, err := efactura.NewSandboxClient(ctx, efactura_oauth2.TokenSource(ctx, token))
if err != nil {
    // Handle error
}
```

If you want to store the token in a store/db and update it everytime it
refreshes use `efactura_oauth2.TokenSourceWithChangedHandler`:

```go
ctx := context.TODO()
onTokenChanged := func(ctx context.Context, token *xoauth.Token) error {
    // Token changed, maybe update/store it in a database/persistent storage.
    return nil
}
client, err := efactura.NewSandboxClient(ctx,
    efactura_oauth2.TokenSourceWithChangedHandler(ctx, token, onTokenChanged))
if err != nil {
    // Handle error
}
```

### Time and dates in Romanian time zone ###

E-factura APIs expect dates to be in Romanian timezone and will return dates
and times in Romanian timezone. This library tries to load the
`Europe/Bucharest` timezone location on init so that creating and parsing dates
will work as expected. **The user of this library is responsible to ensure the
`Europe/Bucharest` location is available**. If you are not sure that the target
system will have system timezone database, you can use in you `main` package:

```go
import _ "time/tzdata"
```

to load the Go embedded copy of the timezone database.

### Upload invoice ###

```go
var invoice efactura.Invoice
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
import (
    efactura_errors "github.com/printesoi/e-factura-go/pkg/errors"
)

resp, err := client.GetMessageState(ctx, uploadIndex)
if err != nil {
    var limitsErr *efactura_errors.LimitExceededError
    var responseErr *efactura_errors.ErrorResponse
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

## RO e-Transport ##

The `etransport` package can be used for interacting with (calling) the
[RO e-Transport v2](https://mfinante.gov.ro/ro/web/etransport/informatii-tehnice) API
via the Client object or to build a declaration v2 XML using the PostingDeclarationV2 objects.

Construct a new simple client for production environment:

```go
import (
    "github.com/printesoi/e-factura-go/pkg/etransport"
)

ctx := context.TODO()
client, err := etransport.NewProductionClient(ctx, efactura_oauth2.TokenSource(ctx, token))
if err != nil {
    // Handle error
}
```

Construct a new simple client for sandbox (test) environment:

```go
ctx := context.TODO()
client, err := etransport.NewSandboxClient(ctx, efactura_oauth2.TokenSource(ctx, token))
if err != nil {
    // Handle error
}
```

### Upload declaration ###

```go
var declaration etransport.PostingDeclarationV2
// Build posting declaration

uploadRes, err := client.UploadPostingDeclarationV2(ctx, declaration, "123456789")
if err != nil {
    // Handle error
}
if uploadRes.IsOk() {
    fmt.Printf("Upload index: %d\n", uploadRes.GetUploadIndex())
} else {
    // The upload was not successful, check uploadRes.Errors
    fmt.Printf("Upload failed: %s\n", uploadRes.GetFirstErrorMessage())
}
```

### Get message state ###

Check the message state for an upload index resulted from an upload:

```go
resp, err := client.GetMessageState(ctx, uploadIndex)
if err != nil {
    // Handle error
}
switch {
case resp.IsOk():
    // Uploaded declaration was processed, we can now use the UIT.
case resp.IsNok():
    // Processing failed
case resp.IsProcessing():
    // The message/declaration is still processing
case resp.IsInvalidXML():
    // The uploaded XML is invalid
```

### Get messages list ###

```go
numDays := 7 // Between 1 and 60
resp, err := client.GetMessagesList(ctx, "123456789", numDays)
if err != nil {
    // Handle error
}
if resp.IsOk() {
    for _, message := range resp.Messages {
        // Process message
    }
} else {
    // Handle error
    fmt.Printf("GetMessagesList failed: %s\n", resp.GetFirstErrorMessage())
}
```

## Contributing ##

Pull requests are more than welcome :)

## License ##

This library is distributed under the Apache License version 2.0 found in the
[LICENSE](./LICENSE) file.

## Commercial support ##

If you need help integrating this library in your software or you need
consulting services regarding e-factura APIs contact me (contact email in my
[Github profile](https://github.com/printesoi)).
