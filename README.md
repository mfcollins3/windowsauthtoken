Windows Authentication Handler for Go
=====================================
Package `windowsauthtoken` implements Go HTTP middleware that will
extract the username of a Windows user when running a Go web application
in IIS with Windows authentication enabled. Package `windowsauthtoken`
works with the [HttpPlatformHandler](http://www.iis.net/downloads/microsoft/httpplatformhandler)
module. When Windows authentication is enabled on the IIS website or
application and HttpPlatformHandler is configured to pass the user's
token to the Go web application, HttpPlatformHandler will pass the
handle to the Go web application in the `X-IIS-WindowsAuthToken` HTTP
header. The middleware providwed by package `windowsauthtoken` will
obtain the full domain name of the Windows user and will make the
username available to the web application.

Package `windowsauthtoken` is only available on a Windows Server when
running the Go web application through IIS using the
[HttpPlatformHandler](http://www.iis.net/downloads/microsoft/httpplatformhandler)
module. When using this package on non-Windows platforms, it is
essentially a no-op passthrough hamdler, so it is safe to use and call
in web applications that can run on multiple platforms.

Usage
-----
First, get the source code for the package:

    $ go get github.com/mfcollins3/windowsauthtoken
    
Next, use the middleware in your web application. HttpPlatformHandler
will pass the Windows authentication token on every request, and it is
the responsibility of your web application to close the handle, so you
will want to apply the middleware globally to all requests. To do this,
you will typically wrap the `http.DefaultServeMux` handler with the
`WindowsAuthTokenHandler` middleware:

```go

package main

import (
    "net/http"
    
    "github.com/mfcollins3/windowsauthtoken"
)

func main() {
    // TODO: register HTTP handlers with net/http
    
    rootHandler := windowsauthtoken.WindowsAuthTokenHandler(
        http.DefaultServeMux,
        func(username string) error {
            // Store username somewhere for the request. For example,
            // you can use Gorilla Toolkit's context package to store
            // the username in the request.
            return nil
        })
    log.Fatal(http.ListenAndServe(":8080", rootHandler))    
}

```

To run your Go web application in IIS, do the following:

1. Install [HttpPlatformHandler](http://www.iis.net/downloads/microsoft/httpplatformhandler)
2. Create a website or application virtual directory in an existing
   website.
3. Create a new handler module mapping that maps all requests (*) to
   the `httpPlatformHandler` module.
4. Copy your Golang application to the physical directory that you
   pointed the website or application virtual directory to in step 2.
5. Create a `web.config` file in the physical directory:

```xml

<?xml version="1.0" encoding="UTF-8"?>
<configuration>
  <system.webServer>
    <handlers>
      <add name="httpPlatformHandler" path="*" verb="*"
           modules="httpPlatformHandler"
           resourceType="Unspecified"/>
    </handlers>
    <httpPlatform processPath="PATH-TO-EXE-HERE"
                  arguments="ANY-COMMAND-LINE-ARGUMENTS-HERE"
                  startupRetryCount="3"
                  stdoutLogEnabled="true"
                  forwardWindowsAuthToken="true"/>
  </system.webServer>
</configuration>

```

HttpPlatformHandler will pass the TCP/IP port for the web application
to listen to for incoming requests in the `HTTP_PLATFORM_PORT`
environment variable, so be sure that you use that in your program or
pass it as an argument to your web application server.
