package api

type Manager interface {
	Reconcile() error
}
