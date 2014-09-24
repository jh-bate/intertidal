package ebb

type Notifier interface {
	Send(addresses []string, subject, content string) error
}
