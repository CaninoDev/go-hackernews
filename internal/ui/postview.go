package ui

import (
	"code.rocketnine.space/tslocum/cview"
	"fmt"
	"github.com/CaninoDev/go-hackernews/internal/store"
	"github.com/dustin/go-humanize"
)

type Post struct {
	app          *App
	commentsTree *cview.TreeView
	debugBar     *cview.TextView
}

func NewPost(app *App) *Post {
	treeRoot := cview.NewTreeNode("")
	tree := cview.NewTreeView()
	tree.SetRoot(treeRoot)
	return &Post{
		app:          app,
		commentsTree: tree,
		debugBar:     NewDebugBar(),
	}
}

func (p *Post) SetPost(post store.Item) {
	p.commentsTree.SetRoot(cview.NewTreeNode(""))
	p.generateCommentTree(post)
}

func (p *Post) generateCommentTree(post store.Item) {
	if len(post.Kids()) > 0 {
		go p.renderTree(post)
	}
}

func (p *Post) addCommentNode(post store.Item) *cview.TreeNode {
	nodeText := treeTextGenerator(post)
	node := cview.NewTreeNode(nodeText)
	node.SetReference(post)
	if len(post.Kids()) > 0 {
		subscription := p.app.store.Subscribe(post.Kids())
		go func() {
			for grandChild := range subscription.Updates() {
				grandChildNode := p.addCommentNode(grandChild)
				node.AddChild(grandChildNode)
			}
		}()
	}
	return node
}

func (p *Post) renderTree(post store.Item) {
	treeRoot := cview.NewTreeNode(post.Title())

	subscription := p.app.store.Subscribe(post.Kids())
	for child := range subscription.Updates() {
		commentNode := p.addCommentNode(child)
		treeRoot.AddChild(commentNode)
		p.app.ui.QueueUpdateDraw(func() {})
	}
	p.commentsTree.SetRoot(treeRoot)
}

func treeTextGenerator(post store.Item) string {
	var postText string
	if len(post.Kids()) == 1 {
		postText = fmt.Sprintf("[-:-:-]%s[::d] (%s) %v child", post.By(), humanize.Time(post.Time()), len(post.Kids()))
	} else {
		postText = fmt.Sprintf("[-:-:-]%s[::d] (%s) %v children", post.By(), humanize.Time(post.Time()), len(post.Kids()))
	}
	return postText
}
