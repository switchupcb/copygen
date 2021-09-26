package loader

import (
	"github.com/switchupcb/copygen/cli/models"
)

// DefineFieldsByFrom defines fields for a to-type and from-type based on a YML from.
// Used when the user specifies the fields to match in the loader.
func DefineFieldsByFrom(from *From, toType *models.Type, fromType *models.Type) ([]models.Field, []models.Field) {
	var toFields, fromFields []models.Field
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
		toFields = append(toFields, toField)
		fromFields = append(fromFields, fromField)
	}
	return toFields, fromFields
}

// Automatch uses an AST to automatically match the fields of a toType by name.
// Used when no field options are specified in the loader.
func Automatch(toType *models.Type, fromType *models.Type) ([]models.Field, []models.Field, error) {
	toFields, err := ASTSearch(toType.Options.Import, toType.Package, toType.Name)
	if err != nil {
		return nil, nil, err
	}

	fromFields, err := ASTSearch(fromType.Options.Import, fromType.Package, fromType.Name)
	if err != nil {
		return nil, nil, err
	}

	// ASTSearch finds all the fields for each type.
	// The name and definition can be used to determine a to-from field pair
	// The field pair requires parents and to point to each other.
	var newToFields, newFromFields []models.Field
	for i := 0; i < len(toFields); i++ {
		for j := 0; j < len(fromFields); j++ {
			if toFields[i].Name == fromFields[j].Name && toFields[i].Definition == fromFields[j].Definition {
				toFields[i].Parent = *toType
				fromFields[j].Parent = *fromType
				toFields[i].From = &fromFields[j]
				fromFields[j].To = &toFields[i]
				newToFields = append(newToFields, toFields[i])
				newFromFields = append(newFromFields, fromFields[j])
			}
		}
	}
	return newToFields, newFromFields, nil
}
