# Pocketbase UI-Pluggable

This is a fork of [PocketBase](https://github.com/pocketbase/pocketbase). The original project can be found at https://github.com/pocketbase/pocketbase.

# UI Plugin Command

The `uiplugin` package adds a new `ui-plugin` command to PocketBase that allows you to easily create new UI plugins with a basic structure.

## Installation

Add the plug to your main.go file:

```go

// register the plug command
uiplugin.MustRegister(app, app.RootCmd, uiplugin.Config{
	Dir: "ui-plugins", // optional: defaults to "ui-plugins"
})
```

## Usage

Once registered, you can use the `plug` command to create new UI plugins:

```bash
# Create a new UI plugin named "my-plugin", or run `go run ./ plug my-plugin` with source code
pb plug my-plugin

# This will create the following structure:
# ui-plugins/my_plugin/
# ├── plugin.go      # Plugin registration file
# └── dist/          # Static files directory
#     └── index.html # Default HTML file
```



## Using the Plugin

To use the created plugin:

1. **Import the plugin** in your main.go file:
   ```go
   import _ "your-project/ui-plugins/my_plugin"
   ```

2. **Build your frontend assets** and place them in the `dist/` directory

3. **Restart your PocketBase application**

4. **Access your plugin** through the PocketBase admin UI sidebar.


## Plugin Structure

Each UI plugin consists of:

- **plugin.go**: Go file that registers the plugin with PocketBase
- **dist/**: Directory containing static frontend assets
  - **index.html**: Main HTML file for the plugin

The plugin will be accessible in the PocketBase admin UI sidebar with the specified icon and name. 

## How it works   

The implementation is based on the changes in commit [f0fb4d4](https://github.com/pocketbase/pocketbase/commit/f0fb4d463d214145ff9d3daa8c584c93a5a7f700), which adds the core UI plugin functionality to PocketBase.
