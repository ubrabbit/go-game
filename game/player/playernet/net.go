package playernet

import (
	"errors"
	"fmt"
	"server/game/net"
	"server/game/player"
)

func SendToPlayerID(pid int, i interface{}) error {
	p := player.GetPlayer(pid)
	if p == nil {
		return errors.New(fmt.Sprintf("player %d is not online!", pid))
	}
	return net.SendToPlayer(p, i)
}

func SendToPlayer(p *player.Player, i interface{}) error {
	return net.SendToPlayer(p, i)
}
