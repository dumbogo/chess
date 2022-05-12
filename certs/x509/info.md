# Certificates, Howto

In this sample, we're going to generate self-signed certificates with our own CA


## Create Root Certificate CA

### Root key
```sh
$ openssl ecparam -out contoso.key -name prime256v1 -genkey
```
outputs: `contoso.key`

### CSR for root certificate
```sh
$ openssl req -new -sha256 -key contoso.key -out contoso.csr
```
outputs: `contoso.csr`

### Root certificate
```sh
$ openssl x509 -req -sha256 -days 365 -in contoso.csr -signkey contoso.key -out contoso.crt
```
outputs: `contoso.crt`


## Create certificates for server:

### Server certificate key
```sh
$ openssl ecparam -out fabrikam.key -name prime256v1 -genkey
```
outputs: `fabrikam.key`

### CSR for Server
```sh
$ openssl req -new \
	-subj "/C=US/ST=Utah/L=Lehi/O=fabrikam, Inc./OU=IT/CN=www.fabrikam.com" \
	-addext "subjectAltName = DNS:www.fabrikam.com" \
	-sha256 -key fabrikam.key -out fabrikam.csr
```
outputs: `fabrikam.csr`


### Server certificate with CSR and signed with CA's root key
```sh
$ openssl x509 -req  \
	-extfile <(printf "subjectAltName=DNS:www.fabrikam.com,DNS:www.fabrikam.com") \
	-in fabrikam.csr -CA  contoso.crt -CAkey contoso.key -CAcreateserial -out fabrikam.crt -days 365 -sha256
```
outputs: `fabrikam.crt`

## how to verify:
```sh
$ openssl x509 -in fabrikam.crt -text -noout
```


One important note, experienced an issue with a deprecated flag needed to be updated, more info on the link below:
https://jfrog.com/knowledge-base/general-what-should-i-do-if-i-get-an-x509-certificate-relies-on-legacy-common-name-field-error/


TODO: Check for doc above on how to configure client CA certificate chain to use with client authentication
https://docs.microsoft.com/en-us/azure/application-gateway/mutual-authentication-certificate-management

Read:
- https://docs.microsoft.com/en-us/azure/application-gateway/ssl-overview
- https://docs.microsoft.com/en-us/windows-hardware/drivers/install/trusted-root-certification-authorities-certificate-store
