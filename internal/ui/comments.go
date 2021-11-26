package ui

import (
	"fmt"
	"github.com/CaninoDev/go-hackernews/internal/api"
	"code.rocketnine.space/tslocum/cview"
	"github.com/dustin/go-humanize"
)

type Comments struct {
	Title *cview.TextView
	Tree  *cview.TreeView
	Text  *cview.TextView
}

func NewComments() *Comments {
	return &Comments{
		cview.NewTextView(),
		cview.NewTreeView(),
		cview.NewTextView(),
	}
}
func (d *Display) generateCommentTree(post api.Post) *cview.TreeNode {
	treeRoot := cview.NewTreeNode(post.Title())

	var addChildNode func(api.Post) *cview.TreeNode

	addChildNode = func(childPost api.Post) *cview.TreeNode {
		childText := treeTextGenerator(childPost)
		childNode := cview.NewTreeNode(childText)
		childNode.SetReference(childPost)
		if len(childPost.Kids()) > 0 {
			subscription := d.DB.Subscribe(childPost.Kids())
			go func() {
				for grandChild := range subscription.Updates() {
					grandChildNode := addChildNode(grandChild)
					childNode.AddChild(grandChildNode)
				}
			}()
		}
		return childNode
	}

	renderTree := func() {
		subscription := d.DB.Subscribe(post.Kids())
		for child := range subscription.Updates() {
			commentNode := addChildNode(child)
			treeRoot.AddChild(commentNode)
			d.App.QueueUpdateDraw(func() {})
		}
	}

	if len(post.Kids()) > 0 {
		go renderTree()
	}

	return treeRoot

}

func treeTextGenerator(post api.Post) string {
	var postText string
	if len(post.Kids()) == 1 {
		postText = fmt.Sprintf("[-:-:-]%s[::d] (%s) %v child", post.By(), humanize.Time(post.Time()), len(post.Kids()))
	} else {
		postText = fmt.Sprintf("[-:-:-]%s[::d] (%s) %v children", post.By(), humanize.Time(post.Time()), len(post.Kids()))
	}
	return postText
}
