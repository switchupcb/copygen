// Package matcher matches fields.
package matcher

import (
	"fmt"

	"github.com/switchupcb/copygen/cli/models"
)

// FieldsMatcher represents a matcher of two fields.
type FieldsMatcher struct {
	// The fields that will be mapped from other fields TO these fields.
	toFields []*models.Field

	// The fields that will be mapped to other fields FROM these fields.
	fromFields []*models.Field
}

// Match matches the fields of a parsed generator.
func Match(gen *models.Generator) error {
	for _, function := range gen.Functions {
		for _, toType := range function.To {
			for _, fromType := range function.From {
				// The main types are not pointed to any fields (i.e domain.Account).
				fm := FieldsMatcher{toType.Field.Fields, fromType.Field.Fields}
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
			fromType.Field.Fields = RelatedFields(fromType.Field.Fields)
		}
		for _, toType := range function.To {
			toType.Field.Fields = RelatedFields(toType.Field.Fields)
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
				fm := fieldMatcher{
					toField:   fm.toFields[i],
					fromField: fm.fromFields[j],
				}
				err := fm.matchFields()
				if err != nil {
					continue
				}
			}
		}
	}
	return nil
}

// fieldMatcher represents a matcher of two fields.
type fieldMatcher struct {
	toField   *models.Field
	fromField *models.Field
}

// matchFields points respective fields to each other.
func (fm fieldMatcher) matchFields() error {
	if fm.toField.Name == fm.fromField.Name && fm.toField.Definition == fm.fromField.Definition {
		fm.fromField.To = fm.toField
		fm.toField.From = fm.fromField
		return nil
	}

	// reminder: AST search only find fields at the depth-level specified.
	// if a field has the same name, but wrong definition (i.e models.User vs. domain.User)
	// there is a chance for it contain a match at the next depth-level.
	//
	// when both fields have nested fields, there an be a direct match between any level.
	if len(fm.toField.Fields) != 0 && len(fm.fromField.Fields) != 0 {
		for i := 0; i < len(fm.toField.Fields); i++ {
			fm.toField = fm.toField.Fields[i]
			for j := 0; j < len(fm.fromField.Fields); j++ {
				fm.fromField = fm.fromField.Fields[j]
				return fm.matchFields()
			}
		}
	}

	// when a toField has fields but a fromField doesn't, there can be a direct match
	// from the fields of the toField to the fromField (see automatch example: User.UserID -> UserID).
	if len(fm.toField.Fields) != 0 {
		for i := 0; i < len(fm.toField.Fields); i++ {
			fm.toField = fm.toField.Fields[i]
			return fm.matchFields()
		}
	} else if len(fm.fromField.Fields) != 0 {
		for i := 0; i < len(fm.fromField.Fields); i++ {
			fm.fromField = fm.fromField.Fields[i]
			return fm.matchFields()
		}
	}
	return fmt.Errorf("The fields %v and %v could not be matched.", fm.toField, fm.fromField)
}

// TODO: MANUAL MATCH
// // Map represents a manual match between fields.
// func Map(from *From, toType *models.Type, fromType *models.Type) ([]*models.Field, []*models.Field) {
// 	var toFields, fromFields []*models.Field
// 	for fieldname, field := range (*from).Fields {
// 		toField := models.Field{
// 			Parent:  *toType,
// 			Name:    field.To,
// 			Convert: field.Convert,
// 			Options: models.FieldOptions{
// 				Deepcopy: field.Deepcopy,
// 				Custom:   field.Options,
// 			},
// 		}

// 		fromField := models.Field{
// 			Parent:  *fromType,
// 			Name:    fieldname,
// 			Convert: field.Convert,
// 			Options: models.FieldOptions{
// 				Deepcopy: field.Deepcopy,
// 				Custom:   field.Options,
// 			},
// 		}

// 		// point the fields
// 		toField.From = &fromField
// 		fromField.To = &toField

// 		// keep track of the pointer for field.To and field.From comparison if required
// 		toFields = append(toFields, &toField)
// 		fromFields = append(fromFields, &fromField)
// 	}
// 	return toFields, fromFields
// }
