#include <chrono>
#include <cassert>
#include <fstream>
#include <iostream>
#include <string>
#include <jwt-cpp/jwt.h>

int main() {

    const std::string MRENCLAVE = "mrenclave";
    const std::string REPORTDATA = "reportdata";

    std::string rsa_pub_key = R"(-----BEGIN CERTIFICATE-----
MIICjzCCAjSgAwIBAgIUImUM1lqdNInzg7SVUr9QGzknBqwwCgYIKoZIzj0EAwIw
aDEaMBgGA1UEAwwRSW50ZWwgU0dYIFJvb3QgQ0ExGjAYBgNVBAoMEUludGVsIENv
cnBvcmF0aW9uMRQwEgYDVQQHDAtTYW50YSBDbGFyYTELMAkGA1UECAwCQ0ExCzAJ
BgNVBAYTAlVTMB4XDTE4MDUyMTEwNDUxMFoXDTQ5MTIzMTIzNTk1OVowaDEaMBgG
A1UEAwwRSW50ZWwgU0dYIFJvb3QgQ0ExGjAYBgNVBAoMEUludGVsIENvcnBvcmF0
aW9uMRQwEgYDVQQHDAtTYW50YSBDbGFyYTELMAkGA1UECAwCQ0ExCzAJBgNVBAYT
AlVTMFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAEC6nEwMDIYZOj/iPWsCzaEKi7
1OiOSLRFhWGjbnBVJfVnkY4u3IjkDYYL0MxO4mqsyYjlBalTVYxFP2sJBK5zlKOB
uzCBuDAfBgNVHSMEGDAWgBQiZQzWWp00ifODtJVSv1AbOScGrDBSBgNVHR8ESzBJ
MEegRaBDhkFodHRwczovL2NlcnRpZmljYXRlcy50cnVzdGVkc2VydmljZXMuaW50
ZWwuY29tL0ludGVsU0dYUm9vdENBLmRlcjAdBgNVHQ4EFgQUImUM1lqdNInzg7SV
Ur9QGzknBqwwDgYDVR0PAQH/BAQDAgEGMBIGA1UdEwEB/wQIMAYBAf8CAQEwCgYI
KoZIzj0EAwIDSQAwRgIhAOW/5QkR+S9CiSDcNoowLuPRLsWGf/Yi7GSX94BgwTwg
AiEA4J0lrHoMs+Xo5o/sX6O9QWxHRAvZUGOdRQ7cvqRXaqI=
-----END CERTIFICATE-----)";

    std::string token = "eyJhbGciOiJSUzM4NCIsImtpZCI6IjBjZWFiOThlMWEyZjZhMDI4NThmZTczNWFkMzFhYWQ4Y2Q0M2Q4YWUiLCJ0eXAiOiJKV1QifQ.eyJtcmVuY2xhdmUiOiIwMWU5ZjhmZWFmZThlYmY3MDE5MjA4MDQ3OGQyY2NkMDFkYWExNmE2MjE1ZTcxZDRkY2RiNzNhNWJmZWRjN2M2IiwibXJzaWduZXIiOiJjOWUyNGE3MTgzMjg4MTE5NDcxYTc1N2MyNzZlZmMwNTEyOWI5MWY3MjE5NjlkNTI3YTM4NGMyZGQ3MGY5OTExIiwicmVwb3J0ZGF0YSI6IjYxNjI2MzY0NjUwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAiLCJpc3Zwcm9kaWQiOjAsImlzdnN2biI6MCwidGVlX2hlbGRfZGF0YSI6IllXSmpaR1U9IiwicG9saWN5X2lkcyI6bnVsbCwidGNiX3N0YXR1cyI6IlNXX0hBUkRFTklOR19ORUVERUQiLCJ0ZWUiOiJTR1giLCJ2ZXIiOiIxLjAiLCJleHAiOjE2NTg0NTk0NjEsImlhdCI6MTY1ODQ1OTEzMSwiaXNzIjoiQVMgQXR0ZXN0YXRpb24gVG9rZW4gSXNzdWVyIn0.fnLZ3CjWaVTu_NKmlDaYYsQhuKlWytxQsHe1RUSPihWxqbLrbSks4vbRrw88eU1r08JUM0KOTQsK6B9Ta4OZp22Fs6DO1oJ5S5GcoeqxZ_Q94vEbn00e7Qc6HAK9-2Q_vFWNKntP4gF_AFktJjT9mh-ii1iIUsNoYcvKJN-ls78sjCJ0et3QgOYMDl2bjLn0t0WsMQJP84af3g8e5pvg27okoDHMCArdzeC8fclkgk5PpAi2ds_Q4GmlDGr_VSQYW7GxFCXFx1NjqsRJd9powdZ9nX-tQxmH-WaSpxa7GJmm3KaUMoV2v7amUucQkx6P_mpbYsTeYtGq9Kgj-kAnTKSgNmkpMI0Iku_vb3YGVexykpS5HOUl9AJjFHRYxz4WC1uVRX5oSIOdZK4BlMcZjl2BGJT7RBMakg8WWJ-0dg9lUjrXkPqxbD22z4sZ4poZPZ4fKhJSIIRjl9EQzpvxiXexavNQ6XHtQy6GCQ73In89W4CniLthpFDIHzTi-u55";
    std::string mrenclave = "01e9f8feafe8ebf70192080478d2ccd01daa16a6215e71d4dcdb73a5bfedc7c6";

    auto verify = jwt::verify().allow_algorithm(jwt::algorithm::rs384(rsa_pub_key, "", "", "")).with_issuer("AS Attestation Token Issuer");

    auto decoded = jwt::decode(token);

    try {
        verify.verify(decoded);
    } catch (...) {
        std::cerr << "verification fail!\n";
    }
    for (auto& e : decoded.get_header_claims())
        std::cout << e.first << " = " << e.second.to_json() << std::endl;


    std::string amber_mrenclave;
    std::string amber_report_data;
    for (auto& e : decoded.get_payload_claims()) {
        if (std::string(e.first) == MRENCLAVE) {
            amber_mrenclave = e.second.to_json().get<std::string>();
        }
        if (std::string(e.first) == REPORTDATA) {
            amber_report_data = e.second.to_json().get<std::string>();
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

    std::cout << amber_report_data << std::endl;
}
