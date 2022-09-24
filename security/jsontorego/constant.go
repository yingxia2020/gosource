/*
 *  Copyright (C) 2022 Intel Corporation
 *  SPDX-License-Identifier: BSD-3-Clause
 */
package main

const (
	OPA_HEADER = `
default matches_sgx_policy = false

matches_sgx_policy = true {
quote = input

`
	OPA_FOOTER = `}

includes_value(policy_values, quote_value) = true {
	policy_value := policy_values[x]
	policy_value == quote_value
}
`
	FUNC_BIGGER = `
bigger_than(value1, value2) = true {
	value1 > value2
}
`
	FUNC_BIGGER_EQUAL = `
bigger_equal_than(value1, value2) = true {
	value1 >= value2
}
`
	FUNC_SMALLER = `
smaller_than(value1, value2) = true {
	value1 < value2
}
`
	FUNC_SMALLER_EQUAL = `
smaller_equal_than(value1, value2) = true {
	value1 <= value2
}
`
	SGX_TYPE         = "sgx"
	TDX_TYPE         = "tdx"
	SGX_MRENCLAVE    = "amber_sgx_mrenclave"
	SGX_MRSINGER     = "amber_sgx_mrsigner"
	SGX_ISVPRODID    = "amber_sgx_isvprodid"
	SGX_ISVSVN       = "amber_sgx_isvsvn"
	SGX_ISDEBUGGABLE = "amber_sgx_is_debuggable"
	TRUST_SCORE      = "amber_trust_score"
	TCB_STATUS       = "amber_tcb_status"
)

const (
	BIG int = 1 << iota
	BIGE
	SMALL
	SMALLE
)
