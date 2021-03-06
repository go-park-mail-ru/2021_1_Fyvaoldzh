package models

import "sync"

type RegisterHandler struct {
	Mu     *sync.Mutex
}