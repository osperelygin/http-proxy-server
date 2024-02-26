# http-proxy-server

## Description:

http-proxy-server proxies http and https requests

webapi implemented command injection vulnerability scanner and API for show request/response history and repeat request

### API

Parsing requests and response:
* HTTP method
* Path and GET params 
* Headers, while separately parsing Cookies
* Request body, in case of application/x-www-form-urlencoded separate POST parameters
* Compression methods are disabled on proxy-server side by header `Accept-Encoding:identity;q=0` 

### API Description

| Method | Path                   | Description
| ---    | ---                    | --- |
| GET    |`/api/v1/requests`      | List of requests; |
| GET    |`/api/v1/response/{id}` | Output response for request with specified id; |
| GET    |`/api/v1/requests/{id}` | Output request; |
| GET    |`/api/v1/repeat/{id}`   | Resubmit request; |
| GET    |`/api/v1/scan/{id}`     | Request vulnerability scanner (param-miner); |

### Command injection vulnerability scanner

Scanner insert into GET/POST/Сookie/HTTP headers one by one next strings:

`;cat /etc/passwd;`

`|cat /etc/passwd|`

`` `cat /etc/passwd\` `` 

and do request by proxy server, then search into given response `"root:"` substring. If response contains this substring, scanner write that request are vulnerable into webapi response.

### Dependencies

* openssl
* docker, docker compose

### How to start

Execute next commands:

1. Generate certificate authority cert and key 

    `make gen-ca`

2. Add generate self-signed cert to system trusted list for correct cURL work

* on macos: 

    `sudo security add-trusted-cert -d -r trustRoot -k /Library/Keychains/System.keychain "$(pwd)/certs/ca.crt" `

* on linux:

    `sudo cp ./certs/ca.crt /usr/local/share/ca-certificates/ && sudo update-ca-certificates`

3. Run proxy server

    `make up-proxy-server`

4. Try examples

5. Shutdown proxy-server

    `make down-proxy-server`

### Examples

Proxy HTTP request:

`curl -v -x http://localhost:8080 'http://example.com'`

Proxy HTTPS request:

`curl -v -x http://localhost:8080 'https://example.com'`
