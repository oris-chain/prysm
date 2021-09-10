package sync

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
	ethpb "github.com/prysmaticlabs/prysm/proto/prysm/v1alpha1"
	"google.golang.org/protobuf/proto"
)

// skipcq: SCC-U1000
func (s *Service) syncCommitteeMessageSubscriber(_ context.Context, msg proto.Message) error {
	m, ok := msg.(*ethpb.SyncCommitteeMessage)
	if !ok {
		return fmt.Errorf("message was not type *eth.SyncCommitteeMessage, type=%T", msg)
	}

	if m == nil {
		return errors.New("nil sync committee message")
	}

	return s.cfg.SyncCommsPool.SaveSyncCommitteeMessage(m)
}
