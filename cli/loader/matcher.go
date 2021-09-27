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
			toField, fromField, err := matchFields(toFields[i], fromFields[j], toType, fromType)
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

// matchFields points respective fields to each other and a parent.
func matchFields(toField *models.Field, fromField *models.Field, toType *models.Type, fromType *models.Type) (*models.Field, *models.Field, error) {
	if toField.Name == fromField.Name && toField.Definition == fromField.Definition {
		toField.Parent = *toType
		fromField.Parent = *fromType
		toField.From = fromField
		fromField.To = toField
		return toField, fromField, nil
	} else {
		// reminder: AST search only find fields at the depth-level specified.
		// if a field has the same name, but wrong definition (i.e models.User vs. domain.User)
		// there is a chance for it contain a match at the next depth-level.
		//
		// when both fields have nested fields, there an be a direct match between any level.
		if len(toField.Fields) != 0 && len(fromField.Fields) != 0 {
			for _, nestedToField := range toField.Fields {
				for _, nestedFromField := range fromField.Fields {
					return matchFields(nestedToField, nestedFromField, toType, fromType)
				}
			}
		}

		// when a toField has fields but a fromField doesn't, there can be a direct match
		// from the fields of the toField to the fromField (see automatch example: User.UserID -> UserID).
		if len(toField.Fields) != 0 {
			for _, nestedToField := range toField.Fields {
				return matchFields(nestedToField, fromField, toType, fromType)
			}
		} else if len(fromField.Fields) != 0 {
			for _, nestedFromField := range fromField.Fields {
				return matchFields(toField, nestedFromField, toType, fromType)
			}
		}
	}
	return nil, nil, fmt.Errorf("The fields %v and %v could not be matched.", toField, fromField)
}
