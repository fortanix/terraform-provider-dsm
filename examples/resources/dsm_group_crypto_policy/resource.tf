// Create a group
resource "dsm_group" "group" {
  name = "group"
}

// Adding cryptographic policy to the group
resource "dsm_group_crypto_policy" "group_crypto_policy" {
  name = dsm_group.group.name
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
      "MACGENERATE",
      "MACVERIFY",
      "EXPORT",
      "APPMANAGEABLE",
      "AGREEKEY",
      "ENCAPSULATE",
      "DECAPSULATE",
      "TRANSFORM"
    ]
    aes = {
      key_sizes = [128, 192, 256]
    }
    des3 = {
      key_sizes = [112, 168]
    }
    hmac = {
      minimum_key_length = 112
    }
    opaque = {}
    rsa = {
      encryption_policy  = []
      signature_policy   = []
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
    dsa         = {}
    secret      = {}
    certificate = {}
    aria        = {}
    seed        = {}
    kcdsa       = {}
    eckcdsa     = {}
    bip32       = {}
    lms         = {}
    mlkem_beta  = {}
    bls         = {}
  })
}

