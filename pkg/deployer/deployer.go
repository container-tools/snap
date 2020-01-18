package deployer

type Deployer interface {
	Deploy(source, destination string) error
}
