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
				// The top-level types are pointed if applicable (i.e domain.Account).
				toFields := toType.Field.AllFields(nil)
				fromFields := fromType.Field.AllFields(nil)

				// each toField is compared to every fromField.
				for i := 0; i < len(toFields); i++ {
					for j := 0; j < len(fromFields); j++ {
						// therefore, don't compare pointed fields.
						if toFields[i].From != nil || fromFields[j].To != nil {
							continue
						}

						// don't compare top-level fields that have subfields.
						// allows type such as `type T int` but not `type User struct` to be matched.
						if (toFields[i].Parent != nil || len(toFields[i].Fields) == 0) && (fromFields[j].Parent != nil || len(fromFields[j].Fields) == 0) {
							match(function, toFields[i], fromFields[j])
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
		case toField.Options.Automatch:
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
		(toField.Definition == fromField.Definition || fromField.Options.Convert != "") {
		fromField.To = toField
		toField.From = fromField
	}
}

// mapmatch manually maps a from-field to a to-field.
// mapmatch is used when a map option is specified.
func mapmatch(toField, fromField *models.Field) {
	if fromField.Options.Map != "" && toField.FullNameWithoutContainer("") == fromField.Options.Map {
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
