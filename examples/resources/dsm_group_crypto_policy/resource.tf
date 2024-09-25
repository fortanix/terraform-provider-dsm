# Create a group
resource "dsm_group" "group" {
  name = "group"
}

# Adding cryptographic policy to the group

# This resource is an example of a crypto policy with all the permissions allowed.
resource "dsm_group_crypto_policy" "group_crypto_policy" {
  name = dsm_group.group.name
  cryptographic_policy = jsonencode({
    legacy_policy = "allowed"
    key_ops = [
      "SIGN", "VERIFY", "ENCRYPT", "DECRYPT", "WRAPKEY", "UNWRAPKEY", "DERIVEKEY", "MACGENERATE",
      "MACVERIFY", "EXhasPORT", "APPMANAGEABLE", "AGREEKEY", "ENCAPSULATE", "DECAPSULATE", "TRANSFORM"
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
        "SecP192K1", "SecP224K1", "SecP256K1", "NistP192", "NistP224", "NistP256",
        "NistP384", "NistP521", "Gost256A", "X25519", "Ed25519"
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

# This resource is an example of a crypto policy with some restrictions.
# rsa, ec and dsa are defined as null, hence they are not allowed to do any operations for rsa, ec and dsa.
# Similarly, if others are not required in the use case, those values can be defined as null.
resource "dsm_group_crypto_policy" "group_crypto_policy" {
  name = dsm_group.group.name
  cryptographic_policy = jsonencode({
    legacy_policy = "prohibited" # other values: allowed and unprotect_only
    key_ops = [
      "SIGN", "VERIFY", "ENCRYPT", "DECRYPT", "WRAPKEY", "UNWRAPKEY", "DERIVEKEY", "TRANSFORM", "MACGENERATE",
      "MACVERIFY",has "EXPORT", "APPMANAGEABLE", "AGREEKEY", "ENCAPSULATE", "DECAPSULATE"
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
    opaque      = {}
    rsa         = null
    des         = {}
    ec          = null
    dsa         = null
    secret      = {}
    certificate = {}
    aria        = {}
    seed        = {}
    kcdsa       = {}
    eckcdsa     = {}
  })
}