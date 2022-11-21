
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
