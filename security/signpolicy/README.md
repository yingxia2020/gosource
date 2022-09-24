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
openssl req -x509 -nodes -days 365 -newkey rsa:3072 -keyout amber-jwt.key -out amber-jwt.crt -config amber-jwt.cnf -sha256
```


### Step two: Use the files generated to sign the policy
```
Usage of ./policy-sign:
  -algorithm string
        Supported algorithm of RSA key pair (default "PS384")
  -certfile string
        Input certificate file to verify the policy
  -policyfile string
        Input policy file to be signed
  -privkeyfile string
        Input private key file to sign the policy
```

#### To create policy JWT without signing algorithm and signature:
```
$ ./policy-sign --policyfile policy.rego 
eyJhbGciOiJub25lIiwidHlwIjoiSldUIn0.eyJBdHRlc3RhdGlvblBvbGljeSI6IlpHVm1ZWFZzZENCdFlYUmphR1Z6WDNObmVGOXdiMnhwWTNrZ1BTQm1ZV3h6WlFwdFlYUmphR1Z6WDNObmVGOXdiMnhwWTNrZ1BTQjBjblZsSUhzS0lDQWdhVzV3ZFhRdVlXMWlaWEpmYzJkNFgyMXlaVzVqYkdGMlpTQTlQU0FpT0RObU5HVTRNVGs0TmpGaFpHVm1ObVptWWpKaE5EZzJOV1ZtWldFNU16TTNZamt4WldRek1HWmhNek0wT1RGaU1UZG1NR1ExWkRsbE9ESXdORFF4TUNJS0lDQWdhVzV3ZFhRdVlXMWlaWEpmYzJkNFgyMXljMmxuYm1WeUlEMDlJQ0k0TTJRM01UbGxOemRrWldGallURTBOekJtTm1KaFpqWXlZVFJrTnpjME16QXpZemc1T1dSaU5qa3dNakJtT1dNM01HVmxNV1JtWXpBNFl6ZGpaVGxsSWdvZ0lDQnBibkIxZEM1aGJXSmxjbDl6WjNoZmFYTjJjSEp2Wkdsa0lEMDlJREFLSUNBZ2FXNXdkWFF1WVcxaVpYSmZjMmQ0WDJsemRuTjJiaUE5UFNBd0NpQWdJR2x1Y0hWMExtRnRZbVZ5WDNObmVGOXBjMTlrWldKMVoyZGhZbXhsSUQwOUlHWmhiSE5sQ24wSyJ9.

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
eyJhbGciOiJQUzM4NCIsInR5cCI6IkpXVCIsIng1YyI6WyJNSUlFTlRDQ0FwMmdBd0lCQWdJQkF6QU5CZ2txaGtpRzl3MEJBUXdGQURCSE1Rc3dDUVlEVlFRR0V3SlZVekVMXG5NQWtHQTFVRUNCTUNVMFl4Q3pBSkJnTlZCQWNUQWxORE1RNHdEQVlEVlFRS0V3VkpUbFJGVERFT01Bd0dBMVVFXG5BeE1GUTAxVFEwRXdIaGNOTWpFd05EQXlNVFF6TWpBeVdoY05Nall3TkRBeU1UUXpNakF5V2pCUU1Rc3dDUVlEXG5WUVFHRXdKVlV6RUxNQWtHQTFVRUNCTUNVMFl4Q3pBSkJnTlZCQWNUQWxORE1RNHdEQVlEVlFRS0V3VkpUbFJGXG5UREVYTUJVR0ExVUVBeE1PUTAxVElGTnBaMjVwYm1jZ1EwRXdnZ0dpTUEwR0NTcUdTSWIzRFFFQkFRVUFBNElCXG5qd0F3Z2dHS0FvSUJnUUNwaTZMRGxiTk5hejVmSjhPSG45QytGSFd2MFRUODhMUlNlbFpwZWxLZG9PVWFXd0d3XG5qRTNlZVk1T0UvMm5PWFZNbStEVjRlaXNxa2ZQNUtwU3o5MWNGaGYwclovclVUYkVpUm9Bd1ltSDBoWmE2YzJZXG5FQTJmN3p2SWwrK1V0aWloanhsK1kvTHVhWnBDMHZTZ3BzcnJCK1RpejlqRHhGWFd3cHorckh1VGE2cEY2Ry8rXG5UclRhQU5RL0VnODEwSXV0cFBjRldyMm9VbHUwbjlZWi81eFBoWXJZdnNDaDdQOWV5V0NremtTbHlVQ0prbEFIXG5oV0VyZUxJNTlieDJGVnhJVzNUTGZBWE4wQ2RsK1hzMkdNdmc0S213L1p6bGFBcXMvcVJWVGJnNUZ6b2JQNjcyXG4zdjhrQzQxT3k3Zi85WDFpY3dEblVvem1RTWR6VEhRcWVEVjlXamtxb0w4U2IwZlRTUmliSWlnK2Z5ZWIvQ2svXG5mRVh1clRNa0tDdmw2YnJEckFTNjFFcXNoRDhVZjBxOFlTcDN2c1B5eUV3OHRkbnFLdTBzTlFNUW1hOXBjRnEyXG5JbTJpWm45NzhmRDQxUDZHZmphWlRUQTJFMUpaS3VVRnJ3WTRkQjlkRVk4K25PdFVDd1dVekVFYkVncm82VUN5XG5KRGVsRU04b2JvUE5NR2NDQXdFQUFhTWpNQ0V3RGdZRFZSMFBBUUgvQkFRREFnRUdNQThHQTFVZEV3RUIvd1FGXG5NQU1CQWY4d0RRWUpLb1pJaHZjTkFRRU1CUUFEZ2dHQkFHUU5rcERtb2JrQXFlNHlhVXkrNFJ6THppNElKdTViXG5lQk8wL25iaWtVc0xWdXNOTTFjazd5WlZnVlBJRVJGMjhvZVpSWStwUW03TTdDM3N0QjRqQkJpR0FtM0FpMzduXG4vZC9JQkhKdTkzaU5pTmxEV3hoRm95c0lwdGg2QXBiS0RqQVdBZ1hqRGwrVU9BelUrOXVPVEV2Zm9wa0Z6dkhnXG5BRGRmYlhwMnhVaWhXZnphaHp4QkdUenFGdUdhL1M3bDc0aE5UWjdORmdmdFZkZzIrZEZyamJlY3pIaXBZY05OXG4rM1JWR3Y1b2grc0VNeHpuc2FJUGtMYm00MTIxNnN4T2phUzVJb2IzZmg2aEM4elRqanJwUzdYZUlvYXV5d2tlXG5BWnQvcEdnTW9kYlljVWVUV0ZlK0JxUEEwbm1aNFlvVi81R3RlRU42SHhvSVJjQkxDVmtIL0ZXd0FGOEFWS1JRXG5tb0NWc2M2R3FtWUdzdVdVU0FWdUpmZXBVOVppREFVdHIva1lPOWw3cjdiNTBkVytWQm1YVVZMN0NYeTkyUkowXG5wYXdmVXVvbWR2bUJJN3FWYVJ4bFNmVGRTZzg2ZDJnQm9oS0hzSndsclJLVTVCclIzaWk2NitLWGNmWjdnMVFJXG5QZU45OWsya2pxVkRqZDBFeEZmKzNrVS9OcGpVY0EybXNBPT0iXX0.eyJBdHRlc3RhdGlvblBvbGljeSI6IlpHVm1ZWFZzZENCdFlYUmphR1Z6WDNObmVGOXdiMnhwWTNrZ1BTQm1ZV3h6WlFwdFlYUmphR1Z6WDNObmVGOXdiMnhwWTNrZ1BTQjBjblZsSUhzS0lDQWdhVzV3ZFhRdVlXMWlaWEpmYzJkNFgyMXlaVzVqYkdGMlpTQTlQU0FpT0RObU5HVTRNVGs0TmpGaFpHVm1ObVptWWpKaE5EZzJOV1ZtWldFNU16TTNZamt4WldRek1HWmhNek0wT1RGaU1UZG1NR1ExWkRsbE9ESXdORFF4TUNJS0lDQWdhVzV3ZFhRdVlXMWlaWEpmYzJkNFgyMXljMmxuYm1WeUlEMDlJQ0k0TTJRM01UbGxOemRrWldGallURTBOekJtTm1KaFpqWXlZVFJrTnpjME16QXpZemc1T1dSaU5qa3dNakJtT1dNM01HVmxNV1JtWXpBNFl6ZGpaVGxsSWdvZ0lDQnBibkIxZEM1aGJXSmxjbDl6WjNoZmFYTjJjSEp2Wkdsa0lEMDlJREFLSUNBZ2FXNXdkWFF1WVcxaVpYSmZjMmQ0WDJsemRuTjJiaUE5UFNBd0NpQWdJR2x1Y0hWMExtRnRZbVZ5WDNObmVGOXBjMTlrWldKMVoyZGhZbXhsSUQwOUlHWmhiSE5sQ24wSyJ9.Z-4zufkwQFdO8vpj6eyGjToBh0zzOGq_fUfRTv2lAUkGzhB_rL8Ma8FDA6s-fVOTB7e_RqnocDBQGOEwgI2GyYFw6hQZgCBB1KiuXn5uJAlj7gxRP24JxJVPhNkiNBPo-xNnunjqw_Gg1zsf4uoa6Wp3ItFtCLe_0EiyW3lVRecMXU1lGIy6BbH3FP7_vM9vi_q8xVs9A8wvbtZd2KncGbZNqKE_oVsTx5JBNS5JOGPMERYukdlGuT_hPbvYIMMtCEBocBARFiaqRgY7nGvCHGI7tRwxYpVF7k8YLTqWz2rwBxA5jhDd9J7AtfMsGhZwOlBG9KWlbOkAreiBcOAPeKOzPczpMInvIZ8nm9n4rdyfbrEhKmE0JmyV2a9CS-sC-X077YnwuWBuoaDYXuyVLTY4SyRn9LrHA2acpdzdZ9c99yM2w5LVZO3c7s2XbOGCWZ9oWmvYwoWFo18RTctCY3w7-ASVixyLpoxdDxYS7gUgATO1VJqLQbcVb9rQ6r_D

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
4. Supported signing algorithms are "RS256", "PS256", "RS384", "PS384", "RS512", "PS512"
5. Include sample codes for how to verify the signed and unsighed policy token and extract rego policy payload
