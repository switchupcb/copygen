package parser

import (
	"go/ast"
	"go/token"
	"strings"

	"github.com/switchupcb/copygen/cli/parser/options"
)

// Removed contains removed ast.Nodes from a setup file's Abstract Syntax Tree.
type Removed struct {
	// Copygen represents ast.Node of the `type Copygen Interface`.
	Copygen *ast.InterfaceType

	// Comments represents ast.Comments parsed in the `type Copygen Interface`.
	Comments []*ast.Comment

	// ConvertOptions represents convert function options.
	ConvertOptions []*options.Option
}

const convertOptionSplitAmount = 3

// Keep removes ast.Nodes from an ast.File that won't be kept in a generated output file.
// modifies the given ast.File (what was kept) and returns what was removed.
func Keep(astFile *ast.File) (Removed, error) {
	var trash Removed

	for i := len(astFile.Decls) - 1; i > -1; i-- {
		switch declaration := astFile.Decls[i].(type) {
		case *ast.GenDecl:

			// keep all declaration objects in the setup file except for the `type Copygen interface`.
			if it, ok := assertCopygenInterface(declaration); ok {

				// remove from the `type Copygen interface` (from the slice).
				trash.Copygen = it
				astFile.Decls[i] = astFile.Decls[len(astFile.Decls)-1]
				astFile.Decls = astFile.Decls[:len(astFile.Decls)-1]

				// remove the `type Copygen interface` function comment.
				trash.Comments = append(trash.Comments, getNodeComments(declaration)...)
			}

		case *ast.FuncDecl:
			comments, options, err := assignConvertOptions(declaration)
			if err != nil {
				return trash, err
			}

			// remove convert option ast.Comments
			trash.Comments = append(trash.Comments, comments...)
			trash.ConvertOptions = append(trash.ConvertOptions, options...)
		}
	}

	// Remove ast.Comments that will be parsed into options from the ast.File.
	astRemoveComments(astFile, trash.Comments)

	return trash, nil
}

// assertCopygenInterface determines if an ast.GenDecl is a Copygen Interface by type assertion.
func assertCopygenInterface(x *ast.GenDecl) (*ast.InterfaceType, bool) {
	if x.Tok == token.TYPE {
		for _, spec := range x.Specs {
			if ts, ok := spec.(*ast.TypeSpec); ok {
				if it, ok := ts.Type.(*ast.InterfaceType); ok && ts.Name.Name == "Copygen" {
					return it, true
				}
			}
		}
	}

	return nil, false
}

// getNodeComments returns all of the ast.Comments in a given node.
func getNodeComments(x ast.Node) []*ast.Comment {
	var optionComments []*ast.Comment

	ast.Inspect(x, func(node ast.Node) bool {
		commentGroup, ok := node.(*ast.CommentGroup)
		if !ok {
			return true
		}

		for i := 0; i < len(commentGroup.List); i++ {
			optionComments = append(optionComments, commentGroup.List[i])
		}

		return true
	})

	return optionComments
}

// assignConvertOptions initializes convert options.
// Used in the context of functions other than the type Copygen interface.
func assignConvertOptions(x *ast.FuncDecl) ([]*ast.Comment, []*options.Option, error) {
	var (
		convertComments []*ast.Comment
		convertOptions  []*options.Option
		assignErr       error
	)

	ast.Inspect(x, func(node ast.Node) bool {
		commentGroup, ok := node.(*ast.CommentGroup)
		if !ok {
			return true
		}

		for _, comment := range commentGroup.List {
			text := comment.Text
			splitcomments := strings.Fields(text[2:])

			// determine if the comment is a convert option.
			if len(splitcomments) == convertOptionSplitAmount {
				category := splitcomments[0]
				value := strings.Join(splitcomments[1:], " ")
				if category == options.CategoryConvert {
					option, err := options.ParseConvert(value, x.Name.Name)
					if err != nil {
						assignErr = err
						return false
					}

					convertOptions = append(convertOptions, option)
					convertComments = append(convertComments, comment)
				}
			}
		}

		return true
	})

	return convertComments, convertOptions, assignErr
}

// astRemoveComments removes ast.Comments from an *ast.File.
func astRemoveComments(file *ast.File, comments []*ast.Comment) {
	cLength := len(comments)

	for i := 0; i < len(file.Comments); i++ {
		commentGroup := file.Comments[i]
		// traverse through the comment group
		for j := 0; j < len(commentGroup.List); j++ {
			// remove the comments by comparison
			for c := 0; c < cLength; c++ {
				if commentGroup.List[j] == comments[c] {
					// remove from the comment group top-down.
					if j > 0 {
						commentGroup.List[j].Text = commentGroup.List[j-1].Text
						commentGroup.List[j-1].Text = "  " // printer: "" and " " give an out of bounds error
					} else {
						commentGroup.List[j].Text = "  " // printer: "" and " " give an out of bounds error
					}

					break
				}
			}
		}
	}
}
