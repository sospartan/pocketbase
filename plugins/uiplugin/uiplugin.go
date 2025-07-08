// Package uiplugin adds a new "ui-plugin" command support to a PocketBase instance.
//
// It provides functionality to create new UI plugins with a basic structure.
//
// Example usage:
//
//	uiplugin.MustRegister(app, app.RootCmd, uiplugin.Config{})
//
//	Note: This plugin creates UI plugins that can be registered and displayed
//	in the admin UI sidebar.
package uiplugin

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tools/inflector"
	"github.com/pocketbase/pocketbase/tools/osutils"
	"github.com/spf13/cobra"
)

// Config defines the config options of the uiplugin plugin.
type Config struct {
	// Dir specifies the directory where UI plugins will be created.
	//
	// If not set it fallbacks to a relative "ui-plugins" directory.
	Dir string
}

// MustRegister registers the uiplugin plugin to the provided app instance
// and panic if it fails.
//
// Example usage:
//
//	uiplugin.MustRegister(app, app.RootCmd, uiplugin.Config{})
func MustRegister(app core.App, rootCmd *cobra.Command, config Config) {
	if err := Register(app, rootCmd, config); err != nil {
		panic(err)
	}
}

// Register registers the uiplugin plugin to the provided app instance.
func Register(app core.App, rootCmd *cobra.Command, config Config) error {
	p := &plugin{App: app, Config: config}

	if p.Config.Dir == "" {
		p.Config.Dir = "ui-plugins"
	}

	// attach the ui-plugin command
	if rootCmd != nil {
		rootCmd.AddCommand(p.createCommand())
	}

	return nil
}

type Plugin struct {
	App    core.App
	Config Config
}

type plugin = Plugin

func (p *plugin) createCommand() *cobra.Command {
	const cmdDesc = `Creates a new UI plugin with the specified name.

The command will create a new directory under the ui-plugins directory
with the following structure:
  ui-plugins/{name}/
  ├── plugin.go      # Plugin registration file
  └── dist/          # Static files directory
      └── index.html # Default HTML file

Example:
  pb plug my-plugin
`

	command := &cobra.Command{
		Use:          "plug",
		Short:        "Creates a new UI plugin",
		Long:         cmdDesc,
		ValidArgs:    []string{"name"},
		Args:         cobra.ExactArgs(1),
		SilenceUsage: true,
		RunE: func(command *cobra.Command, args []string) error {
			name := args[0]
			if _, err := p.uiPluginCreateHandler(name, true); err != nil {
				return err
			}
			return nil
		},
	}

	return command
}

