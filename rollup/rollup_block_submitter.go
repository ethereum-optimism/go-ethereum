package rollup

type RollupBlockSubmitter interface {
  submit(block *RollupBlock) error
}

type BlockSubmitter struct {}
func NewBlockSubmitter() *BlockSubmitter {
  return &BlockSubmitter{}
}
func (d *BlockSubmitter) submit(block *RollupBlock) error {
  return nil
}