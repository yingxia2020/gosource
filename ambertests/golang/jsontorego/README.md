### Utility to help convert JSON format policy to Rego format

#### Usage
```
$ ./policy-tool --help
Usage of ./policy-tool:
  -attType string
        Supported attestation type (default "SGX")
  -inputFile string
        Input JSON format policy file (default "policy.json")
  -outputFile string
        Output Rego format policy file  (default "policy.rego")
```

#### Example
```
$ ./policy-tool
############################### Original Format #################################

default matches_sgx_policy = false

matches_sgx_policy = true {
quote = input

amber_sgx_isvsvn_limit := 1
bigger_equal_than(quote.amber_sgx_isvsvn, amber_sgx_isvsvn_limit)

amber_sgx_mrenclaveValues := ["bab91f200038076ac25f87de0ca67472443c2ebe17ed9ba195314e609038f51ab"]
includes_value(amber_sgx_mrenclaveValues, quote.amber_sgx_mrenclave)

amber_sgx_mrsignerValues := ["d412a4f07ef83892a5915fb2ab584be31e186e5a4f95ab5f6950fd4eb8694d7b"]
includes_value(amber_sgx_mrsignerValues, quote.amber_sgx_mrsigner)

amber_sgx_isvprodidValues := [0, 2, 4]
includes_value(amber_sgx_isvprodidValues, quote.amber_sgx_isvprodid)

amber_sgx_is_debuggableValue := [true]
includes_value(amber_sgx_is_debuggableValue, quote.amber_sgx_is_debuggable)

amber_trust_score_limit := "5"
bigger_equal_than(quote.amber_trust_score, amber_trust_score_limit)

amber_tcb_statusValues := ["OK", "CONFIG_NEEDED", "SW_HARDENING_NEEDED"]
includes_value(amber_tcb_statusValues, quote.amber_tcb_status)

}

includes_value(policy_values, quote_value) = true {
        policy_value := policy_values[x]
        policy_value == quote_value
}

bigger_equal_than(value1, value2) = true {
        value1 >= value2
}

#################################################################################
################################ Escaped Format #################################
"default matches_sgx_policy = false\n\nmatches_sgx_policy = true {\nquote = input\n\namber_sgx_isvsvn_limit := 1\nbigger_equal_than(quote.amber_sgx_isvsvn, amber_sgx_isvsvn_limit)\n\namber_sgx_mrenclaveValues := [\"bab91f200038076ac25f87de0ca67472443c2ebe17ed9ba195314e609038f51ab\"]\nincludes_value(amber_sgx_mrenclaveValues, quote.amber_sgx_mrenclave)\n\namber_sgx_mrsignerValues := [\"d412a4f07ef83892a5915fb2ab584be31e186e5a4f95ab5f6950fd4eb8694d7b\"]\nincludes_value(amber_sgx_mrsignerValues, quote.amber_sgx_mrsigner)\n\namber_sgx_isvprodidValues := [0, 2, 4]\nincludes_value(amber_sgx_isvprodidValues, quote.amber_sgx_isvprodid)\n\namber_sgx_is_debuggableValue := [true]\nincludes_value(amber_sgx_is_debuggableValue, quote.amber_sgx_is_debuggable)\n\namber_trust_score_limit := \"5\"\nbigger_equal_than(quote.amber_trust_score, amber_trust_score_limit)\n\namber_tcb_statusValues := [\"OK\", \"CONFIG_NEEDED\", \"SW_HARDENING_NEEDED\"]\nincludes_value(amber_tcb_statusValues, quote.amber_tcb_status)\n\n}\n\nincludes_value(policy_values, quote_value) = true {\n\tpolicy_value := policy_values[x]\n\tpolicy_value == quote_value\n}\n\nbigger_equal_than(value1, value2) = true {\n\tvalue1 >= value2\n}"
#################################################################################
```
#### Supported fields in JSON file
| Field  | Format  | Example |
| :------------ |:---------------:| -----:|
| amber_sgx_isvprodid  | Array of Integer | [2, 4] |
| amber_sgx_isvsvn   |   >=, <=, >, < Integer | >=1 |
| amber_sgx_mrenclave   |  Array of mrenclaves | ["bab91...", "74724..."] |
| amber_sgx_mrsigner |  Array of mrsigners | ["d412a...", "5a4f9..."] |
| amber_sgx_is_debuggable | Boolean | true (false) |
| amber_trust_score | >=, <=, >, < Integer | >=5 |
| amber_tcb_status | Array of String | ["OK", "CONFIG_NEEDED"]** |

** amber_tcb_status accepted values are limited to:
0:    "OK",
1:    "CONFIG_NEEDED",
2:    "OUT_OF_DATE",
3:    "OUT_OF_DATE_CONFIG_NEEDED",
7:    "SW_HARDENING_NEEDED",
8:    "CONFIG_AND_SW_HARDENING_NEEDED"

#### Note
* Converted policy is displayed on screen in both original and escaped rego policy file format
* Converted policy is also stored in file policy.rego or output file name specified in CLI command
* Output policy.rego file could be verified with the following web site with input data input.json: https://play.openpolicyagent.org/
* Right now, only attestation type SGX is supported
