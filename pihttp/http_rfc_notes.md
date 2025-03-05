## About HTTP request and response formats


https://www.rfc-editor.org/rfc/rfc9110.html

https://www.iana.org/assignments/http-status-codes/http-status-codes.xhtml

// reading HTTP RFC -> https://www.rfc-editor.org/rfc/rfc1945.html#section-5
// SECTION THAT DEFINE A HTTP REQUEST FORMAT

	// if a higher version request is received, the
	//  proxy/gateway must either downgrade the request version or respond
	//  with an error.

  //  After receiving and interpreting a request message, a server responds
	//  in the form of an HTTP response message.
	//  Response        = Simple-Response | Full-Response
	//  Simple-Response = [ Entity-Body ]

	// Simple response should be only from HTTP versions <= 0.9.
	// If a client sends an HTTP/1.0 Full-Request and the server response with
	// with a Status-Line the client should assume that is a Simple-response.

	// Simple response format:
	// Just entity body -> see definition of entity body

	// Full response format:

	// Full-Response   =
	// FirstLine -> Status-Line             ; Section 6.1
	//  General-Header       ; Section 4.3
	//  | Response-Header      ; Section 6.2
	//  | Entity-Header )      ; Section 7.1
	// CRLF
	// [ Entity-Body ]         ; Section 7.2

	// status line
	// Status-Line = HTTP-Version SP Status-Code SP Reason-Phrase CRLF

	// StatusCode element is a 3-digit int

	// he presence of an entity body in a
	//  request is signaled by the inclusion of a Content-Length header field
	//  in the request message headers. HTTP/1.0 requests containing an
	//  entity body must include a valid Content-Length header field.

	// when a request or response has a body it must include the
	// header Content-Type and Content-Encoding.

	// entity-body := Content-Encoding( Content-Type( data ) )

	// Content-Length header must have the lenght of bytes of the entity body

	// if we have a body and dont have the content length header the
	// server must return bad request

	// entity body definition
	// Entity Body

	//  The entity body (if any) sent with an HTTP request or response is in
	//  a format and encoding defined by the Entity-Header fields.
	//      Entity-Body    = *OCTET

	// if a request has a body that means that in the request we have:
	// 1. the http request metthod allows.
	// 2. We have Content-Length header field

	//  For response messages,
	// All responses dependent on request method and response code.
	// All responses to the HEAD request method must not include a body.
	// All status code 1xx (informational), 204 (no content), and 304 (not modified) responses must not include a body.
	// All other responses must include body or a Content-Length header field defined with a value of zero (0).