package configuration

import (
	"errors"
	"time"
)

// NewConfiguration returns a new configuration
func NewConfiguration(settings ...Setting) Configuration {
	config := Configuration{
		macros: make(map[string]MacroTemplate),
	}
	for _, set := range settings {
		set(&config)
	}
	return config
}

// Configuration the configuration used when rendering a document
type Configuration struct {
	InputFile   string
	OutputFile  string
	LastUpdated time.Time
	EDATool     string
	OnePartRows bool
	SortLayer   bool
	macros      map[string]MacroTemplate
}

// Clone return a clone of the current configuration
func (c Configuration) Clone() Configuration {
	return Configuration{
		EDATool:     c.EDATool,
		InputFile:   c.InputFile,
		OutputFile:  c.OutputFile,
		LastUpdated: c.LastUpdated,
	}
}

// MacroTemplate finds and returns a user macro function by specified name.
func (c Configuration) MacroTemplate(name string) (MacroTemplate, error) {
	macro, ok := c.macros[name]
	if ok {
		return macro, nil
	}
	return nil, errors.New("unknown user macro: " + name)
}

const (
	// LastUpdatedFormat key to the time format for the `last updated` document attribute
	LastUpdatedFormat string = "2006-01-02 15:04:05 -0700"
)

// Setting a setting to customize the configuration used during parsing and rendering of a document
type Setting func(config *Configuration)

func WithOnePartRows(value bool) Setting {
	return func(config *Configuration) {
		config.OnePartRows = value
	}
}

func WithSortLayer(value bool) Setting {
	return func(config *Configuration) {
		config.SortLayer = value
	}
}

// WithLastUpdated function to set the `last updated` option in the renderer context (default is `time.Now()`)
func WithLastUpdated(value time.Time) Setting {
	return func(config *Configuration) {
		config.LastUpdated = value
	}
}

// WithEDATool function to set the `EDATool` setting in the config
func WithEDATool(tool string) Setting {
	return func(config *Configuration) {
		config.EDATool = tool
	}
}

// WithInputFile function to set the `filename` setting in the config
func WithInputFile(filename string) Setting {
	return func(config *Configuration) {
		config.InputFile = filename
	}
}

// WithOutputFile function to set the `filename` setting in the config
func WithOutputFile(filename string) Setting {
	return func(config *Configuration) {
		config.OutputFile = filename
	}
}

// WithMacroTemplate defines the given template to a user macro with the given name
func WithMacroTemplate(name string, t MacroTemplate) Setting {
	return func(config *Configuration) {
		config.macros[name] = t
	}
}
