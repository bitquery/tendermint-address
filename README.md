# PubKeyMultisigThreshold tendermint address encoder

Simple web service to encode PubKeyMultisigThreshold tendermint addresses, for example, for 
Cosmos blockchain.

To calculate PubKeyMultisigThreshold address in Tendermint, need to get all public keys, 
encode as structure using amino encoding and calculate SHA256, and truncate to 
20 bytes.

As Amino encoding only implemented on Go, web service is needed to calculate addresses on other languages.


## Running

```
go get github.com/tendermint/tendermint/crypto
go get github.com/tendermint/tendermint/crypto/multisig
go get github.com/tendermint/tendermint/crypto/secp256k1
go build service.go
 ./service -port=8080
```

Commend line arguments are optional:

- -port=port for service to listen
- -bind=IP address to bind to ( default is ANY )

## Query

Use GET request to /multisig path with pubkey parameter, for example:

```
http://localhost:8080/multisig?pubkey={%20%22type%22:%20%22tendermint/PubKeyMultisigThreshold%22,%20%22value%22:%20{%20%22threshold%22:%20%222%22,%20%22pubkeys%22:%20[%20{%20%22type%22:%20%22tendermint/PubKeySecp256k1%22,%20%22value%22:%20%22AvxcNJmBbtUaDMC4gIo/zusV4hBRKXKMZ1DEOARD2FzZ%22%20},%20{%20%22type%22:%20%22tendermint/PubKeySecp256k1%22,%20%22value%22:%20%22A07oUd5VEbNBor8By/2ROFRHzBvSuc9%2BVyg%2BF1F51xy0%22%20}%20]%20}%20}
```

results in 

```json
{
  "address": "cafdee438702afd1a184e43f500864cb0e63acae"
}
```

Address then can be converted to string represnetation using BEP32 encoding
