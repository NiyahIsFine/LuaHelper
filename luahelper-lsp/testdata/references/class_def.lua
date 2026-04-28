---@class RefWidgetModel
---@field x integer
local RefWidgetController = {}

function RefWidgetController:InitializeWidget()
    self.primaryOption = 2
    self.secondaryOption = 1
    self.compactMode = 1
    self.displayMode = 2
end

return RefWidgetController
