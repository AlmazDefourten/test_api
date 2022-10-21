package avl

import (
	"errors"
	"fmt"
)

var (
	ErrDuplicatedNode error = errors.New("bst: found duplicated value on tree")
	ErrNodeNotFound   error = errors.New("bst: node not found")
)

type User struct {
	id   int
	name string
	mail string
}

type node struct {
	height      int
	value       *User
	left, right *node
}

func (n *node) balanceFactor() int {
	if n == nil {
		return 0
	}

	return n.left.Height() - n.right.Height()
}

func (n *node) updateHeight() {
	max := func(a, b int) int {
		if a > b {
			return a
		}

		return b
	}
	n.height = max(n.left.Height(), n.right.Height()) + 1
}

func (n *node) Height() int {
	if n == nil {
		return 0
	}

	return n.height
}

func (n *node) Value() *User {
	return n.value
}

func (n *node) Left() *node {
	return n.left
}

func (n *node) Right() *node {
	return n.right
}

func newNode(val *User) *node {
	return &node{
		height: 1,
		value:  val,
		left:   nil,
		right:  nil,
	}
}

func insertNode(node *node, val *User) (*node, error) {
	// Создание нового узла
	if node == nil {
		return newNode(val), nil
	}

	// Если значение уже существует в дереве, возвращаем ошибку
	if node.value == val {
		return nil, ErrDuplicatedNode
	}

	// если значение больше текущего узла, вставляем его в правое поддерево
	if val.name > node.value.name {
		right, err := insertNode(node.right, val)

		if err != nil {
			return nil, err
		}

		node.right = right
	}

	// если значение меньше текущего узла, вставляем его в левое поддерево
	if val.name < node.value.name {
		left, err := insertNode(node.left, val)

		if err != nil {
			return nil, err
		}

		node.left = left
	}

	return rotateInsert(node, val), nil
}

func removeNode(node *node, val *User) (*node, error) {
	if node == nil {
		return nil, ErrNodeNotFound
	}

	if val.name > node.value.name {
		right, err := removeNode(node.right, val)
		if err != nil {
			return nil, err
		}

		node.right = right
	} else if val.name < node.value.name {
		left, err := removeNode(node.left, val)
		if err != nil {
			return nil, err
		}

		node.left = left
	} else {
		if node.left != nil && node.right != nil {
			// имеется оба поддерева

			// находим самый большой элемент в левом поддереве
			successor := greatest(node.left)
			value := successor.value

			// удаляем самый большой элемент в левом поддереве
			left, err := removeNode(node.left, value)
			if err != nil {
				return nil, err
			}
			node.left = left

			// присваиваем значение самому большому элементу в левом поддереве
			node.value = value
		} else if node.left != nil || node.right != nil {
			// имеется только одно поддерево
			// присваиваем значение самому большому элементу в левом поддереве
			if node.left != nil {
				node = node.left
			} else {
				node = node.right
			}
		} else if node.left == nil && node.right == nil {
			// удаляемый узел не имеет поддеревьев
			// присваиваем значение самому большому элементу в левом поддереве
			node = nil
		}
	}

	if node == nil {
		return nil, nil
	}

	return rotateDelete(node), nil
}

func findNode(node *node, val *User) *node {
	if node == nil {
		return nil
	}

	// если значение больше текущего узла, ищем его в правом поддереве
	if node.value == val {
		return node
	}

	// если значение меньше текущего узла, ищем его в левом поддереве
	if val.name > node.value.name {
		return findNode(node.right, val)
	}

	// если значение больше текущего узла, ищем его в правом поддереве
	if val.name < node.value.name {
		return findNode(node.left, val)
	}

	return nil
}

func rotateInsert(node *node, val *User) *node {
	// обновляем высоту узла
	node.updateHeight()

	// получаем разницу высот между правым и левым поддеревьями
	bFactor := node.balanceFactor()

	// если разница высот между правым и левым поддеревьями больше 1, то выполняем поворот
	if bFactor > 1 && val.name < node.left.value.name {
		return rightRotate(node)
	}

	// если разница высот между правым и левым поддеревьями меньше -1, то выполняем поворот
	if bFactor < -1 && val.name > node.right.value.name {
		return leftRotate(node)
	}

	// если разница высот между правым и левым поддеревьями больше 1, то выполняем поворот
	if bFactor > 1 && val.name > node.left.value.name {
		node.left = leftRotate(node.left)
		return rightRotate(node)
	}

	// если разница высот между правым и левым поддеревьями меньше -1, то выполняем поворот
	if bFactor < -1 && val.name < node.right.value.name {
		node.right = rightRotate(node.right)
		return leftRotate(node)
	}

	return node
}

func rotateDelete(node *node) *node {
	node.updateHeight()
	bFactor := node.balanceFactor()

	// линейно влево
	if bFactor > 1 && node.left.balanceFactor() >= 0 {
		return rightRotate(node)
	}

	// Меньше символа
	if bFactor > 1 && node.left.balanceFactor() < 0 {
		node.left = leftRotate(node.left)
		return rightRotate(node)
	}

	// линейно вправо
	if bFactor < -1 && node.right.balanceFactor() <= 0 {
		return leftRotate(node)
	}

	// больше символа
	if bFactor < -1 && node.right.balanceFactor() > 0 {
		node.right = rightRotate(node.right)
		return leftRotate(node)
	}

	return node
}

func rightRotate(x *node) *node {
	y := x.left
	t := y.right

	y.right = x
	x.left = t

	x.updateHeight()
	y.updateHeight()

	return y
}

func leftRotate(x *node) *node {
	y := x.right
	t := y.left

	y.left = x
	x.right = t

	x.updateHeight()
	y.updateHeight()

	return y
}

func greatest(node *node) *node {
	if node == nil {
		return nil
	}

	if node.right == nil {
		return node
	}

	return greatest(node.right)
}

func traverse(node *node) {
	// условие выхода
	if node == nil {
		return
	}

	fmt.Println(node.value)
	traverse(node.left)
	traverse(node.right)
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
