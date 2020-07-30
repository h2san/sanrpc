package selector

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

// go test -v -coverprofile=cover.out
// go tool cover -func=cover.out

type IPSelectorTestSuite struct {
	suite.Suite
	selector Selector
}

func (suite *IPSelectorTestSuite) SetupSuite() {
}

func (suite *IPSelectorTestSuite) SetupTest() {
	suite.selector = Get("ip")
}

func (suite *IPSelectorTestSuite) TestIpSelectSingleIp() {
	serviceName := "trpc.service.ip.1:1234"
	node, err := suite.selector.Select(serviceName)
	suite.T().Logf("Select return node:{%+v}, err:{%+v}", node, err)

	suite.NoError(err)
	suite.Equal(node.ServiceName, "trpc.service.ip.1:1234")
	suite.Equal(node.Address, "trpc.service.ip.1:1234")
}

func (suite *IPSelectorTestSuite) TestIpSelectMultiIp() {
	serviceName := "trpc.service.ip.1:1234,trpc.service.ip.2:1234"

	node, err := suite.selector.Select(serviceName)
	suite.T().Logf("Select return node:{%+v}, err:{%+v}", node, err)
	suite.NoError(err)
	suite.Equal(node.ServiceName, serviceName)

	node, err = suite.selector.Select(serviceName)
	suite.T().Logf("Select return node:{%+v}, err:{%+v}", node, err)
	suite.NoError(err)
	suite.Equal(node.ServiceName, serviceName)
}

func (suite *IPSelectorTestSuite) TestIpSelectEmpty() {
	serviceName := ""

	node, err := suite.selector.Select(serviceName)
	suite.T().Logf("Select return node:{%+v}, err:{%+v}", node, err)
	suite.Error(err)
	suite.Nil(node, "serviceName is empty")
}

func TestIpSelector(t *testing.T) {
	suite.Run(t, new(IPSelectorTestSuite))
}

func TestIpSelectorSelect(t *testing.T) {
	s := Get("ip")
	n, err := s.Select("trpc.service.ip.1:8888")
	assert.Nil(t, err)
	assert.Equal(t, n.Address, "trpc.service.ip.1:8888")
}

func TestIpSelectorReport(t *testing.T) {
	s := Get("ip")
	n, err := s.Select("trpc.service.ip.1:8888")
	assert.Nil(t, err)
	assert.Equal(t, n.Address, "trpc.service.ip.1:8888")
	assert.Nil(t, s.Report(n, 0, nil))
}
