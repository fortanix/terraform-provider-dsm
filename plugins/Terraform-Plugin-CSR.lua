--
--
-- IMPORT: Name "Terraform Plugin - CSR"
--
--
-- EXAMPLES:
--
-- ```
-- {
--   "subject_key": "my server key",
--   "cert_lifetime": 86400,
--   "subject_dn": { "CN": "localhost", "OU": "Testing" }
-- }
-- ```

function check(input)
   if type(input) ~= 'table' then
      return nil, 'invalid input'
   end
   if not input.subject_dn then
      return nil, 'must provide subject DN'
   end
   if not input.subject_key then
      return nil, 'must provide subject key'
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
   local something_key = assert(Sobject { kid = input.subject_key })

   local subject_dn = load_dn(input.subject_dn)
   -- log the DN here?

   local csr = Pkcs10Csr.new(something_key.name, subject_dn)

   local conv = Blob.from_bytes(csr:to_der())
   return {
      value = format_pem(conv:base64(), "CERTIFICATE REQUEST"),
      kid = something_key.kid,
      id = tostring(Time.now_insecure()[1])
   }
end