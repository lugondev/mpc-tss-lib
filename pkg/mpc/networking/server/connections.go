package server

import (
	"context"
	"fmt"

	"github.com/lugondev/mpc-tss-lib/pkg/mpc/networking/common"
)

func (s *Server) makeNamedConnections(IDs []string) {
	s.connections = make(map[string]*common.Conn)
	for _, id := range IDs {
		if s.Tss == nil || id != s.Tss.MyID() {
			s.logger.Debug().Str("id", id).Msg("Add connection")
			s.connections[id] = nil
		}
	}
}

func (s *Server) waitForConnections(ctx context.Context) error {
	s.logger.Debug().Msg("Wait for connections")

	var connChange context.Context // context will be canceled when connection change
	connChange, s.connChangeNotify = context.WithCancel(context.Background())

	for {
		if !s.isAllConnected() {
			select {
			case <-connChange.Done(): // wait for new connections
				s.logger.Debug().Msg("Connection change")
				connChange, s.connChangeNotify = context.WithCancel(context.Background()) // this ctx done, create a new one
				continue
			case <-ctx.Done():
				notConnected := make([]string, 0)
				for id, conn := range s.connections {
					if conn == nil {
						notConnected = append(notConnected, id)
					}
				}
				return fmt.Errorf("%w. not connected: %v", ctx.Err(), notConnected)
			}
		}

		s.logger.Debug().Msg("All connections established")
		return nil
	}
}

func (s *Server) disconnectAll(err error) {
	for id, conn := range s.connections {
		if conn != nil {
			_ = conn.Close(err)
		}

		s.Lock()
		s.connections[id] = nil
		s.Unlock()
	}
}

func (s *Server) clientConnected(id string, conn *common.Conn) error {
	s.Lock()
	defer s.Unlock()

	oldConn, ok := s.connections[id]
	if !ok {
		return fmt.Errorf("invalid id")
	}
	if oldConn != nil {
		return fmt.Errorf("id already connected")
		//s.logger.Warn().Str("id", id).Msg("id already connected")
	}

	s.connections[id] = conn
	if s.connChangeNotify != nil {
		s.connChangeNotify()
	}

	s.logger.Debug().Str("id", id).Msg("Client connected")

	return nil

}

func (s *Server) IsClientConnected(id string) bool {
	s.Lock()
	defer s.Unlock()

	oldConn, ok := s.connections[id]
	if !ok {
		return false
	}
	if oldConn != nil {
		return true
	}

	return false
}

func (s *Server) isAllConnected() bool {
	s.Lock()
	defer s.Unlock()
	for _, conn := range s.connections {
		if conn == nil {
			return false
		}
	}
	return true
}
