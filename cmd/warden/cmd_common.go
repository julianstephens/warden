package main

type CommonFlags struct {
	Store     string `short:"s" xor:"storefile" required:"" type:"path" help:"Path to your store"`
	StoreFile string `short:"f" xor:"store" required:"" type:"path" help:"Path to your store definition file"`
}
