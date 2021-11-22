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

						// don't compare top-level fields that have subfields
						// this allows type such as `type T int` but not `type User struct` to be matched.
						if (toFields[i].Parent != nil || len(toFields[i].Fields) == 0) && (fromFields[j].Parent != nil || len(fromFields[j].Fields) == 0) {
							if function.Options.Manual {
								if fromFields[j].Options.Map != "" {
									manualmatch(toFields[i], fromFields[j])
								}
							} else {
								automatch(toFields[i], fromFields[j])
							}
						}
					}
				}
			}
		}
	}

	// don't return unpointed fields.
	for _, function := range gen.Functions {
		for _, fromType := range function.From {
			fromType.Field.Fields = RelatedFields(fromType.Field.Fields, nil)
		}

		for _, toType := range function.To {
			toType.Field.Fields = RelatedFields(toType.Field.Fields, nil)
		}
	}

	return nil
}

// automatch automatically matches the fields of a fromType to a toType by name.
// automatch is used when no `map` options apply to a field.
func automatch(toField, fromField *models.Field) {
	if toField.Name == fromField.Name {
		if (toField.Definition == fromField.Definition && toField.Import == fromField.Import) || fromField.Options.Convert != "" || toField.OrigDefinition == fromField.OrigDefinition {
			fromField.To = toField
			toField.From = fromField
		}
	}
}

// manualmatch uses a manual matcher to map a from-field to a to-field.
// manualmatch is used when a map option is specified.
func manualmatch(toField, fromField *models.Field) {
	if toField.FullName("") == fromField.Options.Map {
		fromField.To = toField
		toField.From = fromField
	}
}
