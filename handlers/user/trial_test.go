package user_test

type MockMongo struct {
}

func (m *MockMongo) CreateClient(uri string) error {
	return nil
}
