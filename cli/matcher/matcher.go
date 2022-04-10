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

// match determines which matcher to use for two fields,
// then matches them.
func match(function models.Function, toField *models.Field, fromField *models.Field) {
	if function.Options.Manual {
		manualmatch(toField, fromField)
	} else {
		automatch(toField, fromField)
	}
}

// automatch automatically matches the fields of a fromType to a toType by name.
// automatch is used when no `map` options apply to a field.
func automatch(toField, fromField *models.Field) {
	if toField.Name == fromField.Name &&
		(toField.Definition == fromField.Definition || fromField.Options.Convert != "") {
		fromField.To = toField
		toField.From = fromField
	}
}

// manualmatch uses a manual matcher to map a from-field to a to-field.
// manualmatch is used when a map option is specified.
func manualmatch(toField, fromField *models.Field) {
	if fromField.Options.Map != "" && toField.FullName("") == fromField.Options.Map {
		fromField.To = toField
		toField.From = fromField
	}
}
