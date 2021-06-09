--
--
-- IMPORT: Name "Terraform Plugin"
--
--
-- EXAMPLES:
--
--{
--     "operation":"create",
--     "name":"something"
--     "group_id":"some group id"
--}
--
--{
--     "operation":"drop",
--     "kid":"some key id"
--}

function get_date_from_unix(unix_time)
    local day_count, year, days, month = function(yr) return (yr % 4 == 0 and (yr % 100 ~= 0 or yr % 400 == 0)) and 366 or 365 end, 1970, math.ceil(unix_time/86400)
  
    while days >= day_count(year) do
      days = days - day_count(year) year = year + 1
    end
    local tab_overflow = function(seed, table) for i = 1, #table do if seed - table[i] <= 0 then return i, seed end seed = seed - table[i] end end
    month, days = tab_overflow(days, {31,(day_count(year) == 366 and 29 or 28),31,30,31,30,31,31,30,31,30,31})
    local hours, minutes, seconds = math.floor(unix_time / 3600 % 24), math.floor(unix_time / 60 % 60), math.floor(unix_time % 60)
    local period = hours > 12 and "pm" or "am"
    hours = hours > 12 and hours - 12 or hours == 0 and 12 or hours
    return string.format("%04d%02d%dT%02d%02d%02dZ", year, month, days, hours, minutes, seconds)
  end
  
  function generate_secret()
    local chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*()_+-={}|[]`~"
    local length = 32
    local randomString = ""
  
    charTable = {}
    for c in chars:gmatch"." do
      table.insert(charTable, c)
    end
  
    for i = 1, length do
      randomString = randomString .. charTable[math.random(1, #charTable)]
    end
    return randomString
  end
  
  function http_request(endpoint, authorization_header, port, path, request_body, method)
    local content_type = "application/json"
    if authorization_header ~= nil then
        local headers = { ["Content-Type"] = content_type, ["Authorization"] = authorization_header }
    else
        local headers = { ["Content-Type"] = content_type }
    end
    local request_url = ""
    if port ~= nil then
     	request_url = "https://" .. endpoint .. ":" .. port .. "/" .. path
    else
     	request_url = "https://" .. endpoint .. "/" .. path
    end
    
    if request_body ~= nil then
      response, err = request { method = method, url = request_url, headers = headers, body = json.encode(request_body) }
    else
      response, err = request { method = method, url = request_url, headers = headers }
    end
  
    return response, err
  end
  
function run(input)
	if input.operation == "create" then
	    local new_secret = generate_secret()
    	local new_sobject = assert(Sobject.import { name = input.name, group_id = input.group_id, obj_type = "SECRET", value = Blob.from_bytes(new_secret)})
        local resp_payload = {
            kid      = new_sobject.kid,
            name     = new_sobject.name,
            group_id = new_sobject.group_id
		}
		return resp_payload
    elseif input.operation == "drop" then
        local sobject, err = Sobject { name = input.kid }
        if sobject == nil or err ~= nil then
            err = "[DSM SDK]: No known security object named " .. input.kid .. "."
            return nil, err
        end
		assert(sobject:delete())
		local resp_payload = {}
        return resp_payload
    else
		err = "[DSM SDK]: No operation was specified"
		return nil, err
    end
end