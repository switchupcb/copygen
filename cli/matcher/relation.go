package matcher

import (
	"github.com/switchupcb/copygen/cli/models"
)

// GetRelatedFields returns solely related fields in a list of fields.
func GetRelatedFields(fields []*models.Field) []*models.Field {
	for i := len(fields) - 1; i > -1; i-- {
		if len(fields[i].Fields) != 0 {
			fields[i].Fields = GetRelatedFields(fields[i].Fields)
		} else if fields[i].To == nil && fields[i].From == nil {
			fields = fields[:len(fields)-1]
		}
	}
	return fields
}

// isFieldRelated determines whether there is a relationship between two fields.
func isFieldRelated(toField *models.Field, fromField *models.Field) bool {
	if (*toField).From == fromField && (*fromField).To == toField {
		return true
	} else if (*toField).From == fromField {
		return true
	} else if (*fromField).To == toField {
		return true
	} else {
		if len(toField.Fields) != 0 && len(fromField.Fields) != 0 {
			for i := 0; i < len(toField.Fields); i++ {
				for j := 0; j < len(fromField.Fields); i++ {
					return isFieldRelated(toField.Fields[i], fromField.Fields[j])
				}
			}
		} else if len(toField.Fields) != 0 {
			for i := 0; i < len(toField.Fields); i++ {
				return isFieldRelated(toField.Fields[i], fromField)
			}
		} else if len(fromField.Fields) != 0 {
			for i := 0; i < len(fromField.Fields); i++ {
				return isFieldRelated(toField, fromField.Fields[i])
			}
		} else {
			return false
		}
	}
	return false
}
