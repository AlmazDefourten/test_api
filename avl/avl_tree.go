package avl

type binarySearchTree struct {
  root *node
}

func New() *binarySearchTree {
  return &binarySearchTree{}
}

func (tree *binarySearchTree) Insert(val string) error {
  // always start insert from the root
  root, err := insertNode(tree.root, val)

  if err != nil {
    return err
  }

  tree.root = root
  return nil
}

func (tree *binarySearchTree) Remove(val string) error {
  root, err := removeNode(tree.root, val)

  if err != nil {
    return err
  }

  tree.root = root
  return nil
}

func (tree *binarySearchTree) Find(val string) *node {
  // as always, search from the root
  return findNode(tree.root, val)
}

func (tree *binarySearchTree) Traverse() {
  // traverse from the root
  traverse(tree.root)
}
