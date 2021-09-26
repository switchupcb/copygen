package loader

import (
	"github.com/switchupcb/copygen/cli/models"
)

// DefineFieldsByFromType defines fields for a to-type and from-type based on the from-type.
// Used when the user specifies the fields to match in the loader.
func DefineFieldsByFromType(from *From) ([]models.Field, []models.Field) {
	var toFields, fromFields []models.Field
	for fieldName, field := range (*from).Fields {
		var fromField models.Field
		fromField.Name = fieldName
		fromField.Convert = field.Convert
		fromField.Options = models.FieldOptions{
			Custom: field.Options,
		}

		var toField models.Field
		toField.Name = field.To
		toField.Convert = field.Convert
		toField.Options = models.FieldOptions{
			Custom: field.Options,
		}

		// point the fields
		fromField.To = &toField
		fromFields = append(fromFields, fromField)
		toField.From = &fromField
		toFields = append(toFields, toField)
	}
	return toFields, fromFields
}

// Automatch uses an AST to automatically match the fields of a toType by name.
// Used when no field options are specified in the loader.
func Automatch(to *models.Type, from *models.Type) ([]models.Field, []models.Field, error) {
	var tempToFields, tempFromFields []models.Field
	tempToFields, err := ASTSearch(to.Options.Import, to.Package, to.Name)
	if err != nil {
		return nil, nil, err
	}

	tempFromFields, err = ASTSearch(from.Options.Import, from.Package, from.Name)
	if err != nil {
		return nil, nil, err
	}

	// point the fields
	var toFields, fromFields []models.Field
	for _, toField := range tempToFields {
		newToField := models.Field{
			Parent: *to,
			Name:   toField.Name,
		}
		for _, fromField := range tempFromFields {
			if toField.Name == fromField.Name && toField.Definition == fromField.Definition {
				newFromField := models.Field{
					Parent: *from,
					Name:   fromField.Name,
				}
				newFromField.To = &newToField
				fromFields = append(fromFields, newFromField)
				newToField.From = &newFromField
				toFields = append(toFields, newToField)
			}
		}
	}
	return toFields, fromFields, nil
}