func (p *plugin) uiPluginCreateHandler(name string, interactive bool) (string, error) {
	if name == "" {
		return "", errors.New("missing plugin name")
	}

	// normalize the name
	normalizedName := inflector.Snakecase(name)

	// create the plugin directory
	pluginDir := filepath.Join(p.Config.Dir, normalizedName)

	if interactive {
		confirm := osutils.YesNoPrompt(fmt.Sprintf("Do you really want to create UI plugin %q in %q?", name, pluginDir), false)
		if !confirm {
			fmt.Println("The command has been cancelled")
			return "", nil
		}
	}

	// ensure that the ui-plugins dir exists
	if err := os.MkdirAll(p.Config.Dir, os.ModePerm); err != nil {
		return "", fmt.Errorf("failed to create ui-plugins directory: %w", err)
	}

	// create the plugin directory
	if err := os.MkdirAll(pluginDir, os.ModePerm); err != nil {
		return "", fmt.Errorf("failed to create plugin directory: %w", err)
	}

	// create the dist directory
	distDir := filepath.Join(pluginDir, "dist")
	if err := os.MkdirAll(distDir, os.ModePerm); err != nil {
		return "", fmt.Errorf("failed to create dist directory: %w", err)
	}

	// generate plugin.go file
	pluginGoContent, err := p.generatePluginGo(name, normalizedName)
	if err != nil {
		return "", fmt.Errorf("failed to generate plugin.go template: %w", err)
	}

	pluginGoPath := filepath.Join(pluginDir, "plugin.go")
	if err := os.WriteFile(pluginGoPath, []byte(pluginGoContent), 0644); err != nil {
		return "", fmt.Errorf("failed to save plugin.go file: %w", err)
	}

	// generate index.html file
	indexHTMLContent, err := p.generateIndexHTML(name)
	if err != nil {
		return "", fmt.Errorf("failed to generate index.html template: %w", err)
	}

	indexHTMLPath := filepath.Join(distDir, "index.html")
	if err := os.WriteFile(indexHTMLPath, []byte(indexHTMLContent), 0644); err != nil {
		return "", fmt.Errorf("failed to save index.html file: %w", err)
	}

	if interactive {
		fmt.Printf("Successfully created UI plugin %q in %q\n", name, pluginDir)
		fmt.Printf("Plugin files:\n")
		fmt.Printf("  - %s\n", pluginGoPath)
		fmt.Printf("  - %s\n", indexHTMLPath)
		fmt.Printf("\nTo use this plugin, you need to:\n")
		fmt.Printf("1. Import the plugin in your main.go file\n")
		fmt.Printf("2. Build the dist directory with your frontend assets\n")
		fmt.Printf("3. Restart your PocketBase application\n")
	}

	return pluginDir, nil
}

// UIPluginCreateHandler is a public method for testing purposes
func (p *Plugin) UIPluginCreateHandler(name string, interactive bool) (string, error) {
	return p.uiPluginCreateHandler(name, interactive)
}

func (p *plugin) generatePluginGo(name, normalizedName string) (string, error) {
	const template = `// Package %s handles the PocketBase UI plugin embedding.
package %s

import (
	"embed"
	"io/fs"

	"github.com/pocketbase/pocketbase/apis"
)

//go:embed all:dist
var distDir embed.FS

var PluginDirFS, _ = fs.Sub(distDir, "dist")

func init() {
	apis.RegisterUIPlugin(apis.UiPlugin{
		Name: "%s",
		Base: "%s",
		Icon: "ri-plug-line",
		FS:   PluginDirFS,
	})
}
`

	return fmt.Sprintf(template, normalizedName, normalizedName, name, normalizedName), nil
}

func (p *plugin) generateIndexHTML(name string) (string, error) {
	const template = `<!DOCTYPE html>
<html>
<head>
    <title>%s</title>
    <style>
        body {
            font-family: system-ui, -apple-system, sans-serif;
            max-width: 800px;
            margin: 0 auto;
            padding: 20px;
            line-height: 1.6;
        }
        h1 {
            color: #333;
            margin-bottom: 30px;
        }
        p {
            color: #666;
        }
        .container {
            background: #f8f9fa;
            border-radius: 8px;
            padding: 20px;
            margin: 20px 0;
        }
        .code {
            background: #e9ecef;
            padding: 10px;
            border-radius: 4px;
            font-family: 'Courier New', monospace;
            font-size: 14px;
        }
    </style>
</head>
<body>
    <h1>Welcome to %s</h1>
    <p>This is a UI plugin page for PocketBase.</p>
    
    <div class="container">
        <h2>Getting Started</h2>
        <p>You can customize this page with your own content and styling.</p>
        <p>This plugin is registered with the base path: <span class="code">%s</span></p>
    </div>
    
    <div class="container">
        <h2>Development</h2>
        <p>To develop this plugin:</p>
        <ul>
            <li>Replace the content in the <span class="code">dist/</span> directory with your frontend assets</li>
            <li>Build your frontend application and output to the <span class="code">dist/</span> directory</li>
            <li>Restart your PocketBase application to see changes</li>
        </ul>
    </div>
</body>
</html>`

	normalizedName := inflector.Snakecase(name)
	return fmt.Sprintf(template, name, name, normalizedName), nil
}
