/*
 *  Copyright (C) 2022 Intel Corporation
 *  SPDX-License-Identifier: BSD-3-Clause
 */

package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Description struct {
	PolicyType string `json:"policy_type"`
	Label      string `json:"label"`
}

type Metadata struct {
	Version     float32     `json:"version"`
	Id          string      `json:"id"`
	Description Description `json:"description"`
}

type SGX struct {
	Isvprodid    []int    `json:"amber_sgx_isvprodid,omitempty"`
	Isvsvn       string   `json:"amber_sgx_isvsvn,omitempty"`
	Mrenclave    []string `json:"amber_sgx_mrenclave,omitempty"`
	Mrsigner     []string `json:"amber_sgx_mrsigner,omitempty"`
	IsDebuggable *bool    `json:"amber_sgx_is_debuggable,omitempty"`
	TrustScore   string   `json:"amber_trust_score,omitempty"`
	TcbStatus    []string `json:"amber_tcb_status,omitempty"`
}

type SgxPolicy struct {
	Metadata Metadata `json:"meta"`
	SGX      SGX      `json:"sgx"`
}

type Policy struct {
	SgxPolicy SgxPolicy `json:"policy"`
	Signature string    `json:"signature,omitempty"`
}

var (
	inputFile    = flag.String("inputFile", "policy.json", "Input JSON format policy file")
	outputFile   = flag.String("outputFile", "policy.rego", "Output Rego format policy file ")
	attType      = flag.String("attType", "SGX", "Supported attestation type")
	usedFuncFlag = 0
	usedFuncMap  = make(map[int]string)
)

/*
Usage:
./main -inputFile=JSON_file -outputFile=OPA_file -attType=SGX(TDX)

Input JSON file such as:
{
	"policy": {
		"meta": {
			"version": 1.0,
			"id": "unique identifier for the policy, UUIDv4",
			"description": {
				"policy_type": "SGX Enclave",
				"label": "unique label to be assigned to each policy"
			}
		},
		"sgx": {
			"amber_sgx_isvprodid": [0, 2, 4],
			"amber_sgx_isvsvn": ">=1",
			"amber_sgx_mrenclave": ["bab91f200038076ac25f87de0ca67472443c2ebe17ed9ba195314e609038f51ab"],
			"amber_sgx_mrsigner": ["d412a4f07ef83892a5915fb2ab584be31e186e5a4f95ab5f6950fd4eb8694d7b"],
			"amber_sgx_is_debuggable": true,
			"amber_trust_score": ">=5",
			"amber_tcb_status": ["CONFIG_NEEDED", "SW_HARDENING_NEEDED"]
		}
	}
}

Output file could be verified at:
https://play.openpolicyagent.org/

With input file such as:
{
    "amber_sgx_mrenclave": "bab91f200038076ac25f87de0ca67472443c2ebe17ed9ba195314e609038f51ab",
    "amber_sgx_mrsigner": "d412a4f07ef83892a5915fb2ab584be31e186e5a4f95ab5f6950fd4eb8694d7b",
    "amber_sgx_isvsvn": 1,
    "amber_sgx_isvprodid": 0,
    "amber_sgx_is_debuggable": true,
    "amber_trust_score": "5",
    "amber_tcb_status": "CONFIG_NEEDED"
}
*/

func main() {
	flag.Parse()
	if strings.ToLower(*attType) != SGX_TYPE {
		fmt.Println("Error: Only SGX type is supported now")
		return
	}

	/* set initial values for map */
	usedFuncMap[BIG] = FUNC_BIGGER
	usedFuncMap[BIGE] = FUNC_BIGGER_EQUAL
	usedFuncMap[SMALL] = FUNC_SMALLER
	usedFuncMap[SMALLE] = FUNC_SMALLER_EQUAL

	if len(*inputFile) == 0 || len(*outputFile) == 0 {
		fmt.Println("Error: input|output file not set")
		return
	}

	data, err := os.ReadFile(*inputFile)
	if err != nil {
		fmt.Printf("Error: failed to open %s file\n", *inputFile)
		return
	}

	policy := Policy{}
	err = json.Unmarshal([]byte(data), &policy)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	/* Compose OPA policy string */
	var output strings.Builder

	output.WriteString(OPA_HEADER)

	if len(policy.SgxPolicy.SGX.Isvsvn) > 0 {
		output.WriteString(convertLimitValue(policy.SgxPolicy.SGX.Isvsvn, SGX_ISVSVN, false))
	}

	if len(policy.SgxPolicy.SGX.Mrenclave) > 0 {
		output.WriteString(convertIncludeArray(stringToInterfaceSlice(policy.SgxPolicy.SGX.Mrenclave), SGX_MRENCLAVE))
	}

	if len(policy.SgxPolicy.SGX.Mrsigner) > 0 {
		output.WriteString(convertIncludeArray(stringToInterfaceSlice(policy.SgxPolicy.SGX.Mrsigner), SGX_MRSINGER))
	}

	if len(policy.SgxPolicy.SGX.Isvprodid) > 0 {
		output.WriteString(convertIncludeArray(intToInterfaceSlice(policy.SgxPolicy.SGX.Isvprodid), SGX_ISVPRODID))
	}

	if policy.SgxPolicy.SGX.IsDebuggable != nil {
		output.WriteString(convertIncludeValue(*policy.SgxPolicy.SGX.IsDebuggable, SGX_ISDEBUGGABLE))
	}

	if len(policy.SgxPolicy.SGX.TrustScore) > 0 {
		output.WriteString(convertLimitValue(policy.SgxPolicy.SGX.TrustScore, TRUST_SCORE, true))
	}

	if len(policy.SgxPolicy.SGX.TcbStatus) > 0 {
		output.WriteString(convertIncludeArray(stringToInterfaceSlice(policy.SgxPolicy.SGX.TcbStatus), TCB_STATUS))
	}
	output.WriteString(OPA_FOOTER)

	tester := BIG
	for i := 0; i < 4; i++ {
		if usedFuncFlag&tester == tester {
			output.WriteString(usedFuncMap[tester])
		}
		tester = tester << 1
	}

	/* Write to file */
	err = os.WriteFile(*outputFile, []byte(output.String()), 0644)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	/* And output to terminal */
	fmt.Println("############################### Original Format #################################")
	fmt.Println(output.String())
	fmt.Println("#################################################################################")
	fmt.Println("################################ Escaped Format #################################")
	fmt.Println(strconv.Quote(strings.TrimSpace(output.String())))
	fmt.Println("#################################################################################")
}
