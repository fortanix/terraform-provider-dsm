variable "acct_id" {
  type    = string
  default = "XXXXXXXX-XXXX-XXXX-XXXX-XXXXXXXXXXXX"
}
// Add cryptographic policy to a Fortanix DSM account
resource "dsm_acc_crypto_policy" "name" {
  acct_id = var.acct_id
  cryptographic_policy = jsonencode({
    legacy_policy = "allowed"
    key_ops = [
      "SIGN",
      "VERIFY",
      "ENCRYPT",
      "DECRYPT",
      "WRAPKEY",
      "UNWRAPKEY",
      "DERIVEKEY",
      "TRANSFORM",
      "MACGENERATE",
      "MACVERIFY",
      "EXPORT",
      "APPMANAGEABLE",
      "AGREEKEY",
      "ENCAPSULATE",
      "DECAPSULATE"
    ]
    aes = {
      key_sizes = [
        128,
        192,
        256
      ]
    }
    des3 = {
      key_sizes = [
        112,
        168
      ]
    }
    hmac = {
      minimum_key_length = 112
    }
    opaque = {}
    rsa = {
      encryption_policy = [
        {
          padding = {
            PKCS1_V15 = {}
          }
        },
        {
          padding = {
            RAW_DECRYPT = {}
          }
        },
        {
          padding = {
            OAEP = {
              mgf = {
                mgf1 = {}
              }
            }
          }
        }
      ]
      signature_policy = [
        {
          padding = {
            PKCS1_V15 = {}
          }
        },
        {
          padding = {
            PSS = {
              mgf = {
                mgf1 = {}
              }
            }
          }
        }
      ]
      minimum_key_length = 1024
    }
    des = {}
    ec = {
      elliptic_curves = [
        "SecP192K1",
        "SecP224K1",
        "SecP256K1",
        "NistP192",
        "NistP224",
        "NistP256",
        "NistP384",
        "NistP521",
        "Gost256A",
        "X25519",
        "Ed25519"
      ]
    }
    dsa = {}
    secret = {}
    certificate = {}
    aria = {}
    seed = {}
    kcdsa = {}
    eckcdsa = {}
  })
}