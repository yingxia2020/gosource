## Utility to create policy file in JSON Web Signature (JWS) format

### References:
- Azure MAA:
  - https://learn.microsoft.com/en-us/azure/attestation/policy-examples
  - https://docs.microsoft.com/en-us/azure/attestation/author-sign-policy
- JWS RFC7515:
  - https://www.rfc-editor.org/rfc/rfc7515


### Step one: Create self signed cert for policy JWS creation:
- Generate config file
```
cat << EOF > amber-jwt.cnf
[req]
distinguished_name = req_distinguished_name
x509_extensions = v3_req
prompt = no
[req_distinguished_name]
C = US
ST = California
L = Santa Clara
O = Intel
OU = SATG
CN = amber.intel.com
[v3_req]
keyUsage = critical, digitalSignature, keyAgreement
extendedKeyUsage = serverAuth
EOF
```

- Generate key and cert files
```
openssl req -x509 -nodes -days 365 -newkey rsa:3072 -keyout amber-jwt.key -out amber-jwt.crt -config amber-jwt.cnf
```


### Step two: Use the files generated to sign the policy
```
Usage of ./policy-sign:
  -algorithm string
        Supported algorithm of policy signing key pair (default "PS384")
  -certfile string
        Required. Path to PEM formatted file that contains your signing certificate
  -policyfile string
        Required. Path to text policy file to be signed into a JWT format Amber policy
  -privkeyfile string
        Required. Path to PEM formatted file that contains your 2048 bit RSA private key
```

#### To create policy JWT without signing algorithm and signature:
```
$ ./policy-sign --policyfile policy.rego
eyJhbGciOiJub25lIn0.eyJBdHRlc3RhdGlvblBvbGljeSI6IlpHVm1XXX.

Policy payload:
default matches_sgx_policy = false
matches_sgx_policy = true {
   input.amber_sgx_mrenclave == "83f4e819861adef6ffb2a4865efea9337b91ed30fa33491b17f0d5d9e8204410"
   input.amber_sgx_mrsigner == "83d719e77deaca1470f6baf62a4d774303c899db69020f9c70ee1dfc08c7ce9e"
   input.amber_sgx_isvprodid == 0
   input.amber_sgx_isvsvn == 0
   input.amber_sgx_is_debuggable == false
}

Token is verified: true
```

#### To create policy JWT with default signing algorithm PS384 and signature:
```
$ ./policy-sign --policyfile policy.rego --privkeyfile amber-jwt.key --certfile amber-jwt.crt
eyJhbGciOiJSUzM4NCIsIng1XXX.YyI6WyJNSUlFdGpDQ0F4NmdBd0lCQWdJVWEwQWXXX.FGempmdGpvT2hkbGN3d3dEUVlKS29aSWh2Y05BUUVMXG5CUXXX

Policy payload:
default matches_sgx_policy = false
matches_sgx_policy = true {
   input.amber_sgx_mrenclave == "83f4e819861adef6ffb2a4865efea9337b91ed30fa33491b17f0d5d9e8204410"
   input.amber_sgx_mrsigner == "83d719e77deaca1470f6baf62a4d774303c899db69020f9c70ee1dfc08c7ce9e"
   input.amber_sgx_isvprodid == 0
   input.amber_sgx_isvsvn == 0
   input.amber_sgx_is_debuggable == false
}

Token is verified: true
```

### Note:
1. Signed policy token could be self verified at jwt.io
2. Output file of this program is input policy file with ".signed" extension
3. Policy payload Amber uses rego format which is different from Azure MAA
4. Supported signing algorithms are "RS256", "PS256", "RS384", "PS384", default algorithm is PS384
5. Include sample codes for how to verify the signed and unsighed policy token and extract rego policy payload
