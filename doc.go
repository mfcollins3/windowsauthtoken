// Copyright 2016 Michael F. Collins, III
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this softwar and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING
// FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER
// DEALINGS IN THE SOFTWARE.

// Package windowsauthtoken implements Windows user authentication using
// Windows domain credentials for Go web applications. Package
// windowsauthtoken is designed to work when the web application is being
// hosted by IIS on a Windows Server connected to a Windows domain using
// the HttpPlatformHandler module to forward requests from IIS to the web
// application.
//
// HttpPlatformHandler supports a configuration option that forwards the
// handle for the Windows token for the authenticated user to the Go web
// application by attaching the X-IIS-WindowsAuthToken HTTP header to the
// forwarded request. Package windowsauthtoken implements a middleware
// handler that will execute on every request to extract the handle from
// the HTTP header, obtain the name of the authenticated user, and make
// the Windows username available to the web application.
//
// HttpPlatformHandler will pass the handle on every request that is
// forwarded to the web application, and it is the responsibility of the web
// application to close the handle at the end of the request. This is handled
// automatically by the middleware. But because the handle is passed on every
// request, it is recommended that you wrap the http.DefaultServeMux handler
// or your root handler with the middleware. For example:
//
//     package main
//
//     import (
//         "net/http"
//
//         "github.com/mfcollins3/windowsauthtoken"
//     )
//
//     func main() {
//         // TODO: register HTTP handlers with net/http
//
//         rootHandler := windowsauthtoken.Handler(
//             http.DefaultServeMux,
//             func(username string) error {
//                 // Store username somewhere for the request. For example,
//                 // you can use Gorilla Toolkit's context package to store
//                 // the username in the request.
//                 return nil
//             })
//         log.Fatal(http.ListenAndServe(":8080", rootHandler))
//     }
package windowsauthtoken
