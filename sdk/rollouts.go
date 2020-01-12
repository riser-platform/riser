package sdk

import (
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/riser-platform/riser-server/api/v1/model"
)

var trafficRuleExp = regexp.MustCompile("[0-9]+:[0-9]+")

type RolloutsClient interface {
	Save(deploymentName, stageName string, trafficRule ...string) error
}

type rolloutsClient struct {
	client *Client
}

func (c *rolloutsClient) Save(deploymentName, stageName string, trafficRule ...string) error {
	rolloutRequest := model.RolloutRequest{}
	for _, rule := range trafficRule {
		if !trafficRuleExp.MatchString(rule) {
			return errors.New("Rules must be in the format of \"(rev):(percentage)\" e.g. \"1:100\" routes 100% of traffic to rev 1")
		}
		ruleSplit := strings.Split(rule, ":")
		rolloutRequest.Traffic = append(rolloutRequest.Traffic,
			model.TrafficRule{
				RiserGeneration: mustParseInt(ruleSplit[0]),
				Percent:         int(mustParseInt(ruleSplit[1])),
			})
	}
	request, err := c.client.NewRequest(http.MethodPut, fmt.Sprintf("/api/v1/rollout/%s/%s", deploymentName, stageName), rolloutRequest)
	if err != nil {
		return err
	}

	_, err = c.client.Do(request, nil)
	return err
}

// mustParseInt panics which should never happen - validate input before using!
func mustParseInt(s string) int64 {
	i, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		panic(err)
	}

	return i
}
