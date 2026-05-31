---@class PrivateCompleteClass
---@field public publicField number
---@field private privateField number
---@field public childObj PrivateCompleteClass
local PrivateCompleteClass = {}

---@return PrivateCompleteClass
function PrivateCompleteClass:new()
    return self
end

---@private
function PrivateCompleteClass:privateMethod()
end

function PrivateCompleteClass:publicMethod()
end

function PrivateCompleteClass:testSelf()
    self.childObj = PrivateCompleteClass:new()

end

local privateCompleteObj = PrivateCompleteClass:new()

