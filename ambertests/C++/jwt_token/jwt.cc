#include <chrono>
#include <cassert>
#include <fstream>
#include <string>
#include <iostream>
#include "jwt/jwt.hpp"

/***
 * STEPS TO GENERATE RSA PRIVATE PUBLIC KEYPAIR.
 *
 * 1. openssl genrsa -out jwtRS256.key 1024
 * 2. openssl rsa -in jwtRS256.key -pubout -outform PEM -out jwtRS256.key.pub
 */
/* Amber JWT token payload format is like below:
 * {"exp":1658266968,"iat":1658266638,"iss":"AS Attestation Token Issuer","isvprodid":0,"isvsvn":0,
 * "mrenclave":"01e9f8feafe8ebf70192080478d2ccd01daa16a6215e71d4dcdb73a5bfedc7c6",
 * "mrsigner":"c9e24a7183288119471a757c276efc05129b91f721969d527a384c2dd70f9911",
 * "policy_ids":null,"reportdata":"6162636465000000000000000000000000000000000000000000000000000000",
 * "tcb_status":"SW_HARDENING_NEEDED","tee":"SGX","ver":"1.0"}
 */

std::string read_from_file(const std::string& path)
{
    std::string contents;
    std::ifstream is{path, std::ifstream::binary};

    if (is) {
        // get length of file:
        is.seekg (0, is.end);
        auto length = is.tellg();
        is.seekg (0, is.beg);
        contents.resize(length);

        is.read(&contents[0], length);
        if (!is) {
            is.close();
            return {};
        }
    } else {
        std::cerr << "FILE not FOUND!!" << std::endl;
    }

    is.close();
    return contents;
}

int main(int argc, char** argv) {
    using namespace jwt::params;

    std::string token_path =  "token.txt";
    std::string pub_key_path = "mypubkey.pem";
    bool verify_token = false;

    if ( argc >= 2 ) {
        if (strcmp(argv[1], "true") == 0) {
            verify_token = true;
        }
    }
    if ( argc == 4 ) {
        token_path = std::string(argv[2]);
        pub_key_path = std::string(argv[3]);
    }

    auto jwt_token = read_from_file(token_path);
    auto pub_key = read_from_file(pub_key_path);
    std::string mrenclave = "01e9f8feafe8ebf70192080478d2ccd01daa16a6215e71d4dcdb73a5bfedc7c6";

    try {
        auto dec_obj = jwt::decode(jwt_token, algorithms({"RS384"}), verify(verify_token), secret(pub_key));
        std::cout << dec_obj.header() << std::endl;
        std::cout << dec_obj.payload() << std::endl;
        if (dec_obj.payload().has_claim_with_value("mrenclave", mrenclave)) {
            std::cout << "mrenclave is match\n";
        } else {
            std::cerr << "mrenclave is not match\n";
        }
        std::cout << "JWT token is parsed successfully!\n";
    } catch (const jwt::TokenExpiredError& e) {
        //Handle Token expired exception here
        std::cerr << "Token expired exception found\n";
    } catch (const jwt::SignatureFormatError& e) {
        //Handle invalid signature format error
        std::cerr << "Signature format error exception found\n";
    } catch (const jwt::DecodeError& e) {
        //Handle all kinds of other decode errors
        std::cerr << "Decode error exception found\n";
    } catch (const jwt::VerificationError& e) {
        // Handle the base verification error.
        //NOTE: There are other derived types of verification errors
        // which will be discussed in next topic.
        std::cerr << "Verification error exception found\n";
    } catch (...) {
        std::cerr << "Caught unknown exception\n";
    }

    return 0;
}
