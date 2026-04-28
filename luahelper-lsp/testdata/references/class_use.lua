---@type RefWidgetModel
local widget = {}

widget.x = 1
widget:InitializeWidget()
widget.primaryOption = 3
widget.secondaryOption = 4
widget.compactMode = 5
widget.displayMode = 6

local aliasWidget = widget
aliasWidget:InitializeWidget()

local ServiceRegistry = {}

function ServiceRegistry.GetWidgetModel()
	return widget
end

ServiceRegistry.GetWidgetModel():InitializeWidget()

local returnedWidget = ServiceRegistry.GetWidgetModel()
returnedWidget:InitializeWidget()
