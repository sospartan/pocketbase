package uiplugin_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/sospartan/pocketbase/plugins/uiplugin"
	"github.com/sospartan/pocketbase/tests"
)

func TestUIPluginCreate(t *testing.T) {
	t.Parallel()

	scenarios := []struct {
		name           string
		expectedDir    string
		expectedFiles  []string
		expectedPlugin string
	}{
		{
			"my-plugin",
			"my_plugin",
			[]string{
				"plugin.go",
				"dist/index.html",
			},
			`// Package my_plugin handles the PocketBase UI plugin embedding.
package my_plugin

import (
	"embed"
	"io/fs"

	"github.com/sospartan/pocketbase/apis"
)

//go:embed all:dist
var distDir embed.FS

var PluginDirFS, _ = fs.Sub(distDir, "dist")

func init() {
	apis.RegisterUIPlugin(apis.UiPlugin{
		Name: "my-plugin",
		Base: "my_plugin",
		Icon: "ri-plug-line",
		FS:   PluginDirFS,
	})
}
`,
		},
		{
			"TestPlugin",
			"test_plugin",
			[]string{
				"plugin.go",
				"dist/index.html",
			},
			`// Package test_plugin handles the PocketBase UI plugin embedding.
package test_plugin

import (
	"embed"
	"io/fs"

	"github.com/sospartan/pocketbase/apis"
)

//go:embed all:dist
var distDir embed.FS

var PluginDirFS, _ = fs.Sub(distDir, "dist")

func init() {
	apis.RegisterUIPlugin(apis.UiPlugin{
		Name: "TestPlugin",
		Base: "test_plugin",
		Icon: "ri-plug-line",
		FS:   PluginDirFS,
	})
}
`,
		},
	}

	for _, s := range scenarios {
		t.Run(s.name, func(t *testing.T) {
			app, _ := tests.NewTestApp()
			defer app.Cleanup()

			// create a temporary directory for testing
			tempDir := filepath.Join(app.DataDir(), "_test_ui_plugins")
			defer os.RemoveAll(tempDir)

			// create plugin instance and test the handler directly
			p := &uiplugin.Plugin{App: app, Config: uiplugin.Config{Dir: tempDir}}

			// create the plugin
			pluginDir, err := p.UIPluginCreateHandler(s.name, false)
			if err != nil {
				t.Fatalf("Failed to create UI plugin, got: %v", err)
			}

			expectedDir := filepath.Join(tempDir, s.expectedDir)
			if pluginDir != expectedDir {
				t.Fatalf("Expected plugin directory %q, got %q", expectedDir, pluginDir)
			}

			// check if the directory exists
			if _, err := os.Stat(expectedDir); os.IsNotExist(err) {
				t.Fatalf("Expected plugin directory to exist: %v", err)
			}

			// check if all expected files exist
			for _, file := range s.expectedFiles {
				filePath := filepath.Join(expectedDir, file)
				if _, err := os.Stat(filePath); os.IsNotExist(err) {
					t.Fatalf("Expected file to exist: %v", err)
				}
			}

			// check plugin.go content
			pluginGoPath := filepath.Join(expectedDir, "plugin.go")
			content, err := os.ReadFile(pluginGoPath)
			if err != nil {
				t.Fatalf("Failed to read plugin.go: %v", err)
			}

			contentStr := strings.TrimSpace(string(content))
			expectedPlugin := strings.TrimSpace(s.expectedPlugin)
			if contentStr != expectedPlugin {
				t.Fatalf("Expected plugin.go content:\n%v\ngot:\n%v", expectedPlugin, contentStr)
			}

			// check index.html content
			indexHTMLPath := filepath.Join(expectedDir, "dist", "index.html")
			htmlContent, err := os.ReadFile(indexHTMLPath)
			if err != nil {
				t.Fatalf("Failed to read index.html: %v", err)
			}

			htmlContentStr := string(htmlContent)
			if !strings.Contains(htmlContentStr, s.name) {
				t.Fatalf("Expected index.html to contain plugin name %q", s.name)
			}

			if !strings.Contains(htmlContentStr, "Welcome to") {
				t.Fatalf("Expected index.html to contain welcome message")
			}
		})
	}
}

func TestUIPluginCreateWithEmptyName(t *testing.T) {
	t.Parallel()

	app, _ := tests.NewTestApp()
	defer app.Cleanup()

	// create plugin instance and test the handler directly
	p := &uiplugin.Plugin{App: app, Config: uiplugin.Config{}}

	// try to create plugin with empty name
	_, err := p.UIPluginCreateHandler("", false)
	if err == nil {
		t.Fatalf("Expected error when creating plugin with empty name")
	}

	if !strings.Contains(err.Error(), "missing plugin name") {
		t.Fatalf("Expected error message to contain 'missing plugin name', got: %v", err)
	}
}
