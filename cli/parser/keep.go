package parser

import (
	"fmt"
	"go/ast"
	"go/token"
	"strings"

	"github.com/switchupcb/copygen/cli/parser/options"
)

const convertOptionSplitAmount = 3

// Keep removes ast.Nodes from an ast.File that will be kept in a generated output file.
func Keep(astFile *ast.File) error {
	var trash []*ast.Comment

	for i := len(astFile.Decls) - 1; i > -1; i-- {
		switch declaration := astFile.Decls[i].(type) {
		case *ast.GenDecl:

			// keep all declaration objects in the setup file except for the `type Copygen interface`.
			if _, ok := assertCopygenInterface(declaration); ok {

				// remove from the `type Copygen interface` (from the slice).
				astFile.Decls[i] = astFile.Decls[len(astFile.Decls)-1]
				astFile.Decls = astFile.Decls[:len(astFile.Decls)-1]

				// remove the `type Copygen interface` function ast.Comments.
				comments := getNodeComments(declaration)
				trash = append(trash, comments...)
				err := assignFieldOption(commentOptionMap, comments)
				if err != nil {
					return fmt.Errorf("%w", err)
				}
			}

		case *ast.FuncDecl:
			comments, options, err := assignConvertOptions(declaration)
			if err != nil {
				return fmt.Errorf("%w", err)
			}

			// remove convert option ast.Comments.
			trash = append(trash, comments...)
			convertOptions = append(convertOptions, options...)
		}
	}

	// Remove ast.Comments that will be parsed into options from the ast.File.
	astRemoveComments(astFile, trash)

	return nil
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

// astRemoveComments removes ast.Comments from an *ast.File.
func astRemoveComments(file *ast.File, comments []*ast.Comment) {

	// remove comments starting from the bottom of the file.
	for i := len(file.Comments) - 1; i > -1; i-- {
		fileCommentGroup := file.Comments[i]

		// remove comments from the bottom of each comment group.
		for j := len(fileCommentGroup.List) - 1; j > -1; j-- {
			fileComment := fileCommentGroup.List[j]

			for k := len(comments) - 1; k > -1; k-- {
				comment := comments[k]

				// remove the comment.
				if fileComment == comment {
					// reslice the commentGroup to remove the comment.
					fileCommentGroup.List = append(fileCommentGroup.List[:j], fileCommentGroup.List[j+1:]...)

					// prevent free-floating comments.
					if j != 0 && fileCommentGroup.List[j-1].End()+2 == comment.Slash {
						fileCommentGroup.List[j-1].Slash = comment.Slash
					}

					// prevent the comment from being compared again.
					comments[k] = comments[len(comments)-1]
					comments = comments[:len(comments)-1]
					break
				}
			}
		}
	}
}

// assignFieldOption parses a list of ast.Comments into options
// and places them in a map[text]Option.
func assignFieldOption(optionmap map[string]*options.Option, comments []*ast.Comment) error {
	for _, comment := range comments {
		text := comment.Text
		if optionmap[text] != nil {
			continue
		}

		splitcomments := strings.Fields(text[2:])
		if len(splitcomments) >= 1 {

			category := splitcomments[0]
			if category == options.CategoryConvert {
				continue
			}

			optiontext := strings.Join(splitcomments[1:], " ")
			option, err := options.NewFieldOption(category, optiontext)
			if err != nil {
				return fmt.Errorf("%w", err)
			}

			optionmap[text] = option
		}
	}

	return nil
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
