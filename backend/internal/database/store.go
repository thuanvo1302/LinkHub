package database

import "sync"

type Store struct {
	Mu sync.RWMutex

	nextUserID      int
	nextProfileID   int
	nextProfileLink int
	nextShortLinkID int

	UsersByID         map[string]map[string]any
	UsersByEmail      map[string]string
	ProfilesByUserID  map[string]map[string]any
	ProfilesByName    map[string]string
	ProfileLinksByID  map[string]map[string]any
	ProfileLinksOrder map[string][]string
	ShortLinksByID    map[string]map[string]any
	ShortLinksByCode  map[string]string
	ClickCounts       map[string]int
	RefreshTokens     map[string]map[string]any
	PlansByID         map[string]map[string]any
	SubscriptionsByID map[string]map[string]any
	PaymentsByID      map[string]map[string]any

	nextSubscriptionID int
	nextPaymentID      int
}

func NewStore() *Store {
	return &Store{
		UsersByID:         map[string]map[string]any{},
		UsersByEmail:      map[string]string{},
		ProfilesByUserID:  map[string]map[string]any{},
		ProfilesByName:    map[string]string{},
		ProfileLinksByID:  map[string]map[string]any{},
		ProfileLinksOrder: map[string][]string{},
		ShortLinksByID:    map[string]map[string]any{},
		ShortLinksByCode:  map[string]string{},
		ClickCounts:       map[string]int{},
		RefreshTokens:     map[string]map[string]any{},
		PlansByID:         map[string]map[string]any{},
		SubscriptionsByID: map[string]map[string]any{},
		PaymentsByID:      map[string]map[string]any{},
	}
}

func (s *Store) NextUserID() string {
	s.Mu.Lock()
	defer s.Mu.Unlock()
	s.nextUserID++
	return makeID("usr", s.nextUserID)
}

func (s *Store) NextProfileID() string {
	s.Mu.Lock()
	defer s.Mu.Unlock()
	s.nextProfileID++
	return makeID("pro", s.nextProfileID)
}

func (s *Store) NextProfileLinkID() string {
	s.Mu.Lock()
	defer s.Mu.Unlock()
	s.nextProfileLink++
	return makeID("lnk", s.nextProfileLink)
}

func (s *Store) NextShortLinkID() string {
	s.Mu.Lock()
	defer s.Mu.Unlock()
	s.nextShortLinkID++
	return makeID("shl", s.nextShortLinkID)
}

func (s *Store) NextSubscriptionID() string {
	s.Mu.Lock()
	defer s.Mu.Unlock()
	s.nextSubscriptionID++
	return makeID("sub", s.nextSubscriptionID)
}

func (s *Store) NextPaymentID() string {
	s.Mu.Lock()
	defer s.Mu.Unlock()
	s.nextPaymentID++
	return makeID("pay", s.nextPaymentID)
}

func makeID(prefix string, id int) string {
	return prefix + "_" + itoa(id)
}

func itoa(v int) string {
	if v == 0 {
		return "0"
	}

	buf := [20]byte{}
	i := len(buf)
	for v > 0 {
		i--
		buf[i] = byte('0' + v%10)
		v /= 10
	}
	return string(buf[i:])
}
