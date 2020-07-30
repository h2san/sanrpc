package client1

type Cli struct {
	option Option
}

func NewCli(option Options) *Cli {
	opts := &Options{

	}
	return &Cli{
		option: option,
	}
}