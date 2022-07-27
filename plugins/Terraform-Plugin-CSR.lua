--
--
-- IMPORT: Name "Terraform Plugin - CSR"
--
--
-- EXAMPLES:
--
-- ```
-- {
--     "kid":       dsmsigner.kid,
--     "hash_alg":  "SHA256",
--     "data":      base64.StdEncoding.EncodeToString(digest),
-- }
-- ```

function check(input)
   if type(input) ~= 'table' then
      return nil, 'invalid input'
   end
   if not input.kid then
      return nil, 'must provide the key id'
   end
   if not input.hash_alg then
      return nil, 'must provide the has algorithm'
   end
   if not input.data then
      return nil, 'must provide the digest to sign'
   end
end

function format_pem(b64, type)
   local wrap_at = 64
   local len = string.len(b64)
   local pem = ""

   pem = pem .. "-----BEGIN " .. type .. "-----\n"

   for i = 1, len, wrap_at do
      local stop = i + wrap_at - 1
      pem = pem .. string.sub(b64, i, stop) .. "\n"
   end

   pem = pem .. "-----END " .. type .. "-----\n"

   return pem
end

function load_dn(dn)
   local name = X509Name.new()

   for k,v in pairs(dn)
   do
      name:set(Oid.from_str(k), v, 'utf8')
   end

   return name
end

function run(input)
   local signing_key = assert(Sobject { kid = input.kid })

   --local subject_dn = load_dn(input.subject_dn)
   -- log the DN here?

   --local csr = Pkcs10Csr.new(something_key.name, subject_dn)

   --local conv = Blob.from_bytes(csr:to_der())
   local signature = assert(signing_key:sign { hash = input.data, hash_alg = input.hash_alg })
   return signature
end