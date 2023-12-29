package robot

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/astaxie/beego/httplib"
)

type RobotTestReq struct {
	Canister string `json:"canister"`
	Owner    string `json:"owner"`
}

type RobotTestResp struct {
	Data   bool   `json:"data"`
	Msg    string `json:"msg"`
	Status int    `json:"status"`
}

type Roboter struct {
	host string
}

func NewRoboter(host string) *Roboter {
	return &Roboter{
		host: host,
	}
}

func (p *Roboter) IsRobot(req RobotTestReq) (bool, error) {
	str := fmt.Sprintf("%s/v1/api/robot/test", p.host)

	hreq := httplib.NewBeegoRequest(str, "POST")
	hreq.JSONBody(req)
	hreq.SetTimeout(30*time.Second, 30*time.Second)
	body, err := hreq.Bytes()
	if err != nil {
		return false, err
	}

	var resp RobotTestResp
	if err := json.Unmarshal(body, &resp); err != nil {
		return false, err
	}
	if resp.Status != 0 {
		return false, errors.New(resp.Msg)
	}
	return resp.Data, nil
}
