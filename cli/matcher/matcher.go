// Package matcher matches fields.
package matcher

import (
	"github.com/switchupcb/copygen/cli/models"
)

// FieldsMatcher represents a matcher of two fields.
type FieldsMatcher struct {
	// The fields that will be mapped from other fields TO these fields.
	toFields []*models.Field

	// The fields that will be mapped to other fields FROM these fields.
	fromFields []*models.Field
}

// fieldMatcher represents a matcher of two fields.
type fieldMatcher struct {
	toField   *models.Field
	fromField *models.Field
}

// Match matches the fields of a parsed generator.
func Match(gen *models.Generator) error {
	for _, function := range gen.Functions {
		for _, toType := range function.To {
			for _, fromType := range function.From {
				// The top-level types are not pointed to any fields (i.e domain.Account).
				fm := FieldsMatcher{toFields: toType.Field.AllFields(nil), fromFields: fromType.Field.AllFields(nil)}
				err := fm.Automatch()
				if err != nil {
					return err
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

// Automatch automatically matches the fields of a fromType to a toType by name.
// Automatch is used when no `map` options apply to a field.
func (fm *FieldsMatcher) Automatch() error {
	for i := 0; i < len(fm.toFields); i++ {
		// each toField is compared to every fromField.
		for j := 0; j < len(fm.fromFields); j++ {
			// therefore, don't compare pointed fields.
			if fm.toFields[i].From == nil && fm.fromFields[j].To == nil {
				// don't compare top-level fields that have subfields
				// this allows type such as `type T int` but not `type User struct` to be matched.
				if (fm.toFields[i].Parent != nil || len(fm.toFields[i].Fields) == 0) && (fm.fromFields[j].Parent != nil || len(fm.fromFields[j].Fields) == 0) {
					if fm.toFields[i].Name == fm.fromFields[j].Name && (fm.toFields[i].Definition == fm.toFields[i].Definition || fm.fromFields[j].Options.Convert != "") {
						fm.fromFields[j].To = fm.toFields[i]
						fm.toFields[i].From = fm.fromFields[j]
					}
				}
			}
		}
	}
	return nil
}
