package parser

import (
	"fmt"
	"go/ast"
	"go/token"
	"strings"

	"github.com/switchupcb/copygen/cli/models"
)

// Traverse parses the generator setup file's Abstract Syntax Tree into function and field data.
// Assigns values to the parser's convert options.
func (p *Parser) Traverse(gen *models.Generator) error {
	// traverse the AST to locate the `type Copygen interface` and assign options respectively.
	// store comments for use in analysis and to remove from the keep.
	for i := len(p.SetupFile.Decls) - 1; i > -1; i-- {
		switch x := p.SetupFile.Decls[i].(type) {
		case *ast.GenDecl:
			if it, ok := assertCopygenInterface(x); ok {
				// Keep all declaration objects in the setup file except for the `type Copygen interface`.
				p.Copygen = it
				p.SetupFile.Decls[i] = p.SetupFile.Decls[len(p.SetupFile.Decls)-1]
				p.SetupFile.Decls = p.SetupFile.Decls[:len(p.SetupFile.Decls)-1]

				// assign respective `type Copygen interface` comments.
				comments, err := p.assignOptions(x)
				if err != nil {
					return err
				}

				p.Comments = append(p.Comments, comments...)
			}

		// set convert option values.
		case *ast.FuncDecl:
			comments, err := p.assignConvertOptions(x)
			if err != nil {
				return err
			}

			p.Comments = append(p.Comments, comments...)
		}
	}

	// Analyze the `type Copygen Interface` for function and field data.
	if p.Copygen == nil {
		return fmt.Errorf("the \"type Copygen interface\" could not be found in the setup file")
	}

	var err error
	gen.Functions, err = p.parseFunctions(p.Copygen)

	if err != nil {
		return err
	}

	// Remove option-comments from the AST.
	astRemoveComments(p.SetupFile, p.Comments)

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

// assignOptions initializes function and field-specific options.
// Used in the context of the type Copygen interface.
func (p *Parser) assignOptions(x ast.Node) ([]*ast.Comment, error) {
	var (
		comments  []*ast.Comment
		assignerr error
	)

	ast.Inspect(x, func(node ast.Node) bool {
		xcg, ok := node.(*ast.CommentGroup)
		if !ok {
			return true
		}

		for i := 0; i < len(xcg.List); i++ {
			if xcg.List[i].Slash < x.Pos() {
				comments = append(comments, xcg.List[i])
				continue
			}
			// do not use the Doc above the node as an option
			text := xcg.List[i].Text
			splitcomments := strings.Fields(text[2:])

			// map[comment]map[optionname]map[]
			// determine if the comment is an option.
			if len(splitcomments) >= 1 {
				category := splitcomments[0]
				option := strings.Join(splitcomments[1:], " ")
				switch category {
				case categoryDeepCopy:
					opt, err := parseDeepcopy(option)
					if err != nil {
						assignerr = err
						return false
					}
					p.Options[text] = *opt
				case categoryDepth:
					opt, err := parseDepth(option)
					if err != nil {
						assignerr = err
						return false
					}
					p.Options[text] = *opt
				case categoryMap:
					opt, err := parseMap(option)
					if err != nil {
						assignerr = err
						return false
					}
					p.Options[text] = *opt
				default:
					p.Options[text] = Option{
						Category: categoryCustom,
						Regex:    nil,
						Value:    map[string]string{category: option},
					}
				}
			}
			// all type Copygen interface comments will be removed.
			comments = append(comments, xcg.List[i])
		}

		return true
	})

	return comments, assignerr
}

// assignConvertOptions initializes convert options.
// Used in the context of functions other than the type Copygen interface.
func (p *Parser) assignConvertOptions(x *ast.FuncDecl) ([]*ast.Comment, error) {
	var (
		comments  []*ast.Comment
		assignerr error
	)

	ast.Inspect(x, func(node ast.Node) bool {
		xcg, ok := node.(*ast.CommentGroup)
		if !ok {
			return true
		}

		for i := 0; i < len(xcg.List); i++ {
			text := xcg.List[i].Text
			splitcomments := strings.Fields(text[2:])

			// map[comment]map[optionname]map[]
			// determine if the comment is a convert option.
			if len(splitcomments) == 3 {
				category := splitcomments[0]
				option := strings.Join(splitcomments[1:], " ")
				if category == categoryConvert {
					opt, err := parseConvert(option, x.Name.Name)
					if err != nil {
						assignerr = err
						return false
					}
					p.Options[text] = *opt
					comments = append(comments, xcg.List[i])
				}
			}
		}

		return true
	})

	return comments, assignerr
}

// astRemoveComments removes comments from an *ast.File.
func astRemoveComments(file *ast.File, comments []*ast.Comment) {
	clength := len(comments)

	for i := 0; i < len(file.Comments); i++ {
		cg := file.Comments[i]
		// traverse through the comment group
		for j := 0; j < len(cg.List); j++ {
			// remove the comments by comparison
			for c := 0; c < clength; c++ {
				if cg.List[j] == comments[c] {
					// remove from the comment group top-down.
					if j-1 > -1 {
						cg.List[j].Text = cg.List[j-1].Text
						cg.List[j-1].Text = "  " // printer: "" and " " give an out of bounds error
					} else {
						cg.List[j].Text = "  " // printer: "" and " " give an out of bounds error
					}

					break
				}
			}
		}
	}
}
