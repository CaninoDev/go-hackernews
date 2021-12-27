package ui

import (
	"bytes"
	"fmt"
	"log"
	"net/url"
	"os/exec"
	"sync"

	"code.rocketnine.space/tslocum/cview"
	"github.com/CaninoDev/go-hackernews/internal/store"
	"github.com/dustin/go-humanize"
)

type Post struct {
	app          *App
	ref          store.Item
	content      *cview.TextView
	commentsTree *cview.TreeView
	debugBar     *cview.TextView
	sync.RWMutex
}

func NewPostView(app *App) *Post {
	tree := cview.NewTreeView()
	treeRoot := cview.NewTreeNode("")
	tree.SetRoot(treeRoot)

	contentView := cview.NewTextView()
	contentView.SetDynamicColors(true)
	contentView.SetPadding(1, 1, 0, 0)
	contentView.SetBorder(true)
	contentView.SetBorderColorFocused(Orange)

	return &Post{
		app:          app,
		ref:          store.Item{},
		content:      contentView,
		commentsTree: tree,
		debugBar:     NewDebugBar(),
	}
}

func (p *Post) SetPost(post store.Item) {
	p.Lock()
	defer p.Unlock()

	p.ref = post
	go p.renderContent()
	go p.renderCommentTree()
	p.app.ui.SetFocus(p.content)
}

func (p *Post) renderCommentTree() {
	p.commentsTree.SetRoot(cview.NewTreeNode(""))

	post := p.ref

	if len(post.Kids()) > 0 {
		treeRoot := cview.NewTreeNode(post.Title())
		kidsIDs := post.Kids()
		if len(kidsIDs) > 0 {

			buildTree := func() {
				for _, kidID := range kidsIDs {
					kid, err := p.app.store.Item(kidID)
					if err != nil {
						p.debugBar.SetText(fmt.Sprintf("%v", err))
					}
					childNode := p.addCommentNode(kid)
					treeRoot.AddChild(childNode)
				}
			}

			go buildTree()

			p.commentsTree.SetRoot(treeRoot)

		}
	}
	p.generateCommentTree(p.ref)
}

func (p *Post) renderContent() {
	postURL := p.ref.URL()

	if len(postURL) < 1 {
		p.content.SetTitle(p.ref.Title())
		p.content.Write([]byte(p.ref.Text()))
	} else {
		if parsedURL, err := url.Parse(postURL); err != nil {
			p.debugBar.SetText(fmt.Sprint(err))
		} else {
			parsedURLString := parsedURL.String()
			p.terminalRender(parsedURLString)
			// 	content, err := readability.FromURL(parsedURLString, 7*time.Second)
			// 	if err != nil {
			// 		p.debugBar.SetText(fmt.Sprint(err))
			// 	}
			// 	p.content.Write([]byte(content.Content))
			// }
		}
	}
}

func (p *Post) terminalRender(url string) {
	cmd := exec.Command("w3m", "ansi-color", "-dump", url)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		log.Fatalf("cmd.Run() failed with %s\n", err)
	}
	p.content.Write(stdout.Bytes())
}

func (p *Post) generateCommentTree(post store.Item) {

}

func (p *Post) addCommentNode(post store.Item) *cview.TreeNode {
	nodeText := treeTextGenerator(post)
	node := cview.NewTreeNode(nodeText)
	node.SetReference(post)

	ids := post.Kids()
	if len(ids) > 0 {
		//subscription := p.app.store.Subscribe(postView.Kids())
		//go func() {
		//	for grandChild := range subscription.Updates() {
		//		grandChildNode := p.addCommentNode(grandChild)
		//		node.AddChild(grandChildNode)
		//	}
		//}()
		buildTree := func() {
			for _, id := range ids {
				child, err := p.app.store.Item(id)
				if err != nil {
					p.debugBar.SetText(fmt.Sprintf("%v", err))
				}
				childNode := p.addCommentNode(child)
				node.AddChild(childNode)
			}
		}

		p.app.ui.QueueUpdateDraw(buildTree)
	}
	return node
}

func (p *Post) renderTree(post store.Item) {

	//subscription := p.app.store.Subscribe(postView.Kids())
	//for child := range subscription.Updates() {
	//
	//	commentNode := p.addCommentNode(child)
	//	treeRoot.AddChild(commentNode)
	//	p.app.ui.QueueUpdateDraw(func() {})
	//}

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
