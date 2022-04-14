package parser

import (
	"go/ast"
	"go/token"
)

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

const (
	newline        = 1
	carriagereturn = 2
)

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
					if j != 0 &&
						(fileCommentGroup.List[j-1].End()+newline == comment.Slash ||
							fileCommentGroup.List[j-1].End()+carriagereturn == comment.Slash) {
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
