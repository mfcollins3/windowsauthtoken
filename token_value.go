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

package windowsauthtoken

// TokenValue is an enumeration that specifies which user token value should
// be passed to the callback handler by the Windows authentication token
// handler. The current options are either to send the DOMAIN\Username vale
// or the user's SID to the callback handler.
type TokenValue int

const (
	// TokenUsername indicates that the Windows authentication token handler
	// will invoke the callback handler with the Windows username
	// (DOMAIN\Username) for the authenticated user.
	TokenUsername TokenValue = iota

	// TokenSid indicates that the Windows authentication token handler will
	// invoke the callback handler with the authenticated user's SID as a
	// string.
	TokenSid
)
