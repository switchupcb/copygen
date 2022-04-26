// Package matcher matches fields.
package matcher

import (
	"github.com/switchupcb/copygen/cli/models"
)

// Match matches the fields of a parsed generator.
func Match(gen *models.Generator) error {
	for _, function := range gen.Functions {
		for _, toType := range function.To {
			for _, fromType := range function.From {

				// top-level types can be pointed (i.e domain.Account).
				toFields := toType.Field.AllFields(nil, nil)
				fromFields := fromType.Field.AllFields(nil, nil)

				// each toField is compared to every fromField.
				for i := 0; i < len(toFields); i++ {
					for j := 0; j < len(fromFields); j++ {
						match(function, toFields[i], fromFields[j])
						if toFields[i].From != nil {
							break
						}
					}
				}
			}
		}
	}

	RemoveUnpointedFields(gen)
	return nil
}

// match determines which matcher to use for two fields, then matches them.
func match(function models.Function, toField *models.Field, fromField *models.Field) {
	if function.Options.Manual {
		switch {
		case toField.Options.Automatch || fromField.Options.Automatch:
			automatch(toField, fromField)

		case toField.Options.Tag != "":
			tagmatch(toField, fromField)

		default:
			mapmatch(toField, fromField)
		}
	} else {
		automatch(toField, fromField)
	}
}

// automatch automatically matches the fields of a fromType to a toType by name and definition.
// automatch is used when no `map` or `tag` options apply to a field.
func automatch(toField, fromField *models.Field) {
	if toField.Name == fromField.Name &&
		(toField.FullDefinition() == fromField.FullDefinition() || fromField.Options.Convert != "") {
		fromField.To = toField
		toField.From = fromField

		// prevent parallel matching.
		fromField.Fields = make([]*models.Field, 0)
		toField.Fields = make([]*models.Field, 0)
	}
}

// mapmatch manually maps a from-field to a to-field.
// mapmatch is used when a map option is specified.
func mapmatch(toField, fromField *models.Field) {
	if fromField.Options.Map != "" && toField.FullNameWithoutPointer("") == fromField.Options.Map {
		fromField.To = toField
		toField.From = fromField
	}
}

// tagmatch manually maps a from-field to a to-field using tags.
// tagmatch is used when a tag option is specified.
func tagmatch(toField, fromField *models.Field) {
	if toField.Options.Tag != "" && toField.Options.Tag == fromField.Options.Tag {
		fromField.To = toField
		toField.From = fromField
	}
}
