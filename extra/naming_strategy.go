package extra

import (
	"github.com/json-iterator/go"
	"unicode"
)

func SetNamingStrategy(translate func(string) string) {
	jsoniter.RegisterExtension(&namingStrategyExtension{jsoniter.DummyExtension{}, translate})
}

type namingStrategyExtension struct {
	jsoniter.DummyExtension
	translate func(string) string
}

func (extension *namingStrategyExtension) UpdateStructDescriptor(structDescriptor *jsoniter.StructDescriptor) {
	for _, binding := range structDescriptor.Fields {
		binding.ToNames = []string{extension.translate(binding.Field.Name)}
		binding.FromNames = []string{extension.translate(binding.Field.Name)}
	}
}

func LowerCaseWithUnderscores(name string) string {
	newName := []rune{}
	for i, c := range name {
		if i == 0 {
			newName = append(newName, unicode.ToLower(c))
		} else {
			if unicode.IsUpper(c) {
				newName = append(newName, '_')
				newName = append(newName, unicode.ToLower(c))
			} else {
				newName = append(newName, c)
			}
		}
	}
	return string(newName)
}
