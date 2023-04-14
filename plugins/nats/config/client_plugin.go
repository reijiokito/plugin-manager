package config

type PluginServer struct {
	clients map[string]string
}

func NewPluginServer() *PluginServer {
	return &PluginServer{
		clients: make(map[string]string),
	}
}

func (s *PluginServer) RegisterClient(name string, config string) {
	s.clients[name] = config
}

func (s *PluginServer) GetClient(name string) (string, bool) {
	client, ok := s.clients[name]
	return client, ok
}

func (s *PluginServer) GetPluginServer() map[string]string {
	return s.clients
}
