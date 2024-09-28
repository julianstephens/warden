package main

type InitCmd struct {
	Path string `arg:"" name:"path" help:"Path to new store." type:"path"`
}

func (i *InitCmd) Run(ctx *Globals) error {

	return nil
}
