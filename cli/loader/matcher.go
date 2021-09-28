package loader

import (
	"fmt"

	"github.com/switchupcb/copygen/cli/models"
)

// DefineFieldsByFrom defines fields for a to-type and from-type based on a YML from.
// Used when the user specifies the fields to match in the loader.
func DefineFieldsByFrom(from *From, toType *models.Type, fromType *models.Type) ([]*models.Field, []*models.Field) {
	var toFields, fromFields []*models.Field
	for fieldname, field := range (*from).Fields {
		toField := models.Field{
			Parent:  *toType,
			Name:    field.To,
			Convert: field.Convert,
			Options: models.FieldOptions{
				Deepcopy: field.Deepcopy,
				Custom:   field.Options,
			},
		}

		fromField := models.Field{
			Parent:  *fromType,
			Name:    fieldname,
			Convert: field.Convert,
			Options: models.FieldOptions{
				Deepcopy: field.Deepcopy,
				Custom:   field.Options,
			},
		}

		// point the fields
		toField.From = &fromField
		fromField.To = &toField

		// keep track of the pointer for field.To and field.From comparison if required
		toFields = append(toFields, &toField)
		fromFields = append(fromFields, &fromField)
	}
	return toFields, fromFields
}

// Automatch uses an AST to automatically match the fields of a toType by name.
// Used when no field options are specified in the loader.
func (a *AST) Automatch(toType *models.Type, fromType *models.Type) ([]*models.Field, []*models.Field, error) {
	a.MaxDepth = toType.Options.Depth
	a.Depth = 0
	toFields, err := a.Search(toType.Options.Import, toType.Package, toType.Name)
	if err != nil {
		return nil, nil, err
	}

	a.MaxDepth = fromType.Options.Depth
	a.Depth = 0
	fromFields, err := a.Search(fromType.Options.Import, fromType.Package, fromType.Name)
	if err != nil {
		return nil, nil, err
	}

	// The AST Search finds all the fields for each type at the specified depth level.
	// The name and definition can be used to determine a to-from field pair.
	// The field pair still needs to be pointed to its parent and to its respective field.
	var newToFields, newFromFields []*models.Field
	for i := 0; i < len(toFields); i++ {
		for j := 0; j < len(fromFields); j++ {
			fm := fieldMatcher{
				toField:   toFields[i],
				fromField: fromFields[j],
				toType:    toType,
				fromType:  fromType,
			}
			toField, fromField, err := fm.matchFields()
			if err != nil {
				continue
			}
			// keep track of the pointer for field.To and field.From comparison if required
			newToFields = append(newToFields, toField)
			newFromFields = append(newFromFields, fromField)
		}
	}
	return newToFields, newFromFields, nil
}

// fieldMatcher represets a matcher of two fields.
type fieldMatcher struct {
	toField   *models.Field
	fromField *models.Field
	toType    *models.Type
	fromType  *models.Type
}

// matchFields points respective fields (or their child fields) to each other and a parent.
func (fm fieldMatcher) matchFields() (*models.Field, *models.Field, error) {
	if fm.toField.Name == fm.fromField.Name && fm.toField.Definition == fm.fromField.Definition {
		fm.toField.Parent = *(fm.toType)
		fm.fromField.Parent = *(fm.fromType)
		fm.toField.From = fm.fromField
		fm.fromField.To = fm.toField
		return fm.toField, fm.fromField, nil
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
	return nil, nil, fmt.Errorf("The fields %v and %v with parents %v and %v could not be matched.", fm.toField, fm.fromField, fm.toType, fm.fromType)
}
