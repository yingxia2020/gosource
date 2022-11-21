#include <chrono>
#include <cassert>
#include <fstream>
#include <iostream>
#include <string>
#include <jwt-cpp/jwt.h>

using namespace picojson;

int main() {

    const std::string MRENCLAVE = "amber_sgx_mrenclave";
    const std::string REPORTDATA = "amber_report_data";
    const std::string MATCHEDPOLICY = "amber_matched_policy_ids";

    std::string rsa_pub_key = R"(-----BEGIN PUBLIC KEY-----
MIIBojANBgkqhkiG9w0BAQEFAAOCAY8AMIIBigKCAYEAqYuiw5WzTWs+XyfDh5/Q
vhR1r9E0/PC0UnpWaXpSnaDlGlsBsIxN3nmOThP9pzl1TJvg1eHorKpHz+SqUs/d
XBYX9K2f61E2xIkaAMGJh9IWWunNmBANn+87yJfvlLYooY8ZfmPy7mmaQtL0oKbK
6wfk4s/Yw8RV1sKc/qx7k2uqRehv/k602gDUPxIPNdCLraT3BVq9qFJbtJ/WGf+c
T4WK2L7Aoez/XslgpM5EpclAiZJQB4VhK3iyOfW8dhVcSFt0y3wFzdAnZfl7NhjL
4OCpsP2c5WgKrP6kVU24ORc6Gz+u9t7/JAuNTsu3//V9YnMA51KM5kDHc0x0Kng1
fVo5KqC/Em9H00kYmyIoPn8nm/wpP3xF7q0zJCgr5em6w6wEutRKrIQ/FH9KvGEq
d77D8shMPLXZ6irtLDUDEJmvaXBatiJtomZ/e/Hw+NT+hn42mU0wNhNSWSrlBa8G
OHQfXRGPPpzrVAsFlMxBGxIK6OlAsiQ3pRDPKG6DzTBnAgMBAAE=
-----END PUBLIC KEY-----)";

    std::string token = "eyJhbGciOiJQUzM4NCIsImprdSI6Imh0dHBzOi8vd3d3LmludGVsLmNvbS9hbWJlci9jZXJ0cyIsImtpZCI6IjBjZWFiOThlMWEyZjZhMDI4NThmZTczNWFkMzFhYWQ4Y2Q0M2Q4YWUiLCJ0eXAiOiJKV1QifQ.eyJhbWJlcl90cnVzdF9zY29yZSI6MTAsImFtYmVyX3JlcG9ydF9kYXRhIjoiMzE5MWMwNjhkODkzMzdlZDg5ZGJiNTZiNjlkNzE3Y2RjZWE2ZTg1NDljYWNjM2M0ZTZlMDUyZWIxZjVlNWNmYjAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAiLCJhbWJlcl90ZWVfaGVsZF9kYXRhIjoiTFMwdExTMUNSVWRKVGlCRFJWSlVTVVpKUTBGVVJTMHRMUzB0Q2sxSlNVSlBla05DZHpaQlJFRm5SVU5CYUZKNlZXUXZjMFpCUTBKRlVuZDFSVk5PU0d0b2QyOUhkV1pOZFVSQlMwSm5aM0ZvYTJwUFVGRlJSRUY2UVVFS1RVSTBXRVJVU1hsTlJHZDVUMVJKZUU1VVVURk9iRzlZUkZSSmVVMVVTWGRPZWtsNFRsUlJNVTVzYjNkQlJFSXlUVUpCUjBKNWNVZFRUVFE1UVdkRlJ3cENVM1ZDUWtGQmFVRXlTVUZDUVhGWmFWQnNTR3RZV0ZRdk5Wb3paRGhOTkVWVmMyZHFkRGQ0VUhKeFJWTXdhVVJXY25Sa1QyZEZXRUYyUlZaSk5DdFRDbVZTY1hoaVJYbDJSelpYZDA0MVZuTk9UelExUjFCbVR6aGxVWEJWVVVSR01HUnZRbU14UTFGeldFVnhNVllyYkc5S1UzcFpXamRqUkU5bmNFMHlibk1LVVdGdFIwaDJUMHBDVVhnNE5ucEJTMEpuWjNGb2EycFBVRkZSUkVGM1RtNUJSRUpyUVdwQlZHVkZkRE5OUmt4aGVucEtiRzA1WkM4MWNYRkliREpqVVFvMlVqZDFiaXRrUlM5MFlWcE9lRE5hVkZOelFYQkROVk5PTVRsRU9USjVlbmhPVUhCTlRYZERUVUZpVlZvMGQzZFliVFpzU0VGdmFVeDVPV0p6YXpWMkNqSXZTQ3RrY0RaVmNIVnNNR1ZCTjNGSVEyaEZiazQyTkZWNlZuaGlkVXRrWW1WM0szZFFVM0V4WnowOUNpMHRMUzB0UlU1RUlFTkZVbFJKUmtsRFFWUkZMUzB0TFMwSyIsImFtYmVyX3NneF9tcmVuY2xhdmUiOiI3NzFkZjA3YzEwMjI5Y2FiMDEwNmE1MTg3YmZlZTQ1NDQxYTRjZjQ3ZTQ2YzRlZjdhNTZjNDhlODI0YmJhMzMzIiwiYW1iZXJfc2d4X2lzX2RlYnVnZ2FibGUiOmZhbHNlLCJhbWJlcl9zZ3hfbXJzaWduZXIiOiI5NjVlNTAyNmEwY2VmZDI3NzkyYjk5ZDUwZDIyNWJmNjBiNTNiYWU0MDA1ZDJmYzExZDY3ZDhmZjc0MmMyZWM1IiwiYW1iZXJfc2d4X2lzdnByb2RpZCI6MCwiYW1iZXJfc2d4X2lzdnN2biI6MCwiYW1iZXJfbWF0Y2hlZF9wb2xpY3lfaWRzIjpbIjQ3NThkMTg3LWI5OTAtNGRlNC1hODkxLTQ0OGI3NTRkY2ZhNyJdLCJhbWJlcl90Y2Jfc3RhdHVzIjoiT0siLCJhbWJlcl9ldmlkZW5jZV90eXBlIjoiU0dYIiwiYW1iZXJfY2xpZW50X25vbmNlIjp0cnVlLCJhbWJlcl9jdXN0b21fcG9saWN5Ijp7fSwidmVyIjoiMS4wIiwiZXhwIjoxNjYyNTA0NjAyLCJqdGkiOiI4Y2U0OGRhMy0zNTVkLTQwNTAtODZiZi03MWIwNzExYmRjMjAiLCJpYXQiOjE2NjI1MDQyNzIsImlzcyI6IkFTIEF0dGVzdGF0aW9uIFRva2VuIElzc3VlciJ9.bClFDipBSgOwGvOZDs7Ls51v6Q9wMv4KKufNWI3hRNradQcN53PbeOqS25ugMpBbEOx3xMiSbbfw-3RshRotQWguiICOa0LMeKmExAvbyoZiK5r4UfqHou5fO4La_Q22nTawDGrse-J5gq_bmS1vIndugYlv2MqXAOXnl1li2cBISiLPCXw5qAbJ9nr0JNywDXARMSjTYsJ5E5AbcUqBVgmH-6Ej3W4YFlg61Nc1wQ3GEmX186wsraf5ipPqOh6WrNsfli9uPe5Gslkvt4C-q94XUH2ZrN26lTP4ERYj-k3lmCS5LgTuOKef4x0TNGCbt5SO1rMaSwD_kfAbJV4ebuB_xAkUk7bG7OS6SBG2ITgmVQ5h2VunPSW6kX0FiXTsc-ley0fDgN7VAebfqLrvnK0wX4TO-2BRnFEsBQKiLlng9yg0P4uVBEgfEwX1AYk3EyMuBkRhNzTzfx7tWZlm6AZpjzQt8zuzkDFEHfuHjgEN3FAHA5Jw_m7Uo3M0wcqM";

    std::string mrenclave = "771df07c10229cab0106a5187bfee45441a4cf47e46c4ef7a56c48e824bba333";

    auto decoded = jwt::decode(token);
    auto verify = jwt::verify().allow_algorithm(jwt::algorithm::ps384(rsa_pub_key, "", "", "")).with_issuer("AS Attestation Token Issuer");

    try {
        verify.verify(decoded);
    } catch (...) {
        std::cerr << "verification fail!\n";
    }
    for (auto& e : decoded.get_header_claims())
        std::cout << e.first << " = " << e.second.to_json() << std::endl;

    std::string amber_mrenclave;
    std::string amber_report_data;
    int matched_policy_id_num = 0;

    for (auto& e : decoded.get_payload_claims()) {
        if (std::string(e.first) == MRENCLAVE) {
            amber_mrenclave = e.second.to_json().get<std::string>();
        }
        if (std::string(e.first) == REPORTDATA) {
            amber_report_data = e.second.to_json().get<std::string>();
        }
        if (std::string(e.first) == MATCHEDPOLICY) {
            matched_policy_id_num = e.second.to_json().get<array>().size();
        }
    }

    if ( amber_mrenclave.size() == 0 ) {
        std::cerr << "Invalid Amber Token. MREnclave value not found" << std::endl;
    }

    if ( amber_report_data.size() == 0 ) {
        std::cerr << "Invalid Amber Token. Report data not found" << std::endl;
    }

    if ( amber_mrenclave != mrenclave)
    {
        std::cerr << "Mismatch mrenclave" << std::endl;
    }

    if ( matched_policy_id_num == 0 ) {
        std::cerr << "Invalid Amber Token. No matched policy found" << std::endl;
    }
    std::cout << "Matched policy number: " << matched_policy_id_num << std::endl;
    std::cout << "Report data: " << amber_report_data << std::endl;
}

