package apis

import (
	"fmt"
	"io/fs"
	"net/http"

	"github.com/sospartan/pocketbase/core"
	"github.com/sospartan/pocketbase/tools/router"
)

// UiPlugin defines a UI plugin that can be registered and displayed in the admin UI sidebar.
type UiPlugin struct {
	// The display name of the plugin
	Name string `json:"name"`
	// The base path or URL of the plugin, eg. "static-plugins/plugin1"
	Base string `json:"base"`
	// The icon class name to display in the sidebar (eg. "ri-plug-line")
	Icon string `json:"icon"`
	// The filesystem to serve the plugin from
	FS fs.FS `json:"-"`
	// The flag to ignore route setting
	IgnoreRoute bool `json:"_"`
}

var uiPlugins = []UiPlugin{}

// RegisterUIPlugin registers one or more UI plugins to be displayed in the admin UI sidebar.
// Each plugin must have a unique name, base path/URL, and filesystem to serve its content from.
// The icon field is optional and can be any valid icon class name (eg. "ri-plug-line").
//
// Example:
//
//	apis.RegisterUIPlugin(apis.UiPlugin{
//		Name: "My Plugin",
//		Base: "my-plugin",
//		Icon: "ri-plug-line",
//		FS:   pluginFS,
//	})

func RegisterUIPlugin(plugins ...UiPlugin) {
	for _, p := range plugins {
		if p.Name == "" {
			panic("ui-plugins: name cannot be empty")
		}
		if p.Base == "" {
			panic("ui-plugins: base path/url cannot be empty")
		}
		if p.FS == nil {
			panic("ui-plugins: fs cannot be nil")
		}
	}
	uiPlugins = append(uiPlugins, plugins...)
}

// bindUIPluginsApi registers a route handler that returns a list of registered UI plugins.
// The plugins list is returned as JSON in the format:
// {
//   "plugins": [
//     {
//       "name": "Plugin Name",
//       "base": "plugin-path",
//       "icon": "ri-icon-name"
//     },
//     ...
//   ]
// }

func bindUIPluginsApi(_ core.App, rg *router.RouterGroup[*core.RequestEvent]) {
	subGroup := rg.Group("/ui-plugins")
	subGroup.GET("", func(e *core.RequestEvent) error {
		return e.JSON(http.StatusOK, map[string]any{
			"plugins": uiPlugins,
		})
	})
}

func bindUIPluginServeRoute(g *router.RouterGroup[*core.RequestEvent]) {
	for _, p := range uiPlugins {
		if !p.IgnoreRoute {
			g.GET(fmt.Sprintf("/%s/{path...}", p.Base), Static(p.FS, false))
		}
	}

}
