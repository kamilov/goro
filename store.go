package goro

import (
	"fmt"
	"math"
	"regexp"
	"strings"
)

type (
	store struct {
		root *node
		size int
	}

	node struct {
		static            bool
		path              string
		order             int
		minOrder          int
		handlers          []Handler
		children          []*node
		attributeChildren []*node
		attributeIndex    int
		attributeNames    []string
		requirement       *regexp.Regexp
	}
)

func newStore() *store {
	return &store{
		root: &node{
			static:            true,
			children:          make([]*node, 256),
			attributeChildren: make([]*node, 0),
			attributeIndex:    -1,
			attributeNames:    []string{},
		},
	}
}

func (s *store) String() string {
	return s.root.print(0)
}

func (s *store) add(path string, handlers []Handler) int {
	s.size++
	return s.root.add(path, handlers, s.size)
}

func (s *store) get(path string, values []string) (handlers []Handler, names []string) {
	handlers, names, _ = s.root.get(path, values)
	return
}

func (n *node) add(path string, handlers []Handler, order int) int {
	matched := 0

	for ; matched < len(path) && matched < len(n.path); matched++ {
		if path[matched] != n.path[matched] {
			break
		}
	}

	if matched == len(n.path) {
		if matched == len(path) {
			if n.handlers == nil {
				n.handlers = handlers
				n.order = order
			}
			return n.attributeIndex + 1
		}

		newPath := path[matched:]

		if child := n.children[newPath[0]]; child != nil {
			if index := child.add(newPath, handlers, order); index >= 0 {
				return index
			}
		}

		for _, child := range n.attributeChildren {
			if index := child.add(newPath, handlers, order); index >= 0 {
				return index
			}
		}

		return n.addAttribute(newPath, handlers, order)
	}

	if matched == 0 || !n.static {
		return -1
	}

	child := &node{
		static:            true,
		path:              n.path[matched:],
		order:             n.order,
		minOrder:          n.minOrder,
		handlers:          n.handlers,
		children:          n.children,
		attributeChildren: n.attributeChildren,
		attributeIndex:    n.attributeIndex,
		attributeNames:    n.attributeNames,
	}

	n.path = path[0:matched]
	n.handlers = nil
	n.children = make([]*node, 256)
	n.attributeChildren = make([]*node, 0)
	n.children[child.path[0]] = child

	return n.add(path, handlers, order)
}

func (n *node) addAttribute(path string, handlers []Handler, order int) int {
	openIndex, closeIndex, parentNode := -1, -1, n

	for i := 0; i < len(path); i++ {
		if openIndex < 0 && path[i] == '<' {
			openIndex = i
		}

		if openIndex >= 0 && path[i] == '>' {
			closeIndex = i
			break
		}
	}

	if openIndex > 0 && closeIndex > 0 || closeIndex == -1 {
		child := &node{
			static:            true,
			path:              path,
			minOrder:          order,
			children:          make([]*node, 256),
			attributeChildren: make([]*node, 0),
			attributeIndex:    n.attributeIndex,
			attributeNames:    n.attributeNames,
		}

		n.children[path[0]] = child

		if closeIndex == -1 {
			child.handlers = handlers
			child.order = order

			return child.attributeIndex + 1
		}

		child.path = path[:openIndex]
		parentNode = child
	}

	child := &node{
		static:            false,
		path:              path[openIndex : closeIndex+1],
		minOrder:          order,
		children:          make([]*node, 256),
		attributeChildren: make([]*node, 0),
		attributeIndex:    parentNode.attributeIndex,
		attributeNames:    parentNode.attributeNames,
	}

	requirement := ""
	attributeName := path[openIndex+1 : closeIndex]

	for i := openIndex + 1; i < closeIndex; i++ {
		if path[i] == ':' {
			attributeName = path[openIndex+1 : i]
			requirement = path[i+1 : closeIndex]
		}
	}

	if requirement != "" {
		child.requirement = regexp.MustCompile("^" + requirement)
	}

	attributeNames := make([]string, len(parentNode.attributeNames)+1)

	copy(attributeNames, n.attributeNames)

	attributeNames[len(parentNode.attributeNames)] = attributeName

	child.attributeIndex = len(attributeNames) - 1
	child.attributeNames = attributeNames

	parentNode.attributeChildren = append(parentNode.attributeChildren, child)

	if closeIndex == len(path)-1 {
		child.handlers = handlers
		child.order = order

		return child.attributeIndex + 1
	}

	return child.addAttribute(path[closeIndex+1:], handlers, order)
}

func (n *node) get(path string, values []string) (handlers []Handler, names []string, order int) {
	order = math.MaxInt32
	currentNode := n

find:
	if currentNode.static {
		npl := len(currentNode.path)

		if npl > len(path) {
			return
		}

		for i := npl - 1; i >= 0; i-- {
			if currentNode.path[i] != path[i] {
				return
			}
		}

		path = path[npl:]
	} else if currentNode.requirement != nil {
		if currentNode.requirement.String() == "^.*" {
			values[currentNode.attributeIndex] = path
			path = ""
		} else if match := currentNode.requirement.FindStringIndex(path); match != nil {
			values[currentNode.attributeIndex] = path[0:match[1]]
			path = path[match[1]:]
		} else {
			return
		}
	} else {
		i, kl := 0, len(path)

		for ; i < kl; i++ {
			if path[i] == '/' {
				values[currentNode.attributeIndex] = path[0:i]
				path = path[i:]
				break
			}
		}

		if i == kl {
			values[currentNode.attributeIndex] = path
			path = ""
		}
	}

	if len(path) > 0 {
		if child := currentNode.children[path[0]]; child != nil {
			if len(currentNode.children) == 0 {
				currentNode = child
				goto find
			}

			handlers, names, order = child.get(path, values)
		}
	} else if currentNode.handlers != nil {
		handlers, names, order = currentNode.handlers, currentNode.attributeNames, currentNode.order
	}

	trustedValues := values
	allocated := false

	for _, child := range currentNode.attributeChildren {
		if child.minOrder >= order {
			continue
		}

		if handlers != nil && !allocated {
			trustedValues = make([]string, len(values))
			allocated = true
		}

		if childHandlers, childNames, childOrder := child.get(path, trustedValues); childHandlers != nil && childOrder < order {
			if allocated {
				for i := child.attributeIndex; i < len(childNames); i++ {
					values[i] = trustedValues[i]
				}
			}

			handlers, names, order = childHandlers, childNames, childOrder
		}
	}

	return
}

func (n *node) print(level int) string {
	str := fmt.Sprintf(
		"%v{path: %v, order: %v, minOrder: %v, countHandlers: %v, attributeIndex: %v, attributeNames: %v, requirement: %v}\n",
		strings.Repeat(" ", level<<2),
		n.path,
		n.order,
		n.minOrder,
		len(n.handlers),
		n.attributeIndex,
		n.attributeNames,
		n.requirement,
	)

	for _, child := range n.children {
		if child != nil {
			str += child.print(level + 1)
		}
	}

	for _, child := range n.attributeChildren {
		str += child.print(level + 1)
	}

	return str
}
