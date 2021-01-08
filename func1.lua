local function max(...)
    local params = {...}
    local var, idx
    for i = 1,#params do
        if var == nil  or params[i] > var then
            print(i, var)
            var, idx = params[i], i
        end
    end
end


local function assert(v)
    if not v then failed() end
end


local v1 = max(3, 9, 7, 128, 35)
assert(v1 == 128)
local v2, i2 = max(3, 9, 7, 128, 35)
assert(v2 == 128 and i2 == 4)
local v3, i3 = max(max(3, 9, 7, 128, 35))
assert(v3 == 128 and i3 == 1)
local t = {max(3, 9, 7, 128, 35)}
assert(t[1] == 128 and t[2] == 4)
