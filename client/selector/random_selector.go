package selector

import (
	"errors"
	"github.com/hillguo/sanrpc/client/node"

	"github.com/valyala/fastrand"
)


func init(){
	Register("random",&RandomSelector{})
}
// randomSelector selects randomly.
type RandomSelector struct {

}

func (s *RandomSelector) Select(list []*node.Node) (*node.Node,error) {
	if len(list) == 0 {
		return nil,errors.New("node list empty")
	}
	i := fastrand.Uint32n(uint32(len(list)))
	return list[i],nil
}


